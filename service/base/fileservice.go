package base

import (
	"github.com/astaxie/beego/logs"
	"github.com/xiliangMa/diss-backend/utils"
	"mime/multipart"
	"net/http"
	"strings"
)

type FileService struct {
}

func (this *FileService) CheckFilePost(fh *multipart.FileHeader, fileType string) int {

	// Open File
	f, err := fh.Open()
	if err != nil {
		logs.Error("Open file fail, err: %s", err)
		return utils.CheckFilePostErr
	}
	defer f.Close()

	// Get the content
	datatype, err := this.GetFileContentType(f)
	if err != nil {
		logs.Error("Get the content fail, err: %s", err)
		return utils.CheckFilePostErr
	}

	if !strings.Contains(datatype, fileType) {
		return utils.CheckFileTypeErr
	}

	return http.StatusOK
}

func (this *FileService) GetFileContentType(file multipart.File) (string, error) {

	buffer := make([]byte, 512)

	contentType := ""
	_, err := file.Read(buffer)
	if err != nil {
		return contentType, err
	}

	contentType = http.DetectContentType(buffer)

	return contentType, nil
}
