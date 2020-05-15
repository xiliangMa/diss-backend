package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	uuid "github.com/satori/go.uuid"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
)

type TimeEdgePoint struct {
	TEPointId     string `orm:"pk;" description:"(时间边界点id)"`
	EdgePointCode string `orm:"" description:"(时间边界Code)"`
	EdgePointName string `orm:"" description:"(时间边界名称/说明)"`
	TimePointA    string `orm:"" description:"(时间边界A点)"`
	TimePointB    string `orm:"" description:"(时间边界B点)"`
	Direction     string `orm:"" description:"(时间区方向，回溯/范围/向前/展望/预期等)"`
	ScopeSymbol   string `orm:"" description:"(时间区域图示, 如----|)"`
}

type TimeEdgePointInterface interface {
	Add() Result
	Update() Result
	Get() []*TimeEdgePoint
}

func (this *TimeEdgePoint) Add() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result

	uid, _ := uuid.NewV4()
	this.TEPointId = uid.String()
	_, err := o.Insert(this)
	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddTimeEdgePointErr
		logs.Error("Add TimeEdgePoint failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *TimeEdgePoint) Update() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result

	_, err := o.Update(this)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.EditTimeEdgePointErr
		logs.Error("Update TimeEdgePoint: %s failed, code: %d, err: %s", this.EdgePointCode, ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}

func (this *TimeEdgePoint) Get() []*TimeEdgePoint {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var TEPoint []*TimeEdgePoint = nil
	var ResultData Result
	cond := orm.NewCondition()

	if this.EdgePointCode != "" {
		cond = cond.And("edge_point_code", this.EdgePointCode)
	}

	err := o.QueryTable(utils.TimeEdgePoint).SetCond(cond).One(&TEPoint)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetTimeEdgePointErr
		//logs.Info("Get TimeEdgePoint failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return nil
	}
	return TEPoint
}
