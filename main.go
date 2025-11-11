package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/rs/cors"

	"backend/handlers"
	"backend/middleware"

	"github.com/gorilla/mux"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Database connection
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test database connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	fmt.Println("Connected to MySQL database!")

	// Initialize handlers
	articleHandler := handlers.NewArticleHandler(db)
	projectHandler := handlers.NewProjectHandler(db)

	// Router
	r := mux.NewRouter()

	// Serve static files (uploads)
	r.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	// Public routes
	r.HandleFunc("/api/articles", articleHandler.GetArticles).Methods("GET")
	r.HandleFunc("/api/articles/{slug}", articleHandler.GetArticleBySlug).Methods("GET")
	r.HandleFunc("/api/stats", articleHandler.GetStats).Methods("GET")
	r.HandleFunc("/api/login", articleHandler.Login).Methods("POST")
	r.HandleFunc("/api/projects", projectHandler.GetProjects).Methods("GET")
	r.HandleFunc("/api/projects/{id}", projectHandler.GetProject).Methods("GET")

	// Protected routes
	protected := r.PathPrefix("/api/admin").Subrouter()
	protected.Use(middleware.AuthMiddleware)
	protected.HandleFunc("/articles", articleHandler.CreateArticle).Methods("POST")
	protected.HandleFunc("/articles/{id}", articleHandler.UpdateArticle).Methods("PUT")
	protected.HandleFunc("/articles/{id}", articleHandler.DeleteArticle).Methods("DELETE")
	protected.HandleFunc("/upload", articleHandler.UploadImage).Methods("POST") // Image upload
	protected.HandleFunc("/projects", projectHandler.CreateProject).Methods("POST")
	protected.HandleFunc("/projects/{id}", projectHandler.UpdateProject).Methods("PUT")
	protected.HandleFunc("/projects/{id}", projectHandler.DeleteProject).Methods("DELETE")

	// Configure CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization", "X-Requested-With"},
		AllowCredentials: true,
		Debug:            true,
	})

	// Wrap the router with CORS middleware
	handler := c.Handler(r)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
