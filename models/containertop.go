package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"time"
)

type ContainerPs struct {
	Id          string        `orm:"pk;description(id)"`
	HostId      string        `orm:"description(主机id)"`
	PID         string        `orm:"description(PID)"`
	User        string        `orm:"description(用户)"`
	ContainerId string        `orm:"description(容器id)"`
	CPU         string        `orm:"description(CPU)"`
	Mem         string        `orm:"description(内存)"`
	Time        string        `orm:"description(时间)"`
	Start       string        `orm:"description(运行时长 非mac)"`
	Started     string        `orm:"description(运行时长 mac)"`
	Command     orm.TextField `orm:"description(Command)"`
}

func init() {
	orm.RegisterModel(new(ContainerPs))
}

type ContainerPsInterface interface {
	Add()
	Delete()
	Edit()
	Get()
	List()
}

func (this *ContainerPs) Add() Result {
	o := orm.NewOrm()
	o.Using("default")
	var ResultData Result
	var err error
	var containerTopList []*ContainerPs
	cond := orm.NewCondition()
	cond = cond.And("id", this.Id)
	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}

	_, err = o.QueryTable(utils.ContainerPs).SetCond(cond).All(&containerTopList)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetContainerPsErr
		logs.Error("Get ContainerPs failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
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
		ResultData.Code = utils.AddContainerPsErr
		logs.Error("Add ContainerPs failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *ContainerPs) List(from, limit int) Result {
	o := orm.NewOrm()
	orm.DefaultTimeLoc = time.Local
	o.Using("default")
	var ContainerList []*ContainerPs = nil
	var ResultData Result
	var err error

	cond := orm.NewCondition()
	if this.ContainerId != "" {
		cond = cond.And("container_id", this.ContainerId)
	}
	if this.Command != "" {
		cond = cond.And("command__icontains", this.Command)
	}
	_, err = o.QueryTable(utils.ContainerPs).SetCond(cond).Limit(limit, from).All(&ContainerList)

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetContainerPsErr
		logs.Error("Get ContainerPs List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	total, _ := o.QueryTable(utils.ContainerPs).SetCond(cond).Count()
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

func (this *ContainerPs) Update() Result {
	o := orm.NewOrm()
	o.Using("default")
	var ResultData Result

	_, err := o.Update(this)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.EditContainerPsErr
		logs.Error("Update ContainerPs: %s failed, code: %d, err: %s", this.Id, ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *ContainerPs) Delete() Result {
	o := orm.NewOrm()
	o.Using("default")
	var ResultData Result

	cond := orm.NewCondition()

	// 从agent同步时 依据 host_id && ContainerId 删除该主机上所有容器进程
	if this.HostId != "" {
		cond = cond.And("host_id", this.HostId)
	}
	//if this.ContainerId != "" {
	//	cond = cond.And("container_id", this.ContainerId)
	//}
	_, err := o.QueryTable(utils.ContainerPs).SetCond(cond).Delete()

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.DeleteContainerPsErr
		logs.Error("Delete ContainerPs failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	return ResultData
}
