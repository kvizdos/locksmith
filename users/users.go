package users

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"reflect"
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/kvizdos/locksmith/authentication"
	"github.com/kvizdos/locksmith/authentication/magic"
	"github.com/kvizdos/locksmith/database"
	"github.com/kvizdos/locksmith/roles"
	"github.com/kvizdos/locksmith/tenant"
)

type JWTs struct {
	// Used to authorize the user.
	Access string

	// Used for storing the roles / entitlements,
	// things you may want to use in the frontend
	Profile string

	// Used for the Refresh token
	Refresh string

	RefreshExpiresAt time.Time
}

type BaseValidationClaims struct {
	jwt.RegisteredClaims
	Roles        []string `json:"roles"`
	TenantID     string   `json:"tenant"`
	Entitlements []string `json:"entitlements"`
}

type LocksmithUserInterface interface {
	ValidatePassword(string) (bool, error)
	ValidateSessionToken(token string, db database.DatabaseAccessor) bool
	GeneratePasswordSession() (authentication.PasswordSession, error)
	SavePasswordSession(authentication.PasswordSession, database.DatabaseAccessor) error

	GenerateJWTCookies(issuer string, privateKey *ecdsa.PrivateKey, db database.DatabaseAccessor) (JWTs, error)
	// If you'd like to add custom info to the private Access JWT, do so here.
	FinalizeAccessJWTClaims(forTenant tenant.Tenant, incomingClaims jwt.MapClaims, db database.DatabaseAccessor) (jwt.MapClaims, error)
	// If you'd like to be able to access variables on the frontend, change these.
	FinalizeProfileJWTClaims(forTenant tenant.Tenant, incomingClaims jwt.MapClaims, db database.DatabaseAccessor) (jwt.MapClaims, error)

	ValidateAccessJWT(token string, privateKey *ecdsa.PublicKey, withClaimStruct jwt.Claims) (bool, jwt.Claims)
	ValidateProfileJWT(token string, shouldMatchJit string, publicKey *ecdsa.PublicKey) bool
	ValidateRefreshJWT(refreshJWT string, profileJWT string, publicKey *ecdsa.PublicKey) (bool, string)

	// Finally, create a User from all of the claims.
	FromAccessJWTClaims(claims jwt.Claims) LocksmithUserInterface

	GetTenantID() uuid.UUID

	// Under regular use, this variable
	// is only accessible after a Secure Middleware
	// call that has a Tenant structure attached.
	//
	// It utilizes internalTenant from LocksmithUser.
	GetTenant() tenant.Tenant
	SetTenant(tenant.Tenant) LocksmithUserInterface

	// Unlike role, entitlements should dictate specific
	// *features* that are paywalled / not accessible
	// by default.
	//
	// This is a list of Entitlement ID
	HasEntitlement(entitlementID string) bool
	GetAttachedEntitlements() []string

	GetLastLoginDate() time.Time

	// Read from Database
	ReadFromMap(*LocksmithUserInterface, map[string]interface{})
	ToMap() map[string]interface{}

	// Convert to "public" interface
	// Slimmed down version of this interface
	// with less sensitive information
	ToPublic() (PublicLocksmithUserInterface, error)

	GetRoles() []string
	HasPermission(perm string) bool

	// Magic Auth stuff
	CleanupOldMagicTokens(database.DatabaseAccessor)
	SetMagicPermissions([]string) LocksmithUserInterface
	SetMagic() LocksmithUser // Denotes the user as a "magic only" user
	CreateMagicAuthenticationCode(database.DatabaseAccessor, magic.MagicAuthenticationVariables) (string, error)
	GetMagicPermissions() []string
	IsMagic() bool

	WebAuthnID() []byte
	WebAuthnDisplayName() string
	WebAuthnName() string
	WebAuthnCredentials() []webauthn.Credential

	// WebAuthnIcon is a deprecated option.
	// Deprecated: this has been removed from the specification recommendation. Suggest a blank string.
	WebAuthnIcon() string
	AddNewWebAuthnCredential(*webauthn.Credential)

	// Getters
	GetUsername() string
	GetEmail() string
	GetID() string
	GetPasswordInfo() authentication.PasswordInfo
	GetWebAuthnSessions() []webauthn.SessionData
	GetPasswordSessions() []authentication.PasswordSession
}

