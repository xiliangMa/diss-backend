package k8s

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"time"
)

type NameSpace struct {
	Id        string `orm:"pk;description(命名空间id)"`
	Name      string `orm:"unique;description(命名空间)"`
	ClusterId string `orm:"default(null);description(集群id)"`
}

type NameSpaceInterface interface {
	Add()
	Delete()
	Edit()
	Get()
	List()
}

func init() {
	orm.RegisterModel(new(NameSpace))
}

func (this *NameSpace) Add() models.Result {
	o := orm.NewOrm()
	o.Using("default")
	var ResultData models.Result
	var nameSpaceList []*NameSpace
	var err error
	cond := orm.NewCondition()
	cond = cond.And("id", this.Id)
	if this.Name != "" {
		cond = cond.And("name__contains", this.Name)
	}
	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}
	_, err = o.QueryTable(utils.NameSpace).SetCond(cond).All(&nameSpaceList)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetNameSpaceErr
		logs.Error("Get NameSpace failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	if len(nameSpaceList) != 0 {
		// agent 或者 k8s 数据更新（因为没有diss-backend的关系数据，所以直接更新）
		return this.Update()
	} else {
		_, err = o.Insert(this)
		if err != nil {
			ResultData.Message = err.Error()
			ResultData.Code = utils.AddNameSpaceErr
			logs.Error("Add NameSpace failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
			return ResultData
		}
	}

	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *NameSpace) List(from, limit int) models.Result {
	o := orm.NewOrm()
	orm.DefaultTimeLoc = time.Local
	o.Using("default")
	var nameSpaceList []*NameSpace
	var ResultData models.Result
	var err error
	cond := orm.NewCondition()

	if this.Name != "" {
		cond = cond.And("name__contains", this.Name)
	}
	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}

	if this.ClusterId != "" {
		cond = cond.And("cluster_id", this.ClusterId)
	}
	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}

	_, err = o.QueryTable(utils.NameSpace).SetCond(cond).Limit(limit, from).All(&nameSpaceList)

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetClusterErr
		logs.Error("Get NameSpace List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	total, _ := o.QueryTable(utils.NameSpace).Count()
	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = nameSpaceList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}

func (this *NameSpace) Update() models.Result {
	o := orm.NewOrm()
	o.Using("default")
	var ResultData models.Result

	_, err := o.Update(this)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.EditNameSpaceErr
		logs.Error("Update NameSpace: %s failed, code: %d, err: %s", this.Name, ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}
