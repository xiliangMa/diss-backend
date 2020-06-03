package system

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
	"mime/multipart"
	"net/http"
	"os"
	"path"
)

func CheckK8sFilePost(ext, fName string) int {
	var AllowExtMap map[string]bool = map[string]bool{
		".yml":  true,
		".yaml": true,
	}
	if _, ok := AllowExtMap[ext]; !ok {
		return utils.CheckK8sFilePostErr
	}
	return http.StatusOK
}

func CheckK8sFileIsExist(fpath, fName string) int {
	if _, err := os.Stat(fpath + fName); err == nil {
		return utils.CheckK8sFileIsExistErr
	}
	return http.StatusOK
}

func Check(h *multipart.FileHeader) (models.Result, string) {
	var fpath = getK8sFilePath()
	var result models.Result
	fName := h.Filename
	ext := path.Ext(fName)

	//创建目录
	CreateKubeConfigDir(fpath)

	// 后缀名不符合上传要求
	if code := CheckK8sFilePost(ext, fName); code != http.StatusOK {
		result.Code = code
		result.Message = "CheckK8sFilePostErr"
		return result, fpath
	}

	// 检查文件是否存在
	if code := CheckK8sFileIsExist(fpath, fName); code != http.StatusOK {
		result.Code = code
		result.Message = "CheckK8sFileIsExistErr"
		return result, fpath
	}

	fpath = fpath + h.Filename
	result.Code = http.StatusOK
	return result, fpath
}

func TestClient(params utils.ApiParams) models.Result {
	code := http.StatusOK
	message := ""
	if clientgo := utils.CreateK8sClient(&params); clientgo.ErrMessage != "" {
		switch params.AuthType {
		case models.Api_Auth_Type_BearerToken:
			code = utils.CreateClientByBearerTokenErr
			message = "CreateClientByBearerTokenErr"
		case models.Api_Auth_Type_KubeConfig:
			code = utils.CreateClientByKubeConfigErr
			message = "CreateClientByKubeConfigErr"
		}
		logs.Error("K8s client connect fail, AuthType: %s, err: %s", params.AuthType, clientgo.ErrMessage)
	}
	return models.Result{Code: code, Message: message}
}

func getK8sFilePath() string {
	return beego.AppConfig.String("k8s::KubeCongigPath")
}

func CreateKubeConfigDir(fpath string) {
	_, err := os.Stat(fpath)
	if os.IsNotExist(err) {
		logs.Info("Create KubeConfig Dir success, path: %s", fpath)
		os.MkdirAll(beego.AppConfig.String("k8s::KubeCongigPath"), os.ModePerm)
	}
}
