package base

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"net/http"
)

type TaskService struct {
	*models.Task
	*models.Result
}

func (this *TaskService) UpdateTaskStatus() {
	if this.Task == nil && this.Task.Id == "" {
		return
	}
	level := models.Log_level_Info
	switch this.Task.Status {
	case models.Task_Status_Finished:
		this.Task.RunCount = 1
	case models.Task_Status_Failed:
		level = models.Log_level_Error
	case models.Task_Status_Receive_Failed:
		level = models.Log_level_Error
	}
	this.SaveTaskLog(level)
	this.Task.Update()
}

func (this *TaskService) SaveTaskLog(level string) {
	if this.Result.Code != http.StatusOK {
		logs.Error(this.Result.Message)
	} else {
		logs.Info(this.Result.Message)
	}
	// 更新任务日志
	taskLog := models.TaskLog{}
	taskLog.Account = models.Account_Admin
	taskRawInfo, _ := json.Marshal(this.Task)
	taskLog.Task = string(taskRawInfo)
	taskLog.Level = level
	taskLog.RawLog = this.Result.Message
	taskLog.Add()
}
