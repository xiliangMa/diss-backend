package models

// ws 通信 Tag
const (
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
)
