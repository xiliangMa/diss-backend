package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	uuid "github.com/satori/go.uuid"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"time"
)

type WarningWhiteList struct {
	Id                   string `orm:"pk;size(128)" description:"(Id)"`
	Name                 string `orm:"size(256)" description:"(白名单项名称)"`
	Desc                 string `orm:"" description:"(描述)"`
	WarningInfoType      string `orm:"size(64)" description:"(告警类型)"`
	WarningInfoName      string `orm:"size(64)" description:"(告警名称)"`
	RuleNode             string `orm:"" description:"(节点规则，如主机IP，主机名、容器名、容器ID等)"`
	Rule                 string `orm:"" description:"(规则)"`
	Enabled              bool   `orm:"" description:"(是否启用)"`
	IsAll                bool   `orm:"-" description:"(是否获取全部)"`
	WarningInfoId        string `orm:"size(128)" description:"(外键id)" `
	CreateTime           int64  `orm:"" description:"(创建时间)"`
	RuleNode_IP          string `orm:"-" description:"(节点信息-IP，虚拟字段)"`
	RuleNode_ContainerId string `orm:"-" description:"(节点信息-容器ID，虚拟字段)"`
}

type WarningWhiteListInterface interface {
	List(from, limit int) Result
	Add() Result
	Update() Result
	Delete() Result
	WhiteList() ([]*WarningWhiteList, int64, error)
}

func (this *WarningWhiteList) Add() Result {
	var ResultData Result
	o := orm.NewOrm()
	o.Using(utils.DS_Default)

	uid, _ := uuid.NewV4()
	this.Id = uid.String()
	this.CreateTime = time.Now().UnixNano()
	_, err := o.Insert(this)
	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddWarningWhiteListErr
		logs.Error("Add WarningWhiteList failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *WarningWhiteList) Update() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result

	_, err := o.Update(this)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.EditWarningWhiteListErr
		logs.Error("Edit WarningWhiteList: %s failed, code: %d, err: %s", this.Name, ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *WarningWhiteList) List(from, limit int) Result {
	var ResultData Result

	warnList, total, err := this.WhiteList(from, limit)

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetWarningWhiteListErr
		logs.Error("Get WarningWhiteList failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
	}

	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = warnList
	if total == 0 {
		ResultData.Data = nil
	}

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	return ResultData
}

func (this *WarningWhiteList) WhiteList(from, limit int) (whiteLists []*WarningWhiteList, count int64, err error) {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var warnList []*WarningWhiteList = nil
	cond := orm.NewCondition()

	if this.Id != "" {
		cond = cond.And("id", this.Id)
	}
	if this.WarningInfoType != "" {
		cond = cond.And("warning_info_type", this.WarningInfoType)
	}
	if this.WarningInfoName != "" {
		cond = cond.And("WarningInfoName__contains", this.WarningInfoName)
	}
	if this.Name != "" {
		cond = cond.And("Name__contains", this.Name)
	}
	if this.Rule != "" {
		cond = cond.And("Rule__contains", this.Rule)
	}

	if this.RuleNode_IP != "" {
		cond = cond.And("RuleNode__contains", this.RuleNode_IP)
	}

	if this.RuleNode_ContainerId != "" {
		cond = cond.And("RuleNode__contains", this.RuleNode_ContainerId)
	}

	if !this.IsAll {
		cond = cond.And("enabled", this.Enabled)
	}

	_, err = o.QueryTable(utils.WarningWhiteList).SetCond(cond).Limit(limit, from).OrderBy("-create_time").All(&warnList)

	total, _ := o.QueryTable(utils.WarningWhiteList).SetCond(cond).Count()
	return warnList, total, err
}

func (this *WarningWhiteList) Delete() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	cond := orm.NewCondition()

	if this.Id != "" {
		cond = cond.And("id", this.Id)
	} else {
		ResultData.Message = "No WarningWhiteList Id."
		ResultData.Code = utils.DeleteWarningWhiteListErr
		return ResultData
	}

	_, err := o.QueryTable(utils.WarningWhiteList).SetCond(cond).Delete()

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.DeleteWarningWhiteListErr
		logs.Error("Delete WarningWhiteList failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	return ResultData
}
