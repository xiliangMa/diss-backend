package sysinit

import (
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/service/task"
	"log"
)

func InitGlobalLogConfig() {
	var logConfig models.LogConfig
	logConfig.ConfigName = models.Log_Config_SysLog_Export
	syslogConfig := logConfig.InnerGet()

	if syslogConfig != nil {
		models.GlobalLogConfig[models.Log_Config_SysLog_Export] = syslogConfig[0]
	} else {
		models.GlobalLogConfig[models.Log_Config_SysLog_Export] = &logConfig
	}
}

func InitTimeEdgePoint() {

	syslogHandler := task.GlobalSysLogTaskHandler

	for exType, exGroup := range syslogHandler.ExportTypes {
		log.Printf("extype: %s, exGroup: %#v\n", exType, exGroup)
	}

	TEPoint := new(models.TimeEdgePoint)
	//初始起始时间设置为2018-1-1
	TEPoint.TimePointA = "2018-01-01T0:00:00Z"
	TEPoint.EdgePointCode = ""
}
