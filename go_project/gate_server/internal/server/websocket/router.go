package websocket

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mxxmstar/learning/gate_server/internal/conn"
)

// Router 消息路由器，负责将不同类型的消息路由到对应的处理函数
type Router struct {
	handlers map[string]MessageHandler
}

func NewRouter() *Router {
	return &Router{
		handlers: make(map[string]MessageHandler),
	}
}

func (r *Router) RegisterHandler(msgType string, handler MessageHandler) {
	r.handlers[msgType] = handler
}

func (r *Router) Route(ctx context.Context, conn conn.Connection, envelope *Envelope) error {
	if handler, ok := r.handlers[envelope.Type]; ok {
		return handler.HandleMessage(ctx, conn, envelope)
	}

	// 发送未知消息类型响应到客户端
	rsp := map[string]interface{}{
		"type":    "unknown_message_type",
		"message": fmt.Sprintf("unknown message type: %s", envelope.Type),
	}
	rspBytes, err := json.Marshal(rsp)
	if err != nil {
		return err
	}
	return conn.Send(rspBytes)
}
