package teams

import (
	"log"
	"os"
)

type (
	Webhook interface {
		Publish(format string, a ...interface{})
	}

	simpleWebhook struct {
		url string
	}
)

var (
	logger = log.New(os.Stdout, "HOOK ", log.LstdFlags)
)

func (s *simpleWebhook) Publish(format string, a ...interface{}) {
	logger.Printf(format, a...)
}

func New(webhookUrl string) Webhook {
	return &simpleWebhook{url: webhookUrl}
}
