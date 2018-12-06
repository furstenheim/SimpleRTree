// Simple RTree is a blazingly fast and GC friendly RTree. It handles 1 Million points in 1.2 microseconds for nearest point search
// (measured in an i7 with 16Gb of RAM). It is GC friendly, queries require 0 allocations.
// Building the index requires exactly 8 allocations and approximately 40 bytes per coordinate.
// That is, an index for 1 million points requires approximately 40Mb of RAM.

// To achieve this speed, the index has three restrictions. It is static, once built it cannot be changed.
// It only accepts points coordinates, no bboxes, lines or ids. And it only accepts (for now) one query, closest point to a given coordinate.

// Beware, to achieve top performance one of the hot functions has been rewritten in assembly.
// Library works in x86 but it probably won't work in other architectures. PRs are welcome to fix this deficiency.
//
// Basic Usage
//
// The format of the points is a single array where each too coordinates represent a point.
//   import "SimpleRTree"
//   points := []float64{0.0, 0.0, 1.0, 1.0} // array of two points 0, 0 and 1, 1
//   fp := SimpleRTree.FlatPoints(points)
//   r := SimpleRTree.New().Load(fp)
//   closestX, closestY, distanceSquared, found := r.FindNearestPoint(1.0, 3.0)
//   // 1.0, 1.0, 4.0, true
//
// Algorithm
//
// Index is built using STR (sort tile recursive) method https://en.wikipedia.org/wiki/R*_tree . Points are sorted into buckets using
// Floyd-Rivest algorithm https://en.wikipedia.org/wiki/Floyd%E2%80%93Rivest_algorithm.
// During query time, nodes are traversed using a linear priority queue, so that most probable candidates are visited first.
// Possible distances from the given point to the BBox of the node is computed in Assembly for maximum speed.
//
//
// Benchmarks
//
// These are benchmarks for finding closest point. A uniform distribution of points has been used. If the point distribution is skewed times will vary.
// Nevertheless, STRTrees are known to behave well against skewed geometric information. The number represents the number of points.
//
//    BenchmarkSimpleRTree_FindNearestPoint/10-4      	10000000	       124 ns/op
//    BenchmarkSimpleRTree_FindNearestPoint/1000-4    	 3000000	       465 ns/op
//    BenchmarkSimpleRTree_FindNearestPoint/10000-4   	 2000000	       670 ns/op
//    BenchmarkSimpleRTree_FindNearestPoint/100000-4  	 2000000	       871 ns/op
//    BenchmarkSimpleRTree_FindNearestPoint/200000-4  	 2000000	       961 ns/op
//    BenchmarkSimpleRTree_FindNearestPoint/1000000-4 	 1000000	      1210 ns/op
//    BenchmarkSimpleRTree_FindNearestPoint/10000000-4         	 1000000	      1593 ns/op
//
// Comparison with other RTrees
//
// It is hard to find a good comparison of RTrees, best available can be found at https://github.com/Sizmek/rtree2d which gathers benchmarks for
// RTrees in Java and Scala. According to that benchmark query times for 10M points vary from 2.7 microseconds to over 100 microseconds,
// so SimpleRTree performs fairly better. Times are even comparable with boost, R*-tree layout performs in 1.2 microseconds for 1 Million points https://www.boost.org/doc/libs/1_64_0/libs/geometry/doc/html/geometry/spatial_indexes/introduction.html
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

const MAX_POSSIBLE_SIZE = 9

// SimpleRTree is the main structure of the library
type SimpleRTree struct {
	options Options
	nodes   []rNode
	points  FlatPoints
	built   bool
	queuePool         sync.Pool
	unsafeQueue         searchQueue // Only used in unsafe mode
	sorterBuffer      []int // floyd rivest requires a bucket, we allocate it once and reuse
}

// FlatPoints is the input format for coordinates
// It consists on a flat array, x coordinates are stored in even positions
// y coordinates are stored in odd coordinates
// []float64{0, 0, 2, 4} corresponds to the points (0, 0) and (2, 4)
//
// Note: rtree is assumed to have sole access to the array, it will modify the underlying order and it
// will return wrong results if the elements are modified
type FlatPoints []float64

// SimpleRTree can be slightly tuned
type Options struct {
	UnsafeConcurrencyMode bool // Set this parameter to true if you only intend to access the R tree from one go routine. FindNearestPoint will be faster and it will make no allocations
	MAX_ENTRIES int
	RTreePool *sync.Pool // If a lot of RTrees are being created you can provide a pool to the tree. On destroy the underlying memory space will be saved back to the pool, so next tree can use it
}

