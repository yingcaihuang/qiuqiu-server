# QiuQiu Server Quick Start

快速开始指南 - 在 5 分钟内设置 QiuQiu 推送服务

## 前置条件

- Go 1.18+ 已安装
- macOS 或 Linux 系统
- QiuQiu 应用已在 iPhone 上安装

## 第一步：获取设备 Token

1. 打开 iPhone 上的 QiuQiu 应用
2. 切换到"设置"标签页
3. 在"推送 Token"区域，点击"生成 Token"
4. **复制生成的 Token**（一个长字符串）

## 第二步：编译和运行服务器

```bash
# 进入 qiuqiu-server 目录
cd /Users/betty/QiuQiu/qiuqiu-server

# 编译
go build -o qiuqiu-server

# 运行服务器
./qiuqiu-server --addr 0.0.0.0:8080 --data ./data

# 你应该看到类似的输出：
# 2024/01/04 12:00:00 INF init database [./data]...
# 2024/01/04 12:00:00 INF Starting server on :8080
```

服务器现在监听 `http://localhost:8080`

## 第三步：发送测试消息

打开新的终端窗口，运行以下命令（替换 `YOUR_TOKEN` 为你复制的设备 Token）：

```bash
curl -X POST http://localhost:8080/api/push \
  -H "Content-Type: application/json" \
  -d '{
    "token": "YOUR_TOKEN",
    "title": "Hello QiuQiu",
    "message": "这是第一条测试消息！",
    "timestamp": '$(date +%s)'
  }'
```

✅ 如果成功，你应该在 iPhone 上收到通知！

## 常见用途

### 1. 系统监控告警

```bash
# 磁盘空间告警
curl -X POST http://localhost:8080/api/push \
  -H "Content-Type: application/json" \
  -d '{
    "token": "YOUR_TOKEN",
    "title": "磁盘告警",
    "message": "服务器磁盘使用率 **95%**，剩余容量仅 100GB",
    "url": "https://monitor.example.com/disk",
    "timestamp": '$(date +%s)'
  }'
```

### 2. CI/CD 构建通知

```bash
# 构建成功通知
curl -X POST http://localhost:8080/api/push \
  -H "Content-Type: application/json" \
  -d '{
    "token": "YOUR_TOKEN",
    "title": "✅ 构建成功",
    "message": "Build #456 completed\n- 分支: main\n- 耗时: 5分32秒",
    "url": "https://ci.example.com/builds/456",
    "timestamp": '$(date +%s)'
  }'
```

### 3. 定时任务结果

```bash
# 备份完成通知
curl -X POST http://localhost:8080/api/push \
  -H "Content-Type: application/json" \
  -d '{
    "token": "YOUR_TOKEN",
    "title": "备份完成",
    "message": "数据库备份已完成\n- 文件大小: 2.5GB\n- 耗时: 12分钟",
    "timestamp": '$(date +%s)'
  }'
```

## 编程语言示例

### Python

```python
import requests
import time

TOKEN = "YOUR_TOKEN"

response = requests.post(
    'http://localhost:8080/api/push',
    json={
        "token": TOKEN,
        "title": "Python Alert",
        "message": "Sent from Python",
        "timestamp": int(time.time())
    }
)

print(response.json())
```

### Node.js

```javascript
const fetch = require('node-fetch');

const TOKEN = "YOUR_TOKEN";

fetch('http://localhost:8080/api/push', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    token: TOKEN,
    title: "Node.js Alert",
    message: "Sent from Node.js",
    timestamp: Math.floor(Date.now() / 1000)
  })
})
.then(r => r.json())
.then(data => console.log(data));
```

### Bash Script

```bash
#!/bin/bash

TOKEN="YOUR_TOKEN"
TITLE="${1:-Default Title}"
MESSAGE="${2:-Default Message}"
URL="${3:-}"

curl -X POST http://localhost:8080/api/push \
  -H "Content-Type: application/json" \
  -d "{
    \"token\": \"$TOKEN\",
    \"title\": \"$TITLE\",
    \"message\": \"$MESSAGE\",
    \"url\": \"$URL\",
    \"timestamp\": $(date +%s)
  }"

# 使用方法：
# ./push.sh "告警标题" "告警信息" "https://example.com"
```

## 查看已发送的消息

```bash
# 获取发送给该设备的所有消息
curl http://localhost:8080/qiuqiu/messages/YOUR_TOKEN
```

