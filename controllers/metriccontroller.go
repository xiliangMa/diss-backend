package controllers

import (
	"context"
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
// @Success 200 {object} mertric.HostBasicInfo
// @router /hostinfo [post]
func (this *MetricController) HostInfo() {
	hostname := this.GetString("hostname")
	curhost := models.GetHostInternal(hostname)
	esclient := utils.GetESClient()
	fmt.Print(curhost)
	data := make(map[string]interface{})

	res, _ := esclient.API.Search(esclient.Search.WithContext(context.Background()),
		esclient.Search.WithIndex("metric*"),
		esclient.Search.WithBody(strings.NewReader(`{
    "timeout": "30000ms",
    "query":
    {
        "bool":
        {
            "must": [],
            "filter": [
            {
                "bool":
                {
                    "should": [
                    {
                        "match_phrase":
                        {
                            "host.name": "c5b627e16af7"
                        }
                    }],
                    "minimum_should_match": 1
                }
            },
            {
                "bool":
                {
                    "should": [
                    {
                        "match_phrase":
                        {
                            "event.module": "system"
                        }
                    }],
                    "minimum_should_match": 1
                }
            },{
                "range":
                {
                    "@timestamp":
                    {
                        "format": "strict_date_optional_time",
                        "gte": "2019-11-15T00:29:17.970Z",
                        "lte": "2019-11-15T12:29:17.970Z"
                    }
                }
            }
            ],
            "should": [],
            "must_not": []
        }
    },
    "size": 0,

    "aggs":{
    	"event_dataset_list":{
    		"terms":{
    			"field":"event.dataset",
    			"size" :15
    		},
    		"aggs" : {
    			"metric_data" :{
	    			"top_hits" :{
	    				"sort": [
	                            {
	                                "@timestamp": {
	                                    "order": "desc"
	                                }
	                            }
	                        ],
	                    "size" :1    
	    			}
    			}
    		}
    	}
    }
}`)),
		esclient.Search.WithTrackTotalHits(true),
		esclient.Search.WithPretty())

	data["hostinfo"] = res

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
