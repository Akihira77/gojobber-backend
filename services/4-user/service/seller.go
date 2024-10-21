package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Akihira77/gojobber/services/4-user/types"
	"github.com/Akihira77/gojobber/services/4-user/util"
	"gorm.io/gorm"
)

type SellerService struct {
	db *gorm.DB
}

type SellerServiceImpl interface {
	FindSellerByID(ctx context.Context, id string) (*types.SellerDTO, error)
	FindSellerOverviewByID(ctx context.Context, buyerId, sellerId string) (*types.SellerOverview, error)
	FindSellerByBuyerID(ctx context.Context, buyerId string) (*types.Seller, error)
	FindSellerByUsername(ctx context.Context, username string) (*types.SellerDTO, error)
	GetRandomSellers(ctx context.Context, count int) ([]types.SellerDTO, error)
	Create(ctx context.Context, sellerDataInBuyerDB *types.Buyer, data *types.CreateSellerDTO) (*types.SellerDTO, error)
	Update(ctx context.Context, updatedSellerData *types.Seller, data *types.UpdateSellerDTO) error
	UpdateBalance(ctx context.Context, sellerID string, addedBalance uint64) (*types.SellerIncBalanceDTO, error)
}

func NewSellerService(db *gorm.DB) SellerServiceImpl {
	return &SellerService{
		db: db,
	}
}

func (ss *SellerService) FindSellerByID(ctx context.Context, id string) (*types.SellerDTO, error) {
	dbExec := ss.db.
		Debug().
		WithContext(ctx)

	var seller types.SellerDTO
	result := dbExec.
		Model(&types.Seller{}).
		Select(`
			sellers.id, 
			sellers.full_name, 
			buyers.email, 
			buyers.country, 
			buyers.profile_picture, 
			sellers.bio, 
			sellers.ratings_count, 
			sellers.rating_sum, 
			sellers.rating_categories
		`).
		Joins("INNER JOIN buyers ON buyers.id = sellers.buyer_id").
		First(&seller, "sellers.id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}

	var wg sync.WaitGroup
	goRoutineSize := 5
	errCh := make(chan error, goRoutineSize)

	// INFO: 1. GRAB SELLER LANGUAGES
	var sellerLanguages []types.Language
	wg.Add(1)
	go func() {
		defer wg.Done()
		result = dbExec.
			Model(&types.Language{}).
			Select("languages.id, languages.language").
			Joins("INNER JOIN seller_languages ON seller_languages.language_id = languages.id").
			Where("seller_languages.seller_id = ?", seller.ID).
			Find(&sellerLanguages)
		if result.Error != nil {
			errCh <- result.Error
			return
		}
		errCh <- nil
	}()

	// INFO: 2. GRAB SELLER SKILLS
	var sellerSkills []types.Skill
	wg.Add(1)
	go func() {
		defer wg.Done()
		result = dbExec.
			Model(&types.Skill{}).
			Select("skills.id, skills.name").
			Joins("INNER JOIN seller_skills ON seller_skills.skill_id = skills.id").
			Where("seller_skills.seller_id = ?", seller.ID).
			Find(&sellerSkills)
		if result.Error != nil {
			errCh <- result.Error
			return
		}
		errCh <- nil
	}()

	// INFO: 3. GRAB SELLER EDUCATIONS
	var sellerEducations []types.EducationDTO
	wg.Add(1)
	go func() {
		defer wg.Done()
		result = dbExec.
			Model(&types.Education{}).
			Where("seller_id = ?", seller.ID).
			Find(&sellerEducations)
		if result.Error != nil {
			errCh <- result.Error
			return
		}
		errCh <- nil
	}()

	// INFO: 4. GRAB SELLER CERTIFICATES
	var sellerCertificates []types.CertificateDTO
	wg.Add(1)
	go func() {
		defer wg.Done()
		result = dbExec.
			Model(&types.Certificate{}).
			Where("seller_id = ?", seller.ID).
			Find(&sellerCertificates)
		if result.Error != nil {
			errCh <- result.Error
			return
		}
		errCh <- nil
	}()

	// INFO: 5. GRAB SELLER EXPERIENCES
	var sellerExperiences []types.ExperienceDTO
	wg.Add(1)
	go func() {
		defer wg.Done()
		result = dbExec.
			Model(&types.Experience{}).
			Where("seller_id = ?", seller.ID).
			Find(&sellerExperiences)
		if result.Error != nil {
			errCh <- result.Error
			return
		}
		errCh <- nil
	}()

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		if err != nil {
			return nil, fmt.Errorf("unexpected error: %w", err)
		}
	}

	return &types.SellerDTO{
		ID:               seller.ID,
		FullName:         seller.FullName,
		Email:            seller.Email,
		Bio:              seller.Bio,
		RatingsCount:     seller.RatingsCount,
		RatingSum:        seller.RatingSum,
		RatingCategories: seller.RatingCategories,
		Country:          seller.Country,
		ProfilePicture:   seller.ProfilePicture,
		Languages:        sellerLanguages,
		Skills:           sellerSkills,
		Educations:       sellerEducations,
		Certificates:     sellerCertificates,
		Experiences:      sellerExperiences,
	}, nil
}

