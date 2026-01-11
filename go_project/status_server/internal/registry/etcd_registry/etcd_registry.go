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
	cache  sync.Map // 本地缓存
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
	serviceKey := ServicePrefix + service.ServiceType + "/" + service.ServiceId
	_, err = r.client.Put(ctx, serviceKey, string(serviceBytes), clientv3.WithLease(leaseResp.ID))
	if err != nil {
		return fmt.Errorf("failed to put service info: %v", err)
	}

	// 启动心跳续租
	go r.keepAlive(service.ServiceId, service.ServiceType, leaseResp.ID, service.TTLSeconds)

	// 更新本地缓存
	r.cache.Store(serviceKey, service)

	// TODO: 日志记录
	fmt.Printf("Registered service: %s_%s, lease Id: %d\n", service.ServiceType, service.ServiceId, leaseResp.ID)
	return nil
}

func (r *EtcdRegistry) keepAlive(serviceId string, ServiceType string, leaseId clientv3.LeaseID, ttl int64) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 创建续租
	keepAliveChan, err := r.client.KeepAlive(ctx, leaseId)
	if err != nil {
		// TODO: 日志记录
		fmt.Printf("Failed to keep alive for service, type is %s and id is %s: %v\n", ServiceType, serviceId, err)
		return
	}

	// 定期更新心跳时间
	ticker := time.NewTicker(time.Duration(ttl/3) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case _, ok := <-keepAliveChan:
			if !ok {
				fmt.Printf("Keep alive channel closed for service, type is %s and id is %s\n", ServiceType, serviceId)
				return
			}
		case <-ticker.C:
			// 搜索对应的服务
			serviceKey := ServicePrefix + ServiceType + "/" + serviceId
			service, ok := r.cache.Load(serviceKey)
			if !ok {
				// TODO: 日志记录
				fmt.Printf("Service: %s_%s not found in cache\n", ServiceType, serviceId)
				continue
			}

			// 更新心跳时间
			serviceInfo := service.(*model.ServiceInfo)
			serviceInfo.LastHeartbeat = model.GetCurrentTimestamp()

			serviceBytes, err := json.Marshal(serviceInfo)
			if err != nil {
				// TODO: 日志记录
				fmt.Printf("Failed to marshal service info for service at headtbeat %s_%s: %v\n", ServiceType, serviceId, err)
				continue
			}

			ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			_, err = r.client.Put(ctxWithTimeout, serviceKey, string(serviceBytes))
			cancel()
			if err != nil {
				// TODO: 日志记录
				fmt.Printf("Failed to update heartbeat for service %s_%s: %v\n", ServiceType, serviceId, err)
			}
		}
	}
}

func (r *EtcdRegistry) GetService(serviceType, serviceId string) (*model.ServiceInfo, error) {
	// 先尝试从本地缓存中获取
	cacheKey := ServicePrefix + serviceType + "/" + serviceId
	if cached, ok := r.cache.Load(cacheKey); ok {
		serviceInfo := cached.(*model.ServiceInfo)
		if (!serviceInfo.IsExpired()) && (serviceInfo.Status == "online" || serviceInfo.Status == "active") {
			// 缓存未过期且服务在线
			return serviceInfo, nil
		}
		// 缓存已过期或服务离线，从 Etcd 中重新获取
		r.cache.Delete(cacheKey)
	}

	// 从 Etcd 中获取
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	serviceKey := ServicePrefix + serviceType + "/" + serviceId
	getResp, err := r.client.Get(ctx, serviceKey)
	if err != nil {
		// TODO: 日志记录
		return nil, fmt.Errorf("failed to get service info: %v", err)
	}

	if len(getResp.Kvs) == 0 {
		// 服务不存在
		return nil, fmt.Errorf("service %s_%s not found", serviceType, serviceId)
	}

	var serviceInfo model.ServiceInfo
	err = json.Unmarshal(getResp.Kvs[0].Value, &serviceInfo)
	if err != nil {
		// TODO: 日志记录
		return nil, fmt.Errorf("failed to unmarshal service info: %v", err)
	}

	// 检查服务是否已过期
	if serviceInfo.IsExpired() {
		// 从 Etcd 中删除过期服务
		r.client.Delete(context.Background(), serviceKey)
		return nil, fmt.Errorf("service %s_%s has expired", serviceType, serviceId)
	}

	// 更新本地缓存
	r.cache.Store(serviceKey, &serviceInfo)
	return &serviceInfo, nil
}

func (r *EtcdRegistry) DiscoverServicesByType(serviceType string) ([]*model.ServiceInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	prefix := ServicePrefix + serviceType + "/"
	getResp, err := r.client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		// TODO: 日志记录
		return nil, fmt.Errorf("failed to get services by type: %v", err)
	}

	var services []*model.ServiceInfo
	for _, kv := range getResp.Kvs {
		var serviceInfo model.ServiceInfo
		err := json.Unmarshal(kv.Value, &serviceInfo)
		if err != nil {
			// TODO: 日志记录
			fmt.Printf("Failed to unmarshal service info: %v\n", err)
			continue
		}

		// 检查服务是否已过期
		if serviceInfo.IsExpired() {
			// 从 Etcd 中删除过期服务
			r.client.Delete(context.Background(), string(kv.Key))
			continue
		}

		// 更新本地缓存
		cacheKey := string(kv.Key)
		r.cache.Store(cacheKey, &serviceInfo)
		services = append(services, &serviceInfo)
	}

	return services, nil
}

func (r *EtcdRegistry) GetAllServices() ([]*model.ServiceInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	prefix := ServicePrefix
	getResp, err := r.client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		// TODO: 日志记录
		return nil, fmt.Errorf("failed to get services by type: %v", err)
	}

	var services []*model.ServiceInfo
	for _, kv := range getResp.Kvs {
		var serviceInfo model.ServiceInfo
		err := json.Unmarshal(kv.Value, &serviceInfo)
		if err != nil {
			// TODO: 日志记录
			fmt.Printf("Failed to unmarshal service info: %v\n", err)
			continue
		}

		// 检查服务是否已过期
		if serviceInfo.IsExpired() {
			// 从 Etcd 中删除过期服务
			r.client.Delete(context.Background(), string(kv.Key))
			continue
		}

		// 更新本地缓存
		cacheKey := string(kv.Key)
		r.cache.Store(cacheKey, &serviceInfo)
		services = append(services, &serviceInfo)
	}

	return services, nil
}

func (r *EtcdRegistry) DeregisterService(serviceType, serviceId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	serviceKey := ServicePrefix + serviceType + "/" + serviceId
	_, err := r.client.Delete(ctx, serviceKey)
	if err != nil {
		// TODO: 日志记录
		return fmt.Errorf("failed to deregister service: %v", err)
	}

	// 删除本地缓存
	r.cache.Delete(serviceKey)
	// TODO: 日志记录
	return nil
}

// Close 关闭连接
func (r *EtcdRegistry) Close() error {
	return r.client.Close()
}

// CacheLoad 从缓存加载
func (r *EtcdRegistry) CacheLoad(key string) (interface{}, bool) {
	return r.cache.Load(key)
}

// CacheStore 存储到缓存
func (r *EtcdRegistry) CacheStore(key string, value interface{}) {
	r.cache.Store(key, value)
}

// Client 获取 etcd 客户端
func (r *EtcdRegistry) Client() *clientv3.Client {
	return r.client
}
