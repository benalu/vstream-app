package models

import "time"

type Movie struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TmdbID    string    `gorm:"uniqueIndex;not null" json:"tmdb_id"`
	Title     string    `json:"title"`
	Year      string    `json:"year"`
	Duration  string    `json:"duration"`
	Rating    string    `json:"rating"`
	Genre     string    `json:"genre"`
	Poster    string    `json:"poster"`
	Backdrop  string    `json:"backdrop"`
	LogoPath  string    `json:"logo_path"`
	Overview  string    `json:"overview"`
	URL1      string    `json:"url1"`
	URL2      string    `json:"url2"`
	URL3      string    `json:"url3"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type WatchSession struct {
	Token      string    `gorm:"primaryKey" json:"token"`
	TmdbID     string    `json:"tmdb_id"`
	ClientHash string    `json:"client_hash"`
	ExpiresAt  time.Time `json:"expires_at"`
}
