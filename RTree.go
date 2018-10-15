package SimpleRTree

import (
	"bytes"
	"fmt"
	"log"
	"math"
	"strings"
	"text/template"
	"unsafe"
	"sync"
)

type Interface interface {
	GetPointAt(i int) (x1, y1 float64) // Retrieve point at position i
	Len() int                          // Number of elements
	Swap(i, j int)                     // Swap elements with indexes i and j
}

const (
	DEFAULT = iota
	LEAF
	PRELEAF
)

const MAX_POSSIBLE_SIZE = 9

type nodeType int8

type Options struct {
	MAX_ENTRIES int
}

var NODE_SIZE = unsafe.Sizeof(Node{})

type SimpleRTree struct {
	options Options
	nodes   []Node
	points  FlatPoints
	built   bool
	queuePool         sync.Pool
	sorterBuffer      []int // floyd rivest requires a bucket, we allocate it once and reuse
}
type Node struct {
	nodeType   nodeType
	nChildren  int8
	firstChild *Node
	BBox       VectorBBox
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

func (r *SimpleRTree) FindNearestPoint(x, y float64) (x1, y1, d1 float64, found bool) {
	return r.FindNearestPointWithin(x, y, math.Inf(1))
}
func (r *SimpleRTree) FindNearestPointWithin(x, y, dsquared float64) (x1, y1, d1squared float64, found bool) {
	var minItem searchQueueItem
	distanceLowerBound := math.Inf(1)
	// if bbox is further from this bound then we don't explore it
	distanceUpperBound := dsquared
	sq := r.queuePool.Get().(searchQueue)
	sq = sq[0:0]
	sq.Init()

	rootNode := &r.nodes[0]
	mind, maxd := vectorComputeDistances(rootNode.BBox, x, y)
	if maxd < distanceUpperBound {
		distanceUpperBound = maxd
	}
	// Only start search if it is within bound
	if mind < distanceUpperBound {
		sq.Push(searchQueueItem{node: rootNode, distance: mind})
	}

	for sq.Len() > 0 {
		item := sq.Pop()
		currentDistance := item.distance
		if found && currentDistance > distanceLowerBound {
			break
		}

		switch item.node.nodeType {
		case LEAF:
			// we know it is smaller from the previous test
			distanceLowerBound = currentDistance
			minItem = item
			found = true
		case PRELEAF:
			f := unsafe.Pointer(item.node.firstChild)
			var i int8
			for i = 0; i < item.node.nChildren; i++ {
				n := (*Node)(f)
				d := n.computeLeafDistance(x, y)
				if d <= distanceUpperBound {
					sq.Push(searchQueueItem{node: n, distance: d})
					distanceUpperBound = d
				}
				f = unsafe.Pointer(uintptr(f) + NODE_SIZE)
			}
		default:
			f := unsafe.Pointer(item.node.firstChild)
			var i int8
			for i = 0; i < item.node.nChildren; i++ {
				n := (*Node)(f)
				mind, maxd := vectorComputeDistances(n.BBox, x, y)
				if mind <= distanceUpperBound {
					sq.Push(searchQueueItem{node: n, distance: mind})
				}
				// Distance to one of the corners is lower than the upper bound
				// so there must be a point at most within distanceUpperBound
				if maxd < distanceUpperBound {
					distanceUpperBound = maxd
				}
				f = unsafe.Pointer(uintptr(f) + NODE_SIZE)
			}
		}
	}

	// return heap
	r.queuePool.Put(sq)

	if !found {
		return
	}
	x1 = minItem.node.BBox[VECTOR_BBOX_MAX_X]
	y1 = minItem.node.BBox[VECTOR_BBOX_MAX_Y]
	d1squared = distanceUpperBound
	return
}

func (r *SimpleRTree) load(points FlatPoints, isSorted bool) *SimpleRTree {
	if points.Len() == 0 {
		return r
	}
	if r.built {
		log.Fatal("Tree is static, cannot load twice")
	}
	r.built = true

	r.sorterBuffer = make([]int, 0, r.options.MAX_ENTRIES+1)
	rootNodeConstruct := r.build(points, isSorted)
	r.queuePool = sync.Pool{
		New: func () interface {} {
			return make(searchQueue, rootNodeConstruct.height*r.options.MAX_ENTRIES)
		},
	}
	firstQueue := r.queuePool.Get()
	r.queuePool.Put(firstQueue)
	return r
}

func (r *SimpleRTree) build(points FlatPoints, isSorted bool) nodeConstruct {
	r.points = points
	r.nodes = make([]Node, 0, computeSize(points.Len()))
	r.nodes = append(r.nodes, Node{})
	rootNodeConstruct := nodeConstruct{
		height: int(math.Ceil(math.Log(float64(points.Len())) / math.Log(float64(r.options.MAX_ENTRIES)))),
		start:  0,
		end:    points.Len(),
	}

	r.buildNodeDownwards(&r.nodes[0], rootNodeConstruct, isSorted)
	return rootNodeConstruct
}

func (r *SimpleRTree) buildNodeDownwards(n *Node, nc nodeConstruct, isSorted bool) *VectorBBox {
	N := nc.end - nc.start
	// target number of root entries to maximize storage utilization
	var M float64
	if N <= r.options.MAX_ENTRIES { // Leaf node
		return r.setLeafNode(n, nc)
	}

	M = math.Ceil(float64(N) / float64(math.Pow(float64(r.options.MAX_ENTRIES), float64(nc.height-1))))

	N2 := int(math.Ceil(float64(N) / M))
	N1 := N2 * int(math.Ceil(math.Sqrt(M)))

	// parent node might already be sorted. In that case we avoid double computation
	if !isSorted {
		sortX := xSorter{n: n, points: r.points, start: nc.start, end: nc.end, bucketSize: N1}
		sortX.Sort(r.sorterBuffer)
	}
	nodeConstructs := [MAX_POSSIBLE_SIZE]nodeConstruct{}
	var nodeConstructIndex int8
	firstChildIndex := len(r.nodes)
	for i := 0; i < N; i += N1 {
		right2 := minInt(i+N1, N)
		sortY := ySorter{n: n, points: r.points, start: nc.start + i, end: nc.start + right2, bucketSize: N2}
		sortY.Sort(r.sorterBuffer)
		for j := i; j < right2; j += N2 {
			right3 := minInt(j+N2, right2)
			child := Node{}
			childC := nodeConstruct{
				start:  nc.start + j,
				end:    nc.start + right3,
				height: nc.height - 1,
			}
			r.nodes = append(r.nodes, child)
			nodeConstructs[nodeConstructIndex] = childC
			nodeConstructIndex++
		}
	}
	n.firstChild = &r.nodes[firstChildIndex]
	n.nChildren = nodeConstructIndex
	// compute children
	var i int8
	bbox := r.buildNodeDownwards(&r.nodes[firstChildIndex], nodeConstructs[i], false)
	for i = 1; i < nodeConstructIndex; i++ {
		// TODO check why using (*Node)f here does not work
		bbox2 := r.buildNodeDownwards(&r.nodes[firstChildIndex+int(i)], nodeConstructs[i], false)
		vectorBBoxExtend(bbox, bbox2)
	}
	n.BBox = *bbox
	return bbox
}

func (r *SimpleRTree) setLeafNode(n *Node, nc nodeConstruct) *VectorBBox {
	// Here we follow original rbush implementation.
	firstChildIndex := len(r.nodes)

	x0, y0 := r.points.GetPointAt(nc.start)
	vb := newVectorBBox(x0, y0, x0, y0)
	bbox := &vb
	child0 := Node{
		nodeType: LEAF,
		BBox: [4]float64{
			x0,
			y0,
			x0,
			y0,
		},
	}
	r.nodes = append(r.nodes, child0)

	for i := 1; i < nc.end-nc.start; i++ {
		x1, y1 := r.points.GetPointAt(nc.start + i)
		child := Node{
			nodeType: LEAF,
			BBox: [4]float64{
				x1,
				y1,
				x1,
				y1,
			},
		}
		vectorBBoxExtend(bbox, &child.BBox)
		// Note this is not thread safe. At the moment we are doing it in one goroutine so we are safe
		r.nodes = append(r.nodes, child)
	}
	n.firstChild = &r.nodes[firstChildIndex]
	n.nChildren = int8(nc.end - nc.start)
	n.nodeType = PRELEAF
	n.BBox = vb
	return bbox
}

func (r *SimpleRTree) toJSON() {
	text := make([]string, 0)
	fmt.Println(strings.Join(r.toJSONAcc(&r.nodes[0], text), ","))
}

func (r *SimpleRTree) toJSONAcc(n *Node, text []string) []string {
	t, err := template.New("foo").Parse(`{
	       "type": "Feature",
	       "properties": {},
	       "geometry": {
       "type": "Polygon",
       "coordinates": [
       [
       [
       {{index .BBox 0}},
       {{index .BBox 1}}
       ],
       [
       {{index .BBox 2}},
       {{index .BBox 1}}
       ],
       [
       {{index .BBox 2}},
       {{index .BBox 3}}
       ],
       [
       {{index .BBox 0}},
       {{index .BBox 3}}
       ],
       [
       {{index .BBox 0}},
       {{index .BBox 1}}
       ]
       ]
       ]
       }
       }`)
	if err != nil {
		log.Fatal(err)
	}
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, n); err != nil {
		log.Fatal(err)
	}
	text = append(text, tpl.String())
	f := unsafe.Pointer(n.firstChild)
	var i int8
	for i = 0; i < n.nChildren; i++ {
		cn := (*Node)(f)
		text = r.toJSONAcc(cn, text)
		f = unsafe.Pointer(uintptr(f) + NODE_SIZE)
	}
	return text
}

