package system

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/astaxie/beego"
	"github.com/xiliangMa/diss-backend/models"
)

const (
	probeDriverRepo = "upload/probe"
	falcoDriverRepo = "https://download.falco.org/driver"
)

// Falco Driver Download（驱动下载）
type ProbeDriverController struct {
	beego.Controller
}

// @Title DownloadProbeDriver
// @Description downloads the drivers file
// @Param token header string true "authorized token"
// @Param version path string "" "probe driver version"
// @Param filename path string "" true "probe driver filename"
// @Success 200 {object} models.Result
// @router /downloads/probe/:version/:filename [get]
func (this *ProbeDriverController) DownloadProbeDriver() {
	var (
		version  string
		filename string
		err      error
	)

	version = this.GetString(":version")
	filename = this.GetString(":filename")

	if len(version) == 0 || len(filename) == 0 {
		this.Data["json"] = models.Result{
			Code:    http.StatusBadRequest,
			Message: "probe driver version or filename is missing",
		}

		this.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		this.ServeJSON(false)

		return
	}

	if err = this.download(filename, version); err != nil {
		this.Data["json"] = models.Result{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		}

		log.Println(err)
		this.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		this.ServeJSON(false)

		return
	}

	this.Ctx.Output.Download(path.Join(probeDriverRepo, version, filename))
}

func (this *ProbeDriverController) download(driver, version string) error {
	var (
		deployInternet bool = false
	)

	if v := os.Getenv("ANYI_DEPLOY_INTERNET"); strings.ToUpper(v) == "TRUE" {
		deployInternet = true
	}

	driverName := strings.Replace(driver, "probe", "falco", -1)
	driverPath := path.Join(probeDriverRepo, version, driver)

	if _, err := os.Stat(driverPath); err == nil {
		return nil
	}

	if !deployInternet {
		return errors.New("driver not found, please get in touch with AnYi (Beijing) Company")
	}

	url := fmt.Sprintf("%s/%s/%s", falcoDriverRepo, version, driverName)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(http.StatusText(resp.StatusCode))
	}

	file, err := os.Create(driverPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)

	return err

}
