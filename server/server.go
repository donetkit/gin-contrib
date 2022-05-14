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
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

type config struct {
	ctx             context.Context
	tracer          *tracer.Server
	logger          glog.ILogger
	serviceName     string
	host            string
	port            int
	router          *gin.Engine
	httpServer      http.Server
	clientDiscovery discovery.Discovery
	readTimeout     time.Duration
	writerTimeout   time.Duration
	maxHeaderBytes  int
	version         string
	protocol        string
	pId             int
	environment     string
}

type Server struct {
	options *config
}

func New(opts ...Option) (*Server, error) {
	var cfg = &config{
		ctx:            context.Background(),
		serviceName:    "demo",
		host:           host.GetOutBoundIp(),
		port:           80,
		logger:         glog.New(),
		version:        "V0.1",
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
		options: cfg,
	}
	if cfg.router != nil {
		switch server.options.environment {
		case Dev:
			gin.SetMode(gin.DebugMode)
		case Test:
			gin.SetMode(gin.TestMode)
		case Prod:
			gin.SetMode(gin.ReleaseMode)
		}
	}
	return server, nil
}

func (s *Server) IsDevelopment() bool {
	return s.options.environment == Dev
}

func (s *Server) IsTest() bool {
	return s.options.environment == Test
}

func (s *Server) IsProduction() bool {
	return s.options.environment == Prod
}

func (s *Server) Stop() error {
	s.options.logger.Info("Server is stopping")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5) // 平滑关闭,等待5秒钟处理
	defer cancel()
	if err := s.deregister(); err != nil {
		return errors.Wrap(err, "deregister http server error")
	}
	if err := s.options.httpServer.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "shutdown http server error")
	}
	s.options.logger.Info("Server is stopped.")
	return nil
}

func (s *Server) registerDiscovery() *Server {
	if s.options.clientDiscovery == nil {
		return nil
	}
	err := s.options.clientDiscovery.Register()
	if err != nil {
		s.options.logger.Error(err.Error())
	}
	return s
}

func (s *Server) deregister() error {
	if s.options.clientDiscovery == nil {
		return nil
	}
	err := s.options.clientDiscovery.Deregister()
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
		s.options.logger.Info("receive a signal, " + "signal: " + c.String())
		if err := s.Stop(); err != nil {
			s.options.logger.Error("stop http server error %s", err.Error())
		}
		if s.options.tracer != nil {
			s.options.tracer.Stop(s.options.ctx)
		}
		os.Exit(0)
	}
}

func (s *Server) printLog() {
	s.options.logger.Info("======================================================================")
	s.options.logger.Info(console_colors.Green("Starting server..."))
	s.options.logger.Info(console_colors.Green(fmt.Sprintf("Welcome to %s, starting application ...", s.options.serviceName)))
	s.options.logger.Info(fmt.Sprintf("framework version        :  %s", console_colors.Blue(s.options.version)))
	s.options.logger.Info(fmt.Sprintf("server & protocol        :  %s", console_colors.Green(s.options.protocol)))
	s.options.logger.Info(fmt.Sprintf("machine host ip          :  %s", console_colors.Blue(s.options.host)))
	s.options.logger.Info(fmt.Sprintf("listening on port        :  %s", console_colors.Blue(fmt.Sprintf("%d", s.options.port))))
	s.options.logger.Info(fmt.Sprintf("application running pid  :  %s", console_colors.Blue(strconv.Itoa(s.options.pId))))
	s.options.logger.Info(fmt.Sprintf("application name         :  %s", console_colors.Blue(s.options.serviceName)))
	s.options.logger.Info(fmt.Sprintf("application exec path    :  %s", console_colors.Yellow(files.GetCurrentDirectory())))
	s.options.logger.Info(fmt.Sprintf("application environment  :  %s", console_colors.Yellow(console_colors.Blue(s.options.environment))))
	s.options.logger.Info("running in %s mode , change (Dev,Test,Prod) mode by Environment .", console_colors.Red(s.options.environment))
	s.options.logger.Info(console_colors.Green("Server is Started."))
	s.options.logger.Info("======================================================================")
}

func (s *Server) AddDiscovery(client discovery.Discovery) *Server {
	if client == nil {
		return nil
	}
	s.options.clientDiscovery = client
	return s
}

func (s *Server) AddTrace(tracer *tracer.Server) *Server {
	s.options.tracer = tracer
	return s
}

func (s *Server) AddRouter(router *gin.Engine) *Server {
	s.options.router = router
	return s
}

func (s *Server) Run() {
	addr := fmt.Sprintf("%s:%d", s.options.host, s.options.port)
	s.options.httpServer = http.Server{
		Addr:           addr,
		Handler:        s.options.router,
		ReadTimeout:    s.options.readTimeout,
		WriteTimeout:   s.options.writerTimeout,
		MaxHeaderBytes: s.options.maxHeaderBytes,
	}
	go func() {
		if err := s.options.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return
		}
	}()
	s.registerDiscovery()
	s.printLog()
	s.awaitSignal()
}