// node is point, there is only one distance
func (n *Node) computeLeafDistance(x, y float64) float64 {
	return (x-n.BBox[VECTOR_BBOX_MIN_X])*(x-n.BBox[VECTOR_BBOX_MIN_X]) +
		(y-n.BBox[VECTOR_BBOX_MIN_Y])*(y-n.BBox[VECTOR_BBOX_MIN_Y])
}
func computeDistances(bbox VectorBBox, x, y float64) (mind, maxd float64) {
	// TODO try simd
	minX := bbox[0]
	minY := bbox[1]
	maxX := bbox[2]
	maxY := bbox[3]
	minx, maxx := sortFloats((x-minX)*(x-minX), (x-maxX)*(x-maxX))
	miny, maxy := sortFloats((y-minY)*(y-minY), (y-maxY)*(y-maxY))

	sideX := (maxX - minX) * (maxX - minX)
	sideY := (maxY - minY) * (maxY - minY)

	// Code is a bit cryptic but it is equivalent to the commented code which is clearer
	if maxx >= sideX {
		mind += minx
	}
	if maxy >= sideY {
		mind += miny
	}

	// Given a bbox the distance will be bounded to the two intermediate corners
	maxd = minFloat(maxx+miny, minx+maxy)
	return
	/*
		How to compute mind
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
}
func vectorComputeDistances(bbox VectorBBox, x, y float64) (mind, maxd float64)

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func maxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

type FlatPoints []float64

func (fp FlatPoints) Len() int {
	return len(fp) / 2
}

func (fp FlatPoints) Swap(i, j int) {
	fp[2*i], fp[2*i+1], fp[2*j], fp[2*j+1] = fp[2*j], fp[2*j+1], fp[2*i], fp[2*i+1]
}

func (fp FlatPoints) GetPointAt(i int) (x1, y1 float64) {
	return fp[2*i], fp[2*i+1]
}

func sortFloats(x1, x2 float64) (x3, x4 float64) {
	if x1 > x2 {
		return x2, x1
	}
	return x1, x2
}

func computeSize(n int) (size int) {
	return 2 * n
}
