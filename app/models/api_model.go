package models

import "github.com/google/uuid"

type Cat struct {
	ID          uuid.UUID `db:"id" json:"-" validate:"required,uuid"`
	Name        string    `json:"name" validate:"gt=0,omitempty,required,min=5,max=30"`
	Race        string    `json:"race" validate:"gt=0,omitempty"`
	Sex         string    `json:"sex" validate:"gt=0,omitempty"`
	AgeInMonth  string    `json:"age" validate:"gt=0,omitempty,required,min=1,max=120082"`
	Description string    `json:"desc" validate:"gt=0,omitempty,required,min=1,max=200"`
	ImageUrls   string    `json:"img_link" validate:"gt=0,omitempty,url"`
}
