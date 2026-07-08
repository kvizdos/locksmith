package authenticator_methods

type PasswordValidatorOptions struct {
	MinPasswordLength int
}

func DefaultPasswordValidatorOptions() PasswordValidatorOptions {
	return PasswordValidatorOptions{
		MinPasswordLength: 8,
	}
}

type PasswordValidatorOption func(*PasswordValidatorOptions)

func RequireMinPasswordLength(length int) PasswordValidatorOption {
	return func(opts *PasswordValidatorOptions) {
		opts.MinPasswordLength = length
	}
}
