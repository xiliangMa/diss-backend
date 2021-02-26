package system

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/shirou/gopsutil/host"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

type LicenseService struct {
	LicenseByte []byte
	IsForce     bool
	IsUpdate    bool
	FileType    string
}

func (this *LicenseService) LicenseActive() models.Result {

	result := models.Result{}
	result.Code = http.StatusOK
	licenseObject := new(models.LicenseConfig)
	licenseService := LicenseService{}
	licenseService.FileType = models.EncryptedFileType_License
	result = licenseService.LicenseRSADecrypt(this.LicenseByte)

	plainText := result.Data.([]byte)
	err := json.Unmarshal(plainText, &licenseObject)
	if err != nil {
		result.Code = utils.LicenseUnmarshalErr
		msg := fmt.Sprintf("License file format or content fail, err:  %s", err)
		result.Message = msg
		result.Data = nil
		logs.Error(msg)
		return result
	}

	if licenseObject.Type != models.LicType_TrialLicense {
		result = this.VerifyFeatureCode(licenseObject.Id)
		if result.Data == false {
			result.Code = utils.LicenseFCodeErr
			result.Message = "Feature Code Error"
			result.Data = nil
			return result
		}
	}

	licenseObject.ActiveAt = time.Now().UnixNano()

	message := ""
	licenseQueryObj := models.LicenseConfig{}

	if licenseObject.Type == models.LicType_TrialLicense {
		licenseQueryObj.Id = ""
		licenseQueryObj.Type = models.LicType_TrialLicense
	} else {
		licenseQueryObj.Id = licenseObject.Id
		licenseQueryObj.Type = models.LicType_StandardLicense
	}
	licInDb := licenseQueryObj.Get()

	if licInDb.Data != nil {
		licData := licInDb.Data.([]*models.LicenseConfig)
		if len(licData) > 0 {
			this.IsUpdate = true
			if licenseObject.Type == models.LicType_TrialLicense {
				licenseObject.Id = licData[0].Id
			}
		}
	}
	if this.IsUpdate {
		result = licenseObject.Update()
		message = fmt.Sprintf("Force update license file success, License: %s", string(plainText))
		logs.Info(message)
	} else {
		result = licenseObject.Add()
		message = fmt.Sprintf("License import success, License：%s", string(plainText))
		logs.Info(message)
	}

	licenseService.GetLicensedHostCount()
	//添加license 历史
	licenseHistory := models.LicenseHistory{}
	licenseObjJson, _ := json.Marshal(licenseObject)
	licenseHistory.LicenseJson = string(licenseObjJson)
	licenseHistory.Add()

	return result
}

func (this *LicenseService) CheckLicenseFile(h *multipart.FileHeader) (models.Result, string) {
	licenseService := LicenseService{}
	var fpath = licenseService.GetLicenseFilePath()
	var result models.Result
	fName := h.Filename
	ext := path.Ext(fName)
	licFilename := models.LicType_StandardLicense + models.LicFile_Extension
	if strings.Contains(strings.ToLower(fName), "trial") {
		err := os.Remove(fpath + licFilename)
		if err != nil {
			result.Code = utils.DeleteFileErr
			message := "Delete file Error: " + fpath + licFilename
			result.Message = message
			logs.Info(message)
			return result, fpath
		}
		licFilename = models.LicType_TrialLicense + models.LicFile_Extension
	}

	//创建目录
	licenseService.createLicenseConfigDir(fpath)

	// 后缀名不符合上传要求
	if code := licenseService.checkLicenseFilePost(ext, fName); code != http.StatusOK {
		result.Code = code
		result.Message = "CheckLicenseFilePostErr"
		return result, fpath
	}

	// 非强制更新时，检查文件是否存在
	if !this.IsForce {
		if code := licenseService.checkLicenseFileIsExist(fpath, licFilename); code != http.StatusOK {
			result.Code = code
			result.Message = "CheckLicenseFileIsExistErr"
			return result, fpath
		}
	}

	fpath = fpath + licFilename
	result.Code = http.StatusOK
	return result, fpath
}

func (this *LicenseService) CheckLicenseType() (res models.Result) {
	licenseService := LicenseService{}
	var fpath = licenseService.GetLicenseFilePath()
	var result models.Result
	result.Code = http.StatusOK

	if licenseService.checkLicenseFileIsExist(fpath, models.LicType_StandardLicense+models.LicFile_Extension) != http.StatusOK {
		//取正式版授权
		result.Message = "Standard License"
		result.Data = models.LicType_StandardLicense
	} else if licenseService.checkLicenseFileIsExist(fpath, models.LicType_TrialLicense+models.LicFile_Extension) != http.StatusOK {
		//取测试版授权
		result.Message = "Trail License"
		result.Data = models.LicType_TrialLicense
	} else {
		//无有效的授权文件
		result.Message = "No License File Found"
		result.Code = utils.NoLicenseFileErr
	}
	// todo 到授权中心验证授权文件的有效性

	return result
}

