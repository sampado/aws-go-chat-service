package external

import (
	"github.com/sampado/aws-go-chat-service/chat"
)

const (
	ConnectionsTable = "chat-connections"
)

type AWSDynamoDBConnections struct{}

func (db AWSDynamoDBConnections) Get(ID string) (*chat.ConnectionItem, error) {
	panic("not implemented") // TODO: Implement
}

func (db AWSDynamoDBConnections) GetAll() ([]chat.ConnectionItem, error) {
	panic("not implemented") // TODO: Implement
}

func (db AWSDynamoDBConnections) Save(_ *chat.ConnectionItem) error {
	panic("not implemented") // TODO: Implement
}

func (db AWSDynamoDBConnections) Delete(ID string) error {
	panic("not implemented") // TODO: Implement
}
