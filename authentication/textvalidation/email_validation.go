package textvalidation

import (
	"context"
	"encoding/json"
	"fmt"
)

type EmailValidator interface {
	Validate(ctx context.Context, email string) (ValidationResultEvaluator, error)
}

type ValidationResultEvaluator interface {
	// Returns [DidYouMean, ValidationResult]
	Result(skipDidYouMean bool) (*string, ValidationResult)
	DebugPrint(inputEmail string)
}

func Bool(b bool) *bool {
	return &b
}

func Map(m map[string]string) *map[string]string {
	return &m
}

func StringArray(s []string) *[]string {
	return &s
}

func String(s string) *string {
	return &s
}

type EmailValidationResult struct {
	// Validates against RFC spec
	HasValidSyntax *bool
	// The TLD of the email domain is not a valid public suffix
	HasValidPublicSuffix *bool
	// The domain part of the email is likely registered
	// (confirmed via NS/SOA records)
	IsLikelyRegisteredDomain *bool
	// Confirm that MX DNS records
	// exist for the domain part of the email
	HasValidMXRecords *bool
	// List of MX DNS records for the domain
	MXProviders *[]string
	// If it is a "role" (e.g. "admin@", "support@")
	IsRoleAddress *bool
	// If it is a disposable email address
	IsDisposable *bool
	// If it is a random input (e.g. "fgijdfiogjd@example.com")
	IsRandomInput *bool

	// Misspelling mapping for the domain part of the email
	HasCommonDomainMisspelling *map[string]string
	// Did You Mean for the domain part of the email
	// (this will take a misspelling -> remap to the correct domain)
	DidYouMean *string
}

type ValidationResult string

const (
	// We assume the email address is valid
	ValidationResult_VALID ValidationResult = "valid"
	// CONFIRM will just make users confirm the email they typed
	ValidationResult_CONFIRM ValidationResult = "confirm"
	// used rarely, but will completely DENY the email address
	ValidationResult_REJECT ValidationResult = "reject"
)

func (e EmailValidationResult) Result(skipDidYouMean bool) (*string, ValidationResult) {
	if e.DidYouMean != nil && !skipDidYouMean {
		return e.DidYouMean, ValidationResult_CONFIRM
	}
	if e.HasValidSyntax == nil || !*e.HasValidSyntax {
		return nil, ValidationResult_REJECT
	}
	if e.HasValidPublicSuffix != nil && !*e.HasValidPublicSuffix {
		return nil, ValidationResult_REJECT
	}

	if e.IsRoleAddress != nil && *e.IsRoleAddress {
		return nil, ValidationResult_CONFIRM
	}
	if e.HasCommonDomainMisspelling != nil && len(*e.HasCommonDomainMisspelling) > 0 {
		return nil, ValidationResult_CONFIRM
	}
	if e.IsRandomInput != nil && *e.IsRandomInput {
		return nil, ValidationResult_CONFIRM
	}
	if e.HasValidMXRecords == nil || !*e.HasValidMXRecords {
		return nil, ValidationResult_CONFIRM
	}
	return nil, ValidationResult_VALID
}

func (e EmailValidationResult) DebugPrint(inputEmail string) {
	out := map[string]interface{}{}

	out["input_email"] = inputEmail

	if e.HasValidSyntax != nil {
		out["has_valid_syntax"] = *e.HasValidSyntax
	}
	if e.HasValidPublicSuffix != nil {
		out["has_valid_public_suffix"] = *e.HasValidPublicSuffix
	}
	if e.IsLikelyRegisteredDomain != nil {
		out["is_likely_registered_domain"] = *e.IsLikelyRegisteredDomain
	}
	if e.HasValidMXRecords != nil {
		out["has_valid_mx_records"] = *e.HasValidMXRecords
	}
	if e.MXProviders != nil {
		out["mx_providers"] = *e.MXProviders
	}
	if e.IsRoleAddress != nil {
		out["is_role_address"] = *e.IsRoleAddress
	}
	if e.IsDisposable != nil {
		out["is_disposable"] = *e.IsDisposable
	}
	if e.IsRandomInput != nil {
		out["is_random_input"] = *e.IsRandomInput
	}
	if e.HasCommonDomainMisspelling != nil {
		out["common_domain_misspelling"] = *e.HasCommonDomainMisspelling
	}
	if e.DidYouMean != nil {
		out["did_you_mean"] = *e.DidYouMean
	}

	dym, result := e.Result(false)
	out["final_result"] = result
	if dym != nil {
		out["final_did_you_mean"] = *dym
	}

	pretty, err := json.Marshal(out)
	if err != nil {
		fmt.Println("debug print error:", err)
		return
	}

	fmt.Println(string(pretty))
}
