package xsrf

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/kvizdos/locksmith/authentication"
	"github.com/kvizdos/locksmith/authentication/signing"
)

type XSRFSigningPackages struct {
	Anonymous     *signing.SigningPackage // Used for Login XSRF
	Authenticated *signing.SigningPackage // Used for Authenticated requests
}

var XSRFSigningPackage = XSRFSigningPackages{}

func GenerateXSRFForSession(sessionIDToAttach string, timeToLive time.Duration) (string, error) {
	// HASH & Sign:
	// - ExpiresAt
	// - Attached Session ID
	// - XSRF Token Value
	//
	// Set Cookie To
	// - ExpiresAt + "," + XSRF Token + "," + Signature

	expiresAt := time.Now().UTC().Add(timeToLive)
	XSRFTokenBytes, err := authentication.GenerateRandomBytes(24)
	XSRFToken := base64.RawStdEncoding.EncodeToString(XSRFTokenBytes)

	if err != nil {
		return "", err
	}

	hasher := sha256.New()
	hasher.Write([]byte(fmt.Sprintf("%d.%s.%s", expiresAt.Unix(), XSRFToken, sessionIDToAttach)))
	hashedToken := fmt.Sprintf("%x", hasher.Sum(nil))

	signed, err := XSRFSigningPackage.Anonymous.Sign(hashedToken)

	if err != nil {
		return "", err
	}

	base64Encoded := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d.%s.%s", expiresAt.Unix(), XSRFToken, signed)))

	return base64Encoded, nil
}

func Confirm(xsrfValue string, sessionID string) bool {
	decoded, err := base64.StdEncoding.DecodeString(xsrfValue)

	if err != nil {
		fmt.Println(err)
		return false
	}

	decodedXSRFToken := string(decoded)
	xsrfValues := strings.Split(decodedXSRFToken, ".")

	if len(xsrfValues) != 3 {
		return false
	}

	xsrfExpiresAt := xsrfValues[0]
	epochInt, err := strconv.ParseInt(xsrfExpiresAt, 10, 64)
	if err != nil {
		fmt.Printf("Error converting epoch string to int: %v\n", err)
		return false
	}

	if epochInt < time.Now().UTC().Unix() {
		fmt.Println("XSRF Expired!")
		return false
	}

	hasher := sha256.New()
	hasher.Write([]byte(fmt.Sprintf("%d.%s.%s", epochInt, xsrfValues[1], sessionID)))
	hashedToken := hasher.Sum(nil)

	// Validate the signature.
	isValid := XSRFSigningPackage.Anonymous.Validate(hex.EncodeToString(hashedToken), xsrfValues[2])
	if !isValid {
		fmt.Println("invalid signature used")
		return false
	}

	return true
}
