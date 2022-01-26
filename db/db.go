package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

const OutageLogDBName = "downtimelog"
const SitesDBName = "sitestomonitor"

func Connect(dbname string) (*sql.DB, error) {
	sourceString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("SASHIMI_DB_HOST"),
		os.Getenv("SASHIMI_DB_PORT"),
		os.Getenv("SASHIMI_DB_USER"),
		os.Getenv("SASHIMI_DB_PASSWORD"),
		dbname,
	)

	db, err := sql.Open("postgres", sourceString)
	if err != nil {
		return nil, err
	} else {
		return db, nil
	}
}

func LogOutage(db *sql.DB, o Outage) error {
	command := `SELECT start FROM $1 WHERE service=$2, end=0`
	_, exists := db.Query(command, OutageLogDBName, o.serviceName)

	if exists == sql.ErrNoRows {
		command = `INSERT INTO $1 VALUE service=$2, type=$3, start=$4, end=0`
		_, err := db.Exec(command, OutageLogDBName, o.serviceName, o.outageType, o.outageStart)

		if err == nil {
			return err
		} else {
			return nil
		}
	} else {
		return nil
	}
}

func LogOK(db *sql.DB, o Outage) error {
	command := `SELECT start FROM $1 WHERE service=$2, end=0`
	outageStart, exists := db.Query(command, OutageLogDBName, o.serviceName)

	if exists != sql.ErrNoRows {
		command = `UPDATE $1 SET end=$2 WHERE start=$3`
		_, err := db.Exec(command, OutageLogDBName, o.outageEnd, outageStart)

		if err == nil {
			return err
		} else {
			return nil
		}
	} else {
		return nil
	}
}
