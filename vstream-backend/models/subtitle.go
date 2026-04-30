// models/subtitle.go
package models

import "gorm.io/gorm"

type Subtitle struct {
	gorm.Model
	MovieID  uint   `json:"movie_id" gorm:"not null;index"`
	TmdbID   string `json:"tmdb_id"`
	Type     string `json:"type"`  // "movie" | "series" | "anime"
	Lang     string `json:"lang"`  // "id" | "en"
	Label    string `json:"label"` // "Indonesia" | "English"
	Filename string `json:"filename"`
}
