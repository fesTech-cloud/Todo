package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/festech-cloud/todo/controller"
	"github.com/festech-cloud/todo/database"
	"github.com/festech-cloud/todo/routes"
	"github.com/gin-contrib/cors"
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
	router.Use(cors.Default())

	// Custom CORS middleware configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://example.com"}, // Update with your allowed origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK,
			gin.H{"status": true, "message": "App is running"})
	})

	// Setup routes
	routes.TodoRoutes(router)

	// Start the server
	router.Run(":" + port)
}
