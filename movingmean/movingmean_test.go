package movingmean

import (
	"reflect"
	"testing"
)

func TestMovingMean(t *testing.T) {
	tests := []struct {
		name  string
		setup func() MovingMean
	}{
		{
			name: "calculating moving average",
			setup: func() MovingMean {
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

func TestCalculateTrend(t *testing.T) {
	tests := []struct {
		name      string
		setup     func() MovingMean
		trend     int
		direction Direction
	}{
		{
			name: "successfully calculate trend for continuous increase of spot rate",
			setup: func() MovingMean {
				mm := New(3)
				// Series of spot rate coming in
				mm.Add(2.3)
				mm.Add(3.4)
				mm.Add(4.5)
				return mm
			},
			trend:     2,
			direction: Rise,
		},
		{
			name: "successfully calculate trend for continuous decrease of spot rate",
			setup: func() MovingMean {
				mm := New(3)
				// Series of spot rate coming in
				mm.Add(2.3)
				mm.Add(1.3)
				mm.Add(0.3)
				return mm
			},
			trend:     2,
			direction: Fall,
		},
		{
			name: "spot rate flat",
			setup: func() MovingMean {
				mm := New(3)
				// Series of spot rate coming in
				mm.Add(2)
				mm.Add(2)
				mm.Add(2)
				return mm
			},
			trend:     3,
			direction: Flat,
		},
	}
	for _, tt := range tests {
		mm := tt.setup()
		trend := mm.Trend
		direction := mm.Direction
		if !reflect.DeepEqual(trend, tt.trend) {
			t.Errorf("calculateTrend(newRate)=%v, wanted=%v", trend, tt.trend)
		}
		if !reflect.DeepEqual(direction, tt.direction) {
			t.Errorf("calculateTrend(newRate)=%v, wanted=%v", direction, tt.direction)
		}
	}
}
