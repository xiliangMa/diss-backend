package system

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	sssystem "github.com/xiliangMa/diss-backend/service/system/system"
	"github.com/xiliangMa/diss-backend/utils"
	"io/ioutil"
	"net/http"
	"os"
)

type LicenseController struct {
	beego.Controller
}

// @Title Import License File
// @Description ImportLicense
// @Param token header string true "authToken"
// @Param licenseFile formData file true "licenseFile"
// @Param isForce formData bool false true "force"
// @Success 200 {object} models.Result
// @router /system/license/import [post]
func (this *LicenseController) AddLicenseFile() {
	key := "licenseFile"
	isForce, _ := this.GetBool("isForce", true)
	f, h, _ := this.GetFile(key)

	licenseService := sssystem.LicenseService{}
	result, fpath := licenseService.CheckLicenseFile(h)
	licenseByte := []byte{}

	defer f.Close()

	err := this.SaveToFile(key, fpath)

	if err != nil && result.Code != utils.CheckLicenseFileIsExistErr {
		result.Code = utils.ImportLicenseFileErr
		result.Message = "ImportLicenseFileErr"
		logs.Error("Import license file fail, err: %s", err.Error())
	}

	if isForce { // 强制更新
		// 更新返回值
		result.Code = http.StatusOK
		result.Message = ""

		//更新数据库
		licenseByte, err = ioutil.ReadAll(f)
		if err != nil {
			result.Code = utils.ImportLicenseFileErr
			result.Message = err.Error()
			logs.Error("Read license file fail: %s", err)
		} else {
			licenseService.LicenseByte = licenseByte
			licenseService.IsUpdate = true
			result = licenseService.LicenseActive()
		}
	} else {
		if result.Code == http.StatusOK {
			// 添加license 到数据库
			f, err := os.Open(fpath)
			if err != nil {
				result.Code = utils.ImportLicenseFileErr
				result.Message = err.Error()
				logs.Error("Open license file fail: %s", err)
			} else {
				licenseByte, err = ioutil.ReadAll(f)
				if err != nil {
					result.Code = utils.ImportLicenseFileErr
					result.Message = err.Error()
					logs.Error("Read license file fail: %s", err)
				} else {
					licenseService.LicenseByte = licenseByte
					licenseService.IsUpdate = false
					result = licenseService.LicenseActive()
				}
			}
		}
	}

	this.Data["json"] = result
	this.ServeJSON(false)
}

// @Title Get License
// @Description GetLicense
// @Param token header string true "authToken"
// @Param id query string false "id"
// @Success 200 {object} models.Result
// @router /system/license [get]
func (this *LicenseController) GetLicense() {
	licConfig := new(models.LicenseConfig)
	configId := this.GetString("id")
	licConfig.Id = configId
	json.Unmarshal(this.Ctx.Input.RequestBody, &licConfig)

	licenseService := sssystem.LicenseService{}
	result := licenseService.CheckLicenseType()

	if result.Code != utils.NoLicenseFileErr {
		licConfig.Type = result.Data.(string)
		licData := licConfig.Get()
		result.Data = licData.Data
	}

	this.Data["json"] = result
	this.ServeJSON(false)
}

// @Title Get License History
// @Description Get LicenseHistory
// @Param token header string true "authToken"
// @Param from query int 0 false "from"
// @Param limit query int 50 false "limit"
// @Success 200 {object} models.Result
// @router /system/licensehistory [get]
func (this *LicenseController) GetLicenseHistory() {
	licHistory := new(models.LicenseHistory)
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	this.Data["json"] = licHistory.List(from, limit)
	this.ServeJSON(false)
}
