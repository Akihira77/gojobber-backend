package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Akihira77/gojobber/services/5-gig/service"
	"github.com/Akihira77/gojobber/services/5-gig/types"
	"github.com/Akihira77/gojobber/services/5-gig/util"
	"github.com/Akihira77/gojobber/services/common/genproto/user"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type GigHandler struct {
	gigSvc     service.GigServiceImpl
	validate   *validator.Validate
	cld        *util.Cloudinary
	grpcClient *GRPCClients
}

func NewGigHandler(gigSvc service.GigServiceImpl, cld *util.Cloudinary, grpcServices *GRPCClients) *GigHandler {
	return &GigHandler{
		gigSvc:     gigSvc,
		validate:   validator.New(validator.WithRequiredStructEnabled()),
		cld:        cld,
		grpcClient: grpcServices,
	}
}

func (gh *GigHandler) FindGigByID(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 1*time.Second)
	defer cancel()

	gig, err := gh.gigSvc.FindGigByID(ctx, c.Params("id"))
	if err != nil {
		log.Println("find gig by id", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(http.StatusNotFound, "Gig is not found")
		}
		return fiber.NewError(http.StatusInternalServerError, "Error while searching gig")
	}

	cc, err := gh.grpcClient.GetClient("USER_SERVICE")
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Error while searching gig")
	}

	userGrpcClient := user.NewUserServiceClient(cc)
	res, err := userGrpcClient.FindSeller(ctx, &user.FindSellerRequest{
		SellerId: gig.SellerID,
	})

	if err != nil {
		log.Println("find gig by id", err)
		return fiber.NewError(http.StatusInternalServerError, "Error while searching gig")
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"seller": &types.SellerOverview{
			SellerID:  gig.SellerID,
			FullName:  res.FullName,
			RatingSum: uint64(res.RatingSum),
			RatingCategories: types.RatingCategory{
				One:   uint(res.RatingCategories.One),
				Two:   uint(res.RatingCategories.Two),
				Three: uint(res.RatingCategories.Three),
				Four:  uint(res.RatingCategories.Four),
				Five:  uint(res.RatingCategories.Five),
			},
			RatingsCount: uint64(res.RatingsCount),
		},
		"gig": gig,
	})
}

func (gh *GigHandler) GigQuerySearch(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 1*time.Second)
	defer cancel()

	var p types.GigSearchParams
	err := c.ParamsParser(&p)
	if err != nil {
		log.Println("search-gig query:", err)
		return fiber.NewError(http.StatusBadRequest, "searching error")
	}

	var q types.GigSearchQuery
	err = c.QueryParser(&q)
	if err != nil {
		log.Println("search-gig query:", err)
		return fiber.NewError(http.StatusBadRequest, "searching error")
	}

	result, err := gh.gigSvc.GigQuerySearch(ctx, &p, &q)
	if err != nil {
		log.Println("gig query search:", err)
		return fiber.NewError(http.StatusBadRequest, "searching error")
	}

	cc, err := gh.grpcClient.GetClient("USER_SERVICE")
	if err != nil {
		log.Printf("gig query search:\n%+v", err)
		return fiber.NewError(http.StatusInternalServerError, "Error while validating seller")
	}
	userGrpcClient := user.NewUserServiceClient(cc)
	gigs, err := gh.gigSvc.FindAndMapSellerInGigs(ctx, userGrpcClient, result.Gigs)
	if err != nil {
		log.Println("gig query search", err)
		return fiber.NewError(http.StatusInternalServerError, "Error while searching")
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"total":   result.Total,
		"matched": result.Matched,
		"count":   result.Count,
		"gigs":    gigs,
		"query":   fmt.Sprintf("query=%v&maxPrice=%v&deliveryTime=%v", strings.ReplaceAll(q.Query, " ", "-"), q.Max, q.DeliveryTime),
	})
}

