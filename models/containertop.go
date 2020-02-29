package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"time"
)

type ContainerTop struct {
	Id          string `orm:"pk;description(容器 top id)"`
	HostId      string `orm:"description(主机Id)"`
	ContainerId string `orm:"description(容器id)"`
	PID         string `orm:"description(PID)"`
	User        string `orm:"description(用户)"`
	Time        string `orm:"description(时间)"`
	Command     string `orm:"description(命令)"`
}

func init() {
	orm.RegisterModel(new(ContainerTop))
}

type ContainerTopInterface interface {
	Add()
	Delete()
	Edit()
	Get()
	List()
}

func (this *ContainerTop) Add() Result {
	o := orm.NewOrm()
	o.Using("default")
	var ResultData Result
	var err error
	var containerTopList []*ContainerTop
	cond := orm.NewCondition()
	cond = cond.And("id", this.Id)
	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}

	_, err = o.QueryTable(utils.ContainerTop).SetCond(cond).All(&containerTopList)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetContainerTopErr
		logs.Error("Get ContainerTop failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	if len(containerTopList) != 0 {
		// agent 或者 k8s 数据更新（因为没有diss-backend的关系数据，所以直接删除在添加）
		if result := this.Delete(); result.Code != http.StatusOK {
			return result
		}
	}
	_, err = o.Insert(this)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddContainerTopErr
		logs.Error("Add ContainerTop failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *ContainerTop) List(from, limit int) Result {
	o := orm.NewOrm()
	orm.DefaultTimeLoc = time.Local
	o.Using("default")
	var ContainerList []*ContainerTop = nil
	var ResultData Result
	var err error

	cond := orm.NewCondition()
	if this.Id != "" {
		cond = cond.And("container_id", this.ContainerId)
	}
	if this.Command != "" {
		cond = cond.And("command__icontains", this.Command)
	}
	_, err = o.QueryTable(utils.ContainerTop).SetCond(cond).Limit(limit, from).All(&ContainerList)

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetContainerTopErr
		logs.Error("Get ContainerTop List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	total, _ := o.QueryTable(utils.ContainerTop).Count()
	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = ContainerList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}

func (this *ContainerTop) Update() Result {
	o := orm.NewOrm()
	o.Using("default")
	var ResultData Result

	_, err := o.Update(this)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.EditContainerTopErr
		logs.Error("Update ContainerTop: %s failed, code: %d, err: %s", this.Id, ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *ContainerTop) Delete() Result {
	o := orm.NewOrm()
	o.Using("default")
	var ResultData Result

	cond := orm.NewCondition()

	if this.ContainerId != "" {
		cond = cond.And("container_id", this.ContainerId)
	}
	_, err := o.QueryTable(utils.ContainerTop).SetCond(cond).Delete()

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.DeleteContainerTopErr
		logs.Error("Delete ContainerTop failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	return ResultData
}
