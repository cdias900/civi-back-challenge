package model

// Message content
type MessageContent struct {
	Subject string `json:"subject"`
	Detail  string `json:"detail"`
}

// Message
type Message struct {
	ID        uint32         `json:"id" bson:"_id"`
	Timestamp uint32         `json:"timestamp"`
	Content   MessageContent `json:"content"`
}
