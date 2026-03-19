# x/merchant 模块设计

## 模块概述

商家 API 模块，提供 API Key 管理、支付接口、回调通知功能。

## 核心功能

| 功能 | 说明 |
|------|------|
| API Key 管理 | 申请、查询、撤销 |
| 支付创建 | 创建收款订单 |
| 支付验证 | 验证交易是否到账 |
| 回调通知 | 支付成功后通知商家 |
| 转账接口 | 商家主动转账 |

## 数据结构

### MerchantApiKey

```go
// proto/natchain/merchant/merchant.proto

message MerchantApiKey {
  string id = 1;                 // API Key ID
  string owner = 2;              // 钱包地址
  string api_key = 3;            // API Key (公开)
  string api_secret = 4;         // API Secret (保密)
  string callback_url = 5;       // 回调地址
  CosmosCoin daily_limit = 6;    // API 调用限额
  int64 created_at = 7;          // 创建时间
  int64 expires_at = 8;          // 过期时间
  ApiKeyStatus status = 9;       // 状态
}

enum ApiKeyStatus {
  API_KEY_STATUS_UNSPECIFIED = 0;
  API_KEY_STATUS_ACTIVE = 1;
  API_KEY_STATUS_REVOKED = 2;
  API_KEY_STATUS_EXPIRED = 3;
}
```

### PaymentOrder

```go
message PaymentOrder {
  string order_id = 1;           // 订单 ID
  string merchant_address = 2;   // 商家地址
  CosmosCoin amount = 3;         // 金额
  string memo = 4;               // 备注
  string payer_address = 5;      // 付款人地址（可选）
  int64 created_at = 6;          // 创建时间
  int64 expired_at = 7;          // 过期时间
  PaymentStatus status = 8;      // 状态
  string tx_hash = 9;            // 交易哈希
}

enum PaymentStatus {
  PAYMENT_STATUS_UNSPECIFIED = 0;
  PAYMENT_STATUS_PENDING = 1;
  PAYMENT_STATUS_PAID = 2;
  PAYMENT_STATUS_EXPIRED = 3;
  PAYMENT_STATUS_CANCELLED = 4;
}
```

### ApiUsage（API 调用统计）

```go
message ApiUsage {
  string api_key_id = 1;
  int64 date = 2;                // YYYYMMDD
  uint64 call_count = 3;         // 调用次数
}
```

### GenesisState

```go
message GenesisState {
  repeated MerchantApiKey api_keys = 1;
  repeated PaymentOrder orders = 2;
  Params params = 3;
}

message Params {
  int64 order_expiry = 1;        // 订单过期时间（秒）
  int64 max_api_keys = 2;        // 最大 API Key 数量
}
```

## Keeper 接口

```go
// keeper/merchant.go

type Keeper interface {
  // API Key 管理
  CreateApiKey(ctx sdk.Context, owner, callbackUrl string) (MerchantApiKey, error)
  GetApiKey(ctx sdk.Context, id string) (MerchantApiKey, bool)
  GetApiKeysByOwner(ctx sdk.Context, owner string) []MerchantApiKey
  RevokeApiKey(ctx sdk.Context, id, owner string) error
  
  // 订单管理
  CreateOrder(ctx sdk.Context, merchant, amount, memo string) (PaymentOrder, error)
  GetOrder(ctx sdk.Context, orderId string) (PaymentOrder, bool)
  GetOrdersByMerchant(ctx sdk.Context, merchant string) []PaymentOrder
  UpdateOrderStatus(ctx sdk.Context, orderId string, status PaymentStatus, txHash string) error
  
  // API 用量
  RecordApiUsage(ctx sdk.Context, apiKeyId string) error
  CheckApiLimit(ctx sdk.Context, apiKeyId string) error
  
  // 验证
  ValidateSignature(ctx sdk.Context, apiKey, signature, payload string) bool
}
```

## Msg 服务

```go
// proto/natchain/merchant/tx.proto

service Msg {
  // API Key
  rpc CreateApiKey(MsgCreateApiKey) returns (MsgCreateApiKeyResponse);
  rpc RevokeApiKey(MsgRevokeApiKey) returns (MsgRevokeApiKeyResponse);
  
  // 订单
  rpc CreateOrder(MsgCreateOrder) returns (MsgCreateOrderResponse);
  rpc CancelOrder(MsgCancelOrder) returns (MsgCancelOrderResponse);
}

message MsgCreateApiKey {
  string sender = 1;
  string callback_url = 2;
  CosmosCoin daily_limit = 3;
}

message MsgRevokeApiKey {
  string sender = 1;
  string api_key_id = 2;
}

message MsgCreateOrder {
  string sender = 1;
  CosmosCoin amount = 2;
  string memo = 3;
}

message MsgCancelOrder {
  string sender = 1;
  string order_id = 2;
}
```

## REST API

