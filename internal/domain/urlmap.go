package domain

type URLMap struct {
	ID          string `json:"id"`
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
	CreatedAt   int64  `json:"created_at"`
	Clicks      int64  `json:"clicks"`
}
