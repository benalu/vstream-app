package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"vstream-backend/database"
	"vstream-backend/models"
	"vstream-backend/pkg/verr"
	"vstream-backend/services"

	"github.com/gin-gonic/gin"
)

// ── Helpers ───────────────────────────────────────────────────

var validTypes = map[string]bool{
	"movie":  true,
	"series": true,
	"anime":  true,
}

// buildMovieFromTMDB mengambil metadata TMDB dan mapping ke model Movie
func buildMovieFromTMDB(tmdbID string, contentType string, hasEpisodes bool, url1, url2, url3 string) (*models.Movie, error) {
	meta, err := services.FetchMetadata(tmdbID, contentType)
	if err != nil {
		return nil, err
	}

	var genres []string
	for _, g := range meta.Genres {
		genres = append(genres, g.Name)
	}

	var logo string
	if len(meta.Images.Logos) > 0 {
		logo = "https://image.tmdb.org/t/p/original" + meta.Images.Logos[0].FilePath
	}

	return &models.Movie{
		TmdbID:      tmdbID,
		Type:        contentType,
		HasEpisodes: hasEpisodes,
		Title:       meta.NormalizedTitle(),
		Year:        strings.Split(meta.NormalizedReleaseDate(), "-")[0],
		Duration:    fmt.Sprintf("%d min", meta.NormalizedRuntime()),
		Rating:      fmt.Sprintf("%.1f", meta.VoteAverage),
		Genre:       strings.Join(genres, ", "),
		Poster:      "https://image.tmdb.org/t/p/original" + meta.PosterPath,
		Backdrop:    "https://image.tmdb.org/t/p/original" + meta.BackdropPath,
		LogoPath:    logo,
		Overview:    meta.Overview,
		URL1:        url1,
		URL2:        url2,
		URL3:        url3,
	}, nil
}

