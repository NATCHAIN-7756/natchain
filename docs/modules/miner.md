# x/miner 模块设计

## 模块概述

矿机管理模块，负责矿机的注册、绑定、状态追踪。

## 核心功能

| 功能 | 说明 |
|------|------|
| 矿机注册 | 管理员创建矿机记录 |
| 矿机绑定 | 用户购买并绑定矿机到钱包 |
| 状态查询 | 查询矿机信息和状态 |
| 活跃追踪 | 记录矿机最后出块时间 |
| 离线惩罚 | 超过 180 天未活跃自动解绑 |

## 数据结构

### Miner（矿机）

```go
// proto/natchain/miner/miner.proto

message Miner {
  uint64 index = 1;           // 矿机编号 1-1000
  string cpu_serial = 2;      // CPU 序列号
  string cpu_model = 3;       // CPU 型号
  string owner = 4;           // 所有者地址
  string validator_address = 5; // 验证者地址
  int64 bonded_at = 6;        // 绑定时间戳
  int64 last_active = 7;      // 最后活跃时间
  MinerStatus status = 8;     // 状态
  bool paid = 9;              // 是否已付款
}

enum MinerStatus {
  MINER_STATUS_UNSPECIFIED = 0;
  MINER_STATUS_PENDING = 1;     // 待绑定
  MINER_STATUS_BONDED = 2;      // 已绑定
  MINER_STATUS_UNBONDED = 3;    // 已解绑
}
```

### GenesisState

```go
message GenesisState {
  repeated Miner miners = 1;
  uint64 next_index = 2;       // 下一个可用编号
  Params params = 3;
}

message Params {
  uint64 max_miners = 1;       // 最大矿机数 1000
  string price = 2;            // 价格 100 USDT
  int64 unbonding_time = 3;    // 解绑时间 180天
}
```

## Keeper 接口

```go
// keeper/miner.go

type Keeper interface {
  // 矿机管理
  CreateMiner(ctx sdk.Context, cpuSerial, cpuModel string) (uint64, error)
  GetMiner(ctx sdk.Context, index uint64) (Miner, bool)
  GetMinerByOwner(ctx sdk.Context, owner string) ([]Miner, bool)
  GetMinerByCPU(ctx sdk.Context, cpuSerial string) (Miner, bool)
  
  // 绑定操作
  BondMiner(ctx sdk.Context, index uint64, owner, validatorAddr string) error
  UnbondMiner(ctx sdk.Context, index uint64) error
  
  // 状态更新
  UpdateLastActive(ctx sdk.Context, index uint64) error
  
  // 查询
  GetAllMiners(ctx sdk.Context) []Miner
  GetBondedMiners(ctx sdk.Context) []Miner
  GetPendingMiners(ctx sdk.Context) []Miner
}
```

## Msg 服务

```go
// proto/natchain/miner/tx.proto

service Msg {
  // 创建矿机（管理员）
  rpc CreateMiner(MsgCreateMiner) returns (MsgCreateMinerResponse);
  
  // 绑定矿机（用户）
  rpc BondMiner(MsgBondMiner) returns (MsgBondMinerResponse);
  
  // 解绑矿机
  rpc UnbondMiner(MsgUnbondMiner) returns (MsgUnbondMinerResponse);
}

message MsgCreateMiner {
  string creator = 1;
  string cpu_serial = 2;
  string cpu_model = 3;
}

message MsgBondMiner {
  string sender = 1;
  uint64 index = 2;
  string validator_address = 3;
}

message MsgUnbondMiner {
  string sender = 1;
  uint64 index = 2;
}
```

## Query 服务

```go
// proto/natchain/miner/query.proto

service Query {
  rpc Miner(QueryMinerRequest) returns (QueryMinerResponse);
  rpc Miners(QueryMinersRequest) returns (QueryMinersResponse);
  rpc Params(QueryParamsRequest) returns (QueryParamsResponse);
  rpc MinerByOwner(QueryMinerByOwnerRequest) returns (QueryMinerByOwnerResponse);
}
```

## 业务逻辑

### 1. 矿机创建（管理员）

```
1. 验证调用者是管理员
2. 检查矿机数量是否已达上限（1000）
3. 生成矿机编号
4. 验证 CPU 序列号唯一性
5. 创建矿机记录（状态：PENDING）
6. 发送事件
```

### 2. 矿机绑定（用户）

```
1. 验证矿机存在且状态为 PENDING
2. 验证用户已付款（集成支付模块）
3. 更新矿机信息：
   - owner = 用户地址
   - validator_address = 验证者地址
   - bonded_at = 当前时间
   - status = BONDED
4. 创建验证者（调用 staking 模块）
5. 发送事件
```

### 3. 活跃更新

```
每个区块：
1. 遍历所有已绑定矿机
2. 检查是否是出块验证者
3. 如果是，更新 last_active = 当前时间
```

### 4. 离线惩罚

```
每个区块：
1. 遍历所有已绑定矿机
2. 检查 last_active 是否超过 180 天
3. 如果超过，自动解绑：
   - status = UNBONDED
   - 清除 validator
4. 发送事件
```

## 事件

```go
message EventMinerCreated {
  uint64 index = 1;
  string cpu_serial = 2;
}

message EventMinerBonded {
  uint64 index = 1;
  string owner = 2;
}

message EventMinerUnbonded {
  uint64 index = 1;
  string owner = 2;
  string reason = 3;
}
```

## 单元测试清单

- [ ] TestCreateMiner
- [ ] TestCreateMiner_MaxLimit
- [ ] TestCreateMiner_DuplicateCPU
- [ ] TestBondMiner
- [ ] TestBondMiner_AlreadyBonded
- [ ] TestBondMiner_NotPaid
- [ ] TestUnbondMiner
- [ ] TestUnbondMiner_Unauthorized
- [ ] TestUpdateLastActive
- [ ] TestOfflineUnbonding
