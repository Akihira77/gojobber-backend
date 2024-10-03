package types

import (
	"fmt"

	"gorm.io/gorm"
)

type RatingCategory struct {
	Five  uint `json:"five" gorm:"not null; default:0"`
	Four  uint `json:"four" gorm:"not null; default:0"`
	Three uint `json:"three" gorm:"not null; default:0"`
	Two   uint `json:"two" gorm:"not null; default:0"`
	One   uint `json:"one" gorm:"not null; default:0"`
}

type Language struct {
	ID       uint   `json:"id" gorm:"primaryKey;"`
	Language string `json:"language" gorm:"not null;"`
}

type Skill struct {
	ID   uint   `json:"id" gorm:"primaryKey;"`
	Name string `json:"name" gorm:"not null;"`
}

type Certificate struct {
	ID       uint   `json:"id" gorm:"primaryKey;"`
	Name     string `json:"name" gorm:"not null;"`
	From     string `json:"from" gorm:"not null;"`
	Year     uint   `json:"year" gorm:"not null;"`
	SellerID string `json:"sellerId"`
}

type CertificateDTO struct {
	ID   uint   `json:"id,omitempty"`
	Name string `json:"name" validate:"required"`
	From string `json:"from" validate:"required"`
	Year uint   `json:"year" validate:"required"`
}

type Education struct {
	ID         uint   `json:"id" gorm:"primaryKey;"`
	Title      string `json:"title" gorm:"not null;"`
	Major      string `json:"major" gorm:"not null;"`
	University string `json:"university" gorm:"not null;"`
	Country    string `json:"country" gorm:"not null;"`
	SellerID   string `json:"sellerId"`
}

type EducationDTO struct {
	ID         uint   `json:"id,omitempty"`
	Title      string `json:"title" validate:"required"`
	Major      string `json:"major" validate:"required"`
	University string `json:"university" validate:"required"`
	Country    string `json:"country" validate:"required"`
}

type Experience struct {
	ID                   uint   `json:"id" gorm:"primaryKey;"`
	Title                string `json:"title" gorm:"not null;"`
	Company              string `json:"company" gorm:"not null;"`
	Description          string `json:"description"`
	StartYear            uint   `json:"startYear" gorm:"not null;"`
	EndYear              uint   `json:"endYear"`
	CurrentlyWorkingHere bool   `json:"currentlyWorkingHere" gorm:"not null; default:false;"`
	SellerID             string `json:"sellerId"`
}

type ExperienceDTO struct {
	ID                   uint   `json:"id,omitempty"`
	Title                string `json:"title" validate:"required"`
	Company              string `json:"company" validate:"required"`
	Description          string `json:"description" validate:"required"`
	StartYear            uint   `json:"startYear" validate:"required"`
	EndYear              uint   `json:"endYear" validate:"required"`
	CurrentlyWorkingHere bool   `json:"currentlyWorkingHere" validate:"required"`
}

type CreateSellerDTO struct {
	FullName     string           `json:"fullName" validate:"required"`
	Bio          string           `json:"bio" validate:"required"`
	Languages    []Language       `json:"languages" validate:"required"`
	Skills       []Skill          `json:"skills" validate:"required,min=1"`
	Certificates []CertificateDTO `json:"certificates" validate:"required,min=1"`
	Educations   []EducationDTO   `json:"educations" validate:"required,min=1"`
	Experiences  []ExperienceDTO  `json:"experiences" validate:"required,min=1"`
}

type Seller struct {
	ID      string `json:"id" gorm:"primaryKey;"`
	BuyerID string `json:"buyerId" gorm:"unique;not null;"`
	// Buyer            Buyer          `json:"buyer,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	FullName         string         `json:"fullName" gorm:"not null;"`
	Bio              string         `json:"bio" gorm:"not null;"`
	RatingsCount     uint64         `json:"ratingsCount" gorm:"not null;"`
	RatingSum        uint64         `json:"ratingSum" gorm:"not null;"`
	RatingCategories RatingCategory `json:"ratingCategories" gorm:"type:jsonb;not null;serializer:json;"`
}

