package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-viper/mapstructure/v2"
	"github.com/hunkvm/locator/pkg/security"
	"github.com/hunkvm/locator/pkg/validation"
	"github.com/spf13/viper"
)

type ServerConfig struct {
	Listeners             ListenersConfig
	ClientChannelSecurity security.ServerTLSConfig
	PeerChannelSecurity   security.PeerTLSConfig
	DataDir               string
	Raft                  RaftConfig
	HealthCheck           HealthCheckConfig
	Logger                LogConfig
}

func Load(path string) (ServerConfig, error) {
	v := viper.New()
	v.SetConfigFile(path)

	replacer := strings.NewReplacer("-", "_")
	v.SetEnvKeyReplacer(replacer)
	v.SetEnvPrefix("LOCATOR")

	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		format := "read config file: %w"
		return ServerConfig{}, fmt.Errorf(format, err)
	}

	conf := defaultServerConfig()
	var meta mapstructure.Metadata

	decoderConfig := &mapstructure.DecoderConfig{
		Result:           &conf,
		Metadata:         &meta,
		WeaklyTypedInput: true,
		Squash:           true,
		MatchName: func(mapKey, fieldName string) bool {
			return mapKey == kebabCase(fieldName)
		},
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.TextUnmarshallerHookFunc(),
			mapstructure.StringToTimeDurationHookFunc(),
		),
	}

	decoder, err := mapstructure.NewDecoder(decoderConfig)
	if err != nil {
		return ServerConfig{}, fmt.Errorf("create decoder: %w", err)
	}
	if err := decoder.Decode(v.AllSettings()); err != nil {
		fieldErrors := extractDecodeErrors(err)
		errors := make([]string, len(fieldErrors))
		for i, e := range fieldErrors {
			errors[i] = formatError("", e)
		}

		if msg := strings.Join(errors, "; "); msg != "" {
			format := "decode configuration failed: [%s]"
			return conf, fmt.Errorf(format, msg)
		}
	}

	unused := make([]string, len(meta.Unused))
	for i, u := range meta.Unused {
		unused[i] = kebabCase(u)
	}

	if msg := strings.Join(unused, "; "); len(unused) > 0 {
		format := "unknown configuration: [%s]"
		return conf, fmt.Errorf(format, msg)
	}

	errors := validateServerConfig(conf)
	if msg := strings.Join(errors, "; "); msg != "" {
		format := "invalid configuration: [%s]"
		return conf, fmt.Errorf(format, msg)
	}

	return conf, nil
}

func validateServerConfig(cfg ServerConfig) (errors []string) {
	if path := strings.TrimSpace(cfg.DataDir); path == "" {
		errors = append(errors, "data-dir cannot be empty")
	}

	css, pcs := cfg.ClientChannelSecurity, cfg.PeerChannelSecurity
	errors = append(errors, validate("listeners", cfg.Listeners)...)
	errors = append(errors, validate("client-channel-security", css)...)
	errors = append(errors, validate("peer-channel-security", pcs)...)
	errors = append(errors, validate("raft", cfg.Raft)...)
	errors = append(errors, validate("health-check", cfg.HealthCheck)...)
	errors = append(errors, validate("logger", cfg.Logger)...)
	return
}

func validate(prefix string, conf validation.Validatable) (errs []string) {
	for _, e := range conf.Validate() {
		errs = append(errs, formatError(prefix, e))
	}
	return errs
}

func formatError(prefix string, e validation.Error) string {
	field := kebabCase(e.Field)
	if strings.TrimSpace(prefix) != "" {
		field = prefix + "." + field
	}

	var builder strings.Builder
	fmt.Fprintf(&builder, "%s: %s", field, e.Reason)
	if e.Value != nil {
		fmt.Fprintf(&builder, ", got: %v", e.Value)
	}
	return builder.String()
}

func extractDecodeErrors(err error) []validation.Error {
	if out := collectDecodeErrors(err); len(out) == 0 {
		return []validation.Error{{Reason: err.Error()}}
	} else {
		return out
	}
}

func collectDecodeErrors(err error) []validation.Error {
	if err == nil {
		return nil
	}

	var out []validation.Error

	if multi, ok := err.(interface{ Unwrap() []error }); ok {
		for _, e := range multi.Unwrap() {
			out = append(out, collectDecodeErrors(e)...)
		}
		return out
	}

	if decodeErr, ok := err.(*mapstructure.DecodeError); ok {
		cause := decodeErr.Unwrap()
		out = append(out, validation.Error{
			Field:  decodeErr.Name(),
			Reason: formatDecodeReason(cause),
		})
		return out
	}

	if single, ok := err.(interface{ Unwrap() error }); ok {
		return append(out, collectDecodeErrors(single.Unwrap())...)
	}

	return out
}

func formatDecodeReason(err error) string {
	var parseErr *mapstructure.ParseError
	if ok := errors.As(err, &parseErr); ok {
		typeName := parseErr.Expected.Type().String()
		format := "cannot parse value as '%s'"
		return fmt.Sprintf(format, typeName)
	}

	return err.Error()
}
