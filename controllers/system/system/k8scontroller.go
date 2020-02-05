package system

import (
	"github.com/astaxie/beego"
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
// @Param token header string true "Auth token"
// @Param k8sfile formData file true "k8s file"
// @Success 200 {object} models.Result
// @router /system/k8s/upload [post]
func (this *K8sController) UploadK8sFile() {
	f, h, _ := this.GetFile("k8sfile")
	result, fpath := css.Check(f, h)
	defer f.Close()
	if result.Code == http.StatusOK {
		err := this.SaveToFile("k8sfile", fpath)
		if err != nil {
			result.Code = utils.UploadK8sFileErr
			result.Message = err.Error()
		} else {
			//检查连接是否可用
			if code := css.TestK8sFile(fpath); code != http.StatusOK {
				result.Code = code
				result.Message = "CheckK8sFileTestErr"
				os.Remove(fpath)
			}
		}
	}

	this.Data["json"] = result
	this.ServeJSON(false)
}
