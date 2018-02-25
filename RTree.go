package SimpleRTree

import (
	"log"
	"math"
	"container/heap"
)

const (
	MAX_HEIGHT_TO_SPLIT = 3 // When creating the index we'll split the task into a new goroutine until we reach this height
)


type Interface interface {
	GetPointAt(i int) (x1, y1 float64)        // Retrieve point at position i
	Len() int                                 // Number of elements
	Swap(i, j int)                            // Swap elements with indexes i and j
}


type Options struct {
	MAX_ENTRIES int
}

type SimpleRTree struct {
	options  Options
	rootNode *Node
	points Interface
	built bool
}
type Node struct {
	children   []*Node
	height     int
	isLeaf     bool
	start, end int // index in the underlying array
	parentNode *Node
	BBox       BBox
}

// Create an RBush index from an array of points
func New() *SimpleRTree {
	defaultOptions := Options{
		MAX_ENTRIES: 9,
	}
	return NewWithOptions(defaultOptions)
}

func NewWithOptions(options Options) *SimpleRTree {
	r := &SimpleRTree{
		options: options,
	}
	return r
}

func (r *SimpleRTree) Load(points Interface) *SimpleRTree {
	return r.load(points, false)
}

func (r *SimpleRTree) LoadSortedArray(points Interface) *SimpleRTree {
	return r.load(points, true)
}

func (r *SimpleRTree) FindNearestPoint (x, y float64) (x1, y1 float64){
	var minItem *searchQueueItem
	distanceLowerBound := math.Inf(1)
	// if bbox is further from this bound then we don't explore it
	sq := make(searchQueue, r.rootNode.height * r.options.MAX_ENTRIES)
	heap.Init(&sq)

	mind, maxd := r.rootNode.computeDistances(x, y)
	distanceUpperBound := maxd
	heap.Push(&sq, searchQueueItem{node: r.rootNode, distance: mind})

	for sq.Len() > 0 {
		item := heap.Pop(&sq).(*searchQueueItem)
		currentDistance := item.distance
		if (minItem != nil && currentDistance > distanceLowerBound) {
			break
		}

		if (item.node.isLeaf) {
			// we know it is smaller from the previous test
			distanceLowerBound = currentDistance
			minItem = item
		} else {
			for _, n := range(item.node.children) {
				mind, maxd := n.computeDistances(x, y)
				if (mind < distanceUpperBound) {
					heap.Push(&sq, searchQueueItem{node: n, distance: mind})
				}
				// Distance to one of the corners is lower than the upper bound
				// so there must be a point at most within distanceUpperBound
				if (maxd < distanceUpperBound) {
					distanceUpperBound = maxd
				}
			}
		}
	}
	x1 = minItem.node.BBox.MaxX
	y1 = minItem.node.BBox.MaxY
	return
}

func (r *SimpleRTree) load (points Interface, isSorted bool) *SimpleRTree {
	if points.Len() == 0 {
		return r
	}
	if r.built {
		log.Fatal("Tree is static, cannot load twice")
	}

	node := r.build(points, isSorted)
	r.rootNode = node
	return r
}

func (r *SimpleRTree) build(points Interface, isSorted bool) *Node {

	confirmCh := make(chan int, 1)

	rootNode := &Node{
		height: int(math.Ceil(math.Log(float64(points.Len())) / math.Log(float64(r.options.MAX_ENTRIES)))),
		start: 0,
		end: points.Len(),
	}
	remainingNodes := 1

	go r.buildNodeDownwards(rootNode, confirmCh, true, isSorted)
	for remainingNodes > 0 {
		i := <-confirmCh
		remainingNodes += i
	}
	close(confirmCh)
	rootNode.computeBBoxDownwards()
	return rootNode
}



