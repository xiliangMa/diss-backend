package task

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/service/nats"
	"net/http"
)

type TaskService struct {
	Task *models.Task
}

func (this *TaskService) RemoveTask() models.Result {
	task := this.Task
	natsSubService := nats.NatsSubService{DelTask: task}
	natsSubService.DeleteTask()

	result := task.Delete()
	msg := ""
	if result.Code == http.StatusOK {
		msg = fmt.Sprintf("Delete Task success, status: %s, Id: %s", task.Status, task.Id)
		logs.Info(msg)
	} else {
		msg = fmt.Sprintf("Delet Task fail, Id: %s, err: %s", task.Id, result.Message)
		logs.Error(msg)
	}

	taskRawInfo, _ := json.Marshal(task)
	taskLog := models.TaskLog{Account: task.Account, RawLog: msg, Task: string(taskRawInfo), Level: models.Log_level_Info}
	taskLog.Add()
	return result
}
