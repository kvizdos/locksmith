package textvalidation

import (
	"context"
	"fmt"
	"net"
	"net/mail"
	"slices"
	"strings"

	"golang.org/x/net/publicsuffix"
)

type emailValidator struct {
	commonMisspell     map[string]string
	roles              []string
	lookupMX           func(ctx context.Context, host string) ([]*net.MX, error)
	isICANN            func(ctx context.Context, host string) bool
	isDomainRegistered func(ctx context.Context, host string) (bool, error)
}

type EmailValidatorOptions struct {
	CommonMisspell     map[string]string
	LookupMX           func(ctx context.Context, host string) ([]*net.MX, error)
	IsDomainRegistered func(ctx context.Context, host string) (bool, error)
	IsICANN            func(ctx context.Context, host string) bool
	Roles              []string
}

func NewEmailValidator(opts EmailValidatorOptions) *emailValidator {
	if opts.CommonMisspell == nil {
		opts.CommonMisspell = map[string]string{
			"1-und-1.de":    "1und1.de",
			"aol.cm":        "aol.com",
			"aol.con":       "aol.com",
			"aple.com":      "apple.com",
			"frenet.de":     "freenet.de",
			"g-mail.com":    "gmail.com",
			"g.mail.com":    "gmail.com",
			"gemail.com":    "gmail.com",
			"gimail.com":    "gmail.com",
			"gmai.com":      "gmail.com",
			"gmaik.com":     "gmail.com",
			"gmail.con":     "gmail.com",
			"gmaile.com":    "gmail.com",
			"gmaill.com":    "gmail.com",
			"gmaiol.com":    "gmail.com",
			"gmal.com":      "gmail.com",
			"gmale.com":     "gmail.com",
			"gmall.com":     "gmail.com",
			"gmaol.com":     "gmail.com",
			"gmaul.com":     "gmail.com",
			"gmeil.com":     "gmail.com",
			"gmial.com":     "gmail.com",
			"gmil.com":      "gmail.com",
			"gmmail.com":    "gmail.com",
			"gmsil.com":     "gmail.com",
			"gmx.ta":        "gmx.at",
			"gnail.com":     "gmail.com",
			"goglemail.com": "googlemail.com",
			"googel.com":    "google.com",
			"googl.com":     "google.com",
			"googlmail.com": "googlemail.com",
			"googlr.com":    "google.com",
			"gxm.com":       "gmx.com",
			"gxm.de":        "gmx.de",
			"gxm.net":       "gmx.net",
			"hormail.com":   "hotmail.com",
			"hotmail.con":   "hotmail.com",
			"hotmail.ml":    "hotmail.nl",
			"hotmaill.com":  "hotmail.com",
			"hotmaill.nl":   "hotmail.nl",
			"hotmal.com":    "hotmail.com",
			"hotmal.nl":     "hotmail.nl",
			"hotmale.com":   "hotmail.com",
			"hotmale.nl":    "hotmail.nl",
			"hotmial.com":   "hotmail.com",
			"hotmial.nl":    "hotmail.nl",
			"hotmil.com":    "hotmail.com",
			"hotmil.nl":     "hotmail.nl",
			"hotnail.com":   "hotmail.com",
			"hotnail.nl":    "hotmail.nl",
			"iclaud.com":    "icloud.com",
			"icloud.con":    "icloud.com",
			"icoud.com":     "icloud.com",
			"kpnmail.ml":    "kpnmail.nl",
			"live.con":      "live.com",
			"oulook.com":    "outlook.com",
			"outlok.com":    "outlook.com",
			"outlook.con":   "outlook.com",
			"outlouk.com":   "outlook.com",
			"t.online.de":   "t-online.de",
			"tonline.de":    "t-online.de",
			"wbe.de":        "web.de",
			"we.de":         "web.de",
			"web.de.de":     "web.de",
			"yahho.com":     "yahoo.com",
			"yaho.co.uk":    "yahoo.co.uk",
			"yaho.co":       "yahoo.com",
			"yaho.com":      "yahoo.com",
			"yaho.de":       "yahoo.de",
			"yahoo.cmo":     "yahoo.com",
			"yahoomail.com": "yahoo.com",
			"yahu.com":      "yahoo.com",
			"yhoo.com":      "yahoo.com",
			"yshoo.com":     "yahoo.com",
		}
	}

	if opts.LookupMX == nil {
		opts.LookupMX = func(ctx context.Context, host string) ([]*net.MX, error) {
			return net.DefaultResolver.LookupMX(ctx, host)
		}
	}

	if opts.IsICANN == nil {
		opts.IsICANN = func(ctx context.Context, host string) bool {
			_, isICANN := publicsuffix.PublicSuffix(host)
			return isICANN
		}
	}

	if opts.IsDomainRegistered == nil {
		opts.IsDomainRegistered = func(ctx context.Context, host string) (bool, error) {
			ns, err := net.DefaultResolver.LookupNS(ctx, host)
			if err != nil {
				// NXDOMAIN → likely not registered
				if dnsErr, ok := err.(*net.DNSError); ok && dnsErr.IsNotFound {
					return false, nil
				}
				return false, err // transient DNS error
			}
			return len(ns) > 0, nil
		}
	}

	if opts.Roles == nil {
		opts.Roles = []string{
			// generic / shared
			"admin", "administrator", "root",
			"support", "help", "helpdesk", "service",
			"info", "contact", "hello",
			"team", "office", "staff",

			// business / ops
			"sales", "billing", "payments", "accounts",
			"finance", "legal", "compliance",
			"hr", "careers", "jobs",

			// technical
			"dev", "engineering", "eng", "it", "security",
			"noreply", "no-reply", "donotreply", "do-not-reply",

			// marketing / comms
			"marketing", "press", "media", "newsletter", "news",

			// misc common
			"webmaster", "postmaster", "abuse", "hostmaster",
		}
	}

	return &emailValidator{
		commonMisspell:     opts.CommonMisspell,
		lookupMX:           opts.LookupMX,
		isICANN:            opts.IsICANN,
		roles:              opts.Roles,
		isDomainRegistered: opts.IsDomainRegistered,
	}
}