// This interface is used when user structures are
// sent to the frontend
type PublicLocksmithUserInterface interface {
	FromRegular(LocksmithUserInterface) (PublicLocksmithUserInterface, error)
}

// Only used to show user data to an endpoint
// This should hide any sensitive data like
// password session info, etc
type PublicLocksmithUser struct {
	ID                 string   `json:"id"`
	Username           string   `json:"username"`
	Email              string   `json:"email"`
	ActiveSessionCount int      `json:"sessions"`
	LastActive         int64    `json:"lastActive"`
	Roles              []string `json:"roles"`
	Entitlements       []string `json:"entitlements"`
}

// Convert a LocksmithUser{} into
// the public equivalent.
func (u PublicLocksmithUser) FromRegular(user LocksmithUserInterface) (PublicLocksmithUserInterface, error) {
	publicUser := PublicLocksmithUser{}

	publicUser.Username = user.GetUsername()
	publicUser.Email = user.GetEmail()
	publicUser.ActiveSessionCount = len(user.GetPasswordSessions())
	publicUser.ID = user.GetID()
	publicUser.LastActive = -1
	publicUser.Roles = u.Roles
	publicUser.Entitlements = u.Entitlements

	return publicUser, nil
}

type LocksmithUser struct {
	ID               string                          `bson:"id"`
	TenantID         uuid.UUID                       `json:"tenantID,omitempty"`
	Username         string                          `json:"username" bson:"username"`
	Email            string                          `json:"email" bson:"email"`
	PasswordInfo     authentication.PasswordInfo     `json:"-" bson:"password"`
	WebAuthnSessions []webauthn.SessionData          `json:"-" bson:"websessions"`
	PasswordSessions authentication.PasswordSessions `json:"-" bson:"sessions"`
	Magics           magic.MagicAuthentications      `json:"-" bson:"magic"`

	Roles            []string  `json:"role" bson:"role"`
	Entitlements     []string  `json:"entitlements" bson:"entitlements"`
	MagicPermissions []string  `json:"-" bson:"-"`
	ImMagic          bool      `json:"-" bson:"-"`
	ImRegular        bool      `json:"-" bson:"-"`
	LastLogin        time.Time `json:"-" bson:"-"`

	// Never set internalTenant. Managed by Secure Middleware.
	internalTenant tenant.Tenant `json:"-" bson:"-"`
}

func (u LocksmithUser) HasPermission(perm string) bool {
	for _, roleName := range u.Roles {
		role, err := roles.GetRole(roleName)
		if err != nil {
			continue
		}

		if role.HasPermission(perm) {
			return true
		}
	}

	return false
}

func (u LocksmithUser) GetRoles() []string {
	return u.Roles
}

func (u LocksmithUser) GetAttachedEntitlements() []string {
	return u.Entitlements
}

func (u LocksmithUser) HasEntitlement(entitlementName string) bool {
	for _, entitlement := range u.Entitlements {
		if entitlement == entitlementName {
			return true
		}
	}

	return false
}

func (u LocksmithUser) GetTenant() tenant.Tenant {
	return u.internalTenant
}
func (u LocksmithUser) SetTenant(tenant tenant.Tenant) LocksmithUserInterface {
	u.internalTenant = tenant
	return u
}

func (u LocksmithUser) GetLastLoginDate() time.Time {
	return u.LastLogin.UTC()
}

func (u LocksmithUser) GetMagics() []magic.MagicAuthentication {
	return u.Magics
}

func (u LocksmithUser) ToPublic() (PublicLocksmithUserInterface, error) {
	publicUser, err := PublicLocksmithUser{}.FromRegular(u)

	return publicUser, err
}

func (u LocksmithUser) GetID() string {
	return u.ID
}

func (u LocksmithUser) GetEmail() string {
	return u.Email
}

func (u LocksmithUser) GetPasswordSessions() []authentication.PasswordSession {
	return u.PasswordSessions
}

func (u LocksmithUser) GetTenantID() uuid.UUID {
	return u.TenantID
}

