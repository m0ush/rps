package rpsdb

import (
	"database/sql"
	"log"
	"time"

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
type DLimit struct {
	LimID  int
	Val    float64
	Name   string
	DayLag int
}

func (dl DLimit) Thresh() float64 {
	return dl.Val / 100
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
func (db *DB) InsertAlert(secID, triggerID int) (int, error) {
	stmt := `INSERT INTO alerts (created_on, security_id, trigger_id) VALUES (?, ?, ?) RETURNING alert_id`
	dstr := time.Now().Format("2006-01-02")

	var id int
	if err := db.db.QueryRow(stmt, dstr, secID, triggerID).Scan(&id); err != nil {
		return -1, err
	}
	log.Printf("alert %5d: %9d added with %2d trigger_id\n", id, secID, triggerID)

	return id, nil
}

func (db *DB) AllDLimits() ([]DLimit, error) {
	stmt := `SELECT * FROM limits`

	rows, err := db.db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dls []DLimit
	for rows.Next() {
		var l DLimit
		if err := rows.Scan(&l.LimID, &l.Val, &l.Name, &l.DayLag); err != nil {
			return nil, err
		}
		dls = append(dls, l)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return dls, nil
}
