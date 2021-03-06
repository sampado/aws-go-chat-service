package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/sampado/aws-go-chat-service/chat"
	"github.com/sampado/aws-go-chat-service/external"
)

type Request events.APIGatewayWebsocketProxyRequest

func DisconnectHandler(chat *chat.RoomSession) func(ctx context.Context, request Request) (events.APIGatewayProxyResponse, error) {
	return func(ctx context.Context, request Request) (events.APIGatewayProxyResponse, error) {
		// gets connection id
		connectionId := request.RequestContext.ConnectionID

		// saves connection id into a database
		err := chat.Disconnect(connectionId)
		if err != nil {
			return events.APIGatewayProxyResponse{StatusCode: 500}, err
		}

		// return OK
		return external.NewApiGatewayResponseOK("Disconnect function executed successfully!"), nil
	}
}

func main() {
	chat := &chat.RoomSession{
		Repository: external.NewAWSDynamoDBRepository(),
		Messenger:  external.NewAPIGatewayMessenger(),
	}

	lambda.Start(DisconnectHandler(chat))
}
