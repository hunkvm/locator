package config

import "github.com/hunkvm/locator/pkg/validation"

type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

type LogFormat string

const (
	LogFormatJSON LogFormat = "json"
	LogFormatText LogFormat = "text"
)

type LogConfig struct {
	Level  LogLevel
	Format LogFormat
}

func (c LogConfig) Validate() []validation.Error {
	errors := make([]validation.Error, 0, 2)

	err := validation.ValidateOneOf(c.Level, []LogLevel{
		LogLevelDebug, LogLevelInfo,
		LogLevelWarn, LogLevelError,
	})
	if err != nil {
		errors = append(errors, validation.Error{
			Field: "Level", Value: c.Level, Reason: err.Error(),
		})
	}

	err = validation.ValidateOneOf(c.Format, []LogFormat{
		LogFormatJSON, LogFormatText,
	})
	if err != nil {
		errors = append(errors, validation.Error{
			Field: "Format", Value: c.Format, Reason: err.Error(),
		})
	}

	return errors
}
