package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Match struct {
	ID            uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	User1ID       uuid.UUID  `gorm:"type:uuid;index;not null" json:"user1_id"`
	User2ID       uuid.UUID  `gorm:"type:uuid;index;not null" json:"user2_id"`
	MatchedAt     time.Time  `json:"matched_at"`
	LastMessageAt *time.Time `json:"last_message_at"`

	User1    User      `gorm:"foreignKey:User1ID" json:"user1,omitempty"`
	User2    User      `gorm:"foreignKey:User2ID" json:"user2,omitempty"`
	Messages []Message `gorm:"foreignKey:MatchID" json:"messages,omitempty"`
}

func (m *Match) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	m.MatchedAt = time.Now()
	return nil
}

func (m *Match) GetOtherUser(userID uuid.UUID) *User {
	if m.User1ID == userID {
		return &m.User2
	}
	return &m.User1
}

func (m *Match) GetOtherUserID(userID uuid.UUID) uuid.UUID {
	if m.User1ID == userID {
		return m.User2ID
	}
	return m.User1ID
}
