package validation

var useSigningKeys ValidationSigningKeys

func SetSigningKeys(keys ValidationSigningKeys) {
	useSigningKeys = keys
}

func GetSigningKeys() ValidationSigningKeys {
	return useSigningKeys
}
