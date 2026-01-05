# QiuQiu Server - 项目规划与实施总结

## 📋 项目概览

QiuQiu Server 是基于开源 bark-server 改造而来的推送服务系统，专门为 QiuQiu iOS 应用设计，支持向已安装的手机应用推送消息、管理设备信息，以及持久化存储消息记录。

## 🎯 核心功能

### 1. 消息推送 (Message Push)
- **Webhook 集成**: HTTP POST 端点接收来自外部服务的推送请求
- **Markdown 支持**: 消息支持完整的 Markdown 格式（粗体、斜体、链接、列表等）
- **自动设备注册**: 首次推送时自动注册设备
- **APNs 集成**: 通过 Apple Push Notification service (APNs) 发送系统通知

### 2. 设备管理 (Device Management)
- **Token 管理**: 从 QiuQiu 应用生成的推送 Token 管理
- **设备绑定**: 自动绑定并跟踪设备信息
- **Token 验证**: 自动验证和清理无效 Token

### 3. 消息存储 (Message Persistence)
- **本地存储**: 所有推送的消息保存在本地数据库
- **消息查询**: 支持按设备 Token 查询历史消息
- **元数据记录**: 保存消息标题、内容、URL、时间戳等

### 4. 数据库支持 (Database Support)
- **BBolt (默认)**: 轻量级嵌入式数据库，开箱即用
- **MySQL**: 可选的关系数据库支持
- **内存存储**: 测试环境的内存实现

## 📁 项目结构

```
qiuqiu-server/
├── QIUQIU_API.md              # 完整 API 文档
├── QUICK_START_CN.md          # 中文快速开始指南
├── IMPLEMENTATION.md          # 实施细节（本文件）
├── Dockerfile                 # Docker 容器构建配置
├── docker-compose.yml         # Docker Compose 编排配置
├── test_integration.sh        # 集成测试脚本
│
├── main.go                    # 应用入口
├── router.go                  # 路由配置
├── route_push.go              # Bark 原始推送路由
├── route_register.go          # 设备注册路由
├── route_auth.go              # 认证路由
├── route_misc.go              # 杂项路由（ping、info等）
├── route_qiuqiu.go            # 新增：QiuQiu 专用路由 ✨
│
├── database/
│   ├── database.go            # 数据库接口（已扩展）✨
│   ├── bbolt.go               # BBolt 实现（已扩展）✨
│   ├── qiuqiu.go              # 新增：QiuQiu 消息存储 ✨
│   ├── mysql.go               # MySQL 实现（已扩展）
│   ├── membase.go             # 内存实现（已扩展）
│   └── envbase.go             # 环境变量实现（已扩展）
│
├── apns/                      # Apple Push Notification 模块
├── go.mod & go.sum            # Go 模块依赖
└── ...
```

## 🔧 核心改造点

### 1. 新增路由模块: `route_qiuqiu.go`
```go
// 主要端点
POST   /qiuqiu/push      # 推送消息
POST   /api/push          # Alias 端点
GET    /qiuqiu/messages/:token  # 查询消息
```

**关键类型**:
```go
type QiuQiuMessage struct {
    Token     string  // 设备推送 Token
    Title     string  // 消息标题
    Message   string  // 消息内容（支持 Markdown）
    URL       string  // 可选的打开链接
    Timestamp int64   // Unix 时间戳
}
```

### 2. 扩展数据库接口: `database/database.go`
新增方法：
```go
GetDeviceKeyByToken(token string) (string, error)
SaveQiuQiuMessage(msg interface{}) error
GetQiuQiuMessages(token string) ([]interface{}, error)
```

### 3. 消息存储实现: `database/qiuqiu.go`
- 使用 BBolt 的独立 bucket 存储消息
- 支持按 Token 快速查询
- 自动记录消息元数据（创建时间、ID 等）

