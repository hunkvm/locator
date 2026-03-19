package config

import (
	"time"

	"github.com/hunkvm/locator/pkg/validation"
)

type HealthCheckConfig struct {
	MinBackoffDuration time.Duration
	MaxBackoffDuration time.Duration
	MinConnectTimeout  time.Duration
	MaxBackoffAttempts int
}

func (c HealthCheckConfig) Validate() []validation.Error {
	errors := make([]validation.Error, 0, 4)
	minOfMinBackoffDuration := 500 * time.Millisecond

	if err := validation.ValidateDuration(c.MinBackoffDuration,
		&minOfMinBackoffDuration, nil); err != nil {
		errors = append(errors, validation.Error{
			Field:  "MinBackoffDuration",
			Value:  c.MinBackoffDuration,
			Reason: err.Error(),
		})
	} else {
		minOfMinBackoffDuration = c.MinBackoffDuration
	}

	if err := validation.ValidateDuration(c.MaxBackoffDuration,
		&minOfMinBackoffDuration, new(30*time.Minute)); err != nil {
		errors = append(errors, validation.Error{
			Field:  "MaxBackoffDuration",
			Value:  c.MaxBackoffDuration,
			Reason: err.Error(),
		})
	}

	if err := validation.ValidateDuration(c.MinConnectTimeout,
		new(100*time.Millisecond), new(5*time.Minute)); err != nil {
		errors = append(errors, validation.Error{
			Field:  "MinConnectTimeout",
			Value:  c.MinConnectTimeout,
			Reason: err.Error(),
		})
	}

	if err := validation.ValidateInteger(c.MaxBackoffAttempts,
		new(0), new(500)); err != nil {
		errors = append(errors, validation.Error{
			Field:  "MaxBackoffAttempts",
			Value:  c.MaxBackoffAttempts,
			Reason: err.Error(),
		})
	}

	return errors
}
