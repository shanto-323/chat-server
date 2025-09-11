package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/gocql/gocql"
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
	createdAt := time.Now().Add(time.Millisecond)

	chatPacket := model.ChatPacket{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Message:    message,
		Offline:    offline,
		CreatedAt:  createdAt,
	}

	chat := model.Chat{
		ChatID:         gocql.TimeUUID(),
		ConversationID: conversation_id,
		Payload:        chatPacket,
	}

	if err := m.repo.InsertMessage(ctx, &chat); err != nil {
		return err
	}
	return nil
}

func (m *MessageService) GetMessage(ctx context.Context, senderId, receiverId string, createdAt time.Time) ([]*model.ChatPacket, error) {
	conversation_id := m.conversationIdGenerator(senderId, receiverId)
	createdAt = createdAt.Add(time.Millisecond)
	resp, err := m.repo.GetMessageFromBucket(ctx, conversation_id, createdAt)
	if err != nil {
		slog.Error("SERVICE", "getMessage", err.Error())
		return nil, err
	}

	return resp, nil
}

func (m *MessageService) GetLatestMessage(ctx context.Context, senderId, receiverId string) (*model.ChatPacket, error) {
	conversation_id := m.conversationIdGenerator(senderId, receiverId)
	resp, err := m.repo.GetLatestMessageFromBucket(ctx, conversation_id)
	if err != nil {
		slog.Error("SERVICE", "getMessage", err.Error())
		return nil, err
	}

	return resp, nil
}

func (m *MessageService) conversationIdGenerator(id_1, id_2 string) string {
	if id_1 > id_2 {
		return fmt.Sprintf("%s|%s", id_1, id_2)
	}
	return fmt.Sprintf("%s|%s", id_2, id_1)
}
