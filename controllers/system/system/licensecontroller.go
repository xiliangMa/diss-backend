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
	"strings"
	"time"
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
	licenseService.IsForce = isForce
	result := licenseService.CheckLicenseFile(h)
	licenseByte := []byte{}

	defer f.Close()
	// 文件上传或更新不成功，直接返回（预处理未通过）
	if result.Code != http.StatusOK {
		this.Data["json"] = result
		this.ServeJSON(false)
		return
	}

	fpath := licenseService.LicTrialFile
	// 正式授权先保存为临时名称
	if licenseService.LicType == models.LicType_StandardLicense {
		fpath = licenseService.LicStandTmpFile
	}
	err := this.SaveToFile(key, fpath)

	if err != nil {
		result.Code = utils.SaveLicenseFileErr
		result.Message = "Save LicenseFile Err"
		logs.Error("Save LicenseFile fail, err: %s", err.Error())
		this.Data["json"] = result
		this.ServeJSON(false)
		return
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

	if result.Code != http.StatusOK {
		// 导入授权出错，删除上传的文件
		err := os.Remove(fpath)
		if err != nil {
			result.Code = utils.DeleteFileErr
			result.Message = "Delete file Error: " + fpath
		}
	}
	if strings.Contains(fpath, models.LicType_StandardLicense) {
		os.Rename(licenseService.LicStandTmpFile, licenseService.LicStandFile)
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
	result := licenseService.GetLicenseData(licConfig)
	if result.Data != nil {
		licData := result.Data.([]*models.LicenseConfig)

		if len(licData) > 0 {
			if licData[0].Type != models.LicType_TrialLicense {
				fcode := licData[0].Id
				featurecodeData := licenseService.VerifyFeatureCode(fcode)
				if featurecodeData.Data == false {
					result.Code = utils.LicenseFCodeErr
					result.Message = "Feature Code Error"
					result.Data = nil
				}
			} else {
				fcode := licenseService.GetFeatureCode()
				licData[0].Id = fcode
			}
		}
	}

	this.Data["json"] = result
	this.ServeJSON(false)
}

// @Title Get Licensed Host Count
// @Description Get Licensed Host Count
// @Param token header string true "authToken"
// @Success 200 {object} models.Result
// @router /system/license/hostcount [get]
func (this *LicenseController) GetLicensedHostCount() {
	result := models.Result{Code: http.StatusOK}
	hostObj := models.HostConfig{}
	hostObj.LicCount = true
	licenseHostCount := hostObj.Count()
	result.Data = licenseHostCount

	this.Data["json"] = result
	this.ServeJSON(false)
}

// @Title Set Host License
// @Description Set One Host License
// @Param token header string true "authToken"
// @Param hostId path string "" true "hostId"
// @Param isLicense query string "" true "add or remove host license, bool"
// @Success 200 {object} models.Result
// @router /system/license/hostlicense/:hostId [put]
func (this *LicenseController) SetHostLicense() {
	result := models.Result{Code: http.StatusOK}
	licenseService := sssystem.LicenseService{}

	// 1 - 获取已授权的主机节点个数
	licenseHostCount := licenseService.GetLicensedHostCount()

	// 2 - 获取授权管理数据（源自授权文件的）
	licConfig := new(models.LicenseConfig)
	result = licenseService.GetLicenseData(licConfig)

	if result.Code != 0 && result.Code != http.StatusOK {
		this.Data["json"] = result
		this.ServeJSON(false)
		return
	}
	licData := result.Data.([]*models.LicenseConfig)
	if len(licData) < 1 {
		result.Code = utils.GetLicenseDataErr
		result.Message = "get license data error"
		this.Data["json"] = result
		this.ServeJSON(false)
		return
	}

	// 3 - 主机授权个数校验， 现在校验标准是 小于 基线模块的授权个数
	licenseStatus, _ := this.GetBool("isLicense")
	licModule := models.LicenseModule{}
	licModule.ModuleCode = models.LicModuleType_BenchMark
	licBenchData := licModule.Get()
	licenseBenchCount := licBenchData.LicenseCount
	if licenseStatus && licenseHostCount >= licenseBenchCount {

		result.Code = utils.LicenseHostCountErr
		result.Message = "reach license count , cant add host"
		this.Data["json"] = result
		this.ServeJSON(false)
		return
	}

	// 4 - ，如果是添加授权，检查是否授权过期
	expireTime := licBenchData.LicenseExpireAt
	curTime := time.Now().UnixNano()
	if curTime > expireTime {
		result.Code = utils.LicenseExpiredErr
		result.Message = "License Expired"
		this.Data["json"] = result
		this.ServeJSON(false)
		return
	}

	// 5 - 修改主机授权状态
	hostId := this.GetString(":hostId")
	hostConfig := new(models.HostConfig)
	hostConfig.Id = hostId
	hostConfig = hostConfig.Get()
	if hostConfig == nil {
		result.Code = utils.HostExistError
		result.Message = "Not Fonud host"
		this.Data["json"] = result
		this.ServeJSON(false)
		return
	}
	hostConfig.IsLicensed = licenseStatus
	result = hostConfig.Update()
	if result.Code != http.StatusOK {
		result.Code = utils.LicenseHostErr
		result.Message = "LicenseHostErr"
		this.Data["json"] = result
		this.ServeJSON(false)
		return
	}

	// 再次通过数据库获取最新的授权个数
	licenseHostCount = licenseService.GetLicensedHostCount()

	result.Data = licenseHostCount
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

// @Title Get FeatureCode
// @Description GetFeatureCode
// @Param token header string true "authToken"
// @Success 200 {object} models.Result
// @router /system/featurecode [get]
func (this *LicenseController) GetFeatureCode() {
	result := models.Result{Code: http.StatusOK}

	licenseService := sssystem.LicenseService{}
	fcode := licenseService.GetFeatureCode()
	result.Data = fcode

	this.Data["json"] = result
	this.ServeJSON(false)
}

// @Title Verify Feature Code
// @Description Verify FeatureCode
// @Param token header string true "authToken"
// @Param featurecode query string true "featurecode"
// @Success 200 {object} models.Result
// @router /system/featurecode/verify [get]
func (this *LicenseController) VerifyFeatureCode() {
	fcode := this.GetString("featurecode")
	result := models.Result{}

	if fcode != "" {
		licenseService := sssystem.LicenseService{}
		result = licenseService.VerifyFeatureCode(fcode)
	}

	this.Data["json"] = result
	this.ServeJSON(false)
}

// @Title Get module License Count(获取授权模块授权数)
// @Description GetModuleLicenseCount
// @Param token header string true "authToken"
// @Param body body models.LicenseModule false "license 模块信息"
// @Success 200 {object} models.Result
// @router /system/license/modulecount [post]
func (this *LicenseController) GetModuleLicenseCount() {
	lm := new(models.LicenseModule)
	json.Unmarshal(this.Ctx.Input.RequestBody, &lm)
	this.Data["json"] = lm.List()
	this.ServeJSON(false)
}
