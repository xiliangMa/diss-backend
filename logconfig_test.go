package main

import (
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/service/system/system"
	"github.com/xiliangMa/diss-backend/sysinit"
	"testing"
)

func Test_Syslog(t *testing.T) {

	sysinit.InitGlobalLogConfig()
	system.SendSysLog(models.SysLog_BenchScanLog, models.SysLog_IDSLog, "IDS log  message 079089")

}
