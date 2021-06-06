package external

import (
	"bytes"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/sampado/aws-go-chat-service/chat"
)

const (
	APIGatewayEndpoint = "wss://7lucxu8nrd.execute-api.us-east-1.amazonaws.com/dev"
)

func NewApiGatewayResponseOK(msg string) events.APIGatewayProxyResponse {
	var buf bytes.Buffer
	body, err := json.Marshal(map[string]interface{}{
		"message": msg,
	})
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}
	}

	json.HTMLEscape(&buf, body)
	return events.APIGatewayProxyResponse{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
}

func NewApiGatewayErrorResponse(code int, msg string) events.APIGatewayProxyResponse {
	var buf bytes.Buffer
	body, err := json.Marshal(map[string]interface{}{
		"message": msg,
	})
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}
	}

	json.HTMLEscape(&buf, body)
	return events.APIGatewayProxyResponse{
		StatusCode:      code,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
}

func newAPIGatewaySession() *apigatewaymanagementapi.ApiGatewayManagementApi {
	sess, _ := session.NewSession(&aws.Config{
		Endpoint: aws.String(APIGatewayEndpoint),
	})
	//mySession := session.Must(session.NewSession())
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
