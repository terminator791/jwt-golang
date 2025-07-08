package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// enum untuk status transaksi
type TransactionStatus string

const (
	TransactionStatusCheckin    TransactionStatus = "CHECKIN"
	TransactionStatusCheckout   TransactionStatus = "CHECKOUT"
	TransactionStatusIncomplete TransactionStatus = "INCOMPLETE"
	TransactionStatusCancelled  TransactionStatus = "CANCELLED"
)

// enum untuk tipe transaksi
type TransactionType string

const (
	TransactionTypeRegular TransactionType = "REGULAR"
	TransactionTypePenalty TransactionType = "PENALTY"
	TransactionTypeRefund  TransactionType = "REFUND"
)

// Transaction merepresentasikan transaksi perjalanan dalam sistem e-ticketing
type Transaction struct {
	TransactionID         uuid.UUID         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"transaction_id"`
	CardID                uuid.UUID         `gorm:"type:uuid;not null;index" json:"card_id"`
	OriginTerminalID      *int              `gorm:"index" json:"origin_terminal_id"`
	DestinationTerminalID *int              `gorm:"index" json:"destination_terminal_id"`
	CheckinGateID         *uuid.UUID        `gorm:"type:uuid;index" json:"checkin_gate_id"`
	CheckoutGateID        *uuid.UUID        `gorm:"type:uuid;index" json:"checkout_gate_id"`
	CheckinTime           time.Time         `gorm:"not null" json:"checkin_time"`
	CheckoutTime          *time.Time        `json:"checkout_time"`
	FareAmount            *float64          `gorm:"type:decimal(10,2)" json:"fare_amount"`
	BalanceBefore         float64           `gorm:"type:decimal(12,2);not null" json:"balance_before"`
	BalanceAfter          *float64          `gorm:"type:decimal(12,2)" json:"balance_after"`
	TransactionStatus     TransactionStatus `gorm:"type:varchar(20);not null" json:"transaction_status"`
	TransactionType       TransactionType   `gorm:"type:varchar(20);not null;default:'REGULAR'" json:"transaction_type"`
	ReferenceNumber       string            `gorm:"type:varchar(50);uniqueIndex;not null" json:"reference_number"`
	CreatedAt             time.Time         `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt             time.Time         `gorm:"autoUpdateTime" json:"updated_at"`
	IsSynced              bool              `gorm:"not null;default:false" json:"is_synced"`
}

// BeforeCreate - Set UUID sebelum membuat transaksi baru
func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	if t.TransactionID == uuid.Nil {
		t.TransactionID = uuid.New()
	}
	return nil
}
