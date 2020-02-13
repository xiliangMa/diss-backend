package k8s

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
)

type Pod struct {
	Id           string `orm:"pk;description(pod id)"`
	Name         string `orm:"unique;description(集群名)"`
	PodIP        string `orm:"default(null);description(pod ip)"`
	Status       string `orm:"description(pod状态)"`
	GroupId      string `orm:"default(null);description(租户id)"`
	GroupName    string `orm:"default(null);description(租户名)"`
	ClusterId    string `orm:"default(null);description(集群id)"`
	HostIP       string `orm:"default(null);description(主机ip， 默认内网ip)"`
	NameSpaceId  string `orm:"default(null);description(命名空间id)"`
	NamSpaceName string `orm:"default(null);description(命名空间)"`
}

type PodInterface interface {
	Add()
	Delete()
	Edit()
	Get()
	List() å
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
