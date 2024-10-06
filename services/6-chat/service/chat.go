package service

import (
	"context"
	"log"
	"time"

	"github.com/Akihira77/gojobber/services/6-chat/types"
	"gorm.io/gorm"
)

type ChatService struct {
	db *gorm.DB
}

type ChatServiceImpl interface {
	GetAllMyConversations(ctx context.Context, userID string) ([]types.UserConversationDTO, error)
	GetMessages(ctx context.Context, conversationID string) ([]types.MessageDTO, error)
	InsertMessage(ctx context.Context, senderID string, data *types.CreateMessageDTO) (*types.Message, error)
	CalculateUnreadMessages(ctx context.Context, conversationID, senderID string) int
}

func NewChatService(db *gorm.DB) ChatServiceImpl {
	return &ChatService{
		db: db,
	}
}

func (cs *ChatService) GetMessages(ctx context.Context, conversationID string) ([]types.MessageDTO, error) {
	var messages []types.MessageDTO
	result := cs.db.
		Debug().
		WithContext(ctx).
		Model(&types.Message{}).
		Where("conversation_id = ?", conversationID).
		Find(&messages)

	return messages, result.Error

}

func (cs *ChatService) GetAllMyConversations(ctx context.Context, userID string) ([]types.UserConversationDTO, error) {
	subQuery := cs.db.
		Model(&types.Message{}).
		Select("COUNT(CASE WHEN messages.unread = TRUE THEN 1 END)").
		Where("messages.conversation_id = conversations.id").
		Where("messages.sender_id <> ?", userID)

	if subQuery.Error != nil {
		return nil, subQuery.Error
	}

	var conversations []types.UserConversationDTO
	result := cs.db.
		Debug().
		WithContext(ctx).
		Model(&types.Conversation{}).
		Select("conversations.id, conversations.user_one_id, conversations.user_two_id, (?) AS unread_messages", subQuery).
		Where("conversations.user_one_id = ? OR conversations.user_two_id = ?", userID, userID).
		Scan(&conversations)

	return conversations, result.Error
}

func (cs *ChatService) InsertMessage(ctx context.Context, senderID string, data *types.CreateMessageDTO) (*types.Message, error) {
	tx := cs.db.
		Debug().
		WithContext(ctx).
		Begin()

	conv := types.Conversation{
		UserOneID: senderID,
		UserTwoID: data.ReceiverID,
	}
	result := tx.
		Model(&types.Conversation{}).
		FirstOrCreate(
			&conv,
			`(user_one_id = ? AND user_two_id = ?) OR (user_one_id = ? AND user_two_id = ?)`,
			senderID, data.ReceiverID, senderID, data.ReceiverID)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}

	// mark all unread messages as read
	result = tx.
		Model(&types.Message{}).
		Update("unread", false).
		Where("messages.sender_id = ? AND messages.conversation_id = ?", data.ReceiverID, conv.ID)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}

	msg := types.Message{
		SenderID:       senderID,
		Body:           data.Body,
		FileURL:        data.FileURL,
		Offer:          data.Offer,
		Unread:         true,
		CreatedAt:      time.Now(),
		ConversationID: conv.ID.String(),
	}

	result = tx.
		Model(&types.Message{}).
		Create(&msg)
	if result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}

	result = tx.Commit()
	return &msg, result.Error
}

func (cs *ChatService) CalculateUnreadMessages(ctx context.Context, conversationID, senderID string) int {
	var count int64

	result := cs.db.
		Debug().
		WithContext(ctx).
		Model(&types.Message{}).
		Select("COUNT(CASE WHEN messages.unread = TRUE THEN 1 END) AS count").
		Where("messages.sender_id = ? AND messages.conversation_id = ?", senderID, conversationID).
		Scan(&count)

	if result.Error != nil {
		log.Println("CalculateUnreadMessages error", result.Error)
		return 0
	}

	return int(count)
}
