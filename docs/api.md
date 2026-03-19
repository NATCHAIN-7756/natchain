# NatChain API 文档

## API 概览

NatChain 提供三种 API 接口：

| 类型 | 用途 | 协议 |
|------|------|------|
| **REST API** | 商家对接、查询 | HTTP |
| **gRPC API** | 高性能应用 | gRPC |
| **CLI** | 命令行操作 | 终端 |

## 基础信息

| 项目 | 值 |
|------|-----|
| REST Base URL | http://localhost:1317 |
| gRPC Endpoint | localhost:9090 |
| RPC Endpoint | localhost:26657 |

---

## REST API

### 通用说明

#### 请求头

```
Content-Type: application/json
```

#### 商家 API 额外请求头

```
X-Api-Key: {your_api_key}
X-Signature: {hmac_signature}
X-Timestamp: {unix_timestamp}
```

#### 签名算法

```javascript
// 构造签名字符串
const signStr = timestamp + method + path + body;

// HMAC-SHA256 签名
const signature = CryptoJS.HmacSHA256(signStr, apiSecret).toString();
```

---

## 矿机相关 API

### 查询矿机信息

```
GET /natchain/miner/v1/miners/{index}

Response:
{
  "miner": {
    "index": 1,
    "cpu_serial": "CPU-XXX-123",
    "cpu_model": "Intel i5-12400",
    "owner": "natchain1xxx",
    "validator_address": "natchainvaloper1xxx",
    "bonded_at": "1234567890",
    "last_active": "1234567890",
    "status": "MINER_STATUS_BONDED"
  }
}
```

### 查询所有矿机

```
GET /natchain/miner/v1/miners

Response:
{
  "miners": [...],
  "pagination": {
    "next_key": null,
    "total": "100"
  }
}
```

### 按所有者查询

```
GET /natchain/miner/v1/miners/owner/{address}

Response:
{
  "miners": [...]
}
```

---

## 钱包安全 API

### 查询安全设置

```
GET /natchain/wallet/v1/security/{address}

Response:
{
  "security": {
    "owner": "natchain1xxx",
    "daily_limit": { "denom": "unat", "amount": "1000000000" },
    "single_limit": { "denom": "unat", "amount": "100000000" },
    "threshold": { "denom": "unat", "amount": "5000000000" },
    "auto_transfer_to": "natchain1yyy",
    "enabled": true
  }
}
```

### 查询伙伴列表

```
GET /natchain/wallet/v1/partners/{address}

Response:
{
  "partners": [
    {
      "name": "Partner A",
      "address": "natchain1zzz",
      "memo": "Memo",
      "created_at": "1234567890"
    }
  ]
}
```

---

## 商家 API

### 申请 API Key

```
POST /api/v1/apikey/apply

Request:
{
  "callback_url": "https://example.com/callback",
  "daily_limit": { "denom": "unat", "amount": "1000000000" }
}

Response:
{
  "api_key_id": "ak_xxx",
  "api_key": "pk_xxx",
  "api_secret": "sk_xxx",
  "created_at": 1234567890
}
```

### 查询 API Key

```
GET /api/v1/apikey/info

Response:
{
  "api_key_id": "ak_xxx",
  "status": "ACTIVE",
  "created_at": 1234567890,
  "callback_url": "https://example.com/callback"
}
```

### 创建收款订单

```
POST /api/v1/payment/create

Request:
{
  "amount": { "denom": "unat", "amount": "1000000" },
  "memo": "Order #12345",
  "expire_in": 3600
}

Response:
{
  "order_id": "ord_xxx",
  "payment_address": "natchain1xxx",
  "amount": { "denom": "unat", "amount": "1000000" },
  "memo": "Order #12345",
  "expired_at": 1234567890,
  "qr_code": "natchain:pay?address=xxx&amount=1000000&memo=Order%20%2312345"
}
```

### 查询订单

```
GET /api/v1/payment/query/{order_id}

Response:
{
  "order_id": "ord_xxx",
  "status": "PAID",
  "amount": { "denom": "unat", "amount": "1000000" },
  "memo": "Order #12345",
  "tx_hash": "0x...",
  "created_at": 1234567890,
  "paid_at": 1234568000
}
```

### 查询余额

