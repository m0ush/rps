package app

import (
	"log"

	"github.com/m0ush/rps/rpsdb"
	"github.com/m0ush/rps/trundl"
)

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
				log.Printf("%9d exceeded limit %d: %.2f%%\n", k, l.DayLag(), l.Pctg(fx)*100)
				a := NewAlert(k, l)
				ax = append(ax, a)
			}
		}
	}
	return ax
}

func AddAlerts(as []Alert) error {
	db, err := rpsdb.Open("sqlite", "./rpsdb/rps.db")
	if err != nil {
		return err
	}
	defer db.Close()

	for _, a := range as {
		// TODO: Replace DayLag with trigger_id!
		_, err := db.InsertAlert(a.SecID, a.DayLag())
		if err != nil {
			return err
		}
	}
	return nil
}

func RegisterLimits() ([]Limit, error) {
	db, err := rpsdb.Open("sqlite", "./rpsdb/rps.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	dls, err := db.AllDLimits()
	if err != nil {
		return nil, err
	}

	var lims []Limit
	for _, dl := range dls {
		rl := NewReturnLimit(dl.Thresh(), dl.DayLag)
		lims = append(lims, rl)
	}
	return lims, nil
}
