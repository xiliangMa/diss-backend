package models

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/xiliangMa/diss-backend/utils"
	"strings"
)

type esqueryString map[string]string
func ESString (queryTag string) string{
	QueryDefine := make(esqueryString)
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
{"_source":["system.cpu.cores"],"sort": { "@timestamp": { "order": "desc" }},"size": 1,"query":{"bool":{"must": [{"term": {"event.module": "system"}},{"term": {"host.name": "c5b627e16af7"}},{"exists":{"field": "system.cpu.cores"}}]}}}
{"index":"metricbeat-*","ignore_unavailable":true,"preference":1573547985305}
{"_source":["system.memory.total"],"sort": { "@timestamp": { "order": "desc" }},"size": 1,"query":{"bool":{"must": [{"term": {"event.module": "system"}},{"term": {"host.name": "c5b627e16af7"}},{"exists":{"field": "system.memory.total"}}]}}}`
	QueryDefine["container_metric"] = `{
    "timeout": "30000ms",
    "aggs":
    {
        "container_name":
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
                            "host.name": "c5b627e16af7"
                        }
                    }],
                    "minimum_should_match": 1
                }
            },
            {
                "range":
                {
                    "@timestamp":
                    {
                        "format": "strict_date_optional_time",
                        "gte": "2019-06-14T09:29:17.970Z",
                        "lte": "2019-11-14T09:29:17.970Z"
                    }
                }
            }],
            "should": [],
            "must_not": []
        }
    }
}
`

	return QueryDefine[queryTag]
}

func GetHostMetricInfo_M(hostname string) interface{}{
	curhost := GetHostInternal(hostname)
	//esclient := utils.GetESClient()
	fmt.Print(curhost)

	//res, _ := esclient.API.Msearch(esclient.Msearch.WithContext(context.Background()),
	//	esclient.Msearch.WithIndex("metric*"),
	//	esclient.Msearch.WithBody(strings.NewReader(ESString("hostmetric_msearch"))),
	//	esclient.Search.WithTrackTotalHits(true),
	//	esclient.Search.WithPretty())
	//
	//var hostInfo map[string]interface{}
	//json.NewDecoder(res.Body).Decode(&hostInfo)
	//
	//hostInfoPure := hostInfo["aggregations"]

	return curhost
}

func GetHostMetricInfo(hostname string) interface{}{
	//curhost := GetHostInternal(hostname)
	esclient := utils.GetESClient()
	//fmt.Print(curhost)

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

func GetContainerListMetricInfo(hostname string) interface{}{
	//curhost := GetHostInternal(hostname)
	esclient := utils.GetESClient()
	//fmt.Print(curhost)

	res, _ := esclient.API.Search(esclient.Search.WithContext(context.Background()),
		esclient.Search.WithIndex("metric*"),
		esclient.Search.WithBody(strings.NewReader(ESString("container_metric"))),
		esclient.Search.WithTrackTotalHits(true),
		esclient.Search.WithPretty())

	var hostInfo map[string]interface{}
	json.NewDecoder(res.Body).Decode(&hostInfo)

	hostInfoPure := hostInfo["aggregations"]

	return hostInfoPure
}