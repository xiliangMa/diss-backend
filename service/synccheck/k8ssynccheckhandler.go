package synccheck

import (
	"github.com/xiliangMa/diss-backend/models"
)

type K8SCheckHadler struct {
	ContinerConfig *models.ContainerConfig
	ContainerInfo  *models.ContainerInfo
	NS             *models.NameSpace
	Pod            *models.Pod
}

func (this *K8SCheckHadler) Check(model string) {
	switch model {
	case models.Resource_ContainerConfig:
		this.ContainerConfingCheck()
	case models.Resource_ContainerInfo:
		this.ContainerInfoCheck()
	case models.Resource_NameSpace:
		this.NSCheck()
	case models.Resource_Pod:
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
