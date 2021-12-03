package hooks

import (
	"context"
	"fmt"
	"kubernetes-watcher/internal/hooks/teams"
	"time"
)

type teamsHook struct {
	client *teams.Client
}

func NewTeamsHook(url string) Webhook {
	return &teamsHook{client: teams.NewClient(url)}
}

func (t *teamsHook) Publish(title string, text string, a ...interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	content := teams.Content{Color: "0072C6", Title: title, Text: fmt.Sprintf(text, a...)}

	buf, err := teams.NewPodMessage(content)
	if err != nil {
		return err
	}

	return t.client.Post(ctx, buf.Bytes())
}