func (gh *GigHandler) FindSellerGigs(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 1*time.Second)
	defer cancel()

	userInfo, ok := c.UserContext().Value("current_user").(*types.JWTClaims)
	if !ok {
		log.Println(userInfo)
		return fiber.NewError(http.StatusUnauthorized, "Sign-in first")
	}

	var p types.GigSearchParams
	err := c.ParamsParser(&p)
	if err != nil {
		log.Println("find seller active gigs", err)
		return fiber.NewError(http.StatusBadRequest, "searching error")
	}

	cc, err := gh.grpcClient.GetClient("USER_SERVICE")
	if err != nil {
		log.Printf("FindSellerGigs Error:\n+%v", err)
		return fiber.NewError(http.StatusInternalServerError, "Error while validating seller")
	}
	userGrpcClient := user.NewUserServiceClient(cc)
	s, err := userGrpcClient.FindSeller(ctx, &user.FindSellerRequest{
		BuyerId:  userInfo.UserID,
		SellerId: "",
	})
	if err != nil {
		log.Printf("FindSellerGigs Error:\n+%v", err)
		return fiber.NewError(http.StatusInternalServerError, "Error while validating seller")
	}

	gigs, err := gh.gigSvc.FindSellerGigs(ctx, true, s.Id, &p)
	if err != nil {
		log.Println("find seller active gigs", err)
		return fiber.NewError(http.StatusInternalServerError, "Error while searching")
	}

	result, err := gh.gigSvc.FindAndMapSellerInGigs(ctx, userGrpcClient, gigs)
	if err != nil {
		log.Println("find seller active gigs", err)
		return fiber.NewError(http.StatusInternalServerError, "Error while searching")
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"gigs":  result,
		"total": len(result),
	})
}

func (gh *GigHandler) FindSellerInactiveGigs(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 1*time.Second)
	defer cancel()

	var p types.GigSearchParams
	err := c.ParamsParser(&p)
	if err != nil {
		log.Println("find seller inactive gigs", err)
		return fiber.NewError(http.StatusBadRequest, "searching error")
	}

	gigs, err := gh.gigSvc.FindSellerGigs(ctx, false, c.Params("sellerId"), &p)
	if err != nil {
		log.Println("find seller inactive gigs", err)
		return fiber.NewError(http.StatusInternalServerError, "Error while searching")
	}

	cc, err := gh.grpcClient.GetClient("USER_SERVICE")
	if err != nil {
		log.Println("find seller inactive gigs", err)
		return fiber.NewError(http.StatusInternalServerError, "Error while searching gig")
	}

	userGrpcClient := user.NewUserServiceClient(cc)
	result, err := gh.gigSvc.FindAndMapSellerInGigs(ctx, userGrpcClient, gigs)
	if err != nil {
		log.Println("find seller inactive gigs", err)
		return fiber.NewError(http.StatusInternalServerError, "Error while searching")
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"gigs":  result,
		"total": len(result),
	})
}

func (gh *GigHandler) FindGigByCategory(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 1*time.Second)
	defer cancel()

	var p types.GigSearchParams
	err := c.ParamsParser(&p)
	if err != nil {
		log.Println("find gig by category", err)
		return fiber.NewError(http.StatusBadRequest, "searching error")
	}

	gigs, err := gh.gigSvc.FindGigByCategory(ctx, c.Params("category"), &p)
	if err != nil {
		log.Println("find gig by category", err)
		return fiber.NewError(http.StatusInternalServerError, "Error while searching")
	}

	cc, err := gh.grpcClient.GetClient("USER_SERVICE")
	if err != nil {
		log.Println("find gig by category", err)
		return fiber.NewError(http.StatusInternalServerError, "Error while searching gig")
	}

	userGrpcClient := user.NewUserServiceClient(cc)
	result, err := gh.gigSvc.FindAndMapSellerInGigs(ctx, userGrpcClient, gigs)
	if err != nil {
		log.Println("find gig by category", err)
		return fiber.NewError(http.StatusInternalServerError, "Error while searching")
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"gigs":  result,
		"total": len(result),
	})
}

func (gh *GigHandler) FindSimilarGigs(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 2*time.Second)
	defer cancel()

	var p types.GigSearchParams
	err := c.ParamsParser(&p)
	if err != nil {
		log.Println("find similar gigs", err)
		return fiber.NewError(http.StatusBadRequest, "searching error")
	}

	gig, err := gh.gigSvc.FindGigByID(ctx, c.Params("gigId"))
	if err != nil {
		log.Println("find similar gigs", err)
		return fiber.NewError(http.StatusInternalServerError, "Error while searching")
	}

	gigs, err := gh.gigSvc.FindSimilarGigs(ctx, &p, gig)
	if err != nil {
		log.Println("find similar gigs", err)
		return fiber.NewError(http.StatusInternalServerError, "Error while searching")
	}

	cc, err := gh.grpcClient.GetClient("USER_SERVICE")
	if err != nil {
		log.Println("find similar gigs", err)
		return fiber.NewError(http.StatusInternalServerError, "Error while searching gig")
	}

	userGrpcClient := user.NewUserServiceClient(cc)
	result, err := gh.gigSvc.FindAndMapSellerInGigs(ctx, userGrpcClient, gigs)
	if err != nil {
		log.Println("find similar gigs", err)
		return fiber.NewError(http.StatusInternalServerError, "Error while searching")
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"gigs":  result,
		"total": len(result),
	})
}

