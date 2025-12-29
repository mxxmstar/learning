package registry

import (
	"errors"
	"sync"
)

type ServiceRegistry struct {
	services map[string]string // 服务名称 -> 服务地址
	mutex    sync.RWMutex
}

func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{
		services: make(map[string]string),
	}
}

func (r *ServiceRegistry) RegisterService(serviceName, serviceAddress string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.services[serviceName] = serviceAddress
	return nil
}

func (r *ServiceRegistry) GetServiceAddress(serviceName string) (string, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	address, ok := r.services[serviceName]
	if !ok {
		return "", errors.New("service not found: " + serviceName)
	}
	return address, nil
}

func (r *ServiceRegistry) RemoveService(serviceName string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	delete(r.services, serviceName)
	return nil
}
