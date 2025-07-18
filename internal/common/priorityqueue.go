package common

import (
	"golang.org/x/exp/constraints"
)

// An PQItem is something we manage in a priority queue.
type PQItem[T any, N constraints.Integer] struct {
	Content  T // The content of the item; arbitrary.
	Priority N // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface
	// methods.
	index int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue[T any, N constraints.Integer] []*PQItem[T, N]

func (pq PriorityQueue[T, N]) Len() int { return len(pq) }

func (pq PriorityQueue[T, N]) Less(i, j int) bool {
	// We want Pop to give us the lowest, not highest, priority so we use
	// greater than here.
	return pq[i].Priority < pq[j].Priority
}

func (pq PriorityQueue[T, N]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue[T, N]) Push(x any) {
	n := len(*pq)
	item := x.(*PQItem[T, N])
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue[T, N]) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // don't stop the GC from reclaiming the item eventually
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// Returns item from top of the heap, without removing it from heap
func (pq *PriorityQueue[T, N]) GetTop() *PQItem[T, N] {
	if pq.Len() == 0 {
		return nil
	}
	return (*pq)[0]
}
