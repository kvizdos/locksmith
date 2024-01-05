package pages

import (
	_ "embed"
	"html/template"
)

type LocksmithPageStyling struct {
	StartGradient string
	EndGradient   string
	SubmitColor   string
	LogoURL       string
	ManifestURL   string
	InjectHeader  template.HTML
}

//go:embed admin.html
var AdminPageHTML []byte

//go:embed login.html
var LoginPageHTML []byte

//go:embed register.html
var RegisterPageHTML []byte

//go:embed reset-password-public.html
var ResetPasswordPageHTML []byte
