package sysinit

import "github.com/xiliangMa/diss-backend/models"

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
