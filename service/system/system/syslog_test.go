package system

import (
	"encoding/json"
	"github.com/beego/beego/v2/core/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
	"strconv"
	"testing"
	"time"
)

//func Test_Syslog(t *testing.T) {
//
//	sysinit.InitGlobalLogConfig()
//	syslogHandler := task.GlobalSysLogTaskHandler
//
//	fmt.Printf("Sys export groups and taskids :\n %#v\n", syslogHandler)
//	select {}
//}
//
//func Test_AddTEPoint(t *testing.T) {
//	sysinit.InitGlobalLogConfig()
//	sysinit.InitTimeEdgePoint()
//}

func Test_BenchmarkLog_list_timeranged(t *testing.T) {

	exType := models.LogType_ImageSecLog

	//default from and limit , for 3000 records
	from := 0
	limit := 3000

	logs.Info("Sync syslog data , type:", exType)

	GlobalSyslog.OpenSyslog("init synclog")
	if GlobalSyslog == nil {
		logs.Info("cant connet syslog server, code " + strconv.Itoa(utils.ConnectSyslogErr))
		return
	}

	switch exType { //根据syslog日志的类型，对应获取不同数据，更新对应的时间边界点
	case models.LogType_BenchMarkLog:
		benchMarkLog := new(models.BenchMarkLog)
		TEPointObj := new(models.TimeEdgePoint)
		TEPointObj.EdgePointCode = exType
		TEPinDB := TEPointObj.Get()
		if len(TEPinDB) > 0 {
			benchMarkLog.UpdateTime = TEPinDB[0].TimePointA
			benchMarkLog.IsInfo = false
			loglist := benchMarkLog.List(from, limit)

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
	case models.LogType_IntrudeDetectLog:
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

	case models.LogType_ContainerVirusLog:
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

	case models.LogType_ImageSecLog:
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

	case models.LogType_ContainerSecAuditLog:

	}

}
