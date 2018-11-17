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


type SimpleRTree32 struct {
	options Options
	nodes   []Node32
	points  FlatPoints32
	built   bool
	queuePool         sync.Pool
	unsafeQueue         searchQueue32 // Only used in unsafe mode
	sorterBuffer      []int // floyd rivest requires a bucket, we allocate it once and reuse
}
type Node32 struct {
	nodeType   nodeType
	nChildren  int8
	// Here we save firstChild - firstNode. That means that there is there is a theoretical upper limit to the tree of
	// maxuint32 / node_size = 4294967295 / 40 = 107374182 ~ 100M
	firstChildOffset uint32
	BBox       VectorBBox32
}

var NODE32_SIZE = unsafe.Sizeof(Node32{})
var FLAT32_POINT_SIZE =unsafe.Sizeof([2]float32{})
var FLOAT32_SIZE = uintptr(unsafe.Sizeof([1]float32{}))

type pooledMem32 struct {
	sorterBuffer []int
	sq searchQueue32
	nodes []Node32
}

// Create an RTree index from an array of points
func New32() *SimpleRTree32 {
	defaultOptions := Options{
		MAX_ENTRIES: MAX_POSSIBLE_SIZE,
	}
	return New32WithOptions(defaultOptions)
}

func New32WithOptions(options Options) *SimpleRTree32 {
	r := &SimpleRTree32{
		options: options,
	}
	if options.MAX_ENTRIES > MAX_POSSIBLE_SIZE {
		panic(fmt.Sprintf("Cannot exceed %d for size", MAX_POSSIBLE_SIZE))
	}
	if options.MAX_ENTRIES == 0 {
		r.options.MAX_ENTRIES = MAX_POSSIBLE_SIZE
	}
	return r
}

// Free up resources in case they were lent
func (r *SimpleRTree32) Destroy () {
	if r.options.RTreePool != nil {
		r.options.RTreePool.Put(
			&pooledMem32{
				sorterBuffer: r.sorterBuffer,
				sq: r.unsafeQueue,
				nodes: r.nodes,
			},
		)
	}
}

func (r *SimpleRTree32) Load(points FlatPoints32) *SimpleRTree32 {
	return r.load(points, false)
}

func (r *SimpleRTree32) LoadSortedArray(points FlatPoints32) *SimpleRTree32 {
	return r.load(points, true)
}

func (r *SimpleRTree32) FindNearestPoint(x, y float32) (x1, y1, d1 float32, found bool) {
	return r.FindNearestPointWithin(x, y, float32(math.Inf(1)))
}
func (r *SimpleRTree32) FindNearestPointWithin(x, y, dsquared float32) (x1, y1, d1squared float32, found bool) {
	var minItem searchQueueItem32
	distanceLowerBound := float32(math.Inf(1))
	// if bbox is further from this bound then we don't explore it
	distanceUpperBound := dsquared
	var sq searchQueue32
	if r.options.UnsafeConcurrencyMode {
		sq = r.unsafeQueue
	} else {
		sq = r.queuePool.Get().(searchQueue32)
	}
	sq = sq[0:0]
	sq.Init()

	rootNode := &r.nodes[0]
	unsafeRootLeafNode := uintptr(unsafe.Pointer(&r.points[0]))
	unsafeRootNode := uintptr(unsafe.Pointer(rootNode))
	mind, maxd := vectorComputeDistances32(rootNode.BBox, x, y)
	if maxd < distanceUpperBound {
		distanceUpperBound = maxd
	}
	// Only start search if it is within bound
	if mind < distanceUpperBound {
		sq.Push(searchQueueItem32{node: rootNode, distance: mind})
	}

	for sq.Len() > 0 {
		item := sq.Pop()
		currentDistance := item.distance
		if found && currentDistance > distanceLowerBound {
			break
		}

		if item.node == nil { // Leaf
			// we know it is smaller from the previous test
			distanceLowerBound = currentDistance
			minItem = item
			found = true
			continue
		}
		switch item.node.nodeType {
		case PRELEAF:
			f := unsafe.Pointer(unsafeRootLeafNode + uintptr(item.node.firstChildOffset))
			var i int8
			for i = item.node.nChildren; i>0; i-- {
				px := *(*float32)(f)
				f = unsafe.Pointer(uintptr(f) + FLOAT32_SIZE)
				py := *(*float32)(f)

				d := computeLeafDistance32(px, py, x, y)
				if d <= distanceUpperBound {
					sq.Push(searchQueueItem32{node: nil, px: px, py: py, distance: d})
					distanceUpperBound = d
				}
				f = unsafe.Pointer(uintptr(f) + FLOAT32_SIZE)
			}
		default:
			f := unsafe.Pointer(unsafeRootNode + uintptr(item.node.firstChildOffset))
			var i int8
			for i = item.node.nChildren; i>0; i-- {
				n := (*Node32)(f)
				mind, maxd := vectorComputeDistances32(n.BBox, x, y)
				if mind <= distanceUpperBound {
					sq.Push(searchQueueItem32{node: n, distance: mind})
				}
				// Distance to one of the corners is lower than the upper bound
				// so there must be a point at most within distanceUpperBound
				if maxd < distanceUpperBound {
					distanceUpperBound = maxd
				}
				f = unsafe.Pointer(uintptr(f) + NODE32_SIZE)
			}
		}
	}

	// return heap
	if !r.options.UnsafeConcurrencyMode {
		r.queuePool.Put(sq)
	}

	if !found {
		return
	}
	x1 = minItem.px
	y1 = minItem.py
	d1squared = distanceUpperBound
	return
}

