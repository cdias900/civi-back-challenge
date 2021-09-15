package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Message content
type MessageContent struct {
	Subject string `json:"subject"`
	Detail  string `json:"detail"`
}

// Message
type Message struct {
	ID        uint64              `json:"id" bson:"_id"`
	Timestamp primitive.Timestamp `json:"timestamp"`
	Content   MessageContent      `json:"content"`
}
