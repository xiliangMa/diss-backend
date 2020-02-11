package k8s

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
)

type Cluster struct {
	Id       string `orm:"pk;description(集群id)"`
	Name     string `orm:"description(集群名)"`
	FileName string `orm:"description(k8s 文件)"`
	Status   uint8  `orm:"description(集群状态)"`
	IsSync   bool   `orm:"description(是否同步)"`
}

func init() {
	orm.RegisterModel(new(Cluster))
}

type ClusterInterface interface {
	Add()
	Delete()
	Edit()
	Get()
	List()
}

func (this *Cluster) Add() models.Result {
	o := orm.NewOrm()
	o.Using("default")
	var ResultData models.Result

	_, err := o.Insert(this)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddClusterErr
		logs.Error("Add Cluster failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}
