package register_methods

import "github.com/kvizdos/locksmith/authentication/registrationhints"

type HintOptions struct {
	Hints registrationhints.Service
}

type HintOption func(*HintOptions)

func DefaultHintOptions() HintOptions {
	return HintOptions{}
}

func WithHintService(service registrationhints.Service) HintOption {
	return func(opts *HintOptions) {
		opts.Hints = service
	}
}
