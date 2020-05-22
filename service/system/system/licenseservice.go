package system

import (
	"encoding/base64"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
	"mime/multipart"
	"net/http"
	"os"
	"path"
)

func License_RSA_Decrypt(licenseByte []byte, isUpdate bool) models.Result {
	result := models.Result{}
	result.Code = http.StatusOK

	decodeBytes, err := base64.StdEncoding.DecodeString(string(licenseByte))
	if err != nil {
		result.Code = utils.LicenseBase64DecodeErr
		logs.Error("License import fail, license base64 decode err,  %s", err)
	}
	plainText := utils.RSA_Decrypt(decodeBytes, "conf/private.pem")

	licenseObject := new(models.LicenseConfig)
	err = json.Unmarshal(plainText, &licenseObject)
	if err != nil {
		result.Code = utils.LicenseUnmarshalErr
		logs.Error("License import fail, Unmarshal license err,  %s", err)
	} else {
		if isUpdate {
			result = licenseObject.Update()
			//添加license 历史
			liceseHistory := models.LicenseHistory{}
			liceseHistory.LicenseJson = string(plainText)
			liceseHistory.Add()
			logs.Info("Force update license file success, License: %s", string(plainText))
		} else {
			result = licenseObject.Add()
			//添加license 历史
			liceseHistory := models.LicenseHistory{}
			liceseHistory.LicenseJson = string(plainText)
			liceseHistory.Add()
			logs.Info("License import success, License：", string(plainText))
		}

	}
	return result
}

func CheckLicenseFilePost(ext, fName string) int {
	var AllowExtMap map[string]bool = map[string]bool{
		".lic": true,
	}
	if _, ok := AllowExtMap[ext]; !ok {
		return utils.CheckLicenseFilePostErr
	}
	return http.StatusOK
}

func CheckLicenseFileIsExist(fpath, fName string) int {
	if _, err := os.Stat(fpath + fName); err == nil {
		return utils.CheckLicenseFileIsExistErr
	}
	return http.StatusOK
}

func CheckLicenseFile(h *multipart.FileHeader) (models.Result, string) {
	var fpath = getLicenseFilePath()
	var result models.Result
	fName := h.Filename
	ext := path.Ext(fName)

	//创建目录
	CreateKubeConfigDir(fpath)

	// 后缀名不符合上传要求
	if code := CheckLicenseFilePost(ext, fName); code != http.StatusOK {
		result.Code = code
		result.Message = "CheckLicenseFilePostErr"
		return result, fpath
	}

	// 检查文件是否存在
	if code := CheckLicenseFileIsExist(fpath, fName); code != http.StatusOK {
		result.Code = code
		result.Message = "CheckLicenseFileIsExistErr"
		return result, fpath
	}

	fpath = fpath + h.Filename
	result.Code = http.StatusOK
	return result, fpath
}

func getLicenseFilePath() string {
	return beego.AppConfig.String("license::LicensePath")
}
