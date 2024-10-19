package types

import (
	"database/sql"
	"mime/multipart"
	"time"
)

// type Buyer struct {
// 	ID             string    `json:"id" gorm:"primaryKey;not null"`
// 	Username       string    `json:"username" gorm:"not null"`
// 	Email          string    `json:"email" gorm:"not null"`
// 	Country        string    `json:"country" gorm:"not null"`
// 	ProfilePicture string    `json:"profilePicture" gorm:"not null"`
// 	IsSeller       bool      `json:"isSeller" gorm:"not null; default:false"`
// 	CreatedAt      time.Time `json:"createdAt" gorm:"not null"`
// }

type Auth struct {
	ID                     string         `json:"id" gorm:"primaryKey;column:id;not null"`
	Username               string         `json:"username" gorm:"index:idx_username,unique;column:username;not null"`
	Password               string         `json:"password" gorm:"column:password;not null"`
	Email                  string         `json:"email" gorm:"index:idx_email,unique;column:email;not null"`
	ProfilePublicID        string         `json:"profilePublicId" gorm:"column:profile_public_id;not null"`
	Country                string         `json:"country" gorm:"column:country;not null"`
	ProfilePicture         string         `json:"profilePicture" gorm:"column:profile_picture;not null"`
	EmailVerificationToken sql.NullString `json:"emailVerificationToken"`
	EmailVerified          bool           `json:"emailVerified" gorm:"default:false;column:email_verified;not null"`
	CreatedAt              time.Time      `json:"createdAt" gorm:"column:created_at;not null"`
	PasswordResetExpires   *time.Time     `json:"passwordResetExpires" gorm:"default:null;column:password_reset_expires"`
	PasswordResetToken     sql.NullString `json:"passwordResetToken" gorm:"default:null;column:password_reset_token"`
}

type AuthExcludePassword struct {
	ID                     string     `json:"id"`
	Username               string     `json:"username"`
	Email                  string     `json:"email"`
	ProfilePublicID        string     `json:"profilePublicId,omitempty"`
	Country                string     `json:"country"`
	ProfilePicture         string     `json:"profilePicture"`
	EmailVerificationToken string     `json:"emailVerificationToken,omitempty" gorm:"column:email_verification_token;"`
	EmailVerified          bool       `json:"emailVerified"`
	CreatedAt              *time.Time `json:"createdAt,omitempty"`
	PasswordResetExpires   *time.Time `json:"passwordResetExpires,omitempty"`
	PasswordResetToken     string     `json:"passwordResetToken,omitempty"`
}

type SignIn struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type SignUp struct {
	Username        string `json:"username" form:"username" validate:"required"`
	Password        string `json:"password" form:"password" validate:"required,min=8,max=16"`
	Country         string `json:"country" form:"country" validate:"required"`
	Email           string `json:"email" form:"email" validate:"required,email"`
	ProfilePicture  string `json:"profilePicture"`
	File            multipart.File
	ProfilePublicID string `json:"profilePublicId"`
}

type ResetPassword struct {
	Password        string `json:"password" validate:"required,min=8,max=16"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,min=8,max=16"`
}

type ChangePassword struct {
	CurrentPassword string `json:"currentPassword" validate:"required,min=8,max=16"`
	NewPassword     string `json:"newPassword" validate:"required,min=8,max=16"`
}
