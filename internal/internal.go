package internal

import (
	"context"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"kubernetes-watcher/internal/hooks"
	"log"
	"os"
)

var (
	InfoLogger   = log.New(os.Stdout, "INFO ", log.LstdFlags)
	ErrorLogger  = log.New(os.Stdout, "ERROR ", log.LstdFlags)
	EventsLogger = log.New(os.Stdout, "EVENT ", log.LstdFlags)
)

func getResourceVersion(meta metaV1.ListMetaAccessor, err error) (string, error) {
	if err == nil {
		return meta.GetListMeta().GetResourceVersion(), nil
	} else {
		return "", err
	}
}

func watching(label string, ctx context.Context, ch <-chan watch.Event, onEvent func(watch.Event, hooks.Webhook) error, webhook hooks.Webhook) {
	for {
		select {
		case event := <-ch:
			err := onEvent(event, webhook)
			if err != nil {
				ErrorLogger.Printf("Could not process event: %v", event)
			}
		case <-ctx.Done():
			InfoLogger.Printf("Closing %s watcher\n", label)
			return
		}
	}
}
