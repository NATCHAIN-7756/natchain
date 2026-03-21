# NatChain 商家 API 接入指南

## 概述

NatChain 商家模块提供支付处理功能，商家可以通过 API 接收用户的链上支付。

---

## 1. 商家注册

### CLI 注册

```bash
natchaind tx merchant register \
  --name "您的商家名称" \
  --callback-url "https://your-domain.com/callback" \
  --commission-rate 5 \
  --from <your-key> \
  --chain-id natchain \
  --keyring-backend test \
  -y
```

### 参数说明

| 参数 | 类型 | 说明 |
|------|------|------|
| name | string | 商家名称 |
| callback-url | string | 支付回调地址 |
| commission-rate | int | 佣金比例 (0-100) |

### 返回结果

```json
{
  "merchant": {
    "address": "natchain1...",
    "name": "您的商家名称",
    "api_key": "生成的API密钥",
    "callback_url": "https://your-domain.com/callback",
    "commission_rate": "5",
    "active": true,
    "created_at": "时间戳"
  }
}
```

---

## 2. 创建支付

### CLI 创建

```bash
natchaind tx merchant create-payment \
  --merchant-address "商家地址" \
  --amount "100stake" \
  --callback-data "订单ID或自定义数据" \
  --from <user-key> \
  --chain-id natchain \
  --keyring-backend test \
  -y
```

### 参数说明

| 参数 | 类型 | 说明 |
|------|------|------|
| merchant-address | string | 商家链地址 |
| amount | string | 支付金额（如 100stake） |
| callback-data | string | 可选，自定义回调数据 |

### 返回结果

```json
{
  "payment_record": {
    "id": "PAY-1",
    "merchant_address": "natchain1...",
    "user_address": "natchain1...",
    "amount": "100stake",
    "status": "pending",
    "timestamp": "时间戳"
  }
}
```

---

## 3. 查询支付

### 查询单笔支付

```bash
natchaind q merchant payment-record --id "PAY-1"
```

### 查询商家所有支付

```bash
natchaind q merchant payment-records-by-merchant \
  --merchant-address "natchain1..."
```

### 支付状态

| 状态 | 说明 |
|------|------|
| pending | 待确认 |
| confirmed | 已确认 |
| failed | 失败 |

---

## 4. 确认支付

商家确认收到支付后调用：

```bash
natchaind tx merchant confirm-payment \
  --payment-id "PAY-1" \
  --from <merchant-key> \
  --chain-id natchain \
  --keyring-backend test \
  -y
```

---

## 5. API Key 管理

### 查看 API Key

```bash
natchaind q merchant merchant --address "natchain1..."
```

### 重新生成 API Key

```bash
natchaind tx merchant regenerate-api-key \
  --from <merchant-key> \
  --chain-id natchain \
  --keyring-backend test \
  -y
```

---

## 6. 回调通知

当支付状态变更时，系统会向商家配置的 `callback_url` 发送 POST 请求：

```json
{
  "payment_id": "PAY-1",
  "merchant_address": "natchain1...",
  "user_address": "natchain1...",
  "amount": "100stake",
  "status": "confirmed",
  "timestamp": "1774072145",
  "callback_data": "订单ID"
}
```

---

## 7. 错误处理

| 错误码 | 说明 |
|--------|------|
| 1 | 商家不存在 |
| 2 | 支付不存在 |
| 3 | 无权限操作 |
| 4 | 支付已处理 |

---

## 8. 最佳实践

1. **安全存储 API Key** - 不要在客户端暴露
2. **验证回调签名** - 确认请求来自 NatChain
3. **幂等处理** - 同一支付可能收到多次回调
4. **超时重试** - 设置合理的回调超时时间

---

## 9. 测试环境

- **RPC 节点:** http://38.76.200.118:26657
- **链 ID:** natchain
- **代币:** stake (测试用)

---

## 联系支持

如有问题，请联系 NatChain 团队。