type rNode struct {
	nodeType         nodeType
	nChildren        int8
	// Here we save firstChild - firstNode. That means that there is there is a theoretical upper limit to the tree of
	// maxuint32 / node_size = 4294967295 / 40 = 107374182 ~ 100M
	firstChildOffset uint32
	BBox             rVectorBBox
}
type nodeType int8
const (
	default_node = iota
	preleaf_node
)

var node_size = unsafe.Sizeof(rNode{})
var flat_point_size =unsafe.Sizeof([2]float64{})
var float_size = uintptr(unsafe.Sizeof([1]float64{}))

type pooledMem struct {
	sorterBuffer []int
	sq searchQueue
	nodes []rNode
}

// Structure used to constructing the ndoe
type nodeConstruct struct {
	height     int
	start, end uint32 // index in the underlying array
}

// New returns an instance of an RTree with default options
func New() *SimpleRTree {
	defaultOptions := Options{
		MAX_ENTRIES: MAX_POSSIBLE_SIZE,
	}
	return NewWithOptions(defaultOptions)
}

// NewWithOptions returns an instance of an RTree with given options o
func NewWithOptions(o Options) *SimpleRTree {
	r := &SimpleRTree{
		options: o,
	}
	if o.MAX_ENTRIES > MAX_POSSIBLE_SIZE {
		panic(fmt.Sprintf("Cannot exceed %d for size", MAX_POSSIBLE_SIZE))
	}
	if o.MAX_ENTRIES == 0 {
		r.options.MAX_ENTRIES = MAX_POSSIBLE_SIZE
	}
	return r
}

// Destroy frees up resources that are held within the RTree
// This method is useful if RTree was provided a pool, so used memory is return to the pool
// Example:
//   pool := &sync.Pool{}
//   r1 := SimpleRTree.NewWithOptions(SimpleRTree.Options{RTreePool: pool})
//   /* Perform some queries */
//   /* r1 is no longer needed */
//   r1.Destroy()
//r2 will be used with the same underlying objects, thus avoiding heap allocations
//   r2 := SimpleRTree.NewWithOptions(SimpleRTree.Options{RTreePool: pool})
func (r *SimpleRTree) Destroy () {
	if r.options.RTreePool != nil {
		r.options.RTreePool.Put(
			&pooledMem{
				sorterBuffer: r.sorterBuffer,
				sq: r.unsafeQueue,
				nodes: r.nodes,
			},
		)
	}
}
// Load accepts points, an flat array of coordinates and builds the RTree
//
// Note: rtree is assumed to have sole access to the array, it will modify the underlying order and it
// will return wrong results if the elements are modified
func (r *SimpleRTree) Load(points FlatPoints) *SimpleRTree {
	return r.load(points, false)
}

// LoadSortedArray accepts a sorted flat array of coordinates and builds the RTree
// points must have been sorted lexicographically. That is
// (x1, y1) < (x2, y2) if x1 < x2 or x1 === x2 and y1 < y2
//
// Note: rtree is assumed to have sole access to the array, it will modify the underlying order and it
// will return wrong results if the elements are modified
func (r *SimpleRTree) LoadSortedArray(points FlatPoints) *SimpleRTree {
	return r.load(points, true)
}

// FindNearestPoint will return the coordinates of the closest point
// to the provided coordinates x and y. The function returns three parameters
// x1, y1 coordinates of the point
// d1 distances squared to that point. That is |x1 - x|**2 + |y1 - y|**2
//  x1, y1, d1 := r.FindNearestPointWithin(x, y, 4)
//  (x1 - x) * (x1 - x) + (y1 - y) * (y1 - y) == d1
func (r *SimpleRTree) FindNearestPoint(x, y float64) (x1, y1, d1 float64) {
	x1, y1, d1, _ = r.FindNearestPointWithin(x, y, math.Inf(1))
	return
}

