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
		auth.POST("/login", middleware.LoginLimiter(), handlers.Login)
		auth.POST("/logout", handlers.Logout)
		auth.GET("/me", middleware.RequireAdmin, handlers.Me)
	}

	cache60 := middleware.PublicCache(60)

	// ── Aggregated endpoints ───────────────────────────────────
	api.GET("/featured", handlers.GetFeatured)               // hero section
	api.GET("/hero-slides", cache60, handlers.GetHeroSlides) // hero slider
	api.GET("/top10", cache60, handlers.GetTop10)            // TOP 10 row

	// ── Public listings ───────────────────────────────────────
	api.GET("/movies", cache60, handlers.GetMoviesPublicList) // ?limit=&offset=
	api.GET("/series", cache60, handlers.GetSeriesPublicList)
	api.GET("/anime", cache60, handlers.GetAnimePublicList)

	// ── Public single content (by tmdb_id) ────────────────────
	api.GET("/movies/:id", handlers.GetMoviePublic)
	api.GET("/series/:id", handlers.GetSeriesPublic)
	api.GET("/anime/:id", handlers.GetAnimePublic)

	// ── Watch data (includes stream URLs) ─────────────────────
	api.GET("/watch/:type/:tmdb_id", handlers.GetWatchData)

	// ── Playback error logging ─────────────────────────────────
	api.POST("/logs/playback-error", middleware.PlaybackLogLimiter(), handlers.ReportPlaybackError)

}
