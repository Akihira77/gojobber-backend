package service

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/Akihira77/gojobber/services/5-gig/types"
	"github.com/Akihira77/gojobber/services/common/genproto/user"
	"gorm.io/gorm"
)

type GigServiceImpl interface {
	FindGigByID(ctx context.Context, id string) (*types.GigDTO, error)
	FindGigBySellerIDAndGigID(ctx context.Context, sellerId, id string) (*types.GigDTO, error)
	GigQuerySearch(ctx context.Context, p *types.GigSearchParams, q *types.GigSearchQuery) (types.GigSearchQueryResult, error)
	FindSellerGigs(ctx context.Context, gigStatus bool, sellerId string, p *types.GigSearchParams) ([]types.GigDTO, error)
	FindGigByCategory(ctx context.Context, category string, p *types.GigSearchParams) ([]types.GigDTO, error)
	GetPopularGigs(ctx context.Context, p *types.GigSearchParams) ([]types.GigDTO, error)
	FindSimilarGigs(ctx context.Context, p *types.GigSearchParams, data *types.GigDTO) ([]types.GigDTO, error)
	Create(ctx context.Context, data *types.CreateGigDTO) (*types.GigDTO, error)
	Update(ctx context.Context, data *types.UpdateGigDTO) (*types.GigDTO, error)
	ChangeGigStatus(ctx context.Context, gigId string, s bool) error
	DeleteGigByID(ctx context.Context, gigId string) error
	FindAndMapSellerInGigs(ctx context.Context, userGrpcClient user.UserServiceClient, gigs []types.GigDTO) ([]types.GigSellerDTO, error)
}

type GigService struct {
	db *gorm.DB
}

func NewGigService(db *gorm.DB) GigServiceImpl {
	return &GigService{
		db: db,
	}
}

func (gs *GigService) FindAndMapSellerInGigs(ctx context.Context, userGrpcClient user.UserServiceClient, gigs []types.GigDTO) ([]types.GigSellerDTO, error) {
	var result []types.GigSellerDTO
	var wg sync.WaitGroup
	errCh := make(chan error, len(gigs))

	for _, gig := range gigs {
		wg.Add(1)
		go func() {
			defer wg.Done()

			newCtx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
			defer cancel()

			s, err := userGrpcClient.FindSeller(newCtx, &user.FindSellerRequest{
				SellerId: gig.SellerID,
			})

			if err != nil {
				errCh <- fmt.Errorf("Invalid seller data")
				return
			}

			result = append(result, types.GigSellerDTO{
				Seller: types.SellerOverview{
					SellerID:  gig.SellerID,
					FullName:  s.FullName,
					RatingSum: uint64(s.RatingSum),
					RatingCategories: types.RatingCategory{
						One:   uint(s.RatingCategories.One),
						Two:   uint(s.RatingCategories.Two),
						Three: uint(s.RatingCategories.Three),
						Four:  uint(s.RatingCategories.Four),
						Five:  uint(s.RatingCategories.Five),
					},
					RatingsCount: uint64(s.RatingsCount),
				},
				Gig: gig,
			})

			errCh <- nil
		}()

	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		if err != nil {
			fmt.Println(err)
			return result, err
		}
	}

	return result, nil
}

func (gs *GigService) FindGigByID(ctx context.Context, id string) (*types.GigDTO, error) {
	var gig types.GigDTO
	result := gs.db.WithContext(ctx).
		Model(&types.Gig{}).
		First(&gig, "id = ?", id)
	return &gig, result.Error
}

func (gs *GigService) FindGigBySellerIDAndGigID(ctx context.Context, sellerId, id string) (*types.GigDTO, error) {
	var gig types.GigDTO
	result := gs.db.WithContext(ctx).
		Model(&types.Gig{}).
		First(&gig, "id = ? AND seller_id = ?", id, sellerId)
	return &gig, result.Error
}

