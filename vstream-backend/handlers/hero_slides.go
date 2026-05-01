package handlers

import (
	"net/http"
	"vstream-backend/database"
	"vstream-backend/models"
	"vstream-backend/pkg/verr"

	"github.com/gin-gonic/gin"
)

// GetHeroSlides — public, return max 5 slides urut by order
func GetHeroSlides(c *gin.Context) {
	var slides []models.HeroSlide
	database.DB.
		Preload("Movie").
		Order("sort_order asc").
		Limit(5).
		Find(&slides)

	// Sanitize: strip URL stream dari Movie
	type SlideDTO struct {
		ID    uint        `json:"id"`
		Order int         `json:"order"`
		Movie PublicMovie `json:"movie"`
	}

	result := make([]SlideDTO, 0, len(slides))
	for _, s := range slides {
		result = append(result, SlideDTO{
			ID:    s.ID,
			Order: s.SortOrder,
			Movie: sanitize(s.Movie),
		})
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": result})
}

// GetHeroSlidesAdmin — admin, return semua slides
func GetHeroSlidesAdmin(c *gin.Context) {
	var slides []models.HeroSlide
	database.DB.
		Preload("Movie").
		Order("sort_order asc").
		Find(&slides)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": slides})
}

// AddHeroSlide — admin, tambah movie ke hero slider
func AddHeroSlide(c *gin.Context) {
	var input struct {
		MovieID uint `json:"movie_id" binding:"required"`
		Order   int  `json:"order"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		verr.Handle(c, verr.NewAdminError(http.StatusBadRequest, "Input tidak valid", err))
		return
	}

	// Cek movie exists
	var movie models.Movie
	if err := database.DB.First(&movie, input.MovieID).Error; err != nil {
		verr.Handle(c, verr.NewAdminError(http.StatusNotFound, "Film tidak ditemukan", err))
		return
	}

	// Max 5 slides
	var count int64
	database.DB.Model(&models.HeroSlide{}).Count(&count)
	if count >= 5 {
		verr.Handle(c, verr.NewAdminError(http.StatusBadRequest, "Maksimal 5 slide hero", nil))
		return
	}

	slide := models.HeroSlide{
		MovieID:   input.MovieID,
		SortOrder: input.Order,
	}
	if err := database.DB.Create(&slide).Error; err != nil {
		verr.Handle(c, verr.NewAdminError(http.StatusInternalServerError, "Gagal menyimpan slide", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": slide})
}

// UpdateHeroSlideOrder — admin, update urutan slide
func UpdateHeroSlideOrder(c *gin.Context) {
	id := c.Param("id")

	var input struct {
		Order int `json:"order" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		verr.Handle(c, verr.NewAdminError(http.StatusBadRequest, "Input tidak valid", err))
		return
	}

	var slide models.HeroSlide
	if err := database.DB.First(&slide, id).Error; err != nil {
		verr.Handle(c, verr.NewAdminError(http.StatusNotFound, "Slide tidak ditemukan", err))
		return
	}

	slide.SortOrder = input.Order
	database.DB.Save(&slide)

	c.JSON(http.StatusOK, gin.H{"success": true, "data": slide})
}

// DeleteHeroSlide — admin, hapus slide
func DeleteHeroSlide(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&models.HeroSlide{}, id).Error; err != nil {
		verr.Handle(c, verr.NewAdminError(http.StatusInternalServerError, "Gagal menghapus slide", err))
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Slide berhasil dihapus"})
}
