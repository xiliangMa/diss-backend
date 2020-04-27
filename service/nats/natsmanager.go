package nats

import (
	"github.com/astaxie/beego/logs"
	stan "github.com/nats-io/stan.go"
	"github.com/xiliangMa/diss-backend/utils"
)

type NatsManager struct {
	Conn      stan.Conn
	ClusterId string
	ClientId  string
	Err       error
}

func NewNatsManager() *NatsManager {
	serverUrl := utils.GetNatsServerUrl()
	clusterId := utils.GetNatsClusterId()
	clientId := utils.GetNatsClientId()
	nc, err := stan.Connect(clusterId, clientId, stan.NatsURL(serverUrl))
	if err != nil {
		logs.Error("Create NatsManager fail, ServerUrl: %s ClusterId: %s ClientId: %s, err: %s", serverUrl, clusterId, clientId, err)
	}
	return &NatsManager{
		Conn:      nc,
		Err:       err,
		ClientId:  clientId,
		ClusterId: clusterId,
	}
}
