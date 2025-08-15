package database

import (
	"context"
	"strconv"
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

func (m *MessageService) PushMessage(ctx context.Context, c *model.Chat) error {
	id := m.chatId(c.SenderId, c.ReceiverId)

	return m.repo.InsertMessage(ctx, id, c.Message, c.CreatedAt)
}

func (m *MessageService) GetMessage(ctx context.Context, senderId, receiverId uint, createdAt time.Time) ([]*model.Chat, error) {
	id := m.chatId(senderId, receiverId)
	var chats []*model.Chat

	resp, err := m.repo.GetMessageFromBucket(ctx, id, createdAt)
	if err != nil {
		return nil, err
	}

	for _, r := range resp {
		chat := &model.Chat{
			SenderId:   senderId,
			ReceiverId: receiverId,
			Message:    r,
		}

		chats = append(chats, chat)
	}

	return chats, nil
}

func (m *MessageService) chatId(id_1, id_2 uint) string {
	id1 := strconv.FormatUint(uint64(id_1), 10)
	id2 := strconv.FormatUint(uint64(id_2), 10)
	// ChatId = bigId + | + smallId
	if id_1 > id_2 {
		return id1 + id2
	}
	return id2 + "|" + id1
}
