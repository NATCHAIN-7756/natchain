# NatChain 运维文档

## 服务器信息

| 项目 | 值 |
|------|-----|
| IP | 38.76.200.118 |
| 用户 | root |
| SSH 工具 | `/home/node/clawd/ssh_tool.py` |

---

## 服务管理

### 查看服务状态

```bash
systemctl status natchaind
```

### 启动服务

```bash
systemctl start natchaind
```

### 停止服务

```bash
systemctl stop natchaind
```

### 重启服务

```bash
systemctl restart natchaind
```

### 查看日志

```bash
journalctl -u natchaind -f
```

### 查看最近日志

```bash
journalctl -u natchaind -n 100
```

---

## 健康检查

### 检查进程

```bash
pgrep -a natchaind
```

### 检查端口

```bash
netstat -tlnp | grep natchaind
```

### 检查区块高度

```bash
curl -s localhost:26657/status | jq '.result.sync_info'
```

### 检查是否同步完成

```bash
curl -s localhost:26657/status | jq '.result.sync_info.catching_up'
# false = 已同步
```

---

## 常见问题

### 1. 服务无法启动

**检查日志：**
```bash
journalctl -u natchaind -n 50
```

**常见原因：**
- 数据目录权限问题
- 端口被占用
- 配置文件损坏

### 2. 区块停止增长

**检查共识状态：**
```bash
curl -s localhost:26657/dump_consensus_state | jq
```

**可能原因：**
- 验证者密钥问题
- WAL 损坏

### 3. 内存不足

**检查内存：**
```bash
free -h
```

**解决方案：**
- 增加内存
- 调整 `config.toml` 中的缓存设置

---

## 数据备份

### 备份关键文件

```bash
# 备份验证者密钥
cp ~/.natchain/config/priv_validator_key.json /backup/

# 备份节点密钥
cp ~/.natchain/config/node_key.json /backup/

# 备份钱包
cp -r ~/.natchain/keyring-test /backup/
```

### 恢复验证者

```bash
# 停止服务
systemctl stop natchaind

# 恢复密钥
cp /backup/priv_validator_key.json ~/.natchain/config/

# 重启服务
systemctl start natchaind
```

---

## 链重置

### ⚠️ 警告：此操作会清除所有数据！

```bash
# 停止服务
systemctl stop natchaind

# 清除数据（保留配置）
rm -rf ~/.natchain/data

# 重新初始化
natchaind init <moniker> --chain-id natchain

# 重启服务
systemctl start natchaind
```

---

## 监控告警

### 设置监控脚本

```bash
cat > /usr/local/bin/natchain-monitor.sh << 'EOF'
#!/bin/bash
HEIGHT=$(curl -s localhost:26657/status | jq -r '.result.sync_info.latest_block_height')
PROCESS=$(pgrep natchaind)

if [ -z "$PROCESS" ]; then
  echo "CRITICAL: natchaind not running"
  exit 2
fi

if [ -z "$HEIGHT" ]; then
  echo "CRITICAL: Cannot get block height"
  exit 2
fi

echo "OK: Block height $HEIGHT"
exit 0
EOF
chmod +x /usr/local/bin/natchain-monitor.sh
```

### 添加 Cron 检查

```bash
# 每 5 分钟检查一次
*/5 * * * * /usr/local/bin/natchain-monitor.sh || systemctl restart natchaind
```

---

## 性能优化

### 调整文件描述符限制

```bash
# 查看当前限制
ulimit -n

# 修改 systemd 服务文件
# 在 [Service] 中添加：
LimitNOFILE=65535
```

### 调整内存使用

编辑 `~/.natchain/config/app.toml`:

```toml
[api]
enable = true

[grpc]
enable = true

[state-sync]
snapshot-interval = 1000
snapshot-keep-recent = 2
```

---

## 安全加固

### 防火墙设置

```bash
# 允许 P2P 端口
ufw allow 26656/tcp

# 允许 RPC 端口（仅内网）
ufw allow from 10.0.0.0/8 to any port 26657

# 启用防火墙
ufw enable
```

### 禁用 root SSH 密码登录

```bash
# 编辑 /etc/ssh/sshd_config
PermitRootLogin prohibit-password
PasswordAuthentication no
```

---

## 端口说明

| 端口 | 用途 | 访问 |
|------|------|------|
| 26656 | P2P | 公开 |
| 26657 | RPC | 内网 |
| 1317 | REST API | 内网 |
| 9090 | gRPC | 内网 |

---

## 升级流程

### 1. 编译新版本

```bash
cd /root/natchain
git pull
make build
```

### 2. 替换二进制

```bash
systemctl stop natchaind
cp build/natchaind /usr/local/bin/
systemctl start natchaind
```

### 3. 验证

```bash
natchaind version
systemctl status natchaind
```
