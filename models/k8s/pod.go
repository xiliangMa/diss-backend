package k8s

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
)

type Pod struct {
	Id     string `orm:"pk;description(pod id)"`
	Name   string `orm:"unique;description(集群名)"`
	Status uint8  `orm:"description(集群状态)"`
}

type PodInterface interface {
	Add()
	Delete()
	Edit()
	Get()
	List()
}

func init() {
	orm.RegisterModel(new(Pod))
}

func (this *Pod) Add() models.Result {
	o := orm.NewOrm()
	o.Using("default")
	var ResultData models.Result

	_, err := o.Insert(this)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddPodErr
		logs.Error("Add Pod failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}
