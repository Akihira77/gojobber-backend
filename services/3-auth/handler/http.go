package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	svc "github.com/Akihira77/gojobber/services/3-auth/service"
	"github.com/Akihira77/gojobber/services/3-auth/types"
	"github.com/Akihira77/gojobber/services/3-auth/util"
	"github.com/Akihira77/gojobber/services/common/genproto/notification"
	"github.com/Akihira77/gojobber/services/common/genproto/user"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AuthHttpHandler struct {
	validate   *validator.Validate
	authSvc    svc.AuthServiceImpl
	cld        *util.Cloudinary
	grpcClient *GRPCClients
}

func NewAuthHttpHandler(authSvc svc.AuthServiceImpl, cld *util.Cloudinary, grpcServices *GRPCClients) *AuthHttpHandler {
	return &AuthHttpHandler{
		validate:   validator.New(validator.WithRequiredStructEnabled()),
		authSvc:    authSvc,
		cld:        cld,
		grpcClient: grpcServices,
	}
}

func (ah *AuthHttpHandler) GetUserInfo(c *fiber.Ctx) error {
	userInfo, ok := c.UserContext().Value("current_user").(*types.JWTClaims)
	if !ok {
		return fiber.NewError(http.StatusBadRequest, "invalid data. Please re-signin")
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"user": userInfo,
	})
}

func (ah *AuthHttpHandler) SignIn(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 1*time.Second)
	defer cancel()

	data := new(types.SignIn)
	if err := c.BodyParser(data); err != nil {
		fmt.Printf("signin error: \n%+v", err)
		return fiber.NewError(http.StatusBadRequest, "invalid data. Please correct your data")
	}

	err := ah.validate.Struct(data)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors": util.CustomValidationErrors(err),
		})
	}

	u, err := ah.authSvc.FindUserByUsernameOrEmailIncPassword(ctx, data.Username)
	if err != nil {
		fmt.Printf("signin error: \n%+v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(http.StatusNotFound, "user did not found")
		}
		return fiber.NewError(http.StatusBadRequest, "signin failed")
	}

	err = util.CheckPasswordHash(data.Password, u.Password)
	if err != nil {
		fmt.Printf("signin error: \n%+v", err)
		return fiber.NewError(http.StatusBadRequest, "password did not matched")
	}

	token, err := util.SigningJWT(os.Getenv("JWT_SECRET"), u.ID, u.Email, u.Username)
	if err != nil {
		fmt.Printf("signin error: \n%+v", err)
		return fiber.NewError(http.StatusBadRequest, "signin failed")
	}

	c.Cookie(&fiber.Cookie{
		Name:    "token",
		Value:   token,
		Expires: time.Now().Add(1 * time.Hour),
	})
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"user": types.AuthExcludePassword{
			ID:             u.ID,
			Username:       u.Username,
			Email:          u.Email,
			Country:        u.Country,
			ProfilePicture: u.ProfilePicture,
			EmailVerified:  u.EmailVerified,
			CreatedAt:      &u.CreatedAt,
		},
		"token": token,
	})
}

func (ah *AuthHttpHandler) SignUp(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 10*time.Second)
	defer cancel()

	data := new(types.SignUp)
	if err := c.BodyParser(data); err != nil {
		log.Printf("signup error:\n%+v", err)
		return fiber.NewError(http.StatusBadRequest, "invalid data. Please correct your data")
	}

	err := ah.validate.Struct(data)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors": util.CustomValidationErrors(err),
		})
	}

	formHeader, err := c.FormFile("avatar")
	if err != nil {
		log.Printf("signup error:\n%+v", err)
		return fiber.NewError(http.StatusBadRequest, "failed reading avatar file")
	}

	if formHeader.Size > 2*1024*1024 {
		log.Printf("signup error. File is too large")
		return fiber.NewError(http.StatusBadRequest, "file is larger than 2MB")
	}

	if !util.ValidateImgExtension(formHeader) {
		log.Println("signup error file type is unsupported")
		return fiber.NewError(http.StatusBadRequest, "file type is unsupported")
	}

	u, err := ah.authSvc.FindUserByUsernameOrEmail(ctx, data.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("signup error:\n%+v", err)
		return fiber.ErrInternalServerError
	}

	if u.ID != "" {
		return fiber.NewError(http.StatusBadRequest, "user already exists")
	}

	data.File, err = formHeader.Open()
	if err != nil {
		log.Printf("signup error:\n%+v", err)
		return fiber.NewError(http.StatusBadRequest, "failed reading avatar file")
	}

	str := util.RandomStr(32)
	uploadResult, err := ah.cld.UploadImg(ctx, data.File, str)
	if err != nil {
		log.Printf("signup error:\n%+v", err)
		return fiber.NewError(http.StatusBadRequest, "failed upload file")
	}

	data.ProfilePicture = uploadResult.SecureURL
	data.ProfilePublicID = uploadResult.PublicID
	cc, err := ah.grpcClient.GetClient("USER_SERVICE")
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Error while searching gig")
	}

	userGrpcClient := user.NewUserServiceClient(cc)
	result, err := ah.authSvc.Create(ctx, data, userGrpcClient)
	if err != nil {
		log.Printf("signup error:\n%+v", err)
		ah.cld.Destroy(context.Background(), uploadResult.PublicID)
		return fiber.NewError(http.StatusBadRequest, "signup failed")
	}

	token, err := util.SigningJWT(os.Getenv("JWT_SECRET"), result.ID, result.Email, result.Username)
	if err != nil {
		log.Printf("signup error:\n%+v", err)
		return fiber.NewError(http.StatusBadRequest, "error generating JWT")
	}

	c.Cookie(&fiber.Cookie{
		Name:    "token",
		Value:   token,
		Expires: time.Now().Add(1 * time.Hour),
	})
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"user": types.AuthExcludePassword{
			ID:             result.ID,
			Username:       result.Username,
			Email:          result.Email,
			Country:        result.Country,
			ProfilePicture: result.ProfilePicture,
			EmailVerified:  result.EmailVerified,
			CreatedAt:      result.CreatedAt,
		},
		"token": token,
	})
}

