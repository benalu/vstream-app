package router

import (
	"vstream-backend/handlers"
	"vstream-backend/middleware"

	"github.com/gin-gonic/gin"
)

func setupPublicRoutes(api *gin.RouterGroup) {
	// Auth
	auth := api.Group("/auth")
	{
		auth.POST("/login", handlers.Login)
		auth.POST("/logout", handlers.Logout)
		auth.GET("/me", middleware.RequireAdmin, handlers.Me)
	}

	// Public content
	api.GET("/movies/:id", handlers.GetMoviePublic)
	api.GET("/series/:id", handlers.GetSeriesPublic)
	api.GET("/anime/:id", handlers.GetAnimePublic)
}
