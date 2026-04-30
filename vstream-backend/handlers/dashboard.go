// vstream-backend/handlers/dashboard.go
package handlers

import (
	"net/http"
	"vstream-backend/database"
	"vstream-backend/models"

	"github.com/gin-gonic/gin"
)

func GetDashboardStats(c *gin.Context) {
	var totalMovies int64
	var totalSessions int64
	var recentMovies []models.Movie
	var activeSessions int64

	// 1. Hitung total film
	database.DB.Model(&models.Movie{}).Count(&totalMovies)

	// 2. Hitung total tayangan (dari semua history watch_sessions)
	database.DB.Model(&models.WatchSession{}).Count(&totalSessions)

	// 3. Hitung penonton aktif (sesi yang belum expired)
	// Query ini mengasumsikan kita punya kolom expires_at
	database.DB.Model(&models.WatchSession{}).Where("expires_at > datetime('now')").Count(&activeSessions)

	// 4. Ambil 5 film terbaru
	database.DB.Order("created_at desc").Limit(5).Find(&recentMovies)

	// Kirim response
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"total_movies":  totalMovies,
			"total_views":   totalSessions,
			"active_users":  activeSessions,
			"server_load":   "Normal", // Bisa di-kembangkan dengan package gopsutil nanti
			"recent_movies": recentMovies,
		},
	})
}
