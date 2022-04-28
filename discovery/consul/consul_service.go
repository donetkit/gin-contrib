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
	Options               *discovery.Config
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
		Options: cfg,
	}
	return consulClient, nil
}

func (s *Client) Register() error {
	if s.Options.CheckHTTP == "" {
		return s.serviceRegisterTCP()
	}
	return s.serviceRegisterHttp()
}

func (s *Client) Deregister() error {
	if s.Options.CheckHTTP == "" {
		return s.serviceDeregisterTCP()
	}
	return s.serviceDeregisterHttp()
}

func (s *Client) serviceRegisterTCP() error {

	consulCli, err := consulApi.NewClient(&consulApi.Config{Address: fmt.Sprintf("%s:%d", s.Options.ServiceRegisterAddr, s.Options.ServiceRegisterPort)})
	if err != nil {
		return fmt.Errorf("create consul client error")
	}
	s.client = consulCli

	addr := fmt.Sprintf("%s:%d", s.Options.ServiceCheckAddr, s.Options.ServiceCheckPort)
	check := &consulApi.AgentServiceCheck{
		Interval:                       fmt.Sprintf("%ds", s.Options.IntervalTime),
		DeregisterCriticalServiceAfter: fmt.Sprintf("%ds", s.Options.DeregisterTime),
		TCP:                            addr,
	}
	svcReg := &consulApi.AgentServiceRegistration{
		ID:                s.Options.Id,
		Name:              s.Options.ServiceName,
		Tags:              s.Options.Tags,
		Port:              s.Options.ServiceCheckPort,
		Address:           s.Options.ServiceCheckAddr,
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
	err := s.client.Agent().ServiceDeregister(s.Options.Id)
	if err != nil {
		return errors.Wrapf(err, "deregister service error[key=%s]", s.Options.Id)
	}
	return nil
}

func (s *Client) serviceRegisterHttp() error {
	registryClient, err := NewConsulServiceRegistryAddress(fmt.Sprintf("%s:%d", s.Options.ServiceRegisterAddr, s.Options.ServiceRegisterPort), "")
	if err != nil {
		return err
	}
	s.consulServiceRegistry = registryClient
	serviceInstance := DefaultServiceInstance{
		InstanceId:     s.Options.Id,
		ServiceName:    s.Options.ServiceName,
		Host:           s.Options.ServiceCheckAddr,
		Port:           s.Options.ServiceCheckPort,
		Metadata:       s.Options.Tags,
		Timeout:        s.Options.TimeOut,
		Interval:       s.Options.IntervalTime,
		DeregisterTime: s.Options.DeregisterTime,
		CheckHTTP:      s.Options.CheckHTTP,
	}
	serviceInstanceInfo, err := NewDefaultServiceInstance(&serviceInstance)
	if err != nil {
		return err
	}
	if s.consulServiceRegistry.Register(serviceInstanceInfo) {
		return errors.New("register fail")
	}
	fmt.Println(s.Options.CheckHTTP)
	return nil
}

func (s *Client) serviceDeregisterHttp() error {
	s.consulServiceRegistry.Deregister()
	return nil
}
