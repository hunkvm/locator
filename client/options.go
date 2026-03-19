package client

import (
	"fmt"
	"time"

	"github.com/hunkvm/locator/client/config"
	"github.com/hunkvm/locator/pkg/security"
)

type ClientOption func(*config.ClientConfig) error

func WithAddress(address string) ClientOption {
	return func(c *config.ClientConfig) error {
		c.Address = address
		return nil
	}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *config.ClientConfig) error {
		if timeout < 0 {
			return fmt.Errorf("%w: %d", ErrInvalidTimeout, timeout)
		}
		c.Timeout = timeout
		return nil
	}
}

func WithTLS(tls security.ClientTLSConfig) ClientOption {
	return func(c *config.ClientConfig) error {
		c.TLS = tls
		return nil
	}
}

func WithConfig(cfg config.ClientConfig) ClientOption {
	return func(c *config.ClientConfig) error {
		*c = cfg
		return nil
	}
}
