package server

import (
	"context"
	"fmt"
	"github.com/donetkit/gin-contrib-log/glog"
	"github.com/donetkit/gin-contrib/discovery"
	"github.com/donetkit/gin-contrib/tracer"
	"github.com/donetkit/gin-contrib/utils/console_colors"
	"github.com/donetkit/gin-contrib/utils/files"
	"github.com/donetkit/gin-contrib/utils/host"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

type config struct {
	Ctx             context.Context
	Tracer          *tracer.Server
	Logger          glog.ILogger
	ServiceName     string
	Host            string
	Port            int
	handler         http.Handler
	httpServer      http.Server
	clientDiscovery discovery.Discovery
	readTimeout     time.Duration
	writerTimeout   time.Duration
	maxHeaderBytes  int
	Version         string
	protocol        string
	pId             int
	environment     string
}

type Server struct {
	Options *config
}

func New(opts ...Option) (*Server, error) {
	var cfg = &config{
		Ctx:            context.Background(),
		ServiceName:    "demo",
		Host:           host.GetOutBoundIp(),
		Port:           80,
		Logger:         glog.New(),
		Version:        "V0.1",
		protocol:       "HTTP API",
		pId:            os.Getpid(),
		environment:    EnvName,
		writerTimeout:  time.Second * 120,
		readTimeout:    time.Second * 120,
		maxHeaderBytes: 1 << 20,
	}
	for _, opt := range opts {
		opt(cfg)
	}
	server := &Server{
		Options: cfg,
	}
	return server, nil
}

func (s *Server) IsDevelopment() bool {
	return s.Options.environment == Dev
}

func (s *Server) IsTest() bool {
	return s.Options.environment == Test
}

func (s *Server) IsProduction() bool {
	return s.Options.environment == Prod
}

func (s *Server) Stop() error {
	s.Options.Logger.Info("Server is stopping")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5) // 平滑关闭,等待5秒钟处理
	defer cancel()
	if err := s.deregister(); err != nil {
		return errors.Wrap(err, "deregister http server error")
	}
	if err := s.Options.httpServer.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "shutdown http server error")
	}
	s.Options.Logger.Info("Server is stopped.")
	return nil
}

func (s *Server) registerDiscovery() *Server {
	if s.Options.clientDiscovery == nil {
		return nil
	}
	err := s.Options.clientDiscovery.Register()
	if err != nil {
		s.Options.Logger.Error(err.Error())
	}
	return s
}

func (s *Server) deregister() error {
	if s.Options.clientDiscovery == nil {
		return nil
	}
	err := s.Options.clientDiscovery.Deregister()
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) awaitSignal() {
	c := make(chan os.Signal, 1)
	signal.Reset(syscall.SIGTERM, syscall.SIGINT)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	select {
	case c := <-c:
		s.Options.Logger.Info("receive a signal, " + "signal: " + c.String())
		if err := s.Stop(); err != nil {
			s.Options.Logger.Error("stop http server error %s", err.Error())
		}
		if s.Options.Tracer != nil {
			s.Options.Tracer.Stop(s.Options.Ctx)
		}
		os.Exit(0)
	}
}

func (s *Server) printLog() {
	s.Options.Logger.Info("======================================================================")
	s.Options.Logger.Info(console_colors.Green("Starting server..."))
	s.Options.Logger.Info(console_colors.Green(fmt.Sprintf("Welcome to %s, starting application ...", s.Options.ServiceName)))
	s.Options.Logger.Info(fmt.Sprintf("framework version        :  %s", console_colors.Blue(s.Options.Version)))
	s.Options.Logger.Info(fmt.Sprintf("server & protocol        :  %s", console_colors.Green(s.Options.protocol)))
	s.Options.Logger.Info(fmt.Sprintf("machine host ip          :  %s", console_colors.Blue(s.Options.Host)))
	s.Options.Logger.Info(fmt.Sprintf("listening on port        :  %s", console_colors.Blue(fmt.Sprintf("%d", s.Options.Port))))
	s.Options.Logger.Info(fmt.Sprintf("application running pid  :  %s", console_colors.Blue(strconv.Itoa(s.Options.pId))))
	s.Options.Logger.Info(fmt.Sprintf("application name         :  %s", console_colors.Blue(s.Options.ServiceName)))
	s.Options.Logger.Info(fmt.Sprintf("application exec path    :  %s", console_colors.Yellow(files.GetCurrentDirectory())))
	s.Options.Logger.Info(fmt.Sprintf("application environment  :  %s", console_colors.Yellow(console_colors.Blue(s.Options.environment))))
	s.Options.Logger.Info("running in %s mode , change (Dev,Test,Prod) mode by Environment .", console_colors.Red(s.Options.environment))
	s.Options.Logger.Info(console_colors.Green("Server is Started."))
	s.Options.Logger.Info("======================================================================")
}

func (s *Server) AddDiscovery(client discovery.Discovery) *Server {
	if client == nil {
		return nil
	}
	s.Options.clientDiscovery = client
	return s
}

func (s *Server) AddTrace(tracer *tracer.Server) *Server {
	s.Options.Tracer = tracer
	return s
}

func (s *Server) AddHandler(handler http.Handler) *Server {
	s.Options.handler = handler
	return s
}

func (s *Server) Run() {
	addr := fmt.Sprintf("%s:%d", s.Options.Host, s.Options.Port)
	s.Options.httpServer = http.Server{
		Addr:           addr,
		Handler:        s.Options.handler,
		ReadTimeout:    s.Options.readTimeout,
		WriteTimeout:   s.Options.writerTimeout,
		MaxHeaderBytes: s.Options.maxHeaderBytes,
	}
	go func() {
		if err := s.Options.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return
		}
	}()
	s.registerDiscovery()
	s.printLog()
	s.awaitSignal()
}
