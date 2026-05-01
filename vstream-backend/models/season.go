package models

import (
	"time"

	"gorm.io/gorm"
)

type Season struct {
	ID        uint           `json:"id"         gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	MovieID   uint           `json:"movie_id"   gorm:"not null;index"`
	SeasonNum int            `json:"season_num" gorm:"not null"`
	Episodes  []Episode      `json:"episodes"   gorm:"foreignKey:SeasonID"`
}

type Episode struct {
	ID        uint           `json:"id"         gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	SeasonID  uint           `json:"season_id"  gorm:"not null;index"`
	EpNum     int            `json:"ep_num"     gorm:"not null"`
	Title     string         `json:"title"`
	URL1      string         `json:"url1"`
	URL2      string         `json:"url2"`
}
