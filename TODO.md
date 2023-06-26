# todo:
## sorted:
- [ ] Basic admin panel to view / remove users
- [ ] User Roles (admin, and customizable roles)
    - [ ] Middleware to validate required roles (maybe through context that gets read by token validator middleware?)
        - Rename `ValidateUserTokenMiddleware` to `SecureEndpointMiddleware` which can take a custom `SecureOptions` struct to define allowed roles and such
    - [ ] Let admin UI modify roles
    - Roles will be hardcoded with privileges in the code, modifiable by changing `roles.AvailableRoles map[string][]string`
- [ ] Disallow public registration flag
    - [ ] This will need a way to bootstrap the first user. Maybe push a URL to the CLI on first boot w/ an access token to register an admin user?
    - Set this on `register.RegistrationHandler{}` "DisablePublicRegistration"
- [ ] Create users in the Admin UI, they'll need to setup their own password
    - [ ] This should give them a unique registration link to send to individuals.
    - [ ] Allow user creation through locksmith.InviteUser(email string) (inviteURL string)
- [ ] Make URL redirects dynamic
    - e.g. modify redirect for successful auth, let API endpoints be changed for components
    - Maybe set this on `LoginHandler{}` and `RegistrationHandler{}`?
- [ ] Password Policies set on Server
    - Set on `RegistrationHandler{}`
- [ ] Secure registration & login endpoints
    - Only allow requests from the same domain.
    - Flag to DISABLE this on `LoginHandler{}`
- [ ] Multifactor support
    - [ ] WebAuthn support
    - [ ] TOTP support
    - [ ] Option to require MFA on sign up
- [ ] User lockout policy (invalid passwords, etc)
    - Set on `LoginHandler{}`
- [ ] Password reset
- [ ] Convert `InjectDatabaseIntoContext` into passing the DB into the `RegistrationHandler{}` and `LoginHandler{}`
- [ ] Attach a "Device Cookie" to the token system. If a user logs in on the same device, log them out of their previous session.
- [ ] Save last active time to Session
## unsorted:
- [ ] Encrypt User info
    - Allow specific User interface keys to be encrypted before getting sent to the database
        - Make this dynamic so any struct can also have explicit encryption
        - Maybe do this with custom struct metadata tags, e.g. `Address string 'encrypt:"true"'`
            - The custom Mongo handler would check if encryption is enabled, and if so, use the specified `CryptoEngine{}` to encrypt/decrypt data automatically
    - HashiCorp Vault integration
        - Use "derived keys" so each item has a unique encryption key
- [ ] Track IP of each login
    - Attach to session token, if it changes require a relogin
        - Maybe only require a relogin for specific privilege levels (like administrators) to help aid in times where IPs may be changing frequently (a user on cell data)
    - Once a session is expired, log it as "IPs used" to track suspicion level of new logins.
- [ ] Auto-renew session tokens
    - If the token is "soon-to-be" expired, issue a new token on the refresh after it's been validated. Automatically delete the current token and replace it with the new one.
    - Customizable settings in `SecureEndpointMiddleware`:
        - Enable feature
            - Only enable if they are a specific role (e.g. maybe let `users` refresh automatically, but require `admins` to login)
                - Maybe let `users` disable this if they'd like to enhance their own security
        - Session-based refresh period (how long before the token expires should a new one be issued?)
- [ ] Login sus levels
    - Using IP data, create a "trust" level based on if the IP being used to login has been used before, approx. geolocation difference from commonly used IPs, etc.
- [ ] Maximum Username and Password lengths to prevent overflow attacks (?). It should be large for passwords- like 256, and shorter for usernames (maybe 64 by default)
    - Can be modified on `RegistrationHandler{}`
- [ ] Sign tokens
- [ ] Create a flag for the Login API to return verbose error messages (username incorrect / password incorrect) or secure messages ("username or password incorrect")
    - This needs to also reflect status codes getting sent.
- [ ] Inherited encrypted values on registration
    - e.g. for kanban board, automatically give new users access to the HashiCorp Vault API key (maybe different keys for each user?) and encrypt w/ user password.
        - admin adds a new user -> admin creates a new encryption key in HashiCorp Vault -> admin encrypts their copy of the real API key with ephemeral hashicorp key -> send encrypted key to frontend registration -> frontend registration encrypts key w/ password -> key gets sent with registration
- [ ] Federated Logins (Google, Google Workspace, etc)
- [ ] Become an OAuth provider (allow external apps to authenticate with this service)
- [ ] Prometheus Support (exports basic metrics: # users, # login attempts, # failed logins, etc)
- [ ] Audit Logging system (easily create audit logs, allow apps to push to said audit log)
- [ ] API Token Management
    - Create tokens
    - Delete tokens
    - Track usage
    - API Token validation middleware to only allow specific endpoints to be hit if the token has a role associated

## done:
- [x] Only allow [A-z0-9] in username
- [x] Save token as cookie
- [x] Save username with token in cookie
- [x] Create mongo database tester
- [x] Let registration happen in the Web UI
- [x] Let logins happen in the Web UI
- [x] Deprecate `LocksmithUserFromMap` in favor of `LocksmithUserStruct.ReadFromMap`
    - [x] rename `LocksmithUserStruct` -> `LocksmithUserInterface`
- [x] Token validation middleware
    - Checks if they are expired & valid