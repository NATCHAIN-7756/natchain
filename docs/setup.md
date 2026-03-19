# 开发环境搭建

## 系统要求

| 项目 | 最低配置 | 推荐配置 |
|------|----------|----------|
| CPU | 2 核 | 4 核 |
| 内存 | 4 GB | 8 GB |
| 硬盘 | 40 GB SSD | 80 GB SSD |
| 系统 | Ubuntu 20.04 | Ubuntu 22.04 |

## 1. 安装 Go

```bash
# 下载 Go 1.21
wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz

# 解压
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz

# 配置环境变量
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc
source ~/.bashrc

# 验证
go version
```

## 2. 安装 Docker

```bash
# 安装 Docker
curl -fsSL https://get.docker.com | sh

# 添加当前用户到 docker 组
sudo usermod -aG docker $USER

# 启动 Docker
sudo systemctl start docker
sudo systemctl enable docker

# 验证
docker --version
```

## 3. 安装必要工具

```bash
# 安装编译工具
sudo apt update
sudo apt install -y build-essential git curl wget jq

# 安装 Protobuf 编译器
sudo apt install -y protobuf-compiler

# 安装 Go 相关工具
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## 4. 安装 Cosmos SDK

```bash
# 克隆 Cosmos SDK
git clone https://github.com/cosmos/cosmos-sdk.git
cd cosmos-sdk
git checkout v0.47.5

# 安装
make install

# 验证
simd version
```

## 5. 安装 Ignite CLI（可选）

```bash
# 安装 Ignite
curl https://get.ignite.com/cli! | bash

# 验证
ignite version
```

## 6. 克隆 NatChain 仓库

```bash
# 克隆
git clone https://github.com/NATCHAIN-7756/natchain.git
cd natchain

# 安装依赖
go mod download
```

## 开发工具推荐

| 工具 | 用途 |
|------|------|
| VS Code + Go 插件 | 代码编辑 |
| Postman | API 测试 |
| Docker Compose | 多节点部署 |

## 验证环境

```bash
# 检查所有工具
go version          # Go 1.21+
docker --version    # Docker 24+
protoc --version    # libprotoc 3+
ignite version      # Ignite 28+ (可选)
```
