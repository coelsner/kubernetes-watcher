package hooks

type Webhook interface {
	Publish(title string, text string, a ...interface{}) error
}
