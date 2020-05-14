package system

import (
	"fmt"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
	"log"
	"log/syslog"
	"strconv"
	"strings"
)

func OpenSyslog(tag string) *syslog.Writer {

	syslogServer := models.GetSyslogServerUrl()
	//log.Println("syslogServer: ", syslogServer)
	sysLog, err := syslog.Dial("tcp", syslogServer,
		syslog.LOG_WARNING, tag)
	if err != nil {
		log.Println(err)
	}
	return sysLog
}

func SendSysLog(tag, level, msg string) {
	// 没有启用syslog日志导出，直接退出
	syslogConfig := models.GlobalLogConfig[models.Log_Config_SysLog_Export]
	if syslogConfig.Enabled == false {
		return
	}
	// 打开syslog连接
	sysLog := OpenSyslog(tag)
	if sysLog == nil {
		log.Println("Error : Cant connect syslog Server . ErrorCode: " + strconv.Itoa(utils.ConnectSyslogErr))
		return
	}

	if strings.Contains(syslogConfig.ExportedTypes, tag) {
		var err error
		//sysLog.Emerg("Emerg messsage ---------")
		//fmt.Fprintf(sysLog, "Level %s", level)
		msgWithLevel := fmt.Sprintf("[%s] %s", level, msg)
		switch level {
		case models.Log_level_Warn:
			err = sysLog.Warning(msgWithLevel)
		case models.Log_level_Info:
			err = sysLog.Info(msgWithLevel)
		case models.Log_level_Error:
			err = sysLog.Err(msgWithLevel)
		case models.Log_level_Debug:
			err = sysLog.Debug(msgWithLevel)
		}
		if err != nil {
			log.Println("err ", err)
		}
	}
}

func GetSyncSyslogFunc(exType string) func() {
	return func() {
		fmt.Println("Sync log data , type:", exType)
	}
}
