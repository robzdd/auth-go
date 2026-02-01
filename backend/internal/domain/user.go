package domain

import (
	"time"
)

// User entity
type User struct {
	ID              uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name            string     `gorm:"type:varchar(255);not null;index" json:"name"`
	Email           string     `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password        string     `gorm:"type:varchar(255);not null" json:"-"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// PasswordResetToken entity
type PasswordResetToken struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	Email     string    `gorm:"index;not null" json:"email"`
	Token     string    `gorm:"index;not null" json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// UserRepository interface (Contract)
type UserRepository interface {
	Save(user *User) (*User, error)
	FindByEmail(email string) (*User, error)
	FindByID(id uint64) (*User, error)
	FindAll(page int, limit int, search string) ([]*User, int64, error)
	Update(user *User) (*User, error)
}

// PasswordResetRepository interface
type PasswordResetRepository interface {
	Save(reset *PasswordResetToken) (*PasswordResetToken, error)
	FindByToken(token string) (*PasswordResetToken, error)
	DeleteByEmail(email string) error
}
