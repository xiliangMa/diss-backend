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

func (this *TimeEdgePoint) AddTEPoint() Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Default)
	var ResultData Result

	uid, _ := uuid.NewV4()
	this.TEPointId = uid.String()
	_, err := o.Insert(this)
	if err != nil && utils.IgnoreLastInsertIdErrForPostgres(err) != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.AddTimeEdgePointErr
		logs.Error("Add LogConfig failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	ResultData.Code = http.StatusOK
	ResultData.Data = this
	return ResultData
}
