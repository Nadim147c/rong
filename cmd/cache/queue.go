package cache

type queue struct {
	items []*job
}

// Enqueue adds an element to the back of the queue
func (q *queue) Enqueue(item ...*job) {
	q.items = append(q.items, item...)
}

// Dequeue removes and returns the front element of the queue
func (q *queue) Dequeue() (*job, bool) {
	if len(q.items) == 0 {
		return nil, false
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}

// NextN returns the next n items without removing them
func (q *queue) NextN(n int) []*job {
	if n < 0 {
		return nil
	}
	if n > len(q.items) {
		n = len(q.items)
	}
	return q.items[:n]
}

// Size returns the number of items in the queue
func (q *queue) Size() int {
	return len(q.items)
}
