package models

import (
	"context"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"strings"
)

type ESQueryString map[string]string

func ESString(queryTag string) string {
	QueryDefine := make(ESQueryString)
	QueryDefine["host_metric"] = `{
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
	    				"_source":
						    {
						        "includes": ["system","process"]
						    },
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
}`
	QueryDefine["msearch_host_metric"] = `{"index":"metricbeat-*","ignore_unavailable":true,"preference":1573547985305}
				{"_source":["system.cpu.cores"],"sort": { "@timestamp": { "order": "desc" }},"size": 1,"query":{"bool":{"must": [{"term": {"event.module": "system"}},{"term": {"host.name": "!Param@hostname!"}},{"exists":{"field": "system.cpu.cores"}}]}}}
				{"index":"metricbeat-*","ignore_unavailable":true,"preference":1573547985305}
				{"_source":["system.memory.total"],"sort": { "@timestamp": { "order": "desc" }},"size": 1,"query":{"bool":{"must": [{"term": {"event.module": "system"}},{"term": {"host.name": "!Param@hostname!"}},{"exists":{"field": "system.memory.total"}}]}}}
				{"index":"metricbeat-*","ignore_unavailable":true,"preference":1573547985305}
{				"_source":["system.filesystem.total"],"sort": { "@timestamp": { "order": "desc" }},"size": 1,"query":{"bool":{"must": [{"term": {"event.module": "system"}},{"term": {"host.name": "!Param@hostname!"}},{"exists":{"field": "system.filesystem.total"}}]}}}
`
	QueryDefine["container_metric"] = `{
    "timeout": "30000ms",
    "aggs":
    {
        "container_list":
        {
            "terms":
            {
                "field": "container.name",
                "order":
                {
                    "1": "desc"
                },
                "size": 5
            },
            "aggs":
            {
                "1":
                {
                    "cardinality":
                    {
                        "field": "container.id"
                    }
                },
                "cpu_total_pct":
                {
                    "max":
                    {
                        "field": "docker.cpu.total.pct"
                    }
                },
                "diskio_total":
                {
                    "max":
                    {
                        "field": "docker.diskio.total"
                    }
                },
                "memory_usepct":
                {
                    "max":
                    {
                        "field": "docker.memory.usage.pct"
                    }
                },
                "memory_rsstotal":
                {
                    "max":
                    {
                        "field": "docker.memory.rss.total"
                    }
                }
            }
        }
    },
    "size": 0,
    "_source":
    {
        "excludes": []
    },
    "stored_fields": ["*"],
    "sort": { "@timestamp": { "order": "desc" }},
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
                        "match":
                        {
                            "event.module": "docker"
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
                        "match":
                        {
                            "event.module": "docker"
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
                            "host.name": "!Param@hostname!"
                        }
                    }],
                    "minimum_should_match": 1
                }
            }],
            "should": [],
            "must_not": []
        }
    }
}
`
	QueryDefine["container_summary"] = `{
    "timeout": "30000ms",
    "aggs":
    {
        "running":
        {
            "max":
            {
                "field": "docker.info.containers.running"
            }
        },
        "paused":
        {
            "max":
            {
                "field": "docker.info.containers.paused"
            }
        },
        "stopped":
        {
            "max":
            {
                "field": "docker.info.containers.stopped"
            }
        }
    },
    "size": 0,
    "_source":
    {
        "excludes": []
    },
    "stored_fields": ["*"],
	"sort": { "@timestamp": { "order": "desc" }},
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
                        "match":
                        {
                            "event.module": "docker"
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
                            "host.name": "!Param@hostname!"
                        }
                    }],
                    "minimum_should_match": 1
                }
            }],
            "should": [],
            "must_not": []
        }
    }
}`

	return QueryDefine[queryTag]
}

