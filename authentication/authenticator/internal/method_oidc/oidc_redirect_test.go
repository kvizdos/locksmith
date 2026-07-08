package method_oidc

import "testing"

func TestSanitizeRedirectPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		raw  string
		want string
	}{
		{name: "empty", raw: "", want: ""},
		{name: "simple relative path", raw: "/dashboard", want: "/dashboard"},
		{name: "relative path with query", raw: "/dashboard?tab=billing", want: "/dashboard?tab=billing"},
		{name: "absolute url rejected", raw: "https://evil.example.com/phish", want: ""},
		{name: "protocol-relative url rejected", raw: "//evil.example.com/phish", want: ""},
		{name: "no leading slash rejected", raw: "dashboard", want: ""},
		{name: "mailto scheme rejected", raw: "mailto:someone@example.com", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := sanitizeRedirectPath(tt.raw); got != tt.want {
				t.Fatalf("sanitizeRedirectPath(%q) = %q, want %q", tt.raw, got, tt.want)
			}
		})
	}
}