func (ss *SellerService) FindSellerByBuyerID(ctx context.Context, buyerId string) (*types.Seller, error) {
	var seller types.Seller
	result := ss.db.
		WithContext(ctx).
		Model(&types.Seller{}).
		First(&seller, "buyer_id = ?", buyerId)

	return &seller, result.Error
}

func (ss *SellerService) FindSellerOverviewByID(ctx context.Context, buyerId, sellerId string) (*types.SellerOverview, error) {
	var seller types.SellerOverview
	result := ss.db.
		Debug().
		WithContext(ctx).
		Model(&types.Seller{}).
		Select(`
            sellers.id,
			sellers.full_name, 
			buyers.email, 
			sellers.ratings_count, 
			sellers.rating_sum, 
			sellers.rating_categories,
            sellers.stripe_account_id
		`).
		Joins("INNER JOIN buyers ON buyers.id = sellers.buyer_id").
		First(&seller, "sellers.id = ? OR sellers.buyer_id = ?", sellerId, buyerId)

	return &seller, result.Error
}

func (ss *SellerService) FindSellerByUsername(ctx context.Context, username string) (*types.SellerDTO, error) {
	dbExec := ss.db.
		Debug().
		WithContext(ctx)

	var seller types.SellerDTO
	result := dbExec.
		Model(&types.Seller{}).
		Select(`
			sellers.id, 
			sellers.full_name, 
			buyers.email, 
			buyers.country, 
			buyers.profile_picture, 
			sellers.bio, 
			sellers.ratings_count, 
			sellers.rating_sum, 
			sellers.rating_categories
		`).
		Joins("INNER JOIN buyers ON buyers.id = sellers.buyer_id").
		First(&seller, "buyers.username = ?", username)
	if result.Error != nil {
		return nil, result.Error
	}

	var wg sync.WaitGroup
	goRoutineSize := 5
	errCh := make(chan error, goRoutineSize)

	// INFO: 1. GRAB SELLER LANGUAGES
	var sellerLanguages []types.Language
	wg.Add(1)
	go func() {
		defer wg.Done()
		result = dbExec.
			Model(&types.Language{}).
			Select("languages.id, languages.language").
			Joins("INNER JOIN seller_languages ON seller_languages.language_id = languages.id").
			Where("seller_languages.seller_id = ?", seller.ID).
			Find(&sellerLanguages)
		if result.Error != nil {
			errCh <- result.Error
			return
		}
		errCh <- nil
	}()

	// INFO: 2. GRAB SELLER SKILLS
	var sellerSkills []types.Skill
	wg.Add(1)
	go func() {
		defer wg.Done()
		result = dbExec.
			Model(&types.Skill{}).
			Select("skills.id, skills.name").
			Joins("INNER JOIN seller_skills ON seller_skills.skill_id = skills.id").
			Where("seller_skills.seller_id = ?", seller.ID).
			Find(&sellerSkills)
		if result.Error != nil {
			errCh <- result.Error
			return
		}
		errCh <- nil
	}()

	// INFO: 3. GRAB SELLER EDUCATIONS
	var sellerEducations []types.EducationDTO
	wg.Add(1)
	go func() {
		defer wg.Done()
		result = dbExec.
			Model(&types.Education{}).
			Where("seller_id = ?", seller.ID).
			Find(&sellerEducations)
		if result.Error != nil {
			errCh <- result.Error
			return
		}
		errCh <- nil
	}()

	// INFO: 4. GRAB SELLER CERTIFICATES
	var sellerCertificates []types.CertificateDTO
	wg.Add(1)
	go func() {
		defer wg.Done()
		result = dbExec.
			Model(&types.Certificate{}).
			Where("seller_id = ?", seller.ID).
			Find(&sellerCertificates)
		if result.Error != nil {
			errCh <- result.Error
			return
		}
		errCh <- nil
	}()

	// INFO: 5. GRAB SELLER EXPERIENCES
	var sellerExperiences []types.ExperienceDTO
	wg.Add(1)
	go func() {
		defer wg.Done()
		result = dbExec.
			Model(&types.Experience{}).
			Where("seller_id = ?", seller.ID).
			Find(&sellerExperiences)
		if result.Error != nil {
			errCh <- result.Error
			return
		}
		errCh <- nil
	}()

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		if err != nil {
			return nil, fmt.Errorf("unexpected error: %w", err)
		}
	}

	return &types.SellerDTO{
		ID:               seller.ID,
		FullName:         seller.FullName,
		Email:            seller.Email,
		Bio:              seller.Bio,
		RatingsCount:     seller.RatingsCount,
		RatingSum:        seller.RatingSum,
		RatingCategories: seller.RatingCategories,
		Country:          seller.Country,
		ProfilePicture:   seller.ProfilePicture,
		Languages:        sellerLanguages,
		Skills:           sellerSkills,
		Educations:       sellerEducations,
		Certificates:     sellerCertificates,
		Experiences:      sellerExperiences,
	}, nil
}

