package config

import "github.com/hunkvm/locator/pkg/validation"

type ListenersConfig struct {
	Client string
	Peer   string
}

func (c ListenersConfig) Validate() []validation.Error {
	errors := make([]validation.Error, 0, 2)

	if err := validation.ValidateAddress(c.Client); err != nil {
		errors = append(errors, validation.Error{
			Field: "Client", Value: c.Client, Reason: err.Error(),
		})
	}

	if err := validation.ValidateAddress(c.Peer); err != nil {
		errors = append(errors, validation.Error{
			Field: "Peer", Value: c.Peer, Reason: err.Error(),
		})
	}

	return errors
}
