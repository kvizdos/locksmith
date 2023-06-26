package crypto

import "fmt"

type CryptoEngineInterface interface {
	GetKeyForOperation(operation string) (string, error)

	Encrypt(key string, value string) (string, error)
	Decrypt(key string, value string) (string, error)
	Sign(key string, value string) (string, error)
}

type CryptoTestEngine struct {
	Keys map[string]string
}

func (c CryptoTestEngine) GetKeyForOperation(operation string) (string, error) {
	key, notFound := c.Keys[operation]

	if !notFound {
		return "", fmt.Errorf("invalid operation")
	}

	return key, nil
}

func (c CryptoTestEngine) Encrypt(key string, value string) (string, error) {
	return fmt.Sprintf("encrypted:%s", value), nil
}
func (c CryptoTestEngine) Decrypt(key string, value string) (string, error) {
	return fmt.Sprintf("decrypted:%s", value), nil
}
func (c CryptoTestEngine) Sign(key string, value string) (string, error) {
	return fmt.Sprintf("signed:%s", value), nil
}
