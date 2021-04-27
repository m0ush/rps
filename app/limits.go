package app

type Limit interface {
	IsOver([]float64) bool
	Pctg([]float64) float64
	DayLag() int
}

type ReturnLimit struct {
	Val float64
	Lag int
	Fn  func([]float64) float64
}

func (rl ReturnLimit) Pctg(fs []float64) float64 {
	return rl.Fn(fs)
}

func (rl ReturnLimit) DayLag() int {
	return rl.Lag
}

func (rl ReturnLimit) IsOver(fs []float64) bool {
	return rl.Fn(fs) < rl.Val
}

func NewReturnLimit(v float64, i int) ReturnLimit {
	return ReturnLimit{
		Val: v,
		Lag: i,
		Fn:  PeriodReturn(i),
	}
}

func PeriodReturn(i int) func([]float64) float64 {
	return func(fx []float64) float64 {
		return fx[0]/fx[i] - 1
	}
}
