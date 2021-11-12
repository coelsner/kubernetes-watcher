package internal

import (
	"context"
	"fmt"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	typed "k8s.io/client-go/kubernetes/typed/core/v1"
	"kubernetes-watcher/internal/teams"
)

func WatchEvents(ctx context.Context, client typed.CoreV1Interface, namespace string) error {
	var api = client.Events(namespace)
	var resourceVersion, err = getResourceVersion(api.List(ctx, metaV1.ListOptions{}))
	if err != nil {
		return err
	}

	InfoLogger.Printf("Events ResourceVersion: %v\n", resourceVersion)
	watcher, err := api.Watch(ctx, metaV1.ListOptions{ /*ResourceVersion: resourceVersion*/ })
	if err != nil {
		return err
	}

	go watching("events", ctx, watcher.ResultChan(), onEventEvent, nil)
	return nil
}

func onEventEvent(watcherEvent watch.Event, _ teams.Webhook) error {
	event, ok := watcherEvent.Object.(*coreV1.Event)
	if !ok {
		return fmt.Errorf("Could not cast to Event: %v\n", watcherEvent)
	}

	EventsLogger.Printf("EVENT %v - %v\n", event.Name, event.Message)
	return nil
}
