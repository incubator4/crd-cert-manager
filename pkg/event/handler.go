package event

import (
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	"time"
)

type ResourceEventHandlerFunc func(kind string, ev chan<- Event) cache.ResourceEventHandler

var (
	pool = make([]Event, 0)
	t    = time.NewTimer(time.Second * 5)
)

func eventHandlerFunc(events chan<- Event, provider string, obj interface{}) {
	runtimeObj, ok := obj.(runtime.Object)
	if !ok {
		log.Errorf("Error converting object to runtime.Object: %v", obj)
		return
	}
	e := Event{
		Event: watch.Event{
			Type:   watch.Added,
			Object: runtimeObj,
		},
		ProviderName: provider,
	}
	select {
	case events <- e:
	case <-time.After(time.Second * 5):
		// Timeout occurred, event not sent
		log.Printf("Error sending event to channel: timeout occurred")
	default:
		log.Debug("Error sending event to channel: channel is full")
		pool = append(pool, e)
		select {
		case events <- pool[0]:
			pool = pool[1:]
		case <-time.After(time.Second):
			// Timeout occurred, cannot send integer
			log.Errorf("Timeout occurred, cannot send event")
		case <-t.C:
			// Channel is blocked, cannot send integer
			log.Errorf("Channel is blocked, cannot send event")
		}
	}
}

func objChanged(oldObj, newObj interface{}) bool {
	if oldObj == nil || newObj == nil {
		return true
	}

	if oldObj.(metav1.Object).GetResourceVersion() == newObj.(metav1.Object).GetResourceVersion() {
		return false
	}

	return true
}

func NewTypedEventHandler(kind string, ev chan<- Event) cache.ResourceEventHandler {
	return cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			eventHandlerFunc(ev, kind, obj)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			if objChanged(oldObj, newObj) {
				eventHandlerFunc(ev, kind, newObj)
			}
		},
		DeleteFunc: func(obj interface{}) {
			eventHandlerFunc(ev, kind, obj)
		},
	}
}