func (u LocksmithUser) ToMap() map[string]interface{} {
	out := make(map[string]interface{})

	out["id"] = u.ID
	out["username"] = u.Username
	out["email"] = u.Email
	out["password"] = u.PasswordInfo.ToMap()
	out["websessions"] = map[string]interface{}{} // TODO
	out["sessions"] = u.PasswordSessions.ToMap()
	out["role"] = u.Roles
	out["magic"] = u.Magics.ToMap()
	out["entitlements"] = u.Entitlements

	if u.TenantID.String() != "" {
		out["tenant"] = u.TenantID.String()
	}

	if u.GetLastLoginDate().IsZero() {
		out["last_login"] = time.Now().UTC().Unix()
	} else {
		out["last_login"] = u.GetLastLoginDate().Unix()
	}

	return out
}

func (u LocksmithUser) ReadFromMap(writeTo *LocksmithUserInterface, user map[string]interface{}) {
	sessions := []authentication.PasswordSession{}

	if user["sessions"] != nil {
		rawSessions := user["sessions"]
		if mappedSessions, ok := rawSessions.([]map[string]interface{}); ok {
			sessions = authentication.PasswordSessions{}.FromMap(mappedSessions)
		} else {

			for _, v := range user["sessions"].([]interface{}) {
				session, ok := v.(authentication.PasswordSession)
				if !ok {
					newSession := v.(map[string]interface{})
					sessions = append(sessions, authentication.PasswordSession{
						Token:     newSession["token"].(string),
						ExpiresAt: newSession["expire"].(int64),
					})

					continue
				}
				sessions = append(sessions, session)
			}
		}
	}

	var tenantID uuid.UUID
	if tenantString, exists := user["tenant"].(string); exists {
		tenantID = uuid.MustParse(tenantString)
	}

	var passinfo authentication.PasswordInfo
	switch user["password"].(type) {
	case authentication.PasswordInfo:
		passinfo = user["password"].(authentication.PasswordInfo)
	case map[string]interface{}:
		passinfo = authentication.PasswordInfoFromMap(user["password"].(map[string]interface{}))
	}

	var magics []magic.MagicAuthentication
	if magicValue, ok := user["magic"]; ok {
		switch magicValue.(type) {
		case []magic.MagicAuthentication:
			magics = magicValue.([]magic.MagicAuthentication)
		case []interface{}:
			magics = magic.MagicsFromMap(magicValue.([]interface{}))
		}
	}

	var loginTime time.Time
	if user["last_login"] != nil {
		loginTime = time.Unix(user["last_login"].(int64), 0)
	} else {
		// If they've yet to login, set to now.
		loginTime = time.Now().UTC()
	}

	var roles []string

	switch user["role"].(type) {
	case []string:
		roles = user["role"].([]string)
	case []interface{}:
		roles = []string{}
		for _, role := range user["role"].([]interface{}) {
			roles = append(roles, role.(string))
		}
	}

	var entitlements []string

	switch user["entitlements"].(type) {
	case []string:
		entitlements = user["entitlements"].([]string)
	case []interface{}:
		entitlements = []string{}
		for _, entitlement := range user["entitlements"].([]interface{}) {
			entitlements = append(entitlements, entitlement.(string))
		}
	}

	*writeTo = LocksmithUser{
		ID:               user["id"].(string),
		Username:         user["username"].(string),
		Email:            user["email"].(string),
		Roles:            roles,
		Entitlements:     entitlements,
		PasswordInfo:     passinfo,
		PasswordSessions: sessions,
		Magics:           magics,
		LastLogin:        loginTime,
		TenantID:         tenantID,
	}
}

func (u LocksmithUser) GetUsername() string {
	return u.Username
}

func (u LocksmithUser) GetPasswordInfo() authentication.PasswordInfo {
	return u.PasswordInfo
}

func (u LocksmithUser) ValidatePassword(inputPassword string) (bool, error) {
	return authentication.ValidatePassword(u.PasswordInfo, inputPassword)
}

func (u LocksmithUser) ValidateAccessJWT(token string, publicKey *ecdsa.PublicKey, withClaimStruct jwt.Claims) (bool, jwt.Claims) {
	claimsType := reflect.TypeOf(withClaimStruct).Elem()
	claimsValue := reflect.New(claimsType).Interface().(jwt.Claims)

	parsed, err := jwt.ParseWithClaims(token, claimsValue, func(token *jwt.Token) (interface{}, error) {
		// Verify the token with the expected algorithm
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	}, jwt.WithIssuedAt(), jwt.WithValidMethods([]string{"ES256", "ES384", "ES512"}), jwt.WithExpirationRequired())

	if err != nil {
		fmt.Println(err)
		return false, nil
	}

	return parsed.Valid, parsed.Claims
}

