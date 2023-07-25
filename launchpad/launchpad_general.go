package launchpad

import "github.com/kvizdos/locksmith/database"

type LocksmithLaunchpadUserOptions struct {
	// Name to show on Launchpad
	DisplayName string
	// Email to inherit
	Email string
	// Role that user will obtain
	Role string
	// Redirect Path on sucecssful login
	Redirect string
}

type LocksmithLaunchpadOptions struct {
	// Is the Launchpad enabled?
	Enabled bool
	// Caption to show under title in
	// Web Launchpad UI
	Caption string
	// What users are available to the Launchpad?
	Users map[string]LocksmithLaunchpadUserOptions
	// Access token to view the Launchpad
	// and use the Users
	AccessToken string
	// Bootstrap any Demo Database items
	BootstrapDatabase func(database.DatabaseAccessor)
}
