package routes

import (
	"url-shortener/internal/handlers"

	"github.com/gin-gonic/gin"
)

func InitializeRoutes(router *gin.Engine, urlHandler *handlers.URLHandler) {
	router.POST("/shorten", urlHandler.ShortenURL)
	router.GET("/shorten/:shortened", urlHandler.FetchURL)
	router.GET("/urls/sorted", urlHandler.GetSortedURLs)
	router.PUT("/shorten", urlHandler.UpdateShortURL)
}
