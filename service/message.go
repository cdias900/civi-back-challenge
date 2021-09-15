package service

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"log"
	"time"

	"github.com/cdias900/civi-back-challenge/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

// Generate random uint64
func generateUint64() uint64 {
	buf := make([]byte, 8)
	rand.Read(buf)
	return binary.LittleEndian.Uint64(buf)
}

// Constants
const (
	MESSAGES_PER_PAGE = 20
)

// Send message
func (mS *messageService) Send(ctx context.Context, dbS DBService, msg model.MessageContent) (*model.Message, error) {
	m := model.Message{
		ID:        generateUint64(),
		Timestamp: primitive.Timestamp{T: uint32(time.Now().Unix())},
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
func (mS *messageService) Read(ctx context.Context, dbS DBService, page int64) ([]model.Message, error) {
	docs, err := dbS.ReadDocuments(ctx, "messages", nil, page*MESSAGES_PER_PAGE, page)
	if err != nil {
		log.Println("couldn't read message documents from database:", err)
		return nil, err
	}

	msgs := make([]model.Message, len(docs))
	for i, msg := range docs {
		msgs[i] = (*msg).(model.Message)
	}

	return msgs, nil
}
