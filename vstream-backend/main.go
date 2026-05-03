// vstream-backend/main.go
package main

import (
	"log/slog"
	"os"
	"strings"
	"vstream-backend/config"
	"vstream-backend/database"
	"vstream-backend/jobs"
	"vstream-backend/middleware"
	"vstream-backend/router"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	slog.SetDefault(logger)
}

// main.go — ganti corsMiddleware() menjadi:
func corsMiddleware() gin.HandlerFunc {
	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "http://localhost:5173" // fallback development
		slog.Warn("ALLOWED_ORIGIN tidak diset, fallback ke localhost:5173")
	}

	return func(c *gin.Context) {
		origin := allowedOrigin
		if strings.HasPrefix(c.Request.URL.Path, "/static/subtitles/") {
			origin = "*"
		}

		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
func main() {
	if err := godotenv.Load(); err != nil {
		slog.Warn("File .env tidak ditemukan, menggunakan environment OS")
	}

	config.Validate()
	database.InitDB()
	jobs.Start()

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(corsMiddleware())
	r.Use(middleware.BodyLimit(10 << 20))

	router.Setup(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	slog.Info("VSTREAM Engine started", "port", port)
	r.Run(":" + port)
}
