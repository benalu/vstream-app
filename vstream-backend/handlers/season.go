package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"vstream-backend/database"
	"vstream-backend/models"
	"vstream-backend/pkg/verr"

	"github.com/gin-gonic/gin"
)

// getSeasons adalah helper generik untuk GET seasons by movie id
func getSeasons(c *gin.Context) {
	movieID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		verr.Handle(c, verr.NewAdminError(http.StatusBadRequest, "ID tidak valid", err))
		return
	}

	var seasons []models.Season
	database.DB.
		Where("movie_id = ?", movieID).
		Order("season_num asc").
		Preload("Episodes", "1=1 ORDER BY ep_num asc").
		Find(&seasons)

	c.JSON(http.StatusOK, gin.H{"success": true, "data": seasons})
}

// addSeason adalah helper generik untuk POST season
func addSeason(c *gin.Context) {
	movieID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		verr.Handle(c, verr.NewAdminError(http.StatusBadRequest, "ID tidak valid", err))
		return
	}

	var input struct {
		SeasonNum int `json:"season_num" binding:"required"`
		Episodes  []struct {
			EpNum int    `json:"ep_num"`
			Title string `json:"title"`
			URL1  string `json:"url1"`
			URL2  string `json:"url2"`
		} `json:"episodes"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		verr.Handle(c, verr.NewAdminError(http.StatusBadRequest, "Input tidak valid", err))
		return
	}

	// Cek duplikat season_num
	var existing models.Season
	if err := database.DB.Where("movie_id = ? AND season_num = ?", movieID, input.SeasonNum).First(&existing).Error; err == nil {
		verr.Handle(c, verr.NewAdminError(http.StatusConflict, "Season ini sudah ada", nil))
		return
	}

	// Validasi url1 tiap episode wajib diisi
	for i, ep := range input.Episodes {
		if ep.URL1 == "" {
			verr.Handle(c, verr.NewAdminError(
				http.StatusBadRequest,
				fmt.Sprintf("URL Server 1 episode %d wajib diisi", i+1),
				nil,
			))
			return
		}
	}

	var episodes []models.Episode
	for _, ep := range input.Episodes {
		episodes = append(episodes, models.Episode{
			EpNum: ep.EpNum,
			Title: ep.Title,
			URL1:  ep.URL1,
			URL2:  ep.URL2,
		})
	}

	season := models.Season{
		MovieID:   uint(movieID),
		SeasonNum: input.SeasonNum,
		Episodes:  episodes,
	}

	if err := database.DB.Create(&season).Error; err != nil {
		verr.Handle(c, verr.NewAdminError(http.StatusInternalServerError, "Gagal menyimpan season", err))
		return
	}

	database.DB.Preload("Episodes").First(&season, season.ID)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": season})
}

// ── Series Season Routes ──────────────────────────────────────
func GetSeriesSeasons(c *gin.Context) { getSeasons(c) }
func AddSeriesSeason(c *gin.Context)  { addSeason(c) }

// ── Anime Season Routes ───────────────────────────────────────
func GetAnimeSeasons(c *gin.Context) { getSeasons(c) }
func AddAnimeSeason(c *gin.Context)  { addSeason(c) }

// ── Season Actions ────────────────────────────────────────────

func DeleteSeason(c *gin.Context) {
	seasonID := c.Param("seasonId")

	var season models.Season
	if err := database.DB.First(&season, seasonID).Error; err != nil {
		verr.Handle(c, verr.NewAdminError(http.StatusNotFound, "Season tidak ditemukan", err))
		return
	}

	database.DB.Where("season_id = ?", season.ID).Delete(&models.Episode{})
	database.DB.Delete(&season)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Season berhasil dihapus"})
}

// ── Episode Actions ───────────────────────────────────────────

func UpdateEpisode(c *gin.Context) {
	epID := c.Param("epId")

	var input struct {
		Title string `json:"title"`
		URL1  string `json:"url1" binding:"required"`
		URL2  string `json:"url2"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		verr.Handle(c, verr.NewAdminError(http.StatusBadRequest, "URL Server 1 wajib diisi", err))
		return
	}

	var ep models.Episode
	if err := database.DB.First(&ep, epID).Error; err != nil {
		verr.Handle(c, verr.NewAdminError(http.StatusNotFound, "Episode tidak ditemukan", err))
		return
	}

	ep.Title = input.Title
	ep.URL1 = input.URL1
	ep.URL2 = input.URL2

	if err := database.DB.Save(&ep).Error; err != nil {
		verr.Handle(c, verr.NewAdminError(http.StatusInternalServerError, "Gagal update episode", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": ep})
}

func DeleteEpisode(c *gin.Context) {
	epID := c.Param("epId")
	if err := database.DB.Delete(&models.Episode{}, epID).Error; err != nil {
		verr.Handle(c, verr.NewAdminError(http.StatusInternalServerError, "Gagal menghapus episode", err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Episode berhasil dihapus"})
}