func (r *SimpleRTree32) load(points FlatPoints32, isSorted bool) *SimpleRTree32 {
	if points.Len() == 0 {
		return r
	}
	if points.Len() >= math.MaxUint32 / int(NODE32_SIZE) {
		log.Fatal("Exceded maximum possible size", math.MaxUint32 / int(NODE32_SIZE))
	}
	if r.options.MAX_ENTRIES == 0 {
		panic("MAX entries was 0")
	}
	if r.built {
		log.Fatal("Tree is static, cannot load twice")
	}
	r.built = true

	isPooledMemReceived := false
	var rtreePooledMem *pooledMem32
	if r.options.RTreePool != nil {
		rtreePoolMemCandidate := r.options.RTreePool.Get()
		if rtreePoolMemCandidate != nil {
			rtreePooledMem = rtreePoolMemCandidate.(*pooledMem32)
			isPooledMemReceived = true
		}
	}
	if isPooledMemReceived && cap(rtreePooledMem.sorterBuffer) >= r.options.MAX_ENTRIES+1 {
		r.sorterBuffer = rtreePooledMem.sorterBuffer[0: 0]
	} else {
		r.sorterBuffer = make([]int, 0, r.options.MAX_ENTRIES+1)
	}
	r.points = points
	if isPooledMemReceived && cap(rtreePooledMem.nodes) >= computeSize(points.Len()) {
		r.nodes = rtreePooledMem.nodes[0: 0]
	} else {
		r.nodes = make([]Node32, 0, computeSize(points.Len()))
	}

	rootNodeConstruct := r.build(points, isSorted)
	if isPooledMemReceived && r.options.UnsafeConcurrencyMode && cap(rtreePooledMem.sq) >= rootNodeConstruct.height*r.options.MAX_ENTRIES {
		r.unsafeQueue = rtreePooledMem.sq
	} else {
		if r.options.UnsafeConcurrencyMode {
			r.unsafeQueue = make(searchQueue32, rootNodeConstruct.height*r.options.MAX_ENTRIES)
		} else {
			r.queuePool = sync.Pool{
				New: func () interface {} {
					return make(searchQueue32, rootNodeConstruct.height*r.options.MAX_ENTRIES)
				},
			}
			firstQueue := r.queuePool.Get()
			r.queuePool.Put(firstQueue)
		}
	}
	return r
}

func (r *SimpleRTree32) build(points FlatPoints32, isSorted bool) nodeConstruct {
	r.nodes = append(r.nodes, Node32{})
	rootNodeConstruct := nodeConstruct{
		height: int(math.Ceil(math.Log(float64(points.Len())) / math.Log(float64(r.options.MAX_ENTRIES)))),
		start:  0,
		end:    points.Len(),
	}

	r.buildNodeDownwards(&r.nodes[0], rootNodeConstruct, isSorted)
	return rootNodeConstruct
}

