package main

import (
	"bufio"
	"io"
	"net"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

// HTTP 处理器
type HTTPHandler struct {
	logger *zap.Logger
}

// 新建 HTTP 处理器
func NewHTTPHandler(logger *zap.Logger) *HTTPHandler {
	return &HTTPHandler{logger: logger}
}

// 处理 HTTP 请求
func (h *HTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{}
	req, err := http.NewRequest(r.Method, r.URL.String(), r.Body)
	if err != nil {
		h.logger.Error("请求创建失败", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req.Header = r.Header

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		h.logger.Error("请求转发失败", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// 复制响应头
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

// HTTPS 处理器
type HTTPSHandler struct {
	logger *zap.Logger
}

// 新建 HTTPS 处理器
func NewHTTPSHandler(logger *zap.Logger) *HTTPSHandler {
	return &HTTPSHandler{logger: logger}
}

// 处理 HTTPS 连接
func (h *HTTPSHandler) Handle(conn net.Conn) {
	defer conn.Close()

	// 读取请求
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		h.logger.Warn("连接读取失败", zap.Error(err))
		return
	}

	// 解析 CONNECT 请求
	requestLine := string(buf[:n])
	if strings.HasPrefix(requestLine, "CONNECT") {
		parts := strings.Split(requestLine, " ")
		if len(parts) >= 2 {
			req, err := http.ReadRequest(bufio.NewReader(strings.NewReader(requestLine)))
			if err == nil {
				h.handleConnection(conn, req)
				return
			}
		}
	}

	h.logger.Warn("未知请求类型", zap.String("request", requestLine))
}

// 处理 HTTPS 连接
func (h *HTTPSHandler) handleConnection(conn net.Conn, r *http.Request) {
	destConn, err := net.Dial("tcp", r.URL.Host)
	if err != nil {
		h.logger.Error("连接目标失败", zap.Error(err))
		return
	}
	defer destConn.Close()

	// 发送 200 状态响应
	_, err = conn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))
	if err != nil {
		return
	}

	go io.Copy(destConn, conn)
	io.Copy(conn, destConn)
}
