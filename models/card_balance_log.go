package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// enum untuk tipe perubahan saldo
type ChangeType string

const (
	ChangeTypeCheckin  ChangeType = "CHECKIN"
	ChangeTypeCheckout ChangeType = "CHECKOUT"
	ChangeTypeTopUp    ChangeType = "TOP_UP"
	ChangeTypeRefund   ChangeType = "REFUND"
	ChangeTypePenalty  ChangeType = "PENALTY"
	ChangeTypeAdmin    ChangeType = "ADMIN_ADJUSTMENT"
)

// CardBalanceLog merepresentasikan log perubahan saldo kartu
type CardBalanceLog struct {
	LogID           uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"log_id"`
	CardID          uuid.UUID  `gorm:"type:uuid;not null;index" json:"card_id"`
	PreviousBalance float64    `gorm:"type:decimal(12,2);not null" json:"previous_balance"`
	CurrentBalance  float64    `gorm:"type:decimal(12,2);not null" json:"current_balance"`
	AmountChanged   float64    `gorm:"type:decimal(12,2);not null" json:"amount_changed"`
	ChangeType      ChangeType `gorm:"type:varchar(20);not null" json:"change_type"`
	TransactionID   *uuid.UUID `gorm:"type:uuid;index" json:"transaction_id"`
	LoggedAt        time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"logged_at"`
	Description     string     `gorm:"type:varchar(255)" json:"description"`
}

// BeforeCreate - Set UUID sebelum membuat log baru
func (c *CardBalanceLog) BeforeCreate(tx *gorm.DB) error {
	if c.LogID == uuid.Nil {
		c.LogID = uuid.New()
	}
	if c.LoggedAt.IsZero() {
		c.LoggedAt = time.Now()
	}
	return nil
}
