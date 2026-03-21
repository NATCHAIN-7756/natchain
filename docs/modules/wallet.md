# x/wallet 钱包模块

## 概述

提供钱包安全管理和自动转账功能。

---

## 功能

- 安全额度控制
- 冷钱包地址设置
- 合伙人分成管理
- 自动转账到冷钱包

---

## 交易命令

### set-safety-limit

设置安全额度（超过此金额需额外验证）

```bash
natchaind tx wallet set-safety-limit \
  --limit 1000000 \
  --from <your-key> \
  --chain-id natchain \
  --keyring-backend test \
  -y
```

### set-cold-address

设置冷钱包地址

```bash
natchaind tx wallet set-cold-address \
  --cold-address <冷钱包地址> \
  --from <your-key> \
  --chain-id natchain \
  --keyring-backend test \
  -y
```

### set-auto-transfer

设置自动转账

```bash
natchaind tx wallet set-auto-transfer \
  --enabled true \
  --amount 500000 \
  --from <your-key> \
  --chain-id natchain \
  --keyring-backend test \
  -y
```

### add-partner

添加合伙人

```bash
natchaind tx wallet add-partner \
  --partner-address <地址> \
  --percentage 10 \
  --from <your-key> \
  --chain-id natchain \
  --keyring-backend test \
  -y
```

### update-partner

更新合伙人分成比例

```bash
natchaind tx wallet update-partner \
  --partner-address <地址> \
  --percentage 15 \
  --from <your-key> \
  --chain-id natchain \
  --keyring-backend test \
  -y
```

### remove-partner

移除合伙人

```bash
natchaind tx wallet remove-partner \
  --partner-address <地址> \
  --from <your-key> \
  --chain-id natchain \
  --keyring-backend test \
  -y
```

---

## 查询命令

### wallet

查询钱包信息

```bash
natchaind q wallet wallet --address <地址>
```

返回：

```json
{
  "address": "natchain1...",
  "cold_address": "natchain1...",
  "safety_limit": "1000000",
  "auto_transfer_enabled": true,
  "auto_transfer_amount": "500000"
}
```

### balance

查询余额和阈值信息

```bash
natchaind q wallet balance --address <地址>
```

### partner

查询合伙人信息

```bash
natchaind q wallet partner --address <钱包地址> --partner-address <合伙人地址>
```

### partners

查询所有合伙人

```bash
natchaind q wallet partners --address <钱包地址>
```

### transfer-records

查询转账记录

```bash
natchaind q wallet transfer-records --address <地址>
```

---

## 数据结构

### Wallet

| 字段 | 类型 | 说明 |
|------|------|------|
| address | string | 钱包地址 |
| cold_address | string | 冷钱包地址 |
| safety_limit | int64 | 安全额度 |
| auto_transfer_enabled | bool | 是否启用自动转账 |
| auto_transfer_amount | int64 | 自动转账阈值 |

### Partner

| 字段 | 类型 | 说明 |
|------|------|------|
| address | string | 合伙人地址 |
| percentage | int64 | 分成比例 (%) |

---

## 状态码

| Code | 说明 |
|------|------|
| 1 | 钱包不存在 |
| 2 | 合伙人不存在 |
| 3 | 合伙人已存在 |
| 4 | 分成比例超过100% |
