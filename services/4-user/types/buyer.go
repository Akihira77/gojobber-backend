package types

import "time"

type Buyer struct {
	ID             string    `json:"id" gorm:"primaryKey;not null"`
	Username       string    `json:"username" gorm:"not null" validate:"required"`
	Email          string    `json:"email" gorm:"not null" validate:"required,email"`
	Country        string    `json:"country" gorm:"not null" validate:"required"`
	ProfilePicture string    `json:"profilePicture" gorm:"not null" validate:"required"`
	IsSeller       bool      `json:"isSeller" gorm:"not null; default:false;"`
	CreatedAt      time.Time `json:"createdAt" gorm:"not null" validate:"required"`
}
