package memepool

import "container/heap"

// This is an implementation of priority queue with capacity which uses max heap and min heap with the same items referencing each other
// This allows to remove items from the head and the tail of the queue with O(log(n)) complexity, allowing not to discard
// the lowest priority item when the capacity is reached and keeping the O(log(n)) complexity of insert operation
type queueItem struct {
	pairItem *queueItem
	idx          int
	priorityItem PriorityItem
}

type PriorityItem interface {
	Priority() float64
}

type queue []*queueItem

type maxQueue struct {
	queue
}

func (q maxQueue) Less(i, j int) bool {
	return q.queue[i].priorityItem.Priority() > q.queue[j].priorityItem.Priority()
}

type minQueue struct {
	queue
}

func (q minQueue) Less(i, j int) bool {
	return q.queue[i].priorityItem.Priority() < q.queue[j].priorityItem.Priority()
}

func (q queue) Len() int {
	return len(q)
}

func (q queue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
	q[i].idx = i
	q[j].idx = j
}

type PriorityQueue struct {
	minQueue *minQueue
	maxQueue *maxQueue
	capacity int
}

func (q *queue) Push(x interface{}) {
	n := len(*q)
	item := x.(*queueItem)
	item.idx = n
	*q = append(*q, item)
}

func (q *queue) Pop() interface{} {
	old := *q
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.idx = -1 // for safety
	*q = old[0 : n-1]
	return item
}

// Push adds an item to the priority queue
// When the capacity is surpassed, item with the lowest priority will be dropped
func (m *PriorityQueue) Push(item PriorityItem) {
	minItem := queueItem{priorityItem: item}
	maxItem := queueItem{priorityItem: item}
	maxItem.pairItem = &minItem
	minItem.pairItem = &maxItem

	heap.Push(m.minQueue, &minItem)
	heap.Push(m.maxQueue, &maxItem)

	if m.minQueue.Len() > m.capacity {
		item := heap.Pop(m.minQueue)
		tx := item.(*queueItem)
		heap.Remove(m.maxQueue, tx.pairItem.idx)
	}
}

// Pop retrieves and removes the item with the highest priority
// If the queue is empty, the function will panic
func (m *PriorityQueue) Pop() PriorityItem {
	item := heap.Pop(m.maxQueue)
	if item == nil {
		return nil
	}
	tx := item.(*queueItem)
	heap.Remove(m.minQueue, tx.pairItem.idx)
	return tx.priorityItem
}

func (m PriorityQueue) Len() int {
	return m.maxQueue.Len()
}

func NewPriorityQueue(capacity int) PriorityQueue {
	m := PriorityQueue{
		minQueue: &minQueue{queue: queue{}},
		maxQueue: &maxQueue{queue: queue{}},
		capacity: capacity,
	}

	heap.Init(m.minQueue)
	heap.Init(m.maxQueue)
	return m
}
