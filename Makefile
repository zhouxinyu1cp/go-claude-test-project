.PHONY: all build test clean install

# 默认目标
all: build

# 构建 Web 服务 (issue2md CLI)
build:
	go build -o bin/issue2md ./cmd/issue2md

# 运行所有测试
test:
	go test -v ./...

# 清理构建产物
clean:
	rm -rf bin/

# 安装依赖
install:
	go mod download

# 格式化代码
fmt:
	go fmt ./...

# 代码检查
vet:
	go vet ./...
