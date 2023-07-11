package components

import (
	_ "embed"
	"net/http"
)

//go:embed register.component.js
var RegistrationComponentJS []byte

//go:embed signin.component.js
var SigninComponentJS []byte

//go:embed user-list.component.js
var UserListComponentJS []byte

//go:embed user-tab.component.js
var UserTabComponentJS []byte

func ServeComponents(w http.ResponseWriter, r *http.Request) {
	component := r.URL.Path[len("/components/"):]
	switch component {
	case "register.component.js":
		serveJSComponent(w, RegistrationComponentJS)
	case "signin.component.js":
		serveJSComponent(w, SigninComponentJS)
	case "user-list.component.js":
		serveJSComponent(w, UserListComponentJS)
	case "user-tab.component.js":
		serveJSComponent(w, UserTabComponentJS)
	default:
		http.NotFound(w, r)
	}
}
func serveJSComponent(w http.ResponseWriter, jsData []byte) {
	w.Header().Set("Content-Type", "application/javascript")
	w.Write(jsData)
}
