package handlers

import (
	"net/http"
	"url-shortener/internal/models"
	"url-shortener/internal/services"

	"github.com/gin-gonic/gin"
)

type URLHandler struct {
	service *services.URLService
}

func NewURLHandler(service *services.URLService) *URLHandler {
	return &URLHandler{service: service}
}

func (h *URLHandler) ShortenURL(c *gin.Context) {
	var url models.URL
	if err := c.ShouldBindJSON(&url); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if url.Expiry < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "expiry cannot be negative"})
		return
	}

	result, err := h.service.ShortenURL(url.Original, url.Expiry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *URLHandler) FetchURL(c *gin.Context) {
	shortened := c.Param("shortened")
	target, err := h.service.FetchURL(shortened)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	c.JSON(http.StatusFound, target)
}

func (h *URLHandler) UpdateShortURL(c *gin.Context) {
	var url models.URL
	if err := c.ShouldBindJSON(&url); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if url.Expiry < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "expiry cannot be negative"})
		return
	}

	result, err := h.service.UpdateShortURL(url.Shortened, url.Expiry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
