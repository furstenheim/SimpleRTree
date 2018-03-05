package SimpleRTree

import (
	"log"
	"math"
	"container/heap"
	"text/template"
	"bytes"
	"fmt"
	"strings"
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
	points FlatPoints
	built bool
	// Store pool of pools so that between algorithms it uses a channel (thread safe) within one algorithm it uses array
	queueItemPoolPool * searchQueueItemPoolPool
	queuePool * searchQueuePool
}
type Node struct {
	children   []*Node
	height     int
	isLeaf     bool
	start, end int // index in the underlying array
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

func (r *SimpleRTree) Load(points FlatPoints) *SimpleRTree {
	return r.load(points, false)
}

func (r *SimpleRTree) LoadSortedArray(points FlatPoints) *SimpleRTree {
	return r.load(points, true)
}

func (r *SimpleRTree) FindNearestPoint (x, y float64) (x1, y1, d1 float64, found bool){
	var minItem *searchQueueItem
	distanceLowerBound := math.Inf(1)
	// if bbox is further from this bound then we don't explore it
	sq := r.queuePool.take()
	heap.Init(sq)

	queueItemPool := r.queueItemPoolPool.take()
	mind, maxd := r.rootNode.computeDistances(x, y)
	distanceUpperBound := maxd
	item := queueItemPool.take()
	item.node = r.rootNode
	item.distance = mind
	heap.Push(sq, item)

	for sq.Len() > 0 {
		item := heap.Pop(sq).(*searchQueueItem)
		currentDistance := item.distance
		if (minItem != nil && currentDistance > distanceLowerBound) {
			queueItemPool.giveBack(item);
			break
		}

		if (item.node.isLeaf) {
			// we know it is smaller from the previous test
			distanceLowerBound = currentDistance
			minItem = item
		} else {
			for _, n := range(item.node.children) {
				mind, maxd := n.computeDistances(x, y)
				if (mind <= distanceUpperBound) {
					childItem := queueItemPool.take()
					childItem.node = n
					childItem.distance = mind
					heap.Push(sq, childItem)
				}
				// Distance to one of the corners is lower than the upper bound
				// so there must be a point at most within distanceUpperBound
				if (maxd < distanceUpperBound) {
					distanceUpperBound = maxd
				}
			}
		}
		queueItemPool.giveBack(item)
	}

	// Return all missing items. This could probably be async
	for sq.Len() > 0 {
		item := heap.Pop(sq).(*searchQueueItem)
		queueItemPool.giveBack(item)
	}

	// return pool of items
	r.queueItemPoolPool.giveBack(queueItemPool)
	r.queuePool.giveBack(sq)

	if (minItem == nil) {
		return
	}
	x1 = minItem.node.BBox.MaxX
	y1 = minItem.node.BBox.MaxY
	d1 = distanceUpperBound
	found = true
	return
}

func (r *SimpleRTree) load (points FlatPoints, isSorted bool) *SimpleRTree {
	if points.Len() == 0 {
		return r
	}
	if r.built {
		log.Fatal("Tree is static, cannot load twice")
	}

	node := r.build(points, isSorted)
	r.rootNode = node
	r.queueItemPoolPool = newSearchQueueItemPoolPool(2, r.rootNode.height * r.options.MAX_ENTRIES)
	r.queuePool = newSearchQueuePool(2, r.rootNode.height * r.options.MAX_ENTRIES)
	// Max proportion when not checking max distance 2.3111111111111113
	// Max proportion checking max distance 39 6 9 0.7222222222222222
	return r
}

func (r *SimpleRTree) build(points FlatPoints, isSorted bool) *Node {

	r.points = points
	rootNode := &Node{
		height: int(math.Ceil(math.Log(float64(points.Len())) / math.Log(float64(r.options.MAX_ENTRIES)))),
		start: 0,
		end: points.Len(),
	}

	r.buildNodeDownwards(rootNode, true, isSorted)
	rootNode.computeBBoxDownwards()
	return rootNode
}



func (r *SimpleRTree) buildNodeDownwards(n *Node, isCalledAsync, isSorted bool) {
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
	if (!isSorted) {
		sortX := xSorter{n: n, points: r.points, start: n.start, end: n.end, bucketSize:  N1}
		sortX.Sort()
	}
	for i := 0; i < N; i += N1 {
		right2 := minInt(i+N1, N)
		sortY := ySorter{n: n, points: r.points, start: n.start + i, end: n.start + right2, bucketSize: N2}
		sortY.Sort()
		for j := i; j < right2; j += N2 {
			right3 := minInt(j+N2, right2)
			child := Node{
				start: n.start + j,
				end: n.start + right3,
				height:     n.height - 1,
			}
			n.children = append(n.children, &child)
			// remove reference to interface, we only need it for points

		}
	}
	// compute children
	for _, c := range n.children {
		// Only launch a goroutine for big height. we don't want a goroutine to sort 4 points
		r.buildNodeDownwards(c, true, false)
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
		x1, y1 := r.points.GetPointAt(n.start + i)
		children[i] = &Node{
			start: n.start + i,
			end: n.start + i +1,
			isLeaf: true,
			BBox: BBox{
				MinX: x1,
				MaxX: x1,
				MinY: y1,
				MaxY: y1,
			},
		}
	}
}

func (r *SimpleRTree) toJSON () {
	text := make([]string, 0)
	fmt.Println(strings.Join(r.rootNode.toJSON(text), ","))
}

func (n *Node) toJSON (text []string) []string {
	t, err := template.New("foo").Parse(`{
	       "type": "Feature",
	       "properties": {},
	       "geometry": {
       "type": "Polygon",
       "coordinates": [
       [
       [
       {{.BBox.MinX}},
       {{.BBox.MinY}}
       ],
       [
       {{.BBox.MaxX}},
       {{.BBox.MinY}}
       ],
       [
       {{.BBox.MaxX}},
       {{.BBox.MaxY}}
       ],
       [
       {{.BBox.MinX}},
       {{.BBox.MaxY}}
       ],
       [
       {{.BBox.MinX}},
       {{.BBox.MinY}}
       ]
       ]
       ]
       }
       }`)
	if (err != nil) {
		log.Fatal(err)
	}
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, n); err != nil {
		log.Fatal(err)
	}
	text = append(text, tpl.String())
	for _, c := range(n.children) {
		text = c.toJSON(text)
	}
	return text
}

