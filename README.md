# todo:
## sorted:
- [x] Only allow [A-z0-9] in username
- [x] Save token as cookie
- [x] Save username with token in cookie
- [x] Create mongo database tester
- [x] Let registration happen in the Web UI
- [x] Let logins happen in the Web UI
- [x] Token validation middleware
    - [x] Make sure they are logged in
    - [ ] Restrict access to specific endpoints depending on whether or not the logged in user "owns" the endpoint (e.g. view their own data)
- [ ] Secure registration & login endpoints
    - Only allow requests from the same domain.
- [ ] Basic admin panel to view / remove users
- [ ] User Roles (admin, and customizable roles)
    - [ ] Middleware to validate required roles (maybe through context that gets read by token validator middleware?)
    - [ ] Let admin UI modify roles
- [ ] Disallow public registration flag
    - [ ] This will need a way to bootstrap the first user. Maybe push a URL to the CLI on first boot w/ an access token to register an admin user?
- [ ] Create users in the Admin UI, they'll need to setup their own password
    - [ ] This should give them a unique registration link to send to individuals.
    - [ ] Allow user creation through locksmith.InviteUser(email string) (inviteURL string)
- [ ] Password Policies set on Server
- [ ] Multifactor support
    - [ ] WebAuthn support
    - [ ] TOTP support
    - [ ] Option to require MFA on sign up
- [ ] User lockout policy (invalid passwords, etc)

## unsorted:
- [ ] Inherited encrypted values on registration
    - e.g. for kanban board, automatically give new users access to the HashiCorp Vault API key (maybe different keys for each user?) and encrypt w/ user password.
        - admin adds a new user -> admin creates a new encryption key in HashiCorp Vault -> admin encrypts their copy of the real API key with ephemeral hashicorp key -> send encrypted key to frontend registration -> frontend registration encrypts key w/ password -> key gets sent with registration
- [ ] Federated Logins (Google, Google Workspace, etc)
- [ ] Become an OAuth provider (allow external apps to authenticate with this service)
- [ ] Prometheus Support (exports basic metrics: # users, # login attempts, # failed logins, etc)
- [ ] Audit Logging system (easily create audit logs, allow apps to push to said audit log)
