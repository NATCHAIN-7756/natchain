# NatChain - 支付型公链

基于 Cosmos SDK + CometBFT 构建的支付型公链。

## 特性

- 🏗️ **矿机挖矿** - 1000 台矿机上限，100 USDT/台
- 💰 **3年减半** - 借鉴比特币机制
- 🔐 **钱包安全** - 安全额度 + 自动转冷钱包
- 🛒 **商家 API** - 钱包地址申请 API Key，无需注册
- 🌐 **跨链支持** - Cosmos IBC 协议

## 快速开始

### 环境要求

- Go 1.21+
- Docker
- Ubuntu 22.04 (推荐 2核4G 以上)

### 安装

```bash
git clone https://github.com/NATCHAIN-7756/natchain.git
cd natchain
make install
```

### 运行单节点

```bash
# 初始化
natd init test --chain-id natchain

# 创建账户
natd keys add alice

# 添加创世账户
natd add-genesis-account alice 100000000000unat

# 启动
natd start
```

## 文档

| 文档 | 说明 |
|------|------|
| [开发路线图](docs/roadmap.md) | 开发计划和里程碑 |
| [环境搭建](docs/setup.md) | 开发环境配置 |
| [架构设计](docs/architecture.md) | 系统架构说明 |
| [API 文档](docs/api.md) | REST/gRPC API 接口 |
| [白皮书](docs/natchain-whitepaper.md) | 项目白皮书 |

## 模块设计

| 模块 | 文档 | 状态 |
|------|------|------|
| x/miner | [设计文档](docs/modules/miner.md) | 📝 设计中 |
| x/wallet | [设计文档](docs/modules/wallet.md) | 📝 设计中 |
| x/merchant | [设计文档](docs/modules/merchant.md) | 📝 设计中 |

## 技术栈

- **区块链框架**: Cosmos SDK v0.47+
- **共识引擎**: CometBFT (Tendermint)
- **开发语言**: Go 1.21+
- **API**: REST + gRPC

## 代币经济

| 项目 | 设计 |
|------|------|
| 代币符号 | NAT |
| 总量 | 无上限 |
| 初始奖励 | 50 NAT/块 |
| 减半周期 | 3年 |
| 矿机上限 | 1000 台 |

## 贡献

欢迎提交 Issue 和 Pull Request。

## License

MIT
