package oauth_dummy

import (
	"fmt"
	"net/http"
)

type DummyOAuth struct {
}

func (g DummyOAuth) GetName() string {
	return "google"
}

func (g DummyOAuth) RegisterRoutes(apiMux *http.ServeMux) {
	fmt.Println("Registered dummy oauth!")
}
