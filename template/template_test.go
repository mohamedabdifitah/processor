package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var templates Templates

func TestTemplateInjector(t *testing.T) {
	message := "hello world {{.Name}}"
	Aftermessage := "hello world Mohamed"
	templates := Templates{
		"Test": Template{
			Message: message,
		},
	}
	str, err := templates.TempelateInjector("Test", map[string]string{
		"Name": "Mohamed",
	})
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, Aftermessage, str)
}
func TestMain(m *testing.M) {
	// templates = NewTemplates()
	// templates.LoadTemplates("./text/text.json")
	m.Run()
}
func TestLookUp(t *testing.T) {
	message := "hello world {{.Name}}"
	templates := Templates{
		"Test": Template{
			Message: message,
		},
	}
	str := templates.Lookup("Test")
	assert.Equal(t, message, str)
}
