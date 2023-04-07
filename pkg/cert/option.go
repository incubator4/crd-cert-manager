package cert

import (
	v1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	cmmeta "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Option func(certificate *v1.Certificate) *v1.Certificate

func WithName(name string) Option {
	return func(certificate *v1.Certificate) *v1.Certificate {
		certificate.Name = name
		return certificate
	}
}

func WithNamespace(namespace string) Option {
	return func(certificate *v1.Certificate) *v1.Certificate {
		certificate.Namespace = namespace
		return certificate
	}
}

func WithIssuerRef(ref *cmmeta.ObjectReference) Option {
	return func(certificate *v1.Certificate) *v1.Certificate {
		certificate.Spec.IssuerRef = *ref
		return certificate
	}
}

func WithSecretName(name string) Option {
	return func(certificate *v1.Certificate) *v1.Certificate {
		certificate.Spec.SecretName = name
		return certificate
	}
}

func WithHosts(hosts []string) Option {
	return func(certificate *v1.Certificate) *v1.Certificate {
		certificate.Spec.DNSNames = hosts
		return certificate
	}
}

func NewCert(options ...Option) *v1.Certificate {
	cert := &v1.Certificate{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "cert-manager.io/v1",
			Kind:       "Certificate",
		},
	}

	for _, option := range options {
		cert = option(cert)
	}

	return cert
}
