package models

const (
	// #################### ws 通信资源
	Resource_HostConfigDynamic      = "HostConfigDynamic"
	Resource_HostInfoDynamic        = "HostInfoDynamic"
	Resource_HostConfig             = "HostConfig"
	Resource_HostInfo               = "HostInfo"
	Resource_ImageConfig            = "ImageConfig"
	Resource_ImageInfo              = "ImageInfo"
	Resource_ContainerConfig        = "ContainerConfig"
	Resource_ContainerInfo          = "ContainerInfo"
	Resource_ContainerPs            = "ContainerPs"
	Resource_HostPs                 = "HostPs"
	Resource_DockerBenchMark        = "DockerBenchMark"
	Resource_KubernetesBenchMark    = "KubernetesBenchMark"
	Resource_HostCmdHistory         = "HostCmdHistory"
	Resource_ContainerCmdHistory    = "ContainerCmdHistory"
	Resource_NameSpace              = "NameSpace"
	Resource_Pod                    = "Pod"
	Resource_HeartBeat              = "HeartBeat"
	Resource_Received               = "Received"
	Resource_Task                   = "Task"
	Resource_DockerEvent            = "DockerEvent"
	Resource_CmdHistory_LatestTime  = "CmdHistory_LatestTime"
	Resource_DockerEvent_LatestTime = "DockerEvent_LatestTime"
	Resource_WarningInfo            = "WarningInfo"
	Resource_HostPackage            = "HostPackage"

	// #################### ws 通信类型
	Type_Metric       = "Metric"       //指标类型
	Type_ReceiveState = "ReceiveState" //接收状态类型
	Type_RequestState = "RequestState" //请求状态类型
	Type_Control      = "Control"      //控制类型

	// #################### ws ResourceControlType 资源操作类型
	Resource_Control_Type_Get    = "Get"
	Resource_Control_Type_Post   = "Post"
	Resource_Control_Type_Put    = "Put"
	Resource_Control_Type_Delete = "Delete"
)

type WsData struct {
	Code   int         `json:"code"`
	Msg    string      `json:"msg"`
	Type   string      `json:"type"`    // 资源类型：如指标数据metric，统计数据statistics， 配置信息config
	Tag    string      `json:"tag"`     // 资源
	RCType string      `json:"rc_type"` //ResourceControlType 资源操作类型
	Data   interface{} `json:"data"`    // 指标数据：具体数据内容，可以是聚合型多组
	Config string      `json:"config"`  // 配置相关数据
}
