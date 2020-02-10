package k8s

import (
	"github.com/astaxie/beego/orm"
)

type CLuster struct {
	Id       string `orm:"pk;description(集群id)"`
	Name     string `orm:"description(集群名)"`
	FileName string `orm:"description(k8s 文件)"`
	Status   uint8  `orm:"description(集群状态)"`
	IsSync   bool   `orm:"description(是否同步)"`
}

func init() {
	orm.RegisterModel(new(CLuster))
}

type ClusterInterface interface {
	Add()
	Delete()
	Edit()
	Get()
	List()
}
