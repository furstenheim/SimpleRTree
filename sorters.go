package SimpleRTree

import "github.com/furstenheim/FloydRivest"

type xSorter struct {
	n          *Node
	points Interface
	start, end, bucketSize int
}

func (s xSorter) Less(i, j int) bool {
	x1, _ := s.points.GetPointAt(i + s.start + s.n.start)
	x2, _ := s.points.GetPointAt(j + s.start + s.n.start)
	return x1 < x2
}

func (s xSorter) Swap(i, j int) {
	s.points.Swap(i+s.start + s.n.start, j+s.start + s.n.start)
}

func (s xSorter) Len() int {
	return s.end - s.start
}

func (s xSorter) Sort() {
	FloydRivest.Buckets(s, s.bucketSize)
}

type ySorter struct {
	n          *Node
	points Interface
	start, end, bucketSize int
}

func (s ySorter) Less(i, j int) bool {
	_, y1 := s.points.GetPointAt(i + s.start + s.n.start)
	_, y2 := s.points.GetPointAt(j + s.start + s.n.start)
	return y1 < y2
}

func (s ySorter) Swap(i, j int) {
	s.points.Swap(i+s.start + s.n.start, j+s.start + s.n.start)
}

func (s ySorter) Len() int {
	return s.end - s.start
}
func (s ySorter) Sort() {
	// we already do the shifting on the sort functions
	FloydRivest.Buckets(s, s.bucketSize)
}
