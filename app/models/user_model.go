package models

import (
	"time"

	"github.com/google/uuid"
)

// User struct to describe User object.
type User struct {
	ID           uuid.UUID `db:"id" json:"-" validate:"required,uuid"`
	CreatedAt    time.Time `db:"created_at" json:"-"`
	UpdatedAt    time.Time `db:"updated_at" json:"-"`
	Email        string    `db:"email" json:"email" validate:"required,email,lte=255"`
	PasswordHash string    `db:"password_hash" json:"password_hash,omitempty" validate:"required,lte=255"`
	UserStatus   int       `db:"user_status" json:"-" validate:"required,len=1"`
	UserRole     string    `db:"user_role" json:"-" validate:"required,lte=25"`
	Name         string    `db:"name" json:"name" validate:"required,lte=25"`
	AccessToken  string    `json:"token"`
}
