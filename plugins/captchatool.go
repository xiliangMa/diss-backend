package plugins

import (
	"fmt"
	"github.com/mojocn/base64Captcha"
	"github.com/xiliangMa/diss-backend/models"
)

type CaptchaTool struct {
	Captcha models.Captcha
}

var store = base64Captcha.DefaultMemStore

//  获取验证码
func (this *CaptchaTool) GetCaptcha() {
	// 生成默认数字
	driver := base64Captcha.DefaultDriverDigit
	// 生成base64图片

	c := base64Captcha.NewCaptcha(driver, store)
	// 获取
	id, imgbit, err := c.Generate()
	if err != nil {
		fmt.Println("Get Captcha Photo with base64Captcha err:", err)
	}
	this.Captcha.Id = id
	this.Captcha.ImageData = imgbit
}

// 验证码配置
var captchaConfig = base64Captcha.DriverString{
	Height:          40,
	Width:           80,
	NoiseCount:      0,
	ShowLineOptions: 2 | 4,
	Length:          4,
	Source:          "1234567890QWERTYUPASDFGHJKLZXCVBNM",
	Fonts:           []string{"wqy-microhei.ttc"},
}

func (this *CaptchaTool) GenerateCaptcha() error {
	var driver base64Captcha.Driver
	driver = captchaConfig.ConvertFonts()

	captcha := base64Captcha.NewCaptcha(driver, store)
	id, imgbit, err := captcha.Generate()
	this.Captcha.Id = id
	this.Captcha.ImageData = imgbit

	if err != nil {
		fmt.Println("Get Captcha Photo with base64Captcha err:", err)
		return err
	}
	return nil
}

func (this *CaptchaTool) VerifyCaptcha() bool {
	if this.Captcha.Id == "" {
		return false
	}
	verifyResult := store.Verify(this.Captcha.Id, this.Captcha.Value, false)
	if verifyResult {
		store.Get(this.Captcha.Id, true)
	}
	return verifyResult
}
