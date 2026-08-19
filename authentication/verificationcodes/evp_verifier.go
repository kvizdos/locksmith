package verificationcodes

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	evp "github.com/kvizdos/go-email-verification-protocol"
	"github.com/kvizdos/go-email-verification-protocol/evp_domain"
)

type evpVerifier struct {
	audience string
}

type evpPayload struct {
	Token string `json:"evp_token"`
	Nonce string `json:"-"`
}

var _ AutoVerifier = (*evpVerifier)(nil)

func NewEVPVerifier(audience string) *evpVerifier {
	return &evpVerifier{
		audience: audience,
	}
}

func (v *evpVerifier) Name() string {
	return "evp"
}

func (v *evpVerifier) Load(req *http.Request) AutoVerificationPayload {
	contentType := req.Header.Get("Content-Type")
	if contentType != "application/json" {
		return nil
	}
	var avp evpPayload
	decoder := json.NewDecoder(io.LimitReader(req.Body, 1<<20))
	if err := decoder.Decode(&avp); err != nil {
		fmt.Println(err)
		return nil
	}

	if cookie, err := req.Cookie("evp_nonce"); err == nil {
		if cookie.Value != "" {
			avp.Nonce = cookie.Value
		}
	}

	fmt.Printf("Loading EVP Payload: %+v\n", avp)
	return avp
}

func (v *evpVerifier) Verify(ctx context.Context, email string, r AutoVerificationPayload) error {
	payload, ok := r.(evpPayload)
	if !ok {
		return fmt.Errorf("unsupported payload type: %T", r)
	}
	verified, err := evp.Verify(ctx, payload.Token, evp_domain.VerifyOptions{
		Email:     email,
		Nonce:     payload.Nonce,
		Audience:  v.audience,
		KBMaxAge:  5 * time.Minute,
		EVTMaxAge: 5 * time.Minute,
	})

	if err == nil && verified.Verified {
		return nil
	}
	return fmt.Errorf("verification failed: %w", err)
}