func (gs *GigService) GigQuerySearch(ctx context.Context, p *types.GigSearchParams, q *types.GigSearchQuery) (types.GigSearchQueryResult, error) {
	var gigs []types.GigDTO

	sanitizeNumber := func(num int) int {
		if num <= 0 || num > math.MaxInt32 {
			num = math.MaxInt32
		}
		return num
	}
	//HACK: if the number is too high or too low
	q.Max = sanitizeNumber(q.Max)
	q.DeliveryTime = sanitizeNumber(q.DeliveryTime)

	var total int64
	query := gs.db.
		Debug().
		WithContext(ctx).
		Model(&types.Gig{}).
		Count(&total).
		Where("active = true AND price <= ? AND expected_delivery_days <= ?", q.Max, q.DeliveryTime)

	if q.Query != "" {
		q.Query = strings.ReplaceAll(strings.ToLower(q.Query), "-", " ")
		query = query.
			Where(`
				title_tokens @@ websearch_to_tsquery('english', ?) OR 
				category_tokens @@ websearch_to_tsquery('english', ?) OR 
				sub_categories_tokens @@ websearch_to_tsquery('english', ?) OR 
				tags_tokens @@ websearch_to_tsquery('english', ?) OR
				description_tokens @@ websearch_to_tsquery('english', ?)`,
				q.Query, q.Query, q.Query, q.Query, q.Query)

		orderClause := fmt.Sprintf(`
			ts_rank(title_tokens, websearch_to_tsquery('english', '%s')) +
			ts_rank(category_tokens, websearch_to_tsquery('english', '%s')) +
			ts_rank(sub_categories_tokens, websearch_to_tsquery('english', '%s')) +
			ts_rank(tags_tokens, websearch_to_tsquery('english', '%s')) + 
			ts_rank(description_tokens, websearch_to_tsquery('english', '%s')) DESC`,
			q.Query, q.Query, q.Query, q.Query, q.Query)

		query = query.Order(orderClause)
	}

	var matchedQuery int64
	result := query.
		Order(`rating_sum DESC`).
		Count(&matchedQuery).
		Offset((p.Page - 1) * p.Size).
		Limit(p.Size).
		Find(&gigs)

	return types.GigSearchQueryResult{
		Total:   total,
		Matched: matchedQuery,
		Count:   len(gigs),
		Gigs:    gigs,
	}, result.Error
}

func (gs *GigService) Create(ctx context.Context, data *types.CreateGigDTO) (*types.GigDTO, error) {
	tx := gs.db.
		Debug().
		WithContext(ctx).
		Begin(&sql.TxOptions{
			Isolation: sql.LevelSerializable,
		})

	var gig types.Gig
	result := tx.
		Raw(
			`INSERT INTO gigs (
                id, 
                seller_id, 
                title, 
                description, 
                category, 
                sub_categories, 
                tags, 
                active,
			    expected_delivery_days, 
                ratings_count, 
                rating_sum, 
                rating_categories, 
                price, 
                cover_image, 
                created_at,
			    title_tokens, 
                description_tokens, 
                category_tokens, 
                sub_categories_tokens, 
                tags_tokens
            )
			VALUES (
                uuid_generate_v4(),
                ?,
                ?,
                ?,
                ?,
                ?,
                ?,
                ?,
                ?,
                ?,
                ?,
                ?::jsonb,
                ?,
                ?,
                ?,
                strip(to_tsvector('english', ?)),
                strip(to_tsvector('english', ?)),
                strip(to_tsvector('english', ?)),
                strip(to_tsvector('english', ?)),
			    strip(to_tsvector('english', ?))) 
            RETURNING *`,
			data.SellerID,
			data.Title,
			data.Description,
			data.Category,
			data.SubCategories,
			data.Tags,
			data.Active,
			data.ExpectedDeliveryDays,
			0, // rating count
			0, // rating sum
			// rating categories
			fmt.Sprintf(`{"five": %d,"four": %d,"three": %d,"two": %d,"one": %d}`,
				0, 0, 0, 0, 0),
			data.Price,
			data.CoverImage,
			time.Now(), // created at
			data.Title,
			data.Description,
			data.Category,
			strings.Join(data.SubCategories, ","),
			strings.Join(data.Tags, ","),
		).
		Scan(&gig)

	if result.Error != nil {
		tx.Rollback()
		return &types.GigDTO{}, result.Error
	}

	result = tx.Commit()
	if result.Error != nil {
		tx.Rollback()
	}

	return &types.GigDTO{
		ID:                   gig.ID,
		SellerID:             gig.SellerID,
		Title:                gig.Title,
		Description:          gig.Description,
		Category:             gig.Category,
		SubCategories:        gig.SubCategories,
		Tags:                 gig.Tags,
		Active:               gig.Active,
		ExpectedDeliveryDays: int(gig.ExpectedDeliveryDays),
		RatingsCount:         gig.RatingsCount,
		RatingSum:            gig.RatingSum,
		RatingCategories:     gig.RatingCategories,
		Price:                gig.Price,
		CoverImage:           gig.CoverImage,
		SortID:               gig.SortID,
		CreatedAt:            gig.CreatedAt,
	}, result.Error
}

func (gs *GigService) DeleteGigByID(ctx context.Context, gigID string) error {
	result := gs.db.
		WithContext(ctx).
		Model(&types.Gig{}).
		Where("id = ?", gigID).
		Delete(&types.Gig{})

	return result.Error
}

func (gs *GigService) FindGigByCategory(ctx context.Context, c string, p *types.GigSearchParams) ([]types.GigDTO, error) {
	var gigs []types.GigDTO
	c = strings.ReplaceAll(strings.ToLower(c), "-", " ")
	result := gs.db.
		Debug().
		WithContext(ctx).
		Model(&types.Gig{}).
		Offset((p.Page-1)*p.Size).
		Limit(p.Size).
		Find(&gigs, "active = true AND category_tokens @@ websearch_to_tsquery('english', ?)", strings.ToLower(c))

	return gigs, result.Error
}

