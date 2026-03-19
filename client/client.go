package client

import (
	"net/url"

	"github.com/hunkvm/locator/client/config"
	"github.com/hunkvm/locator/client/service"
	"github.com/hunkvm/locator/pkg/security"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type LocatorClient struct {
	config      config.ClientConfig
	clientConn  *grpc.ClientConn
	credsLoader *security.GRPCCredsLoader
	Registry    *service.RegistryService
	Member      *service.MemberService
}

func NewLocatorClient(options ...ClientOption) (*LocatorClient, error) {
	config := config.DefaultClientConfig()
	for _, option := range options {
		if err := option(&config); err != nil {
			return nil, err
		}
	}
	return newLocatorClient(config)
}

func (c *LocatorClient) Close() error {
	if c.credsLoader != nil {
		c.credsLoader.Stop()
	}
	return c.clientConn.Close()
}

func newLocatorClient(cfg config.ClientConfig) (*LocatorClient, error) {
	var credsLoader *security.GRPCCredsLoader
	credentials := insecure.NewCredentials()

	if cfg.TLS.TLSPolicy != security.TLSPolicyNone {
		credsLoader = &security.GRPCCredsLoader{
			TLSPolicy:         cfg.TLS.TLSPolicy,
			CaCertificateFile: cfg.TLS.CaCertificateFile,
			CertificateFile:   cfg.TLS.CaCertificateFile,
			PrivateKeyFile:    cfg.TLS.PrivateKeyFile,
		}

		if u, err := url.Parse(cfg.Address); err == nil {
			credsLoader.ClientCredentials(u.Hostname())
		} else {
			return nil, err
		}
	}

	transport := grpc.WithTransportCredentials(credentials)
	clientConn, err := grpc.NewClient(cfg.Address, transport)
	if err != nil {
		return nil, err
	}

	registry := service.NewRegistryService(clientConn, cfg.Timeout)
	member := service.NewMemberService(clientConn, cfg.Timeout)

	return &LocatorClient{
		config:      cfg,
		credsLoader: credsLoader,
		clientConn:  clientConn,
		Registry:    registry,
		Member:      member,
	}, nil
}
