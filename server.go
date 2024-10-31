package main

import (
	"go.uber.org/zap"
	"net"
	"net/http"
	"strconv"
)

// 代理服务器结构体
type Server struct {
	config *Config
	logger *zap.Logger
}

// 新建代理服务器
func NewServer(config Config, logger *zap.Logger) *Server {
	return &Server{config: &config, logger: logger}
}

// 启动服务器
func (s *Server) Start() error {
	httpHandler := NewHTTPHandler(s.logger)
	httpsHandler := NewHTTPSHandler(s.logger)

	// 启动 HTTP 服务器，使用不同的端口
	httpPort := s.config.Port + 1 // 假设 HTTPS 在 s.config.Port，HTTP 在 s.config.Port + 1
	go func() {
		s.logger.Info("启动 HTTP 服务器", zap.Int("port", httpPort))
		if err := http.ListenAndServe(":"+strconv.Itoa(httpPort), httpHandler); err != nil {
			s.logger.Fatal("HTTP 服务器启动失败", zap.Error(err))
		}
	}()

	// 启动 HTTPS 代理
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(s.config.Port))
	if err != nil {
		return err
	}
	s.logger.Info("启动 HTTPS 代理", zap.Int("port", s.config.Port))

	for {
		conn, err := listener.Accept()
		if err != nil {
			s.logger.Warn("连接接受失败", zap.Error(err))
			continue
		}

		go func(c net.Conn) {
			defer c.Close() // 确保连接被关闭
			httpsHandler.Handle(c)
		}(conn)
	}
}
