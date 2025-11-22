package session

import (
	"net"
	"time"

	"github.com/gorilla/websocket"
)

type WSCConnection struct {
	Conn *websocket.Conn
}

func (ws *WSCConnection) ReadMessage() ([]byte, error) {
	_, message, err := ws.Conn.ReadMessage()
	if err != nil {
		return nil, err
	}
	return message, nil
}

// 发送消息
func (ws *WSCConnection) WriteMessage(message []byte) error {
	return ws.Conn.WriteMessage(websocket.BinaryMessage, message)
}

// 关闭连接
func (ws *WSCConnection) Close() error {
	return ws.Conn.Close()
}

// 获取远程地址
func (ws *WSCConnection) RemoteAddr() net.Addr {
	return ws.Conn.RemoteAddr()
}

// 设置读取超时时间
func (ws *WSCConnection) SetReadDeadline(t time.Time) error {
	return ws.Conn.SetReadDeadline(t)
}

// 设置写入超时时间
func (ws *WSCConnection) SetWriteDeadline(t time.Time) error {
	return ws.Conn.SetWriteDeadline(t)
}
