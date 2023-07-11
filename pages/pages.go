package pages

import _ "embed"

//go:embed admin.html
var AdminPageHTML []byte

//go:embed login.html
var LoginPageHTML []byte

//go:embed register.html
var RegisterPageHTML []byte
