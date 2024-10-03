package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Akihira77/gojobber/services/3-auth/types"
	"github.com/Akihira77/gojobber/services/3-auth/util"
	"github.com/Akihira77/gojobber/services/common/genproto/user"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type AuthServiceImpl interface {
	FindUserByUsernameOrEmail(ctx context.Context, str string) (*types.AuthExcludePassword, error)
	FindUserByUsernameOrEmailIncPassword(ctx context.Context, str string) (*types.Auth, error)
	Create(ctx context.Context, data *types.SignUp, userGrpcClient user.UserServiceClient) (*types.AuthExcludePassword, error)
	FindUserByID(ctx context.Context, id string) (*types.AuthExcludePassword, error)
	FindUserByIDIncPassword(ctx context.Context, id string) (*types.Auth, error)
	FindUserByVerificationToken(ctx context.Context, token string) (*types.AuthExcludePassword, error)
	FindUserByPasswordToken(ctx context.Context, token string) (*types.AuthExcludePassword, error)
	UpdateEmailVerification(ctx context.Context, userId string, emailStatus bool, emailVerifToken ...string) (*types.AuthExcludePassword, error)
	UpdatePasswordToken(ctx context.Context, userId string, token string, tokenExpiration time.Time) error
	UpdatePassword(ctx context.Context, userId string, password string) error
}

type AuthService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) AuthServiceImpl {
	return &AuthService{
		db: db,
	}
}

func (as *AuthService) FindUserByUsernameOrEmail(ctx context.Context, str string) (*types.AuthExcludePassword, error) {
	var user types.AuthExcludePassword
	result := as.db.WithContext(ctx).
		Model(&types.Auth{}).
		First(&user, "username = ? OR email = ?", str, str)
	return &user, result.Error
}

func (as *AuthService) FindUserByUsernameOrEmailIncPassword(ctx context.Context, str string) (*types.Auth, error) {
	var user types.Auth
	result := as.db.WithContext(ctx).
		Model(&types.Auth{}).
		First(&user, "username = ? OR email = ?", str, str)
	return &user, result.Error
}

func (as *AuthService) Create(ctx context.Context, u *types.SignUp, userGrpcClient user.UserServiceClient) (*types.AuthExcludePassword, error) {
	hashPass, err := util.HashPassword(u.Password)
	if err != nil {
		return nil, fmt.Errorf("hashing password error: %v", err)
	}

	auth := &types.Auth{
		ID:              util.RandomStr(64),
		Username:        u.Username,
		Email:           u.Email,
		Password:        hashPass,
		Country:         u.Country,
		ProfilePicture:  u.ProfilePicture,
		ProfilePublicID: u.ProfilePublicID,
		CreatedAt:       time.Now(),
	}

	tx := as.db.WithContext(ctx).
		Model(&types.Auth{}).
		Begin()

	if err := tx.Create(auth).Error; err != nil {
		tx.Rollback()
		return &types.AuthExcludePassword{}, err
	}

	_, err = userGrpcClient.SaveBuyerData(ctx, &user.SaveBuyerRequest{
		Id:             auth.ID,
		Username:       auth.Username,
		Email:          auth.Email,
		Country:        auth.Country,
		ProfilePicture: auth.ProfilePicture,
		IsSeller:       false,
		CreatedAt:      timestamppb.New(auth.CreatedAt),
	})

	if err != nil {
		tx.Rollback()
		return &types.AuthExcludePassword{}, err
	}

	return &types.AuthExcludePassword{
		ID:              auth.ID,
		Username:        auth.Username,
		Email:           auth.Email,
		Country:         auth.Country,
		ProfilePicture:  auth.ProfilePicture,
		ProfilePublicID: auth.ProfilePublicID,
	}, tx.Commit().Error
}

func (as *AuthService) FindUserByID(ctx context.Context, id string) (*types.AuthExcludePassword, error) {
	var user types.AuthExcludePassword
	result := as.db.WithContext(ctx).
		Model(&types.Auth{}).
		First(&user, "id = ?", id)
	return &user, result.Error
}

func (as *AuthService) FindUserByIDIncPassword(ctx context.Context, id string) (*types.Auth, error) {
	var user types.Auth
	result := as.db.WithContext(ctx).
		Model(&types.Auth{}).
		First(&user, "id = ?", id)
	return &user, result.Error
}

func (as *AuthService) FindUserByVerificationToken(ctx context.Context, token string) (*types.AuthExcludePassword, error) {
	var user types.AuthExcludePassword
	result := as.db.WithContext(ctx).
		Model(&types.Auth{}).
		First(&user, "emailVerificationToken = ?", token)
	return &user, result.Error
}

func (as *AuthService) FindUserByPasswordToken(ctx context.Context, token string) (*types.AuthExcludePassword, error) {
	var user types.AuthExcludePassword
	result := as.db.WithContext(ctx).
		Model(&types.Auth{}).
		First(&user,
			"passwordResetToken = ? OR passwordResetExpires > ?", token, time.Now().String())
	return &user, result.Error
}

func (as *AuthService) UpdateEmailVerification(ctx context.Context,
	userId string, emailStatus bool, emailVerifToken ...string) (*types.AuthExcludePassword, error) {
	var user types.AuthExcludePassword
	var result *gorm.DB
	if len(emailVerifToken) == 0 {
		result = as.db.WithContext(ctx).
			Model(&types.Auth{}).
			Where("id = ?", userId).
			Updates(types.Auth{EmailVerificationToken: util.NewNullString(""), EmailVerified: emailStatus})
	} else if len(emailVerifToken) == 1 {
		result = as.db.WithContext(ctx).
			Model(&types.Auth{}).
			Where("id = ?", userId).
			Updates(types.Auth{EmailVerificationToken: util.NewNullString(emailVerifToken[0]), EmailVerified: emailStatus})
	} else {
		err := fmt.Errorf("BUG!. email verification token is too many")
		return nil, err
	}

	if result.Error != nil {
		return nil, result.Error
	}

	result = result.Model(&types.Auth{}).
		First(&user, "id = ?", userId)

	return &user, result.Error
}

func (as *AuthService) UpdatePasswordToken(ctx context.Context, userId, token string, tokenExpiration time.Time) error {
	result := as.db.WithContext(ctx).
		Model(&types.Auth{}).
		Where("id = ?", userId).
		Updates(types.Auth{PasswordResetToken: util.NewNullString(token), PasswordResetExpires: &tokenExpiration})

	return result.Error
}

func (as *AuthService) UpdatePassword(ctx context.Context, userId string, password string) error {
	now := time.Now()
	result := as.db.WithContext(ctx).
		Model(&types.Auth{}).
		Where("id = ?", userId).
		Updates(types.Auth{PasswordResetToken: util.NewNullString(""), PasswordResetExpires: &now, Password: password})

	return result.Error
}
