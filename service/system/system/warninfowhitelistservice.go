package system

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/service/base"
	"github.com/xiliangMa/diss-backend/utils"
	"mime/multipart"
	"net/http"
	"os"
)

type WarnWhitelistService struct {
	WhitelistData []byte
}

func (this *WarnWhitelistService) CreateWarnWhitelistDir(fpath string) {
	_, err := os.Stat(fpath)
	if os.IsNotExist(err) {
		logs.Info("Create warninginfo whitelist dir success, path: %s", fpath)
		os.MkdirAll(beego.AppConfig.String("system::WarnWhitelistPath"), os.ModePerm)
	}
}

func (this *WarnWhitelistService) Check(h *multipart.FileHeader) (models.Result, string) {
	var fpath = utils.GetWarnWhitelistPath()
	var result models.Result

	//创建目录
	this.CreateWarnWhitelistDir(fpath)

	// 后缀名不符合上传要求
	fileService := base.FileService{}
	fileType := models.TextType
	if code := fileService.CheckFilePost(h, fileType); code != http.StatusOK {
		result.Code = code
		result.Message = "Text Format Incorrect."
		return result, fpath
	}

	fpath = fpath + h.Filename
	result.Code = http.StatusOK
	return result, fpath
}
func (this *WarnWhitelistService) SaveList() models.Result {
	var result models.Result
	whitelistObj := []models.WarningWhiteList{}
	err := json.Unmarshal(this.WhitelistData, &whitelistObj)
	if err != nil {
		result.Code = utils.WarnWhitelistUnmarshalErr
		msg := fmt.Sprintf("WarnWhitelist file format or content fail, err:  %s", err)
		result.Message = msg
		result.Data = nil
		logs.Error(msg)
		return result
	}
	for _, whitelist := range whitelistObj {
		result = whitelist.Add()
		if result.Code != http.StatusOK {
			return result
		}
	}
	result.Message = fmt.Sprintf("Added whitelist success, count: %d ", len(whitelistObj))
	result.Code = http.StatusOK
	return result
}
