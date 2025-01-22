package services

import (
	"crypto/rand"
	"encoding/base64"
	"time"
	"url-shortener/api/models"
	"url-shortener/api/queries"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
)

type URLService struct {
	urlQueries  *queries.URLQueries
	redisClient *redis.Client
	baseURL     string
}

func NewURLService(db *mongo.Database, redisClient *redis.Client, baseURL string) *URLService {
	return &URLService{
		urlQueries:  queries.NewURLQueries(db),
		redisClient: redisClient,
		baseURL:     baseURL,
	}
}

func (s *URLService) ShortenURL(original string, expiry int64, expiredAt int64) (models.URL, error) {
	shortened := generateShortenedURL()
	url := models.URL{
		Original:   original,
		BaseURL:    s.baseURL,
		Shortened:  shortened,
		ClickCount: 0,
		Expiry:     expiry,
		ExpiredAt:  expiredAt,
		CreatedAt:  time.Now().Unix(),
	}

	if err := s.urlQueries.InsertURL(url); err != nil {
		return models.URL{}, err
	}

	// Store in Redis
	s.redisClient.HSet(s.redisClient.Context(), shortened, "original", original)
	s.redisClient.ExpireAt(s.redisClient.Context(), shortened, time.Unix(expiredAt, 0))

	return url, nil
}

func (s *URLService) FetchURL(shortened string) (models.URL, error) {
	// Fetch from Redis first
	original, err := s.redisClient.HGet(s.redisClient.Context(), shortened, "original").Result()
	if err != nil {
		return models.URL{}, err
	}

	// Fetch the URL model from MongoDB
	url, err := s.urlQueries.FindURL(shortened)
	if err != nil {
		return models.URL{}, err
	}

	// Update the original URL in the model
	url.Original = original

	return url, nil
}

func (s *URLService) IncrementClickCounter(shortened string) {
	s.urlQueries.UpdateClickCount(shortened)
}

func (s *URLService) GetSortedURLs(ascending bool, page int, limit int, isExpired *bool) (models.PaginatedURLsResponse, error) {
	urls, totalCount, err := s.urlQueries.GetSortedURLs(ascending, page, limit, isExpired)
	if err != nil {
		return models.PaginatedURLsResponse{}, err
	}

	totalPages := totalCount / limit

	response := models.PaginatedURLsResponse{
		URLs:        urls,
		TotalItems:  totalCount,
		PageSize:    limit,
		TotalPages:  totalPages,
		CurrentPage: page,
	}

	// Calculate next and previous page
	if totalCount > page*limit {
		nextPage := page + 1
		response.NextPage = &nextPage
	}
	if page > 1 && totalPages > 0 {
		prevPage := page - 1
		response.PrevPage = &prevPage
	}

	return response, nil
}

func generateShortenedURL() string {
	b := make([]byte, 6) // 6 bytes will give us 8 characters in base64
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}
