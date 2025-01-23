package services

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
	"url-shortener/internal/models"
	"url-shortener/internal/queries"

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

func (s *URLService) ShortenURL(original string, expiry int64) (models.URL, error) {
	if expiry == 0 {
		expiry = 60 // default expiry in a second
	}
	expiredAt := time.Duration(expiry) * time.Second

	shortened := generateShortenedURL()
	url := models.URL{
		Original:       original,
		BaseURL:        s.baseURL,
		Shortened:      shortened,
		ClickCount:     0,
		Expiry:         expiry,
		ExpiryDuration: int64(expiredAt),
		ExpiredAt:      time.Now().Add(expiredAt).Unix(),
		CreatedAt:      time.Now().Unix(),
	}

	b, err := json.Marshal(url)
	if err != nil {
		return models.URL{}, err
	}

	// Store in Redis
	key := fmt.Sprintf(`shorten:%s`, shortened)
	s.redisClient.Set(s.redisClient.Context(), key, string(b), expiredAt)

	return url, nil
}

func (s *URLService) FetchURL(shortened string) (models.URL, error) {
	// Fetch from Redis first
	key := fmt.Sprintf(`shorten:%s`, shortened)
	original, err := s.redisClient.Get(s.redisClient.Context(), key).Result()
	if err != nil {
		return models.URL{}, err
	}

	var url models.URL
	err = json.Unmarshal([]byte(original), &url)
	if err != nil {
		return models.URL{}, err
	}

	url.ClickCount += 1

	b, err := json.Marshal(url)
	if err != nil {
		return models.URL{}, err
	}

	s.redisClient.Set(s.redisClient.Context(), key, string(b), time.Duration(url.ExpiryDuration))

	return url, nil
}

func generateShortenedURL() string {
	b := make([]byte, 6) // 6 bytes will give us 8 characters in base64
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func (s *URLService) UpdateShortURL(shortened string, expiry int64) (*models.URL, error) {
	if expiry == 0 {
		expiry = 60 // default expiry in a second
	}
	expiredAt := time.Duration(expiry) * time.Second

	key := fmt.Sprintf(`shorten:%s`, shortened)
	original, err := s.redisClient.Get(s.redisClient.Context(), key).Result()
	if err != nil {
		return nil, err
	}

	var url models.URL
	err = json.Unmarshal([]byte(original), &url)
	if err != nil {
		return nil, err
	}

	url.Expiry = expiry
	url.ExpiryDuration = int64(expiredAt)
	url.ExpiredAt = time.Now().Add(expiredAt).Unix()

	b, err := json.Marshal(url)
	if err != nil {
		return nil, err
	}

	s.redisClient.Set(s.redisClient.Context(), key, string(b), expiredAt)

	return &url, nil
}
