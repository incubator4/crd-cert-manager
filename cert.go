package main

import (
	"fmt"
	v1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	cmmeta "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
	"github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/traefik/v1alpha1"
)

const (
	// IngressIssuerNameAnnotationKey holds the issuerNameAnnotation value which can be
	// used to override the issuer specified on the created Certificate resource.
	IngressIssuerNameAnnotationKey = "cert-manager.io/issuer"
	// IngressClusterIssuerNameAnnotationKey holds the clusterIssuerNameAnnotation value which
	// can be used to override the issuer specified on the created Certificate resource. The Certificate
	// will reference the specified *ClusterIssuer* instead of normal issuer.
	IngressClusterIssuerNameAnnotationKey = "cert-manager.io/cluster-issuer"
	// CommonNameAnnotationKey Annotation key for certificate common name.
	CommonNameAnnotationKey = "cert-manager.io/common-name"
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

func GetCertName(ingressRoute *v1alpha1.IngressRoute) string {
	if ingressRoute.Annotations[CommonNameAnnotationKey] != "" {
		return ingressRoute.Annotations[CommonNameAnnotationKey]
	}
	return fmt.Sprintf("%s-cert", ingressRoute.Name)
}

func GetIssuerRef(annotations map[string]string) (*cmmeta.ObjectReference, error) {
	if annotations[IngressClusterIssuerNameAnnotationKey] != "" {
		return &cmmeta.ObjectReference{
			Name: annotations[IngressClusterIssuerNameAnnotationKey],
			Kind: "ClusterIssuer",
		}, nil
	} else if annotations[IngressIssuerNameAnnotationKey] != "" {
		return &cmmeta.ObjectReference{
			Name: annotations[IngressIssuerNameAnnotationKey],
			Kind: "Issuer",
		}, nil
	} else {
		return nil, fmt.Errorf("not found cert issue annotations")
	}
}
