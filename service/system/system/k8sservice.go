package system

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	uuid "github.com/satori/go.uuid"
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
	//fName := fmt.Sprintf("%x", hashName) + ext
	////this.Ctx.WriteString(  fName )
	fpath = fpath + h.Filename
	result.Code = http.StatusOK
	return result, fpath
}

func TestK8sFile(fpath string) int {
	if clientgo := utils.CreateK8sClient(fpath); clientgo.ErrMessage != "" {
		logs.Error("K8s file test not connect, err: %s", clientgo.ErrMessage)
		return utils.CheckK8sFileTestErr
	}
	return http.StatusOK
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

func AddCluster(clusterName, path, clusterType string) {
	var cluster models.Cluster
	uid, _ := uuid.NewV4()
	cluster.Id = uid.String()
	cluster.Name = clusterName
	cluster.Type = clusterType
	cluster.Status = models.Cluster_Status_Active
	cluster.SyncStatus = models.Cluster_Sync_Status_NotSynced
	cluster.FileName = path
	cluster.IsSync = models.Cluster_IsSync
	cluster.Add()
}
