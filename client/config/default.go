package config

import (
	"time"

	"github.com/hunkvm/locator/pkg/security"
)

func DefaultClientConfig() ClientConfig {
	return ClientConfig{
		Address: "localhost:8882",
		Timeout: 60 * time.Second,
		TLS: security.ClientTLSConfig{
			TLSPolicy: security.TLSPolicyNone,
		},
	}
}
