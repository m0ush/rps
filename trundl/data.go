package trundl

import (
	"sort"
)

// TODO: pull these variables from database
const (
	lim1d = -0.01
	lim7d = -0.03
)

type Records struct {
	Rx []Record `json:"records"`
	N  int      `json:"total"`
}

func DataFrame(rs Records) map[int][]Entry {
	df := make(map[int][]Entry, rs.N)
	for _, r := range rs.Rx {
		df[r.ID] = append(df[r.ID], r.Entry)
	}
	return df
}

// TODO: Make Parallel
func Sorter(df map[int][]Entry) map[int][]Entry {
	for _, e := range df {
		sort.Slice(e, func(i, j int) bool { return e[i].Date < e[j].Date })
	}
	return df
}

// TODO: Make Parallel
func Series(df map[int][]Entry) map[int][]float64 {
	series := make(map[int][]float64, len(df))
	for id, data := range df {
		var fx []float64
		for _, e := range data {
			fx = append(fx, e.Price)
		}
		series[id] = fx
	}
	return series
}

// TODO: Make Process Paralell
// TODO: Accept Functional Args for Based on Different Limits
func FindAlerts(series map[int][]float64) []Alert {
	var ax []Alert
	for k, fx := range series {
		// Daily Limit
		if priceReturn(fx[4], fx[3]) < lim1d {
			a := NewAlert(k, "daily", lim1d, fx[4], fx[3])
			ax = append(ax, a)
		}
		// Weekly Limit
		if priceReturn(fx[4], fx[0]) < lim7d {
			a := NewAlert(k, "weekly", lim7d, fx[4], fx[0])
			ax = append(ax, a)
		}
	}
	return ax
}

func Pipeline(rs Records) []Alert {
	return FindAlerts(Series(Sorter(DataFrame(rs))))
}

// TODO: Include Secruity Attributes (e.g. Sedol, Name)
// TODO: Include Limit Struct as an Embedded Field
type Alert struct {
	SecID    int
	Lookback string
	Limit    float64
	Close    float64
	Previous float64
}

func NewAlert(id int, lb string, lim, close, prev float64) Alert {
	return Alert{
		SecID:    id,
		Lookback: lb,
		Limit:    lim,
		Close:    close,
		Previous: prev,
	}
}

type Record struct {
	ID    int `json:"id,string"`
	Entry `json:"data"`
}

type Entry struct {
	Date  string  `json:"price_effective_date"`
	Price float64 `json:"price_local_amount,string"`
}

// TODO: Handle ZeroDivisorError
func priceReturn(close, previous float64) float64 {
	return close/previous - 1
}
