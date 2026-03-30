package saml_config

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"net/http"

	"github.com/kvizdos/locksmith/authentication/saml/saml_discovery"
	"github.com/kvizdos/locksmith/authentication/saml/saml_entities"
	"github.com/kvizdos/locksmith/users"
)

type IdPConfig struct {
	EntityID         string
	SSOURL           string
	SigningCertPEM   string // PEM, public cert only
	NameIDFormats    []string
	SupportEmail     string // optional
	Organization     string // optional
	EnabledProviders []*saml_entities.SAMLProvider
	Signer           *signer

	decodeUserInto users.LocksmithUserInterface
	discovery      saml_discovery.DiscoveryFunc
}

func (cfg IdPConfig) Discover(r *http.Request, sp *saml_entities.SAMLProvider) (*saml_discovery.IdPDiscovery, error) {
	if cfg.discovery == nil {
		return nil, errors.New("no discovery function set")
	}
	return cfg.discovery(r, sp)
}

func (cfg *IdPConfig) WithDiscoveryFunc(discovery saml_discovery.DiscoveryFunc) *IdPConfig {
	cfg.discovery = discovery
	return cfg
}

func (cfg *IdPConfig) WithUserDecoder(u users.LocksmithUserInterface) *IdPConfig {
	cfg.decodeUserInto = u
	return cfg
}

func (cfg *IdPConfig) GetUserDecoder() users.LocksmithUserInterface {
	return cfg.decodeUserInto
}

func (cfg *IdPConfig) WithSigner(privateKeyPem string) *IdPConfig {
	if cfg.EntityID == "" {
		panic("missing IdP EntityID")
	}
	if cfg.SigningCertPEM == "" {
		panic("missing signing certificate PEM")
	}

	key, err := parsePrivateKey(privateKeyPem)
	if err != nil {
		panic(err)
	}

	cert, err := parseCertificate(cfg.SigningCertPEM)
	if err != nil {
		panic(err)
	}

	if key.PublicKey.N.BitLen() < 2048 {
		panic(errors.New("RSA key too small"))
	}

	s := &signer{
		EntityID: cfg.EntityID,
		Key:      key,
		Cert:     cert,
	}

	cfg.Signer = s
	return cfg

}

func parsePrivateKey(pemData string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(pemData))
	if block == nil {
		return nil, errors.New("failed to decode PEM private key")
	}

	switch block.Type {

	case "RSA PRIVATE KEY":
		// PKCS#1
		key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		return key, nil

	case "PRIVATE KEY":
		// PKCS#8
		keyAny, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}

		key, ok := keyAny.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("private key is not RSA")
		}
		return key, nil

	default:
		return nil, errors.New("unsupported private key type: " + block.Type)
	}
}

func parseCertificate(pemData string) (*x509.Certificate, error) {
	block, _ := pem.Decode([]byte(pemData))
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, errors.New("failed to decode PEM certificate")
	}

	return x509.ParseCertificate(block.Bytes)
}

type signer struct {
	EntityID string
	Key      *rsa.PrivateKey
	Cert     *x509.Certificate
}

func (s *signer) GetKeyPair() (*rsa.PrivateKey, []byte, error) {
	if s == nil || s.Key == nil || s.Cert == nil {
		return nil, nil, errors.New("signer not fully initialized")
	}

	return s.Key, s.Cert.Raw, nil
}
