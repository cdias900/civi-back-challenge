package service

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"log"
	"time"

	"github.com/cdias900/civi-back-challenge/model"
)

// Message service
type MessageService interface {
	Send(ctx context.Context, dbS DBService, msg model.MessageContent) (*model.Message, error)
	Read(ctx context.Context, dbS DBService, page int64) ([]model.Message, error)
}

// Message service private struct
type messageService struct{}

// New message service
func NewMessageService() MessageService {
	return &messageService{}
}

// Generate random uint32
func generateUint32() uint32 {
	buf := make([]byte, 4)
	rand.Read(buf)
	return binary.LittleEndian.Uint32(buf)
}

// Constants
const (
	MESSAGES_PER_PAGE = 20
)

// Send message
func (mS *messageService) Send(ctx context.Context, dbS DBService, msg model.MessageContent) (*model.Message, error) {
	m := model.Message{
		ID:        generateUint32(),
		Timestamp: uint32(time.Now().Unix()),
		Content:   msg,
	}
	err := dbS.CreateDocument(ctx, "messages", m)
	if err != nil {
		log.Println("couldn't create message document on database:", err)
		return nil, err
	} else {
		return &m, nil
	}
}

// Read messages
func (mS *messageService) Read(ctx context.Context, dbS DBService, page int64) (output []model.Message, err error) {
	cur, err := dbS.ReadDocuments(ctx, "messages", nil, MESSAGES_PER_PAGE, page*MESSAGES_PER_PAGE)
	if err != nil {
		log.Println("couldn't read message documents from database:", err)
		return nil, err
	}

	msgs := make([]model.Message, 0)
	for cur.Next(context.TODO()) {
		// Create a value into which the single document can be decoded
		msg := model.Message{}
		err := cur.Decode(&msg)
		if err != nil {
			log.Println("couldn't decode document:", err)
			return nil, err
		}
		// Append
		msgs = append(msgs, msg)
	}
	return msgs, nil
}
