#!/bin/bash

# API监控系统测试脚本

BASE_URL="http://localhost:8000/api/v1/monitor"
TOKEN="YOUR_AUTH_TOKEN"  # 替换为实际的认证Token

echo "========================================="
echo "API监控系统功能测试"
echo "========================================="
echo ""

# 1. 测试API接口连通性
echo "1️⃣  测试API接口连通性..."
curl -X POST "$BASE_URL/api-monitors/test" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "url": "http://127.0.0.1:8000/api/v1/cmdb/grouplist",
    "method": "GET",
    "timeout": 5000,
    "expectedStatus": 200
  }' | jq .

echo ""
echo "----------------------------------------"
echo ""

# 2. 创建API监控
echo "2️⃣  创建API监控配置..."
MONITOR_ID=$(curl -s -X POST "$BASE_URL/api-monitors" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "业务组列表API监控",
    "description": "监控业务组列表接口可用性",
    "groupId": 1,
    "groupName": "测试环境",
    "serviceName": "cmdb-service",
    "url": "http://127.0.0.1:8000/api/v1/cmdb/grouplist",
    "method": "GET",
    "timeout": 5000,
    "expectedStatus": 200,
    "checkInterval": 30,
    "retryCount": 3,
    "alertEnabled": true,
    "status": 1
  }' | jq -r '.data.id')

echo "创建的监控ID: $MONITOR_ID"
echo ""
echo "----------------------------------------"
echo ""

# 3. 获取监控详情
echo "3️⃣  获取监控详情..."
curl -s "$BASE_URL/api-monitors/$MONITOR_ID" \
  -H "Authorization: Bearer $TOKEN" | jq .

echo ""
echo "----------------------------------------"
echo ""

# 4. 获取监控列表
echo "4️⃣  获取监控列表..."
curl -s "$BASE_URL/api-monitors?page=1&pageSize=10" \
  -H "Authorization: Bearer $TOKEN" | jq .

echo ""
echo "----------------------------------------"
echo ""

# 5. 获取统计信息
echo "5️⃣  获取统计信息..."
curl -s "$BASE_URL/api-monitors/stats" \
  -H "Authorization: Bearer $TOKEN" | jq .

echo ""
echo "----------------------------------------"
echo ""

# 6. 获取业务组统计
echo "6️⃣  获取业务组统计..."
curl -s "$BASE_URL/api-monitors/group-stats" \
  -H "Authorization: Bearer $TOKEN" | jq .

echo ""
echo "----------------------------------------"
echo ""

# 7. 更新监控状态（禁用）
echo "7️⃣  禁用监控..."
curl -s -X PUT "$BASE_URL/api-monitors/$MONITOR_ID/status?status=0" \
  -H "Authorization: Bearer $TOKEN" | jq .

echo ""
echo "----------------------------------------"
echo ""

# 8. 再次启用监控
echo "8️⃣  启用监控..."
curl -s -X PUT "$BASE_URL/api-monitors/$MONITOR_ID/status?status=1" \
  -H "Authorization: Bearer $TOKEN" | jq .

echo ""
echo "----------------------------------------"
echo ""

echo "✅ 测试完成！"
echo ""
echo "📊 查看 Pushgateway 指标:"
echo "   curl http://localhost:9091/metrics | grep api_monitor"
echo ""
echo "📈 查看 Prometheus 指标:"
echo "   访问 http://localhost:9090/graph"
echo "   查询: api_monitor_success"
echo ""
echo "🗑️  删除测试监控:"
echo "   curl -X DELETE $BASE_URL/api-monitors/$MONITOR_ID -H 'Authorization: Bearer $TOKEN'"
echo ""
