package models

type Article struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	Content     string `json:"content"`
	ContentType string `json:"content_type"`
	Category    string `json:"category"`
	Tags        string `json:"tags"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type Stats struct {
	TotalArticles   int `json:"total_articles"`
	TotalCategories int `json:"total_categories"`
	TotalViews      int `json:"total_views"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}

type Project struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Category     string `json:"category"`
	Technologies string `json:"technologies"` // atau []string dengan custom marshalling
	ThumbnailURL string `json:"thumbnail_url"`
	GithubURL    string `json:"github_url"`
	LiveURL      string `json:"live_url"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type ImageUploadResponse struct {
	URL string `json:"url"`
}
