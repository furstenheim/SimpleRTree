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

/**
Linear heap for performance boost (actually log in len(points)

**/
func (h searchQueue) Init() {
	// We start with empty queues
}

func (sq searchQueue) Swap(i, j int) {
	sq[i], sq[j] = sq[j], sq[i]
}

// Push pushes the element x onto the heap. The complexity is
// O(1)
func (h *searchQueue) Push(x searchQueueItem) {
	*h = append(*h, x)
	h.up(h.Len() - 1)
}

// Pop removes the minimum element (according to Less) from the heap
// and returns it. The complexity is O(n) where n = h.Len().
func (h *searchQueue) Pop() searchQueueItem {
	n := h.Len() - 1
	h.Swap(0, n)
	h.down(0, n)
	arr := *h
	item := arr[n]
	*h = arr[0 : n]
	return item
}

func (h searchQueue) up(j int) {
	if h.Less(j, 0) {
		h.Swap(0, j)
	}
}

func (h searchQueue) down(i0, n int) {
	for j := 1; j < n; j++ {
		if h.Less(j, 0) {
			h.Swap(0, j)
		}
	}
}

func (sq searchQueue) Less(i, j int) bool {
	// We want to pop element with smaller distance
	return sq[i].distance < sq[j].distance
}

func (sq *searchQueue) Empty() {
	arr := *sq
	*sq = arr[0:0]
}
