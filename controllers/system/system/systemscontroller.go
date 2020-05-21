package system

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	css "github.com/xiliangMa/diss-backend/service/system/system"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
)

type SystemController struct {
	beego.Controller
}

// @Title UpLoadLogo
// @Description UpLoad Logo
// @Param token header string true "authToken"
// @Param logo formData file true "logo"
// @Success 200 {object} models.Result
// @router /system/logo [post]
func (this *SystemController) UploadLogo() {
	key := "logo"
	f, h, _ := this.GetFile(key)
	defer f.Close()
	logoService := new(css.LogoService)
	result, fpath := logoService.Check(h)
	if result.Code != http.StatusOK {
		logs.Error("Upload logo  fail, err: %s", result.Message)
	} else {
		err := this.SaveToFile(key, fpath)
		if err != nil {
			result.Code = utils.UploadLogoErr
			result.Message = "UploadLogoErr"
			logs.Error("Upload logo  fail, err: %s", err.Error())
		} else {
			result.Code = http.StatusOK
		}
	}
	this.Data["json"] = result
	this.ServeJSON(false)
}

// @Title Check Logo isexist
// @Description Check kLogo IsExist
// @Param token header string true "authToken"
// @Success 200 {object} models.Result
// @router /system/logo/isexist [get]
func (this *SystemController) CheckLogoIsExist() {
	logoService := new(css.LogoService)
	result := logoService.CheckLogoIsExist()
	this.Data["json"] = result
	this.ServeJSON(false)
}
