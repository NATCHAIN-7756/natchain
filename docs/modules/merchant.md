# x/merchant 商家模块

## 概述

提供商家注册和支付处理功能。

---

## 功能

- 商家注册与管理
- API Key 自动生成
- 支付创建与确认
- Webhook 回调通知

---

## 交易命令

### register

注册成为商家

```bash
natchaind tx merchant register \
  --name "商家名称" \
  --callback-url "https://example.com/callback" \
  --commission-rate 5 \
  --from <your-key> \
  --chain-id natchain \
  --keyring-backend test \
  -y
```

### create-payment

创建支付

```bash
natchaind tx merchant create-payment \
  --merchant-address <商家地址> \
  --amount 100stake \
  --callback-data "订单ID" \
  --from <user-key> \
  --chain-id natchain \
  --keyring-backend test \
  -y
```

### confirm-payment

确认支付

```bash
natchaind tx merchant confirm-payment \
  --payment-id PAY-1 \
  --from <merchant-key> \
  --chain-id natchain \
  --keyring-backend test \
  -y
```

### regenerate-api-key

重新生成 API Key

```bash
natchaind tx merchant regenerate-api-key \
  --from <merchant-key> \
  --chain-id natchain \
  --keyring-backend test \
  -y
```

### update-merchant

更新商家信息

```bash
natchaind tx merchant update-merchant \
  --name "新名称" \
  --callback-url "https://new-url.com/callback" \
  --from <merchant-key> \
  --chain-id natchain \
  --keyring-backend test \
  -y
```

---

## 查询命令

### merchant

查询单个商家

```bash
natchaind q merchant merchant --address <商家地址>
```

返回：

```json
{
  "merchant": {
    "address": "natchain1...",
    "name": "商家名称",
    "api_key": "生成的API密钥",
    "callback_url": "https://...",
    "commission_rate": "5",
    "active": true,
    "created_at": "时间戳"
  }
}
```

### merchants

查询所有商家

```bash
natchaind q merchant merchants
```

### payment-record

查询单笔支付

```bash
natchaind q merchant payment-record --id PAY-1
```

### payment-records-by-merchant

查询商家所有支付

```bash
natchaind q merchant payment-records-by-merchant \
  --merchant-address <商家地址>
```

### params

查询模块参数

```bash
natchaind q merchant params
```

---

## 数据结构

### Merchant

| 字段 | 类型 | 说明 |
|------|------|------|
| address | string | 商家地址 |
| name | string | 商家名称 |
| api_key | string | API 密钥 |
| callback_url | string | 回调地址 |
| commission_rate | string | 佣金比例 |
| active | bool | 是否活跃 |
| created_at | string | 创建时间 |

### PaymentRecord

| 字段 | 类型 | 说明 |
|------|------|------|
| id | string | 支付ID (PAY-xxx) |
| merchant_address | string | 商家地址 |
| user_address | string | 用户地址 |
| amount | string | 支付金额 |
| status | string | 状态 (pending/confirmed/failed) |
| timestamp | string | 时间戳 |
| callback_data | string | 自定义回调数据 |

---

## 支付状态

| 状态 | 说明 |
|------|------|
| pending | 待确认 |
| confirmed | 已确认 |
| failed | 失败 |

---

## Webhook 回调

支付状态变更时，系统向 `callback_url` 发送 POST 请求：

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

## 状态码

| Code | 说明 |
|------|------|
| 1 | 商家不存在 |
| 2 | 商家已存在 |
| 3 | 支付不存在 |
| 4 | 无权限操作 |
| 5 | 支付已处理 |
