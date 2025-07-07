package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserType adalah enum untuk tipe user
type UserType string

const (
	UserTypeAdmin    UserType = "ADMIN"
	UserTypeCustomer UserType = "CUSTOMER"
	UserTypeStaff    UserType = "STAFF"
)

// User merepresentasikan entitas pengguna sistem e-ticketing
type User struct {
	UserID      uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"user_id"`
	FullName    string    `gorm:"type:varchar(100);not null" json:"full_name"`
	Email       string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password    string    `gorm:"type:varchar(255);not null" json:"-"` // Password field ditambahkan untuk autentikasi
	Phone       string    `gorm:"type:varchar(20);uniqueIndex" json:"phone"`
	DateOfBirth time.Time `json:"date_of_birth"`
	UserType    UserType  `gorm:"type:varchar(20);not null;default:'CUSTOMER'" json:"user_type"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// BeforeCreate - Set UUID sebelum membuat user baru
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.UserID == uuid.Nil {
		u.UserID = uuid.New()
	}
	return nil
}
