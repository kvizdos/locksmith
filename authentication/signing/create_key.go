package signing

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
)

func CreateSigningPackage() (SigningPackage, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Println("Error generating key:", err)
		return SigningPackage{}, err
	}

	return SigningPackage{
		Public:  &privateKey.PublicKey,
		Private: privateKey,
	}, nil
}
