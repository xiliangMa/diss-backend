package ws

import (
	"encoding/json"
	"errors"
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
		result := models.WsData{Type: models.Type_Control, Tag: models.Resource_Task, Data: task, RCType: models.Resource_Control_Type_Post}
		data, _ := json.Marshal(result)
		hostName := ""
		if task.Host != nil {
			hostName = task.Host.HostName
		}
		if task.Container != nil {
			hostName = task.Container.HostName
		}
		subject := hostName + `_` + models.Topic_Task
		err := models.NatsManager.Conn.Publish(subject, data)
		if err == nil {
			logs.Info("Deliver Task to Nats Success, Subject: %s Id: %s, data: %v", subject, task.Id, result)
		} else {
			//更新 task 状态
			task.Status = models.Task_Status_Deliver_Failed
			task.Update()
			logs.Error("Deliver Task to Nats Fail,  Subject: %s Id: %s, err: %s", subject, task.Id, err.Error())
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
				taskLog := models.TaskLog{RawLog: msg, Task: task, Account: task.Account}
				taskLog.Add()
			} else {
				//更新 task 状态
				task.Status = models.Task_Status_Deliver_Failed
				task.Update()
				msg := fmt.Sprintf("Deliver Task Fail, Id: %s, err: %s", task.Id, err.Error())
				logs.Error(msg)
				taskLog := models.TaskLog{RawLog: msg, Task: task, Account: task.Account}
				taskLog.Add()
			}
		} else {
			//更新 task 状态
			task.Status = models.Task_Status_Deliver_Failed
			task.Update()
			errMsg := "Agent not connect"
			msg := fmt.Sprintf("Deliver Task Fail, Id: %s, err: %s", task.Id, errMsg)
			logs.Error(msg)
			taskLog := models.TaskLog{RawLog: msg, Task: task, Account: task.Account}
			taskLog.Add()
		}
	}
	logs.Info("################ Deliver Task <<<end>>> ################")
}

func (this *WSDeliverService) DeleteTask() error {
	logs.Info("################ Delete Task <<<start>>> ################")
	task := this.DelTask
	if task.Host == nil {
		return nil
	}
	hostId := task.Host.Id
	if _, ok := this.Hub.DissClient[hostId]; !ok {
		errMsg := "Agent not connect"
		msg := fmt.Sprintf("Deliver Task Fail, Id: %s, err: %s", task.Id, errMsg)
		logs.Error(msg)
		taskLog := models.TaskLog{RawLog: msg, Task: task, Account: task.Account}
		taskLog.Add()
		return errors.New(errMsg)
	}
	client := this.Hub.DissClient[hostId]
	result := models.WsData{Type: models.Type_Control, Tag: models.Resource_Task, RCType: models.Resource_Control_Type_Delete, Data: task}
	data, err := json.Marshal(result)
	if err != nil {
		msg := fmt.Sprintf("Delete Task Fail, Id: %s, err: %s", task.Id, err.Error())
		logs.Error(msg)
		taskLog := models.TaskLog{RawLog: msg, Task: task, Account: task.Account}
		taskLog.Add()
		return err
	}
	err = client.Conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		msg := fmt.Sprintf("Delete Task Fail, Id: %s, err: %s", task.Id, err.Error())
		logs.Error(msg)
		taskLog := models.TaskLog{RawLog: msg, Task: task, Account: task.Account}
		taskLog.Add()
		return err
	}
	logs.Info("Delete Task Success, Id: %s, data: %v", task.Id, result)
	logs.Info("################ Delete Task <<<end>>> ################")
	return nil
}
