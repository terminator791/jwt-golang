package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SyncType adalah enum untuk tipe sinkronisasi
type SyncType string

const (
	SyncTypeCardData      SyncType = "CARD_DATA"
	SyncTypeTransactions  SyncType = "TRANSACTIONS"
	SyncTypeFareMatrix    SyncType = "FARE_MATRIX"
	SyncTypeConfiguration SyncType = "CONFIGURATION"
)

// SyncStatus adalah enum untuk status sinkronisasi
type SyncStatus string

const (
	SyncStatusPending    SyncStatus = "PENDING"
	SyncStatusInProgress SyncStatus = "IN_PROGRESS"
	SyncStatusCompleted  SyncStatus = "COMPLETED"
	SyncStatusFailed     SyncStatus = "FAILED"
	SyncStatusRetrying   SyncStatus = "RETRYING"
)

// JSONData adalah tipe untuk menyimpan data JSON di database
type JSONData map[string]interface{}

// Value mengkonversi JSONData ke format yang dapat disimpan di database
func (j JSONData) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan mengkonversi data dari database ke JSONData
func (j *JSONData) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("tipe data tidak valid untuk JSONData")
	}

	return json.Unmarshal(bytes, &j)
}

// SyncLog merepresentasikan log sinkronisasi data antar terminal dan sistem pusat
type SyncLog struct {
	SyncID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"sync_id"`
	TerminalID      int        `gorm:"index;not null" json:"terminal_id"`
	SyncType        SyncType   `gorm:"type:varchar(20);not null" json:"sync_type"`
	SyncData        JSONData   `gorm:"type:json" json:"sync_data"`
	SyncStatus      SyncStatus `gorm:"type:varchar(15);not null" json:"sync_status"`
	SyncStartedAt   time.Time  `gorm:"not null" json:"sync_started_at"`
	SyncCompletedAt *time.Time `json:"sync_completed_at"`
	ErrorMessage    string     `gorm:"type:varchar(255)" json:"error_message"`
	RetryCount      int        `gorm:"not null;default:0" json:"retry_count"`
}

// BeforeCreate - Set UUID sebelum membuat sync log baru
func (s *SyncLog) BeforeCreate(tx *gorm.DB) error {
	if s.SyncID == uuid.Nil {
		s.SyncID = uuid.New()
	}
	if s.SyncStartedAt.IsZero() {
		s.SyncStartedAt = time.Now()
	}
	return nil
}
