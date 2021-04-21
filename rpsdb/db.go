package rpsdb

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

type DAlert struct {
	ID         int
	CreateDate string
	SecurityID int
	TriggerID  int
}

type DB struct {
	db *sql.DB
}

func ToOpen(driverName, dataSource string) (*DB, error) {
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
