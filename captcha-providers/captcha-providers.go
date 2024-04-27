package captchaproviders

import (
	"fmt"
	"html/template"
	"net/http"
)

var UseProvider CAPTCHAProvider = NoCaptchaEnabledProvider{}

type CAPTCHAProvider interface {
	GetHeadInjection() template.HTML
	RenderCaptchaComponentCode(w http.ResponseWriter, r *http.Request)
	Validate(r *http.Request) (bool, error)
}

type NoCaptchaEnabledProvider struct{}

func (n NoCaptchaEnabledProvider) GetHeadInjection() template.HTML { return template.HTML("") }
func (n NoCaptchaEnabledProvider) RenderCaptchaComponentCode(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
}
func (n NoCaptchaEnabledProvider) Validate(r *http.Request) (bool, error) {
	fmt.Println("NO CAPTCHA PROVIDER ENABLED. Set captchaproviders.UseProvider")
	return false, fmt.Errorf("no captcha provider enabled")
}