## Docker 部署

```bash
# 构建 Docker 镜像
docker build -t qiuqiu-server .

# 运行容器
docker run -d \
  -p 8080:8080 \
  -v $(pwd)/data:/data \
  --name qiuqiu-server \
  qiuqiu-server

# 查看日志
docker logs -f qiuqiu-server
```

## 远程部署

### 在服务器上运行

```bash
# 1. 复制代码到服务器
scp -r qiuqiu-server user@your-server.com:/home/user/

# 2. SSH 连接到服务器
ssh user@your-server.com

# 3. 编译和运行
cd /home/user/qiuqiu-server
go build -o qiuqiu-server
nohup ./qiuqiu-server --addr 0.0.0.0:8080 --data ./data > server.log 2>&1 &

# 4. 验证服务
curl http://localhost:8080/ping
```

### 通过远程服务器发送消息

```bash
# 从本地发送到远程服务器
curl -X POST http://your-server.com:8080/api/push \
  -H "Content-Type: application/json" \
  -d '{
    "token": "YOUR_TOKEN",
    "title": "Remote Server Alert",
    "message": "Message from remote server",
    "timestamp": '$(date +%s)'
  }'
```

## 配置 APNs 证书

为了让 QiuQiu 服务器能够发送推送通知到 iPhone，你需要配置 Apple Push Notification (APNs) 认证。

**APNs 支持两种认证方式：**
- **P8 格式（推荐）** ✅ 苹果推荐，**可设置无限有效期**，无需更新
- **PEM 格式（旧格式）** ⚠️ 需要每年更新，不推荐使用

### 方式一：使用 P8 格式（推荐）

#### 第一步：获取 P8 密钥文件

