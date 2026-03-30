package saml_auth

import (
	"errors"
	"time"

	"github.com/kvizdos/locksmith/authentication/saml/saml_entities"
)

func ValidateAuthnRequest(
	providers []*saml_entities.SAMLProvider,
	req *AuthnRequest,
	now time.Time,
) (*ValidatedRequest, error) {
	if req == nil {
		return nil, errors.New("nil AuthnRequest")
	}

	if req.Issuer == "" {
		return nil, errors.New("missing Issuer")
	}

	// 1. Find matching SP by entityID
	var sp *saml_entities.SAMLProvider
	for _, p := range providers {
		if !p.Enabled {
			continue
		}
		if p.EntityID == req.Issuer {
			sp = p
			break
		}
	}

	if sp == nil {
		return nil, errors.New("unknown service provider")
	}

	// 2. ACS URL must match EXACTLY
	if req.ACSURL != "" && req.ACSURL != sp.ACSURL {
		return nil, errors.New("ACS URL mismatch")
	}

	// 3. Enforce assertion signing expectations
	if sp.WantAssertionsSigned && !req.IsSigned && sp.SigningCertPEM != nil {
		// NOTE: Not supported currently.
	}

	// 4. Time sanity (basic replay protection)
	// Allow small clock skew
	const skew = 5 * time.Minute
	if req.IssueInstant.After(now.Add(skew)) ||
		req.IssueInstant.Before(now.Add(-skew)) {
		return nil, errors.New("AuthnRequest outside allowed time window")
	}

	// 5. Request ID is required
	if req.ID == "" {
		return nil, errors.New("missing AuthnRequest ID")
	}

	return &ValidatedRequest{
		SP:        sp,
		RequestID: req.ID,
		ACSURL:    sp.ACSURL,
	}, nil
}
