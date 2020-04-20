package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"time"
)

type Host struct {
	Id           int       `orm:"auto;" description:"(主机id)"`
	HostLabel    string    `orm:"" description:"(标签)"`
	HostName     string    `orm:"" description:"(主机名)"`
	InternalAddr string    `orm:"" description:"(内部访问主机地址，通常为IP)"`
	PublicAddr   string    `orm:"null;" description:"(外部访问地址，IP或域名)"`
	HostDesc     string    `orm:"null;" description:"(主机描述)"`
	State        string    `orm:"null;" description:"(状态)"`
	CreateTime   time.Time `orm:"auto_now_add;type(datetime)" description:"(创建时间)"`
	UpdateTime   time.Time `orm:"null;auto_now;type(datetime)" description:"(更新时间)"`
	CpuKernel    float64   `orm:"null;" description:"(cpu)"`
	Mem          float64   `orm:"null;" description:"(内存)"`
	Disk         float64   `orm:"null;" description:"(磁盘)"`
}

func GetHostList(name, ip string, from, limit int) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
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
	o.Using(utils.DS_Default)
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

func GetHostWithContainer_Processing(hostname string) Result {
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

func GetHostWithImage_Processing(hostname string) Result {
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
	dockerImageList := Internal_ImageListMetricInfo(hostname)

	// alldata包含：主机基本配置，主机动态指标获取，镜像列表
	alldata["hostConfig"] = hostdata
	alldata["hostMetric"] = pureMetric
	alldata["imageList"] = dockerImageList.Data

	ResultData.Code = http.StatusOK
	ResultData.Data = alldata
	return ResultData
}

func Internal_AddHost(host *Host) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result

	_, err := o.Insert(host)
	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddHostErr
		logs.Error("AddHost failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = host
	return ResultData
}

func Internal_EditHost(host *Host) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result

	_, err := o.Update(host)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.EditHostErr
		logs.Error("EditHost failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = host
	return ResultData
}

func AddHost_Processing(h Host, isEdit int) interface{} {
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
		ResultData.Message = "Host Data cant acquire"
		logs.Error("AddHost failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	h.CpuKernel = pureMetric["cpu.cores"].(float64)
	h.Disk = pureMetric["filesystem.total"].(float64)
	h.Mem = pureMetric["memory.total"].(float64)

	// isEdit = 0, add host
	// isEdit = 1, edit host
	if isEdit == 0 {
		ResultData.Data = Internal_AddHost(&h)
	} else {
		ResultData.Data = Internal_EditHost(&h)
	}

	ResultData.Code = http.StatusOK

	return ResultData
}

func DeleteHost(id int) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
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
