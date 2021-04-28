package app

import (
	"encoding/json"
	"io"
	"io/ioutil"
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

func Runner(r io.Reader) error {
	bs, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	var rx trundl.Records
	if err := json.Unmarshal(bs, &rx); err != nil {
		return err
	}

	db, err := rpsdb.Open("sqlite", "./rpsdb/rps.db")
	if err != nil {
		return err
	}
	defer db.Close()

	lims, err := db.AllLimits()
	if err != nil {
		return err
	}
	ax := Sieve(rx, lims...)
	db.InsertAlerts(ax)

	return nil
}