func (u LocksmithUser) ValidateProfileJWT(token string, shouldMatchJIT string, publicKey *ecdsa.PublicKey) bool {
	parsed, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify the token with the expected algorithm
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	}, jwt.WithIssuedAt(), jwt.WithValidMethods([]string{"ES256", "ES384", "ES512"}), jwt.WithExpirationRequired())

	if err != nil || !parsed.Valid {
		return false
	}

	claims := parsed.Claims.(*jwt.RegisteredClaims)

	issuer, err := claims.GetIssuer()

	if err != nil || !parsed.Valid {
		return false
	}

	if issuer != shouldMatchJIT {
		return false
	}

	return parsed.Valid
}

// returns validated, userID
func (u LocksmithUser) ValidateRefreshJWT(refreshJWT string, profileJWT string, publicKey *ecdsa.PublicKey) (bool, string) {
	// Get Profile JWT Issuer
	parsed, err := jwt.ParseWithClaims(profileJWT, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify the token with the expected algorithm
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	}, jwt.WithIssuedAt(), jwt.WithValidMethods([]string{"ES256", "ES384", "ES512"}), jwt.WithExpirationRequired())

	if err != nil || !parsed.Valid {
		fmt.Println("Failed to validate profile jwt", err.Error())
		return false, ""
	}

	claims := parsed.Claims.(*jwt.RegisteredClaims)

	profileIssuer, err := claims.GetIssuer()

	if err != nil {
		fmt.Println("Failed to get profile issuer", err.Error())
		return false, ""
	}

	// Get Refresh JWT Issuer
	parsedRefresh, err := jwt.ParseWithClaims(refreshJWT, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify the token with the expected algorithm
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	}, jwt.WithIssuedAt(), jwt.WithValidMethods([]string{"ES256", "ES384", "ES512"}), jwt.WithExpirationRequired())

	if err != nil || !parsed.Valid {
		fmt.Println("Failed to validate refresh jwt", err.Error())
		return false, ""
	}

	refreshclaims := parsedRefresh.Claims.(*jwt.RegisteredClaims)

	refreshIssuer, err := refreshclaims.GetIssuer()

	if err != nil {
		fmt.Println("Failed to get refresh issuer", err.Error())
		return false, ""
	}

	if refreshIssuer != profileIssuer {
		fmt.Println("Issuers do NOT Match", refreshIssuer, profileIssuer)
		return false, ""
	}

	return true
}

func (u LocksmithUser) GenerateJWTCookies(issuer string, key *ecdsa.PrivateKey, db database.DatabaseAccessor) (JWTs, error) {
	// Access Claims
	frontendPermissions := []string{}

	for _, roleName := range u.Roles {
		role, err := roles.GetRole(roleName)
		if err != nil {
			continue
		}
		frontendPermissions = append(frontendPermissions, role.FrontendPermissions...)
	}

	iat := time.Now().UTC()
	nbf := iat.Unix()
	exp := iat.Add(5 * time.Minute).Unix()
	jti := uuid.NewString()

	refresh_expires := iat.AddDate(0, 0, 30)

	claims := jwt.MapClaims{
		"iss": issuer,
		"jti": jti,
		"sub": u.GetID(),
		"iat": iat.Unix(),
		"nbf": nbf,
		"exp": exp,
	}

	var tenant tenant.Tenant
	if u.GetTenantID() != uuid.Nil {
		claims["tenant"] = u.GetTenantID()
		tenant = u.GetTenant()
	}

	finalizedClaims, err := u.FinalizeAccessJWTClaims(tenant, claims, db)
	if err != nil {
		return JWTs{}, err
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodES256, finalizedClaims)

	signedAccessToken, err := accessToken.SignedString(key)

	if err != nil {
		return JWTs{}, err
	}

	// Profile Claims
	baseProfileClaims := jwt.MapClaims{
		"iss":   jti, // This needs to align with the users current Access JWT or it will be rejected!
		"jti":   uuid.NewString(),
		"sub":   u.GetID(),
		"roles": frontendPermissions,
		"iat":   iat.Unix(),
		"nbf":   nbf,
		"exp":   refresh_expires.Unix(),
	}

	finalizedProfileClaims, err := u.FinalizeProfileJWTClaims(tenant, baseProfileClaims, db)
	if err != nil {
		return JWTs{}, err
	}

	profileToken := jwt.NewWithClaims(jwt.SigningMethodES256, finalizedProfileClaims)
	signedProfileToken, err := profileToken.SignedString(key)

	if err != nil {
		return JWTs{}, err
	}

	// Refresh Claims
	refreshClaims := jwt.MapClaims{
		"iss": jti, // This needs to align with the users current Profile JWT or it will be rejected!
		"jti": uuid.NewString(),
		"sub": u.GetID(),
		"iat": iat.Unix(),
		"nbf": nbf,
		"exp": refresh_expires.Unix(),
		"aud": u.ID,
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodES256, refreshClaims)
	signedRefreshToken, err := refreshToken.SignedString(key)

	if err != nil {
		return JWTs{}, err
	}

	return JWTs{
		Access:           signedAccessToken,
		Profile:          signedProfileToken,
		Refresh:          signedRefreshToken,
		RefreshExpiresAt: refresh_expires,
	}, nil
}

