package memepool

import "container/heap"

const memPoolCapacity = 5000

type memPoolItem struct {
	pairItem *memPoolItem
	idx          int
	priorityItem PriorityItem
}

type PriorityItem interface {
	Priority() float64
}

type memPoolQueue []*memPoolItem

type memPoolMaxQueue struct {
	memPoolQueue
}

func (q memPoolMaxQueue) Less(i, j int) bool {
	return q.memPoolQueue[i].priorityItem.Priority() > q.memPoolQueue[j].priorityItem.Priority()
}

type memPoolMinQueue struct {
	memPoolQueue
}

func (q memPoolMinQueue) Less(i, j int) bool {
	return q.memPoolQueue[i].priorityItem.Priority() < q.memPoolQueue[j].priorityItem.Priority()
}

func (q memPoolQueue) Len() int {
	return len(q)
}

func (q memPoolQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].idx = i
	q[j].idx = j
}

type MemPool struct {
	minQueue *memPoolMinQueue
	maxQueue *memPoolMaxQueue
}

func (q *memPoolQueue) Push(x interface{}) {
	n := len(*q)
	item := x.(*memPoolItem)
	item.idx = n
	*q = append(*q, item)
}

func (q *memPoolQueue) Pop() interface{} {
	old := *q
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.idx = -1 // for safety
	*q = old[0 : n-1]
	return item
}

// Push adds an item to mempool
// When the capacity is surpassed, item with the lowest priority will be dropped
func (m *MemPool) Push(tx PriorityItem) {
	minItem := memPoolItem{priorityItem: tx}
	maxItem := memPoolItem{priorityItem: tx}
	maxItem.pairItem = &minItem
	minItem.pairItem = &maxItem

	heap.Push(m.minQueue, &minItem)
	heap.Push(m.maxQueue, &maxItem)

	if m.minQueue.Len() > memPoolCapacity {
		item := heap.Pop(m.minQueue)
		tx := item.(*memPoolItem)
		heap.Remove(m.maxQueue, tx.pairItem.idx)
	}
}

// Pop retrieves and removes the item with the highest priority
// If the mempool is empty, the function will panic
func (m *MemPool) Pop() PriorityItem {
	item := heap.Pop(m.maxQueue)
	if item == nil {
		return nil
	}
	tx := item.(*memPoolItem)
	heap.Remove(m.minQueue, tx.pairItem.idx)
	return tx.priorityItem
}

func (m MemPool) Len() int {
	return m.maxQueue.Len()
}

func NewMemPool() MemPool {
	m := MemPool{
		minQueue: &memPoolMinQueue{memPoolQueue: memPoolQueue{}},
		maxQueue: &memPoolMaxQueue{memPoolQueue: memPoolQueue{}},
	}

	heap.Init(m.minQueue)
	heap.Init(m.maxQueue)
	return m
}
