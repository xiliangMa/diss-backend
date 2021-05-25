package securitylog

import (
	"github.com/xiliangMa/diss-backend/models"
	"time"
)

type SensitiveInfoService struct {
	SensitiveInfo models.SensitiveInfo
}

func (this *SensitiveInfoService) AddFileList() {

	sensiInfo := this.SensitiveInfo
	sensiInfo.CreateTime = time.Now().UnixNano()
	for _, fileinfo := range sensiInfo.Files {
		sensiInfo.FileName = fileinfo.Name
		sensiInfo.MD5 = fileinfo.MD5
		sensiInfo.Permission = fileinfo.Permission
		sensiInfo.FileType = fileinfo.Type
		sensiInfo.Size = fileinfo.Size

		sensiInfo.Add()
	}
}
