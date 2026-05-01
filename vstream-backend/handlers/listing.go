package handlers

import (
	"net/http"
	"strconv"
	"vstream-backend/database"
	"vstream-backend/models"

	"github.com/gin-gonic/gin"
)

// ── Public Listing ────────────────────────────────────────────

// GetMoviesPublicList mengembalikan semua movie (tanpa URL stream)
func GetMoviesPublicList(c *gin.Context) {
	getPublicList(c, "movie")
}

// GetSeriesPublicList mengembalikan semua series (tanpa URL stream)
func GetSeriesPublicList(c *gin.Context) {
	getPublicList(c, "series")
}

// GetAnimePublicList mengembalikan semua anime (tanpa URL stream)
func GetAnimePublicList(c *gin.Context) {
	getPublicList(c, "anime")
}

// getPublicList — helper generik, strip URL sebelum kirim ke publik
func getPublicList(c *gin.Context, contentType string) {
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 || limit > 100 {
		limit = 20
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	var items []models.Movie
	var total int64

	database.DB.Model(&models.Movie{}).Where("type = ?", contentType).Count(&total)
	database.DB.
		Where("type = ?", contentType).
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&items)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    sanitizeList(items),
		"meta": gin.H{
			"total":  total,
			"limit":  limit,
			"offset": offset,
		},
	})
}

// ── Featured ──────────────────────────────────────────────────

// GetFeatured mengembalikan 1 konten terbaru sebagai featured hero
func GetFeatured(c *gin.Context) {
	var item models.Movie
	if err := database.DB.Order("created_at desc").First(&item).Error; err != nil {
		// Tidak ada konten sama sekali — kembalikan 204
		c.JSON(http.StatusNoContent, gin.H{"success": true, "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": sanitize(item)})
}

// ── TOP 10 ────────────────────────────────────────────────────

// GetTop10 mengembalikan 10 konten terbaru lintas semua tipe
func GetTop10(c *gin.Context) {
	var items []models.Movie
	database.DB.
		Order("created_at desc").
		Limit(10).
		Find(&items)

	c.JSON(http.StatusOK, gin.H{"success": true, "data": sanitizeList(items)})
}

// ── Helpers ───────────────────────────────────────────────────

// PublicMovie adalah DTO publik — tanpa field URL stream
type PublicMovie struct {
	ID          uint   `json:"id"`
	TmdbID      string `json:"tmdb_id"`
	Type        string `json:"type"`
	HasEpisodes bool   `json:"has_episodes"`
	Title       string `json:"title"`
	Year        string `json:"year"`
	Duration    string `json:"duration"`
	Rating      string `json:"rating"`
	Genre       string `json:"genre"`
	Poster      string `json:"poster"`
	Backdrop    string `json:"backdrop"`
	Logos       string `json:"logos"`
	Overview    string `json:"overview"`
}

func sanitize(m models.Movie) PublicMovie {
	return PublicMovie{
		ID:          m.ID,
		TmdbID:      m.TmdbID,
		Type:        m.Type,
		HasEpisodes: m.HasEpisodes,
		Title:       m.Title,
		Year:        m.Year,
		Duration:    m.Duration,
		Rating:      m.Rating,
		Genre:       m.Genre,
		Poster:      m.Poster,
		Backdrop:    m.Backdrop,
		Logos:       m.Logos,
		Overview:    m.Overview,
	}
}

func sanitizeList(items []models.Movie) []PublicMovie {
	result := make([]PublicMovie, 0, len(items))
	for _, m := range items {
		result = append(result, sanitize(m))
	}
	return result
}