func (gh *GigHandler) GetPopularGigs(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 1*time.Second)
	defer cancel()

	gigs, err := gh.gigSvc.GetPopularGigs(ctx)
	if err != nil {
		log.Println("get popular gigs", err)
		return fiber.NewError(http.StatusInternalServerError, "Error while searching")
	}

	cc, err := gh.grpcClient.GetClient("USER_SERVICE")
	if err != nil {
		log.Println("get popular gigs", err)
		return fiber.NewError(http.StatusInternalServerError, "Error while searching gig")
	}

	userGrpcClient := user.NewUserServiceClient(cc)
	result, err := gh.gigSvc.FindAndMapSellerInGigs(ctx, userGrpcClient, gigs)
	if err != nil {
		log.Println("get popular gigs", err)
		return fiber.NewError(http.StatusInternalServerError, "Error while searching")
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"gigs":  result,
		"total": len(result),
	})
}

func (gh *GigHandler) Create(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 10*time.Second)
	defer cancel()

	userInfo, ok := c.UserContext().Value("current_user").(*types.JWTClaims)
	if !ok {
		log.Println(userInfo)
		return fiber.NewError(http.StatusUnauthorized, "Sign-in first")
	}

	data := new(types.CreateGigDTO)
	err := c.BodyParser(data)
	if err != nil {
		log.Println("create gig", err)
		return fiber.NewError(http.StatusInternalServerError, "Error while parsing body")
	}

	cc, err := gh.grpcClient.GetClient("USER_SERVICE")
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "Error while searching gig")
	}

	userGrpcClient := user.NewUserServiceClient(cc)
	s, err := userGrpcClient.FindSeller(ctx, &user.FindSellerRequest{
		SellerId: "",
		BuyerId:  userInfo.UserID,
	})
	if err != nil {
		log.Println("create gig", err)
		return fiber.NewError(http.StatusInternalServerError, "Invalid seller data")
	}

	data.SellerID = s.Id
	err = gh.validate.Struct(data)
	if err != nil {
		log.Printf("create gig:\n%+v", err)
		return fiber.NewError(http.StatusBadRequest, "invalid data")
	}

	formHeader, err := c.FormFile("imageFile")
	if err != nil {
		log.Printf("create gig:\n%+v", err)
		return fiber.NewError(http.StatusBadRequest, "failed reading image file")
	}

	if formHeader.Size > 2*1024*1024 {
		log.Printf("create gig. File is too large")
		return fiber.NewError(http.StatusBadRequest, "file is larger than 2MB")
	}

	if !util.ValidateImgExtension(formHeader) {
		log.Println("create gig file type is unsupported")
		return fiber.NewError(http.StatusBadRequest, "file type is unsupported")
	}

	data.ImageFile, err = formHeader.Open()
	if err != nil {
		log.Printf("create gig error opening image file:\n%+v", err)
		return fiber.NewError(http.StatusBadRequest, "failed reading image file")
	}

	str := util.RandomStr(32)
	uploadResult, err := gh.cld.UploadImg(ctx, data.ImageFile, str)
	if err != nil {
		log.Printf("create gig error upload image:\n%+v", err)
		return fiber.NewError(http.StatusBadRequest, "failed upload file")
	}

	data.CoverImage = uploadResult.SecureURL
	gig, err := gh.gigSvc.Create(ctx, data)
	if err != nil {
		log.Println("create gig", err)
		return fiber.NewError(http.StatusInternalServerError, "Error while creating")
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"gig": gig,
	})
}

