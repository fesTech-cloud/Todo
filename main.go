package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/festech-cloud/todo/controller"
	"github.com/festech-cloud/todo/database"
	"github.com/festech-cloud/todo/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	// Initialize MongoDB client
	database.Client = database.DBInstance(os.Getenv("DBURL"))

	controller.Init(database.Client)

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := database.Client.Disconnect(ctx); err != nil {
			fmt.Println("Error disconnecting from MongoDB:", err)
		}
		fmt.Println("Mongodb disconnected")
	}()

	// Load environment variables from .env file
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	// Set the port for the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize Gin router
	router := gin.New()
	router.Use(gin.Logger())

	// Setup routes
	routes.TodoRoutes(router)

	// Start the server
	router.Run(":" + port)
}
