package consul

import (
	"fmt"
	"github.com/donetkit/gin-contrib/discovery"
	consulApi "github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	"time"
)

type Client struct {
	client                *consulApi.Client
	options               *discovery.Config
	consulServiceRegistry *ServiceConsul
}

func New(opts ...discovery.Option) (*Client, error) {
	cfg := &discovery.Config{
		Id:                  fmt.Sprintf("%d", time.Now().UnixNano()),
		ServiceName:         "127.0.0.1:80",
		ServiceRegisterAddr: "127.0.0.1",
		ServiceRegisterPort: 8500,
		ServiceCheckAddr:    "127.0.0.1",
		ServiceCheckPort:    80,
		Tags:                []string{"v0.0.1"},
		IntervalTime:        5,
		DeregisterTime:      15,
		TimeOut:             5,
	}
	for _, opt := range opts {
		opt(cfg)
	}
	consulClient := &Client{
		options: cfg,
	}
	return consulClient, nil
}

func (s *Client) Register() error {
	if s.options.CheckHTTP == "" {
		return s.serviceRegisterTCP()
	}
	return s.serviceRegisterHttp()
}

func (s *Client) Deregister() error {
	if s.options.CheckHTTP == "" {
		return s.serviceDeregisterTCP()
	}
	return s.serviceDeregisterHttp()
}

func (s *Client) serviceRegisterTCP() error {

	consulCli, err := consulApi.NewClient(&consulApi.Config{Address: fmt.Sprintf("%s:%d", s.options.ServiceRegisterAddr, s.options.ServiceRegisterPort)})
	if err != nil {
		return fmt.Errorf("create consul client error")
	}
	s.client = consulCli

	addr := fmt.Sprintf("%s:%d", s.options.ServiceCheckAddr, s.options.ServiceCheckPort)
	check := &consulApi.AgentServiceCheck{
		Interval:                       fmt.Sprintf("%ds", s.options.IntervalTime),
		DeregisterCriticalServiceAfter: fmt.Sprintf("%ds", s.options.DeregisterTime),
		TCP:                            addr,
	}
	svcReg := &consulApi.AgentServiceRegistration{
		ID:                s.options.Id,
		Name:              s.options.ServiceName,
		Tags:              s.options.Tags,
		Port:              s.options.ServiceCheckPort,
		Address:           s.options.ServiceCheckAddr,
		EnableTagOverride: true,
		Check:             check,
		Checks:            nil,
	}
	err = s.client.Agent().ServiceRegister(svcReg)
	if err != nil {
		return errors.Wrap(err, "register service error")
	}
	return nil
}

func (s *Client) serviceDeregisterTCP() error {
	err := s.client.Agent().ServiceDeregister(s.options.Id)
	if err != nil {
		return errors.Wrapf(err, "deregister service error[key=%s]", s.options.Id)
	}
	return nil
}

func (s *Client) serviceRegisterHttp() error {
	registryClient, err := NewConsulServiceRegistryAddress(fmt.Sprintf("%s:%d", s.options.ServiceRegisterAddr, s.options.ServiceRegisterPort), "")
	if err != nil {
		return err
	}
	s.consulServiceRegistry = registryClient
	serviceInstance := DefaultServiceInstance{
		InstanceId:     s.options.Id,
		ServiceName:    s.options.ServiceName,
		Host:           s.options.ServiceCheckAddr,
		Port:           s.options.ServiceCheckPort,
		Metadata:       s.options.Tags,
		Timeout:        s.options.TimeOut,
		Interval:       s.options.IntervalTime,
		DeregisterTime: s.options.DeregisterTime,
		CheckHTTP:      s.options.CheckHTTP,
	}
	serviceInstanceInfo, err := NewDefaultServiceInstance(&serviceInstance)
	if err != nil {
		return err
	}
	if s.consulServiceRegistry.Register(serviceInstanceInfo) {
		return errors.New("register fail")
	}
	fmt.Println(s.options.CheckHTTP)
	return nil
}

func (s *Client) serviceDeregisterHttp() error {
	s.consulServiceRegistry.Deregister()
	return nil
}
