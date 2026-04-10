# /monitor/alert/webhook/prometheus 测试脚本
curl -X POST "http://127.0.0.1:8000/api/v1/monitor/alert/webhook/prometheus?type=wx&tpl=prometheus-wx&wxurl=https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=6b93c3ab-8216-417e-bfd2-14ee5d18db33" \
 -H "Content-Type: application/json" \
 -d '{
  "version": "4",
  "groupKey": "{}:{alertname=\"TestWebhook\"}",
  "status": "firing",
  "receiver": "prometheus-wx",
  "externalURL": "http://alertmanager:9093",
  "groupLabels": {
    "alertname": "TestWebhook"
  },
  "alerts": [
    {
      "status": "firing",
      "labels": {
        "alertname": "CPU使用率过高",
        "severity": "critical",
        "instance": "192.168.10.100"
      },
      "annotations": {
        "description": "服务器 CPU 使用率当前已达到 95% 以上"
      },
      "startsAt": "2026-04-06T10:00:00Z",
      "endsAt": "0001-01-01T00:00:00Z",
      "generatorURL": "http://prometheus:9090"
    }
  ]
}'