package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type Response events.APIGatewayProxyResponse
type Request events.APIGatewayWebsocketProxyRequest

// Handler is our lambda handler. It subscribes an user to the pool of connected users.
func Handler(ctx context.Context, request Request) (Response, error) {
	// gets connection id
	connectionId := request.RequestContext.ConnectionID
	// saves connection id into a database

	// return OK
	var buf bytes.Buffer
	body, err := json.Marshal(map[string]interface{}{
		"message": fmt.Sprintf("User %s connected successfully!", connectionId),
	})
	if err != nil {
		return Response{StatusCode: 404}, err
	}
	json.HTMLEscape(&buf, body)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "connect-handler",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
