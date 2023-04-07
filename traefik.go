package main

import (
	"context"
	"encoding/json"
	"fmt"
	v1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	cmVersioned "github.com/cert-manager/cert-manager/pkg/client/clientset/versioned"
	log "github.com/sirupsen/logrus"
	"github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/traefik/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func handleIngressRoute(
	ctx context.Context,
	clientSet cmVersioned.Interface,
	ingressRoute *v1alpha1.IngressRoute) {

	issuer, err := GetIssuerRef(ingressRoute.Annotations)
	if err != nil {
		return
	}

	if ingressRoute.Spec.TLS == nil {
		return
	}

	secretName := ingressRoute.Spec.TLS.SecretName

	certName := GetCertName(ingressRoute)

	cert, err := clientSet.
		CertmanagerV1().
		Certificates(ingressRoute.Namespace).
		Get(ctx, certName, metav1.GetOptions{})
	if err != nil {
		var c = NewCert(
			WithName(certName),
			WithIssuerRef(issuer),
			WithSecretName(secretName),
		)
		yamlData, _ := json.Marshal(c)
		log.Infof("%+v\n", string(yamlData))
		//clientSet.CertmanagerV1().Certificates(ingressRoute.Namespace).Create(ctx, c, metav1.CreateOptions{})
	} else {
		fmt.Println(cert)
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
