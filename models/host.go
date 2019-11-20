package models

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"time"
)

type Host struct {
	Id            int       `xorm:"not null pk autoincr INT(11)"`
	HostLabel     string    `xorm:"not null comment('用于界面显示的标签') VARCHAR(255)"`
	HostName      string    `xorm:"not null comment('主机名称') VARCHAR(30)"`
	HostIp        string    `xorm:"not null comment('主机IP') VARCHAR(30)"`
	HostDesc      string    `xorm:"comment('主机说明') VARCHAR(255)"`
	State         string    `xorm:"comment('状态') VARCHAR(30)"`
	PublicAddress string    `xorm:"comment('外部访问地址') VARCHAR(255)"`
	CreateTime    time.Time `xorm:"comment('记录添加时间') DATETIME"`
	UpdateTime    time.Time `xorm:"comment('记录更改时间') DATETIME"`
}

func init() {
	orm.RegisterModel(new(Host))
}

func GetHostList(name, ip string, from, limit int) Result {
	o := orm.NewOrm()
	o.Using("default")
	var HostList []*Host
	var ResultData Result
	_, err := o.QueryTable("host").Filter("name__icontains", name).Filter("host_ip__icontains", ip).Limit(limit, from).All(&HostList)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetHostListErr
		logs.Error("GetHostList failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	total, _ := o.QueryTable("host").Count()
	data := make(map[string]interface{})
	data["items"] = HostList
	data["total"] = total
	ResultData.Code = utils.Success
	ResultData.Data = data
	return ResultData
}

func GetHost(hostname string) Result{
	o := orm.NewOrm()
	o.Using("default")
	var host Host
	data := make(map[string]interface{})
	var ResultData Result

	err := o.QueryTable("host").Filter("host_name", hostname).One(&host)
	if (err == orm.ErrNoRows){
		fmt.Print(err)
	}

	if host.Id != 0 {
		data["items"] = host
		data["total"] = 1
	}else{
		data = nil
	}


	ResultData.Code = utils.Success
	ResultData.Data = data
	return ResultData
}

func AddHost(host *Host) Result {
	o := orm.NewOrm()
	o.Using("default")
	var ResultData Result
	id, err := o.Insert(host)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddHostErr
		logs.Error("AddHost failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = utils.Success
	ResultData.Data = id
	return ResultData
}

func DeleteHost(id int) Result {
	o := orm.NewOrm()
	o.Using("default")
	var ResultData Result
	_, err := o.Delete(&Host{Id: id})
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.DeleteHostErr
		logs.Error("DeleteHost failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = utils.Success
	return ResultData
}
