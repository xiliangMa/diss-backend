package system

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/service/task"
	"log"
)

type IntegrationController struct {
	beego.Controller
}

// @Title SystemIntegration Config
// @Description SystemIntegration Config
// @Param token header string true "authToken"
// @Param configName query string true "Config Name"
// @Param body body models.LogConfig false "日志配置信息"
// @Success 200 {object} models.Result
// @router /system/logconfig [post]
func (this *IntegrationController) AddLogConfig() {
	logConfig := new(models.LogConfig)
	configname := this.GetString("configName")
	json.Unmarshal(this.Ctx.Input.RequestBody, &logConfig)
	if configname != "" {
		logConfig.ConfigName = configname
	}

	this.Data["json"] = logConfig.Add()
	this.ServeJSON(false)
}

// @Title Update SystemIntegration Config
// @Description Update One SystemIntegration Config
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

////内部Action，目前仅用于测试
// @Title Get SystemIntegration Config
// @Description Get One SystemIntegration Config [Inner内部操作]
// @Param token header string true "authToken"
// @Param configName query string true "Config Name"
// @Success 200 {object} models.Result
// @router /system/logconfig [get]
func (this *IntegrationController) InnerGetLogConfig() {
	logConfig := new(models.LogConfig)
	configName := this.GetString(":configName")
	logConfig.ConfigName = configName
	json.Unmarshal(this.Ctx.Input.RequestBody, &logConfig)

	log.Println("==============global log config============", models.GlobalLogConfig[models.Log_Config_SysLog_Export])
	this.Data["json"] = logConfig.InnerGet()
	this.ServeJSON(false)
}

//// 外部传入日志记录转发
// @Title
