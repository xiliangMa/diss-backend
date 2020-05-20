package task

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/service/system/system"
	"os"
	"strings"
)

type LogExportGroup struct {
	ExportName string
	TaskId     string
	LastTime   string
}

type SyslogTaskHandler struct {
	CurGroup    string
	ExportTypes *map[string]LogExportGroup
}

func NewSyslogTaskHandler(exportTypes *map[string]LogExportGroup) *SyslogTaskHandler {
	return &SyslogTaskHandler{ExportTypes: exportTypes}
}

var GlobalSysLogTaskHandler *SyslogTaskHandler

func InitGlobalSyslogHander() {
	exportGroups := AllExGroups()
	GlobalSysLogTaskHandler = NewSyslogTaskHandler(exportGroups)
	GlobalSysLogTaskHandler.ReGenSyncSyslogTask()
}

func AllExGroups() *map[string]LogExportGroup {
	exportGroups := map[string]LogExportGroup{}
	exportGroups[models.SysLog_ImageSecLog] = LogExportGroup{ExportName: "镜像安全日志"}
	exportGroups[models.SysLog_BenchScanLog] = LogExportGroup{ExportName: "基线扫描日志"}
	exportGroups[models.SysLog_IDSLog] = LogExportGroup{ExportName: "入侵检测日志"}
	exportGroups[models.SysLog_ContainerVirusLog] = LogExportGroup{ExportName: "容器杀毒日志"}
	exportGroups[models.SysLog_SecAuditLog] = LogExportGroup{ExportName: "安全审计日志"}
	return &exportGroups
}

func (this *SyslogTaskHandler) ReGenSyncSyslogTask() {
	syslogConfig := models.GlobalLogConfig[models.Log_Config_SysLog_Export]
	th := NewTaskHandler()

	for _, exGroup := range *this.ExportTypes {
		if exGroup.TaskId != "" {
			th.DelByID(exGroup.TaskId)
			logs.Info("Delete logsync task type - %s \n", exGroup.TaskId)
		}
	}
	if len(*this.ExportTypes) != 0 {
		exgroup := map[string]LogExportGroup{}
		this.ExportTypes = &exgroup
	}

	if syslogConfig.Enabled == true {
		exportTypes := strings.Split(syslogConfig.ExportedTypes, ",")
		for _, exType := range exportTypes {
			logExGroup := (*this.ExportTypes)[exType]
			logExGroup.TaskId = exType
			(*this.ExportTypes)[exType] = logExGroup
			if err := th.AddByFunc(exType, beego.AppConfig.String("syslog::SyslogSyncSpec"), system.GetSyncSyslogFunc(exType)); err != nil {
				logs.Error("error to add syslog TaskHandler task: %s", err)
				os.Exit(-1)
			}
		}
	}

	th.Start()

}
