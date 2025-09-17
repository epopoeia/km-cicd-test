package common

import (
	"encoding/json"
	"fmt"
	"time"
)

// Response 通用响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Time    time.Time   `json:"time"`
}

// NewResponse 创建新的响应
func NewResponse(code int, message string, data interface{}) *Response {
	return &Response{
		Code:    code,
		Message: message,
		Data:    data,
		Time:    time.Now(),
	}
}

// SuccessResponse 成功响应
func SuccessResponse(data interface{}) *Response {
	return NewResponse(200, "success", data)
}

// ErrorResponse 错误响应
func ErrorResponse(code int, message string) *Response {
	return NewResponse(code, message, nil)
}

// ToJSON 转换为JSON字符串
func (r *Response) ToJSON() (string, error) {
	bytes, err := json.Marshal(r)
	if err != nil {
		return "", fmt.Errorf("failed to marshal response: %w", err)
	}
	return string(bytes), nil
}

// Logger 简单的日志接口
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
}

// SimpleLogger 简单日志实现
type SimpleLogger struct {
	prefix string
}

// NewSimpleLogger 创建新的简单日志器
func NewSimpleLogger(prefix string) *SimpleLogger {
	return &SimpleLogger{prefix: prefix}
}

// Info 信息日志
func (l *SimpleLogger) Info(msg string, args ...interface{}) {
	fmt.Printf("[INFO] [%s] %s\n", l.prefix, fmt.Sprintf(msg, args...))
}

// Error 错误日志
func (l *SimpleLogger) Error(msg string, args ...interface{}) {
	fmt.Printf("[ERROR] [%s] %s\n", l.prefix, fmt.Sprintf(msg, args...))
}

// Debug 调试日志
func (l *SimpleLogger) Debug(msg string, args ...interface{}) {
	fmt.Printf("[DEBUG] [%s] %s\n", l.prefix, fmt.Sprintf(msg, args...))
}
