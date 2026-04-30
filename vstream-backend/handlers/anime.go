package handlers

import (
	"net/http"
	"vstream-backend/database"
	"vstream-backend/models"
	"vstream-backend/pkg/verr"

	"github.com/gin-gonic/gin"
)

// ── Anime Handlers ────────────────────────────────────────────

func AddAnime(c *gin.Context)    { addContent(c, "anime") }
func GetAllAnime(c *gin.Context) { getAllContent(c, "anime") }
func UpdateAnime(c *gin.Context) { updateContent(c, "anime") }
func DeleteAnime(c *gin.Context) { deleteContent(c, "anime") }

func GetAnimePublic(c *gin.Context) {
	id := c.Param("id")
	var anime models.Movie
	if err := database.DB.Where("tmdb_id = ? AND type = 'anime'", id).First(&anime).Error; err != nil {
		verr.Handle(c, verr.NewPublicError(http.StatusNotFound, "Anime tidak ditemukan", "Anime lookup failed", err))
		return
	}

	// Kalau anime punya episode, sertakan seasons
	if anime.HasEpisodes {
		var seasons []models.Season
		database.DB.Where("movie_id = ?", anime.ID).
			Order("season_num asc").
			Preload("Episodes", "1=1 ORDER BY ep_num asc").
			Find(&seasons)

		c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{
			"info":    anime,
			"seasons": seasons,
		}})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{
		"info":    anime,
		"seasons": nil,
	}})
}

// ── Anime Season Routes ───────────────────────────────────────

func GetAnimeSeasons(c *gin.Context) { getSeasons(c) }
func AddAnimeSeason(c *gin.Context)  { addSeason(c) }
