package moving_mean

// Direction indicates Rise->0, Fall->1 or Flat->2 of rates coming in.
type Direction int

// enum for Direction
const (
	Rise Direction = iota
	Fall
	Flat
)

// A MovingMean keeps track of the rolling average of a currency pair's rates and its change Direction.
type MovingMean struct {
	Size      int       // Size of the window
	Values    Queue     // Values in the queue for the window
	Total     float64   // Total of rates for the current window
	Count     int       // current length of Values, always <= Size
	Direction Direction // last change Direction
	Trend     int       // number of seconds the current direction has not changed
}

// Average calculates the moving average for the window
func (mm *MovingMean) Average() float64 {
	avg := mm.Total / float64(mm.Count)
	return avg
}

// Add add a incoming rate to the the Values for the window
// and calculates Total sum needed to moving Average.
func (mm *MovingMean) Add(newRate float64) {
	mm.Total += newRate
	mm.calculateTrend(newRate)
	mm.Values.Enqueue(newRate)
	if mm.Values.Size() > mm.Size {
		firstValue, _ := mm.Values.Dequeue()
		mm.Total -= *firstValue
	} else {
		mm.Count = mm.Count + 1
	}
}

// calculateTrend check if there is rise/fall of rate by comparing with last rate and setup the Direction
// and increment the Trend by 1 or set Trend to 1 if there is change in Direction.
func (mm *MovingMean) calculateTrend(newRate float64) {
	lastDirection := mm.Direction
	lastRate, err := mm.Values.PeekLast()
	if err == nil {
		// change Direction by comparing newRate and lastRate from Queue
		switch diff := newRate - *lastRate; {
		case diff > 0:
			mm.Direction = Rise
		case diff < 0:
			mm.Direction = Fall
		default:
			mm.Direction = Flat
		}
	}
	if mm.Direction == lastDirection {
		mm.Trend++
	} else {
		mm.Trend = 1
	}
}

// New initializes MovingMean for a given window
func New(window int) *MovingMean {
	return &MovingMean{
		Size:      window,
		Values:    Queue{},
		Total:     0,
		Count:     0,
		Direction: Flat,
		Trend:     0,
	}
}
