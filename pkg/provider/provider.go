package provider

import (
	"auto-cert/pkg/event"
	"context"
	v1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

type Provider interface {
	Provide(ctx context.Context, event event.Event) (*v1.Certificate, error)
	AddEventHandler(config *rest.Config, handler cache.ResourceEventHandler, stopCh <-chan struct{})
	Init() error
	Name() string
}