1. **访问 Apple Developer 网站**
   - 登录 [developer.apple.com](https://developer.apple.com)
   - 进入 "Certificates, Identifiers & Profiles"
   - 选择 "Keys"

2. **创建新的密钥**
   - 点击 "+" 创建新密钥
   - 选择 "Apple Push Notifications service (APNs)"
   - 勾选该选项后点击 "Continue"
   - **选择有效期：**
     - "No Expiration"（无限有效期）**推荐** ✅
     - 或选择具体年数（例如 1 年、2 年）
   - 完成注册，下载 `.p8` 文件
   - **妥善保管此文件，后续无法重新下载**

3. **记录必要信息**
   - **Key ID**: 在密钥详情页面显示（例如：`ABC123DEF4`）
   - **Team ID**: 在 Apple Developer 账户右上角显示（例如：`ABCD123456`）
   - **Bundle ID**: 你的应用标识（例如：`com.example.qiuqiu`）

#### 第二步：配置服务器（P8 格式）

1. **放置 P8 文件**
   ```bash
   mkdir -p qiuqiu-server/apns
   cp AuthKey_ABC123DEF4.p8 qiuqiu-server/apns/
   ```

2. **配置环境变量**
   ```bash
   export BARK_APNS_P8_PATH="./apns/AuthKey_ABC123DEF4.p8"
   export BARK_APNS_KEY_ID="ABC123DEF4"
   export BARK_APNS_TEAM_ID="ABCD123456"
   export BARK_APNS_BUNDLE_ID="com.example.qiuqiu"
   ```

3. **启动服务器**
   ```bash
   ./qiuqiu-server \
     --addr 0.0.0.0:8080 \
     --data ./data \
     --apns-p8 ./apns/AuthKey_ABC123DEF4.p8 \
     --apns-key-id ABC123DEF4 \
     --apns-team-id ABCD123456 \
     --apns-bundle-id com.example.qiuqiu
   ```

4. **Docker Compose 配置**
   ```yaml
   version: '3.8'
   services:
     qiuqiu-server:
       build:
         context: .
         dockerfile: Dockerfile
       ports:
         - "8080:8080"
       volumes:
         - ./data:/data
         - ./apns:/apns
       environment:
         - BARK_APNS_P8_PATH=/apns/AuthKey_ABC123DEF4.p8
         - BARK_APNS_KEY_ID=ABC123DEF4
         - BARK_APNS_TEAM_ID=ABCD123456
         - BARK_APNS_BUNDLE_ID=com.example.qiuqiu
       restart: unless-stopped
   ```

---

### 方式二：使用 PEM 格式（旧格式，不推荐）

如果你已经有 PEM 格式的证书（或者需要使用旧的 .cer/.p12 证书），也可以使用此方式。

#### 第一步：获取证书文件

1. **访问 Apple Developer 网站**
   - 登录 [developer.apple.com](https://developer.apple.com)
   - 进入 "Certificates, Identifiers & Profiles"

2. **创建 App ID**
   - 选择 "Identifiers"
   - 点击 "+" 创建新的 App ID
   - Bundle ID: `com.example.qiuqiu`（替换为你的应用标识）
   - 启用 "Push Notifications" 功能
   - 点击 "Continue" 并完成注册

3. **创建 APNs 证书**
   - 返回 Certificates 页面
   - 点击 "+" 创建新证书
   - 选择 "Apple Push Notification service (APNs)" SSL Certificate
   - 选择刚创建的 App ID
   - 使用 macOS Keychain 生成证书签名请求 (CSR)
   - 下载证书文件 (`.cer`)

#### 第二步：转换证书格式（P8 → PEM）

如果你已经有 P8 文件，可以转换为 PEM：

```bash
# 将 P8 转换为 PEM
openssl pkcs8 -in AuthKey_ABC123DEF4.p8 -out certificate.pem -nocrypt

# 或者如果 P8 有密码保护
openssl pkcs8 -in AuthKey_ABC123DEF4.p8 -out certificate.pem
```

#### 第三步：转换旧格式证书（.cer/.p12 → PEM）

如果你有 `.cer` 或 `.p12` 格式的旧证书：

```bash
# 1. 将 .cer 文件转换为 .pem
openssl x509 -inform DER -outform PEM -in aps_production.cer -out certificate.pem

# 2. 从 macOS Keychain 导出私钥
# - 打开 Keychain Access
# - 找到 "Apple Push Services" 项
# - 右键 → Export
# - 保存为 `key.p12`
# - 设置密码（记住密码）

# 3. 将 .p12 转换为 PEM（输入第 2 步设置的密码）
openssl pkcs12 -in key.p12 -out key.pem -nocerts -nodes

# 4. 创建合并的证书文件
cat certificate.pem key.pem > combined.pem
```

#### 第四步：配置服务器（PEM 格式）

1. **放置证书文件**
   ```bash
   mkdir -p qiuqiu-server/apns
   cp combined.pem qiuqiu-server/apns/
   ```

2. **启动服务器**
   ```bash
   ./qiuqiu-server \
     --addr 0.0.0.0:8080 \
     --data ./data \
     --apns-cert ./apns/combined.pem
   ```

3. **Docker 部署**
   ```bash
   docker run -d \
     -p 8080:8080 \
     -v $(pwd)/data:/data \
     -v $(pwd)/apns:/apns \
     -e BARK_APNS_CERT="/apns/combined.pem" \
     --name qiuqiu-server \
     qiuqiu-server
   ```

---

## 验证 APNs 配置

### 验证 P8 文件
```bash
# 查看 P8 文件内容
openssl pkey -in AuthKey_ABC123DEF4.p8 -text

# 验证 Key ID 和 Team ID 正确性
echo "Key ID: ABC123DEF4"
echo "Team ID: ABCD123456"
```

### 验证 PEM 文件
```bash
# 验证证书文件
openssl x509 -in certificate.pem -text -noout | head -20

# 验证私钥
openssl pkey -in key.pem -text -noout | head -10

# 验证证书和私钥是否匹配
openssl x509 -noout -modulus -in certificate.pem | openssl md5
openssl pkey -noout -modulus -in key.pem | openssl md5
# 输出应该相同
```

### 测试 APNs 连接
```bash
# P8 格式测试
./qiuqiu-server --test-apns-p8 ./apns/AuthKey_ABC123DEF4.p8 \
  --apns-key-id ABC123DEF4 \
  --apns-team-id ABCD123456

# PEM 格式测试
./qiuqiu-server --test-apns ./apns/combined.pem
```

---

## 常见问题

**Q: 我应该用哪个格式？**
```
✅ 推荐：P8 格式
- 可设置无限有效期（最大优势！）
- 即使设置有效期，也比 PEM 格式更长
- 更安全，使用 Token 认证
- 苹果官方推荐

❌ 不推荐：PEM 格式（旧的 .cer/.p12）
- 需要每年更新证书
- 功能过时，不再推荐
```

**Q: 创建 P8 密钥时，有效期应该选什么？**
```
推荐选择：No Expiration（无限有效期）

优势：
- 一劳永逸，永久可用
- 无需后续维护和更新
- 最大化降低运维成本

如果需要安全轮换：
- 可以选择 1 年有效期
- 但需要定期创建新密钥并更新配置
- 对大多数用户不必要
```

**Q: P8 文件丢失了怎么办？**
```
- P8 文件是一次性下载的，丢失后无法重新下载
- 需要在 Apple Developer 中删除该 Key
- 重新创建新的密钥
- 使用新的 P8 文件重新配置
```

**Q: 如何从 P8 转换为 PEM？**
```bash
# 转换命令
openssl pkcs8 -in AuthKey_ABC123DEF4.p8 -out certificate.pem -nocrypt
```

**Q: 证书过期了怎么办？**
```
P8 格式：
- 如果创建时选择了"No Expiration"，永不过期 ✅
- 如果设置了有效期，到期后需要创建新的 Key

PEM 格式（旧）：
- 需要每年更新
- 在 Apple Developer 重新创建证书
- 下载新证书后转换并更新配置
```

**Q: 如何在 Linux 服务器上配置？**
```bash
# 上传 P8 文件到服务器
scp AuthKey_ABC123DEF4.p8 user@server.com:/opt/qiuqiu-server/apns/

# SSH 连接后设置权限
chmod 600 /opt/qiuqiu-server/apns/AuthKey_ABC123DEF4.p8

# 编辑 systemd 服务
sudo nano /etc/systemd/system/qiuqiu-server.service

# 添加以下内容到 [Service] 部分
[Service]
ExecStart=/opt/qiuqiu-server/qiuqiu-server \
    --addr 0.0.0.0:8080 \
    --data /opt/qiuqiu-server/data \
    --apns-p8 /opt/qiuqiu-server/apns/AuthKey_ABC123DEF4.p8 \
    --apns-key-id ABC123DEF4 \
    --apns-team-id ABCD123456 \
    --apns-bundle-id com.example.qiuqiu

# 重新加载并启动服务
sudo systemctl daemon-reload
sudo systemctl restart qiuqiu-server
```

**Q: 推送显示 "Invalid Certificate" 或 "authentication failed" 错误？**

P8 格式：
```bash
# 检查必要参数是否正确
echo "Key ID: ABC123DEF4"
echo "Team ID: ABCD123456"
echo "Bundle ID: com.example.qiuqiu"

# 验证 P8 文件格式
openssl pkey -in AuthKey_ABC123DEF4.p8 -text -noout | head -3
```

PEM 格式：
```bash
# 检查证书和私钥是否匹配
openssl x509 -noout -modulus -in certificate.pem | openssl md5
openssl rsa -noout -modulus -in key.pem | openssl md5

# 输出应该相同。如果不同，检查是否使用了错误的组合
```

## 支持的消息格式

消息支持 Markdown 格式：

```
**粗体** 或 __粗体__
*斜体* 或 _斜体_
[链接](https://example.com)
`代码块`
# 大标题
## 中标题
### 小标题
- 列表项
* 列表项
+ 列表项
```

## 故障排查

### 问题：收不到消息

**检查清单**：
1. ✅ Token 是否正确复制（从 QiuQiu 设置页获取）
2. ✅ iPhone 上 QiuQiu 是否有通知权限
3. ✅ 服务器是否运行：`curl http://localhost:8080/ping`
4. ✅ 服务器日志中是否有错误

### 问题：连接被拒绝

```bash
# 检查服务器是否监听
netstat -tulpn | grep 8080

# 或者
lsof -i :8080
```

### 问题：Token 无效

```bash
# 重新生成 Token：
# 1. 打开 QiuQiu 应用
# 2. 设置 → 推送 Token
# 3. 点击"重置 Token"，然后"生成 Token"
# 4. 复制新的 Token
```

## 下一步

- 查看完整 API 文档：[QIUQIU_API.md](./QIUQIU_API.md)
- 集成到你的监控系统
- 设置 HTTPS（推荐用于生产环境）
- 配置身份验证和速率限制

## 获取帮助

遇到问题？检查以下资源：
- 查看服务器日志：`server.log`
- 测试 API：`curl http://localhost:8080/info`
- 阅读完整文档：[QIUQIU_API.md](./QIUQIU_API.md)
