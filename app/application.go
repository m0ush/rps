package app

import "github.com/m0ush/rps/trundl"

func Sieve(rs trundl.Records, lims ...Limit) []Alert {
	return Evaluator(trundl.Pipeline(rs), lims...)
}

// TODO: Make Process Paralell
// TODO: Accept Functional Args for Based on Different Limits
func Evaluator(series map[int][]float64, lims ...Limit) []Alert {
	var ax []Alert
	for k, fx := range series {
		for _, l := range lims {
			if l.IsOver(fx) {
				a := NewAlert(k, l)
				ax = append(ax, a)
			}
		}
	}
	return ax
}
