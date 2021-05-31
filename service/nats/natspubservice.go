package nats

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/nats-io/stan.go"
	"github.com/xiliangMa/diss-backend/models"
	"net/http"
)

type NatsPubService struct {
	Message       []byte
	Conn          stan.Conn
	ClientSubject string
	Type          string
}

func (this *NatsPubService) PublishData(data []byte) {
	err := this.Conn.Publish(this.ClientSubject, data)
	logs.Info("Publish success, subject: %s data : %s", this.ClientSubject, data)
	if err != nil {
		logs.Error("Nats ############################ Publish to agent fail ############################", err)
	}
}

func (this *NatsPubService) RuleDefinePub() {

	natsData := models.NatsData{Code: http.StatusOK, Type: models.Type_Config, Tag: models.Config_RuleDefineList, RCType: models.Resource_Control_Type_Post}

	// 获取规则列表，并发送
	ruledefineObj := models.RuleDefine{}
	ruledefineObj.Type = this.Type
	rulelist := ruledefineObj.List(0, 0)
	natsData.Data = rulelist.Data
	rulesData, _ := json.Marshal(natsData)

	if this.ClientSubject != "" {
		this.PublishData(rulesData)
	} else {
		// 循环已授权的主机列表，逐个发送
		hostObj := models.HostConfig{LicCount: true, IsLicensed: true}
		_, _, hostList := hostObj.BaseList(0, 0)
		for _, host := range hostList {
			this.ClientSubject = host.Id
			this.PublishData(rulesData)
		}
	}
}
