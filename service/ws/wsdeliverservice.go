package ws

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	"github.com/xiliangMa/diss-backend/models/job"
	"github.com/xiliangMa/diss-backend/models/ws"
)

type WSDeliverService struct {
	*Hub
	Bath                 int64
	CurrentBatchTaskList []*job.Task
}

func (this *WSDeliverService) DeliverTask() {
	logs.Info("################ Deliver Task <<<start>>> ################")
	for _, task := range this.CurrentBatchTaskList {
		if _, ok := this.Hub.DissClient[task.Host.Id]; ok {
			client := this.Hub.DissClient[task.Host.Id]
			result := ws.WsData{Type: "Task", Tag: task.Type, Data: task, Config: ""}
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
