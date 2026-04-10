#!/bin/bash

# ========================================
# 监控告警自建引擎: 告警群组 & 子规则 & 分类(Style) API 测试脚本
# ========================================

BASE_URL="http://172.22.107.76:8000/api/v1"
TOKEN="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ODksInVzZXJuYW1lIjoiYWRtaW4iLCJuaWNrbmFtZSI6ImFkbWluIiwiaWNvbiI6Imh0dHA6Ly8xOTIuMTY4LjMuNzo4MDgwL2FwaS92MS91cGxvYWQvMjAyNTEyMTMvODYyMzI4MDAwLnBuZyIsImVtYWlsIjoiMTIzNDU2Nzg5QHFxLmNvbSIsInBob25lIjoiMTM3NTQzNTQ1MzYiLCJub3RlIjoi5ZCO56uv56CU5Y-RIiwiZXhwIjoxNzc1NzE5OTU0LCJpc3MiOiJhZG1pbiJ9.fScdJkkKnS1rhDuODHj5udUREFgnyw-4keBGnDHS1Gg"

HEADER_JSON="Content-Type: application/json"
HEADER_AUTH="Authorization: Bearer ${TOKEN}"
DS_ID=1

echo "========================================"
echo "0. 创建告警分类 (POST /monitor/alert/style)"
echo "========================================"
CREATE_STYLE_RES=$(curl -s -X POST "${BASE_URL}/monitor/alert/style" -H "${HEADER_JSON}" -H "${HEADER_AUTH}" -d '{"name": "CPU", "description": "CPU指标告警规则"}')
STYLE_ID=$(echo "${CREATE_STYLE_RES}" | jq -r '.data.ID // .data.id')

echo "========================================"
echo "1. 创建告警群组 (POST /monitor/alert/group)"
echo "========================================"
YAML_RULE="apiVersion: monitoring.coreos.com/v1\nkind: PrometheusRule\nmetadata:\n  name: node-cpu-usage\n  namespace: monitor\n  labels:\n    release: prometheus\nspec:\n  groups:\n    - name: node.cpu.usage.rules\n      rules: []"

CREATE_GROUP_PAYLOAD=$(cat <<PAYLOAD
{
  "data_source_id": ${DS_ID},
  "group_name": "node-system-usage",
  "labels": "{\"cluster\": \"beijing-core\"}",
  "rule_content": "${YAML_RULE}"
}
PAYLOAD
)

CREATE_GROUP_RES=$(curl -s -X POST "${BASE_URL}/monitor/alert/group" -H "${HEADER_JSON}" -H "${HEADER_AUTH}" -d "${CREATE_GROUP_PAYLOAD}")
GROUP_ID=$(echo "${CREATE_GROUP_RES}" | jq -r '.data.ID // .data.id')
echo "Group ID: ${GROUP_ID}"

echo ""
echo "========================================"
echo "2. 直接创建告警规则测试评价引擎 (POST /monitor/alert/rule)"
echo "========================================"
CREATE_RULE_PAYLOAD1=$(cat <<PAYLOAD
{
  "group_id": ${GROUP_ID},
  "alert": "FakeTestCPU",
  "expr": "up == 1",
  "for_duration": "1s",
  "severity": "warning",
  "summary": "自建引擎测试",
  "description": "总是触发的测试用例",
  "constraints": "{\"service\": \"nginx\", \"env\": \"prod\"}",
  "labels": "{\"team\": \"devops\"}",
  "style": "CPU",
  "enabled": 1
}
PAYLOAD
)

CREATE_RULE_RES1=$(curl -s -X POST "${BASE_URL}/monitor/alert/rule" -H "${HEADER_JSON}" -H "${HEADER_AUTH}" -d "${CREATE_RULE_PAYLOAD1}")
echo "新增子规则1(表达式 up==1, 启用)返回结果: ${CREATE_RULE_RES1}"
RULE_ID1=$(echo "${CREATE_RULE_RES1}" | jq -r '.data.ID // .data.id')

echo ""
echo "========================================"
echo "3. 等待引擎 Evaluate 评估并检查后台日志 (Sleep 5s)"
echo "========================================"
echo ">> 自建引擎使用 'up == 1' 应该会很快命中并触发 firing..."
sleep 5

echo ""
echo "========================================"
echo "4. 获取单条子规则验证状态是否变为 pending/firing"
echo "========================================"
curl -s -X GET "${BASE_URL}/monitor/alert/rules_list?alert=FakeTestCPU&page=1&pageSize=1" -H "${HEADER_AUTH}" | jq '.data.list[0] | {alert, status}'

echo ""
echo "========================================"
echo "5. 测试完毕！系统日志中应该会打印 '[状态更新] 自研引擎检查规则'"
echo "========================================"