func (n * Node) computeDistances (x, y float64) (mind, maxd float64) {
       // TODO try reuse array
       // TODO try simd
       if (n.isLeaf) {
	       // node is point, there is only one distance
	       d := (x - n.BBox.MinX) * (x - n.BBox.MinX)  + (y - n.BBox.MinY) * (y - n.BBox.MinY)
	       return d, d
       }

       distances := [4]float64{
	       (x - n.BBox.MinX) * (x - n.BBox.MinX) + (y - n.BBox.MinY) * (y - n.BBox.MinY),
	       (x - n.BBox.MinX) * (x - n.BBox.MinX) + (y - n.BBox.MaxY) * (y - n.BBox.MaxY),
	       (x - n.BBox.MaxX) * (x - n.BBox.MaxX) + (y - n.BBox.MinY) * (y - n.BBox.MinY),
	       (x - n.BBox.MaxX) * (x - n.BBox.MaxX) + (y - n.BBox.MaxY) * (y - n.BBox.MaxY),
       }
       mind, maxd = minmaxFloatArray(distances)

       // Min distance is vertical line
       if (n.BBox.MinX <= x && x <= n.BBox.MaxX) {
	       mind = math.Min((n.BBox.MaxY - y) * (n.BBox.MaxY - y), (n.BBox.MinY - y) * (n.BBox.MinY - y))
       }
       if (n.BBox.MinY <= y && y <= n.BBox.MaxY) {
	       mind = math.Min((n.BBox.MaxX - x) * (n.BBox.MaxX - x), (n.BBox.MinX - x) * (n.BBox.MinX - x))
       }
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

func maxInt(a, b int) int {
       if a > b {
	       return a
       }
       return b
}

func minmaxFloatArray (s [4]float64) (min, max float64) {
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


type FlatPoints []float64

func (fp FlatPoints) Len () int {
	return len(fp) / 2
}

func (fp FlatPoints) Swap (i, j int) {
	fp[2 * i], fp[2 * i + 1], fp[2 * j], fp[2 * j + 1] = fp[2 * j], fp[2 * j + 1], fp[2 * i], fp[2 * i + 1]
}

func (fp FlatPoints) GetPointAt(i int) (x1, y1 float64) {
	return fp[2 * i], fp[2 * i +1]
}