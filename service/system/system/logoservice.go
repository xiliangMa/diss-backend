package system

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
	"mime/multipart"
	"net/http"
	"os"
)

type LogoService struct {
}

func (this *LogoService) getLogoPath() string {
	return beego.AppConfig.String("system::LogoPath")
}

func (this *LogoService) CheckLogoPost(fh *multipart.FileHeader) int {

	// Open File
	f, err := fh.Open()
	if err != nil {
		logs.Error("Open file error, err: %s", err)
		return utils.ChecLogoPostErr
	}
	defer f.Close()

	// Get the content
	datatype, err := this.GetFileContentType(f)
	if err != nil {
		logs.Error("Get the content error, err: %s", err)
		return utils.ChecLogoPostErr
	}

	if datatype != models.PictureType {
		return utils.ChecLogoPostErr
	}

	return http.StatusOK
}

func (this *LogoService) CreateLogoDir(fpath string) {
	_, err := os.Stat(fpath)
	if os.IsNotExist(err) {
		logs.Info("Create logo dir success, path: %s", fpath)
		os.MkdirAll(beego.AppConfig.String("system::LogoPath"), os.ModePerm)
	}
}

func (this *LogoService) Check(h *multipart.FileHeader) (models.Result, string) {
	var fpath = this.getLogoPath()
	var result models.Result

	//创建目录
	this.CreateLogoDir(fpath)

	// 后缀名不符合上传要求
	if code := this.CheckLogoPost(h); code != http.StatusOK {
		result.Code = code
		result.Message = "Png Format Incorrect."
		return result, fpath
	}

	fpath = fpath + beego.AppConfig.String("system::NewLogoName")
	result.Code = http.StatusOK
	return result, fpath
}

func (this *LogoService) CheckLogoIsExist() models.Result {
	newLogoPath := this.getLogoPath() + beego.AppConfig.String("system::NewLogoName")
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

func (this *LogoService) GetFileContentType(file multipart.File) (string, error) {

	buffer := make([]byte, 512)

	contentType := ""
	_, err := file.Read(buffer)
	if err != nil {
		return contentType, err
	}

	contentType = http.DetectContentType(buffer)

	return contentType, nil
}
