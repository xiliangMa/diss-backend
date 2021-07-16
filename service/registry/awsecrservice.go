package registry

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/plugins/proxy"
	"github.com/xiliangMa/diss-backend/utils"
)

type AwsECRService struct {
	ImageConfig *models.ImageConfig
}

type awsRepositories struct {
	RepositoriesData []string `json:"repositories"`
}

func (this *AwsECRService) Auth(url string, user string, pwd string) (authorizationData *ecr.AuthorizationData, error error) {

	region := regexp.MustCompile("dkr.ecr\\.([\\w\\-]+)\\.amazonaws\\.com")
	rs := region.FindStringSubmatch(url)

	newSession, err := session.NewSession()
	config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(user, pwd, ""),
		Region:      aws.String(rs[1]),
	}

	svc := ecr.New(newSession, config)
	result, err := svc.GetAuthorizationToken(nil)
	if err != nil {
		awserr := err.(awserr.RequestFailure)
		return nil, errors.New(awserr.Message())
	}

	authorizationData = result.AuthorizationData[0]
	return authorizationData, nil

}

func (this *AwsECRService) ListRepositories() models.Result {
	var ResultData models.Result
	authorizationData, err := this.Auth(this.ImageConfig.Registry.Url, this.ImageConfig.Registry.User, this.ImageConfig.Registry.Pwd)
	if err != nil {
		ResultData.Message = err.Error()
		ResultData.Code = utils.GetNamespacesErr
		return ResultData
	}
	urls := fmt.Sprintf("%s/v2/_catalog", *authorizationData.ProxyEndpoint)
	proxy := proxy.ProxyServer{TargetUrl: urls, Token: "Basic " + *authorizationData.AuthorizationToken}
	resp, _ := proxy.Request(this.ImageConfig.Registry.User, this.ImageConfig.Registry.Pwd)

	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		body, e := ioutil.ReadAll(resp.Body)
		awsRepo := awsRepositories{}
		json.Unmarshal(body, &awsRepo)
		if e != nil {
			ResultData.Message = e.Error()
			ResultData.Code = utils.GetNamespacesErr
			return ResultData
		}
		ResultData.Data = awsRepo.RepositoriesData
		ResultData.Code = http.StatusOK
	}
	return ResultData
}

func (this *AwsECRService) Imports() (err error) {
	authorizationData, err := this.Auth(this.ImageConfig.Registry.Url, this.ImageConfig.Registry.User, this.ImageConfig.Registry.Pwd)
	if err != nil {
		return err
	}
	this.ImageConfig.Registry.Url = *authorizationData.ProxyEndpoint
	urls := fmt.Sprintf("%s/v2/%s/tags/list", this.ImageConfig.Registry.Url, this.ImageConfig.Namespaces)
	proxy := proxy.ProxyServer{TargetUrl: urls, Token: "Basic " + *authorizationData.AuthorizationToken}
	tags, _ := proxy.Request(this.ImageConfig.Registry.User, this.ImageConfig.Registry.Pwd)
	if tags.StatusCode == 200 {
		defer tags.Body.Close()
		t, _ := ioutil.ReadAll(tags.Body)
		var tagObj map[string]interface{}
		json.Unmarshal(t, &tagObj)

		if tagObj["tags"] != nil {
			go func() {
				for _, tag := range tagObj["tags"].([]interface{}) {
					this.ImageConfig.Name = this.ImageConfig.Namespaces + ":" + tag.(string)
					cs := CommonService{ImageConfig: this.ImageConfig, Token: *authorizationData.AuthorizationToken}
					cs.AddDetail()
				}
			}()
		}
	}
	return
}
