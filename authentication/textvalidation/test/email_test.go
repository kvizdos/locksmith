package textvalidation_test

import (
	"context"
	"net"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kvizdos/locksmith/authentication/textvalidation"
)

func TestEmailValidation(t *testing.T) {
	testcases := []struct {
		name               string
		email              string
		misspellings       map[string]string
		roles              []string
		isICANN            func(context.Context, string) bool
		lookupMX           func(context.Context, string) ([]*net.MX, error)
		isDomainRegistered func(context.Context, string) (bool, error)

		expectedResult        textvalidation.EmailValidationResult
		expectedDYM           *string
		expectedDecision      textvalidation.ValidationResult
		expectedDecisionNoDYM textvalidation.ValidationResult
	}{
		{
			name:  "valid email with mx",
			email: "demo@kv.codes",
			isICANN: func(ctx context.Context, host string) bool {
				return host == "kv.codes"
			},
			lookupMX: func(ctx context.Context, host string) ([]*net.MX, error) {
				return []*net.MX{
					{Host: "ch.protonmail.com."},
				}, nil
			},
			expectedResult: textvalidation.EmailValidationResult{
				HasValidSyntax:           textvalidation.Bool(true),
				HasValidPublicSuffix:     textvalidation.Bool(true),
				HasValidMXRecords:        textvalidation.Bool(true),
				IsLikelyRegisteredDomain: textvalidation.Bool(true),
				MXProviders:              textvalidation.StringArray([]string{"ch.protonmail.com."}),
			},
			expectedDecision:      textvalidation.ValidationResult_VALID,
			expectedDecisionNoDYM: textvalidation.ValidationResult_VALID,
		},
		{
			name:  "misspelled domain suggests correction",
			email: "demo@gmal.com",
			misspellings: map[string]string{
				"gmal.com": "gmail.com",
			},
			isICANN: func(ctx context.Context, host string) bool { return true },
			lookupMX: func(ctx context.Context, host string) ([]*net.MX, error) {
				t.Fatalf("lookupMX should not be called for misspelled domain")
				return nil, nil
			},
			expectedResult: textvalidation.EmailValidationResult{
				HasValidSyntax:       textvalidation.Bool(true),
				HasValidPublicSuffix: textvalidation.Bool(true),
				HasCommonDomainMisspelling: textvalidation.Map(map[string]string{
					"gmal.com": "gmail.com",
				}),
				DidYouMean: textvalidation.String("demo@gmail.com"),
			},
			expectedDYM:           textvalidation.String("demo@gmail.com"),
			expectedDecision:      textvalidation.ValidationResult_CONFIRM,
			expectedDecisionNoDYM: textvalidation.ValidationResult_CONFIRM,
		},
		{
			name:  "missing tld rejects",
			email: "demo@demo",
			isICANN: func(ctx context.Context, host string) bool {
				t.Fatalf("isICANN should not be called when hostname has no dot")
				return false
			},
			lookupMX: func(ctx context.Context, host string) ([]*net.MX, error) {
				t.Fatalf("lookupMX should not be called when syntax is invalid")
				return nil, nil
			},
			expectedResult: textvalidation.EmailValidationResult{
				HasValidSyntax: textvalidation.Bool(false),
			},
			expectedDecision:      textvalidation.ValidationResult_REJECT,
			expectedDecisionNoDYM: textvalidation.ValidationResult_REJECT,
		},
		{
			name:  "invalid public suffix rejects",
			email: "demo@example.invalidtld",
			isICANN: func(ctx context.Context, host string) bool {
				return false
			},
			lookupMX: func(ctx context.Context, host string) ([]*net.MX, error) {
				t.Fatalf("lookupMX should not be called for invalid tld")
				return nil, nil
			},
			expectedResult: textvalidation.EmailValidationResult{
				HasValidSyntax:       textvalidation.Bool(true),
				HasValidPublicSuffix: textvalidation.Bool(false),
			},
			expectedDecision:      textvalidation.ValidationResult_REJECT,
			expectedDecisionNoDYM: textvalidation.ValidationResult_REJECT,
		},
		{
			name:  "invalid syntax rejects",
			email: "not-an-email",
			expectedResult: textvalidation.EmailValidationResult{
				HasValidSyntax: textvalidation.Bool(false),
			},
			expectedDecision:      textvalidation.ValidationResult_REJECT,
			expectedDecisionNoDYM: textvalidation.ValidationResult_REJECT,
		},
		{
			name:  "empty email rejects",
			email: "   ",
			expectedResult: textvalidation.EmailValidationResult{
				HasValidSyntax: textvalidation.Bool(false),
			},
			expectedDecision:      textvalidation.ValidationResult_REJECT,
			expectedDecisionNoDYM: textvalidation.ValidationResult_REJECT,
		},
		{
			name:  "trimmed valid email",
			email: "  demo@kv.codes  ",
			isICANN: func(ctx context.Context, host string) bool {
				return host == "kv.codes"
			},
			lookupMX: func(ctx context.Context, host string) ([]*net.MX, error) {
				return []*net.MX{{Host: "ch.protonmail.com."}}, nil
			},
			expectedResult: textvalidation.EmailValidationResult{
				HasValidSyntax:           textvalidation.Bool(true),
				HasValidPublicSuffix:     textvalidation.Bool(true),
				HasValidMXRecords:        textvalidation.Bool(true),
				IsLikelyRegisteredDomain: textvalidation.Bool(true),
				MXProviders:              textvalidation.StringArray([]string{"ch.protonmail.com."}),
			},
			expectedDecision:      textvalidation.ValidationResult_VALID,
			expectedDecisionNoDYM: textvalidation.ValidationResult_VALID,
		},
		{
			name:  "uppercase domain misspelling still suggests correction",
			email: "Demo@GMAL.COM",
			misspellings: map[string]string{
				"gmal.com": "gmail.com",
			},
			isICANN: func(ctx context.Context, host string) bool { return true },
			lookupMX: func(ctx context.Context, host string) ([]*net.MX, error) {
				t.Fatalf("lookupMX should not be called")
				return nil, nil
			},
			expectedResult: textvalidation.EmailValidationResult{
				HasValidSyntax:       textvalidation.Bool(true),
				HasValidPublicSuffix: textvalidation.Bool(true),
				HasCommonDomainMisspelling: textvalidation.Map(map[string]string{
					"gmal.com": "gmail.com",
				}),
				DidYouMean: textvalidation.String("demo@gmail.com"),
			},
			expectedDYM:           textvalidation.String("demo@gmail.com"),
			expectedDecision:      textvalidation.ValidationResult_CONFIRM,
			expectedDecisionNoDYM: textvalidation.ValidationResult_CONFIRM,
		},
		{
			name:    "role address confirms",
			email:   "support@example.com",
			isICANN: func(ctx context.Context, host string) bool { return true },
			lookupMX: func(ctx context.Context, host string) ([]*net.MX, error) {
				return []*net.MX{{Host: "mx.example.com."}}, nil
			},
			expectedResult: textvalidation.EmailValidationResult{
				HasValidSyntax:           textvalidation.Bool(true),
				HasValidPublicSuffix:     textvalidation.Bool(true),
				IsRoleAddress:            textvalidation.Bool(true),
				HasValidMXRecords:        textvalidation.Bool(true),
				IsLikelyRegisteredDomain: textvalidation.Bool(true),
				MXProviders:              textvalidation.StringArray([]string{"mx.example.com."}),
			},
			expectedDecision:      textvalidation.ValidationResult_CONFIRM,
			expectedDecisionNoDYM: textvalidation.ValidationResult_CONFIRM,
		},
		{
			name:    "role detection disabled with empty roles slice",
			email:   "support@example.com",
			roles:   []string{},
			isICANN: func(ctx context.Context, host string) bool { return true },
			lookupMX: func(ctx context.Context, host string) ([]*net.MX, error) {
				return []*net.MX{{Host: "mx.example.com."}}, nil
			},
			expectedResult: textvalidation.EmailValidationResult{
				HasValidSyntax:           textvalidation.Bool(true),
				HasValidPublicSuffix:     textvalidation.Bool(true),
				HasValidMXRecords:        textvalidation.Bool(true),
				IsLikelyRegisteredDomain: textvalidation.Bool(true),
				MXProviders:              textvalidation.StringArray([]string{"mx.example.com."}),
			},
			expectedDecision:      textvalidation.ValidationResult_VALID,
			expectedDecisionNoDYM: textvalidation.ValidationResult_VALID,
		},
		{
			name:    "no mx but registered domain confirms",
			email:   "demo@example.com",
			isICANN: func(ctx context.Context, host string) bool { return true },
			lookupMX: func(ctx context.Context, host string) ([]*net.MX, error) {
				return nil, nil
			},
			isDomainRegistered: func(ctx context.Context, host string) (bool, error) {
				return true, nil
			},
			expectedResult: textvalidation.EmailValidationResult{
				HasValidSyntax:           textvalidation.Bool(true),
				HasValidPublicSuffix:     textvalidation.Bool(true),
				HasValidMXRecords:        textvalidation.Bool(false),
				IsLikelyRegisteredDomain: textvalidation.Bool(true),
			},
			expectedDecision:      textvalidation.ValidationResult_CONFIRM,
			expectedDecisionNoDYM: textvalidation.ValidationResult_CONFIRM,
		},
		{
			name:    "no mx and unregistered domain",
			email:   "demo@example.com",
			isICANN: func(ctx context.Context, host string) bool { return true },
			lookupMX: func(ctx context.Context, host string) ([]*net.MX, error) {
				return nil, nil
			},
			isDomainRegistered: func(ctx context.Context, host string) (bool, error) {
				return false, nil
			},
			expectedResult: textvalidation.EmailValidationResult{
				HasValidSyntax:           textvalidation.Bool(true),
				HasValidPublicSuffix:     textvalidation.Bool(true),
				HasValidMXRecords:        textvalidation.Bool(false),
				IsLikelyRegisteredDomain: textvalidation.Bool(false),
			},
			expectedDecision:      textvalidation.ValidationResult_CONFIRM,
			expectedDecisionNoDYM: textvalidation.ValidationResult_CONFIRM,
		},
		{
			name:    "mx lookup error and registration unknown",
			email:   "demo@example.com",
			isICANN: func(ctx context.Context, host string) bool { return true },
			lookupMX: func(ctx context.Context, host string) ([]*net.MX, error) {
				return nil, context.DeadlineExceeded
			},
			isDomainRegistered: func(ctx context.Context, host string) (bool, error) {
				return false, context.DeadlineExceeded
			},
			expectedResult: textvalidation.EmailValidationResult{
				HasValidSyntax:       textvalidation.Bool(true),
				HasValidPublicSuffix: textvalidation.Bool(true),
				HasValidMXRecords:    textvalidation.Bool(false),
			},
			expectedDecision:      textvalidation.ValidationResult_CONFIRM,
			expectedDecisionNoDYM: textvalidation.ValidationResult_CONFIRM,
		},
		{
			name:    "default misspelling map works",
			email:   "demo@gmal.com",
			isICANN: func(ctx context.Context, host string) bool { return true },
			lookupMX: func(ctx context.Context, host string) ([]*net.MX, error) {
				t.Fatalf("lookupMX should not be called")
				return nil, nil
			},
			expectedResult: textvalidation.EmailValidationResult{
				HasValidSyntax:       textvalidation.Bool(true),
				HasValidPublicSuffix: textvalidation.Bool(true),
				HasCommonDomainMisspelling: textvalidation.Map(map[string]string{
					"gmal.com": "gmail.com",
				}),
				DidYouMean: textvalidation.String("demo@gmail.com"),
			},
			expectedDYM:           textvalidation.String("demo@gmail.com"),
			expectedDecision:      textvalidation.ValidationResult_CONFIRM,
			expectedDecisionNoDYM: textvalidation.ValidationResult_CONFIRM,
		},
		{
			name:  "multiple at signs rejects",
			email: "a@b@c.com",
			expectedResult: textvalidation.EmailValidationResult{
				HasValidSyntax: textvalidation.Bool(false),
			},
			expectedDecision:      textvalidation.ValidationResult_REJECT,
			expectedDecisionNoDYM: textvalidation.ValidationResult_REJECT,
		},
		{
			name:    "display name is accepted and parsed",
			email:   "Kent <demo@example.com>",
			isICANN: func(ctx context.Context, host string) bool { return true },
			lookupMX: func(ctx context.Context, host string) ([]*net.MX, error) {
				return []*net.MX{{Host: "mx.example.com."}}, nil
			},
			expectedResult: textvalidation.EmailValidationResult{
				HasValidSyntax:           textvalidation.Bool(true),
				HasValidPublicSuffix:     textvalidation.Bool(true),
				HasValidMXRecords:        textvalidation.Bool(true),
				IsLikelyRegisteredDomain: textvalidation.Bool(true),
				MXProviders:              textvalidation.StringArray([]string{"mx.example.com."}),
			},
			expectedDecision:      textvalidation.ValidationResult_VALID,
			expectedDecisionNoDYM: textvalidation.ValidationResult_VALID,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			opts := textvalidation.EmailValidatorOptions{
				CommonMisspell: tc.misspellings,
				LookupMX:       tc.lookupMX,
				IsICANN:        tc.isICANN,
				Roles:          tc.roles,
			}
			if tc.isDomainRegistered != nil {
				opts.IsDomainRegistered = tc.isDomainRegistered
			}

			validator := textvalidation.NewEmailValidator(opts)

			gotEval, err := validator.Validate(context.Background(), tc.email)
			if err != nil {
				t.Fatalf("Validate(%q) error: %v", tc.email, err)
			}

			got, ok := gotEval.(textvalidation.EmailValidationResult)
			if !ok {
				t.Fatalf("Validate(%q) returned %T, want EmailValidationResult", tc.email, gotEval)
			}

			if diff := cmp.Diff(tc.expectedResult, got); diff != "" {
				t.Errorf("Validate(%q) mismatch (-want +got):\n%s", tc.email, diff)
			}

			gotDYM, gotDecision := got.Result(false)
			if diff := cmp.Diff(tc.expectedDYM, gotDYM); diff != "" {
				t.Errorf("Validate(%q) did you mean mismatch (-want +got):\n%s", tc.email, diff)
			}
			if diff := cmp.Diff(tc.expectedDecision, gotDecision); diff != "" {
				t.Errorf("Validate(%q) decision mismatch (-want +got):\n%s", tc.email, diff)
			}

			_, gotDecisionNoDYM := got.Result(true)
			if diff := cmp.Diff(tc.expectedDecisionNoDYM, gotDecisionNoDYM); diff != "" {
				t.Errorf("Validate(%q) skip-did-you-mean decision mismatch (-want +got):\n%s", tc.email, diff)
			}
		})
	}
}
