package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"time"
)

type Host struct {
	Id            int       `orm:"not null;auto;description(主机id)"`
	HostLabel     string    `orm:"not null;description(标签)"`
	HostName      string    `orm:"not null;description(主机名)"`
	HostIp        string    `orm:"not null;description(主机IP)"`
	HostDesc      string    `orm:"null;description(主机名)"`
	State         string    `orm:"null;description(状态)"`
	PublicAddress string    `orm:"description(主机名)"`
	CreateTime    time.Time `orm:"not null;description(创建时间);auto_now_add;type(datetime)"`
	UpdateTime    time.Time `orm:"null;description(更新时间);auto_now;type(datetime)"`
	CpuKernel     float64   `orm:"null;;description(cpu)"`
	Mem           float64   `orm:"null;description(内存)"`
	Disk          float64   `orm:"null;description(磁盘)"`
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

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	return ResultData
}

func Internal_GetHost(hostname string) map[string]interface{} {
	o := orm.NewOrm()
	o.Using("default")
	var host Host
	data := make(map[string]interface{})

	err := o.QueryTable("host").Filter("host_name", hostname).One(&host)
	if err == orm.ErrNoRows {
		logs.Info("Not Get one host , can be forward,  code: %d, warn: %s", utils.GetHostZero, "Get Host Zero")
	}

	if host.Id != 0 {
		data["items"] = host
		data["total"] = 1
	} else {
		data = nil
	}

	return data
}

func GetHostWithMetric(hostname string) Result {
	var ResultData Result

	alldata := make(map[string]interface{})
	hostdata := Internal_GetHost(hostname)
	ResultData = Internal_HostMetricInfo_M(hostname)
	if ResultData.Message != "" {
		ResultData.Code = utils.ElasticConnErr
		ResultData.Message = "Cant Connect ElaticSearch"
		return ResultData
	}
	pureMetric := utils.ExtractHostInfo(ResultData.Data.([]interface{}))
	dockerContainerSummary := Internal_ContainerSummaryInfo(hostname)
	dockerContainerRuning := Internal_ContainerListMetricInfo(hostname)

	// alldata包含：主机基本配置，主机动态指标获取，运行的容器汇总，运行中的容器列表
	alldata["hostConfig"] = hostdata
	alldata["hostMetric"] = pureMetric
	alldata["containerSummary"] = dockerContainerSummary.Data
	alldata["containerRunning"] = dockerContainerRuning.Data

	ResultData.Code = http.StatusOK
	ResultData.Data = alldata
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
	ResultData.Code = http.StatusOK
	ResultData.Data = id
	return ResultData
}

func AddHostProcessing(h Host) interface{} {
	var ResultData Result

	// host exist detect
	existhost := Internal_GetHost(h.HostName)
	if existhost != nil {
		ResultData.Code = utils.HostExistError
		ResultData.Message = "Host Exist"
		ResultData.Data = existhost
		logs.Error("AddHost failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	// get host metric by hostname
	ResultData = Internal_HostMetricInfo_M(h.HostName)
	if ResultData.Code != http.StatusOK {
		if ResultData.Code == utils.ElasticConnErr {
			ResultData.Message = "Cant Connect ElaticSearch, Please retry"
		}
		if ResultData.Code == utils.ElasticSearchErr {
			ResultData.Message = "ElasticSearch fetch data Error, Please retry"
		}

		return ResultData
	}
	pureMetric := utils.ExtractHostInfo(ResultData.Data.([]interface{}))

	if len(pureMetric) == 0 {
		ResultData.Code = utils.GetHostMetricError
		ResultData.Message = "Host Metric cant acquire"
		logs.Error("AddHost failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	h.CpuKernel = pureMetric["cpu.cores"].(float64)
	h.Disk = pureMetric["filesystem.total"].(float64)
	h.Mem = pureMetric["memory.total"].(float64)

	// add host
	hostadded := AddHost(&h)
	ResultData.Data = hostadded
	ResultData.Code = http.StatusOK

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
	ResultData.Code = http.StatusOK
	return ResultData
}
