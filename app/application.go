package app

import (
	"log"

	"github.com/m0ush/rps/rpsdb"
	"github.com/m0ush/rps/trundl"
)

func Sieve(rs trundl.Records, ls ...rpsdb.Limit) []rpsdb.Alert {
	return Evaluator(trundl.Pipeline(rs), ls...)
}

// TODO: Make Process Paralell
// TODO: Accept Functional Args for Based on Different Limits
func Evaluator(series map[int][]float64, ls ...rpsdb.Limit) []rpsdb.Alert {
	var ax []rpsdb.Alert
	for k, fx := range series {
		for _, l := range ls {
			if l.IsOver(fx) {
				log.Printf("%9d exceeded limit %d: %.2f%%\n", k, l.DayLag, l.Pctg(fx)*100)
				ax = append(ax, rpsdb.Alert{
					SecurityID: k,
					Limit:      l,
				})
			}
		}
	}
	return ax
}