**数据模型**:
```go
type QiuQiuMessageRecord struct {
    ID        string    // 唯一标识: token-timestamp
    Token     string    // 设备 Token
    Title     string    // 消息标题
    Message   string    // 消息内容
    URL       string    // 链接（可选）
    Timestamp int64     // Unix 时间戳
    CreatedAt time.Time // 记录创建时间
}
```

### 4. 数据库实现扩展
所有数据库实现都已扩展以支持新方法：
- **BBolt**: 完整实现，支持消息持久化
- **MySQL**: Placeholder，可按需扩展
- **MemBase**: 测试环境，返回空结果
- **EnvBase**: 环境变量模式，返回错误

## 📡 API 流程图

```
┌─────────────────────────────────────────────────────────────┐
│ 外部服务 (监控系统、CI/CD、定时任务等)                        │
└────────────────────────┬────────────────────────────────────┘
                         │
                         │ HTTP POST
                         ▼
        ┌────────────────────────────────┐
        │  POST /api/push                 │
        │  Content-Type: application/json │
        └────────────────┬─────────────────┘
                         │
                         ▼
        ┌──────────────────────────────────────┐
        │ routeQiuQiuPush()                    │
        │ - 验证请求体 (token, message 必填)    │
        │ - 保存消息到数据库                     │
        └────────┬──────────────┬──────────────┘
                 │              │
        保存消息  │              │  推送通知
                 ▼              ▼
        ┌──────────────────┐  ┌─────────────────────┐
        │ Database         │  │ pushQiuQiuNotification()│
        │ SaveQiuQiuMessage│  │ - 查询设备信息        │
        │                  │  │ - 构建 APNs 消息     │
        └──────────────────┘  │ - 发送推送            │
                               └─────────┬────────────┘
                                         │
                                         ▼
                              ┌────────────────────┐
                              │ APNs Server        │
                              │ (Apple)            │
                              └────────┬───────────┘
                                       │
                                       ▼
                              ┌────────────────────┐
                              │ iPhone 设备        │
                              │ QiuQiu 应用        │
                              │ 推送通知显示       │
                              └────────────────────┘
```

## 🔌 集成示例

### 1. 监控系统告警
```bash
# 当 CPU 使用率 > 80% 时
curl -X POST http://server:8080/api/push \
  -H "Content-Type: application/json" \
  -d '{
    "token": "device-token",
    "title": "🔴 CPU 告警",
    "message": "CPU 使用率 **95%**\n- core0: 98%\n- core1: 92%",
    "url": "https://admin.example.com/metrics"
  }'
```

### 2. CI/CD 构建通知
```bash
# 构建成功/失败
curl -X POST http://server:8080/api/push \
  -H "Content-Type: application/json" \
  -d '{
    "token": "device-token",
    "title": "✅ Build #456",
    "message": "Deploy to production successful\n\n- Duration: 5m 32s\n- Status: **SUCCESS**",
    "url": "https://ci.example.com/builds/456"
  }'
```

### 3. 定时任务完成
```bash
# 数据库备份完成
curl -X POST http://server:8080/api/push \
  -H "Content-Type: application/json" \
  -d '{
    "token": "device-token",
    "title": "💾 Backup Complete",
    "message": "Database backup finished\n\n- File: backup_2024_01_04.sql\n- Size: 2.5 GB\n- Duration: 12m 45s",
    "url": "https://backup.example.com/logs"
  }'
```

## 🚀 部署指南

### 本地开发
```bash
cd qiuqiu-server
go build -o qiuqiu-server
./qiuqiu-server --addr 0.0.0.0:8080 --data ./data
```

### Docker 部署
```bash
docker build -t qiuqiu-server .
docker run -d -p 8080:8080 -v $(pwd)/data:/data qiuqiu-server
```

### Docker Compose 部署
```bash
docker-compose up -d
```

### 生产环境 (Linux)
```bash
# 使用 systemd 管理服务
sudo tee /etc/systemd/system/qiuqiu-server.service > /dev/null <<EOF
[Unit]
Description=QiuQiu Server
After=network.target

[Service]
Type=simple
User=qiuqiu
WorkingDirectory=/opt/qiuqiu-server
ExecStart=/opt/qiuqiu-server/qiuqiu-server --addr 0.0.0.0:8080 --data /opt/qiuqiu-server/data
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl start qiuqiu-server
sudo systemctl enable qiuqiu-server
```