func (ss *SellerService) GetRandomSellers(ctx context.Context, total int) ([]types.SellerDTO, error) {
	var sellers []types.SellerDTO

	dbExec := ss.db.
		Debug().
		WithContext(ctx)
	result := dbExec.
		Model(&types.Seller{}).
		Raw(`
			SELECT 
				sellers.id, 
				sellers.full_name, 
				buyers.email, 
				buyers.country, 
				buyers.profile_picture, 
				sellers.bio, 
				sellers.ratings_count, 
				sellers.rating_sum, 
				sellers.rating_categories
			FROM sellers 
			TABLESAMPLE SYSTEM_ROWS(?)
			INNER JOIN buyers ON buyers.id = sellers.buyer_id
		`,
			total).
		Scan(&sellers)

	for i := range sellers {
		var wg sync.WaitGroup
		goRoutineSize := 5
		errCh := make(chan error, goRoutineSize)

		//INFO: 1. GRAB SELLER LANGUAGES
		wg.Add(1)
		var sellerLanguages []types.Language
		go func() {
			defer wg.Done()
			result = dbExec.
				Model(&types.Language{}).
				Select("languages.id, languages.language").
				Joins("INNER JOIN seller_languages ON seller_languages.language_id = languages.id").
				Find(&sellerLanguages)
			if result.Error != nil {
				errCh <- result.Error
				return
			}
			errCh <- nil
		}()

		//INFO: 2. GRAB SELLER SKILLS
		wg.Add(1)
		var sellerSkills []types.Skill
		go func() {
			defer wg.Done()
			result = dbExec.
				Model(&types.Skill{}).
				Select("skills.id, skills.name").
				Joins("INNER JOIN seller_skills ON seller_skills.skill_id = skills.id").
				Find(&sellerSkills)
			if result.Error != nil {
				errCh <- result.Error
				return
			}
			errCh <- nil
		}()

		//INFO: 3. GRAB SELLER EDUCATIONS
		wg.Add(1)
		var sellerEducations []types.EducationDTO
		go func() {
			defer wg.Done()
			result = dbExec.
				Model(&types.Education{}).
				Find(&sellerEducations, "seller_id = ?", sellers[i].ID)
			if result.Error != nil {
				errCh <- result.Error
				return
			}
			errCh <- nil
		}()

		//INFO: 4. GRAB SELLER CERTIFICATES
		wg.Add(1)
		var sellerCertificates []types.CertificateDTO
		go func() {
			defer wg.Done()
			result = dbExec.
				Model(&types.Certificate{}).
				Find(&sellerCertificates, "seller_id = ?", sellers[i].ID)
			if result.Error != nil {
				errCh <- result.Error
				return
			}
			errCh <- nil
		}()

		//INFO: 5. GRAB SELLER EXPERIENCES
		wg.Add(1)
		var sellerExperiences []types.ExperienceDTO
		go func() {
			defer wg.Done()
			result = dbExec.
				Model(&types.Experience{}).
				Find(&sellerExperiences, "seller_id = ?", sellers[i].ID)
			if result.Error != nil {
				errCh <- result.Error
				return
			}
			errCh <- nil
		}()

		// Wait for all goroutines to finish
		go func() {
			wg.Wait()
			close(errCh) // Close the error channel after all goroutines are done
		}()

		for err := range errCh {
			if err != nil {
				log.Println(err)
				return []types.SellerDTO{}, fmt.Errorf("unexpected error")
			}
		}

		sellers[i] = types.SellerDTO{
			ID:               sellers[i].ID,
			FullName:         sellers[i].FullName,
			Email:            sellers[i].Email,
			Bio:              sellers[i].Bio,
			RatingsCount:     sellers[i].RatingsCount,
			RatingSum:        sellers[i].RatingSum,
			RatingCategories: sellers[i].RatingCategories,
			ProfilePicture:   sellers[i].ProfilePicture,
			Country:          sellers[i].Country,
			Languages:        sellerLanguages,
			Skills:           sellerSkills,
			Educations:       sellerEducations,
			Certificates:     sellerCertificates,
			Experiences:      sellerExperiences,
		}

	}
	return sellers, result.Error
}

