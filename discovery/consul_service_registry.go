package discovery

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/consul/api"
)

type ConsulServiceRegistry struct {
	serviceInstances     map[string]map[string]ServiceInstance
	client               api.Client
	localServiceInstance ServiceInstance
}

func (c *ConsulServiceRegistry) Register(serviceInstance ServiceInstance) bool {
	// 创建注册到consul的服务到
	registration := new(api.AgentServiceRegistration)
	registration.ID = serviceInstance.GetInstanceId()
	registration.Name = serviceInstance.GetServiceName()
	registration.Port = serviceInstance.GetPort()
	registration.Tags = serviceInstance.GetMetadata()
	registration.Address = serviceInstance.GetHost()
	// 增加consul健康检查回调函数
	check := new(api.AgentServiceCheck)
	check.HTTP = serviceInstance.GetCheckHTTP()
	check.Timeout = fmt.Sprintf("%d%s", serviceInstance.GetTimeout(), "s")
	check.Interval = fmt.Sprintf("%d%s", serviceInstance.GetInterval(), "s")
	check.DeregisterCriticalServiceAfter = fmt.Sprintf("%d%s", serviceInstance.GetDeregisterTime(), "s") // 故障检查失败30s后 consul自动将注册服务删除
	registration.Check = check
	// 注册服务到consul
	err := c.client.Agent().ServiceRegister(registration)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if c.serviceInstances == nil {
		c.serviceInstances = map[string]map[string]ServiceInstance{}
	}
	services := c.serviceInstances[serviceInstance.GetServiceName()]
	if services == nil {
		services = map[string]ServiceInstance{}
	}
	services[serviceInstance.GetInstanceId()] = serviceInstance
	c.serviceInstances[serviceInstance.GetServiceName()] = services
	c.localServiceInstance = serviceInstance
	return true
}

// deregister a service
func (c *ConsulServiceRegistry) Deregister() {
	if c.serviceInstances == nil {
		return
	}
	services := c.serviceInstances[c.localServiceInstance.GetServiceName()]
	if services == nil {
		return
	}
	delete(services, c.localServiceInstance.GetInstanceId())
	if len(services) == 0 {
		delete(c.serviceInstances, c.localServiceInstance.GetServiceName())
	}
	_ = c.client.Agent().ServiceDeregister(c.localServiceInstance.GetInstanceId())
	c.localServiceInstance = nil
}

// new a consulServiceRegistry instance token is optional
func NewConsulServiceRegistry(host string, port int, token string) (*ConsulServiceRegistry, error) {
	if len(host) < 3 {
		return nil, errors.New("check host")
	}
	if port <= 0 || port > 65535 {
		return nil, errors.New("check port, port should between 1 and 65535")
	}
	config := api.DefaultConfig()
	config.Address = host + ":" + strconv.Itoa(port)
	config.Token = token
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &ConsulServiceRegistry{client: *client}, nil
}

// new a consulServiceRegistry instance  token is optional
func NewConsulServiceRegistryAddress(address string, token string) (*ConsulServiceRegistry, error) {
	config := api.DefaultConfig()
	config.Address = address
	config.Token = token
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &ConsulServiceRegistry{client: *client}, nil
}
