package hooks

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

type teamsHook struct {
	url string
}

func NewTeamsHook(url string) Webhook {
	return &teamsHook{url: url}
}

func (t *teamsHook) Publish(title string, text string, a ...interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tpl, err := template.New("message").Parse(`{
  "@context": "https://schema.org/extensions",
  "@type": "MessageCard",
  "themeColor": "0072C6",
  "title": "{{ .Title }}",
  "text": "{{ .Text }}",
}`)
	content := struct {
		Title string
		Text  string
	}{title, fmt.Sprintf(text, a)}

	var buf bytes.Buffer
	_ = tpl.Execute(&buf, content)

	req, err := http.NewRequestWithContext(ctx, "POST", t.url, bytes.NewReader(buf.Bytes()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	fmt.Printf("Resp: %v\n", resp.Status)
	return nil
}
