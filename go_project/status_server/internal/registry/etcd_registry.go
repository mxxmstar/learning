package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/mxxmstar/learning/status_server/internal/model"
	"go.etcd.io/etcd/clientv3"
)

const (
	ServicePrefix = "/services/"
	HeartbeatKey  = "/heartbeat/"
)

type EtcdRegistry struct {
	client *clientv3.Client
	cache  sync.Map
}

func NewEtcdRegistry(endpoints []string) (*EtcdRegistry, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	reg := &EtcdRegistry{
		client: cli,
	}
	return reg, nil
}

func (r *EtcdRegistry) RegisterService(service *model.ServiceInfo) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 创建租约
	lease := clientv3.NewLease(r.client)
	leaseResp, err := lease.Grant(ctx, int64(service.TTLSeconds))
	if err != nil {
		return fmt.Errorf("failed to grant lease: %v", err)
	}

	// 序列化服务信息
	serviceBytes, err := json.Marshal(service)
	if err != nil {
		return fmt.Errorf("failed to marshal service info: %v", err)
	}

	// 存储服务信息
	serviceKey := ServicePrefix + service.ServiceName
	_, err = r.client.Put(ctx, serviceKey, string(serviceBytes), clientv3.WithLease(leaseResp.ID))
	if err != nil {
		return fmt.Errorf("failed to put service info: %v", err)
	}

	// 启动心跳续租
	go r.keepAlive(service.ServiceName, leaseResp.ID, int64(service.TTLSeconds))

	// 更新本地缓存
	r.cache.Store(serviceKey, service)

	// TODO: 日志记录
	fmt.Printf("Registered service: %s, lease Id: %d\n", service.ServiceName, leaseResp.ID)
	return nil
}

func (r *EtcdRegistry) keepAlive(serviceName string, leaseId clientv3.LeaseID, ttl int64) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 创建续租
	keepAliveChan, err := r.client.KeepAlive(ctx, leaseId)
	if err != nil {
		fmt.Printf("Failed to keep alive for service %s: %v\n", serviceName, err)
		return
	}

	// 定期更新心跳时间
	ticker := time.NewTicker(time.Duration(ttl/3) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case _, ok := <-keepAliveChan:
			if !ok {
				fmt.Printf("Keep alive channel closed for service %s\n", serviceName)
				return
			}
		}
	}
}
