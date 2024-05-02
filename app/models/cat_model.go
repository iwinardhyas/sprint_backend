package models

import (
	"time"

	"github.com/google/uuid"
)

type Cat struct {
	ID     uuid.UUID `json:"id" db:"id"`
	UserID uuid.UUID `json:"user_id" db:"user_id"`
	NewCat
	HasMatched bool      `json:"hasMatched" db:"has_matched"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt  time.Time `json:"updatedAt" db:"updated_at"`
	DeletedAt  time.Time `db:"deleted_at"`
}

type NewCat struct {
	Name        string   `json:"name" db:"name" validate:"required,min=1,max=30"`
	Race        string   `json:"race" db:"race" validate:"required,oneof='Persian' 'Maine Coon' 'Siamese' 'Ragdoll' 'Bengal' 'Sphynx' 'British Short hair' 'Abyssinian' 'Scottish Fold' 'Birman'"`
	Sex         string   `json:"sex" db:"sex" validate:"required,oneof=male female"`
	AgeInMonth  int      `json:"ageInMonth" db:"age" validate:"required,min=1,max=120082"`
	Description string   `json:"description" db:"desc" validate:"required,min=1,max=200"`
	ImageUrls   []string `json:"imageUrls" db:"img_link" validate:"required,min=1,dive,url"`
}

type CatData struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Race        string   `json:"race"`
	Sex         string   `json:"sex"`
	AgeInMonth  int      `json:"ageInMonth"`
	ImageUrls   []string `json:"imageUrls"`
	Description string   `json:"description"`
	HasMatched  bool     `json:"hasMatched"`
	CreatedAt   string   `json:"createdAt"`
}
