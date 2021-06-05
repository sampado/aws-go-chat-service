package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sampado/aws-go-chat-service/chat"
	"github.com/sampado/aws-go-chat-service/external"
)

type Request events.APIGatewayWebsocketProxyRequest

// Handler is our lambda handler. It subscribes an user to the pool of connected users.

// ConnectHandler adapter used to pass extra parameters to the handler without changing the signature
func ConnectHandler(chat *chat.RoomSession) func(ctx context.Context, request Request) (events.APIGatewayProxyResponse, error) {
	return func(ctx context.Context, request Request) (events.APIGatewayProxyResponse, error) {
		// gets connection id
		connectionId := request.RequestContext.ConnectionID

		// saves connection id into a database
		err := chat.Connect(connectionId)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 500}, err
		}

		// return OK
		messageOK := fmt.Sprintf("User %s connected successfully!", connectionId)
		return external.NewApiGatewayResponseOK(messageOK), nil
	}
}

func main() {
	chat := &chat.RoomSession{
		Repository: external.NewAWSDynamoDBRepository(),
		Messenger:  external.NewAPIGatewayMessenger(),
	}
	lambda.Start(ConnectHandler(chat))
}
