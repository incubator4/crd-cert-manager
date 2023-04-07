package traefik

import (
	"auto-cert/pkg/cert"
	"auto-cert/pkg/event"
	"auto-cert/pkg/provider"
	"context"
	v1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	log "github.com/sirupsen/logrus"
	traefikVersioned "github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/generated/clientset/versioned"
	"github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/generated/informers/externalversions"
	"github.com/traefik/traefik/v2/pkg/provider/kubernetes/crd/traefik/v1alpha1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"time"
)

var _ provider.Provider = (*Provider)(nil)

type Provider struct {
	Resync time.Duration
}

func (p *Provider) Name() string {
	return "traefik"
}

func (p *Provider) Provide(ctx context.Context, event event.Event) (*v1.Certificate, error) {

	ing := IngressRoute{IngressRoute: *event.Object.(*v1alpha1.IngressRoute)}

	certName := provider.CertName(&ing)

	issuer, err := provider.GetIssuer(&ing)
	if err != nil {
		return nil, err
	}

	secretName, err := ing.GetSecretName()
	if err != nil {
		return nil, err
	}

	hosts := ing.GetHosts()
	log.Infof("%s | %s | %+v | %s | %+v \n", ing.Name, certName, issuer, secretName, hosts)

	c := cert.NewCert(
		cert.WithName(certName),
		cert.WithIssuerRef(issuer),
		cert.WithSecretName(secretName),
		cert.WithHosts(hosts),
	)

	return c, nil
}

func (p *Provider) Init() error {
	return nil
}

func (p *Provider) AddEventHandler(config *rest.Config, handler cache.ResourceEventHandler, stopCh <-chan struct{}) {
	clientSet, err := traefikVersioned.NewForConfig(config)
	if err != nil {
		return
	}

	factory := externalversions.NewSharedInformerFactory(clientSet, p.Resync)
	factory.Traefik().V1alpha1().IngressRoutes().Informer().AddEventHandler(handler)
	factory.Start(stopCh)

	for t, ok := range factory.WaitForCacheSync(stopCh) {
		if !ok {
			log.Errorf("timed out waiting for controller caches to sync %s", t.String())
		}
	}
}
