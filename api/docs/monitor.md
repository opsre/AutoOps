# 监控中心 (Monitor) API文档

## 概述
包含主机监控、Agent部署管理、以及PrometheusAlert相关功能(Webhooks接入、告警模板管理、告警路由管理、告警记录)。

### Base URL: `/api/v1`

---

## 1. Webhooks接入管理（无需鉴权）
Webhooks端点用于接收外部系统的监控与审计事件, 解析后进行告警处理。

### 1.1 接收Gitlab Webhook
- **URL**: `/monitor/alert/webhook/gitlab`
- **Method**: `POST`
- **Description**: 接收Gitlab推送的代码提交/流水线/Merge Request等事件Webhook。

### 1.2 接收Zabbix Webhook
- **URL**: `/monitor/alert/webhook/zabbix`
- **Method**: `POST`
- **Description**: 接收来自Zabbix系统的告警推送。

### 1.3 接收Prometheus Webhook
- **URL**: `/monitor/alert/webhook/prometheus`
- **Method**: `POST`
- **Description**: 接收来自Prometheus AlertManager的告警推送。

---

## 2. 告警模板管理 (Alert Templates)

### 2.1 创建模板
- **URL**: `/monitor/alert/template`
- **Method**: `POST`
- **Header**: `Authorization: Bearer <token>`
- **Description**: 创建一个新的告警通知模板(支持钉钉、微信、飞书等格式)。

### 2.2 删除模板
- **URL**: `/monitor/alert/template/:id`
- **Method**: `DELETE`
- **Header**: `Authorization: Bearer <token>`
- **Description**: 根据模板ID删除该模板。

### 2.3 更新模板
- **URL**: `/monitor/alert/template/:id`
- **Method**: `PUT`
- **Header**: `Authorization: Bearer <token>`
- **Description**: 更新指定的模板信息。

### 2.4 获取模板列表
- **URL**: `/monitor/alert/templates`
- **Method**: `GET`
- **Header**: `Authorization: Bearer <token>`
- **Description**: 获取已创建的系统告警通知模板。

---

## 3. 告警路由管理 (Alert Routers)

### 3.1 创建路由
- **URL**: `/monitor/alert/router`
- **Method**: `POST`
- **Header**: `Authorization: Bearer <token>`
- **Description**: 建立新的告警路由分发规则。

### 3.2 删除路由
- **URL**: `/monitor/alert/router/:id`
- **Method**: `DELETE`
- **Header**: `Authorization: Bearer <token>`
- **Description**: 根据ID删除告警路由规则。

### 3.3 更新路由
- **URL**: `/monitor/alert/router/:id`
- **Method**: `PUT`
- **Header**: `Authorization: Bearer <token>`
- **Description**: 更新指定的路由规则。

### 3.4 获取路由列表/详情
- **URL**: `/monitor/alert/routers`
- **Method**: `GET`
- **Header**: `Authorization: Bearer <token>`
- **Description**: 获取系统已有的路由规则列表。

---

## 4. 其它运维组件配置 (Component Config Operations)

### 4.1 重新加载配置
- **URL**: `/monitor/alert/router/reload`
- **Method**: `POST`
- **Header**: `Authorization: Bearer <token>`
- **Description**: 触发Prometheus Alert动态重新加载配置。

### 4.2 检查组件健康状态
- **URL**: `/monitor/alert/router/health`
- **Method**: `GET`
- **Header**: `Authorization: Bearer <token>`
- **Description**: 探测监控组件是否运行正常。

---

## 5. 告警记录 (Records Operations)

### 5.1 查看告警记录列表
- **URL**: `/monitor/alert/records`
- **Method**: `GET`
- **Header**: `Authorization: Bearer <token>`
- **Description**: 查询告警的分发历史记录。

### 5.2 清理告警记录
- **URL**: `/monitor/alert/records/clean`
- **Method**: `DELETE`
- **Header**: `Authorization: Bearer <token>`
- **Description**: 清理/归档过期或指定的告警记录减少存储压力。