func (ah *AuthHttpHandler) RefreshToken(c *fiber.Ctx) error {
	userInfo, ok := c.UserContext().Value("current_user").(*types.JWTClaims)
	if !ok {
		return fiber.NewError(http.StatusUnauthorized, "sign in first")
	}

	token, err := util.SigningJWT(os.Getenv("JWT_SECRET"), userInfo.UserID, userInfo.Email, userInfo.Username)
	if err != nil {
		log.Printf("refreshtoken:\n%+v", err)
		return fiber.NewError(http.StatusBadRequest, "failed refresh the token")
	}

	c.Cookie(&fiber.Cookie{
		Name:    "token",
		Value:   token,
		Expires: time.Now().Add(1 * time.Hour),
	})
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"token": token,
		"user":  userInfo,
	})
}

func (ah *AuthHttpHandler) VerifyEmail(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 1*time.Second)
	defer cancel()

	token := c.Params("token", "")
	userInfo, ok := c.UserContext().Value("current_user").(*types.JWTClaims)
	if !ok {
		return fiber.NewError(http.StatusBadRequest, "invalid data. Please re-signin")
	}

	_, err := ah.authSvc.FindUserByVerificationToken(ctx, token)
	if err != nil {
		log.Printf("verifyemail error:\n%+v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(http.StatusNotFound, "user did not found")
		}
		return fiber.ErrInternalServerError
	}

	result, err := ah.authSvc.UpdateEmailVerification(ctx, userInfo.UserID, true, "")
	if err != nil {
		log.Printf("verifyemail error:\n%+v", err)
		return fiber.ErrInternalServerError
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"user": result,
	})
}

func (ah *AuthHttpHandler) SendVerifyEmailURL(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	userInfo, ok := c.UserContext().Value("current_user").(*types.JWTClaims)
	if !ok {
		return fiber.NewError(http.StatusBadRequest, "invalid data. Please re-signin")
	}

	cc, err := ah.grpcClient.GetClient("NOTIFICATION_SERVICE")
	if err != nil {
		log.Printf("sendverifyemail error:\n%+v", err)
		return fiber.NewError(http.StatusInternalServerError, "Error sending email")
	}

	randStr := util.RandomStr(64)
	_, err = ah.authSvc.UpdateEmailVerification(ctx, userInfo.UserID, false, randStr)
	if err != nil {
		log.Printf("sendverifyemail error:\n%+v", err)
		return fiber.NewError(http.StatusInternalServerError, "Error sending email")
	}

	go func() {
		verifURL := fmt.Sprintf("%s/confirm_email?v_token=%s", os.Getenv("CLIENT_URL"), randStr)
		notificationGrpcClient := notification.NewNotificationServiceClient(cc)
		notificationGrpcClient.UserVerifyingEmail(context.TODO(), &notification.VerifyingEmailRequest{
			ReceiverEmail:    userInfo.Email,
			HtmlTemplateName: "verifyEmail",
			VerifyLink:       verifURL,
		})
	}()

	// if err != nil {
	// 	fmt.Printf("sendforgotpasswordurl error:\n%+v", err)
	// 	return fiber.NewError(http.StatusInternalServerError, "Error sending verify email URL to your email")
	// }

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "verify email URL has been send to your email",
	})
}

