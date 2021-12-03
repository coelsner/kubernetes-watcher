package teams

import (
	"bytes"
	"html/template"
	"log"
)

type Content struct {
	Color string
	Title string
	Text  string
}

var templates *template.Template

func init() {
	var err error
	if templates, err = template.ParseGlob("./templates/*.json"); err != nil {
		log.Fatalf("Error during parsing of templates: %v", err)
	}
}

func NewPodMessage(content Content) (bytes.Buffer, error) {
	var buf bytes.Buffer
	err := templates.ExecuteTemplate(&buf, "pod.json", content)
	return buf, err
}
