package system

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/service/base"
	"github.com/xiliangMa/diss-backend/utils"
	"mime/multipart"
	"net/http"
	"os"
)

type LogoService struct {
}

func (this *LogoService) CreateLogoDir(fpath string) {
	_, err := os.Stat(fpath)
	if os.IsNotExist(err) {
		logs.Info("Create logo dir success, path: %s", fpath)
		os.MkdirAll(beego.AppConfig.String("system::LogoPath"), os.ModePerm)
	}
}

func (this *LogoService) Check(h *multipart.FileHeader) (models.Result, string) {
	var fpath = utils.GetLogoPath()
	var result models.Result

	//创建目录
	this.CreateLogoDir(fpath)

	// 后缀名不符合上传要求
	fileService := base.FileService{}
	fileType := models.PictureType
	if code := fileService.CheckFilePost(h, fileType); code != http.StatusOK {
		result.Code = code
		result.Message = "Png Format Incorrect."
		return result, fpath
	}

	fpath = fpath + beego.AppConfig.String("system::NewLogoName")
	result.Code = http.StatusOK
	return result, fpath
}

func (this *LogoService) CheckLogoIsExist() models.Result {
	newLogoPath := utils.GetLogoPath() + beego.AppConfig.String("system::NewLogoName")
	var result models.Result
	if _, err := os.Stat(newLogoPath); err != nil {
		result.Code = utils.CheckLogoIsNotExistErr
		result.Message = "CheckLogoIsNotExistErr"
		return result
	}
	data := make(map[string]string)
	data["url"] = "http://ip:port/" + beego.AppConfig.String("system::LogoUrl")
	result.Data = data
	result.Code = http.StatusOK
	return result
}
