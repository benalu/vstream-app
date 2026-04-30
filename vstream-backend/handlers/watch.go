package handlers

import (
	"net/http"
	"vstream-backend/database"
	"vstream-backend/models"
	"vstream-backend/pkg/verr"

	"github.com/gin-gonic/gin"
)

// GetWatchData — public endpoint untuk halaman watch
// GET /api/watch/:type/:tmdb_id
func GetWatchData(c *gin.Context) {
	contentType := c.Param("type") // "movie" | "series" | "anime"
	tmdbID := c.Param("tmdb_id")

	// Validasi type
	validTypes := map[string]bool{"movie": true, "series": true, "anime": true}
	if !validTypes[contentType] {
		verr.Handle(c, verr.NewPublicError(http.StatusBadRequest, "Tipe konten tidak valid", "invalid type param", nil))
		return
	}

	var item models.Movie
	if err := database.DB.
		Where("tmdb_id = ? AND type = ?", tmdbID, contentType).
		First(&item).Error; err != nil {
		verr.Handle(c, verr.NewPublicError(http.StatusNotFound, "Konten tidak ditemukan", "watch lookup failed", err))
		return
	}

	// Untuk series dan anime dengan episode: sertakan seasons
	if item.Type == "series" || (item.Type == "anime" && item.HasEpisodes) {
		var seasons []models.Season
		database.DB.
			Where("movie_id = ?", item.ID).
			Order("season_num asc").
			Preload("Episodes", "1=1 ORDER BY ep_num asc").
			Find(&seasons)

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": gin.H{
				"info":    item,
				"seasons": seasons,
			},
		})
		return
	}

	// Movie / anime tanpa episode: langsung return item (termasuk URL)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"info":    item,
			"seasons": nil,
		},
	})
}
