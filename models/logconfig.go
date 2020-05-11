package models

type LogConfig struct {
	Id            string `orm:"pk;" description:"(Log配置id)"`
	Enabled       bool   `orm:"" description:"(是否启用)"`
	ServerUrl     string `orm:"" description:"(服务器url)"`
	ServerPort    string `orm:"" description:"(服务器端口)"`
	ExportedTypes string `orm:"" description:"(导出日志类型)"` //日志类型的枚举，多个，以, 分割
}
