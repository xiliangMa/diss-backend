package main

import (
	"encoding/json"
	"fmt"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/service/system/system"
	"github.com/xiliangMa/diss-backend/service/task"
	"github.com/xiliangMa/diss-backend/sysinit"
	"testing"
	"time"
)

func Test_Syslog(t *testing.T) {

	sysinit.InitGlobalLogConfig()
	syslogHandler := task.GlobalSysLogTaskHandler

	fmt.Printf("Sys export groups and taskids :\n %#v\n", syslogHandler)
	select {}
}

func Test_AddTEPoint(t *testing.T) {
	sysinit.InitGlobalLogConfig()
	sysinit.InitTimeEdgePoint()
}

func Test_BenchmarkLog_list_timeranged(t *testing.T) {

	sysinit.InitGlobalLogConfig()
	sysinit.InitTimeEdgePoint()

	exType := models.SysLog_ImageSecLog

	//default from and limit , for 3000 records
	from := 0
	limit := 3000

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
				for index, logitem := range mapdata["items"].([]*models.BenchMarkLog) {
					fmt.Printf(" logs item %d:  %#v \n", index, logitem)
					logitemJson, _ := json.Marshal(logitem)
					system.SendSysLog(exType, models.Log_level_Info, string(logitemJson))
				}
			}

			TEPinDB[0].TimePointA = time.Now().Add(time.Hour * -8).Format("2006-01-02T15:04:05Z")
			TEPinDB[0].Update()

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
				for index, logitem := range mapdata["items"].([]*models.DcokerIds) {
					fmt.Printf(" logs item %d:  %#v \n", index, logitem)
					logitemJson, _ := json.Marshal(logitem)
					system.SendSysLog(exType, models.Log_level_Info, string(logitemJson))
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
				for index, logitem := range mapdata["items"].([]*models.DockerVirus) {
					fmt.Printf(" logs item %d:  %#v \n", index, logitem)
					logitemJson, _ := json.Marshal(logitem)
					system.SendSysLog(exType, models.Log_level_Info, string(logitemJson))
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
				for index, logitem := range mapdata["items"].([]*models.ImageVirus) {
					fmt.Printf(" logs item %d:  %#v \n", index, logitem)
					logitemJson, _ := json.Marshal(logitem)
					system.SendSysLog(exType, models.Log_level_Info, string(logitemJson))
				}
			}

			TEPinDB[0].TimePointA = time.Now().Format("2006-01-02T15:04:05Z")
			TEPinDB[0].Update()
		}

	case models.SysLog_SecAuditLog:

	}

}
