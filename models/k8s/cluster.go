package k8s

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"time"
)

type Cluster struct {
	Id         string    `orm:"pk;description(集群id)"`
	Name       string    `orm:"unique;description(集群名)"`
	FileName   string    `orm:"description(k8s 文件)"`
	Status     uint8     `orm:"default(0);description(集群状态)"`
	IsSync     bool      `orm:"default(false);description(是否同步)"`
	Synced     bool      `orm:"default(false);description(同步状态)"`
	CreateTime time.Time `orm:"description(创建时间);auto_now_add;type(datetime)"`
	UpdateTime time.Time `orm:"null;description(更新时间);auto_now;type(datetime)"`
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

func (this *Cluster) List(from, limit int) models.Result {
	o := orm.NewOrm()
	orm.DefaultTimeLoc = time.Local
	o.Using("default")
	var ClusterList []*Cluster
	var total = 0
	var ResultData models.Result
	var err error

	if this.Name != "" {
		_, err = o.QueryTable(utils.Cluster).Filter("name", this.Name).Limit(limit, from).All(&ClusterList)
	} else {
		_, err = o.QueryTable(utils.Cluster).Limit(limit, from).All(&ClusterList)
	}

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetClusterErr
		logs.Error("Get Cluster List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	if ClusterList != nil {
		total = len(ClusterList)
	}
	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = ClusterList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}

func (this *Cluster) Update() models.Result {
	o := orm.NewOrm()
	o.Using("default")
	var ResultData models.Result

	_, err := o.Update(this)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.EditHostErr
		logs.Error("Update cluster: %s failed, code: %d, err: %s", this.Name, ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}