func (gs *GigService) FindSellerGigs(ctx context.Context, active bool, sellerID string, p *types.GigSearchParams) ([]types.GigDTO, error) {
	var gigs []types.GigDTO
	result := gs.db.
		Debug().
		WithContext(ctx).
		Model(&types.Gig{}).
		Offset((p.Page-1)*p.Size).
		Limit(p.Size).
		Find(&gigs, "seller_id = ? AND active = ?", sellerID, active)

	return gigs, result.Error
}

func (gs *GigService) FindSimilarGigs(ctx context.Context, p *types.GigSearchParams, gig *types.GigDTO) ([]types.GigDTO, error) {
	query := gs.db.
		Debug().
		WithContext(ctx).
		Model(&types.Gig{}).
		Where(`
			id <> ? AND
			active = true AND
			(title_tokens @@ websearch_to_tsquery('english', ?) OR 
			category_tokens @@ websearch_to_tsquery('english', ?) OR 
			sub_categories_tokens @@ websearch_to_tsquery('english', ?) OR 
			tags_tokens @@ websearch_to_tsquery('english', ?))`,
			gig.ID,
			strings.ToLower(gig.Title),
			strings.ToLower(gig.Category),
			strings.Join(gig.SubCategories, ","),
			strings.Join(gig.Tags, ","))

	orderClause := fmt.Sprintf(`
		ts_rank(title_tokens, websearch_to_tsquery('english', '%s')) +
		ts_rank(category_tokens, websearch_to_tsquery('english', '%s')) +
		ts_rank(sub_categories_tokens, websearch_to_tsquery('english', '%s')) +
		ts_rank(tags_tokens, websearch_to_tsquery('english', '%s')) DESC`,
		strings.ToLower(gig.Title),
		strings.ToLower(gig.Category),
		strings.Join(gig.SubCategories, ","),
		strings.Join(gig.Tags, ","))

	var gigs []types.GigDTO
	result := query.
		Order(orderClause).
		Order(`rating_sum DESC`).
		Offset((p.Page - 1) * p.Size).
		Limit(p.Size).
		Find(&gigs)

	return gigs, result.Error
}

func (gs *GigService) GetPopularGigs(ctx context.Context, p *types.GigSearchParams) ([]types.GigDTO, error) {
	var gigs []types.GigDTO
	result := gs.db.
		WithContext(ctx).
		Model(&types.Gig{}).
		Order("rating_sum DESC").
		Offset((p.Page - 1) * p.Size).
		Limit(p.Size).
		Find(&gigs)

	return gigs, result.Error
}

func (gs *GigService) Update(ctx context.Context, data *types.UpdateGigDTO) (*types.GigDTO, error) {
	var gig types.Gig
	result := gs.db.
		Debug().
		Raw(
			`UPDATE gigs 
			SET title = ?, description = ?, category = ?, sub_categories = ?, tags = ?,
			expected_delivery_days = ?, price = ?, cover_image = ?, title_tokens = strip(to_tsvector('english', ?)), 
			description_tokens = strip(to_tsvector('english', ?)), category_tokens = strip(to_tsvector('english', ?)), 
			sub_categories_tokens = strip(to_tsvector('english', ?)), tags_tokens = strip(to_tsvector('english', ?)) RETURNING *`,
			data.Title,
			data.Description,
			data.Category,
			data.SubCategories,
			data.Tags,
			data.ExpectedDeliveryDays,
			data.Price,
			data.CoverImage,
			data.Title,
			data.Description,
			data.Category,
			strings.Join(data.SubCategories, ","),
			strings.Join(data.Tags, ","),
		).Scan(&gig)

	return &types.GigDTO{
		ID:                   gig.ID,
		SellerID:             gig.SellerID,
		Title:                gig.Title,
		Description:          gig.Description,
		Category:             gig.Category,
		SubCategories:        gig.SubCategories,
		Tags:                 gig.Tags,
		Active:               gig.Active,
		ExpectedDeliveryDays: int(gig.ExpectedDeliveryDays),
		RatingsCount:         gig.RatingsCount,
		RatingSum:            gig.RatingSum,
		RatingCategories:     gig.RatingCategories,
		Price:                gig.Price,
		CoverImage:           gig.CoverImage,
		SortID:               gig.SortID,
		CreatedAt:            gig.CreatedAt,
	}, result.Error
}

func (gs *GigService) ChangeGigStatus(ctx context.Context, gigID string, s bool) error {
	result := gs.db.
		WithContext(ctx).
		Model(&types.Gig{}).
		Where("id = ?", gigID).
		Update("active", s)

	return result.Error
}
