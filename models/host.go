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
	HostLabel     string    `xorm:"not null default '' comment('用于界面显示的标签') VARCHAR(255)"`
	HostName      string    `xorm:"not null default '' comment('主机名称') VARCHAR(255)"`
	HostIp        string    `xorm:"not null default '' comment('主机IP') VARCHAR(255)"`
	HostDesc      string    `xorm:"default '' comment('主机说明') VARCHAR(255)"`
	State         string    `xorm:"default '' comment('状态') VARCHAR(255)"`
	PublicAddress string    `xorm:"default '' comment('外部访问地址') VARCHAR(255)"`
	CreateTime    time.Time `xorm:"comment('记录添加时间') DATETIME"`
	UpdateTime    time.Time `xorm:"comment('记录更改时间') DATETIME"`
	CpuKernel     float64   `xorm:"DOUBLE"`
	CpuFrequency  float64   `xorm:"DOUBLE"`
	Mem           float64   `xorm:"DOUBLE"`
	Disk          float64   `xorm:"DOUBLE"`
	Network       float64   `xorm:"DOUBLE"`
}

func init() {
	orm.RegisterModel(new(Host))
}

func GetHostList(name, ip string, from, limit int) Result {
	o := orm.NewOrm()
	o.Using("default")
	var HostList []*Host
	var ResultData Result
	_, err := o.QueryTable("host").Limit(limit, from).All(&HostList)
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

func GetHostInternal(hostname string) map[string]interface{} {
	o := orm.NewOrm()
	o.Using("default")
	var host Host
	data := make(map[string]interface{})

	err := o.QueryTable("host").Filter("host_name", hostname).One(&host)
	if err == orm.ErrNoRows {
		fmt.Print(err)
		logs.Error("GetHost failed, code: %d, err: %s", utils.GetHostZero, "Get Host Zero")
		//data["Message"] = "Get Host Error"
		//data["Code"] = utils.GetHostErr
		//return data
	}

	if host.Id != 0 {
		data["items"] = host
		data["total"] = 1
	} else {
		data = nil
	}

	return data
}

func GetHost(hostname string) Result {
	var ResultData Result

	data := GetHostInternal(hostname)
	ResultData.Code = utils.Success
	ResultData.Data = data
	return ResultData
}

func GetHostWithMetric(hostname string) Result {
	var ResultData Result
	alldata := make(map[string]interface{})
	hostdata := GetHostInternal(hostname)
	hostMetricInfo := GetHostMetricInfo_M(hostname)
	alldata["hostConfig"] = hostdata
	alldata["hostMetric"] = hostMetricInfo

	ResultData.Code = utils.Success
	ResultData.Data = alldata
	return ResultData
}

func AddHost(host *Host) Result {
	o := orm.NewOrm()
	o.Using("default")
	var ResultData Result

	//createTime
	timenow := time.Now().Unix()
	formatTime := time.Unix(timenow, 0)
	host.CreateTime = formatTime

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

func AddHostProcessing(h Host) interface{} {
	var ResultData Result

	// host exist detect
	existhost := GetHostInternal(h.HostName)
	if existhost != nil {
		ResultData.Code = utils.HostExistError
		ResultData.Message = "Host Exist"
		ResultData.Data = existhost
		logs.Error("AddHost failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	// get host metric by hostname
	hostMetricInfo := GetHostMetricInfo_M(h.HostName)
	pureMetric := utils.ExtractHostInfo(hostMetricInfo)

	if len(pureMetric) == 0 {
		ResultData.Code = utils.GetHostMetricError
		ResultData.Message = "Host Metric cant acquire"
		ResultData.Data = h
		logs.Error("AddHost failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	h.CpuKernel = pureMetric["cpu.cores"].(float64)
	h.Disk = pureMetric["filesystem.total"].(float64)
	h.Mem = pureMetric["memory.total"].(float64)

	// add host
	hostadded := AddHost(&h)
	ResultData.Data = hostadded
	ResultData.Code = utils.Success

	return ResultData
}