// addContent adalah handler generik yang dipakai oleh AddMovie, AddSeries, AddAnime
func addContent(c *gin.Context, contentType string) {
	var input struct {
		TmdbID      string `json:"tmdb_id" binding:"required"`
		HasEpisodes bool   `json:"has_episodes"` // hanya dibaca untuk anime
		URL1        string `json:"url1"`
		URL2        string `json:"url2"`
		URL3        string `json:"url3"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		verr.Handle(c, verr.NewAdminError(http.StatusBadRequest, "Input tidak valid", err))
		return
	}

	// Series selalu punya episode, jadi has_episodes tidak relevan
	// Anime: ditentukan dari input
	// Movie: tidak punya episode
	hasEpisodes := false
	switch contentType {
	case "series":
		hasEpisodes = true
	case "anime":
		hasEpisodes = input.HasEpisodes
	}

	// URL1 wajib kalau tidak punya episode
	// (series/anime dengan episode: URL diisi di level episode, bukan konten)
	if !hasEpisodes && input.URL1 == "" {
		verr.Handle(c, verr.NewAdminError(http.StatusBadRequest, "URL Server 1 wajib diisi", nil))
		return
	}

	movie, err := buildMovieFromTMDB(input.TmdbID, contentType, hasEpisodes, input.URL1, input.URL2, input.URL3)
	if err != nil {
		verr.Handle(c, verr.NewAdminError(http.StatusBadGateway, "TMDB fetch gagal", err))
		return
	}

	if err := database.DB.Save(movie).Error; err != nil {
		verr.Handle(c, verr.NewAdminError(http.StatusInternalServerError, "Gagal menyimpan ke database", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": movie})
}

// getAllContent adalah handler generik untuk GET list per tipe
func getAllContent(c *gin.Context, contentType string) {
	var items []models.Movie
	database.DB.Where("type = ?", contentType).Order("created_at desc").Find(&items)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": items})
}

// updateContent update URL (dan has_episodes untuk anime)
func updateContent(c *gin.Context, contentType string) {
	id := c.Param("id")

	var input struct {
		HasEpisodes *bool  `json:"has_episodes"` // pointer: nil = tidak dikirim
		URL1        string `json:"url1"`
		URL2        string `json:"url2"`
		URL3        string `json:"url3"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		verr.Handle(c, verr.NewAdminError(http.StatusBadRequest, "Input tidak valid", err))
		return
	}

	var item models.Movie
	if err := database.DB.Where("id = ? AND type = ?", id, contentType).First(&item).Error; err != nil {
		verr.Handle(c, verr.NewAdminError(http.StatusNotFound, "Konten tidak ditemukan", err))
		return
	}

	// Untuk anime: boleh ubah has_episodes
	if contentType == "anime" && input.HasEpisodes != nil {
		item.HasEpisodes = *input.HasEpisodes
	}

	// URL hanya diupdate kalau tidak punya episode
	willHaveEpisodes := item.HasEpisodes
	if contentType == "anime" && input.HasEpisodes != nil {
		willHaveEpisodes = *input.HasEpisodes
	}

	if !willHaveEpisodes {
		if input.URL1 == "" {
			verr.Handle(c, verr.NewAdminError(http.StatusBadRequest, "URL Server 1 wajib diisi", nil))
			return
		}
		item.URL1 = input.URL1
		item.URL2 = input.URL2
		item.URL3 = input.URL3
	}

	if err := database.DB.Save(&item).Error; err != nil {
		verr.Handle(c, verr.NewAdminError(http.StatusInternalServerError, "Gagal update database", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": item})
}

// deleteContent hapus konten berdasarkan id + type
func deleteContent(c *gin.Context, contentType string) {
	id := c.Param("id")

	var item models.Movie
	if err := database.DB.Where("id = ? AND type = ?", id, contentType).First(&item).Error; err != nil {
		verr.Handle(c, verr.NewAdminError(http.StatusNotFound, "Konten tidak ditemukan", err))
		return
	}

	// Hapus seasons & episodes terkait kalau ada
	if item.HasEpisodes || item.Type == "series" {
		var seasons []models.Season
		database.DB.Where("movie_id = ?", item.ID).Find(&seasons)
		for _, s := range seasons {
			database.DB.Where("season_id = ?", s.ID).Delete(&models.Episode{})
		}
		database.DB.Where("movie_id = ?", item.ID).Delete(&models.Season{})
	}

	database.DB.Delete(&item)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Konten berhasil dihapus"})
}

// ── Movie Handlers ────────────────────────────────────────────

func AddMovie(c *gin.Context)     { addContent(c, "movie") }
func GetAllMovies(c *gin.Context) { getAllContent(c, "movie") }
func UpdateMovie(c *gin.Context)  { updateContent(c, "movie") }
func DeleteMovie(c *gin.Context)  { deleteContent(c, "movie") }

func GetMoviePublic(c *gin.Context) {
	id := c.Param("id")
	var movie models.Movie
	if err := database.DB.Where("tmdb_id = ? AND type = 'movie'", id).First(&movie).Error; err != nil {
		verr.Handle(c, verr.NewPublicError(http.StatusNotFound, "Film tidak ditemukan", "Movie lookup failed", err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": movie})
}

// ── Series Handlers ───────────────────────────────────────────

func AddSeries(c *gin.Context)    { addContent(c, "series") }
func GetAllSeries(c *gin.Context) { getAllContent(c, "series") }
func UpdateSeries(c *gin.Context) { updateContent(c, "series") }
func DeleteSeries(c *gin.Context) { deleteContent(c, "series") }

func GetSeriesPublic(c *gin.Context) {
	id := c.Param("id")
	var series models.Movie
	if err := database.DB.Where("tmdb_id = ? AND type = 'series'", id).First(&series).Error; err != nil {
		verr.Handle(c, verr.NewPublicError(http.StatusNotFound, "Series tidak ditemukan", "Series lookup failed", err))
		return
	}

	// Sertakan seasons + episodes untuk public
	var seasons []models.Season
	database.DB.Where("movie_id = ?", series.ID).
		Order("season_num asc").
		Preload("Episodes", "1=1 ORDER BY ep_num asc").
		Find(&seasons)

	c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{
		"info":    series,
		"seasons": seasons,
	}})
}

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

// ── TMDB Preview (shared) ─────────────────────────────────────

func PreviewTMDB(c *gin.Context) {
	id := c.Param("id")
	contentType := c.DefaultQuery("type", "movie")
	if id == "" {
		verr.Handle(c, verr.NewAdminError(http.StatusBadRequest, "ID TMDB tidak boleh kosong", nil))
		return
	}
	meta, err := services.FetchMetadata(id, contentType)
	if err != nil {
		verr.Handle(c, verr.NewAdminError(http.StatusNotFound, "Tidak ditemukan di TMDB", err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": meta})
}
