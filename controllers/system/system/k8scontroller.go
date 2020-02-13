package system

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	css "github.com/xiliangMa/diss-backend/service/system/system"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"os"
)

type K8sController struct {
	beego.Controller
}

// @Title UpLoadK8sFile
// @Description UpLoad K8s File
// @Param token header string true "auth token"
// @Param k8sFile formData file true "k8s file"
// @Param clusterName formData string true "cluster name"
// @Param isForce formData bool true true "force update file"
// @Success 200 {object} models.Result
// @router /system/k8s/upload [post]
func (this *K8sController) UploadK8sFile() {
	key := "k8sFile"
	isForce, _ := this.GetBool("isForce", true)
	clusterName := this.GetString("clusterName")
	f, h, _ := this.GetFile(key)
	result, fpath := css.Check(f, h)
	defer f.Close()
	if result.Code == http.StatusOK {
		err := this.SaveToFile(key, fpath)
		if err != nil {
			result.Code = utils.UploadK8sFileErr
			result.Message = "UploadK8sFileErr"
			logs.Error("Upload k8s file  fail, err: %s", err.Error())
		} else {
			//检查连接是否可用
			if code := css.TestK8sFile(fpath); code != http.StatusOK {
				result.Code = code
				result.Message = "CheckK8sFileTestErr"
				os.Remove(fpath)
			} else {
				logs.Info("Upload k8s file success, file name: %s", h.Filename)
				// 添加集群记录
				css.AddCluster(clusterName, fpath)
			}
		}
	} else {
		// 强制更新
		if result.Code == utils.CheckK8sFileIsExistErr && isForce {
			// 更新返回值
			result.Code = http.StatusOK
			result.Message = ""

			// 更新上传文件path
			fpath = fpath + h.Filename

			// 删除旧文件
			os.Remove(fpath)

			err := this.SaveToFile(key, fpath)
			if err != nil {
				result.Code = utils.UploadK8sFileErr
				result.Message = err.Error()
			} else {
				//检查连接是否可用
				if code := css.TestK8sFile(fpath); code != http.StatusOK {
					result.Code = code
					result.Message = "CheckK8sFileTestErr"
					os.Remove(fpath)
				} else {
					logs.Info("Force update k8s file success, file name: %s", h.Filename)
					// to do 强制更新后文件名相同、内容不一样
				}
			}
		}
	}
	this.Data["json"] = result
	this.ServeJSON(false)
}
