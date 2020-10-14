package system

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"time"
)

type LicenseService struct {
	LicenseByte []byte
	IsUpdate    bool
}

func (this *LicenseService) LicenseActive() models.Result {

	result := models.Result{}
	result.Code = http.StatusOK
	licenseObject := new(models.LicenseConfig)
	licenseService := LicenseService{}

	result = licenseService.licenseRSADecrypt(this.LicenseByte)
	plainText := result.Data.([]byte)
	err := json.Unmarshal(plainText, &licenseObject)
	if err != nil {
		result.Code = utils.LicenseUnmarshalErr
		logs.Error("License import fail, Unmarshal license err,  %s", err)
	}

	licenseObject.ActiveAt = time.Now()

	message := ""
	if this.IsUpdate {
		result = licenseObject.Update()
		message = fmt.Sprintf("Force update license file success, License: %s", string(plainText))
		logs.Info(message)
	} else {
		result = licenseObject.Add()
		message = fmt.Sprintf("License import success, License：", string(plainText))
		logs.Info(message)
	}
	//添加license 历史
	liceseHistory := models.LicenseHistory{}
	liceseHistory.LicenseJson = string(plainText)
	liceseHistory.Add()

	return result
}

func (this *LicenseService) CheckLicenseFile(h *multipart.FileHeader) (models.Result, string) {
	licenseService := LicenseService{}
	var fpath = licenseService.GetLicenseFilePath()
	var result models.Result
	fName := h.Filename
	ext := path.Ext(fName)
	stanardLicFilename := models.LicType_StandardLicense + models.LicFile_Extension

	//创建目录
	licenseService.createLicenseConfigDir(fpath)

	// 后缀名不符合上传要求
	if code := licenseService.checkLicenseFilePost(ext, fName); code != http.StatusOK {
		result.Code = code
		result.Message = "CheckLicenseFilePostErr"
		return result, fpath
	}

	// 检查文件是否存在
	if code := licenseService.checkLicenseFileIsExist(fpath, stanardLicFilename); code != http.StatusOK {
		result.Code = code
		result.Message = "CheckLicenseFileIsExistErr"
		return result, fpath
	}

	fpath = fpath + stanardLicFilename
	result.Code = http.StatusOK
	return result, fpath
}

func (this *LicenseService) CheckLicenseType() (res models.Result) {
	licenseService := LicenseService{}
	var fpath = licenseService.GetLicenseFilePath()
	var result models.Result

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
		licenseService := LicenseService{}
		licConfig := new(models.LicenseConfig)
		licConfig.Type = models.LicType_TrialLicense
		licData := licConfig.Get()
		if licData.Code == http.StatusOK {
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
	}

	logs.Info("Import trial license file ok. ")
	return result
}

func (this *LicenseService) GetLicenseFilePath() string {
	return beego.AppConfig.String("license::LicensePath")
}

func (this *LicenseService) licenseRSADecrypt(licenseByte []byte) models.Result {
	result := models.Result{}
	result.Code = http.StatusOK

	decodeBytes, err := base64.StdEncoding.DecodeString(string(licenseByte))
	if err != nil {
		result.Code = utils.LicenseBase64DecodeErr
		logs.Error("License import fail, license base64 decode err,  %s", err)
	}
	plainText := utils.RSA_Decrypt(decodeBytes, "conf/private.pem")

	result.Data = plainText
	return result
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
