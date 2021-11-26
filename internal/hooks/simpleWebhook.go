package hooks

import (
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

func (s *stdOutHook) Publish(format string, a ...interface{}) {
	logger.Printf(format, a...)
}
