package models

type URL struct {
	Original       string `json:"original"`
	BaseURL        string `json:"base_url"`
	Shortened      string `json:"shortened"`
	ClickCount     int    `json:"click_count"`
	Expiry         int64  `json:"expiry"` // in seconds
	ExpiryDuration int64  `json:"expiry_duration"`
	ExpiredAt      int64  `json:"expired_at"` // Unix timestamp
	CreatedAt      int64  `json:"created_at"` // Unix timestamp
}
