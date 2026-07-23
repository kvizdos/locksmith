package register_methods

type PasswordOptions struct {
	MinimumLength int
}

type PasswordOption func(*PasswordOptions)

func DefaultPasswordOptions() PasswordOptions {
	return PasswordOptions{}
}

func WithMinimumLength(length int) PasswordOption {
	return func(opts *PasswordOptions) {
		opts.MinimumLength = length
	}
}
