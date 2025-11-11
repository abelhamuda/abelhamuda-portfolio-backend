package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"backend/models"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

// Article Handler
type ArticleHandler struct {
	db *sql.DB
}

func NewArticleHandler(db *sql.DB) *ArticleHandler {
	return &ArticleHandler{db: db}
}

// Get all articles
func (h *ArticleHandler) GetArticles(w http.ResponseWriter, r *http.Request) {
	query := `SELECT id, title, slug, content, category, tags, created_at, updated_at FROM articles ORDER BY created_at DESC`
	rows, err := h.db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var articles []models.Article
	for rows.Next() {
		var article models.Article
		err := rows.Scan(&article.ID, &article.Title, &article.Slug, &article.Content, &article.Category, &article.Tags, &article.CreatedAt, &article.UpdatedAt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		articles = append(articles, article)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(articles)
}

// Get article by slug
func (h *ArticleHandler) GetArticleBySlug(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	var article models.Article
	query := `SELECT id, title, slug, content, category, tags, created_at, updated_at FROM articles WHERE slug = ?`
	err := h.db.QueryRow(query, slug).Scan(
		&article.ID, &article.Title, &article.Slug, &article.Content, &article.Category, &article.Tags, &article.CreatedAt, &article.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Article not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(article)
}

// Get stats
func (h *ArticleHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	var totalArticles int
	err := h.db.QueryRow("SELECT COUNT(*) FROM articles").Scan(&totalArticles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var totalCategories int
	err = h.db.QueryRow("SELECT COUNT(DISTINCT category) FROM articles").Scan(&totalCategories)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stats := models.Stats{
		TotalArticles:   totalArticles,
		TotalCategories: totalCategories,
		TotalViews:      totalArticles * 150, // Simulated views
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// Login
func (h *ArticleHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginReq models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Hardcoded admin credentials for development
	if loginReq.Email == "admin@abelhamuda.com" && loginReq.Password == "admin123" {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"email": loginReq.Email,
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
		})

		tokenString, err := token.SignedString([]byte("your-secret-key"))
		if err != nil {
			http.Error(w, "Error generating token", http.StatusInternalServerError)
			return
		}

		response := models.LoginResponse{
			Token:   tokenString,
			Message: "Login successful",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	http.Error(w, "Invalid credentials", http.StatusUnauthorized)
}

// Create article
func (h *ArticleHandler) CreateArticle(w http.ResponseWriter, r *http.Request) {
	var article models.Article
	if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Generate slug from title
	article.Slug = strings.ToLower(strings.ReplaceAll(article.Title, " ", "-"))

	query := `INSERT INTO articles (title, slug, content, category, tags) VALUES (?, ?, ?, ?, ?)`
	result, err := h.db.Exec(query, article.Title, article.Slug, article.Content, article.Category, article.Tags)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	article.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(article)
}

// Update article
func (h *ArticleHandler) UpdateArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	var article models.Article
	if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	query := `UPDATE articles SET title=?, content=?, category=?, tags=?, updated_at=NOW() WHERE id=?`
	_, err = h.db.Exec(query, article.Title, article.Content, article.Category, article.Tags, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Article updated successfully"})
}

// Delete article
func (h *ArticleHandler) DeleteArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	_, err = h.db.Exec("DELETE FROM articles WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Article deleted successfully"})
}

// Projet Handler
type ProjectHandler struct {
	db *sql.DB
}

func NewProjectHandler(db *sql.DB) *ProjectHandler {
	return &ProjectHandler{db: db}
}

// Get all projects
func (h *ProjectHandler) GetProjects(w http.ResponseWriter, r *http.Request) {
	query := `SELECT id, title, description, category, technologies, thumbnail_url, github_url, live_url, created_at, updated_at FROM projects ORDER BY created_at DESC`
	rows, err := h.db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var project models.Project
		err := rows.Scan(
			&project.ID, &project.Title, &project.Description, &project.Category,
			&project.Technologies, &project.ThumbnailURL, &project.GithubURL,
			&project.LiveURL, &project.CreatedAt, &project.UpdatedAt,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		projects = append(projects, project)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

// Get project by ID
func (h *ProjectHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	var project models.Project
	query := `SELECT id, title, description, category, technologies, thumbnail_url, github_url, live_url, created_at, updated_at FROM projects WHERE id = ?`
	err = h.db.QueryRow(query, id).Scan(
		&project.ID, &project.Title, &project.Description, &project.Category,
		&project.Technologies, &project.ThumbnailURL, &project.GithubURL,
		&project.LiveURL, &project.CreatedAt, &project.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Project not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}

// Create project
func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var project models.Project
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO projects (title, description, category, technologies, thumbnail_url, github_url, live_url) VALUES (?, ?, ?, ?, ?, ?, ?)`
	result, err := h.db.Exec(query,
		project.Title, project.Description, project.Category,
		project.Technologies, project.ThumbnailURL, project.GithubURL, project.LiveURL,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	project.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(project)
}

// Update project
func (h *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	var project models.Project
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	query := `UPDATE projects SET title=?, description=?, category=?, technologies=?, thumbnail_url=?, github_url=?, live_url=?, updated_at=NOW() WHERE id=?`
	_, err = h.db.Exec(query,
		project.Title, project.Description, project.Category, project.Technologies,
		project.ThumbnailURL, project.GithubURL, project.LiveURL, id,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Project updated successfully"})
}

// Delete project
func (h *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	_, err = h.db.Exec("DELETE FROM projects WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Project deleted successfully"})
}
