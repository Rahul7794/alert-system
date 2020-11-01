package moving_mean

import (
	"fmt"
	"sync"
)

// Queue is implementation of Queue DataStructure in Go.
// Elements are stored in items in FIFO order.
// mutex is introduced for thread safety.
type Queue struct {
	items []float64
	mutex sync.Mutex
}

// Enqueue adds an element to the Queue.
func (queue *Queue) Enqueue(item float64) {
	queue.mutex.Lock()
	defer queue.mutex.Unlock()

	queue.items = append(queue.items, item)
}

// Dequeue remove and return the first element in Queue.
func (queue *Queue) Dequeue() (*float64, error) {
	queue.mutex.Lock()
	defer queue.mutex.Unlock()

	// handle the case where Queue is empty
	if len(queue.items) == 0 {
		return nil, fmt.Errorf("the queue is empty")
	}

	lastItem := queue.items[0]
	queue.items = queue.items[1:]

	return &lastItem, nil
}

// PeekLast return the first element in Queue.
func (queue *Queue) PeekLast() (*float64, error) {
	queue.mutex.Lock()
	defer queue.mutex.Unlock()

	// handle the case where Queue is empty
	if len(queue.items) == 0 {
		return nil, fmt.Errorf("the queue is empty")
	}

	return &queue.items[queue.Size()-1], nil
}

// Size returns the size of Queue
func (queue *Queue) Size() int {
	return len(queue.items)
}

// Reset flushes the Queue.
func (queue *Queue) Reset() {
	queue.mutex.Lock()
	defer queue.mutex.Unlock()

	queue.items = nil
}
