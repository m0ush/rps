package rpsdb

import (
	"database/sql"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

type Security struct {
	SecurityID int
	AddedOn    string
	EndedOn    string
}

type Limit struct {
	LimID  int
	Val    float64
	DayLag int
}

type Alert struct {
	SecurityID int
	Limit
}

func (l Limit) Thresh() float64 {
	return l.Val / 100
}

func (l Limit) Pctg(fs []float64) float64 {
	return fs[0]/fs[l.DayLag] - 1
}

func (l Limit) IsOver(fs []float64) bool {
	return l.Pctg(fs) < l.Thresh()
}

func ResetDB(db *sql.DB) error {
	if _, err := db.Exec("DROP TABLE alerts"); err != nil {
		return err
	}

	if _, err := db.Exec("DROP TABLE limits"); err != nil {
		return err
	}

	if _, err := db.Exec("DROP TABLE securities"); err != nil {
		return err
	}

	return nil
}

func CreateTableSecurities(db *sql.DB) error {
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS securities (
		security_id INTEGER PRIMARY KEY,
		added_on TEXT DEFAULT CURRENT_DATE,
		ended_on TEXT
	)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(); err != nil {
		return err
	}
	return nil
}

func CreateTableLimits(db *sql.DB) error {
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS limits (
		limit_id INTEGER PRIMARY KEY,
		value INTEGER NOT NULL,
		name TEXT NOT NULL,
		day_lag INTEGER
	)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(); err != nil {
		return err
	}
	return nil
}

func CreateTableAlerts(db *sql.DB) error {
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS alerts (
		alert_id INTEGER PRIMARY KEY,
		created_on TEXT DEFAULT CURRENT_DATE,
		security_id INTEGER,
		limit_id INTEGER,
		UNIQUE (created_on, security_id, limit_id),
		FOREIGN KEY (security_id) REFERENCES securities (security_id),
		FOREIGN KEY (limit_id) REFERENCES thresholds (limit_id)
	)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(); err != nil {
		return err
	}
	return nil
}

// TODO: Use CreateTable[TableName] with PRAGMA foreign_keys = ON;
// Create using Transactions
func CreateTables(db *sql.DB) error {
	return nil
}

// TODO: Complete Function
func CreateDatabase(db *sql.DB) error {
	return nil
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

func (db *DB) InsertAlert(a Alert) (int, error) {
	stmt := `INSERT INTO alerts (created_on, security_id, limit_id) VALUES (?, ?, ?) RETURNING alert_id`
	dstr := time.Now().Format("2006-01-02")

	var id int
	if err := db.db.QueryRow(stmt, dstr, a.SecurityID, a.LimID).Scan(&id); err != nil {
		return -1, err
	}
	log.Printf("alert %5d: %9d added with %2d limit_id\n", id, a.SecurityID, a.LimID)

	return id, nil
}

func (db *DB) InsertAlerts(as []Alert) error {
	for _, a := range as {
		if _, err := db.InsertAlert(a); err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) AllLimits() ([]Limit, error) {
	stmt := `SELECT limit_id, value, day_lag FROM limits`

	rows, err := db.db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ls []Limit
	for rows.Next() {
		var l Limit
		if err := rows.Scan(&l.LimID, &l.Val, &l.DayLag); err != nil {
			return nil, err
		}
		ls = append(ls, l)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ls, nil
}

func (db *DB) SeedSecurities() error {
	sx := []Security{
		{
			SecurityID: 4320,
			AddedOn:    "2020-11-11",
			EndedOn:    "2021-01-19",
		},
		{
			SecurityID: 5781,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 7160,
			AddedOn:    "2020-11-11",
			EndedOn:    "2021-01-19",
		},
		{
			SecurityID: 7759,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 8062,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 12356,
			AddedOn:    "2020-11-11",
			EndedOn:    "2021-01-19",
		},
		{
			SecurityID: 14284,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 25848,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 32909,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 43434,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 44699,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 45875,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 46182,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 48402,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 55870,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 61146,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 63891,
			AddedOn:    "2021-01-19",
		},
		{
			SecurityID: 71660,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 71661,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 73919,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 79329,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 80219,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 82354,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 83241,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 91934,
			AddedOn:    "2020-11-11",
			EndedOn:    "2021-01-19",
		},
		{
			SecurityID: 95726,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 100618,
			AddedOn:    "2020-11-11",
			EndedOn:    "2021-01-19",
		},
		{
			SecurityID: 101772,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 109974,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 110748,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 117059,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 121635,
			AddedOn:    "2020-11-11",
			EndedOn:    "2021-01-19",
		},
		{
			SecurityID: 123498,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 123890,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 124774,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 127963,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 133453,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 135602,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 138287,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 138691,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 142659,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 149346,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 151861,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 152758,
			AddedOn:    "2020-11-11",
			EndedOn:    "2021-01-19",
		},
		{
			SecurityID: 160462,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 162865,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 170398,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 171570,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 172801,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 175483,
			AddedOn:    "2020-11-11",
			EndedOn:    "2021-01-19",
		},
		{
			SecurityID: 176724,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 176764,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 183226,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 185494,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 185565,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 185566,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 185570,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 186086,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 187318,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 188801,
			AddedOn:    "2020-11-11",
			EndedOn:    "2021-01-19",
		},
		{
			SecurityID: 189361,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 192071,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 197825,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 228616,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 802200,
			AddedOn:    "2020-11-11",
			EndedOn:    "2021-01-19",
		},
		{
			SecurityID: 55574572,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 55574600,
			AddedOn:    "2020-11-11",
		},
		{
			SecurityID: 55579755,
			AddedOn:    "2020-11-11",
			EndedOn:    "2021-01-19",
		},
		{
			SecurityID: 55580129,
			AddedOn:    "2020-11-11",
			EndedOn:    "2021-01-19",
		},
		{
			SecurityID: 55581035,
			AddedOn:    "2020-11-12",
			EndedOn:    "2021-02-03",
		},
	}

	stmt, err := db.db.Prepare(`INSERT INTO securities(security_id, added_on, ended_on) VALUES (?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, s := range sx {
		if _, err := stmt.Exec(s.SecurityID, s.AddedOn, s.EndedOn); err != nil {
			return err
		}
	}

	return nil
}

func (db *DB) SeedLimits() error {
	ls := []struct {
		LimID  int
		Val    int
		Name   string
		DayLag int
	}{
		{
			LimID:  1,
			Val:    -30,
			Name:   "yesterday",
			DayLag: 1,
		},
		{
			LimID:  2,
			Val:    -50,
			Name:   "last_week",
			DayLag: 4,
		},
	}

	stmt, err := db.db.Prepare(`INSERT INTO limits(limit_id, value, name, day_lag) VALUES (?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, l := range ls {
		if _, err := stmt.Exec(l.LimID, l.Val, l.Name, l.DayLag); err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) SeedAlerts() error {
	as := []struct {
		SecurityID int
		CreateDate string
		LimitID    int
	}{
		{
			SecurityID: 138287,
			CreateDate: "2019-10-25",
			LimitID:    1,
		},
		{
			SecurityID: 185566,
			CreateDate: "2020-03-10",
			LimitID:    1,
		},
		{
			SecurityID: 171570,
			CreateDate: "2020-03-10",
			LimitID:    1,
		},
		{
			SecurityID: 186086,
			CreateDate: "2020-03-17",
			LimitID:    1,
		},
		{
			SecurityID: 46182,
			CreateDate: "2020-03-19",
			LimitID:    1,
		},
		{
			SecurityID: 133453,
			CreateDate: "2020-03-19",
			LimitID:    1,
		},
		{
			SecurityID: 55579755,
			CreateDate: "2019-11-26",
			LimitID:    1,
		},
		{
			SecurityID: 138287,
			CreateDate: "2019-10-29",
			LimitID:    2,
		},
		{
			SecurityID: 192071,
			CreateDate: "2020-02-26",
			LimitID:    2,
		},
		{
			SecurityID: 185566,
			CreateDate: "2020-03-09",
			LimitID:    2,
		},
		{
			SecurityID: 171570,
			CreateDate: "2020-03-10",
			LimitID:    2,
		},
		{
			SecurityID: 186086,
			CreateDate: "2020-03-10",
			LimitID:    2,
		},
		{
			SecurityID: 79329,
			CreateDate: "2020-03-12",
			LimitID:    2,
		},
		{
			SecurityID: 133453,
			CreateDate: "2020-03-13",
			LimitID:    2,
		},
		{
			SecurityID: 197825,
			CreateDate: "2020-03-13",
			LimitID:    2,
		},
		{
			SecurityID: 8062,
			CreateDate: "2020-03-13",
			LimitID:    2,
		},
		{
			SecurityID: 185494,
			CreateDate: "2020-03-13",
			LimitID:    2,
		},
		{
			SecurityID: 187318,
			CreateDate: "2020-03-13",
			LimitID:    2,
		},
		{
			SecurityID: 172801,
			CreateDate: "2020-03-17",
			LimitID:    2,
		},
		{
			SecurityID: 12356,
			CreateDate: "2020-03-17",
			LimitID:    2,
		},
		{
			SecurityID: 149346,
			CreateDate: "2020-03-17",
			LimitID:    2,
		},
		{
			SecurityID: 46182,
			CreateDate: "2020-03-17",
			LimitID:    2,
		},
		{
			SecurityID: 71661,
			CreateDate: "2020-03-17",
			LimitID:    2,
		},
		{
			SecurityID: 124774,
			CreateDate: "2020-03-17",
			LimitID:    2,
		},
		{
			SecurityID: 61146,
			CreateDate: "2020-03-17",
			LimitID:    2,
		},
		{
			SecurityID: 189361,
			CreateDate: "2020-03-19",
			LimitID:    2,
		},
		{
			SecurityID: 185565,
			CreateDate: "2020-03-19",
			LimitID:    2,
		},
		{
			SecurityID: 160462,
			CreateDate: "2020-03-20",
			LimitID:    2,
		},
		{
			SecurityID: 127963,
			CreateDate: "2020-03-24",
			LimitID:    2,
		},
		{
			SecurityID: 80219,
			CreateDate: "2020-04-02",
			LimitID:    2,
		},
		{
			SecurityID: 55579755,
			CreateDate: "2019-12-03",
			LimitID:    2,
		},
		{
			SecurityID: 55581035,
			CreateDate: "2020-03-19",
			LimitID:    1,
		},
		{
			SecurityID: 55581035,
			CreateDate: "2020-02-13",
			LimitID:    2,
		},
	}

	stmt, err := db.db.Prepare(`INSERT INTO alerts(security_id, created_on, limit_id) VALUES (?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, a := range as {
		if _, err := stmt.Exec(a.SecurityID, a.CreateDate, a.LimitID); err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) Seed() error {
	if err := db.SeedSecurities(); err != nil {
		return err
	}

	if err := db.SeedLimits(); err != nil {
		return err
	}

	if err := db.SeedAlerts(); err != nil {
		return err
	}
	return nil
}
