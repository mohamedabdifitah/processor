package template

import (
	"bytes"
	htmltemplate "html/template"
	"log"
	"os"
	texttemplate "text/template"
)

type Pattern struct {
	ExpireTime int
	Otp        int
	Unit       string
}
type Tag struct {
	pattern string
	replace string
}

func TemplateInjector() string {
	tmplt, err := htmltemplate.ParseFiles("html/otp.html")
	if err != nil {
		log.Fatal(err)
	}
	// pattern that injects the template into values
	event := Pattern{
		ExpireTime: 30,
		Otp:        55567,
		Unit:       "minutes",
	}
	writer := new(bytes.Buffer)
	// var writer []byte
	// injection template writer
	err = tmplt.Execute(writer, event)
	if err != nil {
		log.Fatal(err)
	}
	return writer.String() + " "
	// fmt.Println(writer)
}
func TextTemplateInjector() {
	// get text template
	tmpl, err := texttemplate.New("test").Parse("{{.Count}} items are made of {{.Material}}")
	if err != nil {

	}
	event := Pattern{
		ExpireTime: 30,
		Otp:        55567,
		Unit:       "minutes",
	}
	err = tmpl.Execute(os.Stdout, event)
	if err != nil {
		log.Fatal()
	}

}

// Text Template
func main() {
	TemplateInjector()
}
