package main

import (
	"fmt"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/service/task"
	"github.com/xiliangMa/diss-backend/sysinit"
	"testing"
)

func Test_Syslog(t *testing.T) {

	sysinit.InitGlobalLogConfig()
	//task.InitGlobalSyslogHander()
	syslogHandler := task.GlobalSysLogTaskHandler

	fmt.Printf("Sys export groups and taskids :\n %#v\n", syslogHandler)
	select {}
}

func Test_AddTEPoing(t *testing.T) {
	sysinit.InitGlobalLogConfig()
	//task.InitGlobalSyslogHander()
	//task.GlobalSysLogTaskHandler
	sysinit.InitTimeEdgePoint()
}

func Test_BenchmarkLog_list_timeranged(t *testing.T) {
	benchMarkLog := new(models.BenchMarkLog)
	from := 0
	limit := 3000
	benchMarkLog.UpdateTime = "2020-05-11T17:30:30Z"

	bmlogs := benchMarkLog.List(from, limit, false)
	fmt.Printf("benchmark logs time ranged: \n %#v", bmlogs)
}
