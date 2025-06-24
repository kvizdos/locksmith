package jwts

import (
	"context"
	"net/http"
	"time"
)

var jwts = map[string]RegisteredJWT{}

type RegisteredJWT struct {
	RequiredPermission string
	AttachPermissions  []string
	ForAudience        string
	Issuer             string
	ExpiresIn          time.Duration
	SigningKey         func(context.Context) (string, error)
	ExtraClaims        func(context.Context, *http.Request) map[string]any
}

func RegisterJWT(name string, jwt RegisteredJWT) {
	jwts[name] = jwt
}

func GetJWT(name string) (RegisteredJWT, bool) {
	jwt, ok := jwts[name]
	return jwt, ok
}
