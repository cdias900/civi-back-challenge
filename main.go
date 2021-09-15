package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cdias900/civi-back-challenge/controller"
	"github.com/cdias900/civi-back-challenge/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Get environment variables
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("error loading .env file")
	}

	// Initialize services
	dbService := service.NewDBService()
	mService := service.NewMessageService()

	// Create context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to database
	err = dbService.Connect(ctx, os.Getenv("ATLAS_URI"))
	if err != nil {
		log.Fatalln("error setting up database:", err)
	}

	// Initialize controllers
	mController := controller.NewMessageController(dbService, mService)

	// Initialize engine
	r := gin.Default()

	// Send message
	r.POST("/send", func(c *gin.Context) {
		msg, err := mController.Send(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": msg})
		}
	})

	// Read messages
	r.GET("/read", func(c *gin.Context) {
		msgs, err := mController.Read(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
		} else {
			c.JSON(http.StatusOK, gin.H{"messages": msgs})
		}
	})

	// Get port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Run engine
	r.Run(":" + port)
}
