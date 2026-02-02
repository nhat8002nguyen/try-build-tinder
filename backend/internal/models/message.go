package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Message struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	MatchID   uuid.UUID `gorm:"type:uuid;index:idx_match_created;not null" json:"match_id"`
	SenderID  uuid.UUID `gorm:"type:uuid;not null" json:"sender_id"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	IsRead    bool      `gorm:"default:false" json:"is_read"`
	CreatedAt time.Time `gorm:"index:idx_match_created" json:"created_at"`

	Match  Match `gorm:"foreignKey:MatchID" json:"-"`
	Sender User  `gorm:"foreignKey:SenderID" json:"sender,omitempty"`
}

func (m *Message) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}
