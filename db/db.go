package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {

	sourceString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("SASHIMI_DB_HOST"),
		os.Getenv("SASHIMI_DB_PORT"),
		os.Getenv("SASHIMI_DB_USER"),
		os.Getenv("SASHIMI_DB_PASSWORD"),
		os.Getenv("SASHIMI_DB_NAME"),
	)

	db, err := sql.Open("postgres", sourceString)
	if err != nil {
		return nil, err
	} else {
		return db, nil
	}

}

func LogOutage(db *sql.DB, o Outage) error {

	_, err := db.Query(`SELECT start FROM OutageList WHERE Service=$1, End=0`, o.serviceName)

	if err == sql.ErrNoRows {

		_, _ = db.Exec(`UPDATE ServiceList SET OK=0 WHERE ServiceName=$1`, o.serviceName)
		_, err = db.Exec(
			`INSERT INTO OutageList VALUE Service=$1, Start=$2, End=0`, o.serviceName, o.outageStart,
		)

		return err

	} else {
		return nil
	}

}

func LogOK(db *sql.DB, o Outage) error {

	_, err := db.Query(`SELECT start FROM OutageList WHERE Service=$1, End=0`, o.serviceName)

	if err != sql.ErrNoRows {

		_, _ = db.Exec(`UPDATE ServiceList SET OK=1 WHERE ServiceName=$1`, o.serviceName)
		_, err = db.Exec(`UPDATE OutageList SET end=$1 WHERE ServiceName=$2, End=0`, o.outageEnd, o.serviceName)
		return err

	} else {
		return nil
	}

}

func QueryStats(db *sql.DB) (Stats, error) {

	var rows int32
	db.QueryRow(`SELECT COUNT(Name) FROM ServiceList`).Scan(&rows)

	if rows == 0 {
		return Stats{}, nil
	}

	data, err := db.Query(`SELECT Name, OK FROM ServiceList`)

	if err != nil {
		return Stats{}, err
	}

	defer data.Close()
	stats := Stats{
		serviceName: make([]string, rows),
		status:      make([]bool, rows),
	}

	for i := 0; data.Next(); i++ {
		data.Scan(&stats.serviceName[i], &stats.status[i])
	}

	return stats, nil

}

func DetailedStats(db *sql.DB, name string) (DetailStats, error) {

	var status bool
	var since int64
	if db.QueryRow(`SELECT OK, Since FROM ServiceList WHERE Name=$1`, name).Scan(&status, &since) == sql.ErrNoRows {
		return DetailStats{}, sql.ErrNoRows
	}

	var rows int32
	db.QueryRow(`SELECT COUNT(Name) FROM OutageList WHERE Name=$1`, name).Scan(&rows)

	data, err := db.Query(`SELECT Start, End FROM OutageList`)

	if err != nil {
		return DetailStats{}, err
	}

	defer data.Close()
	detail := DetailStats{
		status:        status,
		timeMonitored: time.Now().Unix() - since,
		downtime:      0,
		outageStart:   make([]int64, rows),
		outageEnd:     make([]int64, rows),
	}

	for i := 0; data.Next(); i++ {

		data.Scan(&detail.outageStart[i], &detail.outageEnd[i])

		if detail.outageEnd[i] != 0 {
			detail.downtime += detail.outageEnd[i] - detail.outageStart[i]
		} else {
			detail.downtime += time.Now().Unix() - detail.outageStart[i]
		}

	}

	return detail, nil
}
