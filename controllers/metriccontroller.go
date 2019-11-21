package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
	"strings"
)

// Hosts object api list
type MetricController struct {
	beego.Controller
}

// @Title GetHost
// @Description Get Hosts
// @Param token header string true "Auth token"
// @Param name query string false "host name"
// @Param ip query string false "host ip"
// @Param from query int 0 false "from"
// @Param limit query int 20 false "limit"
// @Success 200 {object} models.Result
// @router / [post]
func (this *MetricController) HostList() {
	name := this.GetString("name")
	ip := this.GetString("ip")
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")
	this.Data["json"] = models.GetHostList(name, ip, from, limit)
	this.ServeJSON(false)

}

// @Title HostInfo
// @Description HostMetricBasicInfo
// @Param token header string true "Auth token"
// @Param hostname query string false "Enter hostname"
// @Success 200 {object} models.Result
// @router /hostinfo [post]
func (this *MetricController) HostInfo() {
	hostname := this.GetString("hostname")
	curhost := models.GetHostInternal(hostname)
	esclient := utils.GetESClient()
	fmt.Print(curhost)
	data := make(map[string]interface{})

	res, _ := esclient.API.Search(esclient.Search.WithContext(context.Background()),
		esclient.Search.WithIndex("metric*"),
		esclient.Search.WithBody(strings.NewReader(`{"query" : { "match_all":{} }}`)),
		esclient.Search.WithTrackTotalHits(true),
		esclient.Search.WithPretty())

	var hostInfo map[string]interface{}
	json.NewDecoder(res.Body).Decode(&hostInfo)
	data["hostInfo"] = hostInfo
	this.Data["json"] = data
	this.ServeJSON(false)
}


// @Title GetHost
// @Description Get one Host
// @Param token header string true "Auth token"
// @Param hostname query string false "Enter hostname"
// @Success 200 {object} models.Result
// @router /gethost [post]
func (this *MetricController) GetHost(){
	hostname := this.GetString("hostname")
	this.Data["json"] = models.GetHost(hostname)
	this.ServeJSON(false)
}

// @Title DelHost
// @Description Delete Host
// @Param token header string true "Auth token"
// @Param id path int true "host id"
// @Success 200 {object} models.Result
// @router /:id [delete]
func (this *MetricController) DeleteHost() {
	id, _ := this.GetInt(":id")
	this.Data["json"] = models.DeleteHost(id)
	this.ServeJSON(false)

}