func Internal_HostMetricInfo_M(hostname string) Result {
	var ResultData Result
	esclient := utils.GetESClient()

	if _, errConnect := esclient.Info(); errConnect != nil {
		logs.Error("esclient connect error:", errConnect.Error())
		ResultData.Message = errConnect.Error()
		ResultData.Code = utils.ElasticConnErr
		return ResultData
	}

	esqueryStr := strings.Replace(ESString("msearch_host_metric"), "!Param@hostname!", hostname, 4)
	mres, errSearch := esclient.API.Msearch(strings.NewReader(esqueryStr), esclient.Msearch.WithIndex("metric*"))

	if errSearch != nil {
		logs.Error("host msearch error: ", errSearch.Error())
		ResultData.Message = errSearch.Error()
		ResultData.Code = utils.ElasticSearchErr
		ResultData.Data = nil
		return ResultData
	}

	logs.Info("host metric info is ok , %s", mres.Status())

	var hostInfo map[string]interface{}
	json.NewDecoder(mres.Body).Decode(&hostInfo)

	var hostInfoPure []interface{}
	for _, x := range hostInfo["responses"].([]interface{}) {
		hostInfoRefine1 := x.(map[string]interface{})["hits"]
		hostInfoRefine2 := hostInfoRefine1.(map[string]interface{})["hits"]
		hostInfoRefine2a := hostInfoRefine2.([]interface{})
		if len(hostInfoRefine2a) > 0 {
			hostInfoRefine3 := hostInfoRefine2a[0].(map[string]interface{})["_source"]
			hostInfoPure = append(hostInfoPure, hostInfoRefine3)
		}
	}

	ResultData.Data = hostInfoPure
	ResultData.Code = http.StatusOK
	return ResultData
}

// -- 因为结构复杂不使用的方法 ，通过 search - aggregation 方式获取es数据
func GetHostMetricInfo(hostname string) interface{} {
	var ResultData Result

	esclient := utils.GetESClient()
	if _, err := esclient.Info(); err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.ElasticConnErr
		return ResultData
	}
	res, _ := esclient.API.Search(esclient.Search.WithContext(context.Background()),
		esclient.Search.WithIndex("metric*"),
		esclient.Search.WithBody(strings.NewReader(ESString("host_metric"))),
		esclient.Search.WithTrackTotalHits(true),
		esclient.Search.WithPretty())

	var hostInfo map[string]interface{}
	json.NewDecoder(res.Body).Decode(&hostInfo)

	hostInfoPure := hostInfo["aggregations"]

	return hostInfoPure
}

func Internal_ContainerListMetricInfo(hostname string) Result {
	var ResultData Result

	esclient := utils.GetESClient()
	if _, err := esclient.Info(); err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.ElasticConnErr
		return ResultData
	}
	esqueryStr := strings.Replace(ESString("container_metric"), "!Param@hostname!", hostname, 1)
	res, errSearch := esclient.API.Search(esclient.Search.WithContext(context.Background()),
		esclient.Search.WithIndex("metric*"),
		esclient.Search.WithBody(strings.NewReader(esqueryStr)),
		esclient.Search.WithTrackTotalHits(true),
		esclient.Search.WithPretty())
	if errSearch != nil {
		logs.Error("host msearch error: ", errSearch.Error())
		ResultData.Message = errSearch.Error()
		ResultData.Code = utils.ElasticSearchErr
		ResultData.Data = nil
		return ResultData
	}

	var containerInfo map[string]interface{}
	json.NewDecoder(res.Body).Decode(&containerInfo)

	var hostInfoPure []interface{}
	containerInfoRefine1 := containerInfo["aggregations"].(map[string]interface{})
	containerInfoRefine2 := containerInfoRefine1["container_list"].(map[string]interface{})
	for _, x := range containerInfoRefine2["buckets"].([]interface{}) {
		hostInfoPure = append(hostInfoPure, x)
	}

	ResultData.Data = hostInfoPure
	ResultData.Code = http.StatusOK
	return ResultData
}

func Internal_ContainerSummaryInfo(hostname string) Result {
	var ResultData Result

	esclient := utils.GetESClient()
	if _, err := esclient.Info(); err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.ElasticConnErr
		return ResultData
	}

	esqueryStr := strings.Replace(ESString("container_summary"), "!Param@hostname!", hostname, 1)
	res, errSearch := esclient.API.Search(esclient.Search.WithContext(context.Background()),
		esclient.Search.WithIndex("metric*"),
		esclient.Search.WithBody(strings.NewReader(esqueryStr)),
		esclient.Search.WithTrackTotalHits(true),
		esclient.Search.WithPretty())
	if errSearch != nil {
		logs.Error("host msearch error: ", errSearch.Error())
		ResultData.Message = errSearch.Error()
		ResultData.Code = utils.ElasticSearchErr
		ResultData.Data = nil
		return ResultData
	}

	var containerSummaryInfo map[string]interface{}
	json.NewDecoder(res.Body).Decode(&containerSummaryInfo)

	containerInfoPure := containerSummaryInfo["aggregations"].(map[string]interface{})

	ResultData.Data = containerInfoPure
	ResultData.Code = http.StatusOK
	return ResultData
}
