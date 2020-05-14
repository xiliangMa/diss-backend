package main

import (
	"fmt"
	"github.com/xiliangMa/diss-backend/service/task"
	"github.com/xiliangMa/diss-backend/sysinit"
	"testing"
)

func Test_Syslog(t *testing.T) {

	sysinit.InitGlobalLogConfig()
	//task.InitGlobalSyslogHander()
	syslogHandler := task.GlobalSysLogTaskHandler

	fmt.Printf("------------------sys export groups and taskids :\n %#v\n", syslogHandler.ExportTypes)
	select {}
}
