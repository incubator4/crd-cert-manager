package provider

import (
	"fmt"
	cmmeta "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
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

type Cert interface {
	GetName() string
	GetSecretName() string
	GetHosts() []string
	GetAnnotations() map[string]string
}

//var _ Cert = (*DefaultCert)(nil)

type DefaultCert struct {
}

func (d DefaultCert) GetName() string {
	return "default"
}

func (d DefaultCert) GetSecretName() string {
	//TODO implement me
	panic("implement me")
}

func (d DefaultCert) GetHosts() []string {
	//TODO implement me
	panic("implement me")
}

type ObjectMeta interface {
	GetAnnotations() map[string]string
	GetName() string
}

func GetIssuer(obj ObjectMeta) (*cmmeta.ObjectReference, error) {
	if name := obj.GetAnnotations()[IngressClusterIssuerNameAnnotationKey]; name != "" {
		return &cmmeta.ObjectReference{
			Name: name,
			Kind: "ClusterIssuer",
		}, nil
	} else if name := obj.GetAnnotations()[IngressIssuerNameAnnotationKey]; name != "" {
		return &cmmeta.ObjectReference{
			Name: name,
			Kind: "Issuer",
		}, nil
	} else {
		return nil, fmt.Errorf("not found cert issue annotations")
	}
}

func CertName(obj ObjectMeta) string {
	if name := obj.GetAnnotations()[CommonNameAnnotationKey]; name != "" {
		return obj.GetAnnotations()[CommonNameAnnotationKey]
	}
	return fmt.Sprintf("%s-cert", obj.GetName())
}
