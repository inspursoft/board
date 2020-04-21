package commons

type captchaInfo struct {
	ID string `json:"captcha_id"`
}

type CaptchaController struct {
	BaseController
}

func (c *CaptchaController) Prepare() {

}

func (c *CaptchaController) Get() {
	captchaID, err := Cpt.CreateCaptcha()
	if err != nil {
		c.InternalError(err)
	}
	c.RenderJSON(captchaInfo{ID: captchaID})
}
