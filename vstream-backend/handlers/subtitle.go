package handlers

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"vstream-backend/database"
	"vstream-backend/models"
	"vstream-backend/pkg/verr"

	"github.com/gin-gonic/gin"
)

func UploadSubtitle(c *gin.Context) {
	movieID := c.Param("id")
	lang := c.PostForm("lang")   // "id" atau "en"
	label := c.PostForm("label") // "Indonesia" atau "English"

	if lang == "" || label == "" {
		verr.Handle(c, verr.NewAdminError(http.StatusBadRequest, "lang dan label wajib diisi", nil))
		return
	}

	// Cari movie dulu untuk dapat tmdb_id dan type
	var movie models.Movie
	if err := database.DB.First(&movie, movieID).Error; err != nil {
		verr.Handle(c, verr.NewAdminError(http.StatusNotFound, "Konten tidak ditemukan", err))
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		verr.Handle(c, verr.NewAdminError(http.StatusBadRequest, "File tidak ditemukan", err))
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".vtt" && ext != ".srt" {
		verr.Handle(c, verr.NewAdminError(http.StatusBadRequest, "Hanya .vtt atau .srt yang diizinkan", nil))
		return
	}

	// Pastikan folder ada
	subDir := filepath.Join("data", "subtitles")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		verr.Handle(c, verr.NewAdminError(http.StatusInternalServerError, "Gagal membuat folder", err))
		return
	}

	// Nama file: {type}-{tmdb_id}-{lang}.vtt
	filename := fmt.Sprintf("%s-%s-%s.vtt", movie.Type, movie.TmdbID, lang)
	savePath := filepath.Join(subDir, filename)

	if ext == ".srt" {
		if err := convertSRTtoVTT(file, savePath); err != nil {
			verr.Handle(c, verr.NewAdminError(http.StatusInternalServerError, "Konversi SRT gagal", err))
			return
		}
	} else {
		if err := c.SaveUploadedFile(file, savePath); err != nil {
			verr.Handle(c, verr.NewAdminError(http.StatusInternalServerError, "Gagal menyimpan file", err))
			return
		}
	}

	// Upsert ke DB — kalau lang sama, timpa yang lama
	id, _ := strconv.Atoi(movieID)
	var sub models.Subtitle
	database.DB.Where("movie_id = ? AND lang = ?", id, lang).FirstOrInit(&sub)

	sub.MovieID = uint(id)
	sub.TmdbID = movie.TmdbID
	sub.Type = movie.Type
	sub.Lang = lang
	sub.Label = label
	sub.Filename = filename

	database.DB.Save(&sub)

	c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{
		"filename": filename,
		"url":      "/static/subtitles/" + filename,
	}})
}

func DeleteSubtitle(c *gin.Context) {
	subtitleID := c.Param("subtitleId")

	var sub models.Subtitle
	if err := database.DB.First(&sub, subtitleID).Error; err != nil {
		verr.Handle(c, verr.NewAdminError(http.StatusNotFound, "Subtitle tidak ditemukan", err))
		return
	}

	// Hapus file fisik
	savePath := filepath.Join("data", "subtitles", sub.Filename)
	os.Remove(savePath) // ignore error kalau file sudah tidak ada

	database.DB.Delete(&sub)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Subtitle berhasil dihapus"})
}

func GetSubtitles(c *gin.Context) {
	movieID := c.Param("id")
	var subs []models.Subtitle
	database.DB.Where("movie_id = ?", movieID).Find(&subs)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": subs})
}

func convertSRTtoVTT(file *multipart.FileHeader, savePath string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	content, err := io.ReadAll(src)
	if err != nil {
		return err
	}

	// SRT pakai koma untuk milidetik → VTT pakai titik
	vtt := "WEBVTT\n\n" + strings.ReplaceAll(string(content), ",", ".")
	return os.WriteFile(savePath, []byte(vtt), 0644)
}
