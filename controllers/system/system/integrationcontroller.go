package system

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/service/task"
	"net/http"
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

	result := logConfig.Update()
	this.Data["json"] = result

	// 更新log全局配置
	logConfigObj := result.Data.(*models.LogConfig)
	models.GlobalLogConfig[models.Log_Config_SysLog_Export] = logConfigObj
	// 重新部署任务
	syslogTaskHandler := task.GlobalSysLogTaskHandler
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

	var ResultData models.Result
	ResultData.Code = http.StatusOK
	ResultData.Data = logConfig.InnerGet()

	this.Data["json"] = ResultData
	this.ServeJSON(false)
}

//// 外部传入日志记录转发
// @Title