// If you'd like to add custom info to the private Access JWT, do so here.
func (u LocksmithUser) FinalizeAccessJWTClaims(forTenant tenant.Tenant, incomingClaims jwt.MapClaims, db database.DatabaseAccessor) (jwt.MapClaims, error) {
	entitlements := u.Entitlements

	tenantID := u.GetTenantID()
	if tenantID != uuid.Nil {
		incomingClaims["tenant"] = tenantID.String()

		entitlements = forTenant.ConfirmUserEntitlements(u.GetAttachedEntitlements())
	}

	incomingClaims["roles"] = u.Roles
	incomingClaims["entitlements"] = entitlements
	incomingClaims["aud"] = u.ID

	return incomingClaims, nil
}

// If you'd like to be able to access variables on the frontend, change these.
func (u LocksmithUser) FinalizeProfileJWTClaims(forTenant tenant.Tenant, incomingClaims jwt.MapClaims, db database.DatabaseAccessor) (jwt.MapClaims, error) {
	incomingClaims["username"] = u.Username
	incomingClaims["aud"] = u.ID

	entitlements := u.GetAttachedEntitlements()

	tenantID := u.GetTenantID()
	if tenantID != uuid.Nil {
		entitlements = forTenant.ConfirmUserEntitlements(entitlements)
	}

	incomingClaims["entitlements"] = entitlements

	return incomingClaims, nil
}

func (u LocksmithUser) FromAccessJWTClaims(claims jwt.Claims) LocksmithUserInterface {
	baseClaims := claims.(*BaseValidationClaims)

	tenantID, _ := uuid.Parse(baseClaims.TenantID)

	return LocksmithUser{
		ID:           baseClaims.Subject,
		TenantID:     tenantID,
		Roles:        baseClaims.Roles,
		Entitlements: baseClaims.Entitlements,
	}
}

func (u LocksmithUser) GeneratePasswordSession() (authentication.PasswordSession, error) {
	token, err := authentication.GenerateRandomString(64)

	if err != nil {
		return authentication.PasswordSession{}, err
	}

	now := time.Now()
	futureDate := now.AddDate(0, 0, 30)

	timestamp := futureDate.Unix()

	session := authentication.PasswordSession{
		Token:     token,
		ExpiresAt: timestamp,
	}

	return session, nil
}

func (u LocksmithUser) GenerateCookieValueFromSession(session authentication.PasswordSession) string {
	token := fmt.Sprintf("%s:%s", session.Token, u.GetID())
	return base64.StdEncoding.EncodeToString([]byte(token))
}

