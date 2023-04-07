package event

import "k8s.io/apimachinery/pkg/watch"

type Event struct {
	watch.Event
	ProviderName string
}
