package teams

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
)

type Client struct {
	client  *http.Client
	webhook string
}

func NewClient(url string) *Client {
	return &Client{
		client:  http.DefaultClient,
		webhook: url,
	}
}

func (c *Client) Post(ctx context.Context, content []byte) error {
	req, err := http.NewRequestWithContext(ctx, "POST", c.webhook, bytes.NewReader(content))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}

	fmt.Printf("Resp: %v\n", resp.Status)

	return nil
}
