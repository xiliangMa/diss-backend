package statistics

import (
	"github.com/xiliangMa/diss-backend/models"
)

type PackageStatisticsService struct {
	HostId string
}

func (this *PackageStatisticsService) GetHostPackageStatistics() models.Result {
	hostPackage := new(models.HostPackage)
	hostPackage.HostId = this.HostId
	return hostPackage.GetPackageCountByType()
}

func (this *PackageStatisticsService) GetDBImageStatistics() models.Result {
	image := new(models.ImageConfig)
	image.HostId = this.HostId
	return image.GetDBCountByType()
}
