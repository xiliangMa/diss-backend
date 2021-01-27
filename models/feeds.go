package models

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"strconv"
	"strings"
)

type Feeds struct {
	Name         string `description:"(名称)"`
	Description  string `description:"(描述)"`
	AccessTier   string `description:"(AccessTier)"`
	LastFullSync int64  `description:"(最新同步时间)"`
	LastUpdate   int64  `description:"(更新时间)"`
	CreateAt     int64  `description:"(创建时间)"`
}

type FeedsInterface interface {
	List(from, limit int) Result
}

func (this *Feeds) List(from, limit int) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Diss_Api)
	var feedsList []*Feeds
	var ResultData Result
	var err error
	var total = 0

	sql := `SELECT * FROM ` + utils.Feeds
	countSql := `SELECT "count"(name) FROM ` + utils.Feeds
	filter := ""

	if this.Name != "" {
		filter = filter + ` name = '` + this.Name + `' and `

	}
	//if this.LastFullSync.String() != Null_Time {
	//	filter = filter + `last_full_sync   '` + this.LastFullSync.String() + `' and `
	//}
	if this.LastFullSync != 0 {
		filter = filter + `last_full_sync   ` + fmt.Sprintf("%v", this.LastFullSync) + ` and `
	}

	if filter != "" {
		sql = sql + " where " + filter
		countSql = countSql + " where " + filter
	}
	sql = strings.TrimSuffix(strings.TrimSpace(sql), "and")
	countSql = strings.TrimSuffix(strings.TrimSpace(countSql), "and")
	resultSql := sql
	if from >= 0 && limit > 0 {
		limitSql := " limit " + strconv.Itoa(limit) + " OFFSET " + strconv.Itoa(from)
		resultSql = resultSql + limitSql
	}

	_, err = o.Raw(resultSql).QueryRows(&feedsList)

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetFeedListErr
		logs.Error("Get feed List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	o.Raw(countSql).QueryRow(&total)
	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = feedsList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}
