package server

import (
	"context"
	"fmt"
	"github.com/donetkit/gin-contrib/discovery"
	"github.com/donetkit/gin-contrib/utils/console_colors"
	"github.com/donetkit/gin-contrib/utils/files"
	"github.com/donetkit/gin-contrib/utils/glog"
	"github.com/donetkit/gin-contrib/utils/host"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

type config struct {
	logger          glog.ILogger
	serviceName     string
	host            string
	port            int
	router          *gin.Engine
	httpServer      http.Server
	clientDiscovery interface{}
	readTimeout     time.Duration
	writerTimeout   time.Duration
	maxHeaderBytes  int
}

type Server struct {
	options     *config
	pId         int
	environment string
	version     string
	server      string
}

func New(opts ...Option) (*Server, error) {
	var cfg = &config{
		serviceName: "demo",
		host:        host.GetOutBoundIp(),
		port:        80,
	}
	for _, opt := range opts {
		opt(cfg)
	}
	server := &Server{
		options:     cfg,
		pId:         os.Getpid(),
		version:     "V0.1",
		environment: EnvName,
		server:      "HTTP API",
	}
	if cfg.router != nil {
		switch server.environment {
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

func (s *Server) Start() {
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

func (s *Server) registerDiscovery() *Server {
	if s.options.clientDiscovery == nil {
		return nil
	}
	serverDiscovery, ok := s.options.clientDiscovery.(discovery.Discovery)
	if ok {
		err := serverDiscovery.Register()
		if err != nil {
			s.options.logger.Error(err.Error())
		}
	}
	return s
}

func (s *Server) deregister() error {
	if s.options.clientDiscovery == nil {
		return nil
	}
	serverDiscovery, ok := s.options.clientDiscovery.(discovery.Discovery)
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
	if err := s.options.httpServer.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "shutdown http server error")
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
		if err := s.stop(); err != nil {
			s.options.logger.Error("stop http server error %s", err.Error())
		}
		os.Exit(0)
	}
}

func (s *Server) printLog() {
	if s.options.logger == nil {
		log.Printf("======================================================================")
		log.Println(console_colors.Green("Starting server..."))
		log.Println(console_colors.Green(fmt.Sprintf("Welcome to %s, starting application ...", s.options.serviceName)))
		log.Println(fmt.Sprintf("framework version :  %s", console_colors.Blue(s.version)))
		log.Println(fmt.Sprintf("server & protocol        :  %s", console_colors.Green(s.server)))
		log.Println(fmt.Sprintf("machine host ip          :  %s", console_colors.Blue(s.options.host)))
		log.Println(fmt.Sprintf("listening on port        :  %s", console_colors.Blue(fmt.Sprintf("%d", s.options.port))))
		log.Println(fmt.Sprintf("application running pid  :  %s", console_colors.Blue(strconv.Itoa(s.pId))))
		log.Println(fmt.Sprintf("application name         :  %s", console_colors.Blue(s.options.serviceName)))
		log.Println(fmt.Sprintf("application exec path    :  %s", console_colors.Yellow(files.GetCurrentDirectory())))
		log.Println(fmt.Sprintf("application environment  :  %s", console_colors.Yellow(console_colors.Blue(s.environment))))
		log.Println(fmt.Sprintf("running in %s mode , change (Dev,Test,Prod) mode by HostBuilder.SetEnvironment .", console_colors.Red(s.environment)))
		log.Println(console_colors.Green("Server is Started."))
		log.Printf("======================================================================")
		return
	}
	s.options.logger.Info("======================================================================")
	s.options.logger.Info(console_colors.Green("Starting server..."))
	s.options.logger.Info(console_colors.Green(fmt.Sprintf("Welcome to %s, starting application ...", s.options.serviceName)))
	s.options.logger.Info("framework version :  %s", console_colors.Blue(s.version))
	s.options.logger.Info("server & protocol        :  %s", console_colors.Green(s.server))
	s.options.logger.Info("machine host ip          :  %s", console_colors.Blue(s.options.host))
	s.options.logger.Info("listening on port        :  %s", console_colors.Blue(fmt.Sprintf("%d", s.options.port)))
	s.options.logger.Info("application running pid  :  %s", console_colors.Blue(strconv.Itoa(s.pId)))
	s.options.logger.Info("application name         :  %s", console_colors.Blue(s.options.serviceName))
	s.options.logger.Info("application exec path    :  %s", console_colors.Yellow(files.GetCurrentDirectory()))
	s.options.logger.Info("application environment  :  %s", console_colors.Yellow(console_colors.Blue(s.environment)))
	s.options.logger.Info("running in %s mode , change (Dev,Test,Prod) mode by Environment .", console_colors.Red(s.environment))
	s.options.logger.Info(console_colors.Green("Server is Started."))
	s.options.logger.Info("======================================================================")
}

func (s *Server) IsDevelopment() bool {
	return s.environment == Dev
}

func (s *Server) IsTest() bool {
	return s.environment == Test
}

func (s *Server) IsProduction() bool {
	return s.environment == Prod
}
