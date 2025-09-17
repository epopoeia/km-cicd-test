package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/leon/km-cicd-test/common"
)

// Product 产品结构
type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// Order 订单结构
type Order struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	ProductID int       `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Total     float64   `json:"total"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// ProductService 产品服务
type ProductService struct {
	logger   common.Logger
	products []Product
}

// NewProductService 创建新的产品服务
func NewProductService() *ProductService {
	now := time.Now()
	return &ProductService{
		logger: common.NewSimpleLogger("ProductService"),
		products: []Product{
			{ID: 1, Name: "笔记本电脑", Price: 5999.99, Description: "高性能笔记本电脑", CreatedAt: now},
			{ID: 2, Name: "无线鼠标", Price: 99.99, Description: "蓝牙无线鼠标", CreatedAt: now},
			{ID: 3, Name: "机械键盘", Price: 299.99, Description: "青轴机械键盘", CreatedAt: now},
		},
	}
}

// GetProducts 获取所有产品
func (s *ProductService) GetProducts() []Product {
	s.logger.Info("获取产品列表，共 %d 个产品", len(s.products))
	return s.products
}

// GetProductByID 根据ID获取产品
func (s *ProductService) GetProductByID(id int) (*Product, error) {
	s.logger.Info("查找产品 ID: %d", id)
	for _, product := range s.products {
		if product.ID == id {
			return &product, nil
		}
	}
	return nil, fmt.Errorf("产品不存在: ID %d", id)
}

// OrderService 订单服务
type OrderService struct {
	logger common.Logger
	orders []Order
}

// NewOrderService 创建新的订单服务
func NewOrderService() *OrderService {
	return &OrderService{
		logger: common.NewSimpleLogger("OrderService"),
		orders: []Order{},
	}
}

// CreateOrder 创建订单
func (s *OrderService) CreateOrder(userID, productID, quantity int, price float64) Order {
	newID := len(s.orders) + 1
	total := price * float64(quantity)
	order := Order{
		ID:        newID,
		UserID:    userID,
		ProductID: productID,
		Quantity:  quantity,
		Total:     total,
		Status:    "pending",
		CreatedAt: time.Now(),
	}
	s.orders = append(s.orders, order)
	s.logger.Info("创建新订单: 用户ID=%d, 产品ID=%d, 数量=%d, 总价=%.2f", userID, productID, quantity, total)
	return order
}

// GetOrders 获取所有订单
func (s *OrderService) GetOrders() []Order {
	s.logger.Info("获取订单列表，共 %d 个订单", len(s.orders))
	return s.orders
}

// 处理获取产品列表的HTTP请求
func handleGetProducts(w http.ResponseWriter, r *http.Request, service *ProductService) {
	products := service.GetProducts()
	response := common.SuccessResponse(products)

	w.Header().Set("Content-Type", "application/json")
	jsonStr, err := response.ToJSON()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonStr))
}

// 处理根据ID获取产品的HTTP请求
func handleGetProductByID(w http.ResponseWriter, r *http.Request, service *ProductService) {
	// 这里简化处理，实际应该从URL路径中解析ID
	product, err := service.GetProductByID(1)
	if err != nil {
		response := common.ErrorResponse(404, err.Error())
		jsonStr, _ := response.ToJSON()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(jsonStr))
		return
	}

	response := common.SuccessResponse(product)
	jsonStr, err := response.ToJSON()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonStr))
}

// 处理获取订单列表的HTTP请求
func handleGetOrders(w http.ResponseWriter, r *http.Request, service *OrderService) {
	orders := service.GetOrders()
	response := common.SuccessResponse(orders)

	w.Header().Set("Content-Type", "application/json")
	jsonStr, err := response.ToJSON()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(jsonStr))
}

func main() {
	config := common.LoadConfig()
	logger := common.NewSimpleLogger("Module2")

	logger.Info("启动 Module2 服务")
	logger.Info("配置: Host=%s, Port=%d, Debug=%t", config.Host, config.Port, config.Debug)

	productService := NewProductService()
	orderService := NewOrderService()

	// 设置路由
	http.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			handleGetProducts(w, r, productService)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/product", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			handleGetProductByID(w, r, productService)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/orders", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			handleGetOrders(w, r, orderService)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// 健康检查端点
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		response := common.SuccessResponse(map[string]string{
			"status": "healthy",
			"module": "module2",
		})
		jsonStr, _ := response.ToJSON()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(jsonStr))
	})

	addr := fmt.Sprintf("%s:%d", config.Host, config.Port+1) // 使用不同端口避免冲突
	logger.Info("服务启动在 %s", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		logger.Error("服务启动失败: %v", err)
	}
}
