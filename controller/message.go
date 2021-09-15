package controller

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/cdias900/civi-back-challenge/model"
	"github.com/cdias900/civi-back-challenge/service"
	"github.com/gin-gonic/gin"
)

// Message controller
type MessageController interface {
	Send(c *gin.Context) (*model.Message, error)
	Read(c *gin.Context) ([]model.Message, error)
}

// Message controller private struct
type messageController struct {
	mService  service.MessageService
	dbService service.DBService
}

// New message controller
func NewMessageController(dbService service.DBService, mService service.MessageService) messageController {
	return messageController{mService, dbService}
}

// Send message
func (m *messageController) Send(c *gin.Context) (*model.Message, error) {
	// Bind message content
	msg := model.MessageContent{}
	err := c.Bind(&msg)
	if err != nil {
		log.Println("couldn't get message content from body:", err)
		return nil, err
	} else {
		// Create context
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		// Write message
		return m.mService.Send(ctx, m.dbService, msg)
	}
}

// Read messages
func (m *messageController) Read(c *gin.Context) ([]model.Message, error) {
	// Bind message content
	page, err := strconv.ParseInt(c.Query("page"), 10, 64)
	if err != nil {
		page = 1
	}
	// Create context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// Read messages
	return m.mService.Read(ctx, m.dbService, page-1)
}
