// k:\code_mx\llfc_chat_go\learning\go_project\pkg\config\config_test.go
package config

import (
	"testing"
)

func TestConfigLoad(t *testing.T) {
	// 测试配置加载功能
	// t.Run("LoadConfig", func(t *testing.T) {
	// 	// 这里可以添加配置加载的测试逻辑
	// 	// 由于配置加载依赖外部配置文件，这里只做基本测试
	// 	cfg := &Config{}
	// 	if cfg == nil {
	// 		t.Error("Config should not be nil")
	// 	}
	// })

	t.Run("ServerConfig", func(t *testing.T) {
		// 测试服务器配置结构
		serverCfg := ServerConfig{
			GateServers: []GateServerConfig{
				{
					Name: "gate_server_1",
					GRPCConfig: GateGRPCServerConfig{
						Host: "localhost",
						Port: 50050,
					},
					HttpConfig: GateHttpServerConfig{
						Host: "localhost",
						Port: 8080,
					},
					ServiceConfig: ServiceConfig{
						ClusterId:     "cluster_1",
						Region:        "cn",
						Zone:          "guangzhou",
						Weight:        1,
						Status:        "online",
						ServicePrefix: "gate",
						LoadBalance: LoadBalanceConfig{
							Method:              "round_robin",
							HealthCheck:         true,
							HealthCheckPath:     "/health",
							HealthCheckInterval: 30,
							HealthCheckPort:     8080,
						},
						LogConfig: LogConfig{
							Level:  "debug",
							Format: "json",
							Output: "both",
							FileConfig: LogFileConfig{
								FileName:   "gate_server_1.log",
								FilePath:   "./logs",
								MaxSize:    100,
								MaxBackups: 3,
								MaxAge:     30,
								Compress:   true,
							},
						},
					},
				},
			},
			VerifyServers: []VerifyServerConfig{
				{
					Name: "verify_server_1",
					GRPCConfig: VerifyGRPCServerConfig{
						Host: "localhost",
						Port: 50060,
					},
					HttpConfig: VerifyHttpServerConfig{
						Host: "localhost",
						Port: 8090,
					},
					ServiceConfig: ServiceConfig{
						ClusterId:     "cluster_1",
						Region:        "cn",
						Zone:          "guangzhou",
						Weight:        1,
						Status:        "online",
						ServicePrefix: "verify",
						LoadBalance: LoadBalanceConfig{
							Method:              "round_robin",
							HealthCheck:         true,
							HealthCheckPath:     "/health",
							HealthCheckInterval: 30,
							HealthCheckPort:     8090,
						},
						LogConfig: LogConfig{
							Level:  "debug",
							Format: "json",
							Output: "both",
							FileConfig: LogFileConfig{
								FileName:   "test.log",
								FilePath:   "./logs",
								MaxSize:    100,
								MaxBackups: 3,
								MaxAge:     30,
								Compress:   true,
							},
						},
					},
				},
			},
			GlobalConfig: GlobalConfig{
				Env:             "test",
				DefaultLogLevel: "debug",
			},
		}

		if len(serverCfg.GateServers) != 1 {
			t.Errorf("Expected 1 gate server, got %d", len(serverCfg.GateServers))
		}

		if len(serverCfg.VerifyServers) != 1 {
			t.Errorf("Expected 1 verify server, got %d", len(serverCfg.VerifyServers))
		}

		if serverCfg.GateServers[0].Name != "gate_server_1" {
			t.Errorf("Expected gate server name 'gate_server_1', got %s", serverCfg.GateServers[0].Name)
		}

		if serverCfg.VerifyServers[0].Name != "verify_server_1" {
			t.Errorf("Expected verify server name 'verify_server_1', got %s", serverCfg.VerifyServers[0].Name)
		}
	})

	t.Run("DatabaseConfig", func(t *testing.T) {
		// 测试数据库配置结构
		dbCfg := DatabaseConfig{
			Type: "mysql",
			AuthDB: DatabaseInstanceConfig{
				DSN:          "root:root@tcp(localhost:3308)/auth_db?charset=utf8mb4&parseTime=True&loc=Local",
				CharSet:      "utf8mb4",
				ConnTimeout:  5000000000, // 5s in nanoseconds
				ReadTimeout:  3000000000, // 3s in nanoseconds
				WriteTimeout: 3000000000, // 3s in nanoseconds
				Pool: DatabasePoolConfig{
					MaxOpenConns:    50,
					MaxIdleConns:    10,
					ConnMaxLifetime: 300000000000, // 300s in nanoseconds
					ConnMaxIdleTime: 120000000000, // 120s in nanoseconds
				},
			},

			AutoMigrate: true,
		}

		if dbCfg.Type != "mysql" {
			t.Errorf("Expected database type 'mysql', got %s", dbCfg.Type)
		}

		if dbCfg.AuthDB.DSN != "root:root@tcp(localhost:3308)/auth_db?charset=utf8mb4&parseTime=True&loc=Local" {
			t.Errorf("Expected DSN 'root:root@tcp(localhost:3308)/auth_db?charset=utf8mb4&parseTime=True&loc=Local', got %s", dbCfg.AuthDB.DSN)
		}

		if !dbCfg.AutoMigrate {
			t.Error("Expected AutoMigrate to be true")
		}
	})

	t.Run("RedisConfig", func(t *testing.T) {
		// 测试Redis配置结构
		redisCfg := RedisConfig{
			Mode: "standalone",
			Standalone: RedisInstanceConfig{
				Addr:        "localhost:6379",
				Password:    "",
				DB:          0,
				DialTimeout: 5000000000, // 5s in nanoseconds
			},
			Pool: RedisPoolConfig{
				PoolSize:     10,
				MinIdleConns: 5,
				MaxRetries:   3,
			},
			Usage: RedisUsageConfig{
				Session: RedisInstanceConfig{
					Addr:        "localhost:6379",
					Password:    "",
					DB:          1,
					DialTimeout: 5000000000, // 5s in nanoseconds
				},
				Cache: RedisInstanceConfig{
					Addr:        "localhost:6379",
					Password:    "",
					DB:          2,
					DialTimeout: 5000000000, // 5s in nanoseconds
				},
			},
		}

		if redisCfg.Mode != "standalone" {
			t.Errorf("Expected Redis mode 'standalone', got %s", redisCfg.Mode)
		}

		if redisCfg.Standalone.Addr != "localhost:6379" {
			t.Errorf("Expected Redis address 'localhost:6379', got %s", redisCfg.Standalone.Addr)
		}

		if redisCfg.Pool.PoolSize != 10 {
			t.Errorf("Expected Redis pool size 10, got %d", redisCfg.Pool.PoolSize)
		}
	})

	// t.Run("MessageQueueConfig", func(t *testing.T) {
	// 	// 测试消息队列配置结构
	// 	mqCfg := MessageQueueConfig{
	// 		Type: "redis",
	// 		Redis: RedisInstanceConfig{
	// 			Addr:        "localhost:6379",
	// 			Password:    "",
	// 			DB:          3,
	// 			DialTimeout: 5000000000, // 5s in nanoseconds
	// 		},
	// 		Kafka: KafkaConfig{
	// 			Brokers:  []string{"localhost:9092"},
	// 			Username: "",
	// 			Password: "",
	// 		},
	// 	}

	// 	if mqCfg.Type != "redis" {
	// 		t.Errorf("Expected message queue type 'redis', got %s", mqCfg.Type)
	// 	}

	// 	if mqCfg.Redis.Addr != "localhost:6379" {
	// 		t.Errorf("Expected Redis address 'localhost:6379', got %s", mqCfg.Redis.Addr)
	// 	}

	// 	if len(mqCfg.Kafka.Brokers) != 1 {
	// 		t.Errorf("Expected 1 Kafka broker, got %d", len(mqCfg.Kafka.Brokers))
	// 	}
	// })
}

