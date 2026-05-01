package models

import "time"

type PlaybackLog struct {
	ID         uint      `json:"id"          gorm:"primarykey;autoIncrement"`
	MovieID    uint      `json:"movie_id"    gorm:"index"`
	MovieTitle string    `json:"movie_title"`
	MovieType  string    `json:"movie_type"`
	TmdbID     string    `json:"tmdb_id"`
	Server     string    `json:"server"` // "url1" | "url2" | "url3"
	URL        string    `json:"url"`
	ErrorType  string    `json:"error_type"` // "load_error" | "stalled" | "manual"
	UserAgent  string    `json:"user_agent"`
	ReportedAt time.Time `json:"reported_at"`
}
