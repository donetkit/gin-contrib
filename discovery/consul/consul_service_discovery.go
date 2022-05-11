package consul

import (
	"fmt"
	"github.com/donetkit/gin-contrib/discovery"
	consulApi "github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	"time"
)

type Client struct {
	client  *consulApi.Client
	options *discovery.Config
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
		TimeOut:             3,
	}
	for _, opt := range opts {
		opt(cfg)
	}
	consulCli, err := consulApi.NewClient(&consulApi.Config{Address: fmt.Sprintf("%s:%d", cfg.ServiceRegisterAddr, cfg.ServiceRegisterPort)})
	if err != nil {
		return nil, errors.Wrap(err, "create consul client error")
	}
	consulClient := &Client{
		options: cfg,
		client:  consulCli,
	}
	return consulClient, nil
}

func (s *Client) Register() error {
	check := &consulApi.AgentServiceCheck{
		Timeout:                        fmt.Sprintf("%ds", s.options.TimeOut),        // 超时时间
		Interval:                       fmt.Sprintf("%ds", s.options.IntervalTime),   // 健康检查间隔
		DeregisterCriticalServiceAfter: fmt.Sprintf("%ds", s.options.DeregisterTime), //check失败后30秒删除本服务，注销时间，相当于过期时间
		HTTP:                           s.options.CheckHTTP,
	}
	if s.options.CheckHTTP == "" {
		check.TCP = fmt.Sprintf("%s:%d", s.options.ServiceCheckAddr, s.options.ServiceCheckPort)
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
	err := s.client.Agent().ServiceRegister(svcReg)
	if err != nil {
		return errors.Wrap(err, "register service error")
	}
	return nil
}

func (s *Client) Deregister() error {
	err := s.client.Agent().ServiceDeregister(s.options.Id)
	if err != nil {
		return errors.Wrapf(err, "deregister service error[key=%s]", s.options.Id)
	}
	return nil
}

func (s *Client) Get(key string) ([]byte, error) {
	kv, _, err := s.client.KV().Get(key, nil)
	if err != nil {
		return nil, err
	}
	if kv == nil {
		return nil, errors.New("not found value")
	}
	return kv.Value, nil
}

func (s *Client) Set(key string, value string) error {
	p := &consulApi.KVPair{Key: key, Value: []byte(value)}
	if _, err := s.client.KV().Put(p, nil); err != nil {
		return err
	}
	return nil
}

func (s *Client) Delete(key string) error {
	if _, err := s.client.KV().Delete(key, nil); err != nil {
		return err
	}
	return nil
}

func (s *Client) List(key string) (map[string][]byte, error) {
	p, _, err := s.client.KV().List(key, nil)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, errors.New("not found value")
	}
	values := make(map[string][]byte, len(p))
	for _, v := range p {
		values[v.Key] = v.Value
	}
	return values, nil

}