### 1. 申请 API Key

```
POST /api/v1/apikey/apply

Headers:
  X-Api-Key: {api_key}
  X-Signature: {hmac_sha256(api_secret, payload)}
  X-Timestamp: {unix_timestamp}

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

### 2. 创建收款订单

```
POST /api/v1/payment/create

Headers:
  X-Api-Key: {api_key}
  X-Signature: {signature}
  X-Timestamp: {timestamp}

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
  "qr_code": "natchain:xxx"
}
```

### 3. 查询订单

```
GET /api/v1/payment/query/{order_id}

Headers:
  X-Api-Key: {api_key}
  X-Signature: {signature}
  X-Timestamp: {timestamp}

Response:
{
  "order_id": "ord_xxx",
  "status": "PAID",
  "amount": { "denom": "unat", "amount": "1000000" },
  "tx_hash": "0x...",
  "paid_at": 1234567890
}
```

### 4. 查询余额

```
GET /api/v1/balance

Headers:
  X-Api-Key: {api_key}
  X-Signature: {signature}

Response:
{
  "address": "natchain1xxx",
  "balances": [
    { "denom": "unat", "amount": "100000000" }
  ]
}
```

### 5. 转账

```
POST /api/v1/transfer

Headers:
  X-Api-Key: {api_key}
  X-Signature: {signature}

Request:
{
  "to": "natchain1yyy",
  "amount": { "denom": "unat", "amount": "1000000" },
  "memo": "Payment"
}

Response:
{
  "tx_hash": "0x...",
  "status": "SUCCESS"
}
```

## 回调通知

### 支付成功回调

```
POST {callback_url}

Headers:
  X-Signature: {hmac_sha256(api_secret, payload)}
  X-Timestamp: {timestamp}

Request:
{
  "event": "payment.paid",
  "data": {
    "order_id": "ord_xxx",
    "amount": { "denom": "unat", "amount": "1000000" },
    "tx_hash": "0x...",
    "paid_at": 1234567890
  }
}
```

### 商家响应

```
Response:
{
  "code": 0,
  "message": "success"
}
```

## 签名算法

```
1. 构造待签名字符串
   sign_str = timestamp + method + path + body

2. HMAC-SHA256 签名
   signature = hmac_sha256(api_secret, sign_str)

3. 十六进制编码
   signature_hex = hex(signature)

示例:
   timestamp = "1234567890"
   method = "POST"
   path = "/api/v1/payment/create"
   body = '{"amount":{"denom":"unat","amount":"1000000"}}'
   
   sign_str = "1234567890POST/api/v1/payment/create{\"amount\":{\"denom\":\"unat\",\"amount\":\"1000000\"}}"
   signature = hmac_sha256("sk_xxx", sign_str)
```

## 业务逻辑

### 1. 申请 API Key

```
1. 验证钱包地址有效
2. 检查 API Key 数量限制
3. 生成：
   - id: ak_{random}
   - api_key: pk_{random} (公开)
   - api_secret: sk_{random} (保密)
4. 保存记录
5. 返回 api_key 和 api_secret
```

### 2. 创建订单

```
1. 验证 API Key 有效
2. 验证签名
3. 检查 API 调用限额
4. 生成订单：
   - order_id: ord_{random}
   - payment_address: 商家钱包地址
   - expired_at: now + 1 hour
5. 保存订单
6. 返回订单信息和支付二维码
```

### 3. 支付验证

```
监听链上交易：
1. 过滤到 payment_address 的交易
2. 匹配 memo 找到对应订单
3. 验证金额匹配
4. 更新订单状态
5. 触发回调通知
```

### 4. 回调重试

```
1. 收到支付事件后发送回调
2. 如果失败，重试策略：
   - 第1次: 立即
   - 第2次: 1分钟后
   - 第3次: 5分钟后
   - 第4次: 30分钟后
   - 第5次: 1小时后
3. 超过5次失败，标记为回调失败
```

## 事件

```go
message EventApiKeyCreated {
  string owner = 1;
  string api_key_id = 2;
}

message EventApiKeyRevoked {
  string owner = 1;
  string api_key_id = 2;
}

message EventOrderCreated {
  string order_id = 1;
  string merchant = 2;
  CosmosCoin amount = 3;
}

message EventOrderPaid {
  string order_id = 1;
  string tx_hash = 2;
}

message EventCallbackFailed {
  string order_id = 1;
  string callback_url = 2;
  int64 retry_count = 3;
}
```

## 单元测试清单

- [ ] TestCreateApiKey
- [ ] TestCreateApiKey_MaxLimit
- [ ] TestRevokeApiKey
- [ ] TestCreateOrder
- [ ] TestCancelOrder
- [ ] TestValidateSignature
- [ ] TestCheckApiLimit
- [ ] TestPaymentVerification
- [ ] TestCallbackNotification
- [ ] TestCallbackRetry
