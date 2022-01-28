package db

type Outage struct {
	ServiceName string `sql:"Service"`
	OutageStart int64  `sql:"Start"`
	OutageEnd   int64  `sql:"End"`
}

type Service struct {
	ServiceName     string `sql:"Name"`
	ServiceURL      string `sql:"URL"`
	Status          bool   `sql:"OK"`
	MonitoringSince int64  `sql:"Since"`
}

type Stats struct {
	ServiceName []string `json:"service_name"`
	ServiceURL  []string `json:"service_url"`
	Status      []bool   `json:"status"`
}

type DetailStats struct {
	Status        bool    `json:"status"`
	TimeMonitored int64   `json:"time_monitored"`
	Downtime      int64   `json:"downtime"`
	OutageStart   []int64 `json:"outage_start"`
	OutageEnd     []int64 `json:"outage_end"`
}
