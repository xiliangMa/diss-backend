package bean

import (
	"github.com/xiliangMa/diss-backend/models"
)

type SecurityCheckList struct {
	CheckList []*SecurityCheck
}

type SecurityCheck struct {
	BenchMarkCheck bool                    `orm:"default(false)" description:"(开启基线检测)"`
	VirusScan      bool                    `orm:"default(false)" description:"(开启病毒)"`
	LeakScan       bool                    `orm:"default(false)" description:"(开启漏洞)"`
	Host           *models.HostConfig      `orm:"null" description:"(主机)"`
	Container      *models.ContainerConfig `orm:"null" description:"(容器)"`
	Type           string                  `orm:"default(host)" description:"(类型 host  container)"`
}
