package server

import (
	"context"
	"fmt"
	"github.com/donetkit/gin-contrib/discovery"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type config struct {
	serviceName    string
	host           string
	port           int
	router         *gin.Engine
	httpServer     http.Server
	consulClient   interface{}
	readTimeout    time.Duration
	writerTimeout  time.Duration
	maxHeaderBytes int
}

type Server struct {
	Options *config
}

type InitControllers func(r *gin.Engine)

func New(opts ...Option) (*Server, error) {
	var cfg = &config{
		serviceName: "demo",
		host:        "127.0.0.1",
		port:        80,
	}
	for _, opt := range opts {
		opt(cfg)
	}
	return &Server{Options: cfg}, nil
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.Options.host, s.Options.port)
	s.Options.httpServer = http.Server{
		Addr:           addr,
		Handler:        s.Options.router,
		ReadTimeout:    s.Options.readTimeout,
		WriteTimeout:   s.Options.writerTimeout,
		MaxHeaderBytes: s.Options.maxHeaderBytes,
	}
	go func() {
		if err := s.Options.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return
		}
	}()
	s.register()
	return nil
}

func (s *Server) register() error {
	if s.Options.consulClient == nil {
		return nil
	}
	serverDiscovery, ok := s.Options.consulClient.(discovery.Discovery)
	if ok {
		err := serverDiscovery.Register()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) deregister() error {
	if s.Options.consulClient == nil {
		return nil
	}
	serverDiscovery, ok := s.Options.consulClient.(discovery.Discovery)
	if ok {
		err := serverDiscovery.Deregister()
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5) // 平滑关闭,等待5秒钟处理
	defer cancel()

	if err := s.deregister(); err != nil {
		return errors.Wrap(err, "deregister http server error")
	}

	if err := s.Options.httpServer.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "shutdown http server error")
	}

	return nil
}

func (s *Server) AwaitSignal() {
	c := make(chan os.Signal, 1)
	signal.Reset(syscall.SIGTERM, syscall.SIGINT)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	select {
	case c := <-c:
		fmt.Println("receive a signal", "signal", c.String())
		if err := s.stop(); err != nil {
			fmt.Println("stop http server error", err)
		}
		os.Exit(0)
	}
}
