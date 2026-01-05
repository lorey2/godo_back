package main

import (
	"fmt"
	"goback/config"
	"goback/database"
	"goback/handlers"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}


	config.Load()

	// 1. Connect to DB
	database.Connect()


	r := gin.Default()

	/*
	TODO: check with this code

	corsConfig := cors.DefaultConfig()
	
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://192.168.50.227:3000" // Your React IP and default port
	}

	corsConfig.AllowOrigins = []string{frontendURL}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	*/

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	r.Use(cors.New(config))

	// Public Routes
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)

	// Protected Routes
	protected := r.Group("/")
	protected.Use(handlers.AuthMiddleware())
	{
		protected.GET("/tasks", handlers.GetTasks)
		protected.POST("/tasks", handlers.AddTask)
		protected.PUT("/tasks/:id", handlers.UpdateTask)
		protected.DELETE("/tasks/:id", handlers.DeleteTask)
		protected.PUT("/tasks/reorder", handlers.ReorderTasks)
	}

	fmt.Println("Server running on http://localhost:8080")
	r.Run(":8080")
}
