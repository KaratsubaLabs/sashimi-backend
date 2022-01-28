package sashimi

import (
	"net/http"
	"time"

	"github.com/karatsubalabs/sashimi-backend/db"
)

var Database db.Connection

var PingInterval time.Duration = 60

func Start() error {
	Database.Connect()
	Database.Migrate()

	serviceList, err := Database.GetStats()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(PingInterval * time.Second)
	for t := range ticker.C {
		PingRoutes(serviceList, t.Unix())
	}

	return nil

}

func PingRoutes(list db.Stats, time int64) {
	for i := 0; i < len(list.ServiceName); i++ {
		go Ping(list.ServiceName[i], list.ServiceURL[i], time)
	}
}

func Ping(name string, url string, time int64) error {

	resp, err := http.Get(url)

	if err == nil && resp.StatusCode == 200 {

		// Log OK
		err = Database.LogOK(
			db.Outage{
				ServiceName: name,
				OutageStart: 0,
				OutageEnd:   time,
			},
		)

	} else {

		// Log Outage
		err = Database.LogOK(
			db.Outage{
				ServiceName: name,
				OutageStart: time,
				OutageEnd:   0,
			},
		)

	}

	return err

}
