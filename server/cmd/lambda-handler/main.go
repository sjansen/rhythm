package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-xray-sdk-go/xray"

	"github.com/go-chi/chi"
	"github.com/kelseyhightower/envconfig"

	"github.com/sjansen/rhythm/server/internal/api"
	"github.com/sjansen/rhythm/server/internal/apigw"
)

type App struct {
	env    EnvConfig
	router *chi.Mux
	sess   *session.Session
}

type EnvConfig struct {
	LambdaEnv  string `envconfig:"AWS_EXECUTION_ENV"`
	SecretName string `split_words:"true"`
}

func main() {
	app := &App{}

	err := envconfig.Process("rhythm", &app.env)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	if app.env.LambdaEnv == "" {
		fmt.Fprintln(os.Stderr, "This executable is intended to run on AWS Lambda.")
		os.Exit(1)
	}

	xray.Configure(xray.Config{
		LogLevel: "info",
	})

	app.sess, err = session.NewSession()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	cfg := &api.Config{}
	err = app.readSecrets(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	app.router = chi.NewRouter()
	app.router.Mount("/api", api.New(cfg))
	lambda.Start(app.handleRequest)
}

func (app *App) handleRequest(ctx context.Context, req *events.APIGatewayProxyRequest) (
	*events.APIGatewayProxyResponse, error,
) {
	r, err := apigw.EventToRequestWithContext(ctx, req)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		resp := &events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Headers: map[string]string{
				"Cache-Control": "no-cache, no-store, must-revalidate",
			},
		}
		return resp, nil
	}

	w := apigw.NewResponseWriter()
	app.router.ServeHTTP(w, r)

	resp, err := w.GetResponse()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		resp := &events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Headers: map[string]string{
				"Cache-Control": "no-cache, no-store, must-revalidate",
			},
		}
		return resp, nil
	}

	return resp, nil
}

func (app *App) readSecrets(cfg *api.Config) error {
	svc := ssm.New(app.sess)
	resp, err := svc.GetParameters(&ssm.GetParametersInput{
		Names: []*string{
			aws.String(app.env.SecretName),
		},
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return err
	}

	for _, param := range resp.Parameters {
		switch *param.Name {
		case app.env.SecretName:
			cfg.Secret = *param.Value
		}
	}
	return nil
}
