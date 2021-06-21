package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"time"
)

type Module struct {
	Id         int    `orm:"auto;pk" description:"模块ID"`
	Name       string `orm:"unique" description:"模块名"`
	Code       string `description:"模块密码"`
	Role       *Role  `orm:"rel(fk);null" description:"角色"`
	CreateTime int64  `orm:"default(0);" description:"(创建时间)"`
	UpdateTime int64  `orm:"default(0)" description:"(更新时间)"`
}

type ModuleInterface interface {
	List(from, limit int) Result
	Add() Result
	Update() Result
	Delete() Result
	ModuleList() ([]*Module, int64, error)
}

func (this *Module) Add() Result {
	var ResultData Result
	o := orm.NewOrm()
	o.Using(utils.DS_Default)

	this.CreateTime = time.Now().UnixNano()
	_, err := o.Insert(this)
	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddModuleErr
		logs.Error("Add Module failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *Module) Update() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	moduleObj := Module{}
	moduleObj.Id = this.Id

	moduleList, total, _ := moduleObj.ModuleList(0, 0)
	if total > 0 {
		roleData := moduleList[0]
		this.CreateTime = roleData.CreateTime
		this.UpdateTime = time.Now().UnixNano()
		_, err := o.Update(this)
		if err != nil {
			ResultData.Message = err.Error()
			ResultData.Code = utils.EditModuleErr
			logs.Error("Edit Module %s failed, code: %d, err: %s", this.Name, ResultData.Code, ResultData.Message)
			return ResultData
		}

	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *Module) List(from, limit int) Result {
	var ResultData Result

	roleList, total, err := this.ModuleList(from, limit)

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetModuleErr
		logs.Error("Get Module failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
	}

	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = roleList
	if total == 0 {
		ResultData.Data = nil
	}

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	return ResultData
}

func (this *Module) ModuleList(from, limit int) (moduleLists []*Module, count int64, err error) {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var moduleList []*Module = nil
	cond := orm.NewCondition()

	if this.Id != 0 {
		cond = cond.And("id", this.Id)
	}
	if this.Name != "" {
		cond = cond.And("Name__contains", this.Name)
	}

	_, err = o.QueryTable(utils.Module).SetCond(cond).Limit(limit, from).OrderBy("-create_time").All(&moduleList)

	total, _ := o.QueryTable(utils.Module).SetCond(cond).Count()
	return moduleList, total, err
}

func (this *Module) Delete() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result
	cond := orm.NewCondition()

	if this.Id != 0 {
		cond = cond.And("id", this.Id)
	} else {
		ResultData.Message = "No ModuleList Id."
		ResultData.Code = utils.DeleteModuleErr
		return ResultData
	}

	_, err := o.QueryTable(utils.Module).SetCond(cond).Delete()

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.DeleteModuleErr
		logs.Error("Delete Module failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	return ResultData
}
