package models

import "gorm.io/gorm"

type HeroSlide struct {
	gorm.Model
	MovieID   uint  `json:"movie_id" gorm:"not null;uniqueIndex"`
	SortOrder int   `json:"order" gorm:"column:sort_order;not null;default:0"`
	Movie     Movie `json:"movie" gorm:"foreignKey:MovieID"`
}
