package models

const (
	// #################### Nats 通信资源
	Resource_HostConfigDynamic      = "HostConfigDynamic"      // 主机配置数据
	Resource_HostInfoDynamic        = "HostInfoDynamic"        // 主机详细数据
	Resource_HostConfig             = "HostConfig"             // 主机配置据
	Resource_HostInfo               = "HostInfo"               // 主机详细数据
	Resource_ImageConfig            = "ImageConfig"            // 镜像配置据
	Resource_ImageInfo              = "ImageInfo"              // 镜像详细数据
	Resource_ImageDetail            = "ImageDetail"            // 镜像详情
	Resource_ContainerConfig        = "ContainerConfig"        // 容器配置数据
	Resource_ContainerInfo          = "ContainerInfo"          // 容器详细数据
	Resource_ContainerPs            = "ContainerPs"            // 容器进程数据
	Resource_HostPs                 = "HostPs"                 // 主机进程数据
	Resource_DockerBenchMark        = "DockerBenchMark"        // docker基线数据
	Resource_KubernetesBenchMark    = "KubernetesBenchMark"    // k8s 基线数据
	Resource_HostCmdHistory         = "HostCmdHistory"         // 主机命令历史数据
	Resource_ContainerCmdHistory    = "ContainerCmdHistory"    // 容器命令历史数据
	Resource_HostImageVulnScan      = "HostImageVulnScan"      // 主机镜像扫描
	Resource_NameSpace              = "NameSpace"              // k8s 命名空间
	Resource_Pod                    = "Pod"                    // k8s Pod
	Resource_HeartBeat              = "HeartBeat"              // 心跳
	Resource_Received               = "Received"               // 已接受
	Resource_Task                   = "Task"                   // 任务
	Resource_DockerEvent            = "DockerEvent"            // docker 审计
	Resource_CmdHistory_LatestTime  = "CmdHistory_LatestTime"  // 命令历史最新上报时间
	Resource_DockerEvent_LatestTime = "DockerEvent_LatestTime" // docker 审计最新上报时间
	Resource_WarningInfo            = "WarningInfo"            // 告警
	Resource_HostPackage            = "HostPackage"            // 主机包
	Resource_ClientModuleControl    = "ClientModuleControl"    // 客户端模块控制
	Resource_ContainerVS            = "ContainerVirusScan"     // 容器杀毒数据
	Resource_HostVS                 = "HostVirusScan"          // 主机杀毒数据
	Resource_ImageVS                = "ImageVirusScan"         // 镜像杀毒数据

	// #################### （Type）Nats Type 通信类型
	Type_Metric                     = "Metric"                 //指标类型
	Type_ReceiveState               = "ReceiveState"           //接收状态类型
	Type_RequestState               = "RequestState"           //请求状态类型
	Type_Control                    = "Control"                //控制类型
	Resource_ContainerControl       = "ContainerControl"       // 容器控制
	Resource_ContainerControlStatus = "ContainerControlStatus" // 响应中心状态
	Image_Control                   = "ImageControl"           // 镜像阻断控制
	Image_ControlStatus             = "ImageControlStatus"     // 镜像阻断控制状态

	// #################### （ResourceControlType）Nats RCType  资源操作类型
	Resource_Control_Type_Get    = "Get"
	Resource_Control_Type_Post   = "Post"
	Resource_Control_Type_Put    = "Put"
	Resource_Control_Type_Delete = "Delete"
)

type NatsData struct {
	Code   int         `json:"code"`    // 状态码
	Msg    string      `json:"msg"`     // 响应消息/错误信息
	Type   string      `json:"type"`    // nats 通信类型：如指标数据 metric，统计数据 statistics，控制数据：Control 配置信息config
	Tag    string      `json:"tag"`     // 资源
	RCType string      `json:"rc_type"` // ResourceControlType 资源操作类型 Get Post Delete Put
	Data   interface{} `json:"data"`    // 数据
	Config string      `json:"config"`  // 配置信息
}
