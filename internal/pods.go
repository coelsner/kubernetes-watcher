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

func WatchPods(ctx context.Context, client typed.CoreV1Interface, namespace string, webhook teams.Webhook) error {
	var api = client.Pods(namespace)
	var resourceVersion, err = getResourceVersion(api.List(ctx, metaV1.ListOptions{}))
	if err != nil {
		return err
	}

	InfoLogger.Printf("Pods ResourceVersion: %v\n", resourceVersion)
	watcher, err := api.Watch(ctx, metaV1.ListOptions{ /*ResourceVersion: resourceVersion*/ })
	if err != nil {
		return err
	}

	go watching("pods", ctx, watcher.ResultChan(), onPodsEvent, webhook)
	return nil
}

func onPodsEvent(event watch.Event, webhook teams.Webhook) error {
	pod, ok := event.Object.(*coreV1.Pod)
	if !ok {
		return fmt.Errorf("Could not cast to Pod: %v\n", event)
	}

	switch event.Type {
	case watch.Added:
		webhook.Publish("POD %v was added\n", pod.Name)
	case watch.Modified:
		webhook.Publish("POD %v was modified\n", pod.Name)
	case watch.Deleted:
		webhook.Publish("POD %v was deleted\n", pod.Name)
	}

	switch pod.Status.Phase {
	case coreV1.PodFailed:
		EventsLogger.Printf("POD '%v' has failed: %v\n", pod.Name, pod.Status.Message)
	case coreV1.PodSucceeded:
		EventsLogger.Printf("POD '%v' was successful\n", pod.Name)
	}

	return nil
}
