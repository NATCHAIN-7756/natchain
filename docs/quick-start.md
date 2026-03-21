# NatChain 快速开始

## 1. 环境要求

- Go 1.21+
- Linux/macOS
- 4GB+ RAM

---

## 2. 安装

### 编译源码

```bash
git clone https://github.com/NATCHAIN-7756/natchain.git
cd natchain
make build
```

### 安装二进制

```bash
cp build/natchaind /usr/local/bin/
```

---

## 3. 初始化节点

```bash
# 初始化
natchaind init <your-moniker> --chain-id natchain

# 创建密钥
natchaind keys add validator --keyring-backend test
```

---

## 4. 获取初始代币

联系团队获取初始代币或从测试网水龙头获取。

---

## 5. 启动节点

### 前台运行（测试）

```bash
natchaind start
```

### 后台服务（生产）

```bash
# 创建服务文件
sudo tee /etc/systemd/system/natchaind.service << EOF
[Unit]
Description=NatChain Daemon
After=network.target

[Service]
Type=simple
User=$USER
ExecStart=/usr/local/bin/natchaind start
Restart=always
RestartSec=5
LimitNOFILE=65535

[Install]
WantedBy=multi-user.target
EOF

# 启动服务
sudo systemctl daemon-reload
sudo systemctl enable natchaind
sudo systemctl start natchaind
```

---

## 6. 连接到测试网

### 添加种子节点

编辑 `~/.natchain/config/config.toml`:

```toml
[seeds]
seeds = "节点ID@38.76.200.118:26656"
```

### 或使用 persistent_peers

```toml
[p2p]
persistent_peers = "节点ID@38.76.200.118:26656"
```

---

## 7. 验证连接

```bash
# 查看节点状态
natchaind status

# 查看区块高度
curl -s localhost:26657/status | jq '.result.sync_info.latest_block_height'
```

---

## 8. 测试交易

```bash
# 查询余额
natchaind q bank balances $(natchaind keys show validator -a --keyring-backend test)

# 发送交易
natchaind tx bank send validator <recipient> 100stake \
  --chain-id natchain \
  --keyring-backend test \
  -y
```

---

## 9. 成为验证者（可选）

```bash
natchaind tx staking create-validator \
  --amount 1000000stake \
  --pubkey $(natchaind tendermint show-validator) \
  --moniker "your-validator-name" \
  --chain-id natchain \
  --keyring-backend test \
  --from validator \
  -y
```

---

## 常用命令速查

| 操作 | 命令 |
|------|------|
| 查看状态 | `natchaind status` |
| 查看余额 | `natchaind q bank balances <addr>` |
| 查看密钥 | `natchaind keys list --keyring-backend test` |
| 查看日志 | `journalctl -u natchaind -f` |
| 重启服务 | `systemctl restart natchaind` |

---

## 获取帮助

- GitHub: https://github.com/NATCHAIN-7756/natchain
- 官网: natjen.com
