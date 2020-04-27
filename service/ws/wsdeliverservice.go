package ws

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/models/global"
	"github.com/xiliangMa/diss-backend/models/job"
	"github.com/xiliangMa/diss-backend/models/ws"
)

type WSDeliverService struct {
	*ws.Hub
	Bath                 int64
	CurrentBatchTaskList []*job.Task
	DelTask              *job.Task
}

func (this *WSDeliverService) DeliverTaskToNats() {
	logs.Info("################ Deliver Task <<<start>>> ################")
	for _, task := range this.CurrentBatchTaskList {
		result := ws.WsData{Type: ws.Type_Control, Tag: ws.Resource_Task, Data: task, RCType: ws.Resource_Control_Type_Post}
		data, _ := json.Marshal(result)
		hostName := ""
		if task.Host != nil {
			hostName = task.Host.HostName
		}
		if task.Container != nil {
			hostName = task.Container.HostName
		}
		subject := hostName + `_` + models.Topic_Task
		err := global.NatsManager.Conn.Publish(subject, data)
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
			result := ws.WsData{Type: ws.Type_Control, Tag: ws.Resource_Task, Data: task, RCType: ws.Resource_Control_Type_Post}
			data, _ := json.Marshal(result)
			err := client.Conn.WriteMessage(websocket.TextMessage, data)
			if err == nil {
				logs.Info("Deliver Task Success, Id: %s, data: %v", task.Id, result)
			} else {
				//更新 task 状态
				task.Status = models.Task_Status_Deliver_Failed
				task.Update()
				logs.Error("Deliver Task Fail, Id: %s, err: %s", task.Id, err.Error())
			}
		} else {
			//更新 task 状态
			task.Status = models.Task_Status_Deliver_Failed
			task.Update()
			errMsg := "Agent not connect"
			logs.Error("Deliver Task Fail, Id: %s, err: %s", task.Id, errMsg)
		}
	}
	logs.Info("################ Deliver Task <<<end>>> ################")
}

func (this *WSDeliverService) DeleteTask() error {
	logs.Info("################ Delete Task <<<start>>> ################")
	task := this.DelTask
	hostId := task.Host.Id
	if _, ok := this.Hub.DissClient[hostId]; !ok {
		errMsg := "Agent not connect"
		logs.Error("Delete Task Fail, Id: %s, err: %s", task.Id, errMsg)
		return errors.New(errMsg)
	}
	client := this.Hub.DissClient[hostId]
	result := ws.WsData{Type: ws.Type_Control, Tag: ws.Resource_Task, RCType: ws.Resource_Control_Type_Delete, Data: task}
	data, err := json.Marshal(result)
	if err != nil {
		logs.Error("Delete Task Fail, Id: %s, err: %s", task.Id, err.Error())
		return err
	}
	err = client.Conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		logs.Error("Delete Task Fail, Id: %s, err: %s", task.Id, err.Error())
		return err
	}
	logs.Info("Delete Task Success, Id: %s, data: %v", task.Id, result)
	logs.Info("################ Delete Task <<<end>>> ################")
	return nil
}