func (ss *SellerService) Create(ctx context.Context, sellerDataInBuyerDB *types.Buyer, data *types.CreateSellerDTO) (*types.SellerDTO, error) {
	tx := ss.db.
		WithContext(ctx).
		Begin(&sql.TxOptions{})

	seller := types.Seller{
		ID:           util.RandomStr(64),
		BuyerID:      sellerDataInBuyerDB.ID,
		FullName:     data.FullName,
		Bio:          data.Bio,
		RatingsCount: 0,
		RatingSum:    0,
		RatingCategories: types.RatingCategory{
			One:   0,
			Two:   0,
			Three: 0,
			Four:  0,
			Five:  0,
		},
		StripeAccountID: data.StripeAccountID,
		AccountBalance:  0,
	}
	result := tx.
		Model(&types.Seller{}).
		Create(&seller)
	if result.Error != nil {
		tx.Rollback()
		return &types.SellerDTO{}, result.Error
	}

	var experiences []types.ExperienceDTO
	var educations []types.EducationDTO
	var certificates []types.CertificateDTO

	//INFO: 1. SAVE LANGUAGES
	for _, lang := range data.Languages {
		result = tx.
			Model(&types.Language{}).
			First(&lang, "id = ?", lang.ID)
		if result.Error != nil {
			tx.Rollback()
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return &types.SellerDTO{}, errors.New("invalid language")
			}

			return &types.SellerDTO{}, fmt.Errorf("Error finding language. %v", result.Error)
		}
		result = tx.
			Model(&types.SellerLanguage{}).
			Create(&types.SellerLanguage{
				SellerID:   seller.ID,
				LanguageID: lang.ID,
			})
		if result.Error != nil {
			tx.Rollback()
			return &types.SellerDTO{}, fmt.Errorf("Error saving languages %v", result.Error)
		}
	}

	//INFO: 2. SAVE SKILLS
	for _, skill := range data.Skills {
		result = tx.
			Model(&types.Skill{}).
			FirstOrCreate(&skill, "name = ?", skill.Name)
		if result.Error != nil {
			return &types.SellerDTO{}, fmt.Errorf("Error saving skills %v", result.Error)
		}

		result = tx.
			Model(&types.SellerSkill{}).
			Create(&types.SellerSkill{
				SellerID: seller.ID,
				SkillID:  skill.ID,
			})
		if result.Error != nil {
			tx.Rollback()
			return &types.SellerDTO{}, fmt.Errorf("Error saving sellerSkills %v", result.Error)
		}
	}

	//INFO: 3. SAVE EDUCATIONS
	for _, education := range data.Educations {
		edu := types.Education{
			Title:      education.Title,
			Major:      education.Major,
			Country:    education.Country,
			University: education.University,
			SellerID:   seller.ID,
		}
		result = tx.
			Model(&types.Education{}).
			Create(&edu)
		if result.Error != nil {
			tx.Rollback()
			return &types.SellerDTO{}, fmt.Errorf("Error saving educations %v", result.Error)
		}

		educations = append(educations, types.EducationDTO{
			ID:         edu.ID,
			Title:      edu.Title,
			Major:      edu.Major,
			University: edu.University,
			Country:    edu.Country,
		})
	}

	//INFO: 4. SAVE CERTIFICATES
	for _, certif := range data.Certificates {
		cert := types.Certificate{
			Name:     certif.Name,
			From:     certif.From,
			Year:     certif.Year,
			SellerID: seller.ID,
		}
		result = tx.
			Model(&types.Certificate{}).
			Create(&cert)
		if result.Error != nil {
			tx.Rollback()
			return &types.SellerDTO{}, fmt.Errorf("Error saving certificates %v", result.Error)
		}

		certificates = append(certificates, types.CertificateDTO{
			ID:   cert.ID,
			Name: cert.Name,
			From: cert.From,
			Year: cert.Year,
		})
	}

	//INFO: 5. SAVE EXPERIENCES
	for _, experience := range data.Experiences {
		exp := types.Experience{
			Title:                experience.Title,
			Description:          experience.Description,
			Company:              experience.Company,
			StartYear:            experience.StartYear,
			CurrentlyWorkingHere: experience.CurrentlyWorkingHere,
			SellerID:             seller.ID,
		}

		if experience.CurrentlyWorkingHere {
			exp.EndYear = uint(time.Now().Year())
		} else {
			if experience.StartYear > experience.EndYear {
				return &types.SellerDTO{}, errors.New("StartYear is greater than EndYear. Which is forbidden")
			}
			exp.EndYear = experience.EndYear
		}
		result = tx.
			Model(&types.Experience{}).
			Create(&exp)
		if result.Error != nil {
			tx.Rollback()
			return &types.SellerDTO{}, fmt.Errorf("Error saving experiences %v", result.Error)
		}

		experiences = append(experiences, types.ExperienceDTO{
			ID:                   exp.ID,
			Title:                exp.Title,
			Company:              exp.Company,
			Description:          exp.Description,
			StartYear:            exp.StartYear,
			EndYear:              exp.EndYear,
			CurrentlyWorkingHere: exp.CurrentlyWorkingHere,
		})
	}

	result = tx.
		Model(&types.Buyer{}).
		Where("id", seller.BuyerID).
		Update("is_seller", true)
	if result.Error != nil {
		tx.Rollback()
		return &types.SellerDTO{}, fmt.Errorf("Error updating buyer is_seller status %v", result.Error)
	}

	result = tx.Commit()

	return &types.SellerDTO{
		ID:               seller.ID,
		FullName:         seller.FullName,
		Username:         sellerDataInBuyerDB.Username,
		Email:            sellerDataInBuyerDB.Email,
		Country:          sellerDataInBuyerDB.Country,
		Bio:              seller.Bio,
		RatingSum:        seller.RatingSum,
		RatingsCount:     seller.RatingsCount,
		RatingCategories: seller.RatingCategories,
		ProfilePicture:   sellerDataInBuyerDB.ProfilePicture,
		Languages:        data.Languages,
		Skills:           data.Skills,
		Educations:       educations,
		Certificates:     certificates,
		Experiences:      experiences,
	}, result.Error
}

