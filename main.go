package main

import (
	"context"
	"flag"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"log"
	"os"
	"os/signal"
)

var (
	InfoLogger   = log.New(os.Stdout, "INFO ", log.LstdFlags)
	ErrorLogger  = log.New(os.Stdout, "ERROR ", log.LstdFlags)
	EventsLogger = log.New(os.Stdout, "EVENT ", log.LstdFlags)
)

var (
	namespace    = flag.String("namespace", "default", "namespace to be watched")
	withServices = flag.Bool("enable-svc", false, "enabling watching services")
)

func main() {
	flag.Parse()

	InfoLogger.Printf("Using Namespace: %v\n", *namespace)

	// AUTHENTICATE
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Panicf("Could not load config: %v\n", err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	ctx, cancel := context.WithCancel(context.Background())

	err = clientset.CoreV1().RESTClient().Get().Do(ctx).Error()
	if err != nil {
		cancel()
		ErrorLogger.Fatalf("Could not query api server: %v\n", err)
	}

	go func() {
		if err = events(ctx, clientset.CoreV1(), *namespace); err != nil {
			ErrorLogger.Printf("Could not start watching events: %v\n", err)
		}
	}()

	go func() {
		if err = pods(ctx, clientset.CoreV1(), *namespace); err != nil {
			ErrorLogger.Printf("Could not start watching pods: %v\n", err)
		}
	}()

	if *withServices {
		go func() {
			if err = services(ctx, clientset.CoreV1(), *namespace); err != nil {
				ErrorLogger.Printf("Could not start watching services: %v\n", err)
			}
		}()
	}

	InfoLogger.Println("Started watching kubernetes ... Cancel with CTRL+C")

	osCh := make(chan os.Signal, 1)
	signal.Notify(osCh, os.Interrupt)
	select {
	case <-osCh:
		cancel()
	case <-ctx.Done():
		signal.Stop(osCh)
	}

	InfoLogger.Println("Stopped watching kubernetes")
}

func getResourceVersion(meta metaV1.ListMetaAccessor, err error) (string, error) {
	if err == nil {
		return meta.GetListMeta().GetResourceVersion(), nil
	} else {
		return "", err
	}
}

func watching(label string, ctx context.Context, ch <-chan watch.Event, onEvent func(watch.Event) error) {
	for {
		select {
		case event := <-ch:
			err := onEvent(event)
			if err != nil {
				ErrorLogger.Printf("Could not process event: %v", event)
			}
		case <-ctx.Done():
			InfoLogger.Printf("Closing %s watcher\n", label)
			return
		}
	}
}
