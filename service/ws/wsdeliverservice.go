package ws

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	"github.com/xiliangMa/diss-backend/models"
)

type WSDeliverService struct {
	*models.Hub
	Bath                 int64
	CurrentBatchTaskList []*models.Task
	DelTask              *models.Task
}

func (this *WSDeliverService) DeliverTaskToNats() {
	logs.Info("################ Deliver Task <<<start>>> ################")
	for _, task := range this.CurrentBatchTaskList {
		subject := ""
		// 根据是host还是container类型的任务获得topic
		if task.Host != nil {
			hostCount := task.Host.Count()
			if hostCount > 0 {
				subject = task.Host.Id
			}
		} else if task.Container != nil {
			containerCount := task.Container.Count()
			if containerCount > 0 && task.Container.HostName != "" {
				subject = task.Container.HostName
			}
		}
		if subject != "" {
			result := models.WsData{Type: models.Type_Control, Tag: models.Resource_Task, Data: task, RCType: models.Resource_Control_Type_Post}
			data, _ := json.Marshal(result)

			err := models.Nats.Conn.Publish(subject, data)
			if err == nil {
				logs.Info("Deliver Task to Nats Success, Subject: %s Id: %s, data: %v", subject, task.Id, result)
			} else {
				//更新 task 状态
				task.Status = models.Task_Status_Deliver_Failed
				task.Update()
				logs.Error("Deliver Task to Nats Fail,  Subject: %s Id: %s, err: %s", subject, task.Id, err.Error())
			}
		}
	}
	logs.Info("################ Deliver Task <<<end>>> ################")
}

func (this *WSDeliverService) DeliverTask() {
	logs.Info("################ Deliver Task <<<start>>> ################")
	for _, task := range this.CurrentBatchTaskList {
		if _, ok := this.Hub.DissClient[task.Host.Id]; ok {
			client := this.Hub.DissClient[task.Host.Id]
			result := models.WsData{Type: models.Type_Control, Tag: models.Resource_Task, Data: task, RCType: models.Resource_Control_Type_Post}
			data, _ := json.Marshal(result)
			err := client.Conn.WriteMessage(websocket.TextMessage, data)
			if err == nil {
				msg := fmt.Sprintf("Deliver Task Success, Id: %s, data: %v", task.Id, result)
				logs.Info(msg)
				taskRawInfo, _ := json.Marshal(task)
				taskLog := models.TaskLog{RawLog: msg, Task: string(taskRawInfo), Account: task.Account, Level: models.Log_level_Info}
				taskLog.Add()
			} else {
				//更新 task 状态
				task.Status = models.Task_Status_Deliver_Failed
				task.Update()
				msg := fmt.Sprintf("Deliver Task Fail, Id: %s, err: %s", task.Id, err.Error())
				logs.Error(msg)
				taskRawInfo, _ := json.Marshal(task)
				taskLog := models.TaskLog{RawLog: msg, Task: string(taskRawInfo), Account: task.Account, Level: models.Log_level_Error}
				taskLog.Add()
			}
		} else {
			//更新 task 状态
			task.Status = models.Task_Status_Deliver_Failed
			task.Update()
			errMsg := "Agent not connect"
			msg := fmt.Sprintf("Deliver Task Fail, Id: %s, err: %s", task.Id, errMsg)
			logs.Error(msg)
			taskRawInfo, _ := json.Marshal(task)
			taskLog := models.TaskLog{RawLog: msg, Task: string(taskRawInfo), Account: task.Account, Level: models.Log_level_Error}
			taskLog.Add()
		}
	}
	logs.Info("################ Deliver Task <<<end>>> ################")
}
