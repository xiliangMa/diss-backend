package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/xiliangMa/diss-backend/utils"
	"math/rand"
	"strconv"
)

type NatsManager struct {
	Conn      stan.Conn
	ClusterId string
	ClientId  string
	Err       error
}

func NewNatsManager() *NatsManager {
	servers := utils.GetNatsServes()
	//serverUrl := utils.GetNatsServerUrl()
	clusterId := utils.GetNatsClusterId()
	clientId := utils.GetNatsClientId()
	sc, err := createStanConnect(clusterId, clientId, servers)
	return &NatsManager{
		Conn:      sc,
		Err:       err,
		ClientId:  clientId,
		ClusterId: clusterId,
	}
}

func createStanConnect(clusterId, clientId, servers string) (stan.Conn, error) {
	nc, err := nats.Connect(servers)
	if err != nil {
		logs.Error("Create Nats Connect failed, Servers: %s, err: %s", servers, err)
	}
	var sc stan.Conn
	for {
		id := clientId + strconv.Itoa(rand.Intn(1000))
		sc, err = stan.Connect(clusterId, id, stan.NatsConn(nc))
		if err != nil {
			logs.Error("Create stan Connect failed, Servers: %s ClusterId: %s ClientId: %s, err: %s", servers, clusterId, clientId, err)
			sc, err = stan.Connect(clusterId, id, stan.NatsConn(nc))
		} else {
			logs.Info("Create stan Connect success.")
			break
		}
	}
	return sc, err
}
