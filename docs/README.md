# NatChain 文档中心

## 概述

NatChain 是基于 Cosmos SDK 构建的区块链网络，支持矿机挖矿、商家支付和钱包管理功能。

---

## 文档索引

| 文档 | 说明 |
|------|------|
| [快速开始](./quick-start.md) | 安装、初始化、启动节点 |
| [CLI 命令手册](./cli-reference.md) | 所有 CLI 命令详解 |
| [商家 API 指南](./api/merchant-api.md) | 商家接入 API 文档 |
| [运维文档](./operations.md) | 服务管理、监控、故障排除 |
| [白皮书](./natchain-whitepaper.md) | 项目白皮书 |
| [架构设计](./architecture.md) | 系统架构文档 |
| [开发路线图](./roadmap.md) | 项目规划 |

---

## 核心模块

### x/miner - 矿机模块
- 矿机注册与管理
- CPU 序列号绑定
- 心跳检测
- 离线自动解绑（180天）

### x/wallet - 钱包模块
- 安全额度控制
- 冷钱包地址设置
- 合伙人分成管理
- 自动转账功能

### x/merchant - 商家模块
- 商家注册
- API Key 管理
- 支付创建与确认
- Webhook 回调

---

## 快速链接

- **测试网 RPC:** http://38.76.200.118:26657
- **链 ID:** natchain
- **GitHub:** https://github.com/NATCHAIN-7756/natchain
- **官网:** natjen.com

---

## 模块详细文档

- [miner 模块](./modules/miner.md)
- [wallet 模块](./modules/wallet.md)
- [merchant 模块](./modules/merchant.md)
