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
	// Is the Launchpad enabled? This will be
	// ignored if the build tag is not present.
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
	// What label the "Refresh Environment" button
	// will show.
	RefreshButtonText string
	// Setting this to TRUE will make the
	// launchpad buttons RED to notify
	// users that it is an EARLY preview
	// before it hits an official staging
	// environment.
	IsEarlyDevelopmentEnvironment bool
}
