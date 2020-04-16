package ws

const (
	// #################### ws 通信 Tag
	Tag_HostConfig             = "HostConfig"
	Tag_HostInfo               = "HostInfo"
	Tag_ImageConfig            = "ImageConfig"
	Tag_ImageInfo              = "ImageInfo"
	Tag_ContainerConfig        = "ContainerConfig"
	Tag_ContainerInfo          = "ContainerInfo"
	Tag_ContainerPs            = "ContainerPs"
	Tag_HostPs                 = "HostPs"
	Tag_DockerBenchMarkLog     = "DockerBenchMarkLog"
	Tag_KubernetesBenchMarkLog = "KubernetesBenchMarkLog"
	Tag_HostCmdHistory         = "HostCmdHistory"
	Tag_ContainerCmdHistory    = "ContainerCmdHistory"
	Tag_NameSpace              = "NameSpace"
	Tag_Pod                    = "Pod"
	Tag_HeartBeat              = "HeartBeat"
	Tag_Received               = "Received"
	Tag_Task                   = "Task"

	// diss-backend 下发的数据
	Tag_DockerBenchMark     = "DockerBenchMark"
	Tag_KubernetesBenchMark = "KubernetesBenchMark"

	// #################### ws 通信类型
	Type_Metric       = "Metric"
	Type_ReceiveState = "ReceiveState"
	Type_SyncData     = "SyncData"
	Type_Response     = "Response"
)

type MetricsResult struct {
	ResType string      `json:"res_type"` // 资源类型：如指标数据metric，统计数据statistics， 配置信息config
	ResTag  string      `json:"res_tag"`  // 资源标记：如hostinfo
	Metric  interface{} `json:"metric"`   // 指标数据：具体数据内容，可以是聚合型多组
	Config  string      `json:"config"`   // 配置相关数据
}
