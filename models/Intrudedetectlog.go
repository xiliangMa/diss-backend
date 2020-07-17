package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// 针对 es 解析使用的结构体
type IntrudeDetectLog struct {
	HostId      string `description:"(主机Id)"`
	HostName    string `description:"(主机名)"`
	TargeType   string `description:"(类型)"`
	ContainerId string `description:"(容日Id 如果是主机该字段为：host， 如果是容器为：容器的实际ID)"`
	Output      string `description:"(事件信息)"`
	StartTime   string `description:"(开始时间)"`
	ToTime      string `description:"(结束时间)"`
	AccountName string `description:"(租户)"`
	Priority    string `description:"(安全等级)"`
}

// 入侵检测日志（IntrudeDetectLog） 保存于 timescaledb
type DcokerIds struct {
	HostId       string    `orm:"pk" description:"(主机id)"`
	HostName     string    `description:"(主机名)"`
	MachineId    string    `description:"(Machine_id)"`
	ContainerId  string    `description:"(容器id)"`
	Time         time.Time `description:"(日志生成时间)"`
	Priority     string    `description:"(安全等级)"`
	Rule         string    `description:"(规则)"`
	Output       string    `description:"(事件信息)"`
	OutputFields string    `description:"(Output json)"`
	CreatedAt    int       `description:"(日志保存时间)"`
}

type DcokerIdsInterface interface {
	List(from, limit int) Result
	List1(from, limit int) Result
	GetIntrudeDetectLogStatistics(timeCycle int) Result
}

func (this *DcokerIds) GetIntrudeDetectLogStatistics(timeCycle int) Result {
	o := orm.NewOrm()

	o.Using(utils.DS_Security_Log)
	dcokerIdsCountList := make(map[int]*int)
	var ResultData Result
	var err error
	now := time.Now()
	for i := 0; i < timeCycle; i++ {
		count := 0
		m := i
		c1, _ := time.ParseDuration("-1h")
		c2, _ := time.ParseDuration("-" + strconv.Itoa(m+1) + "h")
		hStart := now.Add(c2).Unix()
		hEnd := now.Unix()
		if i != 0 {
			c1, _ = time.ParseDuration("-" + strconv.Itoa(i) + "h")
			hEnd = now.Add(c1).Unix()
		}
		err = o.Raw("select count(host_id) as count from docker_ids where created_at > ? and created_at < ?", hStart, hEnd).QueryRow(&count)
		if err != nil {
			ResultData.Message = err.Error()
			ResultData.Code = utils.GetIntrudeDetectLogErr
			logs.Error("Get IntrudeDetectLog List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
			return ResultData
		}
		dcokerIdsCountList[i] = &count
	}
	//QueryRows(&dcokerIdsList)
	//_, err = o.QueryTable(utils.DcokerIds).SetCond(cond).Limit(1, from).OrderBy("-time").All(&dcokerIdsList)

	data := make(map[string]interface{})
	data["items"] = dcokerIdsCountList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	return ResultData
}

func (this *IntrudeDetectLog) List(from, limit int) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Security_Log)
	var dcokerIdsList []*DcokerIds = nil
	var ResultData Result
	var err error
	containerId := this.ContainerId
	//cond := orm.NewCondition()
	st, _ := time.ParseInLocation("2006-01-02T15:04:05", this.StartTime, time.UTC)
	tt, _ := time.ParseInLocation("2006-01-02T15:04:05", this.ToTime, time.UTC)

	var total int64
	if this.TargeType == IDLT_Docker {
		containerId = string([]byte(this.ContainerId)[:12])
		total, err = o.Raw("select * from docker_ids where container_id = ? and created_at > ? and created_at < ? order by created_at desc limit ? OFFSET ?", containerId, st.Unix(), tt.Unix(), limit, from).QueryRows(&dcokerIdsList)
	} else {
		total, err = o.Raw("select * from docker_ids where host_id = ? and container_id = ? and created_at > ? and created_at < ? order by created_at desc limit ? OFFSET ?", this.HostId, containerId, st.Unix(), tt.Unix(), limit, from).QueryRows(&dcokerIdsList)
	}

	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetIntrudeDetectLogErr
		logs.Error("Get IntrudeDetectLo List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}

	//total, _ := o.QueryTable(utils.DcokerIds).SetCond(cond).Count()
	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = dcokerIdsList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}

