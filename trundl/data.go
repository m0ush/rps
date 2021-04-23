package trundl

import (
	"sort"
)

type Records struct {
	Rx []Record `json:"records"`
	N  int      `json:"total"`
}

type Record struct {
	ID    int `json:"id,string"`
	Entry `json:"data"`
}

type Entry struct {
	Date  string  `json:"price_effective_date"`
	Price float64 `json:"price_local_amount,string"`
}

func Pipeline(rs Records) map[int][]float64 {
	return Series(Sorter(DataFrame(rs)))
}

// TODO: Parallel
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

// TODO: Parallel
func Sorter(df map[int][]Entry) map[int][]Entry {
	for _, e := range df {
		sort.Slice(e, func(i, j int) bool { return e[i].Date > e[j].Date })
	}
	return df
}

func DataFrame(rs Records) map[int][]Entry {
	df := make(map[int][]Entry, rs.N)
	for _, r := range rs.Rx {
		df[r.ID] = append(df[r.ID], r.Entry)
	}
	return df
}
