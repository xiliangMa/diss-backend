package sysinit

import (
	"github.com/astaxie/beego/logs"
	uuid "github.com/satori/go.uuid"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/service/task"
)

func InitGlobalLogConfig() {
	var logConfig models.LogConfig
	logConfig.ConfigName = models.Log_Config_SysLog_Export
	syslogConfig := logConfig.InnerGet()

	logs.Info("syslog Config ", syslogConfig)
	if len(syslogConfig) > 0 {
		models.GlobalLogConfig[models.Log_Config_SysLog_Export] = syslogConfig[0]
	} else {
		models.GlobalLogConfig[models.Log_Config_SysLog_Export] = &logConfig
	}

	InitTimeEdgePoint()
}

func InitTimeEdgePoint() {

	for exType, exGroup := range *task.AllExGroups() {

		TEPoint := new(models.TimeEdgePoint)
		TEPoint.EdgePointCode = exType
		TEPinDB := TEPoint.Get()
		if len(TEPinDB) == 0 {
			//初始起始时间设置为2018-1-1
			uid, _ := uuid.NewV4()
			TEPoint.TimePointA = "2018-01-01T0:00:00Z"

			TEPoint.EdgePointName = exGroup.ExportName
			TEPoint.TEPointId = uid.String()
			TEPoint.Direction = "lookback"
			TEPoint.ScopeSymbol = "----|"
			TEPoint.Add()
		}
	}
}
