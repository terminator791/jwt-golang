package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// GateType adalah enum untuk tipe gate
type GateType string

const (
	GateTypeEntry GateType = "ENTRY"
	GateTypeExit  GateType = "EXIT"
	GateTypeBoth  GateType = "BOTH"
)

// GateStatus adalah enum untuk status gate
type GateStatus string

const (
	GateStatusActive      GateStatus = "ACTIVE"
	GateStatusInactive    GateStatus = "INACTIVE"
	GateStatusMaintenance GateStatus = "MAINTENANCE"
)

// Gate merepresentasikan pintu masuk/keluar di terminal
type Gate struct {
	GateID        uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"gate_id"`
	TerminalID    int        `gorm:"index;not null" json:"terminal_id"`
	GateCode      string     `gorm:"type:varchar(20);uniqueIndex;not null" json:"gate_code"`
	GateType      GateType   `gorm:"type:varchar(10);not null" json:"gate_type"`
	Status        GateStatus `gorm:"type:varchar(20);not null;default:'ACTIVE'" json:"status"`
	IPAddress     string     `gorm:"type:varchar(45)" json:"ip_address"`
	LastHeartbeat time.Time  `json:"last_heartbeat"`
	IsOnline      bool       `gorm:"default:true" json:"is_online"`
	CreatedAt     time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

// BeforeCreate - Set UUID sebelum membuat gate baru
func (g *Gate) BeforeCreate(tx *gorm.DB) error {
	if g.GateID == uuid.Nil {
		g.GateID = uuid.New()
	}
	return nil
}
