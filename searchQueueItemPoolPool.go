package SimpleRTree

// pool of pools
type searchQueueItemPoolPool struct {
	channel   chan *searchQueueItemPool
	poolsSize int
}

func newSearchQueueItemPoolPool(nPools, poolsSize int) *searchQueueItemPoolPool {
	p := new(searchQueueItemPoolPool)
	p.poolsSize = poolsSize
	p.channel = make(chan *searchQueueItemPool, nPools)
	for i := 0; i < nPools; i++ {
		item := newSearchQueuItemPool(poolsSize)
		p.channel <- item
	}
	return p
}

func (p searchQueueItemPoolPool) take() *searchQueueItemPool {
	select {
	case obj := <-p.channel:
		return obj
	default:
		return newSearchQueuItemPool(p.poolsSize)
	}
}

func (p searchQueueItemPoolPool) giveBack(item *searchQueueItemPool) {
	p.channel <- item
}
