package SimpleRTree

type xSorter32 struct {
	n                      *Node32
	points                 FlatPoints32
	start, end, bucketSize int
}

func (s xSorter32) Less(i, j int) bool {
	x1, _ := s.points.GetPointAt(i + s.start)
	x2, _ := s.points.GetPointAt(j + s.start)
	return x1 < x2
}

func (s xSorter32) Swap(i, j int) {
	s.points.Swap(i+s.start, j+s.start)
}

func (s xSorter32) Len() int {
	return s.end - s.start
}

func (s xSorter32) Sort(buffer []int) {
	bucketsX32(s, s.bucketSize, buffer)
}

type ySorter32 struct {
	n                      *Node32
	points                 FlatPoints32
	start, end, bucketSize int
}

func (s ySorter32) Less(i, j int) bool {
	_, y1 := s.points.GetPointAt(i + s.start)
	_, y2 := s.points.GetPointAt(j + s.start)
	return y1 < y2
}

func (s ySorter32) Swap(i, j int) {
	s.points.Swap(i+s.start, j+s.start)
}

func (s ySorter32) Len() int {
	return s.end - s.start
}
func (s ySorter32) Sort(buffer []int) {
	// we already do the shifting on the sort functions
	bucketsY32(s, s.bucketSize, buffer)
}
