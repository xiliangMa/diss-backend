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

func TestClient(params models.ApiParams) models.Result {
	ResultData := models.Result{}
	ResultData.Code = http.StatusOK
	// 检测连接是否可用
	clientgo := models.CreateK8sClient(&params)
	if clientgo.ErrMessage != "" {
		switch params.AuthType {
		case models.Api_Auth_Type_BearerToken:
			ResultData.Code = utils.CreateClientByBearerTokenErr
			ResultData.Message = "CreateClientByBearerTokenErr"
		case models.Api_Auth_Type_KubeConfig:
			ResultData.Code = utils.CreateClientByKubeConfigErr
			ResultData.Message = "CreateClientByKubeConfigErr"
		}
		logs.Error("K8s client connect fail, AuthType: %s, err: %s", params.AuthType, clientgo.ErrMessage)
		return ResultData
	} else {
		// 检测集群是否有主机
		nodes, err := clientgo.GetNodes()
		if err != nil || len(nodes.Items) == 0 {
			ResultData.Code = utils.ClusterNotAvailableOrNoHostErr
			ResultData.Message = "ClusterNotAvailableOrNoHostErr"
		}
	}
	return ResultData

}

func getK8sFilePath() string {
	return beego.AppConfig.String("k8s::KubeConfigPath")
}

func CreateKubeConfigDir(fpath string) {
	_, err := os.Stat(fpath)
	if os.IsNotExist(err) {
		logs.Info("Create KubeConfig Dir success, path: %s", fpath)
		os.MkdirAll(beego.AppConfig.String("k8s::KubeConfigPath"), os.ModePerm)
	}
}
