# NatChain 白皮书 v1.1

## 一、项目概述

| 项目 | 内容 |
|------|------|
| 名称 | NatChain |
| 定位 | 支付型公链 |
| 技术栈 | Cosmos SDK + CometBFT |
| 代币 | NAT |
| 官网 | natjen.com |

---

## 二、代币经济模型

| 项目 | 设计 |
|------|------|
| 总量 | 无上限 |
| 产出方式 | 矿机挖矿 |
| 初始区块奖励 | 50 NAT |
| 减半周期 | 每 3 年减半 |
| 手续费 | 归开发者地址（程序启动时自动生成） |

---

## 三、矿机挖矿机制

### 核心规则
| 规则 | 内容 |
|------|------|
| 矿机上限 | **1000 台（永久固定）** |
| 价格 | 100 USDT / 台 |
| 绑定方式 | 首次通电联网 → 绑定钱包 → 自动生成绑定 |
| 硬件验证 | CPU序列号 + CPU型号 |
| 挖矿奖励 | 按出块分配，3年减半 |
| 离线惩罚 | 超过 180 天自动解绑 |

### 绑定流程
```
1. 购买矿机（100 USDT）
2. 矿机首次通电联网
3. 用户绑定钱包地址
4. 系统自动生成绑定关系
5. 开始挖矿
```

### 数据结构
```go
type Miner struct {
    Index         uint64       // 矿机编号 1-1000
    CpuSerial     string       // CPU序列号
    CpuModel      string       // CPU型号
    Owner         string       // 所有者地址
    ValidatorAddr string       // 验证者地址
    BondedAt      int64        // 绑定时间戳
    LastActive    int64        // 最后出块时间
    Status        MinerStatus  // 待绑定/已绑定/已解绑
    Paid          bool         // 是否已付款
}
```

---

## 四、钱包安全机制

### 设计理念
- 风险可控：设置转出额度，钱包丢失时损失有限
- 自动转出：余额超过阈值自动转到冷钱包

### 安全额度
| 功能 | 说明 |
|------|------|
| 单笔限额 | 单笔转出上限 |
| 每日限额 | 每日转出累计上限 |
| 自动转出阈值 | 余额超过此值触发自动转出 |
| 冷钱包地址 | 自动转出的目标地址 |

### 数据结构
```go
type WalletSecurity struct {
    Owner           string    // 钱包地址
    DailyLimit      sdk.Coin  // 每日转出上限
    SingleLimit     sdk.Coin  // 单笔转出上限
    Threshold       sdk.Coin  // 自动转出阈值
    AutoTransferTo  string    // 冷钱包地址
    AutoTransferAmt sdk.Coin  // 自动转出金额
    Enabled         bool      // 是否启用
}
```

### 风险控制效果
| 场景 | 结果 |
|------|------|
| 钱包私钥泄露 | 损失 ≤ 每日额度 |
| 大额收款 | 自动转到冷钱包 |
| 冷钱包丢失 | 离线存储，风险可控 |

### 伙伴管理
```go
type Partner struct {
    Owner      string    // 用户地址
    Name       string    // 伙伴名称
    Address    string    // 伙伴地址
    Memo       string    // 备注
    CreatedAt  int64     // 添加时间
}
```

---

## 五、商家 API

### 核心设计
- **无需注册：** 钱包地址即可申请 API Key
- **功能完整：** 收款、转账、查询、回调

### 申请流程
```
创建钱包 → 申请 API Key → 直接对接
```

### 数据结构
```go
type MerchantApiKey struct {
    Id           string       // API Key ID
    Owner        string       // 钱包地址
    ApiKey       string       // API Key
    ApiSecret    string       // API Secret
    CallbackUrl  string       // 回调地址
    DailyLimit   sdk.Coin     // API 调用限额
    CreatedAt    int64
    Status       ApiKeyStatus
}
```

### REST API
```
POST /api/v1/apikey/apply         # 申请 API Key
GET  /api/v1/apikey/info          # 查询 API Key 信息
POST /api/v1/payment/create       # 创建收款订单
GET  /api/v1/payment/query/:id    # 查询订单
POST /api/v1/payment/verify       # 验证交易
GET  /api/v1/balance              # 查询余额
POST /api/v1/transfer             # 转账
GET  /api/v1/transactions         # 交易记录
```

### 安全机制
- API Key + HMAC 签名
- 受钱包安全额度限制
- IP 白名单（可选）

---

## 六、跨链功能

- 原生支持 Cosmos IBC 协议
- 与 Cosmos 生态互通（ATOM, OSMO 等）
- 通过桥接扩展到 ETH、BSC

---

## 七、治理机制

**暂不开启链上治理**

理由：
1. 项目早期，需要快速迭代
2. 开发者决策更高效
3. 等社区成熟后再考虑

---

## 八、模块架构

```
NatChain
├── x/miner      # 矿机绑定、出块验证、离线检测
├── x/wallet     # 安全额度、自动转出、伙伴管理
└── x/merchant   # API Key、收款订单、回调通知
```

---

## 九、开发进度

### ✅ 已完成
- 项目初始化（Cosmos SDK + CometBFT）
- 编译构建（`natchaind` 可执行文件）
- 节点初始化、创建钱包、创世账户
- 启动节点（CometBFT 共识运行）
- P2P 网络（端口 26656）
- RPC 服务（端口 26657）

### 🔄 待开发
| 序号 | 任务 | 优先级 | 模块 |
|------|------|--------|------|
| 1 | 创建 x/miner 模块 | 高 | miner |
| 2 | 实现矿机绑定（CPU验证） | 高 | miner |
| 3 | 集成到共识层验证 | 高 | miner |
| 4 | 实现离线检测与自动解绑 | 中 | miner |
| 5 | 创建 x/wallet 模块 | 高 | wallet |
| 6 | 实现安全额度设置 | 高 | wallet |
| 7 | 实现转出额度检查 | 高 | wallet |
| 8 | 实现自动转出逻辑 | 中 | wallet |
| 9 | 实现伙伴管理 | 低 | wallet |
| 10 | 创建 x/merchant 模块 | 高 | merchant |
| 11 | 实现 API Key 申请 | 高 | merchant |
| 12 | 实现收款订单 API | 高 | merchant |
| 13 | 实现转账与查询 API | 中 | merchant |
| 14 | 实现回调通知 | 中 | merchant |
| 15 | CLI 命令完善 | 低 | all |

---

## 十、技术要点

1. **项目位置：** `~/projects/natchain/`
2. **Go 环境：** 需要 `~/local/go/bin` 在 PATH 中
3. **编译命令：** `export PATH=$HOME/local/go/bin:$PATH && go build -o natchaind ./cmd/natchaind`
4. **数据持久化：** Docker 卷 `natjen` 挂载到 `/home/node/projects`

---

*文档版本：v1.1*
*更新时间：2026-03-16*
