package conn

import "errors"

var (
	ErrInvalidConnection  = errors.New("invalid connection")
	ErrConnectionNotFound = errors.New("connection not found")
	ErrConnectionClosed   = errors.New("connection closed")
)

// Connection 定义连接接口
// 连接接口定义了连接的基本操作，包括获取连接ID、用户ID、发送消息和关闭连接
type Connection interface {
	ID() string
	UserID() uint64
	Send(msg []byte) error
	Close(reason string) error
}

// ConnectionManager 定义连接管理器接口
// 连接管理器接口定义了连接的注册、注销、获取连接和根据用户ID获取连接的操作
type ConnectionManager interface {
	Register(conn Connection) error
	UnRegister(conn Connection) error
	GetConnection(connID string) (Connection, error)
	GetConnectionsByUserID(userID uint64) []Connection
}
