package discovery

import (
	"math/rand"
	"strconv"
	"time"
)

type ServiceInstance interface {
	GetInstanceId() string

	GetServiceName() string

	GetHost() string

	GetPort() int

	GetMetadata() []string

	GetCheckHTTP() string

	GetTimeout() int

	GetInterval() int

	GetDeregisterTime() int
}

type DefaultServiceInstance struct {
	InstanceId     string
	ServiceName    string
	Host           string
	Port           int
	Metadata       []string
	CheckHTTP      string
	Timeout        int
	Interval       int
	DeregisterTime int //故障检查失败30s后 consul自动将注册服务删除
}

func NewDefaultServiceInstance(instance *DefaultServiceInstance) (*DefaultServiceInstance, error) {
	if len(instance.InstanceId) == 0 {
		instance.InstanceId = instance.ServiceName + "-" + strconv.FormatInt(time.Now().Unix(), 10) + "-" + strconv.Itoa(rand.Intn(9000)+1000)
	}
	return instance, nil
}

func (serviceInstance DefaultServiceInstance) GetInstanceId() string {
	return serviceInstance.InstanceId
}

func (serviceInstance DefaultServiceInstance) GetServiceName() string {
	return serviceInstance.ServiceName
}

func (serviceInstance DefaultServiceInstance) GetHost() string {
	return serviceInstance.Host
}

func (serviceInstance DefaultServiceInstance) GetPort() int {
	return serviceInstance.Port
}

func (serviceInstance DefaultServiceInstance) GetMetadata() []string {
	return serviceInstance.Metadata
}

func (serviceInstance DefaultServiceInstance) GetCheckHTTP() string {
	return serviceInstance.CheckHTTP
}

func (serviceInstance DefaultServiceInstance) GetTimeout() int {
	return serviceInstance.Timeout
}

func (serviceInstance DefaultServiceInstance) GetInterval() int {
	return serviceInstance.Interval
}

func (serviceInstance DefaultServiceInstance) GetDeregisterTime() int {
	return serviceInstance.DeregisterTime
}

type ServiceRegistry interface {
	Register(serviceInstance ServiceInstance) bool
	Deregister()
}
