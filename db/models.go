package db

type Outage struct {
	serviceName string `sql:"Service"`
	outageStart int64  `sql:"Start"`
	outageEnd   int64  `sql:"End"`
}

type Service struct {
	serviceName     string `sql:"Name"`
	serviceURL      string `sql:"URL"`
	status          bool   `sql:"OK"`
	monitoringSince int64  `sql:"Since"`
}

type Stats struct {
	serviceName []string
	status      []bool
}

type DetailStats struct {
	status        bool
	timeMonitored int64
	downtime      int64
	outageStart   []int64
	outageEnd     []int64
}
