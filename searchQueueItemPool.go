package SimpleRTree

// Reuse items to avoid allocation

type searchPool chan *searchQueueItem

func newSearchPool (total int) * searchPool {
	p := make(searchPool, total)
	for i:= 0; i < total; i++ {
		p <- new(searchQueueItem)
	}
	return &p
}

func (p searchPool) take () *searchQueueItem {
	select {
	case obj := <-p:
		return obj
	default:
		return new(searchQueueItem)
	}
}

func (p searchPool) giveBack (item *searchQueueItem) {
	p <- item
}