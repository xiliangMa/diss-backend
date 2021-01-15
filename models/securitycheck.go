package models

type SecurityCheckList struct {
	CheckList []*SecurityCheck
}

type SecurityCheck struct {
	DockerCIS         bool             `description:"(required: false, 开启Docker基线检测)"`
	KubenetesCIS      bool             `description:"(required: false, 开启K8s基线检测)"`
	VirusScan         bool             `description:"(required: false, 开启病毒)"`
	LeakScan          bool             `description:"(required: false, 开启漏洞)"`
	HostImageVulnScan bool             `description:"(required: false, 开启漏洞)"`
	Host              *HostConfig      `description:"(required: false, 主机)"`
	Container         *ContainerConfig `description:"(required: false, 容器)"`
	Image             *ImageConfig     `description:"(required: false, 镜像)"`
	Type              string           `description:"(required: false, 类型 host、container、image)"`
	Job               *Job             `description:"(required: false, 来源任务)"`
}

type SecurityCheckParams struct {
	DockerCIS         bool   `description:"(required: false, 开启Docker基线检测)"`
	KubenetesCIS      bool   `description:"(required: false, 开启K8s基线检测)"`
	VirusScan         bool   `description:"(required: false, 开启病毒)"`
	LeakScan          bool   `description:"(required: false, 开启漏洞)"`
	HostImageVulnScan bool   `description:"(required: false, 开启主机镜像漏洞扫描)"`
	HostIds           string `description:"(required: false, 主机id列表 ID1,ID2,ID3)"`
	ContainerIds      string `description:"(required: false, 容器id列表 ID1,ID2,ID3)"`
	ImageIds          string `description:"(required: false, 镜像id列表 ID1,ID2,ID3)"`
	Type              string `description:"(required: false, 类型 host、container、image)"`
	JobId             string `description:"(required: false, JobId)"`
}
