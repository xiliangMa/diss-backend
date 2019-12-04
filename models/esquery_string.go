package models

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



