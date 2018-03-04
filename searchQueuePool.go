package SimpleRTree

// Reuse items to avoid allocation

type searchQueuePool chan *searchQueue

func newSearchQueuePool (total, queueSize int) * searchQueuePool {
	p := make(searchQueuePool, total)
	for i:= 0; i < total; i++ {
		item := make(searchQueue, queueSize)
		p <- &item
	}
	return &p
}

func (p searchQueuePool) take () *searchQueue {
	select {
	case obj := <-p:
		obj.Empty()
		return obj
	default:
		return new(searchQueue)
	}
}

func (p searchQueuePool) giveBack (item *searchQueue) {
	p <- item
}