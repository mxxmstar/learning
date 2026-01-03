package config

// 服务器总配置
type ServerConfig struct {
	GateServers   []GateServerConfig   `mapstructure:"gate_servers"`   // 网关服务配置(集群)
	VerifyServers []VerifyServerConfig `mapstructure:"verify_servers"` // 验证服务配置(集群)
	StatusServer  StatusServerConfig   `mapstructure:"status_server"`  // 状态服务配置(单例)
	GlobalConfig  GlobalConfig         `mapstructure:"global_config"`  // 全局配置
}

// GateServerConfig 网关服务配置
type GateServerConfig struct {
	Name          string               `mapstructure:"name"`
	GRPCConfig    GateGRPCServerConfig `mapstructure:"grpc_config"`
	HttpConfig    GateHttpServerConfig `mapstructure:"http_config"`
	ServiceConfig ServiceConfig        `mapstructure:"service_config"`
}

type GateGRPCServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type GateHttpServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// VerifyServerConfig 验证服务配置
type VerifyServerConfig struct {
	Name          string                 `mapstructure:"name"`           // 服务名称
	GRPCConfig    VerifyGRPCServerConfig `mapstructure:"grpc_config"`    // grpc 服务配置
	HttpConfig    VerifyHttpServerConfig `mapstructure:"http_config"`    // http 服务配置
	ServiceConfig ServiceConfig          `mapstructure:"service_config"` // 服务配置(日志，负载均衡，集群)
}

type VerifyGRPCServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type VerifyHttpServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type StatusServerConfig struct {
	Name          string                 `mapstructure:"name"`
	GRPCConfig    StatusGRPCServerConfig `mapstructure:"grpc_config"`
	HttpConfig    StatusHttpServerConfig `mapstructure:"http_config"`
	ServiceConfig ServiceConfig          `mapstructure:"service_config"`
}

type StatusGRPCServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type StatusHttpServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// 服务配置 包含集群、区域、权重、负载均衡、日志、等服务配置
type ServiceConfig struct {
	ClusterId     string            `mapstructure:"cluster_id"`     // 集群 id
	Region        string            `mapstructure:"region"`         // 地理区域
	Zone          string            `mapstructure:"zone"`           // 可用区域
	Weight        int               `mapstructure:"weight"`         // 权重
	Status        string            `mapstructure:"status"`         // 状态
	ServicePrefix string            `mapstructure:"service_prefix"` // 服务前缀
	LoadBalance   LoadBalanceConfig `mapstructure:"load_balance"`   // 负载均衡配置
	LogConfig     LogConfig         `mapstructure:"log_config"`     // 日志配置
}

// 全局配置 仅包含全局默认值
type GlobalConfig struct {
	Env             string `mapstructure:"env"`
	DefaultLogLevel string `mapstructure:"default_log_level"` // 默认日志级别
}

// 负载均衡配置
type LoadBalanceConfig struct {
	Method              string `mapstructure:"method"`                // 负载均衡方法
	HealthCheck         bool   `mapstructure:"health_check"`          // 是否开启健康检查
	HealthCheckPath     string `mapstructure:"health_check_path"`     // 健康检查路径
	HealthCheckInterval int    `mapstructure:"health_check_interval"` // 健康检查间隔
	HealthCheckPort     int    `mapstructure:"health_check_port"`     // 健康检查端口
}

// 日志配置
type LogConfig struct {
	Level      string        `mapstructure:"level"`       // 日志级别 debug/info/warn/error
	Format     string        `mapstructure:"format"`      // 日志格式 json/text
	Output     string        `mapstructure:"output"`      // 日志输出 stdout/file/both
	FileConfig LogFileConfig `mapstructure:"file_config"` // 日志文件配置
}

// 日志文件配置
type LogFileConfig struct {
	FileName   string `mapstructure:"file_name"`   // 文件名
	FilePath   string `mapstructure:"file_path"`   // 文件路径
	MaxSize    int    `mapstructure:"max_size"`    // 文件大小
	MaxBackups int    `mapstructure:"max_backups"` // 备份数量
	MaxAge     int    `mapstructure:"max_age"`     // 文件保存天数
	Compress   bool   `mapstructure:"compress"`    // 是否压缩
}