// FindNearestPointWithin will return the closest point
// to the provided coordinates x and y, within the distance squared dsquared.
// The method returns coordinates of the point x1, y1.
// Distance squared to the point d1 and a boolean found
// In case there is no point satisfying |x1 - x|**2 + |y1 - y| ** 2 < dsquared
// found will return false
//  x1, y1, d1, found := r.FindNearestPointWithin(x, y, 4)
// (x1 - x) * (x1 - x) + (y1 - y) * (y1 - y) < 4
func (r *SimpleRTree) FindNearestPointWithin(x, y, dsquared float64) (x1, y1, d1 float64, found bool) {
	var minItem searchQueueItem
	distanceLowerBound := math.Inf(1)
	// if bbox is further from this bound then we don't explore it
	distanceUpperBound := dsquared
	var sq searchQueue
	if r.options.UnsafeConcurrencyMode {
		sq = r.unsafeQueue
	} else {
		sq = r.queuePool.Get().(searchQueue)
	}
	sq = sq[0:0]

	rootNode := &r.nodes[0]
	unsafeRootLeafNode := uintptr(unsafe.Pointer(&r.points[0]))
	unsafeRootNode := uintptr(unsafe.Pointer(rootNode))
	sq = append(sq, searchQueueItem{node: uintptr(unsafe.Pointer(rootNode)), distance: 0}) // we don't need distance for first node

	for sq.Len() > 0 {
		sq.PreparePop()
		item := sq[sq.Len() - 1]
		sq = sq[0: sq.Len() - 1]
		currentDistance := item.distance
		if found && currentDistance > distanceLowerBound {
			break
		}

		node := (*rNode)(unsafe.Pointer(item.node))
		if node == nil { // Leaf
			// we know it is smaller from the previous test
			distanceLowerBound = currentDistance
			minItem = item
			found = true
			continue
		}
		switch node.nodeType {
		case preleaf_node:
			f := unsafeRootLeafNode + uintptr(node.firstChildOffset)
			var i int8
			for i = node.nChildren; i>0; i-- {
				px := *(*float64)(unsafe.Pointer(f))
				f = f + float_size
				py := *(*float64)(unsafe.Pointer(f))

				d := computeLeafDistance(px, py, x, y)
				if d <= distanceUpperBound {
					sq = append(sq, searchQueueItem{node: uintptr(unsafe.Pointer(nil)), px: px, py: py, distance: d})
					distanceUpperBound = d
				}
				f = f + float_size
			}
		default:
			f := unsafeRootNode + uintptr(node.firstChildOffset)
			var i int8
			for i = node.nChildren; i>0; i-- {
				n := (*rNode)(unsafe.Pointer(f))
				mind, maxd := vectorComputeDistances(n.BBox, x, y)
				if mind <= distanceUpperBound {
					sq = append(sq, searchQueueItem{node: uintptr(unsafe.Pointer(n)), distance: mind})
					// Distance to one of the corners is lower than the upper bound
					// so there must be a point at most within distanceUpperBound
					if maxd < distanceUpperBound {
						distanceUpperBound = maxd
					}
				}
				f = f + node_size
			}
		}
	}

	// return heap
	if !r.options.UnsafeConcurrencyMode {
		r.queuePool.Put(sq)
	} else {
		r.unsafeQueue = sq
	}

	if !found {
		return
	}
	x1 = minItem.px
	y1 = minItem.py
	d1 = distanceUpperBound
	return
}

