package model

type Registry interface {
	RegisterService(service *ServiceInfo) error
	GetService(serviceType, serviceId string) (*ServiceInfo, error)
	DiscoverServicesByType(serviceType string) ([]*ServiceInfo, error)
	GetAllServices() ([]*ServiceInfo, error)
	DeregisterService(serviceType, serviceId string) error
}