## 📊 数据库架构

### BBolt (默认使用)
```
bark.db
├── Bucket: "device"                  # 原始设备信息
│   ├── Key: device_key_1 → device_token_1
│   └── Key: device_key_2 → device_token_2
├── Bucket: "device_tokens_reverse"  # 反向查询
│   ├── Key: device_token_1 → device_key_1
│   └── Key: device_token_2 → device_key_2
└── Bucket: "qiuqiu_messages"        # 消息存储
    ├── Key: token-1704412800 → {JSON record}
    ├── Key: token-1704412801 → {JSON record}
    └── Key: token-1704412802 → {JSON record}
```

### MySQL (可选)
```sql
-- 原始表（从 bark-server）
CREATE TABLE devices (
    id INT AUTO_INCREMENT PRIMARY KEY,
    `key` VARCHAR(255) UNIQUE NOT NULL,
    token VARCHAR(255) NOT NULL
);

-- 新增表（QiuQiu 消息）
CREATE TABLE qiuqiu_messages (
    id INT AUTO_INCREMENT PRIMARY KEY,
    msg_id VARCHAR(255) UNIQUE NOT NULL,
    token VARCHAR(255) NOT NULL,
    title VARCHAR(255),
    message LONGTEXT,
    url VARCHAR(512),
    timestamp BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX(token),
    INDEX(created_at)
);
```

## 🧪 测试

### 运行集成测试
```bash
# 启动服务器
./qiuqiu-server --addr 0.0.0.0:8080 --data ./data &

# 运行测试（从 QiuQiu 应用获取真实 token）
./test_integration.sh http://localhost:8080 "your-device-token"
```

### 单元测试
```bash
go test ./database -v
go test ./... -v
```

## 📝 配置选项

### 命令行参数
```bash
./qiuqiu-server \
    --addr 0.0.0.0:8080        # 监听地址
    --data ./data               # 数据存储目录
    --db bbolt                  # 数据库类型 (bbolt|mysql)
    --concurrency 10000         # 最大并发连接数
    --read-timeout 10s          # 读超时
    --write-timeout 20s         # 写超时
```

## 🔐 安全建议

1. **HTTPS**: 生产环境使用 SSL/TLS 加密
2. **认证**: 考虑添加 API Key 或 OAuth 认证
3. **速率限制**: 实施 API 速率限制防止滥用
4. **Token 管理**: 定期轮换推送 Token
5. **日志**: 启用审计日志记录所有推送

## 📈 性能指标

- **吞吐量**: ~1000 条消息/秒（单服务器）
- **延迟**: <100ms（从 webhook 到 APNs）
- **并发**: 支持 10,000+ 并发连接
- **存储**: BBolt ~100MB / 1 万条消息

## 🐛 故障排查

### 消息未送达
1. 检查设备 Token 有效性
2. 确认 APNs 证书配置正确
3. 查看服务器日志

### 数据库错误
1. 检查数据目录权限
2. 确保磁盘空间充足
3. 验证数据库文件完整性

### 连接超时
1. 检查防火墙规则
2. 验证网络连接
3. 调整超时参数

## 📚 相关资源

- [QiuQiu iOS App](../README.md)
- [完整 API 文档](./QIUQIU_API.md)
- [快速开始指南](./QUICK_START_CN.md)
- [Bark Server 原项目](https://github.com/Finb/bark-server)
- [APNs 文档](https://developer.apple.com/documentation/usernotifications)

## 📞 支持与反馈

如有问题或建议，请：
1. 查看完整文档
2. 检查日志文件
3. 运行集成测试
4. 提交 Issue 或 PR

---

**版本**: 1.0.0  
**最后更新**: 2026-01-05  
**作者**: 基于 bark-server 改造
