package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sampado/aws-go-chat-service/chat"
	"github.com/sampado/aws-go-chat-service/external"
)

type Request events.APIGatewayWebsocketProxyRequest

// RequestPayload represents the request body sent by the Socket
type RequestPayload struct {
	Message string `json:"message"`
}

func OnMessageHandler(chat *chat.RoomSession) func(ctx context.Context, request Request) (events.APIGatewayProxyResponse, error) {
	return func(ctx context.Context, request Request) (events.APIGatewayProxyResponse, error) {
		// gets connection id
		connectionId := request.RequestContext.ConnectionID
		// Parse the request
		var requestPayload RequestPayload
		if err := json.Unmarshal([]byte(request.Body), &requestPayload); err != nil {
			return external.NewApiGatewayErrorResponse(400, "bad request"), err
		}

		// saves connection id into a database
		err := chat.Broadcast(connectionId, requestPayload.Message)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 500}, err
		}

		// return OK
		return external.NewApiGatewayResponseOK("Message function executed successfully!"), nil
	}
}

func main() {
	chat := &chat.RoomSession{
		Repository: external.NewAWSDynamoDBRepository(),
		Messenger:  external.NewAPIGatewayMessenger(),
	}

	lambda.Start(OnMessageHandler(chat))
}