func (r *SimpleRTree) buildNodeDownwards(n *Node, confirmCh chan int, isCalledAsync, isSorted bool) {
	if isCalledAsync {
		defer func() {
			confirmCh <- -1
		}()
	}

	N := n.end - n.start
	// target number of root entries to maximize storage utilization
	var M float64
	if N <= r.options.MAX_ENTRIES { // Leaf node
		r.setLeafNode(n)
		return
	}

	M = math.Ceil(float64(N) / float64(math.Pow(float64(r.options.MAX_ENTRIES), float64(n.height-1))))

	N2 := int(math.Ceil(float64(N) / M))
	N1 := N2 * int(math.Ceil(math.Sqrt(M)))

	// parent node might already be sorted. In that case we avoid double computation
	if (n.parentNode != nil || !isSorted) {
		sortX := xSorter{n: n, points: r.points, start: n.start, end: n.end, bucketSize:  N1}
		sortX.Sort()
	}
	// runtime.Breakpoint()
	for i := 0; i < N; i += N1 {
		right2 := minInt(i+N1, N)
		sortY := ySorter{n: n, points: r.points, start: i, end: right2, bucketSize: N2}
		sortY.Sort()
		for j := i; j < right2; j += N2 {
			right3 := minInt(j+N2, right2)
			child := Node{
				start: n.start + j,
				end: n.start + right3,
				height:     n.height - 1,
				parentNode: n,
			}
			n.children = append(n.children, &child)
			// remove reference to interface, we only need it for points

		}
	}
	// compute children
	for _, c := range n.children {
		// Only launch a goroutine for big height. we don't want a goroutine to sort 4 points
		if n.height > MAX_HEIGHT_TO_SPLIT {
			confirmCh <- 1
			go r.buildNodeDownwards(c, confirmCh, true, false)
		} else {
			r.buildNodeDownwards(c, confirmCh, false, false)
		}
	}
}



// Compute bbox of all tree all the way to the bottom
func (n *Node) computeBBoxDownwards() BBox {

	var bbox BBox
	if n.isLeaf {
		bbox = n.BBox
	} else {
		bbox = n.children[0].computeBBoxDownwards()

		for i := 1; i < len(n.children); i++ {
			bbox = bbox.extend(n.children[i].computeBBoxDownwards())
		}
	}
	n.BBox = bbox
	return bbox
}


func (r *SimpleRTree) setLeafNode(n * Node) {
	// Here we follow original rbush implementation.
 	children := make([]*Node, n.end - n.start)
 	n.children = children
	n.height = 1

	for i := 0; i < n.end - n.start; i++ {
		x1, y1 := r.points.GetPointAt(i)
		children[i] = &Node{
			start: i,
			end: i +1,
			isLeaf: true,
			BBox: BBox{
				MinX: x1,
				MaxX: x1,
				MinY: y1,
				MaxY: y1,
			},
			parentNode: n,
		}
	}
}

func (n * Node) computeDistances (x, y float64) (mind, maxd float64) {
	// TODO try reuse array
	// TODO try simd
	if (n.isLeaf) {
		// node is point, there is only one distance
		d := math.Pow(x - n.BBox.MinX, 2)  + math.Pow(y - n.BBox.MinY, 2)
		return d, d
	}

	distances := [4]float64{
		math.Pow(x - n.BBox.MinX, 2) + math.Pow(y - n.BBox.MinY, 2),
		math.Pow(x - n.BBox.MinX, 2) + math.Pow(y - n.BBox.MaxY, 2),
		math.Pow(x - n.BBox.MaxX, 2) + math.Pow(y - n.BBox.MinY, 2),
		math.Pow(x - n.BBox.MaxX, 2) + math.Pow(y - n.BBox.MaxY, 2),
	}
	mind, maxd = minmaxFloatArray(distances)

	if (n.BBox.containsPoint(x, y)) {
		mind = 0
	}
	return
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func minmaxFloatArray (s [4]float64) (min, max float64) {
	// TODO try min of four
	min = s[0]
	max = s[0]
	for _, e := range s {
		if e < min {
			min = e
		}
		if e > min {
			max = e
		}
	}
	return min, max
}
