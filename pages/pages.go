package pages

import _ "embed"

type LocksmithPageStyling struct {
	StartGradient string
	EndGradient   string
	SubmitColor   string
	LogoURL       string
}

//go:embed admin.html
var AdminPageHTML []byte

//go:embed login.html
var LoginPageHTML []byte

//go:embed register.html
var RegisterPageHTML []byte
