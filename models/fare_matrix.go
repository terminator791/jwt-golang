package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// FareMatrix merepresentasikan matriks tarif antar terminal
type FareMatrix struct {
	FareID             uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"fare_id"`
	FromTerminalID     int       `gorm:"index;not null" json:"from_terminal_id"`
	ToTerminalID       int       `gorm:"index;not null" json:"to_terminal_id"`
	BaseFare           float64   `gorm:"type:decimal(10,2);not null" json:"base_fare"`
	PeakHourMultiplier float64   `gorm:"type:decimal(4,2);default:1.0" json:"peak_hour_multiplier"`
	EffectiveFrom      time.Time `gorm:"not null" json:"effective_from"`
	EffectiveUntil     time.Time `json:"effective_until"`
	IsActive           bool      `gorm:"not null;default:true" json:"is_active"`
	CreatedAt          time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// BeforeCreate - Set UUID sebelum membuat fare matrix baru
func (f *FareMatrix) BeforeCreate(tx *gorm.DB) error {
	if f.FareID == uuid.Nil {
		f.FareID = uuid.New()
	}
	return nil
}
