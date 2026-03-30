package saml_discovery

import (
	"fmt"

	"github.com/kvizdos/locksmith/error_svc"
)

type DiscoveryError struct {
	ErrorCode error_svc.ErrorCode

	// This is what will be logged.
	SystemError error
}

func (e DiscoveryError) Error() string {
	return fmt.Sprintf("%s", e.SystemError)
}
