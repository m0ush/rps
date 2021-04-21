package trundl

import (
	"sort"
)

const (
	lim1d = -0.01
	lim7d = -0.03
)

type Packet struct {
	Payload []Record `json:"records"`
	DFrame  map[int][]Entry
	Alerts  []Alert
}

func (p *Packet) DataFrame() {
	m := make(map[int][]Entry)
	for _, r := range p.Payload {
		m[r.ID] = append(m[r.ID], r.Entry)
	}
	for _, e := range m {
		sort.Slice(e, func(i, j int) bool { return e[i].Date < e[j].Date })
	}
	p.DFrame = m
}

func (p *Packet) FindAlerts() {
	for k, e := range p.DFrame {
		close := e[4].Price
		day := e[3].Price
		wk := e[0].Price

		// Daily Limit
		if priceReturn(close, day) < lim1d {
			a := NewAlert(k, "daily", lim1d, close, day)
			p.Alerts = append(p.Alerts, a)
		}

		// Weekly Limit
		if priceReturn(close, wk) < lim7d {
			a := NewAlert(k, "weekly", lim1d, close, wk)
			p.Alerts = append(p.Alerts, a)

		}
	}
}

// still likely need alert number (sourced from db),
// security sedol, and name.. Possibly percentage.
type Alert struct {
	ID       int
	Lookback string
	Limit    float64
	Close    float64
	Previous float64
}

func NewAlert(id int, lb string, lim, close, prev float64) Alert {
	return Alert{
		ID:       id,
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

func priceReturn(close, previous float64) float64 {
	return close/previous - 1

}