func (this *IntrudeDetectLog) List1(from, limit int) Result {
	o := orm.NewOrm()
	o.Using(utils.DS_Security_Log)
	var dcokerIdsList []*DcokerIds = nil
	var ResultData Result
	var err error
	var total int64 = 0

	filterSql := ""
	countSql := "select " + `"count"(host_id)` + " from " + utils.DcokerIds
	sql := "select * from " + utils.DcokerIds

	// 根据 TargeType = host 如果快速查询所有主机日志可以设置 HostId=All
	// 根据 TargeType = container 如果快速查询所有容器日志可以设置 ContianerId =All

	if this.TargeType == IDLT_Host {
		filterSql = filterSql + "container_id = '" + IDLT_Host + "' and "

		if this.HostId != "" && this.HostId != All {
			filterSql = filterSql + "host_id = '" + this.HostId + "' and "
		}
		if this.HostName != "" {
			filterSql = filterSql + "host_name = '" + this.HostName + "' and "
		}
		if this.StartTime != "" {
			st, _ := time.ParseInLocation("2006-01-02T15:04:05", this.StartTime, time.UTC)
			filterSql = filterSql + "created_at > '" + strconv.FormatInt(st.Unix(), 10) + "' and "
		}
		if this.ToTime != "" {
			tt, _ := time.ParseInLocation("2006-01-02T15:04:05", this.ToTime, time.UTC)
			filterSql = filterSql + "created_at < '" + strconv.FormatInt(tt.Unix(), 10) + "' and "
		}
		if this.Priority != "" {
			filterSql = filterSql + "priority = '" + this.Priority + "' and "
		}
		if this.Output != "" {
			filterSql = filterSql + "output like '%" + this.Output + "%' and "
		}
	}

	if this.TargeType == IDLT_Docker {
		if this.ContainerId == All {
			filterSql = filterSql + "container_id != '" + IDLT_Host + "' and "
		} else {
			containerId := this.ContainerId
			// 如果是安全日志入口查询 不需要截取12位
			//containerId = string([]byte(this.ContainerId)[:12])
			filterSql = filterSql + "container_id = '" + containerId + "' and "
		}
		if this.HostId != "" {
			filterSql = filterSql + "host_id = '" + this.HostId + "' and "
		}
		if this.HostName != "" {
			filterSql = filterSql + "host_name = '" + this.HostName + "' and "
		}
		if this.StartTime != "" {
			st, _ := time.ParseInLocation("2006-01-02T15:04:05", this.StartTime, time.UTC)
			filterSql = filterSql + "created_at > '" + strconv.FormatInt(st.Unix(), 10) + "' and "
		}
		if this.ToTime != "" {
			tt, _ := time.ParseInLocation("2006-01-02T15:04:05", this.ToTime, time.UTC)
			filterSql = filterSql + "created_at < '" + strconv.FormatInt(tt.Unix(), 10) + "' and "
		}
		if this.Priority != "" {
			filterSql = filterSql + "priority = '" + this.Priority + "' and "
		}
		if this.Output != "" {
			filterSql = filterSql + "output like '%" + this.Output + "%' and "
		}
	}

	//if (this.ContainerId != "" && this.ContainerId != models.All) || (this.HostId != "" && this.HostId != models.All) {
	//	if this.ContainerId != models.IDLT_Host && this.TargeType == models.IDLT_Docker {
	//		containerId := this.ContainerId
	//		// 如果是安全日志入口查询 不需要截取12位
	//		//containerId = string([]byte(this.ContainerId)[:12])
	//		sql = sql + "container_id = '" + containerId + "' and "
	//	}
	//
	//	if this.TargeType == models.IDLT_Host {
	//		sql = sql + "container_id = '" + models.IDLT_Host + "' and "
	//	}
	//	if this.HostId != "" {
	//		sql = sql + "host_id = '" + this.HostId + "' and "
	//	}
	//	if this.HostName != "" {
	//		sql = sql + "host_name = '" + this.HostName + "' and "
	//	}
	//	if this.StartTime != "" {
	//		st, _ := time.ParseInLocation("2006-01-02T15:04:05", this.StartTime, time.UTC)
	//		sql = sql + "created_at > '" + strconv.FormatInt(st.Unix(), 10) + "' and "
	//	}
	//	if this.ToTime != "" {
	//		tt, _ := time.ParseInLocation("2006-01-02T15:04:05", this.ToTime, time.UTC)
	//		sql = sql + "created_at < '" + strconv.FormatInt(tt.Unix(), 10) + "' and "
	//	}
	//	if this.Priority != "" {
	//		sql = sql + "priority = '" + this.Priority + "' and "
	//	}
	//}

	if filterSql != "" {
		sql = sql + ` where ` + filterSql
		countSql = countSql + ` where ` + filterSql
	}
	countSql = strings.TrimSuffix(strings.TrimSpace(countSql), "and")
	sql = strings.TrimSuffix(strings.TrimSpace(sql), "and")
	resultSql := sql
	if from >= 0 && limit > 0 {
		limitSql := " limit " + strconv.Itoa(limit) + " OFFSET " + strconv.Itoa(from)
		resultSql = resultSql + limitSql
	}

	_, err = o.Raw(resultSql).QueryRows(&dcokerIdsList)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetIntrudeDetectLogErr
		logs.Error("Get IntrudeDetectLo List failed, code: %d, err: %s", ResultData.Code, ResultData.Message)
		return ResultData
	}
	o.Raw(countSql).QueryRow(&total)
	data := make(map[string]interface{})
	data["total"] = total
	data["items"] = dcokerIdsList

	ResultData.Code = http.StatusOK
	ResultData.Data = data
	if total == 0 {
		ResultData.Data = nil
	}
	return ResultData
}
