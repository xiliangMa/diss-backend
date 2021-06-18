package registry

import (
	"bytes"
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/xiliangMa/diss-backend/models"
	"net/http"
	"regexp"
	"strings"
)

type CommonService struct {
	ImageConfig *models.ImageConfig
}

var scheme = regexp.MustCompile("(https|http)://([-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|])")

func (this *CommonService) AddDetail() {
	var ref name.Reference
	rs := scheme.FindStringSubmatch(this.ImageConfig.Registry.Url)
	ref, _ = name.ParseReference(this.ImageConfig.Name)

	if rs != nil && this.ImageConfig.Registry.Type != models.Registry_Type_DockerHub {
		this.ImageConfig.Name = strings.Replace(rs[2], "/", "", 1) + "/" + this.ImageConfig.Name
		ref, _ = name.ParseReference(this.ImageConfig.Name)
		if rs[1] == "http" {
			ref, _ = name.ParseReference(this.ImageConfig.Name, name.Insecure)
		}
	}
	if ic := this.ImageConfig.Get(); ic == nil {
		want := authn.AuthConfig{Username: this.ImageConfig.Registry.User, Password: this.ImageConfig.Registry.Pwd}
		img, errs := remote.Image(ref, remote.WithAuth(authn.FromConfig(want)))

		if errs != nil {
			logs.Error("remote Image err : %s", errs)
		}

		hash, _ := img.ConfigName()

		cf, _ := img.ConfigFile()

		digest, _ := img.Digest()

		this.ImageConfig.Id = ""

		this.ImageConfig.ImageId = hash.String()
		this.ImageConfig.CreateTime = cf.Created.UnixNano()
		this.ImageConfig.Add()

		imageDetail := models.ImageDetail{}

		imageDetail.ImageId = hash.String()

		imageDetail.Name = this.ImageConfig.Name

		if imd := imageDetail.Get(); imd == nil {
			Layer, _ := img.Layers()
			imageDetail.Layers = len(Layer)
			imageDetail.RepoDigests = digest.String()

			imageDetail.CreateTime = cf.Created.UnixNano()

			var buffer bytes.Buffer
			re := regexp.MustCompile(`[\s\p{Zs}]{2,}`)
			for _, h := range cf.History {
				str := re.ReplaceAllString(h.CreatedBy, "")
				buffer.WriteString(str + "\n")
			}
			imageDetail.Dockerfile = strings.TrimSpace(buffer.String())
			if result := imageDetail.Add(); result.Code != http.StatusOK {
				logs.Error("ImageDetail err: %s", errors.New(result.Message))
			}
		}
	}
}
