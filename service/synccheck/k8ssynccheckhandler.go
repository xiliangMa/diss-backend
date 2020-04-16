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
	case ws.Tag_ContainerConfig:
		this.ContainerConfingCheck()
	case ws.Tag_ContainerInfo:
		this.ContainerInfoCheck()
	case ws.Tag_NameSpace:
		this.NSCheck()
	case ws.Tag_Pod:
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
