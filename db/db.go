package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Connection struct {
	db        *sql.DB
	connected bool
}

func (r *Connection) Connect() error {

	sourceString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("SASHIMI_DB_HOST"),
		os.Getenv("SASHIMI_DB_PORT"),
		os.Getenv("SASHIMI_DB_USER"),
		os.Getenv("SASHIMI_DB_PASSWORD"),
		os.Getenv("SASHIMI_DB_NAME"),
	)

	var err error
	r.db, err = sql.Open("postgres", sourceString)

	if err == nil {
		r.connected = true
	}

	return err

}

func (r *Connection) Migrate() error {

	_, err1 := r.db.Exec(`CREATE TABLE ServiceList()`)
	_, err2 := r.db.Exec(`CREATE TABLE OutageList()`)

	if err1 == sql.ErrConnDone || err2 == sql.ErrConnDone {
		return sql.ErrConnDone
	}

	r.db.Exec(`ALTER TABLE ServiceList ADD Name VARCHAR`)
	r.db.Exec(`ALTER TABLE ServiceList ADD URL VARCHAR`)
	r.db.Exec(`ALTER TABLE ServiceList ADD OK BOOL`)
	r.db.Exec(`ALTER TABLE ServiceList ADD Since BIGINT`)

	r.db.Exec(`ALTER TABLE OutageList ADD Service VARCHAR`)
	r.db.Exec(`ALTER TABLE OutageList ADD Start BIGINT`)
	r.db.Exec(`ALTER TABLE OutageList ADD End BIGINT`)

	return nil

}

func (r *Connection) AddService(s Service) error {

	_, err := r.db.Exec(
		`INSERT INTO ServiceList VALUE Name=$1, URL=$2, OK=1, Since=$3`, s.serviceName, s.serviceURL, time.Now().Unix(),
	)

	return err

}

func (r *Connection) RemoveService(s Service) error {

	_, err := r.db.Exec(`DELETE FROM ServiceList WHERE Name=$1`, s.serviceName)

	if err != nil {
		_, err = r.db.Exec(`DELETE FROM ServiceList WHERE URL=$1`, s.serviceURL)
	}

	return err

}

func (r *Connection) LogOutage(o Outage) error {

	_, err := r.db.Query(`SELECT start FROM OutageList WHERE Service=$1, End=0`, o.serviceName)

	if err == sql.ErrNoRows {

		_, _ = r.db.Exec(`UPDATE ServiceList SET OK=0 WHERE ServiceName=$1`, o.serviceName)
		_, err = r.db.Exec(
			`INSERT INTO OutageList VALUE Service=$1, Start=$2, End=0`, o.serviceName, o.outageStart,
		)

		return err

	} else {
		return nil
	}

}

func (r *Connection) LogOK(o Outage) error {

	_, err := r.db.Query(`SELECT start FROM OutageList WHERE Service=$1, End=0`, o.serviceName)

	if err != sql.ErrNoRows {

		_, _ = r.db.Exec(`UPDATE ServiceList SET OK=1 WHERE ServiceName=$1`, o.serviceName)
		_, err = r.db.Exec(`UPDATE OutageList SET end=$1 WHERE ServiceName=$2, End=0`, o.outageEnd, o.serviceName)
		return err

	} else {
		return nil
	}

}

func (r *Connection) QueryStats() (Stats, error) {

	var rows int32
	r.db.QueryRow(`SELECT COUNT(Name) FROM ServiceList`).Scan(&rows)

	if rows == 0 {
		return Stats{}, nil
	}

	data, err := r.db.Query(`SELECT Name, OK FROM ServiceList`)

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

func (r *Connection) DetailedStats(name string) (DetailStats, error) {

	var status bool
	var since int64
	if r.db.QueryRow(`SELECT OK, Since FROM ServiceList WHERE Name=$1`, name).Scan(&status, &since) == sql.ErrNoRows {
		return DetailStats{}, sql.ErrNoRows
	}

	var rows int32
	r.db.QueryRow(`SELECT COUNT(Name) FROM OutageList WHERE Name=$1`, name).Scan(&rows)

	data, err := r.db.Query(`SELECT Start, End FROM OutageList`)

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
