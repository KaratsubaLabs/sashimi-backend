package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Connection struct {
	DB        *sql.DB
	Connected bool
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
	r.DB, err = sql.Open("postgres", sourceString)

	if err == nil {
		r.Connected = true
	}

	return err

}

func (r *Connection) Migrate() error {

	_, err1 := r.DB.Exec(`CREATE TABLE ServiceList()`)
	_, err2 := r.DB.Exec(`CREATE TABLE OutageList()`)

	if err1 == sql.ErrConnDone || err2 == sql.ErrConnDone {
		return sql.ErrConnDone
	}

	r.DB.Exec(`ALTER TABLE ServiceList ADD Name VARCHAR`)
	r.DB.Exec(`ALTER TABLE ServiceList ADD URL VARCHAR`)
	r.DB.Exec(`ALTER TABLE ServiceList ADD OK BOOL`)
	r.DB.Exec(`ALTER TABLE ServiceList ADD Since BIGINT`)

	r.DB.Exec(`ALTER TABLE OutageList ADD Service VARCHAR`)
	r.DB.Exec(`ALTER TABLE OutageList ADD Start BIGINT`)
	r.DB.Exec(`ALTER TABLE OutageList ADD End BIGINT`)

	return nil

}

func (r *Connection) AddService(s Service) error {

	_, err := r.DB.Exec(
		`INSERT INTO ServiceList VALUE Name=$1, URL=$2, OK=1, Since=$3`, s.ServiceName, s.ServiceURL, time.Now().Unix(),
	)

	return err

}

func (r *Connection) RemoveService(s Service) error {

	_, err := r.DB.Exec(`DELETE FROM ServiceList WHERE Name=$1`, s.ServiceName)

	if err != nil {
		_, err = r.DB.Exec(`DELETE FROM ServiceList WHERE URL=$1`, s.ServiceURL)
	}

	return err

}

func (r *Connection) LogOutage(o Outage) error {

	// Check for ongoing outage
	_, err := r.DB.Query(`SELECT start FROM OutageList WHERE Service=$1, End=0`, o.ServiceName)

	if err == sql.ErrNoRows {

		// Log Outage if no ongoing outage exists
		_, _ = r.DB.Exec(`UPDATE ServiceList SET OK=0 WHERE ServiceName=$1`, o.ServiceName)
		_, err = r.DB.Exec(
			`INSERT INTO OutageList VALUE Service=$1, Start=$2, End=0`, o.ServiceName, o.OutageStart,
		)

		return err

	} else {
		return nil
	}

}

func (r *Connection) LogOK(o Outage) error {

	// Check for ongoing outage
	_, err := r.DB.Query(`SELECT start FROM OutageList WHERE Service=$1, End=0`, o.ServiceName)

	if err != sql.ErrNoRows {

		// Log OK if ongoing outage exists
		_, _ = r.DB.Exec(`UPDATE ServiceList SET OK=1 WHERE ServiceName=$1`, o.ServiceName)
		_, err = r.DB.Exec(`UPDATE OutageList SET end=$1 WHERE ServiceName=$2, End=0`, o.OutageEnd, o.ServiceName)
		return err

	} else {
		return nil
	}

}

func (r *Connection) GetStats() (Stats, error) {

	var rows int32
	r.DB.QueryRow(`SELECT COUNT(Name) FROM ServiceList`).Scan(&rows)

	if rows == 0 {
		return Stats{}, nil
	}

	data, err := r.DB.Query(`SELECT Name, URL, OK FROM ServiceList`)

	if err != nil {
		return Stats{}, err
	}

	defer data.Close()
	stats := Stats{
		ServiceName: make([]string, rows),
		ServiceURL:  make([]string, rows),
		Status:      make([]bool, rows),
	}

	for i := 0; data.Next(); i++ {
		data.Scan(&stats.ServiceName[i], &stats.ServiceURL[i], &stats.Status[i])
	}

	return stats, nil

}

func (r *Connection) GetDetails(name string) (DetailStats, error) {

	var Status bool
	var since int64
	if r.DB.QueryRow(`SELECT OK, Since FROM ServiceList WHERE Name=$1`, name).Scan(&Status, &since) == sql.ErrNoRows {
		return DetailStats{}, sql.ErrNoRows
	}

	var rows int32
	r.DB.QueryRow(`SELECT COUNT(Name) FROM OutageList WHERE Name=$1`, name).Scan(&rows)

	data, err := r.DB.Query(`SELECT Start, End FROM OutageList`)

	if err != nil {
		return DetailStats{}, err
	}

	defer data.Close()
	detail := DetailStats{
		Status:        Status,
		TimeMonitored: time.Now().Unix() - since,
		Downtime:      0,
		OutageStart:   make([]int64, rows),
		OutageEnd:     make([]int64, rows),
	}

	for i := 0; data.Next(); i++ {

		data.Scan(&detail.OutageStart[i], &detail.OutageEnd[i])

		if detail.OutageEnd[i] != 0 {
			detail.Downtime += detail.OutageEnd[i] - detail.OutageStart[i]
		} else {
			detail.Downtime += time.Now().Unix() - detail.OutageStart[i]
		}

	}

	return detail, nil
}
