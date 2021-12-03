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
	var title, text string

	switch event.Type {
	case watch.Added:
		title, text = "Added", fmt.Sprintf("New State: %v", event.Object)
	case watch.Modified:
		title, text = "Modified", fmt.Sprintf("New State: %v", event.Object)
	case watch.Deleted:
		title, text = "Deleted", fmt.Sprintf("State before deletion: %v", event.Object)
	case watch.Error:
		title, text = "Error", fmt.Sprintf("Error: %v", event.Object)
	}

	pod, ok := event.Object.(*coreV1.Pod)
	if ok {
		switch pod.Status.Phase {
		case coreV1.PodFailed:
			message := fmt.Sprintf("POD '%v' has failed: %v\n", pod.Name, pod.Status.Message)
			text = fmt.Sprintf("%s (%s)", message, text)
		case coreV1.PodSucceeded:
			message := fmt.Sprintf("POD '%v' was successful\n", pod.Name)
			text = fmt.Sprintf("%s (%s)", message, text)
		}
	}

	return webhook.Publish(title, text)
}
