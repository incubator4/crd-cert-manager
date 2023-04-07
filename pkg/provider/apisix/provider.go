package apisix

import (
	"auto-cert/pkg/event"
	"auto-cert/pkg/provider"
	"context"
	"github.com/apache/apisix-ingress-controller/pkg/kube/apisix/apis/config/v2beta3"
	apisixVersioned "github.com/apache/apisix-ingress-controller/pkg/kube/apisix/client/clientset/versioned"
	"github.com/apache/apisix-ingress-controller/pkg/kube/apisix/client/informers/externalversions"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"time"
)

var _ provider.Provider = (*Provider)(nil)

type Provider struct {
	time.Duration `json:",inline"`
}

func (p Provider) Provide(ctx context.Context, event event.Event) error {

	ing := ApisixTls{ApisixTls: *event.Object.(*v2beta3.ApisixTls)}
	issuer, err := provider.GetIssuer(&ing)
	certName := provider.CertName(&ing)

	if err == nil {
		log.Infof("%s | %s | %+v | %s | %+v \n", ing.Name, certName, issuer, ing.GetSecretName(), ing.GetHosts())
	}

	return nil
}

func (p Provider) AddEventHandler(config *rest.Config, handler cache.ResourceEventHandler, stopCh <-chan struct{}) {
	clientSet, err := apisixVersioned.NewForConfig(config)
	if err != nil {
		return
	}

	factory := externalversions.NewSharedInformerFactory(clientSet, p.Duration)
	factory.Apisix().V2beta3().ApisixTlses().Informer().AddEventHandler(handler)
	factory.Start(stopCh)

	for t, ok := range factory.WaitForCacheSync(stopCh) {
		if !ok {
			log.Errorf("timed out waiting for controller caches to sync %s", t.String())
		}
	}
}

func (p Provider) Init() error {
	return nil
}

func (p Provider) Name() string {
	return "apisix"
}
