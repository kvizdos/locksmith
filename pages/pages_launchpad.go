//go:build enable_launchpad
// +build enable_launchpad

package pages

import _ "embed"

//go:embed launchpad.html
var LaunchpadPageHTML []byte
