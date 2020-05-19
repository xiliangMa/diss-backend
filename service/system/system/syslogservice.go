package system

import (
	"encoding/json"
	"fmt"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
	"log"
	"log/syslog"
	"strconv"
	"strings"
	"time"
)

type SyslogHandler struct {
	syslog *syslog.Writer
}

var GlobalSyslog = new(SyslogHandler)

func (this SyslogHandler) OpenSyslog(tag string) error {

	syslogServer := models.GetSyslogServerUrl()
	//log.Println("syslogServer: ", syslogServer)
	sysLog, err := syslog.Dial("tcp", syslogServer,
		syslog.LOG_WARNING, tag)
	if err != nil {
		log.Println("ErrorCode: "+strconv.Itoa(utils.ConnectSyslogErr), err)
	}

	this.syslog = sysLog
	return err
}

func (this SyslogHandler) SendSysLog(tag, level, msg string) {
	// 没有启用syslog日志导出，直接退出
	syslogConfig := models.GlobalLogConfig[models.Log_Config_SysLog_Export]
	if syslogConfig.Enabled == false {
		return
	}

	if strings.Contains(syslogConfig.ExportedTypes, tag) {
		var err error

		msgWithLevel := fmt.Sprintf("[%s] %s", level, msg)
		switch level {
		case models.Log_level_Warn:
			err = this.syslog.Warning(msgWithLevel)
		case models.Log_level_Info:
			err = this.syslog.Info(msgWithLevel)
		case models.Log_level_Error:
			err = this.syslog.Err(msgWithLevel)
		case models.Log_level_Debug:
			err = this.syslog.Debug(msgWithLevel)
		}
		if err != nil {
			log.Println("err ", err)
		}
	}
}

func GetSyncSyslogFunc(exType string) func() {
	return func() {
		log.Println("Sync syslog data , type:", exType)

		//default from and limit , for 3000 records
		from := 0
		limit := 3000

		GlobalSyslog.OpenSyslog("init synclog")
		if GlobalSyslog.syslog == nil {
			log.Println("cant connet syslog server, code " + strconv.Itoa(utils.ConnectSyslogErr))
			return
		}

		switch exType { //根据syslog日志的类型，对应获取不同数据，更新对应的时间边界点
		case models.SysLog_BenchScanLog:
			benchMarkLog := new(models.BenchMarkLog)
			TEPointObj := new(models.TimeEdgePoint)
			TEPointObj.EdgePointCode = exType
			TEPinDB := TEPointObj.Get()
			if len(TEPinDB) > 0 {
				benchMarkLog.UpdateTime = TEPinDB[0].TimePointA
				loglist := benchMarkLog.List(from, limit, false)

				if loglist.Code == 200 && loglist.Data != nil {
					mapdata := loglist.Data.(map[string]interface{})
					for _, logitem := range mapdata["items"].([]*models.BenchMarkLog) {
						logitemJson, _ := json.Marshal(logitem)
						GlobalSyslog.SendSysLog(exType, models.Log_level_Info, string(logitemJson))
					}
				}

				//TEPinDB[0].TimePointA = time.Now().Add(time.Hour * -8).Format("2006-01-02T15:04:05Z")
				//TEPinDB[0].Update()
			}
		case models.SysLog_IDSLog:
			intrudeDetectLog := new(models.IntrudeDetectLog)
			TEPointObj := new(models.TimeEdgePoint)
			TEPointObj.EdgePointCode = exType
			TEPinDB := TEPointObj.Get()
			if len(TEPinDB) > 0 {
				intrudeDetectLog.ToTime = TEPinDB[0].TimePointA
				loglist := intrudeDetectLog.List1(from, limit)
				if loglist.Code == 200 && loglist.Data != nil {
					mapdata := loglist.Data.(map[string]interface{})
					for _, logitem := range mapdata["items"].([]*models.DcokerIds) {
						logitemJson, _ := json.Marshal(logitem)
						GlobalSyslog.SendSysLog(exType, models.Log_level_Info, string(logitemJson))
					}
				}
				TEPinDB[0].TimePointA = time.Now().Add(time.Hour * -8).Format("2006-01-02T15:04:05Z")
				TEPinDB[0].Update()
			}

		case models.SysLog_ContainerVirusLog:
			dockerVirus := new(models.DockerVirus)
			TEPointObj := new(models.TimeEdgePoint)
			TEPointObj.EdgePointCode = exType
			TEPinDB := TEPointObj.Get()
			if len(TEPinDB) > 0 {
				timeTemplate1 := "2006-01-02T15:04:05Z"
				stamp, _ := time.ParseInLocation(timeTemplate1, TEPinDB[0].TimePointA, time.Local)
				dockerVirus.CreatedAt = stamp.Unix()
				loglist := dockerVirus.List(from, limit)
				limit = 30
				if loglist.Code == 200 && loglist.Data != nil {
					mapdata := loglist.Data.(map[string]interface{})
					for _, logitem := range mapdata["items"].([]*models.DockerVirus) {
						logitemJson, _ := json.Marshal(logitem)
						GlobalSyslog.SendSysLog(exType, models.Log_level_Info, string(logitemJson))
					}
				}

				TEPinDB[0].TimePointA = time.Now().Format("2006-01-02T15:04:05Z")
				TEPinDB[0].Update()
			}

		case models.SysLog_ImageSecLog:
			imageVirus := new(models.ImageVirus)
			TEPointObj := new(models.TimeEdgePoint)
			TEPointObj.EdgePointCode = exType
			TEPinDB := TEPointObj.Get()
			if len(TEPinDB) > 0 {
				timeTemplate1 := "2006-01-02T15:04:05Z"
				stamp, _ := time.ParseInLocation(timeTemplate1, TEPinDB[0].TimePointA, time.Local)
				imageVirus.CreatedAt = stamp.Unix()
				loglist := imageVirus.List(from, limit)
				limit = 30
				if loglist.Code == 200 && loglist.Data != nil {
					mapdata := loglist.Data.(map[string]interface{})
					for _, logitem := range mapdata["items"].([]*models.ImageVirus) {
						logitemJson, _ := json.Marshal(logitem)
						GlobalSyslog.SendSysLog(exType, models.Log_level_Info, string(logitemJson))
					}
				}

				TEPinDB[0].TimePointA = time.Now().Format("2006-01-02T15:04:05Z")
				TEPinDB[0].Update()
			}

		case models.SysLog_SecAuditLog:

		}
	}
}
