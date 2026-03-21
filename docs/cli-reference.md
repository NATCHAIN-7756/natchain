# NatChain CLI 命令手册

## 基础命令

### 查看版本

```bash
natchaind version
```

### 查看状态

```bash
natchaind status
```

---

## 密钥管理

### 创建密钥

```bash
natchaind keys add <name> --keyring-backend test
```

### 列出密钥

```bash
natchaind keys list --keyring-backend test
```

### 删除密钥

```bash
natchaind keys delete <name> --keyring-backend test
```

---

## 银行模块

### 查询余额

```bash
natchaind q bank balances <address>
```

### 转账

```bash
natchaind tx bank send <from> <to> <amount> \
  --chain-id natchain \
  --keyring-backend test \
  -y
```

---

## 矿机模块 (x/miner)

### 注册矿机

```bash
natchaind tx miner register-miner \
  --cpu-sn <CPU序列号> \
  --from <your-key> \
  --chain-id natchain \
  --keyring-backend test \
  -y
```

### 矿机心跳

```bash
natchaind tx miner heartbeat \
  --cpu-sn <CPU序列号> \
  --from <your-key> \
  --chain-id natchain \
  --keyring-backend test \
  -y
```

### 查询参数

```bash
natchaind q miner params
```

---

## 钱包模块 (x/wallet)

### 查询钱包

```bash
natchaind q wallet wallet --address <address>
```

### 设置安全额度

```bash
natchaind tx wallet set-safety-limit \
  --limit <金额> \
  --from <your-key> \
  --chain-id natchain \
  --keyring-backend test \
  -y
```

### 设置冷钱包地址

```bash
natchaind tx wallet set-cold-address \
  --cold-address <冷钱包地址> \
  --from <your-key> \
  --chain-id natchain \
  --keyring-backend test \
  -y
```

### 设置自动转账

```bash
natchaind tx wallet set-auto-transfer \
  --enabled true \
  --amount <转账阈值> \
  --from <your-key> \
  --chain-id natchain \
  --keyring-backend test \
  -y
```

### 添加合伙人

```bash
natchaind tx wallet add-partner \
  --partner-address <地址> \
  --percentage <分成比例> \
  --from <your-key> \
  --chain-id natchain \
  --keyring-backend test \
  -y
```

### 查询转账记录

```bash
natchaind q wallet transfer-records --address <address>
```

---

## 商家模块 (x/merchant)

### 注册商家

```bash
natchaind tx merchant register \
  --name <商家名称> \
  --callback-url <回调地址> \
  --commission-rate <佣金比例> \
  --from <your-key> \
  --chain-id natchain \
  --keyring-backend test \
  -y
```

### 查询商家

```bash
natchaind q merchant merchant --address <商家地址>
```

### 查询所有商家

```bash
natchaind q merchant merchants
```

### 创建支付

```bash
natchaind tx merchant create-payment \
  --merchant-address <商家地址> \
  --amount <金额> \
  --callback-data <自定义数据> \
  --from <user-key> \
  --chain-id natchain \
  --keyring-backend test \
  -y
```

### 查询支付记录

```bash
# 单笔支付
natchaind q merchant payment-record --id <支付ID>

# 商家所有支付
natchaind q merchant payment-records-by-merchant --merchant-address <地址>
```

### 确认支付

```bash
natchaind tx merchant confirm-payment \
  --payment-id <支付ID> \
  --from <merchant-key> \
  --chain-id natchain \
  --keyring-backend test \
  -y
```

### 重新生成 API Key

```bash
natchaind tx merchant regenerate-api-key \
  --from <merchant-key> \
  --chain-id natchain \
  --keyring-backend test \
  -y
```

### 更新商家信息

```bash
natchaind tx merchant update-merchant \
  --name <新名称> \
  --callback-url <新回调地址> \
  --from <merchant-key> \
  --chain-id natchain \
  --keyring-backend test \
  -y
```

---

## 查询命令

### 查询区块

```bash
natchaind q block <height>
```

### 查询交易

```bash
natchaind q tx <txhash>
```

### 查询节点状态

```bash
curl -s localhost:26657/status | jq
```

---

## 常用组合命令

### 检查链高度

```bash
natchaind status | jq '.sync_info.latest_block_height'
```

### 等待交易确认

```bash
natchaind q tx <txhash> --wait-tx
```

---

## 全局参数

| 参数 | 说明 |
|------|------|
| `--chain-id` | 链 ID (natchain) |
| `--home` | 数据目录 |
| `--keyring-backend` | 密钥存储后端 (test/os/file) |
| `--node` | RPC 节点地址 |
| `--output` | 输出格式 (json/text) |
| `-y` | 跳过确认 |

---

## 远程连接

```bash
natchaind q bank balances <address> \
  --node tcp://38.76.200.118:26657
```
