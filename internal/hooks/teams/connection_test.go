package teams

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_Post(t *testing.T) {
	var (
		headers http.Header
		body    []byte
	)

	srv := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			headers = r.Header
			body, _ = io.ReadAll(r.Body)
		}),
	)
	defer srv.Close()

	client := NewClient(srv.URL)
	expected := []byte("Hello World")

	if err := client.Post(context.TODO(), expected); err != nil {
		t.Errorf("Post should not return error: %v", err)
	}

	if headers.Get("Content-Type") == "" {
		t.Errorf("Content-Type not set")
	}

	if bytes.Compare(body, expected) != 0 {
		t.Errorf("Bodies not equal (expected=%v, got=%v", expected, body)
	}
}
