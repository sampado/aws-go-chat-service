package external

import (
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/sampado/aws-go-chat-service/chat"
)

const (
	EnvAccessKeyID        = "ACCESS_KEY"
	EnvSecretAccessKey    = "SECRET_KEY"
	EnvAPIGatewayEndpoint = "API_GATEWAY_ENDPOINT"
	EnvRegion             = "REGION"
)

func newAPIGatewaySession() *apigatewaymanagementapi.ApiGatewayManagementApi {
	sess, _ := session.NewSession(&aws.Config{
		Region:      aws.String(os.Getenv(EnvRegion)),
		Credentials: credentials.NewStaticCredentials(os.Getenv(EnvAccessKeyID), os.Getenv(EnvSecretAccessKey), ""),
		Endpoint:    aws.String(os.Getenv(EnvAPIGatewayEndpoint)),
	})
	return apigatewaymanagementapi.New(sess)
}

type APIGatewayMessenger struct {
	session *apigatewaymanagementapi.ApiGatewayManagementApi
}

func NewAPIGatewayMessenger() *APIGatewayMessenger {
	return &APIGatewayMessenger{
		session: newAPIGatewaySession(),
	}
}

func (m *APIGatewayMessenger) Send(msg chat.MessageData, toID string) error {
	jsonData, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	input := &apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: aws.String(toID),
		Data:         jsonData,
	}

	_, errPost := m.session.PostToConnection(input)
	return errPost
}