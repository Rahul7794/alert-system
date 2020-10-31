package moving_mean

type MovingMean struct {
	size   int
	values Queue
	sum    float64
	count  int
}

func (mm *MovingMean) Average() float64 {
	avg := mm.sum / float64(mm.count)
	return avg
}

func (mm *MovingMean) Add(val float64) {
	mm.sum += val
	mm.values.Enqueue(val)
	if mm.values.Size() > mm.size {
		mm.sum -= mm.values.Dequeue()
	} else {
		mm.count = mm.count + 1
	}
}

func New(window int) *MovingMean {
	return &MovingMean{
		size:   window,
		values: Queue{},
		sum:    0,
		count:  0,
	}
}
