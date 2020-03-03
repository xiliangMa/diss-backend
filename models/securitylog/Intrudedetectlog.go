package securitylog

type IntrudeDetectLog struct {
	HostId      string
	TargeType   string
	ContainerId string
	StartTime   string
	ToTime      string
	Limit       int
}
