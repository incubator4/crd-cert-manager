package aggregator

import (
	"auto-cert/pkg/config"
	"auto-cert/pkg/event"
	"auto-cert/pkg/provider"
	"context"
	"fmt"
	v1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"time"
)

var _ provider.Provider = (*ProviderAggregator)(nil)

type ProviderAggregator struct {
	providers                 []provider.Provider
	providersThrottleDuration time.Duration
}

func (p *ProviderAggregator) Name() string {
	return ""
}

func (p *ProviderAggregator) quietAddProvider(provider provider.Provider) {
	err := p.AddProvider(provider)
	if err != nil {
		log.WithError(err).Errorf("Error while initializing provider %T", provider)
	}
}

func (p *ProviderAggregator) AddProvider(provider provider.Provider) error {
	err := provider.Init()
	if err != nil {
		return err
	}

	switch provider.(type) {

	default:
		p.providers = append(p.providers, provider)
	}

	return nil
}

func NewAggregatorProvider(config config.Configurations) ProviderAggregator {
	p := ProviderAggregator{}

	if config.Traefik != nil {
		p.quietAddProvider(config.Traefik)
	}

	return p
}

func (p *ProviderAggregator) Provide(ctx context.Context, event event.Event) (*v1.Certificate, error) {

	prd, err := p.findProviderByName(event.ProviderName)

	if err != nil {
		log.Errorf(err.Error())
		return nil, nil
	}
	prd.Provide(ctx, event)

	return nil, nil
}

func (p *ProviderAggregator) Init() error {
	return nil
}

func (p *ProviderAggregator) AddEventHandler(config *rest.Config, handler cache.ResourceEventHandler, stopCh <-chan struct{}) {

}

func (p *ProviderAggregator) WatchAll(config *rest.Config, stopCh <-chan struct{}) (<-chan event.Event, error) {
	eventChan := make(chan event.Event, 1)

	for _, prd := range p.providers {
		log.Infof("init %s provider \n", prd.Name())
		eventHandler := event.NewTypedEventHandler(prd.Name(), eventChan)
		prd.AddEventHandler(config, eventHandler, stopCh)
	}

	return eventChan, nil
}

func (p *ProviderAggregator) findProviderByName(name string) (provider.Provider, error) {
	for _, prd := range p.providers {
		if prd.Name() == name {
			return prd, nil
		}
	}
	return nil, fmt.Errorf("not found provider with name %s", name)
}
