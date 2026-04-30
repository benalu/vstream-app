package router

import (
	"github.com/gin-gonic/gin"
)

// Setup mendaftarkan semua route ke engine Gin yang diberikan.
func Setup(r *gin.Engine) {
	api := r.Group("/api")
	setupPublicRoutes(api)
	setupAdminRoutes(api)
}
