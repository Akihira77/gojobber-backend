package types

import (
	"database/sql"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	jwt.RegisteredClaims
	UserID   string `json:"userId"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

// NOTE: Buyer types must be exactly the same as Buyer types on User Service
type Buyer struct {
	ID             string    `json:"id" gorm:"primaryKey;not null"`
	Username       string    `json:"username" gorm:"not null" validate:"required"`
	Email          string    `json:"email" gorm:"not null" validate:"required,email"`
	Country        string    `json:"country" gorm:"not null" validate:"required"`
	ProfilePicture string    `json:"profilePicture" gorm:"not null" validate:"required"`
	IsSeller       bool      `json:"isSeller" gorm:"not null; default:false;"`
	CreatedAt      time.Time `json:"createdAt" gorm:"not null" validate:"required"`
}

// NOTE: Auth types must be exactly the same as Auth types on Auth Service
type Auth struct {
	ID                     string         `json:"id" gorm:"primaryKey;column:id;not null"`
	Username               string         `json:"username" gorm:"index:idx_username,unique;column:username;not null"`
	Password               string         `json:"password" gorm:"column:password;not null"`
	Email                  string         `json:"email" gorm:"index:idx_email,unique;column:email;not null"`
	ProfilePublicID        string         `json:"profilePublicId" gorm:"column:profile_public_id;not null"`
	Country                string         `json:"country" gorm:"column:country;not null"`
	ProfilePicture         string         `json:"profilePicture" gorm:"column:profile_picture;not null"`
	EmailVerificationToken sql.NullString `json:"emailVerificationToken" gorm:"index:idx_email_verification_token,unique;column:email_verification_token"`
	EmailVerified          bool           `json:"emailVerified" gorm:"default:false;column:email_verified;not null"`
	CreatedAt              time.Time      `json:"createdAt" gorm:"column:created_at;not null"`
	PasswordResetExpires   *time.Time     `json:"passwordResetExpires" gorm:"default:null;column:password_reset_expires"`
	PasswordResetToken     sql.NullString `json:"passwordResetToken" gorm:"default:null;column:password_reset_token"`
}
