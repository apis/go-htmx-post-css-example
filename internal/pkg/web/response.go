package web

import (
	"io"
	"net/http"
)

type Response struct {
	Status      int
	ContentType string
	Content     io.Reader
	Headers     Headers
}

func (response *Response) Write(responseWriter http.ResponseWriter) {
	if response != nil {
		if response.ContentType != "" {
			responseWriter.Header().Set("Content-Type", response.ContentType)
		}
		for k, v := range response.Headers {
			responseWriter.Header().Set(k, v)
		}
		responseWriter.WriteHeader(response.Status)
		_, err := io.Copy(responseWriter, response.Content)

		if err != nil {
			responseWriter.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		responseWriter.WriteHeader(http.StatusOK)
	}
}