func (ss *SellerService) Update(ctx context.Context, updatedSellerData *types.Seller, data *types.UpdateSellerDTO) error {
	tx := ss.db.
		WithContext(ctx).
		Begin()

	var experiences []types.ExperienceDTO
	var educations []types.EducationDTO
	var certificates []types.CertificateDTO
	var result *gorm.DB

	//INFO: 1. UPDATE SELLER'S DATA
	result = tx.
		Save(updatedSellerData)
	if result.Error != nil {
		tx.Rollback()
		return fmt.Errorf("Error updating seller data. %v", result.Error)
	}

	//INFO: 2. UPDATE SELLER'S LANGUAGES
	result = tx.
		Model(&types.SellerLanguage{}).
		Delete(&types.SellerLanguage{}, "seller_id = ?", updatedSellerData.ID)

	if result.Error != nil {
		tx.Rollback()
		return fmt.Errorf("Error deleting sellerLanguages data. %v", result.Error)
	}

	for _, lang := range data.Languages {
		result = tx.
			Model(&types.SellerLanguage{}).
			Create(&types.SellerLanguage{
				SellerID:   updatedSellerData.ID,
				LanguageID: lang.ID,
			})
		if result.Error != nil {
			tx.Rollback()
			return fmt.Errorf("Error updating sellerLanguages data. %v", result.Error)
		}
	}

	//INFO: 3. UPDATE SELLER'S SKILLS
	result = tx.
		Model(&types.SellerSkill{}).
		Delete(&types.SellerSkill{}, "seller_id = ?", updatedSellerData.ID)

	if result.Error != nil {
		tx.Rollback()
		return fmt.Errorf("Error deleting sellerSkills data. %v", result.Error)
	}

	for _, skill := range data.Skills {
		result = tx.
			Model(&types.Skill{}).
			FirstOrCreate(&skill, "name = ?", skill.Name)
		if result.Error != nil {
			tx.Rollback()
			return fmt.Errorf("Error updating skills data. %v", result.Error)
		}

		result = tx.
			Model(&types.SellerSkill{}).
			Create(&types.SellerSkill{
				SellerID: updatedSellerData.ID,
				SkillID:  skill.ID,
			})
		if result.Error != nil {
			tx.Rollback()
			return fmt.Errorf("Error updating sellerSkills data. %v", result.Error)

		}
	}

	//INFO: 4. UPDATE SELLER'S EDUCATIONS
	result = tx.
		Model(&types.Education{}).
		Delete(&types.Education{}, "seller_id = ?", updatedSellerData.ID)

	if result.Error != nil {
		tx.Rollback()
		return fmt.Errorf("Error deleting educations data. %v", result.Error)
	}

	for _, education := range data.Educations {
		edu := types.Education{
			Title:      education.Title,
			Major:      education.Major,
			Country:    education.Country,
			University: education.University,
			SellerID:   updatedSellerData.ID,
		}
		result = tx.
			Model(&types.Education{}).
			Create(&edu)
		if result.Error != nil {
			tx.Rollback()
			return fmt.Errorf("Error updating sellerEducations data. %v", result.Error)
		}

		educations = append(educations, types.EducationDTO{
			ID:         edu.ID,
			Title:      edu.Title,
			Major:      edu.Major,
			University: edu.University,
			Country:    edu.Country,
		})
	}

	//INFO: 5. UPDATE SELLER'S CERTIFICATES
	result = tx.
		Model(&types.Certificate{}).
		Delete(&types.Certificate{}, "seller_id = ?", updatedSellerData.ID)

	if result.Error != nil {
		tx.Rollback()
		return fmt.Errorf("Error deleting certificates data. %v", result.Error)
	}

	for _, certif := range data.Certificates {
		cert := types.Certificate{
			Name:     certif.Name,
			From:     certif.From,
			Year:     certif.Year,
			SellerID: updatedSellerData.ID,
		}
		result = tx.
			Model(&types.Certificate{}).
			Create(&cert)
		if result.Error != nil {
			tx.Rollback()
			return fmt.Errorf("Error updating sellerCertificates data. %v", result.Error)
		}

		certificates = append(certificates, types.CertificateDTO{
			ID:   cert.ID,
			Name: cert.Name,
			From: cert.From,
			Year: cert.Year,
		})
	}

	//INFO: 6. UPDATE SELLER'S EXPERIENCES
	result = tx.
		Model(&types.Experience{}).
		Delete(&types.Experience{}, "seller_id = ?", updatedSellerData.ID)

	if result.Error != nil {
		tx.Rollback()
		return fmt.Errorf("Error deleting experience data. %v", result.Error)
	}

	for _, experience := range data.Experiences {
		exp := types.Experience{
			Title:                experience.Title,
			Description:          experience.Description,
			Company:              experience.Company,
			StartYear:            experience.StartYear,
			CurrentlyWorkingHere: experience.CurrentlyWorkingHere,
			SellerID:             updatedSellerData.ID,
		}

		if experience.CurrentlyWorkingHere {
			exp.EndYear = uint(time.Now().Year())
		} else {
			exp.EndYear = experience.EndYear
		}
		result = tx.
			Model(&types.Experience{}).
			Create(&exp)
		if result.Error != nil {
			tx.Rollback()
			return fmt.Errorf("Error updating sellerExperiences data. %v", result.Error)
		}

		experiences = append(experiences, types.ExperienceDTO{
			ID:                   exp.ID,
			Title:                exp.Title,
			Company:              exp.Company,
			Description:          exp.Description,
			StartYear:            exp.StartYear,
			EndYear:              exp.EndYear,
			CurrentlyWorkingHere: exp.CurrentlyWorkingHere,
		})
	}

	result = tx.Commit()

	return result.Error
}

