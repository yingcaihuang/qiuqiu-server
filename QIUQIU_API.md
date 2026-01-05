# QiuQiu Server API Documentation

This document describes the REST API for pushing messages to QiuQiu iOS app via this modified bark-server.

## Architecture Overview

QiuQiu Server is built on top of bark-server with the following additions:
- **QiuQiu Message Model**: Structured message format for the QiuQiu app
- **Webhook Integration**: HTTP endpoints for external services to push notifications
- **Message Persistence**: All messages are stored in local database
- **Device Token Management**: Automatic device registration and token tracking

## API Endpoints

### 1. Push Message to QiuQiu App

**Endpoint**: `POST /qiuqiu/push` or `POST /api/push`

**Description**: Send a message to a registered QiuQiu device

**Request Headers**:
```
Content-Type: application/json
```

**Request Body**:
```json
{
  "token": "device_token_from_qiuqiu_app",
  "title": "Alert Title",
  "message": "Alert message body with **markdown** support",
  "url": "https://example.com/details",
  "timestamp": 1704412800
}
```

**Request Fields**:
- `token` (required): Device push token obtained from QiuQiu app settings
- `title` (optional): Alert title, default is "Alert" if not provided
- `message` (required): Alert message content, supports markdown formatting
- `url` (optional): URL to open when user taps on the notification
- `timestamp` (optional): Unix timestamp, default is current time if not provided

**Response - Success (200)**:
```json
{
  "code": 200,
  "data": {
    "code": 200,
    "message": "success",
    "token": "device_token_from_qiuqiu_app"
  },
  "message": "success",
  "timestamp": 1704412800
}
```

**Response - Failure (400/500)**:
```json
{
  "code": 400,
  "message": "token is required",
  "timestamp": 1704412800
}
```

**Error Codes**:
- `400`: Invalid request (missing required fields, invalid format)
- `500`: Server error (database or APNs push failure)

### 2. Retrieve Messages by Device Token

**Endpoint**: `GET /qiuqiu/messages/:token`

**Description**: Retrieve all messages sent to a specific device

**URL Parameters**:
- `token`: Device push token

**Response - Success (200)**:
```json
{
  "code": 200,
  "data": {
    "messages": [
      {
        "id": "token-1704412800",
        "token": "device_token",
        "title": "Alert Title",
        "message": "Alert message",
        "url": "https://example.com",
        "timestamp": 1704412800,
        "created_at": "2024-01-04T12:00:00Z"
      }
    ],
    "token": "device_token"
  },
  "message": "success",
  "timestamp": 1704412801
}
```

**Response - Failure (400/500)**:
```json
{
  "code": 500,
  "message": "failed to retrieve messages: database error",
  "timestamp": 1704412801
}
```

## Integration Guide

### Step 1: Get Device Token from QiuQiu App

1. Open QiuQiu app on your iPhone
2. Go to Settings tab
3. Turn off "Follow System Language" if needed
4. Under "Push Token" section, tap "Generate Token"
5. Copy the generated token

### Step 2: Setup QiuQiu Server

```bash
# Clone the qiuqiu-server repository
git clone <repo-url>
cd qiuqiu-server

# Build the server
go build -o qiuqiu-server

# Run the server
./qiuqiu-server --addr 0.0.0.0:8080 --data ./data
```

### Step 3: Configure Webhook in Your Application

For external services to push messages, use the `/api/push` endpoint:

```javascript
// Example: Node.js with fetch
const pushMessage = async (token, title, message, url = null) => {
  const payload = {
    token: token,
    title: title,
    message: message,
    url: url,
    timestamp: Math.floor(Date.now() / 1000)
  };

  const response = await fetch('http://your-server.com:8080/api/push', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(payload)
  });

  return await response.json();
};

// Usage
pushMessage(
  'your-device-token',
  'Disk Alert',
  'Server disk usage **95%**',
  'https://monitor.example.com/disk'
);
```

### Step 4: Webhook Integration Examples

#### Monitoring System Alert
```bash
curl -X POST http://localhost:8080/api/push \
  -H "Content-Type: application/json" \
  -d '{
    "token": "your-device-token",
    "title": "Disk Space Alert",
    "message": "Server disk usage is at **98%**, only 2GB remaining",
    "url": "https://admin.example.com/storage",
    "timestamp": '$(date +%s)'
  }'
```

