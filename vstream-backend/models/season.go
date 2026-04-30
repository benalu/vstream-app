package models

import "gorm.io/gorm"

type Season struct {
	gorm.Model
	MovieID   uint      `json:"movie_id" gorm:"not null;index"`
	SeasonNum int       `json:"season_num" gorm:"not null"`
	Episodes  []Episode `json:"episodes" gorm:"foreignKey:SeasonID"`
}

type Episode struct {
	gorm.Model
	SeasonID uint   `json:"season_id" gorm:"not null;index"`
	EpNum    int    `json:"ep_num" gorm:"not null"`
	Title    string `json:"title"`
	URL1     string `json:"url1"`
	URL2     string `json:"url2"`
}
