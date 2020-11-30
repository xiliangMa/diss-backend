package system

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/service/task"
)

type IntegrationController struct {
	beego.Controller
}

// @Title Log Config
//// @Description SystemIntegration Config
//// @Param token header string true "authToken"
//
//// @Param body body models.LogConfig false "日志配置信息"
//// @Success 200 {object} models.Result
//// @router /system/logconfig [post]
func (this *IntegrationController) AddLogConfig() {
	logConfig := new(models.LogConfig)
	json.Unmarshal(this.Ctx.Input.RequestBody, &logConfig)

	this.Data["json"] = logConfig.Add()
	this.ServeJSON(false)
}

// @Title Update Log Config
// @Description Update Log Config
// @Param token header string true "authToken"
// @Param body body models.LogConfig false "日志配置信息"
// @Success 200 {object} models.Result
// @router /system/logconfig [put]
func (this *IntegrationController) UpdateLogConfig() {
	logConfig := new(models.LogConfig)
	json.Unmarshal(this.Ctx.Input.RequestBody, &logConfig)
	var result models.Result

	logConfig.ConfigName = models.Log_Config_SysLog_Export
	chkconfig := logConfig.Get()
	chkconfigData := chkconfig.Data.([]*models.LogConfig)
	if len(chkconfigData) > 0 {
		logConfig.Id = chkconfigData[0].Id
		result = logConfig.Update()
	} else {
		result = logConfig.Add()
	}

	this.Data["json"] = result

	// 更新log全局配置
	logConfigObj := result.Data.(*models.LogConfig)
	models.GlobalLogConfig[models.Log_Config_SysLog_Export] = logConfigObj
	// 重新部署任务
	syslogTaskHandler := task.GlobalSysLogManager
	syslogTaskHandler.ReGenSyncSyslogTask()
	this.ServeJSON(false)
}

// @Title Get Log Config
// @Description Get Log Config
// @Param token header string true "authToken"
// @Param configName query string true "配置项名，支持的值:  SysLogExport"
// @Success 200 {object} models.Result
// @router /system/logconfig [get]
func (this *IntegrationController) GetLogConfig() {
	logConfig := new(models.LogConfig)
	configName := this.GetString(":configName")
	logConfig.ConfigName = configName
	json.Unmarshal(this.Ctx.Input.RequestBody, &logConfig)

	this.Data["json"] = logConfig.Get()
	this.ServeJSON(false)
}
