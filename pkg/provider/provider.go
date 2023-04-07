package provider

import (
	"auto-cert/pkg/event"
	"context"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

type Provider interface {
	Provide(ctx context.Context, event event.Event) error
	AddEventHandler(config *rest.Config, handler cache.ResourceEventHandler, stopCh <-chan struct{})
	Init() error
	Name() string
}
