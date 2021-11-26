package hooks

import (
	"fmt"
	"log"
	"os"
)

type stdOutHook struct{}

var (
	logger = log.New(os.Stdout, "HOOK ", log.LstdFlags)
)

func NewStdOutHook() Webhook {
	return &stdOutHook{}
}

func (s *stdOutHook) Publish(title string, text string, a ...interface{}) error {
	logger.Printf("%s: %s", title, fmt.Sprintf(text, a...))
	return nil
}
