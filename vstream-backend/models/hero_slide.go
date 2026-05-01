package models

import (
	"time"

	"gorm.io/gorm"
)

type HeroSlide struct {
	ID        uint           `json:"id"         gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	MovieID   uint           `json:"movie_id"   gorm:"not null;uniqueIndex"`
	SortOrder int            `json:"order"      gorm:"column:sort_order;not null;default:0"`
	Movie     Movie          `json:"movie"      gorm:"foreignKey:MovieID"`
}
