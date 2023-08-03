package template

import (
	"bytes"
	htmltemplate "html/template"
	"log"
	"os"
	texttemplate "text/template"
)

type Pattern struct {
	ExpireTime string
	Otp        string
	Unit       string
}
type Template struct {
	Definition interface{}
	path       string
	Name       string
}
type Definition struct {
	Key   string
	Value interface{}
}

func (p *Template) String() {
}
func TemplateInjector(definition map[string]string, path string) string {
	tmplt, err := htmltemplate.ParseFiles(path)
	if err != nil {
		log.Fatal(err)
	}
	writer := new(bytes.Buffer)
	err = tmplt.Execute(writer, definition)
	if err != nil {
		log.Fatal(err)
	}
	return writer.String()
}
func TextTemplateInjector() {
	// get text template
	tmpl, err := texttemplate.New("test").Parse("{{.Count}} items are made of {{.Material}}")
	if err != nil {

	}
	event := Pattern{
		ExpireTime: "30",
		Otp:        "55567",
		Unit:       "minutes",
	}
	err = tmpl.Execute(os.Stdout, event)
	if err != nil {
		log.Fatal()
	}

}
func CUstomTemplateInjector() {
	_, err := os.ReadFile("./template/html/reset_password.html")
	if err != nil {
		log.Fatal(err)
	}
	// regexp.Compile()
}

// Text Template
func main() {
	CUstomTemplateInjector()
	// TemplateInjector()
}
