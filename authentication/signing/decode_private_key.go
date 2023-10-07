package signing

import (
	"crypto/x509"
	"encoding/base64"
)

func DecodePrivateKey(key string) (SigningPackage, error) {
	// 1. Decode the base64 string to get the byte slice
	privKeyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return SigningPackage{}, err
	}

	// 2. Deserialize the byte slice to retrieve the private key object
	privateKey, err := x509.ParseECPrivateKey(privKeyBytes)
	if err != nil {
		return SigningPackage{}, err
	}

	return SigningPackage{
		Public:  &privateKey.PublicKey,
		Private: privateKey,
	}, nil
}
