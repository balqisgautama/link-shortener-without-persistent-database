package handlers

import (
	"net/http"
	"strconv"
	"time"
	"url-shortener/api/models"
	"url-shortener/api/services"

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

	expiredAt := time.Now().Add(time.Duration(url.Expiry) * time.Second).Unix()
	result, err := h.service.ShortenURL(url.Original, url.Expiry, expiredAt)
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

	h.service.IncrementClickCounter(shortened)
	target.ClickCount++
	c.JSON(http.StatusFound, target)
}

func (h *URLHandler) GetSortedURLs(c *gin.Context) {
	ascending := c.Query("order") == "asc"
	page, err := strconv.Atoi(c.Query("page"))
	var isExpired *bool
	if expiredQuery := c.Query("is_expired"); expiredQuery != "" {
		isExpiredValue, err := strconv.ParseBool(expiredQuery)
		if err == nil {
			isExpired = &isExpiredValue
		}
	}
	if err != nil || page < 1 {
		page = 1 // Default to page 1 if invalid
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil || limit < 1 {
		limit = 10 // Default limit
	}

	response, err := h.service.GetSortedURLs(ascending, page, limit, isExpired)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
