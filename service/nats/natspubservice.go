package nats

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/nats-io/stan.go"
	"github.com/xiliangMa/diss-backend/models"
	"net/http"
	"strings"
)

type NatsPubService struct {
	Message       []byte
	Conn          stan.Conn
	ClientSubject string
	Type          string
}

func (this *NatsPubService) PublishData(data []byte) error {
	err := this.Conn.Publish(this.ClientSubject, data)
	logs.Info("Publish success, subject: %s data : %s", this.ClientSubject, data)
	if err != nil {
		logs.Error("Nats ############################ Publish to agent fail ############################", err)
	}
	return err
}

func (this *NatsPubService) SendToManyHost() []string {
	targethosts := []string{}
	if this.ClientSubject != "" {
		// 发送到指定的主机，支持分号分割的多个Id
		hostids := strings.Split(this.ClientSubject, ",")
		for _, hostid := range hostids {
			this.ClientSubject = hostid
			targethosts = append(targethosts, hostid)
			err := this.PublishData(this.Message)
			if err == nil {
				logs.Info("Deliver UpdateMetricData to Nats Success,  HostId: %s, data: %v", hostid, this.Message)
			}
		}
	} else {
		// 循环已授权的主机列表，逐个发送
		hostObj := models.HostConfig{LicCount: true, IsLicensed: true}
		_, _, hostList := hostObj.BaseList(0, 0)
		for _, host := range hostList {
			this.ClientSubject = host.Id
			targethosts = append(targethosts, host.Id)
			err := this.PublishData(this.Message)
			if err == nil {
				logs.Info("Deliver UpdateMetricData to Nats Success,  HostId: %s, data: %v", host.Id, this.Message)
			}
		}
	}

	return targethosts
}

func (this *NatsPubService) RuleDefinePub() {

	natsData := models.NatsData{Code: http.StatusOK, Type: models.Type_Config, Tag: models.Config_RuleDefineList, RCType: models.Resource_Control_Type_Post}

	// 获取规则列表，并发送
	ruledefineObj := models.RuleDefine{}
	ruledefineObj.Type = this.Type
	rulelist, _, _ := ruledefineObj.RuleDefineList(0, 0)

	// 追加漏洞库规则
	if strings.Contains(this.Type, models.RuleType_DockerVulnerability) {
		vulnerLibObj := models.VulnerabilityLib{}
		vulnerLibObj.Category = models.RuleType_DockerVulnerability
		vulnerLibList, _, _ := vulnerLibObj.VulnerabilityList(0, 0)
		if len(vulnerLibList) > 0 {
			for _, vulnerLib := range vulnerLibList {
				ruleDefine := models.RuleDefine{}
				ruleDefine.Type = vulnerLib.Category
				vulnerRule := models.VulnerRule{}
				vulnerRule.AffectVersion = vulnerLib.SuiteVersions
				vulnerRule.Component = vulnerLib.SubCategory
				vulnerRule.CVE = vulnerLib.CVEId
				vulnerRuleJson, _ := json.Marshal(vulnerRule)
				ruleDefine.Info = string(vulnerRuleJson)

				rulelist = append(rulelist, &ruleDefine)
			}
		}
	}

	natsData.Data = rulelist
	rulesData, _ := json.Marshal(natsData)
	this.Message = rulesData
	this.SendToManyHost()
}
