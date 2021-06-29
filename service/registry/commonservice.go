package registry

import (
	"bytes"
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/xiliangMa/diss-backend/models"
	"github.com/xiliangMa/diss-backend/utils"
	"net/http"
	"regexp"
	"strings"
)

type CommonService struct {
	ImageConfig *models.ImageConfig
	Token       string
}

var scheme = regexp.MustCompile("(https|http)://([-A-Za-z0-9+&@#/%?=~_|!:,.;]+[-A-Za-z0-9+&@#/%=~_|])")

func (this *CommonService) AddDetail() {
	var ref name.Reference
	rs := scheme.FindStringSubmatch(this.ImageConfig.Registry.Url)
	ref, _ = name.ParseReference(this.ImageConfig.Name)

	if rs != nil && this.ImageConfig.Registry.Type != models.Registry_Type_DockerHub {
		if this.ImageConfig.Namespaces != "" && this.ImageConfig.Registry.Type != models.Registry_Type_AwsECR {
			this.ImageConfig.Name = strings.Replace(rs[2], "/", "", 1) + "/" + this.ImageConfig.Namespaces + "/" + this.ImageConfig.Name
		} else {
			this.ImageConfig.Name = strings.Replace(rs[2], "/", "", 1) + "/" + this.ImageConfig.Name
		}
		ref, _ = name.ParseReference(this.ImageConfig.Name)
		if rs[1] == "http" {
			ref, _ = name.ParseReference(this.ImageConfig.Name, name.Insecure)
		}
	}
	if ic := this.ImageConfig.Get(); ic == nil {
		want := authn.AuthConfig{Username: this.ImageConfig.Registry.User, Password: this.ImageConfig.Registry.Pwd}
		if this.Token != "" {
			want = authn.AuthConfig{Auth: this.Token}
		}

		img, err := remote.Image(ref, remote.WithAuth(authn.FromConfig(want)))

		if err != nil {
			logs.Error("remote Image err : %s", err)
		}

		hash, _ := img.ConfigName()

		cf, _ := img.ConfigFile()

		digest, _ := img.Digest()

		layer, _ := img.Layers()

		var sum int64

		for _, l := range layer {
			size, _ := l.Size()
			sum += size
		}

		this.ImageConfig.Id = ""
		this.ImageConfig.ImageId = hash.String()
		this.ImageConfig.Size = utils.FormatFileSize(sum)
		this.ImageConfig.CreateTime = cf.Created.UnixNano()
		this.ImageConfig.Add()

		imageDetail := models.ImageDetail{}
		imageDetail.ImageId = hash.String()
		imageDetail.Name = this.ImageConfig.Name
		imageDetail.ImageConfigId = this.ImageConfig.Id

		if imd := imageDetail.Get(); imd == nil {
			imageDetail.Layers = len(layer)
			imageDetail.RepoDigests = digest.String()
			imageDetail.Size = this.ImageConfig.Size
			imageDetail.CreateTime = cf.Created.UnixNano()

			var buffer bytes.Buffer
			var trimnop = regexp.MustCompile(`^/bin/sh\s+-c\s+#\(nop\)\s+`)
			var trimrun = regexp.MustCompile(`^(RUN\s+){0,1}/bin/sh\s+-c\s+`)
			for _, h := range cf.History {
				tmpstr := trimnop.ReplaceAllString(h.CreatedBy, "")
				tmpstrs := trimrun.ReplaceAllString(tmpstr, "RUN ")
				buffer.WriteString(tmpstrs + "\n")
			}
			imageDetail.Dockerfile = strings.TrimSpace(buffer.String())
			if result := imageDetail.Add(); result.Code != http.StatusOK {
				logs.Error("ImageDetail err: %s", errors.New(result.Message))
			}
		}
	}
}
