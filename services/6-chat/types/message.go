package types

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OfferStatus string

const (
	PENDING  OfferStatus = "PENDING"
	CANCELED OfferStatus = "CANCELED"
	ACCEPTED OfferStatus = "ACCEPTED"
)

func (p *OfferStatus) Scan(value interface{}) error {
	*p = OfferStatus(value.([]byte))
	return nil
}

func (p OfferStatus) Value() (driver.Value, error) {
	return string(p), nil
}

type Offer struct {
	GigTitle             string      `json:"gigTitle" form:"gigTitle" validate:"required"`
	Price                uint        `json:"price" form:"price" validate:"required,gt=0"`
	ExpectedDeliveryDays uint        `json:"expectedDeliveryDays" form:"expectedDeliveryDays" validate:"required,gt=0,lte=365"`
	Description          string      `json:"description" form:"description"`
	Status               OfferStatus `json:"status" gorm:"not null;"`
	CreatedAt            time.Time   `json:"createdAt" gorm:"not null;"`
}

func (o *Offer) Scan(value interface{}) error {
	if value == nil {
		*o = Offer{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("could not scan type %T into Offer", value)
	}
	return json.Unmarshal(bytes, &o)
}

func (o *Offer) Value() (driver.Value, error) {
	if o == nil {
		return nil, nil
	}
	return json.Marshal(o)
}

type Message struct {
	ID             uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4();"`
	Body           string    `json:"body,omitempty"`
	FileURL        string    `json:"fileUrl,omitempty"`
	Offer          *Offer    `json:"offer,omitempty" gorm:"type:jsonb;serializer:json;default:null;"`
	Unread         bool      `json:"unread" gorm:"not null;default:true"`
	CreatedAt      time.Time `json:"createdAt" gorm:"not null;"`
	ConversationID string    `json:"conversationId" gorm:"not null;"`
	SenderID       string    `json:"senderId"`
}

type MessageDTO struct {
	ID        uuid.UUID `json:"id"`
	Body      string    `json:"body,omitempty"`
	FileURL   string    `json:"fileUrl,omitempty"`
	Offer     *Offer    `json:"offer,omitempty" gorm:"type:jsonb;"`
	Unread    bool      `json:"unread"`
	CreatedAt time.Time `json:"createdAt"`
	SenderID  string    `json:"senderId"`
}

type CreateMessageDTO struct {
	ReceiverID    string         `json:"receiverId" form:"receiverId" validate:"required"`
	ReceiverEmail string         `json:"receiverEmail" validate:"required,email"`
	Body          string         `json:"body" form:"body"`
	File          multipart.File `json:"file" form:"file"`
	FileURL       string         `json:"fileUrl"`
	Offer         *Offer         `json:"offer,omitempty" form:"offer"`
}

type Conversation struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	UserOneID string    `json:"userOneId"`
	UserTwoID string    `json:"userTwoId"`
	Messages  []Message `json:"messages" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type UserConversationDTO struct {
	ID                   uuid.UUID `json:"id"`
	UserOneID            string    `json:"userOneId" gorm:"user_one_id"`
	UserTwoID            string    `json:"userTwoId" gorm:"user_two_id"`
	SenderName           string    `json:"senderName"`
	SenderEmail          string    `json:"senderEmail"`
	SenderProfilePicture string    `json:"senderProfilePicture"`
	UnreadMessages       int       `json:"unreadMessages" gorm:"unread_messages"`
}

type ConversationDTO struct {
	ID                    uuid.UUID `json:"id" gorm:"id"`
	UserOneID             string    `json:"userOneId" gorm:"user_one_id"`
	UserOneName           string    `json:"userOneName,omitempty" gorm:"user_one_name"`
	UserOneProfilePicture string    `json:"userOneProfilePicture,omitempty" gorm:"user_one_profile_picture"`
	UserTwoID             string    `json:"userTwoId" gorm:"user_two_id"`
	UserTwoName           string    `json:"userTwoName,omitempty" gorm:"user_two_name"`
	UserTwoProfilePicture string    `json:"userTwoProfilePicture,omitempty" gorm:"user_two_profile_picture"`
	LastMessage           string    `json:"lastMessage" gorm:"last_message"`
	MessageSentDate       time.Time `json:"messageSentDate" gorm:"message_sent_date"`
}

type ChatMessagesDTO struct {
	ConversationID        uuid.UUID `json:"conversationId" gorm:"conversation_id"`
	UserOneName           string    `json:"userOneName" gorm:"user_one_name"`
	UserOneProfilePicture string    `json:"userOneProfilePicture" gorm:"user_one_profile_picture"`
	UserTwoName           string    `json:"userTwoName" gorm:"user_two_name"`
	UserTwoProfilePicture string    `json:"userTwoProfilePicture" gorm:"user_two_profile_picture"`
	Messages              []Message `json:"messages"`
}

func ApplyDBSetup(db *gorm.DB) error {
	//INFO: FOREIGN KEY
	result := db.Debug().Exec(`
		ALTER TABLE conversations
		ADD FOREIGN KEY (user_one_id) REFERENCES auths(id) ON DELETE RESTRICT ON UPDATE CASCADE;
		`)
	if result.Error != nil {
		return result.Error
	}

	result = db.Debug().Exec(`
		ALTER TABLE conversations
		ADD FOREIGN KEY (user_two_id) REFERENCES auths(id) ON DELETE RESTRICT ON UPDATE CASCADE;
		`)
	if result.Error != nil {
		return result.Error
	}

	result = db.Debug().Exec(`
		ALTER TABLE messages
		ADD FOREIGN KEY (sender_id) REFERENCES auths(id) ON DELETE RESTRICT ON UPDATE CASCADE;
		`)
	if result.Error != nil {
		return result.Error
	}
	//
	// result = db.Debug().Exec(`
	// 	ALTER TABLE messages
	// 	ADD FOREIGN KEY (conversation_id) REFERENCES conversations(id) ON DELETE CASCADE ON UPDATE CASCADE;
	// 	`)
	// if result.Error != nil {
	// 	return result.Error
	// }

	//INFO: INDEXES
	result = db.Debug().Exec(`
		CREATE INDEX IF NOT EXISTS idx_sender_id
		ON messages USING btree (sender_id);
		`)
	if result.Error != nil {
		return result.Error
	}

	result = db.Debug().Exec(`
		CREATE INDEX IF NOT EXISTS idx_conversation_id
		ON messages USING btree (conversation_id);
		`)
	if result.Error != nil {
		return result.Error
	}

	result = db.Debug().Exec(`
		CREATE INDEX IF NOT EXISTS idx_user_one
		ON conversations USING btree (user_one_id);
		`)
	if result.Error != nil {
		return result.Error
	}

	result = db.Debug().Exec(`
		CREATE INDEX IF NOT EXISTS idx_user_two
		ON conversations USING btree (user_two_id);
		`)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
