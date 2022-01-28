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
	ServiceName []string
	ServiceURL  []string
	Status      []bool
}

type DetailStats struct {
	Status        bool
	TimeMonitored int64
	Downtime      int64
	OutageStart   []int64
	OutageEnd     []int64
}
