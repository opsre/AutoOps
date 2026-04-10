#!/bin/bash

# ========================================
# 告警完整模块测试脚本 (Template, Router, Webhook, Record)
# ========================================

BASE_URL="http://172.22.107.76:8000/api/v1"
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ODksInVzZXJuYW1lIjoiYWRtaW4iLCJuaWNrbmFtZSI6ImFkbWluIiwiaWNvbiI6Imh0dHA6Ly8xOTIuMTY4LjMuNzo4MDgwL2FwaS92MS91cGxvYWQvMjAyNTEyMTMvODYyMzI4MDAwLnBuZyIsImVtYWlsIjoiMTIzNDU2Nzg5QHFxLmNvbSIsInBob25lIjoiMTM3NTQzNTQ1MzYiLCJub3RlIjoi5ZCO56uv56CU5Y-RIiwiZXhwIjoxNzc1NzE5OTU0LCJpc3MiOiJhZG1pbiJ9.fScdJkkKnS1rhDuODHj5udUREFgnyw-4keBGnDHS1Gg"

HEADER_JSON="Content-Type: application/json"
HEADER_AUTH="Authorization: Bearer ${TOKEN}"

echo "=========================================================="
echo "1. 告警模版 (Template) 模块测试"
echo "=========================================================="

echo "[1.1] 创建告警模版 (POST /monitor/alert/template)"
CREATE_TPL_PAYLOAD=$(cat <<PAYLOAD
{
  "Tplname": "TestWXTemplate2",
  "Tpltype": "wx",
  "Tpluse": "Prometheus",
  "Tpl": "## [告警]：收到新的节点告警\n {{.Alerts}}",
  "WebhookContentType": "application/json"
}
PAYLOAD
)
curl -s -X POST "${BASE_URL}/monitor/alert/template" -H "${HEADER_JSON}" -H "${HEADER_AUTH}" -d "${CREATE_TPL_PAYLOAD}" > /dev/null

echo -e "\n[1.2] 获取告警模版列表 (GET /monitor/alert/templates)"
LIST_TPL=$(curl -s -X GET "${BASE_URL}/monitor/alert/templates" -H "${HEADER_AUTH}")
TPL_ID=$(echo "$LIST_TPL" | jq -r '.data.list[]? | select(.Tplname=="TestWXTemplate2") | .Id' | head -n 1)

echo -e "\n[1.2.1] 多条件查询告警模版测试"
curl -s -X GET "${BASE_URL}/monitor/alert/templates?page=1&pageSize=10&tplname=TestWXTemplate2&tpltype=wx&tpluse=Prometheus" -H "${HEADER_AUTH}" | jq . || true

echo "获取到的刚创建的模版 ID: ${TPL_ID}"

echo -e "\n[1.3] 更新告警模版 (PUT /monitor/alert/template/${TPL_ID})"
UPDATE_TPL_PAYLOAD=$(cat <<PAYLOAD
{
  "Id": ${TPL_ID},
  "Tplname": "TestWXTemplate2-Updated",
  "Tpltype": "wx",
  "Tpluse": "Prometheus",
  "Tpl": "## [严重告警]：\n {{.alerts}}",
  "WebhookContentType": "application/json"
}
PAYLOAD
)
curl -s -X PUT "${BASE_URL}/monitor/alert/template/${TPL_ID}" -H "${HEADER_JSON}" -H "${HEADER_AUTH}" -d "${UPDATE_TPL_PAYLOAD}" > /dev/null
echo "更新成功"

echo -e "\n=========================================================="
echo "2. 告警路由 (Router) 模块测试"
echo "=========================================================="

echo "[2.1] 创建告警路由关联模版 (POST /monitor/alert/router)"
CREATE_ROUTER_PAYLOAD=$(cat <<PAYLOAD
{
  "RouterName": "DevOps-WX-Router-2",
  "RouterTplId": "${TPL_ID}",
  "RouterPurl": "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=6b93c3ab-8216-417e-bfd2-14ee5d18db33",
  "RouterPat": "@all",
  "RouterPatRR": false,
  "RouterSendResolved": true,
  "Rules": [
    {
      "Name": "severity",
      "Value": "critical",
      "Regex": false
    }
  ]
}
PAYLOAD
)
RES_ROUTER=$(curl -s -X POST "${BASE_URL}/monitor/alert/router" -H "${HEADER_JSON}" -H "${HEADER_AUTH}" -d "${CREATE_ROUTER_PAYLOAD}")
echo "创建路由响应: $RES_ROUTER"

