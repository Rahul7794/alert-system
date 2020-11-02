package movingmean

import (
	"reflect"
	"testing"
)

func TestQueue(t *testing.T) {
	tests := []struct {
		name string
		push func() *Queue
	}{
		{
			name: "Test basic Operations of Queue",
			push: func() *Queue {
				queue := Queue{}
				queue.Enqueue(10.22)
				queue.Enqueue(20.22)
				queue.Enqueue(30.22)
				return &queue
			},
		},
	}
	for _, tt := range tests {
		queue := tt.push()
		dequeued, _ := queue.Dequeue()
		if !reflect.DeepEqual(*dequeued, 10.22) {
			t.Errorf("Dequeue()=%v, wanted %v", dequeued, 10.22)
		}
		if !reflect.DeepEqual(queue.Size(), 2) {
			t.Errorf("Size()=%v, wanted %v", 2, queue.Size())
		}
		queue.Reset()
	}
}
