package sysinit

import "github.com/xiliangMa/diss-backend/service/system/system"

func InitTrialLicense() {
	licenseService := system.LicenseService{}
	licenseService.InitTrialLicense()
}
