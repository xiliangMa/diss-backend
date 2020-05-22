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

	if result.Code == http.StatusOK {
		err := this.SaveToFile(key, fpath)
		if err != nil {
			result.Code = utils.ImportLicenseFileErr
			result.Message = "ImportLicenseFileErr"
			logs.Error("Import license file fail, err: %s", err.Error())
		} else {
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
	} else {
		// 强制更新
		if result.Code == utils.CheckLicenseFileIsExistErr && isForce {
			// 更新返回值
			result.Code = http.StatusOK
			result.Message = ""

			// 更新上传文件path
			fpath = fpath + h.Filename

			// 删除旧文件
			os.Remove(fpath)

			err := this.SaveToFile(key, fpath)
			if err != nil {
				result.Code = utils.SaveLicenseFileErr
				result.Message = err.Error()
			} else {
				//更新数据库
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
// @Param LicenseUuid query string true "授权文件uuid"
// @Success 200 {object} models.Result
// @router /system/license [get]
func (this *IntegrationController) GetLicenseData() {
	licConfig := new(models.LicenseConfig)
	configId := this.GetString(":LicenseUuid")
	licConfig.LicenseUuid = configId
	json.Unmarshal(this.Ctx.Input.RequestBody, &licConfig)

	var ResultData models.Result
	ResultData.Code = http.StatusOK
	ResultData.Data = licConfig.Get()

	this.Data["json"] = ResultData
	this.ServeJSON(false)
}
