package SimpleRTree

import (
	"log"
	"math"
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

const MAX_POSSIBLE_SIZE = 9


type Options struct {
	MAX_ENTRIES int
}

type SimpleRTree struct {
	options  Options
	nodes []Node
	points FlatPoints
	built bool
	// Store pool of pools so that between algorithms it uses a channel (thread safe) within one algorithm it uses array
	queueItemPoolPool * searchQueueItemPoolPool
	queuePool * searchQueuePool
}
type Node struct {
	isLeaf     bool
	BBox       BBox
	firstChildIndex, endChildren int
}

// Structure used to constructing the ndoe
type nodeConstruct struct {
	height     int
	start, end int // index in the underlying array
}

// Create an RTree index from an array of points
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
	if options.MAX_ENTRIES > MAX_POSSIBLE_SIZE {
		panic(fmt.Sprintf("Cannot exceed %d for size", MAX_POSSIBLE_SIZE))
	}
	return r
}

func (r *SimpleRTree) Load(points FlatPoints) *SimpleRTree {
	return r.load(points, false)
}

func (r *SimpleRTree) LoadSortedArray(points FlatPoints) *SimpleRTree {
	return r.load(points, true)
}

func (r *SimpleRTree) FindNearestPointWithin(x, y, d float64) (x1, y1, d1 float64, found bool) {
	sqd := d * d // we work with squared distances
	return r.findNearestPointWithin(x, y, sqd)
}

