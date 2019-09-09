package apigw_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"

	"github.com/sjansen/rhythm/server/internal/apigw"
)

func TestEvenToRequest(t *testing.T) {
	require := require.New(t)

	e := &events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		Headers: map[string]string{
			"Host":              "example.com",
			"X-Forwarded-Port":  "443",
			"X-Forwarded-Proto": "https",
		},
		MultiValueHeaders: map[string][]string{
			"Host":              {"example.com"},
			"X-Forwarded-Port":  {"443"},
			"X-Forwarded-Proto": {"https"},
		},
		Path: "/",
		QueryStringParameters: map[string]string{
			"foo": "bar",
		},
		MultiValueQueryStringParameters: map[string][]string{
			"foo": {"bar", "baz"},
		},
		IsBase64Encoded: false,
	}

	actual, err := apigw.EventToRequest(e)

	u, _ := url.Parse("https://example.com/?foo=bar&foo=baz")
	expected := &http.Request{
		Method: "GET",
		URL:    u,
		Header: map[string][]string{
			"Host":              {"example.com"},
			"X-Forwarded-Port":  {"443"},
			"X-Forwarded-Proto": {"https"},
		},
		Host: "example.com",
		// TODO RemoteAddr
		// TODO RequestURI
	}

	require.NoError(err)
	require.Equal(expected.Header, actual.Header)
	require.Equal(expected.Host, actual.Host)
	require.Equal(expected.Method, actual.Method)
	require.Equal(expected.URL, actual.URL)
}