func (ss *SellerService) UpdateBalance(ctx context.Context, sellerID string, addedBalance uint64) (*types.SellerIncBalanceDTO, error) {
	type Result struct {
		types.Seller
		Email string `json:"email" gorm:"email"`
	}
	var resultData Result

	result := ss.db.
		Debug().
		WithContext(ctx).
		Raw(`
            UPDATE sellers
            SET account_balance = account_balance + ?
            FROM buyers
            WHERE sellers.id = ? AND buyers.id = sellers.buyer_id
            RETURNING sellers.*, buyers.email AS email
        `, addedBalance, sellerID).
		Scan(&resultData)

	if resultData.Email == "" {
		return nil, fmt.Errorf("Seller data is not found")
	}

	return &types.SellerIncBalanceDTO{
		ID:               resultData.Seller.ID,
		FullName:         resultData.Seller.FullName,
		Email:            resultData.Email,
		Bio:              resultData.Seller.Bio,
		AccountBalance:   resultData.Seller.AccountBalance,
		RatingSum:        resultData.Seller.RatingSum,
		RatingsCount:     resultData.Seller.RatingsCount,
		RatingCategories: resultData.Seller.RatingCategories,
		StripeAccountID:  resultData.Seller.StripeAccountID,
	}, result.Error
}
