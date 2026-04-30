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

func AddMovie(c *gin.Context) {
	var input struct {
		TmdbID string `json:"tmdb_id"`
		URL1   string `json:"url1"`
		URL2   string `json:"url2"`
		URL3   string `json:"url3"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		verr.Handle(c, verr.NewAdminError(http.StatusBadRequest, "Invalid JSON input", err))
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

func GetAllMovies(c *gin.Context) {
	var movies []models.Movie
	database.DB.Order("created_at desc").Find(&movies)
	c.JSON(200, gin.H{"success": true, "data": movies})
}

func GetMoviePublic(c *gin.Context) {
	id := c.Param("id")
	var movie models.Movie

	if err := database.DB.First(&movie, "tmdb_id = ?", id).Error; err != nil {
		verr.Handle(c, verr.NewPublicError(http.StatusNotFound, "Film tidak ditemukan", "Movie lookup failed", err))
		return
	}

	c.JSON(200, gin.H{"success": true, "data": movie})
}

// Handler untuk menghapus film
func DeleteMovie(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&models.Movie{}, id).Error; err != nil {
		verr.Handle(c, verr.NewAdminError(http.StatusInternalServerError, "Gagal menghapus film", err))
		return
	}
	c.JSON(200, gin.H{"success": true, "message": "Film berhasil dihapus"})
}

// Handler untuk update URL film (tanpa fetch ulang TMDB agar cepat)
func UpdateMovie(c *gin.Context) {
	id := c.Param("id")
	var input struct {
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
		verr.Handle(c, verr.NewAdminError(http.StatusNotFound, "Film tidak ditemukan", err))
		return
	}

	// Update field yang diizinkan saja
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

	// Gunakan service TMDB yang sudah kamu buat sebelumnya
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
