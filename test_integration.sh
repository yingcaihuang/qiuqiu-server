#!/bin/bash

# QiuQiu Server Integration Test Script
# 测试 QiuQiu 服务器的 webhook 推送功能

set -e

SERVER_URL="${1:-http://localhost:8080}"
DEVICE_TOKEN="${2:-test-token-12345}"

echo "=========================================="
echo "QiuQiu Server Integration Tests"
echo "=========================================="
echo ""
echo "Server URL: $SERVER_URL"
echo "Device Token: $DEVICE_TOKEN"
echo ""

# Test 1: Server Health Check
echo "Test 1: Server Health Check"
if curl -s "$SERVER_URL/ping" | grep -q "pong"; then
  echo "✅ Server is running"
else
  echo "❌ Server is not responding"
  exit 1
fi
echo ""

# Test 2: Get Server Info
echo "Test 2: Get Server Info"
curl -s "$SERVER_URL/info" | python3 -m json.tool
echo ""

# Test 3: Send Simple Message
echo "Test 3: Send Simple Message"
RESPONSE=$(curl -s -X POST "$SERVER_URL/api/push" \
  -H "Content-Type: application/json" \
  -d "{
    \"token\": \"$DEVICE_TOKEN\",
    \"title\": \"Test Alert\",
    \"message\": \"This is a simple test message\",
    \"timestamp\": $(date +%s)
  }")
echo "$RESPONSE" | python3 -m json.tool
echo ""

# Test 4: Send Message with Markdown
echo "Test 4: Send Message with Markdown Formatting"
RESPONSE=$(curl -s -X POST "$SERVER_URL/api/push" \
  -H "Content-Type: application/json" \
  -d "{
    \"token\": \"$DEVICE_TOKEN\",
    \"title\": \"Markdown Test\",
    \"message\": \"**Bold** and *italic* text\\n\\n- Item 1\\n- Item 2\\n- Item 3\",
    \"url\": \"https://example.com\",
    \"timestamp\": $(date +%s)
  }")
echo "$RESPONSE" | python3 -m json.tool
echo ""

# Test 5: Send Message with URL
echo "Test 5: Send Message with URL"
RESPONSE=$(curl -s -X POST "$SERVER_URL/api/push" \
  -H "Content-Type: application/json" \
  -d "{
    \"token\": \"$DEVICE_TOKEN\",
    \"title\": \"System Alert\",
    \"message\": \"Server disk usage is at **95%**\",
    \"url\": \"https://admin.example.com/storage\",
    \"timestamp\": $(date +%s)
  }")
echo "$RESPONSE" | python3 -m json.tool
echo ""

# Test 6: Retrieve Messages
echo "Test 6: Retrieve Messages for Device"
curl -s "$SERVER_URL/qiuqiu/messages/$DEVICE_TOKEN" | python3 -m json.tool
echo ""

# Test 7: Batch Message Send (Multiple Tokens - simulated)
echo "Test 7: Sending to Multiple Devices"
for i in {1..3}; do
  TOKEN="token-$i"
  echo "Sending to token: $TOKEN"
  curl -s -X POST "$SERVER_URL/api/push" \
    -H "Content-Type: application/json" \
    -d "{
      \"token\": \"$TOKEN\",
      \"title\": \"Batch Message $i\",
      \"message\": \"This is batch message #$i\",
      \"timestamp\": $(date +%s)
    }" | python3 -m json.tool
  echo ""
done

echo "=========================================="
echo "✅ All tests completed successfully!"
echo "=========================================="