func (ah *AuthHttpHandler) SendForgotPasswordURL(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	email := c.Params("email", "")
	user, err := ah.authSvc.FindUserByUsernameOrEmail(ctx, email)
	if err != nil {
		fmt.Printf("sendforgotpasswordurl error:\n%+v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(http.StatusNotFound, "user did not found")
		}
		return fiber.ErrInternalServerError
	}

	cc, err := ah.grpcClient.GetClient("NOTIFICATION_SERVICE")
	if err != nil {
		fmt.Printf("sendforgotpasswordurl error:\n%+v", err)
		return fiber.NewError(http.StatusInternalServerError, "Unexpected error happened. Please try again.")
	}

	randStr := util.RandomStr(32)
	err = ah.authSvc.UpdatePasswordToken(ctx, user.ID, randStr, time.Now().Add(1*time.Hour))
	if err != nil {
		fmt.Printf("sendforgotpasswordurl error:\n%+v", err)
		return fiber.NewError(http.StatusInternalServerError, "Unexpected error happened. Please try again.")
	}

	go func() {
		resetPassURL := fmt.Sprintf("%s/reset-password?token=%s", os.Getenv("CLIENT_URL"), randStr)
		notificationGrpcClient := notification.NewNotificationServiceClient(cc)
		notificationGrpcClient.UserForgotPassword(context.TODO(), &notification.ForgotPasswordRequest{
			ReceiverEmail:    user.Email,
			HtmlTemplateName: "resetPassword",
			Username:         user.Username,
			ResetLink:        resetPassURL,
		})
	}()

	// if err != nil {
	// 	fmt.Printf("sendforgotpasswordurl error:\n%+v", err)
	// 	return fiber.NewError(http.StatusInternalServerError, "Error sending forgot password URL to your email")
	// }

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "email has been sent to your email",
	})
}

func (ah *AuthHttpHandler) ResetPassword(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 5*time.Second)
	defer cancel()

	token := c.Params("token", "")
	obj := new(types.ResetPassword)

	if err := c.BodyParser(obj); err != nil {
		fmt.Printf("resetpassword error:\n%+v", err)
		return fiber.NewError(http.StatusBadRequest, "invalid data")
	}

	err := ah.validate.Struct(obj)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors": util.CustomValidationErrors(err),
		})
	}

	if obj.Password != obj.ConfirmPassword {
		return fiber.NewError(http.StatusBadRequest, "password not matched with confirm password")
	}

	user, err := ah.authSvc.FindUserByPasswordToken(ctx, token)
	if err != nil {
		fmt.Printf("resetpassword error:\n%+v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(http.StatusNotFound, "user did not found")
		}
		return fiber.ErrInternalServerError
	}

	g := user.PasswordResetExpires.After(time.Now())
	if !g {
		return fiber.NewError(http.StatusBadRequest, "reset password token already expired")
	}

	hashedPass, err := util.HashPassword(obj.Password)
	if err != nil {
		fmt.Printf("resetpasswordsuccess error:\n%+v", err)
		return fiber.ErrInternalServerError
	}

	err = ah.authSvc.UpdatePassword(ctx, user.ID, hashedPass)
	if err != nil {
		fmt.Printf("resetpasswordsuccess error:\n%+v", err)
		return fiber.ErrInternalServerError
	}

	cc, err := ah.grpcClient.GetClient("NOTIFICATION_SERVICE")
	if err != nil {
		fmt.Printf("resetpasswordsuccess error:\n%+v", err)
		return fiber.NewError(http.StatusInternalServerError, "Unexpected error happened. Please try again.")
	}

	go func() {
		notificationGrpcClient := notification.NewNotificationServiceClient(cc)
		notificationGrpcClient.UserSucessResetPassword(context.TODO(), &notification.SuccessResetPasswordRequest{
			ReceiverEmail:    user.Email,
			HtmlTemplateName: "resetPasswordSuccess",
			Username:         user.Username,
		})

	}()

	// if err != nil {
	// 	fmt.Printf("resetpasswordsuccess error:\n%+v", err)
	// 	return fiber.NewError(http.StatusInternalServerError, "Error sending email")
	// }

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "password reseted",
	})
}

func (ah *AuthHttpHandler) ChangePassword(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 2*time.Second)
	defer cancel()

	obj := new(types.ChangePassword)
	if err := c.BodyParser(obj); err != nil {
		fmt.Printf("changepassword error:\n%+v", err)
		return fiber.NewError(http.StatusBadRequest, "invalid data")
	}

	err := ah.validate.Struct(obj)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"errors": util.CustomValidationErrors(err),
		})
	}

	userInfo, ok := c.UserContext().Value("current_user").(*types.JWTClaims)
	if !ok {
		return fiber.NewError(http.StatusUnauthorized, "sign in first")
	}

	user, err := ah.authSvc.FindUserByIDIncPassword(ctx, userInfo.UserID)
	if err != nil {
		fmt.Printf("changepassword error:\n%+v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(http.StatusNotFound, "user did not found")
		}
		return fiber.ErrInternalServerError
	}

	err = util.CheckPasswordHash(obj.CurrentPassword, user.Password)
	if err != nil {
		fmt.Printf("changepassword error:\n%+v", err)
		return fiber.NewError(http.StatusBadRequest, "password did not matched")
	}

	hashedPass, err := util.HashPassword(obj.NewPassword)
	if err != nil {
		fmt.Printf("changepassword error:\n%+v", err)
		return fiber.NewError(http.StatusBadRequest, "failed storing password")
	}

	err = ah.authSvc.UpdatePassword(ctx, user.ID, hashedPass)
	if err != nil {
		fmt.Printf("changepassword error:\n%+v", err)
		return fiber.ErrInternalServerError
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "password changed",
	})
}