echo -e "\n[2.2] 获取告警路由列表 (GET /monitor/alert/routers)"
LIST_ROUTER=$(curl -s -X GET "${BASE_URL}/monitor/alert/routers" -H "${HEADER_AUTH}")
ROUTER_ID=$(echo "$LIST_ROUTER" | jq -r '.data.list[]? | select(.Name=="DevOps-WX-Router-2") | .Id' | head -n 1)
echo "获取到的新建路由ID: ${ROUTER_ID}"

echo -e "\n[2.2.1] 多条件查询告警路由测试"
curl -s -X GET "${BASE_URL}/monitor/alert/routers?page=1&pageSize=10&name=DevOps-WX-Router-2&urlOrPhone=6b93c3ab" -H "${HEADER_AUTH}" | jq . || true

echo -e "\n[2.3] 重载告警路由配置到内存 (POST /monitor/alert/router/reload)"
curl -s -X POST "${BASE_URL}/monitor/alert/router/reload" -H "${HEADER_AUTH}" > /dev/null

echo -e "\n=========================================================="
echo "3. 发送 Webhook 触发告警 (利用刚才的路由与模版)"
echo "=========================================================="

WEBHOOK_PAYLOAD=$(cat <<PAYLOAD
{
  "version": "4",
  "groupKey": "{}:{alertname=\"SystemHighLoad\"}",
  "status": "firing",
  "receiver": "prometheus-wx",
  "externalURL": "http://alertmanager.local",
  "groupLabels": {
    "alertname": "SystemHighLoad"
  },
  "alerts": [
    {
      "status": "firing",
      "labels": {
        "alertname": "SystemHighLoad",
        "severity": "critical",
        "instance": "192.168.10.100"
      },
      "annotations": {
        "summary": "系统负载异常",
        "description": "服务器 CPU 使用率当前已达到 95% 以上"
      },
      "startsAt": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
      "endsAt": "0001-01-01T00:00:00Z",
      "generatorURL": "http://prometheus.local"
    }
  ]
}
PAYLOAD
)
echo "发送 Prometheus Alertmanager 标准格式告警 (由于我们不传查询参数，它将走到全局路由匹配)..."
RES_WH=$(curl -s -X POST "${BASE_URL}/monitor/alert/webhook/prometheus" -H "${HEADER_JSON}" -H "${HEADER_AUTH}" -d "${WEBHOOK_PAYLOAD}")
echo "Webhook 响应: $RES_WH"


echo -e "\n=========================================================="
echo "4. 告警历史记录 (Records) 测试"
echo "=========================================================="
sleep 2
echo "[4.1] 查询告警历史一览 (GET /monitor/alert/records)"
curl -s -X GET "${BASE_URL}/monitor/alert/records?page=1&pageSize=10" -H "${HEADER_AUTH}" | jq . || true

echo -e "\n[4.2] 多条件查询告警历史测试"
curl -s -X GET "${BASE_URL}/monitor/alert/records?page=1&pageSize=10&alertname=SystemHighLoad&alertLevel=critical&instance=192.168.10.100&alertStatus=firing" -H "${HEADER_AUTH}" | jq . || true

echo -e "\n[4.3] 清空全部告警历史 (DELETE /monitor/alert/records/clean)"
curl -s -X DELETE "${BASE_URL}/monitor/alert/records/clean" -H "${HEADER_AUTH}" > /dev/null

echo -e "\n=========================================================="
echo "5. 测试尾声：清理本次创建的路由和模型"
echo "=========================================================="
curl -s -X DELETE "${BASE_URL}/monitor/alert/router/${ROUTER_ID}" -H "${HEADER_AUTH}" > /dev/null
curl -s -X DELETE "${BASE_URL}/monitor/alert/template/${TPL_ID}" -H "${HEADER_AUTH}" > /dev/null
echo "清理完毕"
