package memregistry

import (
	"errors"
	"sync"
	"time"

	status_model "github.com/mxxmstar/learning/pkg/model/status"
)

const (
	CleanupInterval = 30 * time.Second
)

// MemRegistry 内存注册中心
type MemRegistry struct {
	services map[string]*status_model.ServiceInfo // 服务信息存储
	mutex    sync.RWMutex
}

func NewMemRegistry() *MemRegistry {
	reg := &MemRegistry{
		services: make(map[string]*status_model.ServiceInfo),
	}
	// 启动一个 goroutine，用于定时清理过期的服务
	go reg.StartCleanupTask()
	return reg
}

func (r *MemRegistry) StartCleanupTask() {
	ticker := time.NewTicker(CleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		r.CleanupExpiredServices()
	}
}

func (r *MemRegistry) RegisterService(service *status_model.ServiceInfo) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// 检查服务是否已存在
	key := r.getServiceKey(service.ServiceType, service.ServiceId)
	if _, exists := r.services[key]; exists {
		return errors.New("service already exists")
	}

	// service.Status = "active"
	service.LastHeartbeat = status_model.GetCurrentTimestamp()
	r.services[key] = service

	return nil
}

func (r *MemRegistry) GetService(serviceType, serviceId string) (*status_model.ServiceInfo, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	key := r.getServiceKey(serviceType, serviceId)
	service, exists := r.services[key]
	if !exists {
		return nil, errors.New("service not found")
	}

	if service.IsExpired() {
		delete(r.services, key)
		return nil, errors.New("service expired")
	}
	return service, nil
}

func (r *MemRegistry) DiscoverServicesByType(serviceType string) ([]*status_model.ServiceInfo, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var services []*status_model.ServiceInfo
	for key, service := range r.services {
		if service.ServiceType == serviceType {
			if !service.IsExpired() {
				services = append(services, service)
			} else {
				// 标记为过期，启动一个 goroutine ，在写锁中删除
				go func(k string) {
					r.mutex.Lock()
					defer r.mutex.Unlock()

					service, exists := r.services[k]
					if exists && service.IsExpired() {
						delete(r.services, k)
					}
				}(key)
			}
		}
	}

	return services, nil
}

func (r *MemRegistry) GetAllServices() ([]*status_model.ServiceInfo, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var services []*status_model.ServiceInfo
	for key, service := range r.services {
		if !service.IsExpired() {
			services = append(services, service)
		} else {
			// 标记为过期，启动一个 goroutine ，在写锁中删除
			go func(k string) {
				r.mutex.Lock()
				defer r.mutex.Unlock()

				service, exists := r.services[k]
				if exists && service.IsExpired() {
					delete(r.services, k)
				}
			}(key)
		}
	}

	return services, nil
}

func (r *MemRegistry) DeregisterService(serviceType, serviceId string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	key := r.getServiceKey(serviceType, serviceId)
	if _, exists := r.services[key]; !exists {
		return errors.New("service not found")
	}
	delete(r.services, key)
	return nil
}

func (r *MemRegistry) KeepAlive(serviceType, serviceId string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	key := r.getServiceKey(serviceType, serviceId)
	service, exists := r.services[key]
	if !exists {
		return errors.New("service not found")
	}

	service.LastHeartbeat = status_model.GetCurrentTimestamp()
	return nil
}

// getServiceKey 生成服务键
func (r *MemRegistry) getServiceKey(serviceType, serviceId string) string {
	return serviceType + ":" + serviceId
}

func (r *MemRegistry) CleanupExpiredServices() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	for key, service := range r.services {
		if service.IsExpired() {
			delete(r.services, key)
		}
	}
}
