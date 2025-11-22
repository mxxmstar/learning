// pkg/session/tcp_connection.go
package session

import (
	"bufio"
	"net"
	"time"
)

// TCPConnection TCP连接适配器
type TCPConnection struct {
	Conn   *net.TCPConn
	reader *bufio.Reader
}

func (tcp *TCPConnection) ReadMessage() ([]byte, error) {
	if tcp.reader == nil {
		tcp.reader = bufio.NewReader(tcp.Conn)
	}

	// 这里假设使用简单的长度前缀协议
	// 实际应用中可能需要根据具体协议实现
	// 例如读取4字节长度，再读取对应长度的数据
	// 此处仅为示例
	return tcp.reader.ReadBytes('\n')
}

func (tcp *TCPConnection) WriteMessage(data []byte) error {
	_, err := tcp.Conn.Write(data)
	return err
}

func (tcp *TCPConnection) Close() error {
	return tcp.Conn.Close()
}

func (tcp *TCPConnection) RemoteAddr() net.Addr {
	return tcp.Conn.RemoteAddr()
}

func (tcp *TCPConnection) SetReadDeadline(t time.Time) error {
	return tcp.Conn.SetReadDeadline(t)
}

func (tcp *TCPConnection) SetWriteDeadline(t time.Time) error {
	return tcp.Conn.SetWriteDeadline(t)
}

// NewTCPConnection 创建TCP连接适配器
func NewTCPConnection(conn *net.TCPConn) *TCPConnection {
	return &TCPConnection{Conn: conn}
}
