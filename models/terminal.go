package models

import (
	"time"
)

// Terminal merepresentasikan terminal fisik dalam sistem transportasi
type Terminal struct {
	TerminalID   int       `gorm:"primaryKey;autoIncrement" json:"terminal_id"`
	TerminalName string    `gorm:"type:varchar(100);not null" json:"terminal_name"`
	TerminalCode string    `gorm:"type:varchar(20);uniqueIndex;not null" json:"terminal_code"`
	Location     string    `gorm:"type:varchar(255)" json:"location"`
	Latitude     float64   `gorm:"type:decimal(10,7)" json:"latitude"`
	Longitude    float64   `gorm:"type:decimal(10,7)" json:"longitude"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