func (gh *GigHandler) Update(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 10*time.Second)
	defer cancel()

	gig, err := gh.gigSvc.FindGigBySellerIDAndGigID(ctx, c.Params("sellerId"), c.Params("gigId"))
	if err != nil {
		log.Println("update gig", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(http.StatusNotFound, "Gig did not found")
		}
		return fiber.NewError(http.StatusInternalServerError, "Error finding gig")
	}

	data := new(types.UpdateGigDTO)
	err = c.BodyParser(data)
	if err != nil {
		log.Println("update gig", err)
		return fiber.NewError(http.StatusInternalServerError, "Error while parsing body")
	}

	err = gh.validate.Struct(data)
	if err != nil {
		log.Printf("update gig:\n%+v", err)
		return fiber.NewError(http.StatusBadRequest, "invalid data")
	}

	if data.CoverImage == "" {
		formHeader, err := c.FormFile("imageFile")
		if err != nil {
			log.Printf("update gig:\n%+v", err)
			return fiber.NewError(http.StatusBadRequest, "failed reading image file")
		}

		if formHeader.Size > 2*1024*1024 {
			log.Printf("update gig. File is too large")
			return fiber.NewError(http.StatusBadRequest, "file is larger than 2MB")
		}

		if !util.ValidateImgExtension(formHeader) {
			log.Println("update gig file type is unsupported")
			return fiber.NewError(http.StatusBadRequest, "file type is unsupported")
		}

		data.ImageFile, err = formHeader.Open()
		if err != nil {
			log.Printf("update gig error opening image file:\n%+v", err)
			return fiber.NewError(http.StatusBadRequest, "failed reading image file")
		}

		str := util.RandomStr(32)
		uploadResult, err := gh.cld.UploadImg(ctx, data.ImageFile, str)
		if err != nil {
			log.Printf("update gig error upload image:\n%+v", err)
			return fiber.NewError(http.StatusBadRequest, "failed upload file")
		}

		// Find the index where "jobber" starts
		startIdxPublicID := strings.Index(gig.CoverImage, "jobber/gig/")
		if startIdxPublicID == -1 {
			fmt.Println("Substring publicId directory did not found")
			return fiber.NewError(http.StatusInternalServerError, "Error processing gig cover image")
		}

		// Find the index where ".webp, .png, .jpeg, .jpg" ends
		endIdxPublicID := strings.LastIndex(gig.CoverImage, ".")
		if endIdxPublicID == -1 {
			fmt.Println("Substring image extension did not found")
			return fiber.NewError(http.StatusInternalServerError, "Error processing gig cover image")
		}

		// Extract the gigPublicID from "jobber" to ".webp"
		gigPublicID := gig.CoverImage[startIdxPublicID:endIdxPublicID]
		_, err = gh.cld.Destroy(ctx, gigPublicID)
		if err != nil {
			fmt.Println("Error updating gig cover image")
			return fiber.NewError(http.StatusInternalServerError, "Error updating gig cover image")
		}

		data.CoverImage = uploadResult.SecureURL
	}

	gig, err = gh.gigSvc.Update(ctx, data)
	if err != nil {
		log.Println("update gig", err)
		return fiber.NewError(http.StatusInternalServerError, "Error while updating")
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"gig": gig,
	})
}

func (gh *GigHandler) ActivateGigStatus(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 1*time.Second)
	defer cancel()

	gig, err := gh.gigSvc.FindGigBySellerIDAndGigID(ctx, c.Params("sellerId"), c.Params("gigId"))
	if err != nil {
		log.Println("activate gig status", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(http.StatusNotFound, "Gig did not found")
		}
		return fiber.NewError(http.StatusInternalServerError, "Error finding gig")
	}

	err = gh.gigSvc.ChangeGigStatus(ctx, gig.ID.String(), true)
	if err != nil {
		log.Println("activate gig status", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(http.StatusNotFound, "gig is not found")
		}
		return fiber.NewError(http.StatusInternalServerError, "Error while changing gig status to active")
	}

	return c.SendStatus(http.StatusOK)
}

func (gh *GigHandler) DeactivateGigStatus(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.UserContext(), 1*time.Second)
	defer cancel()

	gig, err := gh.gigSvc.FindGigBySellerIDAndGigID(ctx, c.Params("sellerId"), c.Params("gigId"))
	if err != nil {
		log.Println("deactivate gig status", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(http.StatusNotFound, "Gig did not found")
		}
		return fiber.NewError(http.StatusInternalServerError, "Error validating gig")
	}

	err = gh.gigSvc.ChangeGigStatus(ctx, gig.ID.String(), false)
	if err != nil {
		log.Println("deactivate gig status", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fiber.NewError(http.StatusNotFound, "gig is not found")
		}
		return fiber.NewError(http.StatusInternalServerError, "Error while changing gig status to deactive")
	}

	return c.SendStatus(http.StatusOK)
}
