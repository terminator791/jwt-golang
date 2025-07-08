package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// enum untuk metode pembayaran
type PaymentMethod string

const (
	PaymentMethodCash       PaymentMethod = "CASH"
	PaymentMethodCreditCard PaymentMethod = "CREDIT_CARD"
	PaymentMethodDebitCard  PaymentMethod = "DEBIT_CARD"
	PaymentMethodEWallet    PaymentMethod = "E_WALLET"
	PaymentMethodTransfer   PaymentMethod = "BANK_TRANSFER"
)

// enum untuk status top up
type TopUpStatus string

const (
	TopUpStatusPending   TopUpStatus = "PENDING"
	TopUpStatusCompleted TopUpStatus = "COMPLETED"
	TopUpStatusFailed    TopUpStatus = "FAILED"
	TopUpStatusRefunded  TopUpStatus = "REFUNDED"
)

// TopUp merepresentasikan transaksi pengisian saldo kartu
type TopUp struct {
	TopUpID          uuid.UUID     `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"topup_id"`
	CardID           uuid.UUID     `gorm:"type:uuid;not null;index" json:"card_id"`
	Amount           float64       `gorm:"type:decimal(12,2);not null" json:"amount"`
	PaymentMethod    PaymentMethod `gorm:"type:varchar(20);not null" json:"payment_method"`
	PaymentReference string        `gorm:"type:varchar(100)" json:"payment_reference"`
	Status           TopUpStatus   `gorm:"type:varchar(10);not null" json:"status"`
	ProcessedAt      *time.Time    `json:"processed_at"`
	CreatedAt        time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
}

// BeforeCreate - Set UUID sebelum membuat top up baru
func (t *TopUp) BeforeCreate(tx *gorm.DB) error {
	if t.TopUpID == uuid.Nil {
		t.TopUpID = uuid.New()
	}
	return nil
}
