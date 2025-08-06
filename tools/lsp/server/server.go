package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
)

// Server LSP 服务器
type Server struct {
	reader    io.Reader
	writer    io.Writer
	documents map[string]string // 文档缓存，key 为 URI，value 为文档内容
}

// NewServer 创建新的 LSP 服务器
func NewServer() *Server {
	return &Server{
		documents: make(map[string]string),
	}
}

// Start 启动 LSP 服务器
func (s *Server) Start(ctx context.Context, reader io.Reader, writer io.Writer) error {
	s.reader = reader
	s.writer = writer

	log.Println("Origami LSP 服务器启动")

	// 处理 LSP 消息循环
	return s.messageLoop(ctx)
}

// messageLoop 处理消息循环
func (s *Server) messageLoop(ctx context.Context) error {
	buffer := make([]byte, 4096)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			n, err := s.reader.Read(buffer)
			if err != nil {
				if err == io.EOF {
					return nil
				}
				return fmt.Errorf("读取消息失败: %w", err)
			}

			if n > 0 {
				if err := s.handleMessage(buffer[:n]); err != nil {
					log.Printf("处理消息失败: %v", err)
				}
			}
		}
	}
}

// handleMessage 处理单个 LSP 消息
func (s *Server) handleMessage(data []byte) error {
	message := string(data)

	// 解析 LSP 消息头
	parts := strings.Split(message, "\r\n\r\n")
	if len(parts) < 2 {
		return fmt.Errorf("无效的 LSP 消息格式")
	}

	headers := parts[0]
	content := parts[1]

	// 解析 Content-Length
	var contentLength int
	for _, line := range strings.Split(headers, "\r\n") {
		if strings.HasPrefix(line, "Content-Length:") {
			lengthStr := strings.TrimSpace(strings.TrimPrefix(line, "Content-Length:"))
			var err error
			contentLength, err = strconv.Atoi(lengthStr)
			if err != nil {
				return fmt.Errorf("解析 Content-Length 失败: %w", err)
			}
			break
		}
	}

	// 确保内容长度正确
	if len(content) < contentLength {
		return fmt.Errorf("消息内容不完整")
	}

	// 解析 JSON-RPC 消息
	var request map[string]interface{}
	if err := json.Unmarshal([]byte(content[:contentLength]), &request); err != nil {
		return fmt.Errorf("解析 JSON-RPC 消息失败: %w", err)
	}

	// 处理请求
	return s.handleRequest(request)
}

// handleRequest 处理 LSP 请求
func (s *Server) handleRequest(request map[string]interface{}) error {
	method, ok := request["method"].(string)
	if !ok {
		return fmt.Errorf("无效的方法名")
	}

	log.Printf("收到 LSP 请求: %s", method)

	switch method {
	case "initialize":
		return s.handleInitialize(request)
	case "initialized":
		return s.handleInitialized(request)
	case "textDocument/didOpen":
		return s.handleTextDocumentDidOpen(request)
	case "textDocument/didChange":
		return s.handleTextDocumentDidChange(request)
	case "textDocument/completion":
		return s.handleTextDocumentCompletion(request)
	case "textDocument/hover":
		return s.handleTextDocumentHover(request)
	case "shutdown":
		return s.handleShutdown(request)
	case "exit":
		return s.handleExit(request)
	default:
		log.Printf("未处理的方法: %s", method)
		return nil
	}
}

// sendResponse 发送响应
func (s *Server) sendResponse(id interface{}, result interface{}) error {
	response := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"result":  result,
	}

	return s.sendMessage(response)
}

// sendMessage 发送消息
func (s *Server) sendMessage(message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("序列化消息失败: %w", err)
	}

	header := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(data))

	if _, err := s.writer.Write([]byte(header)); err != nil {
		return fmt.Errorf("写入消息头失败: %w", err)
	}

	if _, err := s.writer.Write(data); err != nil {
		return fmt.Errorf("写入消息体失败: %w", err)
	}

	return nil
}
