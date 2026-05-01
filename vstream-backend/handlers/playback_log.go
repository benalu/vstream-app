package handlers

import (
	"net/http"
	"time"
	"vstream-backend/database"
	"vstream-backend/models"
	"vstream-backend/pkg/verr"

	"github.com/gin-gonic/gin"
)

// ReportPlaybackError — public, dipanggil frontend saat video gagal
// POST /api/logs/playback-error
func ReportPlaybackError(c *gin.Context) {
	var input struct {
		MovieID   uint   `json:"movie_id"   binding:"required"`
		Server    string `json:"server"     binding:"required"`
		ErrorType string `json:"error_type"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		verr.Handle(c, verr.NewPublicError(http.StatusBadRequest, "Input tidak valid", "playback log bind failed", err))
		return
	}

	var movie models.Movie
	if err := database.DB.First(&movie, input.MovieID).Error; err != nil {
		// Movie tidak ditemukan — tetap OK, jangan ganggu user
		c.JSON(http.StatusOK, gin.H{"success": true})
		return
	}

	serverURL := map[string]string{
		"url1": movie.URL1,
		"url2": movie.URL2,
		"url3": movie.URL3,
	}[input.Server]

	errorType := input.ErrorType
	if errorType == "" {
		errorType = "load_error"
	}

	entry := models.PlaybackLog{
		MovieID:    movie.ID,
		MovieTitle: movie.Title,
		MovieType:  movie.Type,
		TmdbID:     movie.TmdbID,
		Server:     input.Server,
		URL:        serverURL,
		ErrorType:  errorType,
		UserAgent:  c.GetHeader("User-Agent"),
		ReportedAt: time.Now(),
	}

	database.DB.Create(&entry)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

// GetPlaybackLogs — admin
// GET /api/admin/playback-logs
func GetPlaybackLogs(c *gin.Context) {
	var logs []models.PlaybackLog
	database.DB.Order("reported_at desc").Limit(100).Find(&logs)

	type Summary struct {
		MovieID    uint   `json:"movie_id"`
		MovieTitle string `json:"movie_title"`
		MovieType  string `json:"movie_type"`
		TmdbID     string `json:"tmdb_id"`
		ErrorCount int    `json:"error_count"`
	}

	var summary []Summary
	database.DB.Model(&models.PlaybackLog{}).
		Select("movie_id, movie_title, movie_type, tmdb_id, count(*) as error_count").
		Group("movie_id").
		Order("error_count desc").
		Limit(10).
		Scan(&summary)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"logs":    logs,
			"summary": summary,
		},
	})
}

// ClearPlaybackLogs — admin
// DELETE /api/admin/playback-logs
func ClearPlaybackLogs(c *gin.Context) {
	database.DB.Where("1 = 1").Delete(&models.PlaybackLog{})
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Log berhasil dihapus"})
}
