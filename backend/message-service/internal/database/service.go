package database

import (
	"context"
	"fmt"
	"time"

	"github.com/shanto-323/Chat-Server-1/message-service/internal/database/model"
)

type MessageService struct {
	repo MessageRepository
}

func NewUserService(repo MessageRepository) *MessageService {
	return &MessageService{
		repo: repo,
	}
}

func (m *MessageService) PushMessage(ctx context.Context, senderID, receiverID, message string, offline bool) error {
	conversation_id := m.conversationIdGenerator(senderID, receiverID)
	createdAt := time.Now()
	chat := model.Chat{
		ConversationID: conversation_id,
		SenderID:       senderID,
		ReceiverID:     receiverID,
		Message:        message,
		CreatedAt:      createdAt,
		Offline:        offline,
	}

	return m.repo.InsertMessage(ctx, &chat)
}

func (m *MessageService) GetMessage(ctx context.Context, senderId, receiverId string, createdAt time.Time) ([]*model.Chat, error) {
	conversation_id := m.conversationIdGenerator(senderId, receiverId)
	var chats []*model.Chat

	resp, err := m.repo.GetMessageFromBucket(ctx, conversation_id, createdAt)
	if err != nil {
		return nil, err
	}

	for _, r := range resp {
		chat := &model.Chat{
			SenderID:   senderId,
			ReceiverID: receiverId,
			Message:    r,
		}

		chats = append(chats, chat)
	}

	return chats, nil
}

func (m *MessageService) conversationIdGenerator(id_1, id_2 string) string {
	if id_1 > id_2 {
		return fmt.Sprintf("%s|%s", id_1, id_2)
	}
	return fmt.Sprintf("%s|%s", id_2, id_1)
}
