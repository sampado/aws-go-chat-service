package chat

import (
	"errors"
	"log"
)

type Session interface {
	Connect(ID string) error
	Disconnect(ID string) error
	Broadcast(senderID, msg string) error
}

type MessageData struct {
	Message      string `json:"message"`
	ConnectionID string `json:"connectionId"`
}

type Messenger interface {
	Send(msg MessageData, to string) error
}

type ConnectionRepository interface {
	Get(ID string) (*ConnectionItem, error)
	GetAll() ([]ConnectionItem, error)
	Save(*ConnectionItem) error
	Delete(ID string) error
}

type ConnectionItem struct {
	ID string `json:"connectionID"`
}

type RoomSession struct {
	Repository ConnectionRepository
	Messenger  Messenger
}

func (s *RoomSession) Connect(ID string) error {
	return s.Repository.Save(&ConnectionItem{ID})
}

var (
	ErrConnectionNotFound = errors.New("connection not found")
)

func (s *RoomSession) Disconnect(ID string) error {
	if err := s.Repository.Delete(ID); err != nil {
		if err == ErrConnectionNotFound {
			log.Printf("WARN: user %s not found", ID)
		} else {
			return err
		}
	}

	return nil
}

func (s *RoomSession) Broadcast(senderID, msg string) error {
	connections, err := s.Repository.GetAll()
	if err != nil {
		return err
	}

	for _, con := range connections {
		message := MessageData{
			Message:      msg,
			ConnectionID: senderID,
		}

		err := s.Messenger.Send(message, con.ID)
		if err != nil {
			log.Printf("WARN: failed to deliver the message to %s. %v", con.ID, err)
		}
	}

	return nil
}
