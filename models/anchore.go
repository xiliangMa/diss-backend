package models

import (
	"encoding/base64"
	"github.com/anchore/client-go/pkg/external"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type AnchoreEngineManager struct {
	AnchoreClient *external.APIClient
}

type ImageSearchParams struct {
	ImageId     string `description:"(required: true)"`
	ImageDigest string `description:"(required: true)"`
	Ctype       string `description:"(required: false, value: os、files、npm、gem、python、java、binary、go、malware, default: os)"`
	Vtype       string `description:"(required: false, value: os、all、non-os, default: all)"`
	Mtype       string `description:"(required: false, value: manifest、docker_history、dockerfile, default: manifest)"`
}

func NewAnchoreEngineManager() *AnchoreEngineManager {
	//todo mv app.conf key to utils
	anchoreUrl := beego.AppConfig.String("anchore-engine::Url")
	userName := beego.AppConfig.String("anchore-engine::User")
	pwd := beego.AppConfig.String("anchore-engine::Pwd")
	cfg := external.NewConfiguration()
	basic := "Basic "
	cfg.AddDefaultHeader("Authorization", basic+base64.StdEncoding.EncodeToString([]byte(userName+`:`+pwd)))
	cfg.AddDefaultHeader("accept", "application/json")
	cfg.BasePath = anchoreUrl
	anchoreClient := external.NewAPIClient(cfg)
	_, _, err := anchoreClient.DefaultApi.Ping(nil)
	if err != nil {
		logs.Error("Create anchore engine client failed: err: %s.", err.Error())
		return nil
	}
	return &AnchoreEngineManager{AnchoreClient: anchoreClient}
}
