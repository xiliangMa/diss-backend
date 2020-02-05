package system

import (
	"github.com/astaxie/beego"
	css "github.com/xiliangMa/diss-backend/service/system/system"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
)

type K8sController struct {
	beego.Controller
}

// @Title UpLoadK8sFile
// @Description UpLoad K8s File
// @Param token header string true "Auth token"
// @Param k8sfile formData file true "k8s file"
// @Success 200 {object} models.Result
// @router /system/k8s/upload [post]
func (this *K8sController) UploadK8sFile() {
	f, h, _ := this.GetFile("k8sfile")
	result, fpath := css.UploadK8sFile(f, h)
	defer f.Close() //关闭上传的文件，不然的话会出现临时文件不能清除的情况
	if result.Code == http.StatusOK {
		err := this.SaveToFile("k8sfile", fpath)
		if err != nil {
			result.Code = utils.UploadK8sFileErr
			result.Message = err.Error()
		}
	}
	this.Data["json"] = result
	this.ServeJSON(false)
}
