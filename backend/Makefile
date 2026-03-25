.PHONY: build run test clean docker

# 变量
BINARY_NAME=aniya-blog-server
GO=go
GOFLAGS=-v

# 默认目标
all: build

# 下载依赖
deps:
	$(GO) mod download
	$(GO) mod tidy

# 构建
build:
	$(GO) build $(GOFLAGS) -o bin/$(BINARY_NAME) cmd/server/main.go

# 运行
run:
	$(GO) run cmd/server/main.go

# 测试
test:
	$(GO) test -v ./...

# 清理
clean:
	rm -rf bin/
	$(GO) clean

# 生成 Swagger 文档
swagger:
	swag init -g cmd/server/main.go -o ./api/docs

# Docker 构建
docker:
	docker build -t $(BINARY_NAME):latest .

# 格式化代码
fmt:
	$(GO) fmt ./...

# 代码检查
lint:
	golangci-lint run

# 数据库迁移
migrate:
	$(GO) run cmd/server/main.go -migrate

# 创建初始管理员
init-admin:
	$(GO) run cmd/server/init_admin.go

# 帮助
help:
	@echo "Makefile commands:"
	@echo "  make deps     - Download dependencies"
	@echo "  make build    - Build the binary"
	@echo "  make run      - Run the server"
	@echo "  make test     - Run tests"
	@echo "  make clean    - Clean build artifacts"
	@echo "  make swagger  - Generate Swagger docs"
	@echo "  make docker   - Build Docker image"
	@echo "  make fmt      - Format code"
	@echo "  make lint     - Run linter"
	@echo "  make migrate  - Run database migrations"
	@echo "  make init-admin - Create initial admin user"
