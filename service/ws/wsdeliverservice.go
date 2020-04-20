package ws

import (
	"encoding/json"
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	"github.com/xiliangMa/diss-backend/models/job"
	"github.com/xiliangMa/diss-backend/models/ws"
)

type WSDeliverService struct {
	*Hub
	Bath                 int64
	CurrentBatchTaskList []*job.Task
	DelTask              *job.Task
}

func (this *WSDeliverService) DeliverTask() {
	logs.Info("################ Deliver Task <<<start>>> ################")
	for _, task := range this.CurrentBatchTaskList {
		if _, ok := this.Hub.DissClient[task.Host.Id]; ok {
			client := this.Hub.DissClient[task.Host.Id]
			result := ws.WsData{Type: ws.Type_Control, Tag: ws.Resource_Task, Data: task, RCType: ws.Resource_Control_Type_Post}
			data, err := json.Marshal(result)
			err = client.Conn.WriteMessage(websocket.TextMessage, data)
			if err == nil {
				logs.Info("Deliver Task Success, Id: %s, data: %v", task.Id, result)
			} else {
				logs.Error("Deliver Task Fail, Id: %s, err: %s", task.Id, err.Error())
			}
		} else {
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
