package models

type SecurityCheckList struct {
	CheckList []*SecurityCheck
}

type SecurityCheck struct {
	DockerCIS         bool             `description:"(required: false, 开启Docker基线检测)"`
	KubenetesCIS      bool             `description:"(required: false, 开启K8s基线检测)"`
	DockerScan        bool             `description:"(required: false, 开启Docker漏扫)"`
	KubenetesScan     bool             `description:"(required: false, 开启K8s漏扫)"`
	VirusScan         bool             `description:"(required: false, 开启病毒)"`
	LeakScan          bool             `description:"(required: false, 开启漏洞)"`
	HostImageVulnScan bool             `description:"(required: false, 主机镜像漏洞)"`
	ImageVulnScan     bool             `description:"(required: false, 仓库镜像扫描)"`
	Host              *HostConfig      `description:"(required: false, 主机)"`
	Container         *ContainerConfig `description:"(required: false, 容器)"`
	Image             *ImageConfig     `description:"(required: false, 镜像)"`
	Cluster           *Cluster         `description:"(required: false, 集群)"`
	Type              string           `description:"(required: false, 类型 host、container、image)"`
	Job               *Job             `description:"(required: false, 来源任务)"`
	PathList          string           `description:"(required: false, 目标目录列表 Path1,Path2,Path3)"`
}

type SecurityCheckParams struct {
	DockerCIS         bool   `description:"(required: false, 开启Docker基线检测)"`
	KubenetesCIS      bool   `description:"(required: false, 开启K8s基线检测)"`
	KubenetesScan     bool   `description:"(required: false, 开启K8s漏洞扫描)"`
	DockerScan        bool   `description:"(required: false, 开启Docker漏洞扫描)"`
	VirusScan         bool   `description:"(required: false, 开启病毒)"`
	ImageVulnScan     bool   `description:"(required: false, 开启主仓库像漏洞扫描)"`
	HostImageVulnScan bool   `description:"(required: false, 开启主机镜像漏洞扫描)"`
	FileImageVulnScan bool   `description:"(required: false, 开启文件镜像漏洞扫描)"`
	RepoImageVulnScan bool   `description:"(required: false, 开启代码仓库镜像漏洞扫描)"`
	HostIds           string `description:"(required: false, 主机id列表 ID1,ID2,ID3)"`
	ContainerIds      string `description:"(required: false, 容器id列表 ID1,ID2,ID3)"`
	ImageIds          string `description:"(required: false, 镜像id列表 ID1,ID2,ID3)"`
	ClusterIds        string `description:"(required: false, 集群id列表 ID1,ID2,ID3)"`
	PathList          string `description:"(required: false, 目标目录列表 Path1,Path2,Path3)"`
	Type              string `description:"(required: false, 类型 registry、 host、container、image)"`
	JobId             string `description:"(required: false, JobId)"`
	TemplateId        string `description:"(required: false, 自定义模板Id)"`
}
