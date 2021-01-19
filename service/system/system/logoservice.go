package system

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
	"mime/multipart"
	"net/http"
	"os"
	"path"
)

type LogoService struct {
}

func (this *LogoService) getLogoPath() string {
	path, _ := web.AppConfig.String("system::LogoPath")
	return path
}

func (this *LogoService) CheckLogoPost(ext, fName string) int {
	var AllowExtMap map[string]bool = map[string]bool{
		".png": true,
	}
	if _, ok := AllowExtMap[ext]; !ok {
		return utils.ChecLogoPostErr
	}
	return http.StatusOK
}

func (this *LogoService) CreateLogoDir(fpath string) {
	_, err := os.Stat(fpath)
	if os.IsNotExist(err) {
		logs.Info("Create logo Dir success, path: %s", fpath)
		path, _ := web.AppConfig.String("system::LogoPath")
		os.MkdirAll(path, os.ModePerm)
	}
}

func (this *LogoService) Check(h *multipart.FileHeader) (models.Result, string) {
	var fpath = this.getLogoPath()
	var result models.Result
	fName := h.Filename
	ext := path.Ext(fName)

	//创建目录
	this.CreateLogoDir(fpath)

	// 后缀名不符合上传要求
	if code := this.CheckLogoPost(ext, fName); code != http.StatusOK {
		result.Code = code
		result.Message = "ChecLogoPostErr"
		return result, fpath
	}

	logoName, _ := web.AppConfig.String("system::NewLogoName")
	fpath = fpath + logoName
	result.Code = http.StatusOK
	return result, fpath
}

func (this *LogoService) CheckLogoIsExist() models.Result {
	logoName, _ := web.AppConfig.String("system::NewLogoName")
	newLogoPath := this.getLogoPath() + logoName
	var result models.Result
	if _, err := os.Stat(newLogoPath); err != nil {
		result.Code = utils.CheckLogoIsNotExistErr
		result.Message = "CheckLogoIsNotExistErr"
		return result
	}
	data := make(map[string]string)
	url, _ := web.AppConfig.String("system::LogoUrl")
	data["url"] = "http://ip:port/" + url
	result.Data = data
	result.Code = http.StatusOK
	return result
}
