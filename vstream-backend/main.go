// vstream-backend/main.go
package main

import (
	"log/slog"
	"os"
	"vstream-backend/database"
	"vstream-backend/router"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	opts := &slog.HandlerOptions{Level: slog.LevelInfo}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	slog.SetDefault(logger)
}

func corsMiddleware() gin.HandlerFunc {
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
	r.Use(corsMiddleware())

	router.Setup(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	slog.Info("VSTREAM Engine started", "port", port)
	r.Run(":" + port)
}
