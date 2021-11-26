package hooks

type Webhook interface {
	Publish(format string, a ...interface{})
}
