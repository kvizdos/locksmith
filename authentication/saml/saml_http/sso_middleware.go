package saml_http

import (
	"bytes"
	"compress/flate"
	"context"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	saml_auth "github.com/kvizdos/locksmith/authentication/saml/internal/auth"
	saml_handlers "github.com/kvizdos/locksmith/authentication/saml/internal/handlers"
	"github.com/kvizdos/locksmith/authentication/saml/saml_entities"
)

func handleSSORequest(
	providers []*saml_entities.SAMLProvider,
) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			var (
				rawXML   []byte
				relay    string
				isSigned bool
				err      error
			)

			switch r.Method {

			case http.MethodGet:
				rawXML, isSigned, err = decodeRedirectAuthnRequest(r)
				relay = r.URL.Query().Get("RelayState")

			case http.MethodPost:
				if err := r.ParseForm(); err != nil {
					http.Error(w, "Invalid form", http.StatusBadRequest)
					return
				}
				rawXML, isSigned, err = decodePostAuthnRequest(r)
				relay = r.FormValue("RelayState")

			default:
				http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
				return
			}

			if err != nil {
				fmt.Println(err)
				http.Error(w, "Invalid SAMLRequest", http.StatusBadRequest)
				return
			}

			req, err := parseAuthnRequest(rawXML, isSigned)
			if err != nil {
				http.Error(w, "Malformed AuthnRequest", http.StatusBadRequest)
				return
			}

			validated, err := saml_auth.ValidateAuthnRequest(
				providers,
				req,
				time.Now(),
			)
			if err != nil {
				http.Error(w, "AuthnRequest rejected", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), saml_handlers.SAMLCtxKey{}, &saml_handlers.SAMLContext{
				AuthnRequest: req,
				RelayState:   relay,
				Validated:    validated,
			})

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

type authnRequestXML struct {
	XMLName      xml.Name `xml:"AuthnRequest"`
	ID           string   `xml:"ID,attr"`
	IssueInstant string   `xml:"IssueInstant,attr"`
	AssertionURL string   `xml:"AssertionConsumerServiceURL,attr"`
	Issuer       string   `xml:"Issuer"`
}

func parseAuthnRequest(raw []byte, isSigned bool) (*saml_auth.AuthnRequest, error) {
	var ar authnRequestXML
	if err := xml.Unmarshal(raw, &ar); err != nil {
		return nil, err
	}

	ts, err := time.Parse(time.RFC3339, ar.IssueInstant)
	if err != nil {
		return nil, err
	}

	return &saml_auth.AuthnRequest{
		ID:           ar.ID,
		Issuer:       strings.TrimSpace(ar.Issuer),
		ACSURL:       ar.AssertionURL,
		IssueInstant: ts,
		IsSigned:     isSigned,
		RawXML:       raw,
	}, nil
}

func decodeRedirectAuthnRequest(r *http.Request) ([]byte, bool, error) {
	// 1. URL decode FIRST
	reqParam := r.URL.Query().Get("SAMLRequest")
	if reqParam == "" {
		return nil, false, errors.New("missing SAMLRequest")
	}

	reqParam = strings.ReplaceAll(reqParam, " ", "+")

	// 2. Base64 decode (RAW, URL-safe)
	compressed, err := base64.StdEncoding.DecodeString(reqParam)
	if err != nil {
		return nil, false, fmt.Errorf("error decoding base64: %w", err)
	}

	// 2) DEFLATE (raw, no zlib header)
	reader := flate.NewReader(bytes.NewReader(compressed))
	defer reader.Close()

	xmlBytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, false, fmt.Errorf("inflate failed: %w", err)
	}

	// signature is URL-based for Redirect binding
	isSigned := r.URL.Query().Get("Signature") != ""

	return xmlBytes, isSigned, nil
}

func decodePostAuthnRequest(r *http.Request) ([]byte, bool, error) {
	req := r.FormValue("SAMLRequest")
	if req == "" {
		return nil, false, errors.New("missing SAMLRequest")
	}

	decoded, err := base64.StdEncoding.DecodeString(req)
	if err != nil {
		return nil, false, err
	}

	// POST binding uses XML signature
	isSigned := bytes.Contains(decoded, []byte("<Signature"))

	return decoded, isSigned, nil
}
