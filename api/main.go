package main

import (
	"context"
	"log"
	"url-shortener/api/handlers"
	"url-shortener/api/routes"
	"url-shortener/api/services"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// MongoDB setup
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://user:password@mongo:27017"))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(context.Background())

	// Redis setup
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "password",
		DB:       3,
	})

	// Initialize services
	urlService := services.NewURLService(mongoClient.Database("urlshortener"), redisClient, "http://localhost:8080")
	urlHandler := handlers.NewURLHandler(urlService)

	// Set up Gin router
	router := gin.Default()
	routes.InitializeRoutes(router, urlHandler)

	// Start server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