func (this *LicenseService) InitTrialLicense() models.Result {
	result := models.Result{Code: http.StatusOK}
	licenseService := LicenseService{}
	var fpath = licenseService.GetLicenseFilePath()
	licenseFileType := licenseService.CheckLicenseType()
	if licenseFileType.Code == utils.NoLicenseFileErr {
		// 没有发现授权文件
		return licenseFileType
	}

	licenseByte, err := ioutil.ReadFile(fpath + models.LicType_TrialLicense + models.LicFile_Extension)
	if err != nil {
		result.Code = utils.ImportLicenseFileErr
		result.Message = err.Error()
		logs.Error("Read license file fail: %s", err)
		return result
	} else {
		licConfig := new(models.LicenseConfig)
		licConfig.Type = models.LicType_TrialLicense
		licData := licConfig.Get()
		if licData.Code != http.StatusOK {
			return licData
		}

		if licData.Data == nil {
			return result
		}

		if licData.Data != nil {
			licList := licData.Data.([]*models.LicenseConfig)
			if len(licList) < 1 {
				licenseService.LicenseByte = licenseByte
				result = licenseService.LicenseActive()
			}
		}

		if result.Code != http.StatusOK {
			return result
		}
	}

	logs.Info("Import trial license file ok. ")
	return result
}

func (this *LicenseService) GetLicenseFilePath() string {
	return beego.AppConfig.String("license::LicensePath")
}

func (this *LicenseService) GetLicenseData(licConfig *models.LicenseConfig) models.Result {
	result := models.Result{Code: http.StatusOK}
	result = this.CheckLicenseType()

	if result.Code != utils.NoLicenseFileErr {
		licConfig.Type = result.Data.(string)
		licData := licConfig.Get()
		result.Data = licData.Data
	}
	return result
}

func (this *LicenseService) GetLicensedHostCount() int64 {
	hostObj := models.HostConfig{}
	hostObj.LicCount = true
	licenseHostCount := hostObj.Count()
	this.UpdateLicenseCountOfModules(licenseHostCount)
	return licenseHostCount
}

func (this *LicenseService) UpdateLicenseCountOfModules(count int64) models.Result {
	result := models.Result{Code: http.StatusOK}
	lm := models.LicenseModule{ModuleCode: models.LicModuleType_BenchMark}
	dbLm := lm.Get()
	if dbLm != nil && dbLm.Id != "" {
		dbLm.IsLicensedCount = count
		result = dbLm.Update()
	}
	lm = models.LicenseModule{ModuleCode: models.TMP_Type_HostImageVulnScan}
	dbLm = lm.Get()
	if dbLm != nil && dbLm.Id != "" && dbLm.LicenseCount > 0 {
		dbLm.IsLicensedCount = count
		result = dbLm.Update()
	}
	return result
}

func (this *LicenseService) LicenseRSADecrypt(licenseByte []byte) models.Result {
	result := models.Result{}
	result.Code = http.StatusOK

	decodeBytes, err := base64.StdEncoding.DecodeString(string(licenseByte))
	if err != nil {
		result.Code = utils.LicenseBase64DecodeErr
		logs.Error("License import fail, license base64 decode err,  %s", err)
	}
	privateKey := "conf/private.pem"
	if this.FileType == models.EncryptedFileType_FeatureCode {
		privateKey = "conf/featurePrivate.pem"
	}
	plainText := utils.RSA_Decrypt(decodeBytes, privateKey)

	result.Data = plainText
	return result
}

func (this *LicenseService) VerifyFeatureCode(fcode string) models.Result {
	result := models.Result{}
	result.Code = http.StatusOK

	hostFeatureCode := this.GetFeatureCode()
	isVerified := hostFeatureCode == fcode
	result.Data = isVerified
	return result
}

func (this *LicenseService) GenerateFeatureCode() models.Result {
	result := models.Result{}
	result.Code = http.StatusOK

	hInfo, _ := host.Info()
	publicKey := "conf/featurePublic.pem"
	cipherText := utils.RSA_Encrypt([]byte(hInfo.HostID), publicKey)
	plainText := base64.StdEncoding.EncodeToString(cipherText)
	result.Data = plainText
	return result
}

func (this *LicenseService) GetFeatureCode() string {
	featureCode := ""
	sysConfig := models.SysConfig{}
	sysConfig.Key = models.FeatureCode
	featureCodeCfg := sysConfig.Get()
	if featureCodeCfg != nil {
		featureCode = featureCodeCfg.Value
	} else {
		result := this.GenerateFeatureCode()
		if result.Code == http.StatusOK {
			featureCode = result.Data.(string)
			sysConfig.Value = featureCode
			sysConfig.Add()
			logs.Info("Generated FeatureCode :", featureCode)
		}
	}

	return featureCode
}

func (this *LicenseService) createLicenseConfigDir(fpath string) {
	_, err := os.Stat(fpath)
	if os.IsNotExist(err) {
		logs.Info("Create License Dir success, path: %s", fpath)
		os.MkdirAll(beego.AppConfig.String("license::LicensePath"), os.ModePerm)
	}
}

func (this *LicenseService) checkLicenseFileIsExist(fpath, fName string) int {
	if _, err := os.Stat(fpath + fName); err == nil {
		return utils.CheckLicenseFileIsExistErr
	}
	return http.StatusOK
}

func (this *LicenseService) checkLicenseFilePost(ext, fName string) int {
	var AllowExtMap map[string]bool = map[string]bool{
		models.LicFile_Extension: true,
	}
	if _, ok := AllowExtMap[ext]; !ok {
		return utils.CheckLicenseFilePostErr
	}
	return http.StatusOK
}
