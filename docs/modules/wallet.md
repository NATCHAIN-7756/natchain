# x/wallet 模块设计

## 模块概述

钱包安全模块，提供安全额度、自动转账、伙伴管理功能。

## 核心功能

| 功能 | 说明 |
|------|------|
| 安全额度 | 设置单笔/每日转出上限 |
| 自动转账 | 余额超过阈值自动转冷钱包 |
| 伙伴管理 | 常用地址管理 |
| 转账拦截 | 检查额度，保护资金安全 |

## 数据结构

### WalletSecurity（钱包安全设置）

```go
// proto/natchain/wallet/wallet.proto

message WalletSecurity {
  string owner = 1;              // 钱包地址
  CosmosCoin daily_limit = 2;    // 每日转出上限
  CosmosCoin single_limit = 3;   // 单笔转出上限
  CosmosCoin threshold = 4;      // 自动转账阈值
  string auto_transfer_to = 5;   // 冷钱包地址
  CosmosCoin auto_transfer_amount = 6; // 自动转账金额
  bool enabled = 7;              // 是否启用
}

message CosmosCoin {
  string denom = 1;
  string amount = 2;
}
```

### Partner（伙伴）

```go
message Partner {
  string owner = 1;       // 用户地址
  string name = 2;        // 伙伴名称
  string address = 3;     // 伙伴地址
  string memo = 4;        // 备注
  int64 created_at = 5;   // 添加时间
}
```

### DailyTransfer（每日转账记录）

```go
message DailyTransfer {
  string owner = 1;           // 用户地址
  int64 date = 2;             // 日期（YYYYMMDD）
  CosmosCoin transferred = 3; // 已转出金额
}
```

### GenesisState

```go
message GenesisState {
  repeated WalletSecurity wallet_securities = 1;
  repeated Partner partners = 2;
  repeated DailyTransfer daily_transfers = 3;
  Params params = 4;
}

message Params {
  int64 max_partners = 1;      // 最大伙伴数量
}
```

## Keeper 接口

```go
// keeper/wallet.go

type Keeper interface {
  // 安全设置
  SetSecurity(ctx sdk.Context, security WalletSecurity) error
  GetSecurity(ctx sdk.Context, owner string) (WalletSecurity, bool)
  DeleteSecurity(ctx sdk.Context, owner string) error
  
  // 伙伴管理
  AddPartner(ctx sdk.Context, partner Partner) error
  GetPartner(ctx sdk.Context, owner, address string) (Partner, bool)
  GetPartnersByOwner(ctx sdk.Context, owner string) []Partner
  RemovePartner(ctx sdk.Context, owner, address string) error
  
  // 每日额度
  GetDailyTransfer(ctx sdk.Context, owner string) (CosmosCoin, error)
  AddDailyTransfer(ctx sdk.Context, owner string, amount CosmosCoin) error
  ResetDailyTransfer(ctx sdk.Context, owner string) error
  
  // 检查
  CheckTransferLimit(ctx sdk.Context, owner string, amount CosmosCoin) error
  CheckSingleLimit(ctx sdk.Context, owner string, amount CosmosCoin) error
}
```

## Msg 服务

```go
// proto/natchain/wallet/tx.proto

service Msg {
  // 安全设置
  rpc SetSecurity(MsgSetSecurity) returns (MsgSetSecurityResponse);
  rpc DeleteSecurity(MsgDeleteSecurity) returns (MsgDeleteSecurityResponse);
  
  // 伙伴管理
  rpc AddPartner(MsgAddPartner) returns (MsgAddPartnerResponse);
  rpc RemovePartner(MsgRemovePartner) returns (MsgRemovePartnerResponse);
}

message MsgSetSecurity {
  string sender = 1;
  CosmosCoin daily_limit = 2;
  CosmosCoin single_limit = 3;
  CosmosCoin threshold = 4;
  string auto_transfer_to = 5;
  CosmosCoin auto_transfer_amount = 6;
  bool enabled = 7;
}

message MsgAddPartner {
  string sender = 1;
  string name = 2;
  string address = 3;
  string memo = 4;
}

message MsgRemovePartner {
  string sender = 1;
  string address = 2;
}
```

## Query 服务

```go
// proto/natchain/wallet/query.proto

service Query {
  rpc Security(QuerySecurityRequest) returns (QuerySecurityResponse);
  rpc Partner(QueryPartnerRequest) returns (QueryPartnerResponse);
  rpc Partners(QueryPartnersRequest) returns (QueryPartnersResponse);
  rpc DailyTransfer(QueryDailyTransferRequest) returns (QueryDailyTransferResponse);
}
```

## 业务逻辑

### 1. 设置安全额度

```
1. 验证调用者是钱包所有者
2. 验证参数合法性：
   - daily_limit > 0
   - single_limit > 0
   - threshold > 0
   - auto_transfer_to 是有效地址
3. 保存设置
4. 发送事件
```

### 2. 转账拦截（Hook）

```
在 Bank 模块转账前：
1. 检查发送者是否有安全设置
2. 如果启用：
   a. 检查单笔限额
   b. 检查每日限额
   c. 如果超过，拒绝转账
3. 更新每日转账记录
4. 允许转账
```

### 3. 自动转账触发

```
每个区块结束：
1. 遍历所有启用了安全设置的钱包
2. 查询余额
3. 如果余额 > 阈值：
   a. 计算超出金额
   b. 发起转账到冷钱包
   c. 发送事件
```

### 4. 添加伙伴

```
1. 验证调用者
2. 检查伙伴数量是否超限
3. 检查地址是否已存在
4. 添加伙伴记录
5. 发送事件
```

## 事件

```go
message EventSecuritySet {
  string owner = 1;
}

message EventTransferBlocked {
  string owner = 1;
  string reason = 2;
  CosmosCoin attempted = 3;
}

message EventAutoTransfer {
  string owner = 1;
  string to = 2;
  CosmosCoin amount = 3;
}

message EventPartnerAdded {
  string owner = 1;
  string partner_address = 2;
  string name = 3;
}
```

## 与 Bank 模块集成

```go
// 在 x/bank/keeper/keeper.go 中添加钩子

type BankKeeper interface {
  // 原有方法...
  
  // 添加安全检查
  SetWalletKeeper(wk WalletKeeper)
}

// 转账前检查
func (k BankKeeper) SendCoins(ctx sdk.Context, from, to sdk.AccAddress, amt sdk.Coins) error {
  // 调用 wallet 模块检查额度
  if err := k.walletKeeper.CheckTransferLimit(ctx, from.String(), amt); err != nil {
    return err
  }
  
  // 原有转账逻辑...
}
```

## 单元测试清单

- [ ] TestSetSecurity
- [ ] TestSetSecurity_InvalidParams
- [ ] TestDeleteSecurity
- [ ] TestAddPartner
- [ ] TestAddPartner_MaxLimit
- [ ] TestRemovePartner
- [ ] TestCheckTransferLimit
- [ ] TestCheckTransferLimit_Exceeded
- [ ] TestCheckSingleLimit
- [ ] TestAutoTransfer
- [ ] TestDailyTransferReset
