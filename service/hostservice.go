package service

import (
	"github.com/xiliangMa/diss-backend/models"
	"net/http"
)

type HostService struct {
	Host *models.HostConfig
}

// todo add result
func (this *HostService) Delete() models.Result {
	result := models.Result{Code: http.StatusOK}
	hostId := this.Host.Id

	// delete ImageConfig
	ic := models.ImageConfig{}
	ic.HostId = hostId
	result = ic.Delete()

	// delete ImageInfo
	ii := models.ImageInfo{}
	ii.HostId = hostId
	result = ii.Delete()

	// delete ContainerConfig
	cc := models.ContainerConfig{}
	cc.HostId = hostId
	result = cc.Delete()

	// delete ContainerInfo
	ci := models.ContainerInfo{}
	ci.HostId = hostId
	result = ci.Delete()

	// delete HostInfo
	hi := models.HostInfo{}
	hi.Id = hostId
	result = hi.Delete()

	// delete HostConfig
	result = this.Host.Delete()
	return result
}
