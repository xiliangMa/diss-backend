package system

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	css "github.com/xiliangMa/diss-backend/service/system/system"
	"github.com/xiliangMa/diss-backend/utils"
	"io/ioutil"
	"net/http"
	"os"
)

// @Title Import License File
// @Description ImportLicense
// @Param token header string true "authToken"
// @Param licenseFile formData file true "licenseFile"
// @Param isForce formData bool false true "force"
// @Success 200 {object} models.Result
// @router /system/license/import [post]
func (this *IntegrationController) AddLicenseFile() {
	// to do 生成功能需要移动到许可证管理系统
	key := "licenseFile"
	isForce, _ := this.GetBool("isForce", true)
	f, h, _ := this.GetFile(key)
	result, fpath := css.CheckLicenseFile(h)
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
			result = css.License_RSA_Decrypt(licenseByte, true)
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
					result = css.License_RSA_Decrypt(licenseByte, false)
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
func (this *IntegrationController) GetLicense() {
	licConfig := new(models.LicenseConfig)
	configId := this.GetString("id")
	licConfig.Id = configId
	json.Unmarshal(this.Ctx.Input.RequestBody, &licConfig)

	this.Data["json"] = licConfig.Get()
	this.ServeJSON(false)
}

// @Title Get License History
// @Description Get LicenseHistory
// @Param token header string true "authToken"
// @Param from query int 0 false "from"
// @Param limit query int 50 false "limit"
// @Success 200 {object} models.Result
// @router /system/licensehistory [get]
func (this *IntegrationController) GetLicenseHistory() {
	licHistory := new(models.LicenseHistory)
	limit, _ := this.GetInt("limit")
	from, _ := this.GetInt("from")

	this.Data["json"] = licHistory.List(from, limit)
	this.ServeJSON(false)
}
