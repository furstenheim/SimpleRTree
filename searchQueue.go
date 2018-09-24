package SimpleRTree

type searchQueueItem struct {
	nodeIndex int
	distance float64
}


type searchQueue []*searchQueueItem

func (sq searchQueue) Len () int {
	return len(sq)
}

/**
Inlined heap for improved performance

 */
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package heap provides heap operations for any type that implements
// heap.Interface. A heap is a tree with the property that each node is the
// minimum-valued node in its subtree.
//
// The minimum element in the tree is the root, at index 0.
//
// A heap is a common way to implement a priority queue. To build a priority
// queue, implement the Heap interface with the (negative) priority as the
// ordering for the Less method, so Push adds items while Pop removes the
// highest-priority item from the queue. The Examples include such an
// implementation; the file example_pq_test.go has the complete source.
//

// The Interface type describes the requirements
// for a type using the routines in this package.
// Any type that implements it may be used as a
// min-heap with the following invariants (established after
// Init has been called or if the data is empty or sorted):
//
//	!h.Less(j, i) for 0 <= i < h.Len() and 2*i+1 <= j <= 2*i+2 and j < h.Len()
//
// Note that Push and Pop in this interface are for package heap's
// implementation to call. To add and remove things from the heap,
// use heap.Push and heap.Pop.


// Init establishes the heap invariants required by the other routines in this package.
// Init is idempotent with respect to the heap invariants
// and may be called whenever the heap invariants may have been invalidated.
// Its complexity is O(n) where n = h.Len().
func (h searchQueue) Init() {
	// heapify
	n := h.Len()
	for i := n/2 - 1; i >= 0; i-- {
		h.down(i, n)
	}
}

func (sq searchQueue) Swap (i, j int) {
	sq[i], sq[j] = sq[j], sq[i]
}

// Push pushes the element x onto the heap. The complexity is
// O(log(n)) where n = h.Len().
//
func (h *searchQueue) Push(x *searchQueueItem) {
	*h = append(*h, x)
	h.up(h.Len()-1)
}

// Pop removes the minimum element (according to Less) from the heap
// and returns it. The complexity is O(log(n)) where n = h.Len().
// It is equivalent to Remove(h, 0).
//
func (h *searchQueue) Pop() *searchQueueItem {
	n := h.Len() - 1
	h.Swap(0, n)
	h.down(0, n)
	arr := *h
	item := arr[h.Len() - 1]
	*h = arr[0: h.Len() - 1]
	return item
}

// Remove removes the element at index i from the heap.
// The complexity is O(log(n)) where n = h.Len().
//
func (h searchQueue) Remove(i int) interface{} {
	n := h.Len() - 1
	if n != i {
		h.Swap(i, n)
		if !h.down(i, n) {
			h.up(i)
		}
	}
	return h.Pop()
}

// Fix re-establishes the heap ordering after the element at index i has changed its value.
// Changing the value of the element at index i and then calling Fix is equivalent to,
// but less expensive than, calling Remove(h, i) followed by a Push of the new value.
// The complexity is O(log(n)) where n = h.Len().
func (h searchQueue) Fix(i int) {
	if !h.down(i, h.Len()) {
		h.up(i)
	}
}

func (h searchQueue) up(j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !h.Less(j, i) {
			break
		}
		h.Swap(i, j)
		j = i
	}
}

func (h searchQueue) down(i0, n int) bool {
	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && h.Less(j2, j1) {
			j = j2 // = 2*i + 2  // right child
		}
		if !h.Less(j, i) {
			break
		}
		h.Swap(i, j)
		i = j
	}
	return i > i0
}

func (sq searchQueue) Less(i, j int) bool{
	// We want to pop element with smaller distance
	return sq[i].distance < sq[j].distance
}


func (sq *searchQueue) Empty() {
	arr := *sq
	*sq = arr[0: 0]
}