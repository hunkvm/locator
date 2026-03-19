package security

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc/credentials"
)

type GRPCCredsLoader struct {
	initOnce sync.Once
	stopOnce sync.Once
	stopMu   sync.Mutex
	stopWG   sync.WaitGroup
	stopCh   chan struct{}

	cert atomic.Pointer[tls.Certificate]
	pool atomic.Pointer[x509.CertPool]

	TLSPolicy         TLSPolicy
	CaCertificateFile string
	CertificateFile   string
	PrivateKeyFile    string
	RefreshInterval   time.Duration
	OnRefreshError    func(error)
}

func (l *GRPCCredsLoader) Init() error {
	var err error
	l.initOnce.Do(func() {
		stop := l.ensureStop()

		err = l.load()
		if err != nil {
			return
		}

		if l.RefreshInterval > 0 {
			l.stopMu.Lock()
			select {
			case <-stop:
				l.stopMu.Unlock()
				return
			default:
				l.stopWG.Add(1)
				l.stopMu.Unlock()
			}
			go l.reloadLoop(stop)
		}
	})
	return err
}

func (l *GRPCCredsLoader) ServerCredentials() (
	credentials.TransportCredentials, error,
) {
	if err := l.Init(); err != nil {
		return nil, err
	}
	tlsConfig := &tls.Config{
		GetConfigForClient: func(info *tls.ClientHelloInfo) (*tls.Config, error) {
			cfg := &tls.Config{
				ClientAuth: tls.VerifyClientCertIfGiven,
			}

			if cert := l.cert.Load(); cert != nil {
				cfg.Certificates = []tls.Certificate{*cert}
			}

			if l.TLSPolicy == TLSPolicyMTLS {
				cfg.ClientCAs = l.pool.Load()
				cfg.ClientAuth = tls.RequireAndVerifyClientCert
			}

			return cfg, nil
		},
	}

	return credentials.NewTLS(tlsConfig), nil
}

func (l *GRPCCredsLoader) ClientCredentials(
	serverName string,
) (credentials.TransportCredentials, error) {
	if err := l.Init(); err != nil {
		return nil, err
	}
	tlsConfig := &tls.Config{
		ServerName:         serverName,
		InsecureSkipVerify: true,
		VerifyConnection: func(cs tls.ConnectionState) error {
			if len(cs.PeerCertificates) == 0 {
				return fmt.Errorf("no certificates provided by peer")
			}

			opts := x509.VerifyOptions{
				Roots:         l.pool.Load(),
				DNSName:       serverName,
				Intermediates: x509.NewCertPool(),
			}
			for _, cert := range cs.PeerCertificates[1:] {
				opts.Intermediates.AddCert(cert)
			}

			_, err := cs.PeerCertificates[0].Verify(opts)
			return err
		},
		GetClientCertificate: func(
			info *tls.CertificateRequestInfo,
		) (*tls.Certificate, error) {
			if l.TLSPolicy == TLSPolicyMTLS {
				if cert := l.cert.Load(); cert != nil {
					return cert, nil
				}
			}
			return &tls.Certificate{}, nil
		},
	}

	return credentials.NewTLS(tlsConfig), nil
}

func (l *GRPCCredsLoader) Stop() {
	stop := l.ensureStop()

	l.stopMu.Lock()
	l.stopOnce.Do(func() {
		close(stop)
	})
	l.stopMu.Unlock()

	l.stopWG.Wait()
}

func (l *GRPCCredsLoader) ensureStop() chan struct{} {
	l.stopMu.Lock()
	defer l.stopMu.Unlock()

	if l.stopCh == nil {
		l.stopCh = make(chan struct{})
	}
	return l.stopCh
}

func (l *GRPCCredsLoader) load() error {
	cert, err := loadCert(l.CertificateFile, l.PrivateKeyFile)
	if err != nil {
		return fmt.Errorf("tls: load cert/key: %w", err)
	}
	pool, err := loadCAPool(l.CaCertificateFile)
	if err != nil {
		return fmt.Errorf("tls: load CA cert: %w", err)
	}
	l.cert.Store(cert)
	l.pool.Store(pool)
	return nil
}

func loadCert(certPath, keyPath string) (*tls.Certificate, error) {
	certPEM, err := os.ReadFile(certPath)
	if err != nil {
		format := "read cert file %q: %w"
		return nil, fmt.Errorf(format, certPath, err)
	}
	keyPEM, err := os.ReadFile(keyPath)
	if err != nil {
		format := "read key file %q: %w"
		return nil, fmt.Errorf(format, keyPath, err)
	}

	if cert, err := tls.X509KeyPair(certPEM, keyPEM); err != nil {
		format := "parse cert/key file (%q, %q): %w"
		return nil, fmt.Errorf(format, certPath, keyPath, err)
	} else {
		return &cert, nil
	}
}

func loadCAPool(caPath string) (*x509.CertPool, error) {
	caPEM, err := os.ReadFile(caPath)
	if err != nil {
		format := "read CA cert file %q: %w"
		return nil, fmt.Errorf(format, caPath, err)
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(caPEM) {
		format := "parse CA cert file %q: no valid PEM certs found"
		return nil, fmt.Errorf(format, caPath)
	}
	return pool, nil
}

func (l *GRPCCredsLoader) reloadLoop(stop <-chan struct{}) {
	defer l.stopWG.Done()

	ticker := time.NewTicker(l.RefreshInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := l.load(); err != nil {
				if l.OnRefreshError != nil {
					l.OnRefreshError(err)
				}
			}
		case <-stop:
			return
		}
	}
}
