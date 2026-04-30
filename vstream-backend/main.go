// vstream-backend/main.go
package main

import (
	"log/slog"
	"os"
	"vstream-backend/database"
	"vstream-backend/handlers"
	"vstream-backend/middleware"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	slog.SetDefault(logger)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ganti * dengan origin frontend spesifik saat production
		// contoh: "http://localhost:5173" atau "https://yourdomain.com"
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
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

	database.InitDB()

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(CORSMiddleware())

	api := r.Group("/api")
	{
		// ── Auth (publik, tidak butuh token) ────────────────────
		auth := api.Group("/auth")
		{
			auth.POST("/login", handlers.Login)
			auth.POST("/logout", handlers.Logout)
			auth.GET("/me", middleware.RequireAdmin, handlers.Me)
		}

		// ── Admin (semua route dilindungi JWT middleware) ────────
		admin := api.Group("/admin")
		admin.Use(middleware.RequireAdmin)
		{
			admin.GET("/dashboard", handlers.GetDashboardStats)
			admin.GET("/tmdb/:id", handlers.PreviewTMDB)
			admin.POST("/movies", handlers.AddMovie)
			admin.GET("/movies", handlers.GetAllMovies)
			admin.PUT("/movies/:id", handlers.UpdateMovie)
			admin.DELETE("/movies/:id", handlers.DeleteMovie)
		}

		// ── Public (tanpa auth) ──────────────────────────────────
		api.GET("/movies/:id", handlers.GetMoviePublic)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	slog.Info("VSTREAM Engine started", "port", port)
	r.Run(":" + port)
}