func (r *SimpleRTree) load(points FlatPoints, isSorted bool) *SimpleRTree {
	if points.Len() == 0 {
		return r
	}
	if points.Len() >= math.MaxUint32 / int(node_size) {
		log.Fatal("Exceded maximum possible size", math.MaxUint32 / int(node_size))
	}
	if r.options.MAX_ENTRIES == 0 {
		panic("MAX entries was 0")
	}
	if r.built {
		log.Fatal("Tree is static, cannot load twice")
	}
	r.built = true

	isPooledMemReceived := false
	var rtreePooledMem *pooledMem
	if r.options.RTreePool != nil {
		rtreePoolMemCandidate := r.options.RTreePool.Get()
		if rtreePoolMemCandidate != nil {
			rtreePooledMem = rtreePoolMemCandidate.(*pooledMem)
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
		r.nodes = make([]rNode, 0, computeSize(points.Len()))
	}

	rootNodeConstruct := r.build(points, isSorted)
	if isPooledMemReceived && r.options.UnsafeConcurrencyMode && cap(rtreePooledMem.sq) >= rootNodeConstruct.height*r.options.MAX_ENTRIES {
		r.unsafeQueue = rtreePooledMem.sq
	} else {
		if r.options.UnsafeConcurrencyMode {
			r.unsafeQueue = make(searchQueue, rootNodeConstruct.height*r.options.MAX_ENTRIES)
		} else {
			r.queuePool = sync.Pool{
				New: func () interface {} {
					return make(searchQueue, rootNodeConstruct.height*r.options.MAX_ENTRIES)
				},
			}
			firstQueue := r.queuePool.Get()
			r.queuePool.Put(firstQueue)
		}
	}
	return r
}

func (r *SimpleRTree) build(points FlatPoints, isSorted bool) nodeConstruct {
	r.nodes = append(r.nodes, rNode{})
	rootNodeConstruct := nodeConstruct{
		height: int(math.Ceil(math.Log(float64(points.Len())) / math.Log(float64(r.options.MAX_ENTRIES)))),
		start:  uint32(0),
		end:    uint32(points.Len()),
	}

	r.buildNodeDownwards(&r.nodes[0], rootNodeConstruct, isSorted)
	return rootNodeConstruct
}

func (r *SimpleRTree) buildNodeDownwards(n *rNode, nc nodeConstruct, isSorted bool) rVectorBBox {
	N := int(nc.end - nc.start)
	// target number of root entries to maximize storage utilization
	var M float64
	if N <= r.options.MAX_ENTRIES { // Leaf node
		return r.setLeafNode(n, nc)
	}

	M = math.Ceil(float64(N) / float64(math.Pow(float64(r.options.MAX_ENTRIES), float64(nc.height-1))))

	N2 := int(math.Ceil(float64(N) / M))
	N1 := N2 * int(math.Ceil(math.Sqrt(M)))

	start := int(nc.start)
	// parent node might already be sorted. In that case we avoid double computation
	if !isSorted {
		sortX := xSorter{n: n, points: r.points, start: start, end: int(nc.end), bucketSize: N1}
		sortX.Sort(r.sorterBuffer)
	}
	nodeConstructs := [MAX_POSSIBLE_SIZE]nodeConstruct{}
	var nodeConstructIndex int8
	firstChildIndex := len(r.nodes)
	for i := 0; i < N; i += N1 {
		right2 := minInt(i+N1, N)
		sortY := ySorter{n: n, points: r.points, start: start+ i, end: start+ right2, bucketSize: N2}
		sortY.Sort(r.sorterBuffer)
		for j := i; j < right2; j += N2 {
			right3 := minInt(j+N2, right2)
			child := rNode{}
			childC := nodeConstruct{
				start:  nc.start + uint32(j),
				end:    nc.start + uint32(right3),
				height: nc.height - 1,
			}
			r.nodes = append(r.nodes, child)
			nodeConstructs[nodeConstructIndex] = childC
			nodeConstructIndex++
		}
	}
	n.firstChildOffset = uint32(firstChildIndex) * uint32(node_size)
	n.nChildren = nodeConstructIndex
	// compute children
	var i int8
	bbox := r.buildNodeDownwards(&r.nodes[firstChildIndex], nodeConstructs[i], false)
	for i = 1; i < nodeConstructIndex; i++ {
		// TODO check why using (*Node)f here does not work
		bbox2 := r.buildNodeDownwards(&r.nodes[firstChildIndex+int(i)], nodeConstructs[i], false)
		bbox = vectorBBoxExtend(bbox, bbox2)
	}
	n.BBox = bbox
	return bbox
}

func (r *SimpleRTree) setLeafNode(n *rNode, nc nodeConstruct) rVectorBBox {
	// Here we follow original rbush implementation.
	start := int(nc.start)
	end := int(nc.end)
	firstChildIndex := start

	x0, y0 := r.points.GetPointAt(start)
	vb := rVectorBBox{x0, y0, x0, y0}

	for i := end-start - 1; i > 0; i-- {
		x1, y1 := r.points.GetPointAt(start + i)
		vb1 := [4]float64{
			x1,
			y1,
			x1,
			y1,
		}
		vb = vectorBBoxExtend(vb, vb1)
	}
	n.firstChildOffset = uint32(firstChildIndex) * uint32(flat_point_size) // We access leafs on the original array
	n.nChildren = int8(nc.end - nc.start)
	n.nodeType = preleaf_node
	n.BBox = vb
	return vb
}

func (r *SimpleRTree) toJSON() {
	text := make([]string, 0)
	fmt.Println(strings.Join(r.toJSONAcc(&r.nodes[0], text), ","))
}

func (r *SimpleRTree) toJSONAcc(n *rNode, text []string) []string {
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
		cn := (*rNode)(f)
		text = r.toJSONAcc(cn, text)
		f = unsafe.Pointer(uintptr(f) + node_size)
	}
	return text
}

// node is point, there is only one distance
func computeLeafDistance(px, py, x, y float64) float64 {
	return (x-px)*(x-px) +
		(y-py)*(y-py)
}
func computeDistances(bbox rVectorBBox, x, y float64) (mind, maxd float64) {
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
func vectorComputeDistances(bbox rVectorBBox, x, y float64) (mind, maxd float64)

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
	return n
}