type TestSeller struct {
	ID               string           `json:"id" gorm:"primaryKey;"`
	BuyerID          string           `json:"buyerId"`
	Buyer            Buyer            `json:"buyer,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	FullName         string           `json:"fullName" gorm:"not null;"`
	Bio              string           `json:"bio" gorm:"not null;"`
	RatingsCount     uint64           `json:"ratingsCount" gorm:"not null;"`
	RatingSum        uint64           `json:"ratingSum" gorm:"not null;"`
	RatingCategories RatingCategory   `json:"ratingCategories" gorm:"type:jsonb;not null;serializer:json;"`
	Languages        []SellerLanguage `json:"languages" gorm:"foreignKey:SellerID;references:ID"`
	Skills           []SellerSkill    `json:"skills" gorm:"foreignKey:SellerID;references:ID"`
	Certificates     []Certificate    `json:"certificates" gorm:"foreignKey:SellerID;references:ID"`
	Educations       []Education      `json:"educations" gorm:"foreignKey:SellerID;references:ID"`
	Experiences      []Experience     `json:"experiences" gorm:"foreignKey:SellerID;references:ID"`
}

type SellerOverview struct {
	FullName         string         `json:"fullName" gorm:"column:full_name"`
	RatingsCount     uint64         `json:"ratingsCount" gorm:"column:ratings_count"`
	RatingSum        uint64         `json:"ratingSum" gorm:"column:rating_sum"`
	RatingCategories RatingCategory `json:"ratingCategories" gorm:"serializer:json;column:rating_categories;"`
}

type SellerDTO struct {
	ID               string           `json:"id"`
	FullName         string           `json:"fullName"`
	Username         string           `json:"username,omitempty"`
	Email            string           `json:"email"`
	Country          string           `json:"country"`
	Bio              string           `json:"bio"`
	ProfilePicture   string           `json:"profilePicture"`
	Languages        []Language       `json:"languages"`
	Skills           []Skill          `json:"skills"`
	Certificates     []CertificateDTO `json:"certificates"`
	Educations       []EducationDTO   `json:"educations"`
	Experiences      []ExperienceDTO  `json:"experiences"`
	RatingsCount     uint64           `json:"ratingsCount"`
	RatingSum        uint64           `json:"ratingSum"`
	RatingCategories RatingCategory   `json:"ratingCategories" gorm:"serializer:json"`
}

type UpdateSellerDTO struct {
	FullName     string           `json:"fullName" validate:"required"`
	Bio          string           `json:"bio" validate:"required"`
	Languages    []Language       `json:"languages" validate:"required"`
	Skills       []Skill          `json:"skills" validate:"required"`
	Certificates []CertificateDTO `json:"certificates" validate:"required"`
	Educations   []EducationDTO   `json:"educations" validate:"required"`
	Experiences  []ExperienceDTO  `json:"experiences" validate:"required"`
}

func ApplyDBSetup(db *gorm.DB) error {
	err := db.
		Debug().
		Exec(`
		ALTER TABLE seller_languages
		ADD FOREIGN KEY (seller_id) REFERENCES sellers(id) ON UPDATE CASCADE ON DELETE CASCADE;
		`).Error
	if err != nil {
		fmt.Println("Error applying foreign key to seller_languages", err)
		return err
	}

	err = db.
		Debug().
		Exec(`
		ALTER TABLE seller_languages
		ADD FOREIGN KEY (language_id) REFERENCES languages(id) ON UPDATE CASCADE ON DELETE RESTRICT;
		`).Error
	if err != nil {
		fmt.Println("Error applying foreign key to seller_languages", err)
		return err
	}

	err = db.
		Debug().
		Exec(`
		ALTER TABLE seller_skills
		ADD FOREIGN KEY (seller_id) REFERENCES sellers(id) ON UPDATE CASCADE ON DELETE CASCADE;
		`).Error
	if err != nil {
		fmt.Println("Error applying foreign key to seller_skills", err)
		return err
	}

	err = db.
		Debug().
		Exec(`
		ALTER TABLE seller_skills
		ADD FOREIGN KEY (skill_id) REFERENCES skills(id) ON UPDATE CASCADE ON DELETE SET NULL;
		`).Error
	if err != nil {
		fmt.Println("Error applying foreign key to seller_skills", err)
		return err
	}

	err = db.
		Debug().
		Exec(`
		ALTER TABLE experiences
		ADD FOREIGN KEY (seller_id) REFERENCES sellers(id) ON UPDATE CASCADE ON DELETE CASCADE;
		`).Error
	if err != nil {
		fmt.Println("Error applying foreign key to experiences", err)
		return err
	}

	err = db.
		Debug().
		Exec(`
		ALTER TABLE educations
		ADD FOREIGN KEY (seller_id) REFERENCES sellers(id) ON UPDATE CASCADE ON DELETE CASCADE;
		`).Error
	if err != nil {
		fmt.Println("Error applying foreign key to educations", err)
		return err
	}

	err = db.
		Debug().
		Exec(`
		ALTER TABLE certificates
		ADD FOREIGN KEY (seller_id) REFERENCES sellers(id) ON UPDATE CASCADE ON DELETE CASCADE;
		`).Error
	if err != nil {
		fmt.Println("Error applying foreign key to certificates", err)
		return err
	}

	err = db.
		Debug().
		Exec(`
		ALTER TABLE sellers
		ADD FOREIGN KEY (buyer_id) REFERENCES buyers(id) ON UPDATE CASCADE ON DELETE CASCADE;
		`).Error
	if err != nil {
		fmt.Println("Error applying foreign key to certificates", err)
		return err
	}

	return nil
}

type SellerLanguage struct {
	SellerID   string `json:"sellerId"`
	LanguageID uint   `json:"languageId"`
}

type SellerLanguageDTO struct {
	SellerID   string `json:"sellerId"`
	LanguageID uint   `json:"languageId"`
	Language   string `json:"language"`
}

type SellerSkill struct {
	SellerID string `json:"sellerId"`
	SkillID  uint   `json:"skillId"`
}

type SellerSkillDTO struct {
	SellerID string `json:"sellerId"`
	SkillID  uint   `json:"skillId"`
	Skill    string `json:"skill"`
}
