package system

import (
	"fmt"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
	"log"
	"log/syslog"
)

func OpenSyslog(tag string) *syslog.Writer {

	syslogServer := utils.GetSyslogServerUrl()
	sysLog, err := syslog.Dial("tcp", syslogServer,
		syslog.LOG_WARNING, tag)
	if err != nil {
		log.Println(err)
	}
	return sysLog
}

func SendSysLog(tag, level, msg string) {
	sysLog := OpenSyslog(tag)
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
	fmt.Println("err ", err)

}
