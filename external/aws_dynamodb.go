package external

import (
	"errors"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/sampado/aws-go-chat-service/chat"
)

var (
	connectionsTable = os.Getenv("CHAT_CONNECTIONS_TABLE_NAME")
)

func NewDynamoDBSession() *dynamodb.DynamoDB {
	// staticCred:=credentials.NewStaticCredentials(os.Getenv(EnvAccessKeyID), os.Getenv(EnvSecretAccessKey), "")
	// sess, _ := session.NewSession(&aws.Config{
	// 	Region:      aws.String(os.Getenv(EnvRegion)),
	// 	Credentials: staticCred,
	// })
	mySession := session.Must(session.NewSession())
	return dynamodb.New(mySession)
}

type AWSDynamoDBRepository struct {
	db *dynamodb.DynamoDB
}

func NewAWSDynamoDBRepository() *AWSDynamoDBRepository {
	return &AWSDynamoDBRepository{db: NewDynamoDBSession()}
}

func getPrimaryKey(ID string) map[string]*dynamodb.AttributeValue {
	return map[string]*dynamodb.AttributeValue{
		"connectionID": {
			S: aws.String(ID),
		},
	}
}

func (r *AWSDynamoDBRepository) Get(ID string) (*chat.ConnectionItem, error) {
	input := &dynamodb.GetItemInput{
		Key:       getPrimaryKey(ID),
		TableName: aws.String(connectionsTable),
	}
	result, err := r.db.GetItem(input)
	if err != nil {
		return nil, err
	}

	connection := &chat.ConnectionItem{}
	dynamodbattribute.UnmarshalMap(result.Item, &connection)
	return connection, nil
}

func (r *AWSDynamoDBRepository) GetAll() ([]chat.ConnectionItem, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(connectionsTable),
	}
	result, err := r.db.Scan(input)
	if err != nil {
		return nil, err
	}

	items := result.Items
	connections := []chat.ConnectionItem{}
	if len(items) == 0 {
		return connections, nil
	}

	if err := dynamodbattribute.UnmarshalListOfMaps(result.Items, &connections); err != nil {
		return nil, err
	}

	return connections, nil
}

func (r *AWSDynamoDBRepository) Save(connection *chat.ConnectionItem) error {
	if connection == nil {
		return errors.New("can not store a nil connection")
	}

	item, err := dynamodbattribute.MarshalMap(connection)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(connectionsTable),
	}

	r.db.PutItem(input)

	if err != nil {
		log.Printf("got error calling PutItem: %v", err)
		return err
	}

	return nil
}

func (r *AWSDynamoDBRepository) Delete(ID string) error {
	input := &dynamodb.DeleteItemInput{
		Key:       getPrimaryKey(ID),
		TableName: aws.String(connectionsTable),
	}

	_, err := r.db.DeleteItem(input)
	return err
}
