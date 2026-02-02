package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationType string

const (
	NotificationTypeMatch   NotificationType = "match"
	NotificationTypeMessage NotificationType = "message"
	NotificationTypeLike    NotificationType = "like"
)

type NotificationPayload map[string]interface{}

type Notification struct{
	ID        uuid.UUID           `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    uuid.UUID           `gorm:"type:uuid;index;not null" json:"user_id"`
	Type      NotificationType    `gorm:"type:varchar(50);not null" json:"type"`
	Payload   NotificationPayload `gorm:"type:jsonb;serializer:json" json:"payload"`
	IsRead    bool                `gorm:"default:false" json:"is_read"`
	CreatedAt time.Time           `gorm:"index" json:"created_at"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}

func (n *Notification) BeforeCreate(tx *gorm.DB) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	return nil
}
