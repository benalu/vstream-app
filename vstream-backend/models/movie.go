package models

import (
	"time"

	"gorm.io/gorm"
)

type Movie struct {
	gorm.Model
	TmdbID      string `json:"tmdb_id"`
	Type        string `json:"type" gorm:"not null;default:'movie'"` // "movie" | "series" | "anime"
	HasEpisodes bool   `json:"has_episodes" gorm:"default:false"`    // hanya relevan untuk anime
	Title       string `json:"title"`
	Year        string `json:"year"`
	Duration    string `json:"duration"`
	Rating      string `json:"rating"`
	Genre       string `json:"genre"`
	Poster      string `json:"poster"`
	Backdrop    string `json:"backdrop"`
	Logos       string `json:"logos"`
	Overview    string `json:"overview"`
	URL1        string `json:"url1"`
	URL2        string `json:"url2"`
	URL3        string `json:"url3"`
}

type WatchSession struct {
	Token      string    `gorm:"primaryKey" json:"token"`
	TmdbID     string    `json:"tmdb_id"`
	ClientHash string    `json:"client_hash"`
	ExpiresAt  time.Time `json:"expires_at"`
}
