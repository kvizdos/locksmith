package error_svc

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/kvizdos/locksmith/pages"
)

type Error struct {
	Header      string
	Description string
}

type ErrorCode string

type ErrorService interface {
	RegisterError(code ErrorCode, err Error)
	HandleHTTP(appName string, styling pages.LocksmithPageStyling) http.HandlerFunc
}

type errorSvc struct {
	appName          string
	styling          pages.LocksmithPageStyling
	registeredErrors map[ErrorCode]Error
}

func (es *errorSvc) RegisterError(code ErrorCode, err Error) {
	if _, ok := es.registeredErrors[code]; ok {
		panic(fmt.Errorf("Error code '%s' already registered.", code))
	}
	es.registeredErrors[code] = err
}

func (es *errorSvc) HandleHTTP(appName string, styling pages.LocksmithPageStyling) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		tmpl, err := template.New("error.html").Parse(string(pages.ErrorPageHTML))

		if err != nil {
			fmt.Println(err)
		}

		type PageData struct {
			Title            string
			Styling          pages.LocksmithPageStyling
			ErrorHeader      string
			ErrorDescription string
			ErrorCode        string
		}

		data := PageData{
			Title:            appName,
			Styling:          styling,
			ErrorHeader:      "You hit an error.",
			ErrorDescription: "That's all we know.",
			ErrorCode:        "UNKNOWN",
		}

		if data.Styling.SubmitColor == "" {
			data.Styling.SubmitColor = "#476ade"
		}

		if data.Styling.StartGradient == "" {
			data.Styling.StartGradient = "#476ade"
		}

		if data.Styling.EndGradient == "" {
			data.Styling.EndGradient = "#2744a3"
		}

		if data.Title == "" {
			data.Title = "Locksmith"
		}

		errCode := r.URL.Query().Get("code")

		if errDetails, ok := es.registeredErrors[ErrorCode(errCode)]; ok {
			data.ErrorCode = errCode
			data.ErrorHeader = errDetails.Header
			data.ErrorDescription = errDetails.Description
		}

		err = tmpl.Execute(w, data)

		if err != nil {
			log.Println("Error executing template :", err)
			return
		}
	}
}

func NewErrorSvc() *errorSvc {
	return &errorSvc{
		registeredErrors: make(map[ErrorCode]Error),
	}
}
