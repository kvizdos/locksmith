package jwts

import (
	"context"
	"maps"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type IssueJWTOptions struct {
	Sub         string
	Now         time.Time // optional, default to time.Now().UTC()
	ExtraClaims map[string]any
	Req         *http.Request // optional
}

func (j RegisteredJWT) IssueJWT(ctx context.Context, opts IssueJWTOptions) (string, error) {
	signingKey, err := j.SigningKey(ctx)
	if err != nil {
		return "", err
	}

	claims := map[string]any{}

	// Pull from request if provided
	if j.ExtraClaims != nil && opts.Req != nil {
		maps.Copy(claims, j.ExtraClaims(ctx, opts.Req))
	}

	// Explicit extras override request extras
	maps.Copy(claims, opts.ExtraClaims)

	now := opts.Now
	if now.IsZero() {
		now = time.Now().UTC()
	}

	// Standard claims
	claims["jti"] = uuid.New().String()
	claims["sub"] = opts.Sub
	claims["aud"] = j.ForAudience
	claims["iss"] = j.Issuer
	claims["scope"] = j.AttachPermissions
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()
	claims["exp"] = now.Add(j.ExpiresIn).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))
	return token.SignedString([]byte(signingKey))
}
