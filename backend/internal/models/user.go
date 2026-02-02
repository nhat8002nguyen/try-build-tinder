package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
	GenderOther  Gender = "other"
)

type UserPreferences struct {
	MinAge           int      `json:"min_age"`
	MaxAge           int      `json:"max_age"`
	MaxDistance      int      `json:"max_distance"` // in kilometers
	GenderPreference []Gender `json:"gender_preference"`
}

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type User struct {
	ID           uuid.UUID       `gorm:"type:uuid;primaryKey" json:"id"`
	Email        string          `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string          `gorm:"type:varchar(255)" json:"-"`
	Name         string          `gorm:"not null" json:"name"`
	Gender       Gender          `gorm:"type:varchar(20)" json:"gender"`
	Birthdate    *time.Time      `json:"birthdate"`
	Bio          string          `gorm:"type:text" json:"bio"`
	Location     Location        `gorm:"type:jsonb;serializer:json" json:"location"`
	Preferences  UserPreferences `gorm:"type:jsonb;serializer:json" json:"preferences"`
	IsVerified   bool            `gorm:"default:false" json:"is_verified"`
	IsActive     bool            `gorm:"default:true" json:"is_active"`
	LastActiveAt *time.Time      `gorm:"index" json:"last_active_at"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`

	Photos        []UserPhoto     `gorm:"foreignKey:UserID" json:"photos,omitempty"`
	OAuthAccounts []OAuthAccount  `gorm:"foreignKey:UserID" json:"-"`
	Notifications []Notification  `gorm:"foreignKey:UserID" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	now := time.Now()
	u.LastActiveAt = &now
	return nil
}

func (u *User) Age() int {
	if u.Birthdate == nil {
		return 0
	}
	now := time.Now()
	age := now.Year() - u.Birthdate.Year()
	if now.YearDay() < u.Birthdate.YearDay() {
		age--
	}
	return age
}

type UserPhoto struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID       uuid.UUID `gorm:"type:uuid;index;not null" json:"user_id"`
	PhotoURL     string    `gorm:"not null" json:"photo_url"`
	DisplayOrder int       `gorm:"default:0" json:"display_order"`
	IsApproved   bool      `gorm:"default:true" json:"is_approved"`
	CreatedAt    time.Time `json:"created_at"`
}

func (p *UserPhoto) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}

type OAuthAccount struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID         uuid.UUID `gorm:"type:uuid;index;not null" json:"user_id"`
	Provider       string    `gorm:"not null" json:"provider"`
	ProviderUserID string    `gorm:"not null" json:"provider_user_id"`
	AccessToken    string    `json:"-"`
	RefreshToken   string    `json:"-"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}

func (o *OAuthAccount) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}