func (r *SimpleRTree32) buildNodeDownwards(n *Node32, nc nodeConstruct, isSorted bool) VectorBBox32 {
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
		sortX := xSorter32{n: n, points: r.points, start: nc.start, end: nc.end, bucketSize: N1}
		sortX.Sort(r.sorterBuffer)
	}
	nodeConstructs := [MAX_POSSIBLE_SIZE]nodeConstruct{}
	var nodeConstructIndex int8
	firstChildIndex := len(r.nodes)
	for i := 0; i < N; i += N1 {
		right2 := minInt(i+N1, N)
		sortY := ySorter32{n: n, points: r.points, start: nc.start + i, end: nc.start + right2, bucketSize: N2}
		sortY.Sort(r.sorterBuffer)
		for j := i; j < right2; j += N2 {
			right3 := minInt(j+N2, right2)
			child := Node32{}
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
	n.firstChildOffset = uint32(firstChildIndex) * uint32(NODE32_SIZE)
	n.nChildren = nodeConstructIndex
	// compute children
	var i int8
	bbox := r.buildNodeDownwards(&r.nodes[firstChildIndex], nodeConstructs[i], false)
	for i = 1; i < nodeConstructIndex; i++ {
		// TODO check why using (*Node32)f here does not work
		bbox2 := r.buildNodeDownwards(&r.nodes[firstChildIndex+int(i)], nodeConstructs[i], false)
		bbox = vectorBBoxExtend32(bbox, bbox2)
	}
	n.BBox = bbox
	return bbox
}

func (r *SimpleRTree32) setLeafNode(n *Node32, nc nodeConstruct) VectorBBox32 {
	// Here we follow original rbush implementation.
	firstChildIndex := nc.start

	x0, y0 := r.points.GetPointAt(nc.start)
	vb := VectorBBox32{x0, y0, x0, y0}

	for i := 1; i < nc.end-nc.start; i++ {
		x1, y1 := r.points.GetPointAt(nc.start + i)
		vb1 := [4]float32{
			x1,
			y1,
			x1,
			y1,
		}
		vb = vectorBBoxExtend32(vb, vb1)
	}
	n.firstChildOffset = uint32(firstChildIndex) * uint32(FLAT32_POINT_SIZE) // We access leafs on the original array
	n.nChildren = int8(nc.end - nc.start)
	n.nodeType = PRELEAF
	n.BBox = vb
	return vb
}

func (r *SimpleRTree32) toJSON() {
	text := make([]string, 0)
	fmt.Println(strings.Join(r.toJSONAcc(&r.nodes[0], text), ","))
}

func (r *SimpleRTree32) toJSONAcc(n *Node32, text []string) []string {
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
	f := unsafe.Pointer(uintptr(unsafe.Pointer(&r.nodes[0])) + uintptr(n.firstChildOffset))
	var i int8
	for i = 0; i < n.nChildren; i++ {
		cn := (*Node32)(f)
		text = r.toJSONAcc(cn, text)
		f = unsafe.Pointer(uintptr(f) + NODE32_SIZE)
	}
	return text
}

// node is point, there is only one distance
func computeLeafDistance32(px, py, x, y float32) float32 {
	return (x-px)*(x-px) +
		(y-py)*(y-py)
}
func computeDistances32(bbox VectorBBox32, x, y float32) (mind, maxd float32) {
	// TODO try simd
	minX := bbox[0]
	minY := bbox[1]
	maxX := bbox[2]
	maxY := bbox[3]
	minx, maxx := sortFloats32((x-minX)*(x-minX), (x-maxX)*(x-maxX))
	miny, maxy := sortFloats32((y-minY)*(y-minY), (y-maxY)*(y-maxY))

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
	maxd = minFloat32(maxx+miny, minx+maxy)
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
func vectorComputeDistances32(bbox VectorBBox32, x, y float32) (mind, maxd float32)

func minFloat32(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

func maxFloat32(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

type FlatPoints32 []float32

func (fp FlatPoints32) Len() int {
	return len(fp) / 2
}

func (fp FlatPoints32) Swap(i, j int) {
	fp[2*i], fp[2*i+1], fp[2*j], fp[2*j+1] = fp[2*j], fp[2*j+1], fp[2*i], fp[2*i+1]
}

func (fp FlatPoints32) GetPointAt(i int) (x1, y1 float32) {
	return fp[2*i], fp[2*i+1]
}

func sortFloats32(x1, x2 float32) (x3, x4 float32) {
	if x1 > x2 {
		return x2, x1
	}
	return x1, x2
}

func computeSize32(n int) (size int) {
	return n
}
