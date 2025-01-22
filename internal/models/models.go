package models

type URL struct {
	Original   string `json:"original"`
	BaseURL    string `json:"base_url"`
	Shortened  string `json:"shortened"`
	ClickCount int    `json:"click_count"`
	Expiry     int64  `json:"expiry"`     // in seconds
	ExpiredAt  int64  `json:"expired_at"` // Unix timestamp
	CreatedAt  int64  `json:"created_at"` // Unix timestamp
}

type PaginatedURLsResponse struct {
	URLs        []URL `json:"urls"`
	TotalItems  int   `json:"total_items"`
	TotalPages  int   `json:"total_pages"`
	PageSize    int   `json:"page_size"`
	CurrentPage int   `json:"current_page"`
	NextPage    *int  `json:"next_page,omitempty"`
	PrevPage    *int  `json:"prev_page,omitempty"`
}
