package k8s

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"time"
)

type Pod struct {
	Id           string `orm:"pk;description(pod id)"`
	Name         string `orm:"unique;description(集群名)"`
	PodIp        string `orm:"default(null);description(pod ip)"`
	Status       string `orm:"description(pod状态)"`
	GroupId      string `orm:"default(null);description(租户id)"`
	GroupName    string `orm:"default(null);description(租户名)"`
	HostIp       string `orm:"default(null);description(主机ip， 默认内网ip)"`
	HostName     string `orm:"default(null);description(host name)"`
	NamSpaceName string `orm:"default(null);description(命名空间)"`
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

func (this *Pod) List(hostName string, from, limit int) models.Result {
	o := orm.NewOrm()
	orm.DefaultTimeLoc = time.Local
	o.Using("default")
	var PodList []*Pod = nil
	var total = 0
	var ResultData models.Result
	var err error
	// to do Multiple conditions
	if this.Name != "" {
		_, err = o.QueryTable(utils.Pod).Filter("name", this.Name).Limit(limit, from).Filter("host_name", hostName).All(&PodList)
	} else {
		_, err = o.QueryTable(utils.Pod).Limit(limit, from).Filter("host_name", hostName).All(&PodList)
	}
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetPodErr
		logs.Error("Get Pod List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	if PodList != nil {
		total = len(PodList)
	}
	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = PodList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}