func (e emailValidator) isMisspelledDomain(domain string) (string, bool) {
	if misspelled, ok := e.commonMisspell[domain]; ok {
		return misspelled, true
	}
	return "", false
}

func (e emailValidator) Validate(ctx context.Context, email string) (ValidationResultEvaluator, error) {
	email = strings.TrimSpace(email)
	validationResult := EmailValidationResult{}

	validationResult.HasValidSyntax = Bool(true)

	parsed, err := mail.ParseAddress(email)
	if err != nil {
		validationResult.HasValidSyntax = Bool(false)
		return validationResult, nil
	}

	at := strings.LastIndex(parsed.Address, "@")
	if at <= 0 || at == len(parsed.Address)-1 {
		validationResult.HasValidSyntax = Bool(false)
		return validationResult, nil
	}
	username := strings.ToLower(parsed.Address[:at])
	hostname := strings.ToLower(parsed.Address[at+1:])

	if slices.Contains(e.roles, username) {
		validationResult.IsRoleAddress = Bool(true)
	}

	// Make sure hostname contains a dot (.)
	if !strings.Contains(hostname, ".") {
		validationResult.HasValidSyntax = Bool(false)
		return validationResult, nil
	}

	if !e.isICANN(ctx, hostname) {
		validationResult.HasValidPublicSuffix = Bool(false)
		return validationResult, nil
	}
	validationResult.HasValidPublicSuffix = Bool(true)

	if correction, ok := e.isMisspelledDomain(hostname); ok {
		validationResult.HasCommonDomainMisspelling = &map[string]string{hostname: correction}
		out := fmt.Sprintf("%s@%s", username, correction)
		validationResult.DidYouMean = &out
		return validationResult, nil
	}

	validationResult.HasValidMXRecords = Bool(true)
	mxRecords, err := e.lookupMX(ctx, hostname)
	if err != nil || len(mxRecords) == 0 {
		validationResult.HasValidMXRecords = Bool(false)
		out, err := e.isDomainRegistered(ctx, hostname)
		if err == nil {
			validationResult.IsLikelyRegisteredDomain = Bool(out)
		}
		return validationResult, nil
	} else {
		validationResult.IsLikelyRegisteredDomain = Bool(true)
	}
	stringedMX := make([]string, len(mxRecords))
	for i, mx := range mxRecords {
		stringedMX[i] = mx.Host
	}
	validationResult.MXProviders = &stringedMX

	return validationResult, nil
}
