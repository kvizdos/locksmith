package components

import (
	_ "embed"
	"net/http"

	captchaproviders "github.com/kvizdos/locksmith/captcha-providers"
	"github.com/kvizdos/locksmith/launchpad"
)

//go:embed register.component.js
var RegistrationComponentJS []byte

//go:embed reset-password.component.js
var ResetPasswordComponentJS []byte

//go:embed signin.component.js
var SigninComponentJS []byte

//go:embed user-list.component.js
var UserListComponentJS []byte

//go:embed user-tab.component.js
var UserTabComponentJS []byte

//go:embed persona-switcher.component.js
var PersonaSwitcherJS []byte

//go:embed ephemeral_tokens.js
var EphemeralTokensJS []byte

func ServeComponents(w http.ResponseWriter, r *http.Request) {
	component := r.URL.Path[len("/components/"):]
	switch component {
	case "ephemeral_tokens.js":
		serveJSComponent(w, EphemeralTokensJS)
	case "register.component.js":
		serveJSComponent(w, RegistrationComponentJS)
	case "signin.component.js":
		serveJSComponent(w, SigninComponentJS)
	case "user-list.component.js":
		serveJSComponent(w, UserListComponentJS)
	case "user-tab.component.js":
		serveJSComponent(w, UserTabComponentJS)
	case "reset-password.component.js":
		serveJSComponent(w, ResetPasswordComponentJS)
	case "persona-switcher.component.js":
		if launchpad.IS_ENABLED {
			serveJSComponent(w, PersonaSwitcherJS)
		} else {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	case "captcha.component.js":
		w.Header().Set("Content-Type", "application/javascript")
		captchaproviders.UseProvider.RenderCaptchaComponentCode(w, r)
		return
	default:
		http.NotFound(w, r)
	}
}
func serveJSComponent(w http.ResponseWriter, jsData []byte) {
	w.Header().Set("Content-Type", "application/javascript")
	w.Write(jsData)
}
