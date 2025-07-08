package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// enum untuk status kartu
type CardStatus string

const (
	CardStatusActive  CardStatus = "ACTIVE"
	CardStatusBlocked CardStatus = "BLOCKED"
	CardStatusExpired CardStatus = "EXPIRED"
	CardStatusLost    CardStatus = "LOST"
)

// Card merepresentasikan kartu prepaid untuk e-ticketing
type Card struct {
	CardID         uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"card_id"`
	CardNumber     string     `gorm:"type:varchar(20);uniqueIndex;not null" json:"card_number"`
	UserID         *uuid.UUID `gorm:"type:uuid;index" json:"user_id"`
	CurrentBalance float64    `gorm:"type:decimal(12,2);not null;default:0" json:"current_balance"`
	Status         CardStatus `gorm:"type:varchar(10);not null;default:'ACTIVE'" json:"status"`
	IssuedAt       time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"issued_at"`
	ExpiresAt      time.Time  `gorm:"not null" json:"expires_at"`
	CreatedAt      time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	IsActive       bool       `gorm:"not null;default:true" json:"is_active"`
}

// BeforeCreate - Set UUID sebelum membuat card baru
func (c *Card) BeforeCreate(tx *gorm.DB) error {
	if c.CardID == uuid.Nil {
		c.CardID = uuid.New()
	}
	return nil
}
