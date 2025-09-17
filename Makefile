# Go 参数
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# 项目信息
PROJECT_NAME=km-cicd-test
VERSION=1.0.0

# 构建目录
BUILD_DIR=build
BIN_DIR=$(BUILD_DIR)/bin

# 模块路径
MODULE1_PATH=./module1
MODULE2_PATH=./module2
COMMON_PATH=./common

.PHONY: all build clean test deps run-module1 run-module2 help

# 默认目标
all: deps build

# 安装依赖
deps:
	@echo "安装依赖..."
	$(GOMOD) download
	$(GOMOD) tidy

# 构建所有模块
build: build-module1 build-module2

# 构建 module1
build-module1:
	@echo "构建 module1..."
	@mkdir -p $(BIN_DIR)
	cd $(MODULE1_PATH) && $(GOBUILD) -o ../$(BIN_DIR)/module1 main.go

# 构建 module2
build-module2:
	@echo "构建 module2..."
	@mkdir -p $(BIN_DIR)
	cd $(MODULE2_PATH) && $(GOBUILD) -o ../$(BIN_DIR)/module2 main.go

# 运行 module1
run-module1: build-module1
	@echo "启动 module1..."
	./$(BIN_DIR)/module1

# 运行 module2
run-module2: build-module2
	@echo "启动 module2..."
	./$(BIN_DIR)/module2

# 并行运行两个模块
run-all: build
	@echo "启动所有模块..."
	./$(BIN_DIR)/module1 &
	./$(BIN_DIR)/module2 &
	@echo "所有模块已启动"
	@echo "module1: http://localhost:8080"
	@echo "module2: http://localhost:8081"

# 运行测试
test:
	@echo "运行测试..."
	$(GOTEST) -v ./...

# 清理构建文件
clean:
	@echo "清理构建文件..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

# 格式化代码
fmt:
	@echo "格式化代码..."
	$(GOCMD) fmt ./...

# 代码检查
vet:
	@echo "代码检查..."
	$(GOCMD) vet ./...

# 安装工具
install-tools:
	@echo "安装开发工具..."
	$(GOGET) -u golang.org/x/tools/cmd/goimports
	$(GOGET) -u github.com/golangci/golangci-lint/cmd/golangci-lint

# 代码质量检查
lint:
	@echo "代码质量检查..."
	golangci-lint run

# CI/CD 相关命令
create-temp-branch:
	@echo "创建临时分支..."
	@if [ -z "$(ENV)" ]; then \
		echo "错误: 请指定环境 (ENV=dev 或 ENV=prod)"; \
		exit 1; \
	fi
	./scripts/create-temp-branch.sh -e $(ENV)

deploy:
	@echo "部署到指定环境..."
	@if [ -z "$(ENV)" ]; then \
		echo "错误: 请指定环境 (ENV=staging 或 ENV=production)"; \
		exit 1; \
	fi
	./scripts/deploy.sh -e $(ENV)

deploy-status:
	@echo "检查部署状态..."
	@if [ -z "$(ENV)" ]; then \
		echo "错误: 请指定环境 (ENV=staging 或 ENV=production)"; \
		exit 1; \
	fi
	./scripts/deploy.sh -e $(ENV) --status

rollback:
	@echo "回滚部署..."
	@if [ -z "$(ENV)" ]; then \
		echo "错误: 请指定环境 (ENV=staging 或 ENV=production)"; \
		exit 1; \
	fi
	./scripts/deploy.sh -e $(ENV) --rollback

# Docker 相关命令
docker-build:
	@echo "构建Docker镜像..."
	docker build -t $(PROJECT_NAME):$(VERSION) .
	docker build -t $(PROJECT_NAME):latest .

docker-run:
	@echo "运行Docker容器..."
	docker-compose up -d

docker-stop:
	@echo "停止Docker容器..."
	docker-compose down

docker-logs:
	@echo "查看Docker日志..."
	docker-compose logs -f

# 显示帮助信息
help:
	@echo "可用的命令:"
	@echo "  all          - 安装依赖并构建所有模块"
	@echo "  build        - 构建所有模块"
	@echo "  build-module1- 构建 module1"
	@echo "  build-module2- 构建 module2"
	@echo "  run-module1  - 运行 module1"
	@echo "  run-module2  - 运行 module2"
	@echo "  run-all      - 并行运行所有模块"
	@echo "  test         - 运行测试"
	@echo "  clean        - 清理构建文件"
	@echo "  fmt          - 格式化代码"
	@echo "  vet          - 代码检查"
	@echo "  lint         - 代码质量检查"
	@echo "  deps         - 安装依赖"
	@echo ""
	@echo "CI/CD 命令:"
	@echo "  create-temp-branch ENV=dev|prod - 创建临时分支"
	@echo "  deploy ENV=dev|prod - 部署到指定环境"
	@echo "  deploy-status ENV=dev|prod - 检查部署状态"
	@echo "  rollback ENV=dev|prod - 回滚部署"
	@echo ""
	@echo "Docker 命令:"
	@echo "  docker-build - 构建Docker镜像"
	@echo "  docker-run   - 运行Docker容器"
	@echo "  docker-stop  - 停止Docker容器"
	@echo "  docker-logs  - 查看Docker日志"
	@echo ""
	@echo "  help         - 显示此帮助信息"
