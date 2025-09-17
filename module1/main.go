package main

import (
	"fmt"
	"net/http"

	"github.com/leon/km-cicd-test/common"
)

// User 用户结构
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UserService 用户服务
type UserService struct {
	logger common.Logger
	users  []User
}

// NewUserService 创建新的用户服务
func NewUserService() *UserService {
	return &UserService{
		logger: common.NewSimpleLogger("UserService"),
		users: []User{
			{ID: 1, Name: "张三", Email: "zhangsan@example.com"},
			{ID: 2, Name: "李四", Email: "lisi@example.com"},
		},
	}
}

// GetUsers 获取所有用户
func (s *UserService) GetUsers() []User {
	s.logger.Info("获取用户列表，共 %d 个用户", len(s.users))
	return s.users
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(id int) (*User, error) {
	s.logger.Info("查找用户 ID: %d", id)
	for _, user := range s.users {
		if user.ID == id {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("用户不存在: ID %d", id)
}

// AddUser 添加用户
func (s *UserService) AddUser(name, email string) User {
	newID := len(s.users) + 1
	user := User{
		ID:    newID,
		Name:  name,
		Email: email,
	}
	s.users = append(s.users, user)
	s.logger.Info("添加新用户: %s (%s)", name, email)
	return user
}

// 处理获取用户列表的HTTP请求
func handleGetUsers(w http.ResponseWriter, r *http.Request, service *UserService) {
	users := service.GetUsers()
	response := common.SuccessResponse(users)

	w.Header().Set("Content-Type", "application/json")
	jsonStr, err := response.ToJSON()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonStr))
}

// 处理根据ID获取用户的HTTP请求
func handleGetUserByID(w http.ResponseWriter, r *http.Request, service *UserService) {
	// 这里简化处理，实际应该从URL路径中解析ID
	user, err := service.GetUserByID(1)
	if err != nil {
		response := common.ErrorResponse(404, err.Error())
		jsonStr, _ := response.ToJSON()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(jsonStr))
		return
	}

	response := common.SuccessResponse(user)
	jsonStr, err := response.ToJSON()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonStr))
}

func main() {
	config := common.LoadConfig()
	logger := common.NewSimpleLogger("Module1")

	logger.Info("启动 Module1 服务")
	logger.Info("配置: Host=%s, Port=%d, Debug=%t", config.Host, config.Port, config.Debug)

	userService := NewUserService()

	// 设置路由
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			handleGetUsers(w, r, userService)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			handleGetUserByID(w, r, userService)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// 健康检查端点
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		response := common.SuccessResponse(map[string]string{
			"status": "healthy",
			"module": "module1",
		})
		jsonStr, _ := response.ToJSON()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(jsonStr))
	})

	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	logger.Info("服务启动在 %s", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		logger.Error("服务启动失败: %v", err)
	}
}