func TestInit(t *testing.T) {
	t.Run("InitConfig", func(t *testing.T) {
		// 测试初始化
		cfg, err := Init()
		if err != nil {
			t.Errorf("Expected no error during config init, got: %v", err)
		}
		if cfg == nil {
			t.Error("Expected non-nil config, got nil")
		}

		// 验证项目配置
		if cfg.Project.Name != "learn_mx" {
			t.Errorf("Expected project name 'learn_mx', got '%s'", cfg.Project.Name)
		}

		if cfg.Project.Env != "test" {
			t.Errorf("Expected project env 'test', got '%s'", cfg.Project.Env)
		}

		if cfg.Project.Version != "1.0.0" {
			t.Errorf("Expected project version '1.0.0', got '%s'", cfg.Project.Version)
		}

		// 验证服务器配置
		if len(cfg.Server.GateServers) != 1 {
			t.Errorf("Expected 1 gate servers, got %d", len(cfg.Server.GateServers))
		} else {
			gateServer1 := cfg.Server.GateServers[0]
			if gateServer1.Name != "gate_server_1" {
				t.Errorf("Expected first gate server name 'gate_server_1', got '%s'", gateServer1.Name)
			}
			if gateServer1.GRPCConfig.Host != "localhost" {
				t.Errorf("Expected first gate server host 'localhost', got '%s'", gateServer1.GRPCConfig.Host)
			}
			if gateServer1.GRPCConfig.Port != 50050 {
				t.Errorf("Expected first gate server port 50050, got %d", gateServer1.GRPCConfig.Port)
			}

			if gateServer1.HttpConfig.Host != "localhost" {
				t.Errorf("Expected first gate server HTTP host 'localhost', got '%s'", gateServer1.HttpConfig.Host)
			}
			if gateServer1.HttpConfig.Port != 8080 {
				t.Errorf("Expected first gate server HTTP port 8080, got %d", gateServer1.HttpConfig.Port)
			}

			if gateServer1.ServiceConfig.ClusterId != "cluster_1" {
				t.Errorf("Expected first gate server cluster ID 'cluster_1', got '%s'", gateServer1.ServiceConfig.ClusterId)
			}
			if gateServer1.ServiceConfig.LoadBalance.HealthCheck != true {
				t.Error("Expected first gate server load balance health check to be true")
			}
		}

		if len(cfg.Server.VerifyServers) != 1 {
			t.Errorf("Expected 1 verify servers, got %d", len(cfg.Server.VerifyServers))
		} else {
			verifyServer1 := cfg.Server.VerifyServers[0]
			if verifyServer1.Name != "verify_server_1" {
				t.Errorf("Expected first verify server name 'verify_server_1', got '%s'", verifyServer1.Name)
			}
			if verifyServer1.GRPCConfig.Host != "localhost" {
				t.Errorf("Expected first verify server gRPC host 'localhost', got '%s'", verifyServer1.GRPCConfig.Host)
			}
			if verifyServer1.GRPCConfig.Port != 50060 {
				t.Errorf("Expected first verify server gRPC port 50060, got %d", verifyServer1.GRPCConfig.Port)
			}
			if verifyServer1.HttpConfig.Host != "localhost" {
				t.Errorf("Expected first verify server HTTP host 'localhost', got '%s'", verifyServer1.HttpConfig.Host)
			}
			if verifyServer1.HttpConfig.Port != 8090 {
				t.Errorf("Expected first verify server HTTP port 8090, got %d", verifyServer1.HttpConfig.Port)
			}

			if verifyServer1.ServiceConfig.ClusterId != "cluster_1" {
				t.Errorf("Expected first verify server cluster ID 'cluster_1', got '%s'", verifyServer1.ServiceConfig.ClusterId)
			}
			if verifyServer1.ServiceConfig.LoadBalance.HealthCheck != true {
				t.Error("Expected first verify server load balance health check to be true")
			}

		}

		// 验证全局配置
		if cfg.Server.GlobalConfig.Env != "test" {
			t.Errorf("Expected global env 'test', got '%s'", cfg.Server.GlobalConfig.Env)
		}

		if cfg.Server.GlobalConfig.DefaultLogLevel != "debug" {
			t.Errorf("Expected default log level 'debug', got '%s'", cfg.Server.GlobalConfig.DefaultLogLevel)
		}

		// 验证数据库配置
		if cfg.Database.Type != "mysql" {
			t.Errorf("Expected database type 'mysql', got '%s'", cfg.Database.Type)
		}

		expectedDSN := "root:root@tcp(localhost:3308)/auth_db?charset=utf8mb4&parseTime=True&loc=Local"
		if cfg.Database.AuthDB.DSN != expectedDSN {
			t.Errorf("Expected database DSN '%s', got '%s'", expectedDSN, cfg.Database.AuthDB.DSN)
		}

		if cfg.Database.AuthDB.CharSet != "utf8mb4" {
			t.Errorf("Expected database charset 'utf8mb4', got '%s'", cfg.Database.AuthDB.CharSet)
		}

		if !cfg.Database.AutoMigrate {
			t.Error("Expected AutoMigrate to be true")
		}

		authDB := cfg.Database.AuthDB
		if authDB.Pool.MaxOpenConns != 50 {
			t.Errorf("Expected database max open connections 50, got %d", authDB.Pool.MaxOpenConns)
		}
		if authDB.Pool.MaxIdleConns != 10 {
			t.Errorf("Expected database max idle connections 5, got %d", authDB.Pool.MaxIdleConns)
		}

		// 验证Redis配置
		if cfg.Redis.Mode != "standalone" {
			t.Errorf("Expected Redis mode 'standalone', got '%s'", cfg.Redis.Mode)
		}

		if cfg.Redis.Standalone.Addr != "localhost:6379" {
			t.Errorf("Expected Redis address 'localhost:6379', got '%s'", cfg.Redis.Standalone.Addr)
		}

		if cfg.Redis.Pool.PoolSize != 10 {
			t.Errorf("Expected Redis pool size 10, got %d", cfg.Redis.Pool.PoolSize)
		}

		if cfg.Redis.Pool.MinIdleConns != 5 {
			t.Errorf("Expected Redis min idle connections 5, got %d", cfg.Redis.Pool.MinIdleConns)
		}

		if cfg.Redis.Usage.Session.Addr != "localhost:6379" {
			t.Errorf("Expected Redis session address 'localhost:6379', got '%s'", cfg.Redis.Usage.Session.Addr)
		}

		if cfg.Redis.Usage.Session.DB != 1 {
			t.Errorf("Expected Redis session DB 1, got %d", cfg.Redis.Usage.Session.DB)
		}

		if cfg.Redis.Usage.Cache.Addr != "localhost:6379" {
			t.Errorf("Expected Redis cache address 'localhost:6379', got '%s'", cfg.Redis.Usage.Cache.Addr)
		}

		if cfg.Redis.Usage.Cache.DB != 2 {
			t.Errorf("Expected Redis cache DB 2, got %d", cfg.Redis.Usage.Cache.DB)
		}

	})
}
