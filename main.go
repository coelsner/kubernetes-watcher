package main

import (
	"context"
	"flag"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"kubernetes-watcher/internal"
	"kubernetes-watcher/internal/hooks"
	"log"
	"os"
	"os/signal"
)

var (
	namespace = flag.String("namespace", "default", "namespace to be watched")
)

func main() {
	flag.Parse()

	log.Printf("Using Namespace: %v\n", *namespace)

	var webhook hooks.Webhook

	webhookUrl, isPresent := os.LookupEnv("WEBHOOK_URL")
	if isPresent && webhookUrl != "" {
		log.Printf("Using webhook: '%v'\n", webhookUrl)
		webhook = hooks.NewTeamsHook(webhookUrl)
	} else {
		log.Printf("No webhook url found. Using stdout")
		webhook = hooks.NewStdOutHook()
	}

	// AUTHENTICATE
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Panicf("Could not load config: %v\n", err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panicf("Could not instantiate config: %v\n", err.Error())
	}

	ctx, cancel := context.WithCancel(context.Background())

	err = clientset.CoreV1().RESTClient().Get().Do(ctx).Error()
	if err != nil {
		cancel()
		log.Panicf("Could not query api server: %v\n", err)
	}

	go func() {
		if err = internal.WatchEvents(ctx, clientset.CoreV1(), *namespace); err != nil {
			log.Printf("ERROR: Could not start watching events: %v\n", err)
		}
	}()

	go func() {
		if err = internal.WatchPods(ctx, clientset.CoreV1(), *namespace, webhook); err != nil {
			log.Printf("ERROR: Could not start watching pods: %v\n", err)
		}
	}()

	log.Println("Started watching kubernetes ... Cancel with CTRL+C")

	osCh := make(chan os.Signal, 1)
	signal.Notify(osCh, os.Interrupt)
	select {
	case <-osCh:
		cancel()
	case <-ctx.Done():
		signal.Stop(osCh)
	}

	log.Println("Stopped watching kubernetes")
}