```
GET /api/v1/balance

Response:
{
  "address": "natchain1xxx",
  "balances": [
    { "denom": "unat", "amount": "100000000" }
  ]
}
```

### 转账

```
POST /api/v1/transfer

Request:
{
  "to": "natchain1yyy",
  "amount": { "denom": "unat", "amount": "1000000" },
  "memo": "Payment"
}

Response:
{
  "tx_hash": "0x...",
  "status": "SUCCESS",
  "height": 12345
}
```

---

## 回调通知

### 支付成功回调

```
POST {callback_url}

Request:
{
  "event": "payment.paid",
  "data": {
    "order_id": "ord_xxx",
    "amount": { "denom": "unat", "amount": "1000000" },
    "memo": "Order #12345",
    "tx_hash": "0x...",
    "paid_at": 1234567890
  }
}

Expected Response:
{
  "code": 0,
  "message": "success"
}
```

---

## 错误码

| Code | 说明 |
|------|------|
| 0 | 成功 |
| 1001 | API Key 无效 |
| 1002 | 签名验证失败 |
| 1003 | 请求过期 |
| 1004 | API 调用限额超限 |
| 2001 | 订单不存在 |
| 2002 | 订单已过期 |
| 2003 | 订单已支付 |
| 3001 | 余额不足 |
| 3002 | 转账超限 |
| 4001 | 矿机不存在 |
| 4002 | 矿机已绑定 |
| 4003 | 矿机数量超限 |

---

## SDK 示例

### JavaScript

```javascript
const axios = require('axios');
const CryptoJS = require('crypto-js');

class NatChainClient {
  constructor(apiKey, apiSecret, baseUrl = 'http://localhost:1317') {
    this.apiKey = apiKey;
    this.apiSecret = apiSecret;
    this.baseUrl = baseUrl;
  }

  sign(method, path, body) {
    const timestamp = Math.floor(Date.now() / 1000).toString();
    const signStr = timestamp + method + path + JSON.stringify(body);
    const signature = CryptoJS.HmacSHA256(signStr, this.apiSecret).toString();
    return { timestamp, signature };
  }

  async request(method, path, body = {}) {
    const { timestamp, signature } = this.sign(method, path, body);
    const res = await axios({
      method,
      url: this.baseUrl + path,
      data: body,
      headers: {
        'X-Api-Key': this.apiKey,
        'X-Signature': signature,
        'X-Timestamp': timestamp,
        'Content-Type': 'application/json'
      }
    });
    return res.data;
  }

  async createPayment(amount, memo) {
    return this.request('POST', '/api/v1/payment/create', { amount, memo });
  }

  async queryOrder(orderId) {
    return this.request('GET', `/api/v1/payment/query/${orderId}`);
  }
}

// 使用示例
const client = new NatChainClient('pk_xxx', 'sk_xxx');
const order = await client.createPayment(
  { denom: 'unat', amount: '1000000' },
  'Order #12345'
);
console.log(order);
```

### Go

```go
package main

import (
    "bytes"
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "strconv"
    "time"
)

type NatChainClient struct {
    APIKey    string
    APISecret string
    BaseURL   string
}

func (c *NatChainClient) Sign(method, path, body string) (timestamp, signature string) {
    timestamp = strconv.FormatInt(time.Now().Unix(), 10)
    signStr := timestamp + method + path + body
    
    h := hmac.New(sha256.New, []byte(c.APISecret))
    h.Write([]byte(signStr))
    signature = hex.EncodeToString(h.Sum(nil))
    return
}

func (c *NatChainClient) CreatePayment(amount, memo string) (map[string]interface{}, error) {
    body := map[string]interface{}{
        "amount": map[string]string{"denom": "unat", "amount": amount},
        "memo":   memo,
    }
    bodyBytes, _ := json.Marshal(body)
    
    timestamp, signature := c.Sign("POST", "/api/v1/payment/create", string(bodyBytes))
    
    req, _ := http.NewRequest("POST", c.BaseURL+"/api/v1/payment/create", bytes.NewBuffer(bodyBytes))
    req.Header.Set("X-Api-Key", c.APIKey)
    req.Header.Set("X-Signature", signature)
    req.Header.Set("X-Timestamp", timestamp)
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var result map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&result)
    return result, nil
}
```
