package main

import (
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/service/system/system"
	"github.com/xiliangMa/diss-backend/sysinit"
	"testing"
)

func Test_LogConfig(t *testing.T) {
	sysinit.InitGlobalLogConfig()
}

func Test_Syslog(t *testing.T) {
	system.SendSysLog(models.SysLog_BenchScanLog, models.Log_level_Warn, "Test message +088l LLLds")
}
