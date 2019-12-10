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
	QueryDefine["dockerimage_metric"] = `{
    "timeout": "30000ms",
    "aggs":
    {
        "image_list":
        {
            "terms":
            {
                "field": "container.image.name",
                "order":
                {
                    "_count": "desc"
                },
                "size": 15
            },
            "aggs":
            {
                "container":
                {
                    "terms":
                    {
                        "field": "container.name",
                        "order":
                        {
                            "_count": "desc"
                        },
                        "size": 12
                    }
                }
            }
        }
    },
    "size": 0,
    "_source":
    {
        "excludes": ["service", "agent", "ecs", "host"]
    },
    "stored_fields": ["*"],
    "script_fields":
    {},
    "docvalue_fields": [
    {
        "field": "@timestamp",
        "format": "date_time"
    },
    {
        "field": "ceph.monitor_health.last_updated",
        "format": "date_time"
    },
    {
        "field": "docker.container.created",
        "format": "date_time"
    },
    {
        "field": "docker.healthcheck.event.end_date",
        "format": "date_time"
    },
    {
        "field": "docker.healthcheck.event.start_date",
        "format": "date_time"
    },
    {
        "field": "docker.image.created",
        "format": "date_time"
    },
    {
        "field": "event.created",
        "format": "date_time"
    },
    {
        "field": "event.end",
        "format": "date_time"
    },
    {
        "field": "event.start",
        "format": "date_time"
    },
    {
        "field": "file.accessed",
        "format": "date_time"
    },
    {
        "field": "file.created",
        "format": "date_time"
    },
    {
        "field": "file.ctime",
        "format": "date_time"
    },
    {
        "field": "file.mtime",
        "format": "date_time"
    },
    {
        "field": "kubernetes.container.start_time",
        "format": "date_time"
    },
    {
        "field": "kubernetes.event.metadata.timestamp.created",
        "format": "date_time"
    },
    {
        "field": "kubernetes.event.timestamp.first_occurrence",
        "format": "date_time"
    },
    {
        "field": "kubernetes.event.timestamp.last_occurrence",
        "format": "date_time"
    },
    {
        "field": "kubernetes.node.start_time",
        "format": "date_time"
    },
    {
        "field": "kubernetes.pod.start_time",
        "format": "date_time"
    },
    {
        "field": "kubernetes.system.start_time",
        "format": "date_time"
    },
    {
        "field": "mongodb.replstatus.server_date",
        "format": "date_time"
    },
    {
        "field": "mongodb.status.background_flushing.last_finished",
        "format": "date_time"
    },
    {
        "field": "mongodb.status.local_time",
        "format": "date_time"
    },
    {
        "field": "mssql.transaction_log.stats.backup_time",
        "format": "date_time"
    },
    {
        "field": "nats.server.time",
        "format": "date_time"
    },
    {
        "field": "php_fpm.pool.start_time",
        "format": "date_time"
    },
    {
        "field": "php_fpm.process.start_time",
        "format": "date_time"
    },
    {
        "field": "postgresql.activity.backend_start",
        "format": "date_time"
    },
    {
        "field": "postgresql.activity.query_start",
        "format": "date_time"
    },
    {
        "field": "postgresql.activity.state_change",
        "format": "date_time"
    },
    {
        "field": "postgresql.activity.transaction_start",
        "format": "date_time"
    },
    {
        "field": "postgresql.bgwriter.stats_reset",
        "format": "date_time"
    },
    {
        "field": "postgresql.database.stats_reset",
        "format": "date_time"
    },
    {
        "field": "process.start",
        "format": "date_time"
    },
    {
        "field": "system.process.cpu.start_time",
        "format": "date_time"
    },
    {
        "field": "zookeeper.server.version_date",
        "format": "date_time"
    }],
    "query":
    {
        "bool":
        {
            "must": [],
            "filter": [
            {
                "match_all":
                {}
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
}`

	return QueryDefine[queryTag]
}