#### CI/CD Pipeline Notification
```bash
# Send build success notification
curl -X POST http://localhost:8080/api/push \
  -H "Content-Type: application/json" \
  -d '{
    "token": "your-device-token",
    "title": "Build Success",
    "message": "Build #123 completed successfully\n\n- Branch: main\n- Duration: 5m 32s",
    "url": "https://ci.example.com/builds/123"
  }'
```

#### Python Script
```python
import requests
import time

def send_qiuqiu_notification(token, title, message, url=None):
    payload = {
        "token": token,
        "title": title,
        "message": message,
        "url": url,
        "timestamp": int(time.time())
    }
    
    response = requests.post(
        'http://your-server.com:8080/api/push',
        json=payload,
        headers={'Content-Type': 'application/json'}
    )
    
    return response.json()

# Usage
result = send_qiuqiu_notification(
    token="your-device-token",
    title="Temperature Alert",
    message="CPU temperature is **85°C**, above threshold",
    url="https://monitor.example.com"
)
print(result)
```

#### cURL with Response Processing
```bash
#!/bin/bash

TOKEN="your-device-token"
TITLE="System Alert"
MESSAGE="Server reboot completed successfully"
URL="https://admin.example.com"

RESPONSE=$(curl -s -X POST http://localhost:8080/api/push \
  -H "Content-Type: application/json" \
  -d "{
    \"token\": \"$TOKEN\",
    \"title\": \"$TITLE\",
    \"message\": \"$MESSAGE\",
    \"url\": \"$URL\",
    \"timestamp\": $(date +%s)
  }")

echo "Response: $RESPONSE"
```

## Message Format Details

### Markdown Support

Messages support markdown formatting as follows:
- **Bold**: `**text**` or `__text__`
- *Italic*: `*text*` or `_text_`
- [Links](https://example.com): `[text](url)`
- Code: `` `code` ``
- Headers: `# Header 1`, `## Header 2`, etc.
- Lists: `- item` or `* item` for bullets

Example message:
```
Server **disk usage** is critical!

Please review the following:
- Node 1: 98% full
- Node 2: 95% full

Visit [Dashboard](https://admin.example.com) for more info
```

## Testing

### Test Local Setup

```bash
# 1. Start server
./qiuqiu-server --addr 0.0.0.0:8080 --data ./data

# 2. In another terminal, send test message
curl -X POST http://localhost:8080/api/push \
  -H "Content-Type: application/json" \
  -d '{
    "token": "test-token-12345",
    "title": "Test Alert",
    "message": "This is a test message",
    "timestamp": '$(date +%s)'
  }'

# 3. Retrieve messages
curl http://localhost:8080/qiuqiu/messages/test-token-12345
```

### Monitoring API Health

```bash
# Check server is running
curl http://localhost:8080/ping

# Get server info
curl http://localhost:8080/info
```

## Configuration

### Environment Variables

- `BARK_ADDR`: Server address (default: `0.0.0.0:8080`)
- `BARK_DATA`: Data directory for storage (default: `./data`)

### Command Line Options

```bash
./qiuqiu-server --help

Options:
  --addr              Server address (default "0.0.0.0:8080")
  --data              Data directory (default "./data")
  --concurrency       Max concurrent connections
  --read-timeout      Read timeout duration
  --write-timeout     Write timeout duration
```

## Troubleshooting

### Issue: Message not received on device

**Solution**:
1. Verify device token is correct and copied from QiuQiu app
2. Ensure QiuQiu app has notification permission enabled on iPhone
3. Check if server is running: `curl http://server:8080/ping`
4. Check server logs for error messages

### Issue: Invalid device token error

**Solution**:
1. Re-generate token in QiuQiu app settings
2. Ensure token is not longer than 128 characters
3. Token should be copied exactly without spaces

### Issue: CORS or connection refused

**Solution**:
1. Verify server is listening on correct port: `netstat -tulpn | grep 8080`
2. Check firewall rules allow access to port 8080
3. For remote connections, ensure using full domain/IP: `http://your-domain:8080/api/push`

## Security Considerations

- Keep your device token private; it grants access to push notifications to your device
- Use HTTPS in production environments
- Consider implementing authentication if server is exposed to the internet
- Store device tokens securely on your server

## Support

For issues or feature requests, please visit the repository or contact support.
