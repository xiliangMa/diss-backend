package models

type SecurityCheckList struct {
	CheckList []*SecurityCheck
}

type SecurityCheck struct {
	DockerBenchMarkCheck bool             `orm:"default(false)" description:"(开启Docker基线检测)"`
	KubeBenchMarkCheck   bool             `orm:"default(false)" description:"(开启K8s基线检测)"`
	VirusScan            bool             `orm:"default(false)" description:"(开启病毒)"`
	LeakScan             bool             `orm:"default(false)" description:"(开启漏洞)"`
	Host                 *HostConfig      `orm:"null" description:"(主机)"`
	Container            *ContainerConfig `orm:"null" description:"(容器)"`
	Type                 string           `orm:"default(host)" description:"(类型 host  container)"`
	Job                  *Job             `orm:"null" description:"(来源任务)"`
}