func (r *SimpleRTree) FindNearestPoint (x, y float64) (x1, y1, d1 float64, found bool) {
	return r.findNearestPointWithin(x, y, math.Inf(1))
}
func (r *SimpleRTree) findNearestPointWithin (x, y, d float64) (x1, y1, d1 float64, found bool){
	var minItem *searchQueueItem
	distanceLowerBound := math.Inf(1)
	distanceUpperBound := d
	// if bbox is further from this bound then we don't explore it
	sq := r.queuePool.take()
	sq.Init()

	queueItemPool := r.queueItemPoolPool.take()
	rootNode := &r.nodes[0]
	mind, maxd := rootNode.computeDistances(x, y)
	if (maxd < distanceUpperBound) {
		distanceUpperBound = maxd
	}
	// Only start search if it is within bound
	if (mind < distanceUpperBound) {
		item := queueItemPool.take()
		item.node = rootNode
		item.distance = mind
		sq.Push(item)
	}

	for sq.Len() > 0 {
		item := sq.Pop()
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
			for nodeIndex := item.node.firstChildIndex; nodeIndex < item.node.endChildren; nodeIndex++ {
				n := &r.nodes[nodeIndex]
				mind, maxd := n.computeDistances(x, y)
				if (mind <= distanceUpperBound) {
					childItem := queueItemPool.take()
					childItem.node = n
					childItem.distance = mind
					sq.Push(childItem)
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
		item := sq.Pop()
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
	// Only do sqrt at the end
	d1 = math.Sqrt(distanceUpperBound)
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
	r.built = true

	rootNodeConstruct := r.build(points, isSorted)
	r.queueItemPoolPool = newSearchQueueItemPoolPool(2, rootNodeConstruct.height * r.options.MAX_ENTRIES)
	r.queuePool = newSearchQueuePool(2, rootNodeConstruct.height * r.options.MAX_ENTRIES)
	// Max proportion when not checking max distance 2.3111111111111113
	// Max proportion checking max distance 39 6 9 0.7222222222222222
	return r
}

func (r *SimpleRTree) build(points FlatPoints, isSorted bool) nodeConstruct {
	r.points = points
	r.nodes = make([]Node, 0, computeSize(points.Len()))
	r.nodes = append(r.nodes, Node{})
	rootNodeConstruct := nodeConstruct{
		height: int(math.Ceil(math.Log(float64(points.Len())) / math.Log(float64(r.options.MAX_ENTRIES)))),
		start: 0,
		end: points.Len(),
	}

	r.buildNodeDownwards(0, rootNodeConstruct, isSorted)
	r.computeBBoxDownwards(0)
	return rootNodeConstruct
}



func (r *SimpleRTree) buildNodeDownwards(nodeIndex int, nc nodeConstruct, isSorted bool) {
	n := &r.nodes[nodeIndex]
	N := nc.end - nc.start
	// target number of root entries to maximize storage utilization
	var M float64
	if N <= r.options.MAX_ENTRIES { // Leaf node
		r.setLeafNode(n, nc)
		return
	}

	M = math.Ceil(float64(N) / float64(math.Pow(float64(r.options.MAX_ENTRIES), float64(nc.height-1))))

	N2 := int(math.Ceil(float64(N) / M))
	N1 := N2 * int(math.Ceil(math.Sqrt(M)))

	// parent node might already be sorted. In that case we avoid double computation
	if (!isSorted) {
		sortX := xSorter{n: n, points: r.points, start: nc.start, end: nc.end, bucketSize:  N1}
		sortX.Sort()
	}
	nodeConstructs := [MAX_POSSIBLE_SIZE]nodeConstruct{}
	nodeConstructIndex := 0
	n.firstChildIndex = len(r.nodes)
	for i := 0; i < N; i += N1 {
		right2 := minInt(i+N1, N)
		sortY := ySorter{n: n, points: r.points, start: nc.start + i, end: nc.start + right2, bucketSize: N2}
		sortY.Sort()
		for j := i; j < right2; j += N2 {
			right3 := minInt(j+N2, right2)
			child := Node{
			}
			childC := nodeConstruct{
				start: nc.start + j,
				end: nc.start + right3,
				height: nc.height - 1,
			}
			r.nodes = append(r.nodes, child)
			nodeConstructs[nodeConstructIndex] = childC
			nodeConstructIndex++
		}
	}
	n.endChildren = len(r.nodes)
	// compute children
	for i:= 0; i < nodeConstructIndex; i++ {
		r.buildNodeDownwards(n.firstChildIndex + i, nodeConstructs[i], false)
	}
}



// Compute bbox of all tree all the way to the bottom
func (r *SimpleRTree) computeBBoxDownwards(nodeIndex int) BBox {
	n := &r.nodes[nodeIndex]
	var bbox BBox
	if n.isLeaf {
		bbox = n.BBox
	} else {
		bbox = r.computeBBoxDownwards(n.firstChildIndex)

		for i := n.firstChildIndex + 1; i < n.endChildren; i++ {
			bbox = bbox.extend(r.computeBBoxDownwards(i))
		}
	}
	n.BBox = bbox
	return bbox
}


func (r *SimpleRTree) setLeafNode(n * Node, nc nodeConstruct) {
	// Here we follow original rbush implementation.
	n.firstChildIndex = len(r.nodes)
	for i := 0; i < nc.end - nc.start; i++ {
		x1, y1 := r.points.GetPointAt(nc.start + i)
		child := Node{
			isLeaf: true,
			BBox: BBox{
				MinX: x1,
				MaxX: x1,
				MinY: y1,
				MaxY: y1,
			},
		}
		// Note this is not thread safe. At the moment we are doing it in one goroutine so we are safe
		r.nodes = append(r.nodes, child)
	}
	n.endChildren = len(r.nodes)
}

func (r *SimpleRTree) toJSON () {
	text := make([]string, 0)
	fmt.Println(strings.Join(r.toJSONAcc(0, text), ","))
}

func (r *SimpleRTree) toJSONAcc (nodeIndex int, text []string) []string {
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
	n := &r.nodes[nodeIndex]
	if err := t.Execute(&tpl, n); err != nil {
		log.Fatal(err)
	}
	text = append(text, tpl.String())
	for i := n.firstChildIndex; i < n.endChildren; i++ {
		text = r.toJSONAcc(i, text)
	}
	return text
}

func (n * Node) computeDistances (x, y float64) (mind, maxd float64) {
	// TODO try simd
	if (n.isLeaf) {
	       // node is point, there is only one distance
	       d := (x - n.BBox.MinX) * (x - n.BBox.MinX)  + (y - n.BBox.MinY) * (y - n.BBox.MinY)
	       return d, d
	}
	minx, maxx := sortFloats((x - n.BBox.MinX) * (x - n.BBox.MinX), (x - n.BBox.MaxX) * (x - n.BBox.MaxX))
	miny, maxy := sortFloats((y - n.BBox.MinY) * (y - n.BBox.MinY), (y - n.BBox.MaxY) * (y - n.BBox.MaxY))

	sideX := (n.BBox.MaxX - n.BBox.MinX) * (n.BBox.MaxX - n.BBox.MinX)
	sideY := (n.BBox.MaxY - n.BBox.MinY) * (n.BBox.MaxY - n.BBox.MinY)

	// Code is a bit cryptic but it is equivalent to the commented code which is clearer
	if (maxx >= sideX) {
		mind += minx
	}
	if (maxy >= sideY) {
		mind += miny
	}
	/*
	// point is inside because max distances in both axis are smaller than sides of the square
	if (maxx < sideX && maxy < sideY) {
		// do nothing mind is already 0
	} else if (maxx < sideX && maxy >= sideY) {
		// point is in vertical stripe. Hence distance to the bbox is minimum vertical distance
		mind = miny
	} else if (maxx >= sideX && maxy < sideY) {
		// point is in horizontal stripe, Hence distance is least distance to one of the sides (vertical distance is 0
		mind = minx
	} else if (maxx >= sideX && maxy >= sideY){
		// point is not inside bbox. closest vertex is that one with closest x and y
		mind = minx + miny
	}*/
	maxd = maxx + maxy
	return
}

func minInt(a, b int) int {
       if a < b {
	       return a
       }
       return b
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

func sortFloats (x1, x2 float64) (x3, x4 float64) {
	if (x1 > x2) {
		return x2, x1
	}
	return x1, x2
}

func computeSize (n int) (size int) {
	return 2 * n
}