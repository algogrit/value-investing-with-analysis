package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"codermana.com/go/pkg/value_analysis/internal/nse"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context) (Response, error) {
	downloader := nse.NewDownloader("./statements")

	err := downloader.Nifty50List()
	if err != nil {
		return Response{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, err
	}

	errs := downloader.PopulateAllStatementsList() // Blocking for all goroutines

	if errs != nil {
		return Response{
			StatusCode: http.StatusInternalServerError,
			Body:       errs.Error(),
		}, errs
	}

	var buf bytes.Buffer

	body, err := json.Marshal(downloader.Scripts)
	if err != nil {
		return Response{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, err
	}
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":      "application/json",
			"X-CoderMana-Reply": "get-statements-handler",
			"X-CoderMana-Team":  "kratika",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
