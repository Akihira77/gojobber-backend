package types

import (
	"database/sql"
	"mime/multipart"
	"time"
)

type Auth struct {
	ID                     string         `json:"id" gorm:"primaryKey;not null"`
	Username               string         `json:"username" gorm:"index:idx_username,unique;not null"`
	Password               string         `json:"password" gorm:"not null"`
	Email                  string         `json:"email" gorm:"index:idx_email,unique;not null"`
	ProfilePublicID        string         `json:"profilePublicId" gorm:"not null"`
	Country                string         `json:"country" gorm:"not null"`
	ProfilePicture         string         `json:"profilePicture" gorm:"not null"`
	EmailVerificationToken sql.NullString `json:"emailVerificationToken"`
	EmailVerified          bool           `json:"emailVerified" gorm:"default:false;not null"`
	CreatedAt              time.Time      `json:"createdAt" gorm:"not null"`
	PasswordResetExpires   *time.Time     `json:"passwordResetExpires" gorm:"default:null;"`
	PasswordResetToken     sql.NullString `json:"passwordResetToken" gorm:"default:null;"`
}

type AuthExcludePassword struct {
	ID                     string     `json:"id"`
	Username               string     `json:"username"`
	Email                  string     `json:"email"`
	ProfilePublicID        string     `json:"profilePublicId,omitempty"`
	Country                string     `json:"country"`
	ProfilePicture         string     `json:"profilePicture"`
	EmailVerificationToken string     `json:"emailVerificationToken,omitempty"`
	EmailVerified          bool       `json:"emailVerified"`
	CreatedAt              *time.Time `json:"createdAt,omitempty"`
	PasswordResetExpires   *time.Time `json:"passwordResetExpires,omitempty"`
	PasswordResetToken     string     `json:"passwordResetToken,omitempty"`
}

type SignIn struct {
	Username string `json:"username" validate:"required,alphanum"`
	Password string `json:"password" validate:"required,alphanum"`
}

type SignUp struct {
	Username        string         `json:"username" form:"username" validate:"required,alphanum"`
	Password        string         `json:"password" form:"password" validate:"required,min=8,max=16,alphanum"`
	Country         string         `json:"country" form:"country" validate:"required,alpha"`
	Email           string         `json:"email" form:"email" validate:"required,email"`
	ProfilePicture  string         `json:"profilePicture"`
	Avatar          multipart.File `json:"avatar" form:"avatar"`
	ProfilePublicID string         `json:"profilePublicId"`
}

type ResetPassword struct {
	Password        string `json:"password" validate:"required,min=8,max=16,alphanum"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,min=8,max=16,alphanum"`
}

type ChangePassword struct {
	CurrentPassword string `json:"currentPassword" validate:"required,min=8,max=16,alphanum"`
	NewPassword     string `json:"newPassword" validate:"required,min=8,max=16,alphanum"`
}
