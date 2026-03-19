package config

import (
	"time"

	"github.com/hunkvm/locator/pkg/security"
)

type ClientConfig struct {
	Address string
	Timeout time.Duration
	TLS     security.ClientTLSConfig
}
