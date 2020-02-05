package system

import (
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
	"mime/multipart"
	"net/http"
	"os"
	"path"
)

func CheckK8sFile(ext, fileName string) int {
	// 后缀名不符合上传要求
	var AllowExtMap map[string]bool = map[string]bool{
		".yml":  true,
		".yaml": true,
	}
	if _, ok := AllowExtMap[ext]; !ok {
		return utils.CheckK8sFilePostErr
	}

	// 检查文件是否存在
	if _, err := os.Stat("kubeconfig/" + fileName); err == nil {
		return utils.CheckK8sFileExistErr
	}
	return http.StatusOK
}

func UploadK8sFile(f multipart.File, h *multipart.FileHeader) (models.Result, string) {
	fpath := "kubeconfig/"
	var result models.Result
	fileName := h.Filename
	ext := path.Ext(fileName)
	if code := CheckK8sFile(ext, fileName); code != http.StatusOK {
		result.Code = code
		result.Message = "CheckK8sFilePostErr"
		return result, fpath
	}
	//创建目录
	//uploadDir := "kubeconfig/" + time.Now().Format("2006/01/02/")
	//err := os.MkdirAll(uploadDir, 777)
	//if err != nil {
	//	//this.Ctx.WriteString(fmt.Sprintf("%v", err))
	//	result.Code = 1
	//	return result, ""
	//}
	//构造文件名称
	//rand.Seed(time.Now().UnixNano())
	//randNum := fmt.Sprintf("%d", rand.Intn(9999)+1000)
	//hashName := md5.Sum([]byte(time.Now().Format("2006_01_02_15_04_05_") + randNum))
	//
	//fileName := fmt.Sprintf("%x", hashName) + ext
	////this.Ctx.WriteString(  fileName )
	fpath = fpath + h.Filename
	result.Code = http.StatusOK
	return result, fpath
}
