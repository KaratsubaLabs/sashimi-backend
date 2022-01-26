package db

type Outage struct {
	serviceName string
	outageType  int16
	outageStart uint64
	outageEnd   uint64
}
