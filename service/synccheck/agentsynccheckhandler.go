package synccheck

import (
	"github.com/xiliangMa/diss-backend/models"
)

type AgentCheckHadler struct {
	ContinerConfig *models.ContainerConfig
	ContainerInfo  *models.ContainerInfo
}

func (this *AgentCheckHadler) Check(model string) {
	switch model {
	case models.Tag_ContainerConfig:
		this.ContainerConfingCheck()
	case models.Tag_ContainerInfo:
		this.ContainerInfoCheck()
	}

}

func (this *AgentCheckHadler) ContainerConfingCheck() {
	this.ContinerConfig.EmptyDirtyDataForAgent()
}

func (this *AgentCheckHadler) ContainerInfoCheck() {
	this.ContainerInfo.EmptyDirtyDataForAgent()
}
