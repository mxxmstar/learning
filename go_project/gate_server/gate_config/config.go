package gate_config

import (
	"fmt"
	"time"

	"github.com/mxxmstar/learning/pkg/config"
)

type Config struct {
	// 服务器配置
	ServerConfig *config.ServerConfig `mapstructure:"server"`
	// GateServer 特有配置
	WebSocketConfig WebSocketConfig `mapstructure:"websocket_config"`
}

type WebSocketConfig struct {
	AuthTimeout     time.Duration `mapstructure:"auth_timeout"`      // 验证 token/session 超时时间 (ValidateTokenOrSession)
	ReadBufferSize  int           `mapstructure:"read_buffer_size"`  // 读缓冲区大小
	WriteBufferSize int           `mapstructure:"write_buffer_size"` // 写缓冲区大小
	PingWait        time.Duration `mapstructure:"ping_wait"`         // ping 超时
	PongWait        time.Duration `mapstructure:"pong_wait"`         // pong 超时
	MaxMessageSize  int           `mapstructure:"max_message_size"`  // 最大消息长度
}

func Init() (*Config, error) {
	baseCfg, err := config.Init()
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		ServerConfig: &baseCfg.Server,
		WebSocketConfig: WebSocketConfig{
			AuthTimeout:     5 * time.Second,
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			PingWait:        60 * time.Second,
			PongWait:        60 * time.Second,
			MaxMessageSize:  1024 * 1024, // 1M
		},
	}
	return cfg, nil
}

func (c *Config) GetGateServer(serverName string) *config.GateServerConfig {
	for _, server := range c.ServerConfig.GateServers {
		if server.Name == serverName {
			return &server
		}
	}
	return nil
}

func (c *Config) GetAllGateServersAddress() []string {
	var gateServerAddress []string
	for _, server := range c.ServerConfig.GateServers {
		gateServerAddress = append(gateServerAddress, fmt.Sprintf("%s:%d", server.GRPCConfig.Host, server.GRPCConfig.Port))
	}
	return gateServerAddress
}

func (c *Config) GetVerifyGRPCAddress(serverName string) string {
	for _, server := range c.ServerConfig.VerifyServers {
		if server.Name == serverName {
			return fmt.Sprintf("%s:%d", server.GRPCConfig.Host, server.GRPCConfig.Port)
		}
	}
	return ""
}

func (c *Config) GetVerifyHttpAddress(serverName string) string {
	for _, server := range c.ServerConfig.VerifyServers {
		if server.Name == serverName {
			return fmt.Sprintf("http://%s:%d", server.HttpConfig.Host, server.HttpConfig.Port)
		}
	}
	return ""
}

func (c *Config) GetActiveVerifyServers() []config.VerifyServerConfig {
	var activeVerifyServers []config.VerifyServerConfig
	for _, server := range c.ServerConfig.VerifyServers {
		if server.ServiceConfig.Status == "active" {
			activeVerifyServers = append(activeVerifyServers, server)
		}
	}
	return activeVerifyServers
}

func (c *Config) GetStatusServerAddress() string {
	return fmt.Sprintf("%s:%d", c.ServerConfig.StatusServer.HttpConfig.Host, c.ServerConfig.StatusServer.HttpConfig.Port)
}
