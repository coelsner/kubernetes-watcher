package internal

import (
	"context"
	"fmt"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	typed "k8s.io/client-go/kubernetes/typed/core/v1"
	"kubernetes-watcher/internal/hooks"
)

func WatchPods(ctx context.Context, client typed.CoreV1Interface, namespace string, webhook hooks.Webhook) error {
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

func onPodsEvent(event watch.Event, webhook hooks.Webhook) error {
	pod, ok := event.Object.(*coreV1.Pod)
	if !ok {
		return fmt.Errorf("Could not cast to Pod: %v\n", event)
	}

	var err error

	switch event.Type {
	case watch.Added:
		err = webhook.Publish("Pod added", "POD %v was added\n", pod.Name)
	case watch.Modified:
		err = webhook.Publish("Pod modified", "POD %v was modified\n", pod.Name)
	case watch.Deleted:
		err = webhook.Publish("Pod deleted", "POD %v was deleted\n", pod.Name)
	}

	if err != nil {
		ErrorLogger.Printf("Could not publish: %v\n", err)
	}

	switch pod.Status.Phase {
	case coreV1.PodFailed:
		EventsLogger.Printf("POD '%v' has failed: %v\n", pod.Name, pod.Status.Message)
	case coreV1.PodSucceeded:
		EventsLogger.Printf("POD '%v' was successful\n", pod.Name)
	}

	return nil
}
