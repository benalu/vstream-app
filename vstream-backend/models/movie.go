package models

import (
	"time"

	"gorm.io/gorm"
)

// ContentType mendefinisikan jenis konten
type ContentType int

const (
	TypeMovie  ContentType = 1
	TypeSeries ContentType = 2
	TypeAnime  ContentType = 3
)

// Label untuk display di frontend
func (t ContentType) Label() string {
	switch t {
	case TypeMovie:
		return "Movie"
	case TypeSeries:
		return "Series"
	case TypeAnime:
		return "Anime"
	default:
		return "Unknown"
	}
}

// IsValid memastikan value type valid
func (t ContentType) IsValid() bool {
	return t == TypeMovie || t == TypeSeries || t == TypeAnime
}

type Movie struct {
	gorm.Model
	TmdbID   string      `json:"tmdb_id"`
	Type     ContentType `json:"type" gorm:"not null;default:1"` // 1=movie, 2=series, 3=anime
	Title    string      `json:"title"`
	Year     string      `json:"year"`
	Duration string      `json:"duration"`
	Rating   string      `json:"rating"`
	Genre    string      `json:"genre"`
	Poster   string      `json:"poster"`
	Backdrop string      `json:"backdrop"`
	LogoPath string      `json:"logo_path"`
	Overview string      `json:"overview"`
	URL1     string      `json:"url1"`
	URL2     string      `json:"url2"`
	URL3     string      `json:"url3"`
}

type WatchSession struct {
	Token      string    `gorm:"primaryKey" json:"token"`
	TmdbID     string    `json:"tmdb_id"`
	ClientHash string    `json:"client_hash"`
	ExpiresAt  time.Time `json:"expires_at"`
}
