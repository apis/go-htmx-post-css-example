package web

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
)

func RenderResponse(status int, templates *template.Template, templateName string, data any, headers Headers) *Response {
	var buffer bytes.Buffer
	if err := templates.ExecuteTemplate(&buffer, templateName, data); err != nil {
		log.Println(err)
		return GetEmptyResponse(http.StatusInternalServerError)
	}

	return &Response{
		Status:      status,
		ContentType: "text/html",
		Content:     buffer.Bytes(),
		Headers:     headers,
	}
}

func GetEmptyResponse(status int) *Response {
	return GetResponse(status, []byte(""), nil)
}

func GetResponse(status int, content []byte, headers Headers) *Response {
	return &Response{
		Status:  status,
		Content: content,
		Headers: headers,
	}
}
