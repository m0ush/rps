package rpsdb

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

type DAlert struct {
	AlertID    int
	CreateDate string
	SecurityID int
	TriggerID  int
}

// TODO: Create an Inferface for Limit
// Limits can be defined any way but they
// should always perform similar actions
// TODO: Add Lookback/DayLag period to thresholds table
// Yesterday = 2
// LastWeek = 8
type Limit struct {
	LimID int
	Val   float64
	Name  string
	DayLag int
}

func(l Limit) Thresh() float64 {
	return l.Val / 100
}

type DB struct {
	db *sql.DB
}

func Open(driverName, dataSource string) (*DB, error) {
	db, err := sql.Open(driverName, dataSource)
	if err != nil {
		return nil, err
	}
	return &DB{db}, nil
}

func (db *DB) Close() error {
	return db.db.Close()
}

func (db *DB) AllSecurities() ([]int, error) {
	stmt := `SELECT security_id FROM securities WHERE ended_on IS NULL`

	rows, err := db.db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ids, nil
}

// TODO: Possibly Accept Alert Struct
func (db *DB) InsertAlert(secID, trigger_id int) (int, error) {
	return -1, nil
}

func (db *DB) AllLimits() ([]Limit, error) {
	stmt := `SELECT * FROM thresholds`

	rows, err := db.db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lims []Limit
	for rows.Next() {
		var l Limit
		if err := rows.Scan(&l.LimID, &l.Val, &l.Name); err != nil {
			return nil, err
		}
		lims = append(lims, l)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return lims, nil
}
