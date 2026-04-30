package router

import (
	"vstream-backend/handlers"
	"vstream-backend/middleware"

	"github.com/gin-gonic/gin"
)

func setupAdminRoutes(api *gin.RouterGroup) {
	admin := api.Group("/admin")
	admin.Use(middleware.RequireAdmin)
	{
		admin.GET("/dashboard", handlers.GetDashboardStats)
		registerMovieRoutes(admin)
		registerSeriesRoutes(admin)
		registerAnimeRoutes(admin)
		registerSeasonEpisodeRoutes(admin)
	}
}

func registerMovieRoutes(admin *gin.RouterGroup) {
	admin.POST("/movies", handlers.AddMovie)
	admin.GET("/movies", handlers.GetAllMovies)
	admin.PUT("/movies/:id", handlers.UpdateMovie)
	admin.DELETE("/movies/:id", handlers.DeleteMovie)
	admin.GET("/tmdb/:id", handlers.PreviewTMDB)
}

func registerSeriesRoutes(admin *gin.RouterGroup) {
	admin.POST("/series", handlers.AddSeries)
	admin.GET("/series", handlers.GetAllSeries)
	admin.PUT("/series/:id", handlers.UpdateSeries)
	admin.DELETE("/series/:id", handlers.DeleteSeries)
	admin.GET("/series/:id/seasons", handlers.GetSeriesSeasons)
	admin.POST("/series/:id/seasons", handlers.AddSeriesSeason)
}

func registerAnimeRoutes(admin *gin.RouterGroup) {
	admin.POST("/anime", handlers.AddAnime)
	admin.GET("/anime", handlers.GetAllAnime)
	admin.PUT("/anime/:id", handlers.UpdateAnime)
	admin.DELETE("/anime/:id", handlers.DeleteAnime)
	admin.GET("/anime/:id/seasons", handlers.GetAnimeSeasons)
	admin.POST("/anime/:id/seasons", handlers.AddAnimeSeason)
}

func registerSeasonEpisodeRoutes(admin *gin.RouterGroup) {
	admin.DELETE("/seasons/:seasonId", handlers.DeleteSeason)
	admin.POST("/seasons/:seasonId/episodes", handlers.AddEpisode)
	admin.PUT("/episodes/:epId", handlers.UpdateEpisode)
	admin.DELETE("/episodes/:epId", handlers.DeleteEpisode)
}
