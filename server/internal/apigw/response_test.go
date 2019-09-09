package apigw_test

import (
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/require"

	"github.com/sjansen/rhythm/server/internal/apigw"
)

func TestResponseWriter(t *testing.T) {
	require := require.New(t)

	w := apigw.NewResponseWriter()
	h := w.Header()
	for k, vs := range map[string][]string{
		"Content-Type": {"text/text"},
		"Set-Cookie":   {"foo=bar", "baz=qux; Max-Age=42"},
	} {
		for _, v := range vs {
			h.Add(k, v)
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Spoon!"))

	expected := &events.APIGatewayProxyResponse{
		StatusCode: 200,
		MultiValueHeaders: map[string][]string{
			"Content-Type": {"text/text"},
			"Set-Cookie":   {"foo=bar", "baz=qux; Max-Age=42"},
		},
		IsBase64Encoded: false,
		Body:            "Spoon!",
	}

	actual, err := w.GetResponse()
	require.NoError(err)
	require.Equal(expected, actual)
}
