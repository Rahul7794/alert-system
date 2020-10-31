package moving_mean

import (
	"reflect"
	"testing"
)

func TestMovingMean(t *testing.T) {
	tests := []struct {
		name  string
		setup func() *MovingMean
	}{
		{
			name: "calculating moving average",
			setup: func() *MovingMean {
				return New(3)
			},
		},
	}

	for _, tt := range tests {
		movingMean := tt.setup()
		movingMean.Add(2)
		if !reflect.DeepEqual(movingMean.Average(), float64(2)) {
			t.Errorf("Average()=%v, wanted=%v", movingMean.Average(), 2)
		}
		movingMean.Add(1)
		if !reflect.DeepEqual(movingMean.Average(), 1.5) {
			t.Errorf("Average()=%v, wanted=%v", movingMean.Average(), 1.5)
		}
		movingMean.Add(4)
		if !reflect.DeepEqual(movingMean.Average(), 2.3333333333333335) {
			t.Errorf("Average()=%v, wanted=%v", movingMean.Average(), 2.3333333333333335)
		}
		// As the window is 3, it will eliminate the first value i.e 2
		// and add new element 4 and calculate the average in linear time
		movingMean.Add(4)
		if !reflect.DeepEqual(movingMean.Average(), float64(3)) {
			t.Errorf("Average()=%v, wanted=%v", movingMean.Average(), 3)
		}
	}

}
