package main

import (
	"context"
	"fmt"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	typed "k8s.io/client-go/kubernetes/typed/core/v1"
)

func services(ctx context.Context, client typed.CoreV1Interface, namespace string) error {
	var api = client.Services(namespace)
	var resourceVersion, err = getResourceVersion(api.List(ctx, metaV1.ListOptions{}))
	if err != nil {
		return err
	}

	InfoLogger.Printf("Services ResourceVersion: %v\n", resourceVersion)
	watcher, err := api.Watch(ctx, metaV1.ListOptions{ /*ResourceVersion: resourceVersion*/ })
	if err != nil {
		return err
	}

	go watching("pods", ctx, watcher.ResultChan(), onServiceEvent)
	return nil
}

func onServiceEvent(event watch.Event) error {
	service, ok := event.Object.(*coreV1.Service)
	if !ok {
		return fmt.Errorf("Could not cast to service: %v\n", event)
	}

	EventsLogger.Printf("SVC '%v': %v\n", service.Name, service.Status)
	return nil
}
