package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func Connect(dbname string) (*sql.DB, error) {
	sourceString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=downtimelog sslmode=disable",
		os.Getenv("SASHIMI_DB_HOST"),
		os.Getenv("SASHIMI_DB_PORT"),
		os.Getenv("SASHIMI_DB_USER"),
		os.Getenv("SASHIMI_DB_PASSWORD"),
	)

	db, err := sql.Open("postgres", sourceString)
	if err != nil {
		return nil, err
	} else {
		return db, nil
	}
}

func LogOutage(db *sql.DB, o Outage) error {
	command := `SELECT start FROM downtimelog WHERE end=0`
	_, exists := db.Query(command, o.serviceName, o.outageType, o.outageStart)

	if exists == sql.ErrNoRows {
		command = `INSERT INTO downtimelog VALUE service=$1, type=$2, start=$3, end=0`
		_, err := db.Exec(command, o.serviceName, o.outageType, o.outageStart)

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
	command := `SELECT start FROM downtimelog WHERE type=$1, end=0`
	outageStart, exists := db.Query(command, o.serviceName, o.outageType, o.outageStart)

	if exists != sql.ErrNoRows {
		command = `UPDATE downtimelog SET end=$1 WHERE end=$2`
		_, err := db.Exec(command, outageStart, o.outageEnd)

		if err == nil {
			return err
		} else {
			return nil
		}
	} else {
		return nil
	}
}
