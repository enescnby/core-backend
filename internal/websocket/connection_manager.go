package websocket

import "github.com/gofiber/contrib/websocket"

type ConnectionManager interface {
	Register(userID string, conn *websocket.Conn)
	ReadPump(userID string, conn *websocket.Conn)
	Unregister(userID string)
	SendToUser(receiverID string, payload []byte) error
}

type connectionManager struct {
	//TODO
}

func NewConnectionManager() ConnectionManager {
	return &connectionManager{}
}

func (m *connectionManager) Register(userID string, conn *websocket.Conn) {
	//TODO
}

func (m *connectionManager) ReadPump(userID string, conn *websocket.Conn) {
	//TODO
}

func (m *connectionManager) Unregister(userID string) {
	//TODO
}

func (m *connectionManager) SendToUser(receiverID string, payload []byte) error {
	//TODO
	return nil
}
