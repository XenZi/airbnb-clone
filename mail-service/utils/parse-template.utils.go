package utils

import (
	"bytes"
	"html/template"
	"log"
)

func ParseTemplate(file string, data interface{}) (string, error) {
	tmpl, errParseFiles := template.ParseFiles(file)
	if errParseFiles != nil {
		log.Println(errParseFiles)
		return "", errParseFiles
	}
	buffer := new(bytes.Buffer)
	if errExecute := tmpl.Execute(buffer, data); errExecute != nil {
		log.Println(errExecute)
		return "", errExecute
	}
	return buffer.String(), nil
}
