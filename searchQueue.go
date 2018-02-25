package SimpleRTree

type searchQueueItem struct {
	node *Node
	distance float64
}


type searchQueue []*searchQueueItem

func (sq searchQueue) Len () int {
	return len(sq)
}

func (sq searchQueue) Swap (i, j int) {
	sq[i], sq[j] = sq[j], sq[i]
}

func (sq *searchQueue) Push (x interface{}) {
	item := x.(*searchQueueItem)
	*sq = append(*sq, item)
}

func (sq searchQueue) Less(i, j int) bool{
	// We want to pop element with smaller distance
	return sq[i].distance < sq[j].distance
}

func (sq *searchQueue) Pop() interface{} {
	arr := *sq
	item := arr[sq.Len() - 1]
	*sq = arr[0: sq.Len() - 1]
	return item
}