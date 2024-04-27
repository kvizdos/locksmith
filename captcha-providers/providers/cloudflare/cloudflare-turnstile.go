package providers

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/kvizdos/locksmith/logger"
)

//go:embed captcha.component.js
var captchaComponent []byte

type TurnstileCaptcha struct {
	SiteKey   string
	SecretKey string
}

func (t TurnstileCaptcha) GetHeadInjection() template.HTML {
	return template.HTML(`
		<script src="https://challenges.cloudflare.com/turnstile/v0/api.js?render=explicit" defer></script>
		<script type="module" src="/components/captcha.component.js"></script>
		`)
}

func (t TurnstileCaptcha) RenderCaptchaComponentCode(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("turnstile-captcha.html").Parse(string(captchaComponent))

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	type CaptchaData struct {
		SiteKey string
	}

	err = tmpl.Execute(w, CaptchaData{
		SiteKey: t.SiteKey,
	})

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (t TurnstileCaptcha) Validate(r *http.Request) (bool, error) {
	type requestBody struct {
		Response string `json:"cf-turnstile-response"`
	}

	rawBytes, err := io.ReadAll(r.Body)
	if err != nil {
		// handle the error
		return false, fmt.Errorf("failed to parse body")
	}

	var captchaBody requestBody
	err = json.Unmarshal(rawBytes, &captchaBody)

	ip := logger.GetIPFromRequest(*r)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("secret", t.SecretKey)
	_ = writer.WriteField("response", captchaBody.Response)
	_ = writer.WriteField("remoteip", ip)

	err = writer.Close()
	if err != nil {
		return false, fmt.Errorf("failed to close writer")
	}

	req, err := http.NewRequest("POST", "https://challenges.cloudflare.com/turnstile/v0/siteverify", body)
	if err != nil {
		return false, fmt.Errorf("failed to create POST request for Cloudflare")
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to send POST request to Cloudflare")
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read Cloudflare response")
	}

	type cfResponse struct {
		Success bool `json:"success"`
	}

	var captchaResponse cfResponse
	json.Unmarshal(respBody, &captchaResponse)

	return captchaResponse.Success, nil
}
