package models

type MetricsResult struct {
	ResType string      `json:"res_type"` // 资源类型：如指标数据metric，统计数据statistics， 配置信息config
	ResTag  string      `json:"res_tag"`  // 资源标记：如hostinfo
	Metric  interface{} `json:"metric"`   // 指标数据：具体数据内容，可以是聚合型多组
	Config  string      `json:"config"`   // 配置相关数据
}
