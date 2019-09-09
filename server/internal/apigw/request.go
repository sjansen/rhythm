package apigw

import (
	"bytes"
	"context"
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

func EventToRequest(e *events.APIGatewayProxyRequest) (*http.Request, error) {
	u, err := url.Parse(e.Path)
	if err != nil {
		return nil, err
	}

	q := url.Values(e.MultiValueQueryStringParameters)
	u.RawQuery = q.Encode()

	u.Host = e.Headers["Host"]
	u.Scheme = e.Headers["X-Forwarded-Proto"]
	if u.Host != "" && u.Scheme != "" {
		port, ok := e.Headers["X-Forwarded-Port"]
		switch {
		case ok && u.Scheme == "http" && port != "80":
			u.Host += ":" + port
		case ok && u.Scheme == "https" && port != "443":
			u.Host += ":" + port
		}
	}

	body := []byte(e.Body)
	if e.IsBase64Encoded {
		body, err = base64.StdEncoding.DecodeString(e.Body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(
		strings.ToUpper(e.HTTPMethod),
		u.String(),
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, err
	}

	if e.MultiValueHeaders != nil {
		req.Header = http.Header(e.MultiValueHeaders)
	} else {
		for k, v := range e.Headers {
			req.Header.Set(k, v)
		}
	}

	return req, nil
}

func EventToRequestWithContext(ctx context.Context, e *events.APIGatewayProxyRequest) (*http.Request, error) {
	req, err := EventToRequest(e)
	return req.WithContext(ctx), err
}
