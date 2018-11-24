package SimpleRTree

type searchQueueItem struct {
	node     *Node // if nil item carries node
	px, py float64 // points are not stored in nodes so we need to track them explicitely
	distance float64
}

const LOG_HEAP_SIZE = 3
const HEAP_SIZE = 1 << LOG_HEAP_SIZE

type searchQueue []searchQueueItem

func (sq searchQueue) Len() int {
	return len(sq)
}

/**
Inlined heap for improved performance
**/
func (h searchQueue) Init() {
	// heapify
	n := h.Len()
	for i := (n - 1) >> LOG_HEAP_SIZE; i >= 0; i-- {
		h.down(i, n)
	}
}

func (sq searchQueue) Swap(i, j int) {
	sq[i], sq[j] = sq[j], sq[i]
}

// Push pushes the element x onto the heap. The complexity is
// O(log(n)) where n = h.Len().
//
func (h *searchQueue) Push(x searchQueueItem) {
	*h = append(*h, x)
	h.up(h.Len() - 1)
}

// Pop removes the minimum element (according to Less) from the heap
// and returns it. The complexity is O(log(n)) where n = h.Len().
// It is equivalent to Remove(h, 0).
//
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
	for {
		i := (j - 1) >> LOG_HEAP_SIZE // parent
		if i < 0 || !h.Less(j, i) {
			break
		}
		h.Swap(j, i)
		j = i
	}
}

func (h searchQueue) down(i0, n int) {
	// child array to store indexes of all the children of given node
	i := i0
	for {
		minChildIndex := i
		// childrenIndexes[i] = -1 if node is leaf children
		for k, j:= 1, i << LOG_HEAP_SIZE + 1; k <=HEAP_SIZE; k, j = k+1, j+1 {
			if j >= n {
				break
			}
			if h.Less(j, minChildIndex) {
				minChildIndex = j
			}
		}
		if minChildIndex == i { // leaf node or not smaller
			break
		}
		// Swap if min child is smaller than the key of node
		h.Swap(i, minChildIndex)
		i = minChildIndex
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
