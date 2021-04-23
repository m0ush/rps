package app

type Limit interface {
	IsOver([]float64) bool
}

type ReturnLimit struct {
	Val float64
	Fn  func([]float64) float64
}

func (rl ReturnLimit) IsOver(fs []float64) bool {
	return rl.Fn(fs) < rl.Val
}

func NewReturnLimit(v float64, i int) ReturnLimit {
	return ReturnLimit{
		Val: v,
		Fn:  PeriodReturn(i),
	}
}

func PeriodReturn(i int) func([]float64) float64 {
	return func(fx []float64) float64 {
		return fx[0]/fx[i] - 1
	}
}
