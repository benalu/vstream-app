package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"vstream-backend/database"
	"vstream-backend/models"
	"vstream-backend/pkg/verr"
	"vstream-backend/services"

	"github.com/gin-gonic/gin"
)

// ── Shared Helpers ────────────────────────────────────────────

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

	logos, _ := services.FetchLogos(tmdbID, contentType)
	logosJSON := "[]"
	if len(logos) > 0 {
		if b, err := json.Marshal(logos); err == nil {
			logosJSON = string(b)
		}
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
		Poster:      "https://image.tmdb.org/t/p/w500" + meta.PosterPath,
		Backdrop:    "https://image.tmdb.org/t/p/w1280" + meta.BackdropPath,
		Logos:       logosJSON,
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

// updateContent update URL, metadata, dan has_episodes
func updateContent(c *gin.Context, contentType string) {
	id := c.Param("id")

	var input struct {
		// Metadata fields — semua optional, hanya diupdate kalau dikirim
		Title    *string `json:"title"`
		Overview *string `json:"overview"`
		Poster   *string `json:"poster"`
		Backdrop *string `json:"backdrop"`
		Genre    *string `json:"genre"`
		Year     *string `json:"year"`
		Duration *string `json:"duration"`
		Rating   *string `json:"rating"`

		// Episode flag (anime only)
		HasEpisodes *bool `json:"has_episodes"`

		// URL fields
		URL1 *string `json:"url1"`
		URL2 *string `json:"url2"`
		URL3 *string `json:"url3"`

		// Refresh metadata dari TMDB
		RefreshFromTMDB bool `json:"refresh_from_tmdb"`
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

	// ── Refresh dari TMDB (re-fetch semua metadata) ───────────
	if input.RefreshFromTMDB {
		meta, err := services.FetchMetadata(item.TmdbID, item.Type)
		if err != nil {
			verr.Handle(c, verr.NewAdminError(http.StatusBadGateway, "Gagal mengambil data dari TMDB", err))
			return
		}

		var genres []string
		for _, g := range meta.Genres {
			genres = append(genres, g.Name)
		}

		logos, _ := services.FetchLogos(item.TmdbID, item.Type)
		logosJSON := "[]"
		if len(logos) > 0 {
			if b, err := json.Marshal(logos); err == nil {
				logosJSON = string(b)
			}
		}

		item.Title = meta.NormalizedTitle()
		item.Year = strings.Split(meta.NormalizedReleaseDate(), "-")[0]
		item.Duration = fmt.Sprintf("%d min", meta.NormalizedRuntime())
		item.Rating = fmt.Sprintf("%.1f", meta.VoteAverage)
		item.Genre = strings.Join(genres, ", ")
		item.Poster = "https://image.tmdb.org/t/p/w500" + meta.PosterPath
		item.Backdrop = "https://image.tmdb.org/t/p/w1280" + meta.BackdropPath
		item.Logos = logosJSON
		item.Overview = meta.Overview
	}

	// ── Update metadata manual (tiap field hanya kalau dikirim) ──
	if !input.RefreshFromTMDB {
		if input.Title != nil {
			item.Title = strings.TrimSpace(*input.Title)
		}
		if input.Overview != nil {
			item.Overview = strings.TrimSpace(*input.Overview)
		}
		if input.Poster != nil {
			item.Poster = strings.TrimSpace(*input.Poster)
		}
		if input.Backdrop != nil {
			item.Backdrop = strings.TrimSpace(*input.Backdrop)
		}
		if input.Genre != nil {
			item.Genre = strings.TrimSpace(*input.Genre)
		}
		if input.Year != nil {
			item.Year = strings.TrimSpace(*input.Year)
		}
		if input.Duration != nil {
			item.Duration = strings.TrimSpace(*input.Duration)
		}
		if input.Rating != nil {
			item.Rating = strings.TrimSpace(*input.Rating)
		}
	}

	// ── has_episodes (anime only) ─────────────────────────────
	if contentType == "anime" && input.HasEpisodes != nil {
		item.HasEpisodes = *input.HasEpisodes
	}

	// ── URL (hanya kalau tidak punya episode) ─────────────────
	willHaveEpisodes := item.HasEpisodes
	if contentType == "anime" && input.HasEpisodes != nil {
		willHaveEpisodes = *input.HasEpisodes
	}

	if !willHaveEpisodes {
		// Validasi url1 wajib — dari input baru atau sudah ada di DB
		newURL1 := item.URL1
		if input.URL1 != nil {
			newURL1 = strings.TrimSpace(*input.URL1)
		}
		if newURL1 == "" {
			verr.Handle(c, verr.NewAdminError(http.StatusBadRequest, "URL Server 1 wajib diisi", nil))
			return
		}
		item.URL1 = newURL1
		if input.URL2 != nil {
			item.URL2 = strings.TrimSpace(*input.URL2)
		}
		if input.URL3 != nil {
			item.URL3 = strings.TrimSpace(*input.URL3)
		}
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
