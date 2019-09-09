package apigw

import (
	"bytes"
	"encoding/base64"
	"errors"
	"net/http"
	"unicode/utf8"

	"github.com/aws/aws-lambda-go/events"
)

var _ http.ResponseWriter = &ResponseWriter{}

type ResponseWriter struct {
	headers http.Header
	body    bytes.Buffer
	status  int
}

func NewResponseWriter() *ResponseWriter {
	return &ResponseWriter{
		headers: make(http.Header),
	}
}

func (w *ResponseWriter) GetResponse() (*events.APIGatewayProxyResponse, error) {
	if w.status == 0 {
		return &events.APIGatewayProxyResponse{}, errors.New("status code not set")
	}

	body := w.body.Bytes()
	if w.headers.Get("Content-Type") == "" {
		w.headers.Add("Content-Type", http.DetectContentType(body))
	}

	resp := &events.APIGatewayProxyResponse{
		StatusCode:        w.status,
		MultiValueHeaders: w.headers,
	}
	if utf8.Valid(body) {
		resp.Body = string(body)
	} else {
		resp.IsBase64Encoded = true
		resp.Body = base64.StdEncoding.EncodeToString(body)
	}

	return resp, nil
}

func (w *ResponseWriter) Header() http.Header {
	return w.headers
}

func (w *ResponseWriter) Write(body []byte) (int, error) {
	if w.status == 0 {
		w.WriteHeader(http.StatusOK)
	}
	return w.body.Write(body)
}

func (w *ResponseWriter) WriteHeader(status int) {
	w.status = status
}
