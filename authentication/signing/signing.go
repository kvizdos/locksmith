package signing

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"math/big"
)

type SigningPackageInterface interface {
	Sign(string) (string, error)
	Validate(string, string) bool
	MarshalPrivate() (string, error)
}

type SigningPackage struct {
	Public  *ecdsa.PublicKey
	Private *ecdsa.PrivateKey
}

func (sp SigningPackage) Sign(input string) (string, error) {
	hashedData := sha256.Sum256([]byte(input))
	r, s, err := ecdsa.Sign(rand.Reader, sp.Private, hashedData[:])
	if err != nil {
		return "", err
	}
	signature := append(r.Bytes(), s.Bytes()...)

	sigBase64 := base64.StdEncoding.EncodeToString(signature)

	return sigBase64, nil
}

func (sp SigningPackage) Validate(input string, signature string) bool {
	// Decode the base64 signature
	sigBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false
	}

	// Extract R and S values from the signature
	r := new(big.Int).SetBytes(sigBytes[:len(sigBytes)/2])
	s := new(big.Int).SetBytes(sigBytes[len(sigBytes)/2:])

	// Hash the data
	hashedData := sha256.Sum256([]byte(input))

	// Verify the signature
	return ecdsa.Verify(sp.Public, hashedData[:], r, s)
}

func (s SigningPackage) MarshalPrivate() (string, error) {
	// 2. Serialize the private key into a byte slice
	privKeyBytes, err := x509.MarshalECPrivateKey(s.Private)
	if err != nil {
		return "", err
	}

	// 3. Encode the byte slice to a base64 string
	privKeyBase64 := base64.StdEncoding.EncodeToString(privKeyBytes)

	return privKeyBase64, nil
}
