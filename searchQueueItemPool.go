package SimpleRTree

// Reuse items to avoid allocation

type searchQueueItemPool []*searchQueueItem

func newSearchQueuItemPool(total int) *searchQueueItemPool {
	backingItemPool := make([]searchQueueItem, total)
	p := make(searchQueueItemPool, total)

	for i := 0; i < total; i++ {
		p[i] = &backingItemPool[i]
	}
	return &p
}

func (p *searchQueueItemPool) take() *searchQueueItem {
	arr := *p
	if len(arr) == 0 {
		return new(searchQueueItem)
	}

	item := arr[len(arr)-1]
	*p = arr[0 : len(arr)-1]
	return item
}

func (p *searchQueueItemPool) giveBack(item *searchQueueItem) {
	arr := *p
	*p = append(arr, item)
}
