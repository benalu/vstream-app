package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"vstream-backend/database"
	"vstream-backend/models"
	"vstream-backend/pkg/verr"
	"vstream-backend/services"

	"github.com/gin-gonic/gin"
)

func AddMovie(c *gin.Context) {
	var input struct {
		TmdbID string `json:"tmdb_id"`
		Type   int    `json:"type"` // 1=movie, 2=series, 3=anime
		URL1   string `json:"url1"`
		URL2   string `json:"url2"`
		URL3   string `json:"url3"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		verr.Handle(c, verr.NewAdminError(http.StatusBadRequest, "Invalid JSON input", err))
		return
	}

	// Validasi type
	contentType := models.ContentType(input.Type)
	if !contentType.IsValid() {
		verr.Handle(c, verr.NewAdminError(http.StatusBadRequest, "Type tidak valid (1=Movie, 2=Series, 3=Anime)", nil))
		return
	}

	// 1. Ambil data dari TMDB
	meta, err := services.FetchMetadata(input.TmdbID)
	if err != nil {
		verr.Handle(c, verr.NewAdminError(http.StatusBadGateway, "TMDB fetch failed", err))
		return
	}

	// 2. Mapping ke Model GORM
	var genres []string
	for _, g := range meta.Genres {
		genres = append(genres, g.Name)
	}

	var logo string
	if len(meta.Images.Logos) > 0 {
		logo = "https://image.tmdb.org/t/p/original" + meta.Images.Logos[0].FilePath
	}

	movie := models.Movie{
		TmdbID:   input.TmdbID,
		Type:     contentType,
		Title:    meta.Title,
		Year:     strings.Split(meta.ReleaseDate, "-")[0],
		Duration: fmt.Sprintf("%d min", meta.Runtime),
		Rating:   fmt.Sprintf("%.1f", meta.VoteAverage),
		Genre:    strings.Join(genres, ", "),
		Poster:   "https://image.tmdb.org/t/p/original" + meta.PosterPath,
		Backdrop: "https://image.tmdb.org/t/p/original" + meta.BackdropPath,
		LogoPath: logo,
		Overview: meta.Overview,
		URL1:     input.URL1,
		URL2:     input.URL2,
		URL3:     input.URL3,
	}

	// 3. Simpan/Update ke Database (Upsert)
	if err := database.DB.Save(&movie).Error; err != nil {
		verr.Handle(c, verr.NewAdminError(http.StatusInternalServerError, "DB Save failed", err))
		return
	}

	c.JSON(200, gin.H{"success": true, "data": movie})
}

// GetAllMovies mendukung filter via query param ?type=1
// Contoh: GET /admin/movies        → semua konten
//
//	GET /admin/movies?type=1 → hanya movie
//	GET /admin/movies?type=2 → hanya series
//	GET /admin/movies?type=3 → hanya anime
func GetAllMovies(c *gin.Context) {
	var movies []models.Movie

	query := database.DB.Order("created_at desc")

	// Filter by type jika ada query param
	if typeParam := c.Query("type"); typeParam != "" {
		typeInt, err := strconv.Atoi(typeParam)
		if err != nil || !models.ContentType(typeInt).IsValid() {
			verr.Handle(c, verr.NewAdminError(http.StatusBadRequest, "Query param type tidak valid", nil))
			return
		}
		query = query.Where("type = ?", typeInt)
	}

	query.Find(&movies)
	c.JSON(200, gin.H{"success": true, "data": movies})
}

func GetMoviePublic(c *gin.Context) {
	id := c.Param("id")
	var movie models.Movie

	if err := database.DB.First(&movie, "tmdb_id = ?", id).Error; err != nil {
		verr.Handle(c, verr.NewPublicError(http.StatusNotFound, "Konten tidak ditemukan", "Movie lookup failed", err))
		return
	}

	c.JSON(200, gin.H{"success": true, "data": movie})
}

// DeleteMovie menghapus konten berdasarkan ID
func DeleteMovie(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&models.Movie{}, id).Error; err != nil {
		verr.Handle(c, verr.NewAdminError(http.StatusInternalServerError, "Gagal menghapus konten", err))
		return
	}
	c.JSON(200, gin.H{"success": true, "message": "Konten berhasil dihapus"})
}

// UpdateMovie update URL + boleh update type juga
func UpdateMovie(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		Type *int   `json:"type"` // pointer agar bisa detect "tidak dikirim"
		URL1 string `json:"url1"`
		URL2 string `json:"url2"`
		URL3 string `json:"url3"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		verr.Handle(c, verr.NewAdminError(http.StatusBadRequest, "Input tidak valid", err))
		return
	}

	var movie models.Movie
	if err := database.DB.First(&movie, id).Error; err != nil {
		verr.Handle(c, verr.NewAdminError(http.StatusNotFound, "Konten tidak ditemukan", err))
		return
	}

	// Update type jika dikirim
	if input.Type != nil {
		contentType := models.ContentType(*input.Type)
		if !contentType.IsValid() {
			verr.Handle(c, verr.NewAdminError(http.StatusBadRequest, "Type tidak valid (1=Movie, 2=Series, 3=Anime)", nil))
			return
		}
		movie.Type = contentType
	}

	movie.URL1 = input.URL1
	movie.URL2 = input.URL2
	movie.URL3 = input.URL3

	if err := database.DB.Save(&movie).Error; err != nil {
		verr.Handle(c, verr.NewAdminError(http.StatusInternalServerError, "Gagal update database", err))
		return
	}

	c.JSON(200, gin.H{"success": true, "data": movie})
}

func PreviewTMDB(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		verr.Handle(c, verr.NewAdminError(http.StatusBadRequest, "ID TMDB tidak boleh kosong", nil))
		return
	}

	meta, err := services.FetchMetadata(id)
	if err != nil {
		verr.Handle(c, verr.NewAdminError(http.StatusNotFound, "Film tidak ditemukan di TMDB", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    meta,
	})
}
