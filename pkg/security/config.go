package security

import (
	"time"

	"github.com/hunkvm/locator/pkg/validation"
)

type TLSPolicy string

const (
	TLSPolicyNone TLSPolicy = "none"
	TLSPolicyMTLS TLSPolicy = "mtls"
	TLSPolicyTLS  TLSPolicy = "tls"
)

type ServerTLSConfig struct {
	TLSPolicy         TLSPolicy
	CaCertificateFile string
	CertificateFile   string
	PrivateKeyFile    string
	RefreshInterval   time.Duration
}

func (c ServerTLSConfig) Validate() []validation.Error {
	errList := make([]validation.Error, 0, 4)

	err := validation.ValidateOneOf(c.TLSPolicy, []TLSPolicy{
		TLSPolicyNone, TLSPolicyMTLS, TLSPolicyTLS,
	})
	if err != nil {
		errList = append(errList, validation.Error{
			Field: "TLSPolicy", Value: c.TLSPolicy, Reason: err.Error(),
		})
	}

	if c.TLSPolicy != TLSPolicyNone {
		if c.TLSPolicy == TLSPolicyMTLS && c.CaCertificateFile == "" {
			errList = append(errList, validation.Error{
				Field: "CaCertificateFile", Value: c.CaCertificateFile,
				Reason: "cannot be empty when policy is not disabled",
			})
		}

		if c.CertificateFile == "" {
			errList = append(errList, validation.Error{
				Field: "CertificateFile", Value: c.CertificateFile,
				Reason: "cannot be empty when policy is not disabled",
			})
		}

		if c.PrivateKeyFile == "" {
			errList = append(errList, validation.Error{
				Field: "PrivateKeyFile", Value: c.PrivateKeyFile,
				Reason: "cannot be empty when policy is not disabled",
			})
		}
	}

	return errList
}

type ClientTLSConfig struct {
	TLSPolicy         TLSPolicy
	CaCertificateFile string
	CertificateFile   string
	PrivateKeyFile    string
	AllowedHostname   string
	RefreshInterval   time.Duration
}

func (c ClientTLSConfig) Validate() []validation.Error {
	errList := make([]validation.Error, 0, 4)

	err := validation.ValidateOneOf(c.TLSPolicy, []TLSPolicy{
		TLSPolicyNone, TLSPolicyMTLS, TLSPolicyTLS,
	})
	if err != nil {
		errList = append(errList, validation.Error{
			Field: "TLSPolicy", Value: c.TLSPolicy, Reason: err.Error(),
		})
	}

	if c.TLSPolicy != TLSPolicyNone {
		if c.CaCertificateFile == "" {
			errList = append(errList, validation.Error{
				Field: "CaCertificateFile", Value: c.CaCertificateFile,
				Reason: "cannot be empty when policy is not disabled",
			})
		}

		if c.TLSPolicy == TLSPolicyMTLS {
			if c.CertificateFile == "" {
				errList = append(errList, validation.Error{
					Field: "CertificateFile", Value: c.CertificateFile,
					Reason: "cannot be empty when policy is mtls",
				})
			}

			if c.PrivateKeyFile == "" {
				errList = append(errList, validation.Error{
					Field: "PrivateKeyFile", Value: c.PrivateKeyFile,
					Reason: "cannot be empty when policy is mtls",
				})
			}

			if c.AllowedHostname == "" {
				errList = append(errList, validation.Error{
					Field: "AllowedHostname", Value: c.AllowedHostname,
					Reason: "cannot be empty when policy is mtls",
				})
			}
		}
	}

	return errList
}

type PeerTLSConfig ClientTLSConfig

func (c PeerTLSConfig) Validate() []validation.Error {
	errList := make([]validation.Error, 0, 4)

	err := validation.ValidateOneOf(c.TLSPolicy, []TLSPolicy{
		TLSPolicyNone, TLSPolicyMTLS, TLSPolicyTLS,
	})
	if err != nil {
		errList = append(errList, validation.Error{
			Field: "TLSPolicy", Value: c.TLSPolicy, Reason: err.Error(),
		})
	}

	if c.TLSPolicy != TLSPolicyNone {
		if c.CaCertificateFile == "" {
			errList = append(errList, validation.Error{
				Field: "CaCertificateFile", Value: c.CaCertificateFile,
				Reason: "cannot be empty when policy is not disabled",
			})
		}

		if c.CertificateFile == "" {
			errList = append(errList, validation.Error{
				Field: "CertificateFile", Value: c.CertificateFile,
				Reason: "cannot be empty when policy is not disabled",
			})
		}

		if c.PrivateKeyFile == "" {
			errList = append(errList, validation.Error{
				Field: "PrivateKeyFile", Value: c.PrivateKeyFile,
				Reason: "cannot be empty when policy is not disabled",
			})
		}

		if c.TLSPolicy == TLSPolicyMTLS && c.AllowedHostname == "" {
			errList = append(errList, validation.Error{
				Field: "AllowedHostname", Value: c.AllowedHostname,
				Reason: "cannot be empty when policy is mtls",
			})
		}
	}

	return errList
}
