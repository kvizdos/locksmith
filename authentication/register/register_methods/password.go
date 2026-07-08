package register_methods

type PasswordOptions struct{}

type PasswordOption func(*PasswordOptions)

func DefaultPasswordOptions() PasswordOptions {
	return PasswordOptions{}
}
