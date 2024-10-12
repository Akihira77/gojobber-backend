package handler

import (
	"context"
	"log"

	"github.com/Akihira77/gojobber/services/4-user/service"
	"github.com/Akihira77/gojobber/services/4-user/types"
	"github.com/Akihira77/gojobber/services/common/genproto/user"
	"google.golang.org/grpc"
)

type UserGRPCHandler struct {
	buyerSvc  service.BuyerServiceImpl
	sellerSvc service.SellerServiceImpl
	user.UnimplementedUserServiceServer
}

func NewUserGRPCHandler(grpc *grpc.Server, buyerSvc service.BuyerServiceImpl, sellerSvc service.SellerServiceImpl) {
	gRPCHandler := &UserGRPCHandler{
		buyerSvc:  buyerSvc,
		sellerSvc: sellerSvc,
	}

	// register the BuyerServiceServer
	user.RegisterUserServiceServer(grpc, gRPCHandler)
}

func (h *UserGRPCHandler) SaveBuyerData(ctx context.Context, req *user.SaveBuyerRequest) (*user.SaveBuyerResponse, error) {
	log.Println("SaveBuyerData receive data", req)
	res := &user.SaveBuyerResponse{
		Success: false,
		Message: "",
	}

	b := &types.Buyer{
		ID:             req.Id,
		Username:       req.Username,
		Email:          req.Email,
		Country:        req.Country,
		IsSeller:       req.IsSeller,
		ProfilePicture: req.ProfilePicture,
		CreatedAt:      req.CreatedAt.AsTime(),
	}

	err := h.buyerSvc.Create(ctx, *b)
	if err != nil {
		res.Message = err.Error()
		log.Println("Error", err)
		return res, err
	}

	res.Success = true
	return res, nil
}

func (h *UserGRPCHandler) FindSeller(ctx context.Context, req *user.FindSellerRequest) (*user.FindSellerResponse, error) {
	log.Println("FindSeller receive data", req)
	seller, err := h.sellerSvc.FindSellerOverviewByID(ctx, req.BuyerId, req.SellerId)
	if err != nil {
		return nil, err
	}

	return &user.FindSellerResponse{
		FullName:     seller.FullName,
		Email:        seller.Email,
		RatingsCount: int64(seller.RatingsCount),
		RatingSum:    int64(seller.RatingSum),
		RatingCategories: &user.RatingCategory{
			One:   int32(seller.RatingCategories.One),
			Two:   int32(seller.RatingCategories.Two),
			Three: int32(seller.RatingCategories.Three),
			Four:  int32(seller.RatingCategories.Four),
			Five:  int32(seller.RatingCategories.Five),
		},
	}, nil
}

func (h *UserGRPCHandler) UpdateSellerBalance(ctx context.Context, req *user.UpdateSellerBalanceRequest) (*user.UpdateSellerBalanceResponse, error) {
	log.Println("UpdateSellerBalance receive data", req)
	seller, err := h.sellerSvc.UpdateBalance(ctx, req.SellerId, req.Amount)
	if err != nil {
		return nil, err
	}

	return &user.UpdateSellerBalanceResponse{
		Id:           seller.ID,
		Bio:          seller.Bio,
		FullName:     seller.FullName,
		Email:        seller.Email,
		RatingsCount: int64(seller.RatingsCount),
		RatingSum:    int64(seller.RatingSum),
		RatingCategories: &user.RatingCategory{
			One:   int32(seller.RatingCategories.One),
			Two:   int32(seller.RatingCategories.Two),
			Three: int32(seller.RatingCategories.Three),
			Four:  int32(seller.RatingCategories.Four),
			Five:  int32(seller.RatingCategories.Five),
		},
		StripeAccountID: seller.StripeAccountID,
		AccountBalance:  seller.AccountBalance,
	}, nil

}
