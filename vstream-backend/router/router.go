package router

import (
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// Setup mendaftarkan semua route ke engine Gin yang diberikan.
func Setup(r *gin.Engine) {
	// Serve subtitle files dengan CORS header
	r.GET("/static/subtitles/:filename", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		filename := c.Param("filename")
		c.File(filepath.Join("data", "subtitles", filename))
	})

	api := r.Group("/api")
	setupPublicRoutes(api)
	setupAdminRoutes(api)
}
