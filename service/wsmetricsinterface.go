package service

type WSMetricsinterface interface {
	Save() error
	Update() error
}
