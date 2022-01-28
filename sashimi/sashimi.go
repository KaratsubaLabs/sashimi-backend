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
		pingRoutes(serviceList, t.Unix())
	}

	return nil

}

func pingRoutes(list db.Stats, time int64) {
	for i := 0; i < len(list.ServiceName); i++ {
		go ping(list.ServiceName[i], list.ServiceURL[i], time)
	}
}

func ping(name string, url string, time int64) {

	resp, err := http.Get(url)

	if err == nil && resp.StatusCode == 200 {

		// Log OK
		Database.LogOK(
			db.Outage{
				ServiceName: name,
				OutageStart: 0,
				OutageEnd:   time,
			},
		)

	} else {

		// Log Outage
		Database.LogOK(
			db.Outage{
				ServiceName: name,
				OutageStart: time,
				OutageEnd:   0,
			},
		)

	}

}
