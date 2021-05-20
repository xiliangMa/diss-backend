package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"strings"
	"time"
)

type RuleDefine struct {
	Id         int64  `orm:"pk;auto" description:"(Id)"`
	RuleType   string `orm:"size(64)" description:"(规则类型：如敏感信息、告警信息)"`
	RuleName   string `orm:"size(256)" description:"(规则名称)"`
	RuleCode   string `orm:"size(64)" description:"(规则Code)"`
	Desc       string `orm:"" description:"(规则描述)"`
	RuleInfo   string `orm:"size(256)" description:"(规则定义)"`
	Enabled    bool   `orm:"" description:"(是否启用)"`
	IsAll      bool   `orm:"-" description:"(是否获取全部)"`
	SourceInfo string `orm:"size(256)" description:"(来源信息:来源的网址、漏洞编码等)" `
	CreateTime int64  `orm:"" description:"(添加时间)"`
	UpdateTime int64  `orm:"" description:"(修改时间)"`
	RiskLevel  string `orm:"size(64)" description:"(风险等级)"`
}

type RuleDefineInterface interface {
	List(from, limit int) Result
	Add() Result
	Update() Result
	Delete() Result
	RuleDefineList() ([]*RuleDefine, int64, error)
}

func (this *RuleDefine) Add() Result {
	var ResultData Result
	o := orm.NewOrm()
	o.Using(utils.DS_Default)

	this.CreateTime = time.Now().UnixNano()
	_, err := o.Insert(this)
	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddRuleDefineErr
		logs.Error("Add RuleDefine failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *RuleDefine) Update() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	ruleDefineObj := RuleDefine{}
	ruleDefineObj.Id = this.Id
	ruleDefineObj.IsAll = true

	ruleList, total, _ := ruleDefineObj.RuleDefineList(0, 0)
	if total > 0 {
		ruleData := ruleList[0]
		this.CreateTime = ruleData.CreateTime
		this.UpdateTime = time.Now().UnixNano()
		_, err := o.Update(this)
		if err != nil {
			ResultData.Message = err.Error()
			ResultData.Code = utils.EditRuleDefineErr
			logs.Error("Edit RuleDefine %s failed, code: %d, err: %s", this.RuleName, ResultData.Code, ResultData.Message)
			return ResultData
		}

	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *RuleDefine) List(from, limit int) Result {
	var ResultData Result

	ruleList, total, err := this.RuleDefineList(from, limit)

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetRuleDefineErr
		logs.Error("Get RuleDefine failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
	}

	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = ruleList
	if total == 0 {
		ResultData.Data = nil
	}

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	return ResultData
}

func (this *RuleDefine) RuleDefineList(from, limit int) (ruleLists []*RuleDefine, count int64, err error) {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ruleList []*RuleDefine = nil
	cond := orm.NewCondition()

	if this.Id != 0 {
		cond = cond.And("id", this.Id)
	}
	if this.RuleType != "" {
		ruleType := strings.Split(this.RuleType, ",")
		cond = cond.And("RuleType__in", ruleType)
	}
	if this.RuleName != "" {
		cond = cond.And("RuleName__contains", this.RuleName)
	}
	if this.RuleCode != "" {
		cond = cond.And("RuleCode__contains", this.RuleCode)
	}
	if this.RuleInfo != "" {
		cond = cond.And("RuleInfo__contains", this.RuleInfo)
	}
	if !this.IsAll {
		cond = cond.And("enabled", this.Enabled)
	}

	_, err = o.QueryTable(utils.RuleDefine).SetCond(cond).Limit(limit, from).OrderBy("-create_time").All(&ruleList)

	total, _ := o.QueryTable(utils.RuleDefine).SetCond(cond).Count()
	return ruleList, total, err
}

func (this *RuleDefine) Delete() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	cond := orm.NewCondition()

	if this.Id != 0 {
		cond = cond.And("id", this.Id)
	} else {
		ResultData.Message = "No RuleDefineList Id."
		ResultData.Code = utils.DeleteRuleDefineErr
		return ResultData
	}

	_, err := o.QueryTable(utils.RuleDefine).SetCond(cond).Delete()

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.DeleteRuleDefineErr
		logs.Error("Delete RuleDefine failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	return ResultData
}
