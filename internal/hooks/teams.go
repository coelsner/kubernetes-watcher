package hooks

type teamsHook struct {
	url string
}

func NewTeamsHook(url string) Webhook {
	return &teamsHook{url: url}
}

func (t *teamsHook) Publish(_ string, _ ...interface{}) {
	// TODO
}
