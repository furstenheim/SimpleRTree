package SimpleRTree

type searchQueueItem struct {
	node     *Node // if nil item carries node
	px, py float64 // points are not stored in nodes so we need to track them explicitely
	distance float64
}

type searchQueue []searchQueueItem

func (sq searchQueue) Len() int {
	return len(sq)
}

func (sq searchQueue) Swap(i, j int) {
	sq[i], sq[j] = sq[j], sq[i]
}

func (sq searchQueue) PreparePop () {
	n := sq.Len() - 1
	for j := 0; j < n; j++ {
		if sq.Less(j, n) {
			sq.Swap(n, j)
		}
	}
}

func (sq searchQueue) Less(i, j int) bool {
	// We want to pop element with smaller distance
	return sq[i].distance < sq[j].distance
}
