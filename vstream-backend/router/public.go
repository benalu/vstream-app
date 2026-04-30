package router

import (
	"vstream-backend/handlers"
	"vstream-backend/middleware"

	"github.com/gin-gonic/gin"
)

func setupPublicRoutes(api *gin.RouterGroup) {
	// ── Auth ──────────────────────────────────────────────────
	auth := api.Group("/auth")
	{
		auth.POST("/login", handlers.Login)
		auth.POST("/logout", handlers.Logout)
		auth.GET("/me", middleware.RequireAdmin, handlers.Me)
	}

	// ── Aggregated endpoints ───────────────────────────────────
	api.GET("/featured", handlers.GetFeatured) // hero section
	api.GET("/top10", handlers.GetTop10)       // TOP 10 row

	// ── Public listings ───────────────────────────────────────
	api.GET("/movies", handlers.GetMoviesPublicList) // ?limit=&offset=
	api.GET("/series", handlers.GetSeriesPublicList)
	api.GET("/anime", handlers.GetAnimePublicList)

	// ── Public single content (by tmdb_id) ────────────────────
	api.GET("/movies/:id", handlers.GetMoviePublic)
	api.GET("/series/:id", handlers.GetSeriesPublic)
	api.GET("/anime/:id", handlers.GetAnimePublic)
}
