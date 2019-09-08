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

	"github.com/kelseyhightower/envconfig"

	"github.com/sjansen/rhythm/server/internal/api"
)

type App struct {
	api  api.Config
	env  EnvConfig
	sess *session.Session
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

	err = app.readSecrets()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	lambda.Start(app.handleRequest)
}

func (app *App) handleRequest(ctx context.Context, req *events.APIGatewayProxyRequest) (
	*events.APIGatewayProxyResponse, error,
) {
	resp := &events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Cache-Control": "no-cache, no-store, must-revalidate",
			"Content-Type":  "text/text; charset=utf-8",
		},
		Body: app.api.Secret,
	}
	return resp, nil
}

func (app *App) readSecrets() error {
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
			app.api.Secret = *param.Value
		}
	}
	return nil
}
