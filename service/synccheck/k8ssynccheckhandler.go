package synccheck

import (
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/models/k8s"
	"github.com/xiliangMa/diss-backend/models/ws"
)

type K8SCheckHadler struct {
	ContinerConfig *models.ContainerConfig
	ContainerInfo  *models.ContainerInfo
	NS             *k8s.NameSpace
	Pod            *k8s.Pod
}

func (this *K8SCheckHadler) Check(model string) {
	switch model {
	case ws.Resource_ContainerConfig:
		this.ContainerConfingCheck()
	case ws.Resource_ContainerInfo:
		this.ContainerInfoCheck()
	case ws.Resource_NameSpace:
		this.NSCheck()
	case ws.Resource_Pod:
		this.PodCheck()
	}

}

func (this *K8SCheckHadler) ContainerConfingCheck() {
	this.ContinerConfig.EmptyDirtyDataForK8s()
}

func (this *K8SCheckHadler) ContainerInfoCheck() {
	this.ContainerInfo.EmptyDirtyDataForK8s()
}

func (this *K8SCheckHadler) NSCheck() {
	this.NS.EmptyDirtyData()
}

func (this *K8SCheckHadler) PodCheck() {
	this.Pod.EmptyDirtyData()
}
