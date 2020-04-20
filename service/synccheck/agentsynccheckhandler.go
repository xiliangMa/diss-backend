package synccheck

import (
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/models/ws"
)

type AgentCheckHadler struct {
	ContinerConfig *models.ContainerConfig
	ContainerInfo  *models.ContainerInfo
}

func (this *AgentCheckHadler) Check(model string) {
	switch model {
	case ws.Resource_ContainerConfig:
		this.ContainerConfingCheck()
	case ws.Resource_ContainerInfo:
		this.ContainerInfoCheck()
	}

}

func (this *AgentCheckHadler) ContainerConfingCheck() {
	this.ContinerConfig.EmptyDirtyDataForAgent()
}

func (this *AgentCheckHadler) ContainerInfoCheck() {
	this.ContainerInfo.EmptyDirtyDataForAgent()
}
