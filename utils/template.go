package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

type Templates map[string]Template
type Template struct {
	Message string `json:"message"`
	Mime    string `json:"mime"`
	Path    string `json:"path"`
}

func (templates Templates) LoadTemplates(path string, base string) {
	fileContent, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer fileContent.Close()
	byteResult, err := ioutil.ReadAll(fileContent)
	if err != nil {
		fmt.Println(err)
	}
	var res Templates
	err = json.Unmarshal([]byte(byteResult), &res)
	if err != nil {
		log.Fatal(err)
	}
	for key, value := range res {
		if value.Mime == "text/html" {
			path := filepath.Join(base, value.Path)
			fileContent, err := os.Open(path)
			if err != nil {
				log.Fatal(err)
			}
			defer fileContent.Close()
			writer := new(bytes.Buffer)
			byteResult, err := ioutil.ReadAll(fileContent)
			if err != nil {
				fmt.Println(err)
			}
			writer.Write(byteResult)
			value.Message = writer.String()
			templates[key] = value
		}
		templates[key] = value
	}
}
func (templates Templates) TempelateInjector(key string, definition map[string]string) (string, error) {
	temp := templates.Lookup(key)
	tmpl, err := template.New(key).Parse(temp)
	if err != nil {
		return "", err
	}
	writer := new(bytes.Buffer)
	err = tmpl.Execute(writer, definition)
	if err != nil {
		log.Fatal(err)
	}
	return writer.String(), nil
}
func (templates Templates) Lookup(title string) string {
	return templates[title].Message
}
func NewTemplates() Templates {
	return Templates{}
}
func CurrentTemplates() Templates {
	templates := NewTemplates()
	return templates
}
