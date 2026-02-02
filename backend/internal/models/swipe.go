package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SwipeDirection string

const (
	SwipeLike    SwipeDirection = "like"
	SwipeDislike SwipeDirection = "dislike"
)

type Swipe struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	SwiperID  uuid.UUID      `gorm:"type:uuid;index:idx_swiper_target;not null" json:"swiper_id"`
	TargetID  uuid.UUID      `gorm:"type:uuid;index:idx_swiper_target;not null" json:"target_id"`
	Direction SwipeDirection `gorm:"type:varchar(20);not null" json:"direction"`
	CreatedAt time.Time      `gorm:"index" json:"created_at"`

	Swiper User `gorm:"foreignKey:SwiperID" json:"-"`
	Target User `gorm:"foreignKey:TargetID" json:"-"`
}

func (s *Swipe) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
