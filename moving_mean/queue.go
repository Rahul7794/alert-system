package moving_mean

import "sync"

type Queue struct {
	items []float64
	mutex sync.Mutex
}

func (queue *Queue) Enqueue(item float64) {
	queue.mutex.Lock()
	defer queue.mutex.Unlock()

	queue.items = append(queue.items, item)
}

func (queue *Queue) Dequeue() float64 {
	queue.mutex.Lock()
	defer queue.mutex.Unlock()

	if len(queue.items) == 0 {
		return 0.0
	}

	lastItem := queue.items[0]
	queue.items = queue.items[1:]

	return lastItem
}

func (queue *Queue) Size() int {
	return len(queue.items)
}

func (queue *Queue) Reset() {
	queue.mutex.Lock()
	defer queue.mutex.Unlock()

	queue.items = nil
}
