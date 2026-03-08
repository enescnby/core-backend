package websocket

import (
	"core-backend/internal/models"
	"core-backend/internal/repositories"
	"core-backend/pb"
	"core-backend/pkg/logger"
	"fmt"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type ConnectionManager interface {
	Register(userID string, conn *websocket.Conn)
	ReadPump(userID string, conn *websocket.Conn)
	Unregister(userID string)
	SendToUser(receiverID string, payload []byte) error
}

type connectionManager struct {
	clients map[string]*websocket.Conn
	mu      sync.RWMutex
	repo    repositories.MessageRepository
}

func NewConnectionManager(repo repositories.MessageRepository) ConnectionManager {
	return &connectionManager{
		clients: make(map[string]*websocket.Conn),
		repo:    repo,
	}
}

func (m *connectionManager) Register(userID string, conn *websocket.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.clients[userID] = conn

	logger.Log.Info("user connected to WebSocket", zap.String("core_guard_id", userID))
}

func (m *connectionManager) ReadPump(userID string, conn *websocket.Conn) {
	for {
		messageType, rawPayload, err := conn.ReadMessage()
		if err != nil {
			logger.Log.Info("user tunnel disconnected", zap.String("core_guard_id", userID))
			break
		}

		if messageType != websocket.BinaryMessage {
			continue
		}

		var wrapper pb.WebSocketMessage
		if err := proto.Unmarshal(rawPayload, &wrapper); err != nil {
			logger.Log.Error("Protobuf can not decode, broken data", zap.Error(err))
			continue
		}

		switch msg := wrapper.Content.(type) {
		case *pb.WebSocketMessage_Payload:
			payload := msg.Payload
			receiverID := payload.ReceiverId

			err := m.SendToUser(receiverID, rawPayload)
			if err != nil {
				msgUUID, _ := uuid.Parse(payload.MessageId)
				senderUUID, _ := uuid.Parse(payload.SenderId)
				receiverUUID, _ := uuid.Parse(payload.ReceiverId)

				offlineMsg := &models.EncryptedMessages{
					MessageID:   msgUUID,
					SenderID:    senderUUID,
					ReceiverID:  receiverUUID,
					Ciphertext:  payload.Ciphertext,
					Nonce:       payload.Nonce,
					AuthTag:     payload.AuthTag,
					MessageType: int(payload.Type),
				}

				saveErr := m.repo.SaveMessage(offlineMsg)
				if saveErr == nil {
					logger.Log.Info("user is offline, message saved", zap.String("msg_id", payload.MessageId))
				}
			}
		case *pb.WebSocketMessage_Receipt:
			receipt := msg.Receipt
			receiverID := receipt.ReceiverId

			if receipt.Status == pb.ReceiptStatus_DELIVERED {
				msgUUID, _ := uuid.Parse(receipt.MessageId)
				_ = m.repo.MarkAsDelivered([]uuid.UUID{msgUUID})
			}

			_ = m.SendToUser(receiverID, rawPayload)
		}
	}
}

func (m *connectionManager) Unregister(userID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if conn, exists := m.clients[userID]; exists {
		_ = conn.Close()

		delete(m.clients, userID)

		logger.Log.Info("connection closed, user cleaned from RAM", zap.String("core_guard_id", userID))
	}
}

func (m *connectionManager) SendToUser(receiverID string, payload []byte) error {
	m.mu.RLock()
	conn, exists := m.clients[receiverID]

	m.mu.RUnlock()

	if !exists {
		return fmt.Errorf("user %s is offline", receiverID)
	}

	err := conn.WriteMessage(websocket.BinaryMessage, payload)
	if err != nil {
		return fmt.Errorf("failed to send message to user %s: %w", receiverID, err)
	}

	return nil
}
