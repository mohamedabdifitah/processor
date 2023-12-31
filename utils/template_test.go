package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
