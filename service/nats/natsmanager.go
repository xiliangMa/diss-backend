package nats

import (
	"github.com/astaxie/beego/logs"
	stan "github.com/nats-io/stan.go"
	"github.com/xiliangMa/diss-backend/utils"
)

type NatsManager struct {
	Conn stan.Conn
	Err  error
}

func NewNatsManager() *NatsManager {
	serverUrl := utils.GetNatsServerUrl()
	nc, err := stan.Connect("test-cluster", "diss-server", stan.NatsURL(serverUrl))
	if err != nil {
		logs.Error("Create NatsManager fail, ServerUrl: %s, err: %s", serverUrl, err)
	}
	return &NatsManager{
		Conn: nc,
		Err:  err,
	}
}
