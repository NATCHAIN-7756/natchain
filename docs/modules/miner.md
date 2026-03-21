# x/miner 矿机模块

## 概述

管理矿机的注册、心跳和状态。

---

## 功能

- 矿机注册（绑定 CPU 序列号）
- 矿机心跳
- 离线检测（180天自动解绑）
- 矿机上限：1000 台

---

## 交易命令

### register-miner

注册新矿机

```bash
natchaind tx miner register-miner \
  --cpu-sn <CPU序列号> \
  --from <your-key> \
  --chain-id natchain \
  --keyring-backend test \
  -y
```

### heartbeat

发送矿机心跳

```bash
natchaind tx miner heartbeat \
  --cpu-sn <CPU序列号> \
  --from <your-key> \
  --chain-id natchain \
  --keyring-backend test \
  -y
```

---

## 查询命令

### params

查询模块参数

```bash
natchaind q miner params
```

---

## 状态码

| Code | 说明 |
|------|------|
| 1 | CPU SN 已存在 |
| 2 | 矿机数量已达上限 |
| 3 | 矿机不存在 |
| 4 | 非矿机所有者 |

---

## 事件

| 事件类型 | 属性 |
|----------|------|
| register_miner | creator, cpu_sn |
| heartbeat | cpu_sn, timestamp |