func (u LocksmithUser) SavePasswordSession(session authentication.PasswordSession, db database.DatabaseAccessor) error {
	hasher := sha256.New()
	hasher.Write([]byte(session.Token))
	hashedCode := hasher.Sum(nil)
	hashedToken := fmt.Sprintf("%x", hashedCode)

	session.Token = hashedToken

	_, err := db.UpdateOne("users", map[string]interface{}{
		"username": u.Username,
	}, map[database.DatabaseUpdateActions]map[string]interface{}{
		database.PUSH: {
			"sessions": session,
		},
	})

	_, err = db.UpdateOne("users", map[string]interface{}{
		"username": u.Username,
	}, map[database.DatabaseUpdateActions]map[string]interface{}{
		database.SET: {
			"last_login": time.Now().UTC().Unix(),
		},
	})

	return err
}

func (u LocksmithUser) ValidateSessionToken(token string, db database.DatabaseAccessor) bool {
	found := false

	nonexpiredTokens := []authentication.PasswordSession{}

	hasher := sha256.New()
	hasher.Write([]byte(token))
	hashedCode := hasher.Sum(nil)
	hashedToken := fmt.Sprintf("%x", hashedCode)

	for _, session := range u.PasswordSessions {
		// Maybe renew the token here if its getting soon to expiring..
		if !session.IsExpired() {
			nonexpiredTokens = append(nonexpiredTokens, session)
			// Helps prevent timing attacks
			if subtle.ConstantTimeCompare([]byte(session.Token), []byte(hashedToken)) == 1 {
				found = true
			}
		}
	}

	updateMap := map[database.DatabaseUpdateActions]map[string]interface{}{
		database.SET: {},
	}

	if len(nonexpiredTokens) != len(u.PasswordSessions) {
		// Create a new slice of type []interface{}
		interfaces := make([]interface{}, len(nonexpiredTokens))

		// Convert each PasswordSession to an interface{}
		for i, session := range nonexpiredTokens {
			interfaces[i] = session
		}

		updateMap[database.SET]["sessions"] = interfaces
	}

	if len(updateMap[database.SET]) > 0 {
		db.UpdateOne("users", map[string]interface{}{
			"id": u.ID,
		}, updateMap)
	}

	return found
}

func (u LocksmithUser) SetMagicPermissions(permissions []string) LocksmithUserInterface {
	u.MagicPermissions = permissions
	return u
}

func (u LocksmithUser) SetMagic() LocksmithUser {
	u.ImMagic = true
	return u
}

func (u LocksmithUser) IsMagic() bool {
	return u.ImMagic && !u.ImRegular
}

func (u LocksmithUser) CreateMagicAuthenticationCode(db database.DatabaseAccessor, vars magic.MagicAuthenticationVariables) (string, error) {
	mac, identifier, err := magic.CreateMagicAuthentication(vars)

	if err != nil {
		return "", err
	}

	db.UpdateOne("users", map[string]interface{}{
		"id": u.ID,
	}, map[database.DatabaseUpdateActions]map[string]interface{}{
		database.PUSH: {
			"magic": mac.ToMap(),
		},
	})

	return identifier, nil
}

func (u LocksmithUser) GetMagicPermissions() []string {
	return u.MagicPermissions
}

func (u LocksmithUser) CleanupOldMagicTokens(db database.DatabaseAccessor) {
	if magic.MagicSigningPackage == nil {
		return
	}
	activeMagics := make(chan magic.MagicAuthentications)
	go magic.FilterActive(activeMagics, u.Magics)
	keep := <-activeMagics

	if len(keep) != len(u.Magics) {
		db.UpdateOne("users", map[string]interface{}{
			"id": u.ID,
		}, map[database.DatabaseUpdateActions]map[string]interface{}{
			database.SET: {
				"magic": keep.ToMap(),
			},
		})
	}
}

func (u LocksmithUser) WebAuthnID() []byte {
	return []byte(u.ID)
}
func (u LocksmithUser) WebAuthnDisplayName() string {
	return u.Username
}
func (u LocksmithUser) WebAuthnName() string {
	return u.Username
}
func (u LocksmithUser) WebAuthnCredentials() []webauthn.Credential {
	return u.PasswordInfo.WebAuthnCredentials
}
func (u LocksmithUser) WebAuthnIcon() string {
	// IS DEPRECATED.
	return ""
}
func (u LocksmithUser) GetWebAuthnSessions() []webauthn.SessionData {
	return u.WebAuthnSessions
}
func (u LocksmithUser) AddNewWebAuthnCredential(cred *webauthn.Credential) {
	fmt.Println("Adding credential", cred)
}
