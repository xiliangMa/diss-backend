package task

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/service/nats"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
)

type TaskService struct {
	Task *models.Task
}

func (this *TaskService) RemoveTask() models.Result {
	task := this.Task
	var result models.Result
	natsSubService := nats.NatsSubService{DelTask: task}
	err := natsSubService.DeleteTask()
	if err == nil {
		// 更新数据库状态 设置为删除锁定
		task.Status = models.Task_Status_Removing
		task.Update()
		msg := fmt.Sprintf("Update Task success, status: %s, Id: %s", task.Status, task.Id)
		logs.Info(msg)
		taskRawInfo, _ := json.Marshal(task)
		taskLog := models.TaskLog{Account: task.Account, RawLog: msg, Task: string(taskRawInfo), Level: models.Log_level_Info}
		taskLog.Add()
	} else {
		//如果操作资源不存在（无法给nats下发任务） 直接删除
		if err.Error() == string(utils.ResourceNotFoundErr) {
			task.Delete()
			msg := fmt.Sprintf("Delete Task success, status: %s, Id: %s", task.Status, task.Id)
			logs.Error(msg)
			taskRawInfo, _ := json.Marshal(task)
			taskLog := models.TaskLog{Account: task.Account, RawLog: msg, Task: string(taskRawInfo), Level: models.Log_level_Info}
			taskLog.Add()
		} else {
			result.Code = utils.DeleteTaskErr
			result.Message = "DeleteTaskErr"
			result.Data = nil
			msg := fmt.Sprintf("Delet Task fail, Id: %s, err: %s", task.Id, result.Message)
			taskRawInfo, _ := json.Marshal(task)
			taskLog := models.TaskLog{Account: task.Account, RawLog: msg, Task: string(taskRawInfo), Level: models.Log_level_Info}
			taskLog.Add()
		}
	}
	result.Code = http.StatusOK
	return result
}
