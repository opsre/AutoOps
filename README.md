# AutoOps 天枢

**把基础设施、可观测、自动化和日常运维流程放进同一个工作台。**

AutoOps 天枢是一套开源、支持私有化部署的运维管理平台，面向运维与 SRE 团队。平台以 CMDB 作为资源上下文，将主机与云资源、Kubernetes、日志告警、数据库、任务工单和 AIOps 接入统一控制面，并共享身份、权限、凭据与审计能力。

| 在线测试环境 | 测试账号 | 当前部署版本 |
| --- | --- | --- |
| [打开 AutoOps 测试环境](https://autoops.xin/dashboard) | `test` / `123456` | `v6.0` |

官网文档：[功能说明](https://autoops.xin/features.html) · [部署指南](https://autoops.xin/deploy.html) · [更新记录](https://autoops.xin/updates.html)

## 当前能力

以下清单依据当前 `autoops-web` 页面路由、`autoops-api` 接口和官网功能说明整理。

### CMDB、资产与堡垒机

- **主机资产**：主机列表、分组、详情、批量导入导出、资源使用监控、健康状态和资产大屏。
- **云资源同步**：支持阿里云、腾讯云、百度云、华为云、AWS 和火山引擎等云账号及云主机管理。
- **机房与网络设备**：IDC、机房机柜、交换机资产、SNMP 监控和网络设备 SSH 终端。
- **Windows 主机**：Windows 主机详情与远程连接相关管理。
- **业务模型**：服务与资产绑定、业务拓扑、健康快照、关联日志分析和监控告警。
- **远程访问**：SSH 终端、SFTP 文件传输、RDP 远程桌面；支持 JumpServer 资产接入及主机授权。
- **Agent 管理**：探针部署、状态与心跳管理；支持在线部署和离线二进制上传部署。
- **巡检与治理**：主机巡检、资产分析和资源治理视图。
- **凭据与密钥**：主机凭据、通用账号凭据、云资源密钥和访问授权。

### Kubernetes

- **集群与资源**：多集群、节点、命名空间、事件、Deployment、Pod、Service、Ingress、PVC、PV、StorageClass、ConfigMap 和 Secret。
- **工作负载操作**：扩缩容、重启、暂停、恢复、回滚、驱逐等管理操作。
- **容器运维**：Pod 实时日志、终端、文件管理、Bash / Sh 连接和测试 Pod。
- **配置与网络**：资源 YAML 查看和校验、网络资源、存储资源及配置管理。
- **监控与诊断**：集群监控、集群巡检、应用诊断和 Arthas 诊断。

### 监控、日志与告警

- **统一监控**：主机、Kubernetes、交换机及数据库和中间件监控；主机监控详情、仪表盘和监控大屏。
- **指标采集**：Agent 采集和监控数据源管理；支持 All-in-One 采集主机、MySQL、PostgreSQL 与 Redis 数据。
- **日志管理**：日志项目与存储库、解析模板、Vector Agent、采集任务、Kubernetes 日志采集和采集器管理。
- **日志查询**：日志检索、业务访问日志、数据库慢日志、日志链路、日志预览和采集日志详情。
- **告警管理**：告警数据源、规则、模板、事件、静默、通知、日志告警和告警大屏。
- **域名与 API**：域名监控、故障管理和 API 监控；提供证书申请及部署相关能力。
- **通知渠道**：支持飞书、钉钉、企业微信和邮件等通知配置。

### 数据库与云成本

- **数据库资产**：MySQL、PostgreSQL、Redis、Elasticsearch 和 MongoDB 资产管理与连接管理。
- **SQL 运维**：SQL 查询、补全、执行历史和结果导出；支持 DML、DDL 与查询审核、审批、执行及回滚。
- **数据库治理**：数据库监控、慢日志、备份、巡检和数据库 AI 助手。
- **云成本**：多云账号、账单同步、账单总览、费用趋势与分布、余额告警和账单明细。

### 任务、作业与工单

- **任务调度**：脚本任务、任务模板、定时调度、Ansible Playbook、Kubernetes 任务和运行监控。
- **作业管理**：运维作业管理、执行状态、队列指标和运行日志。
- **工单审批**：审批工作台、日常事务工单、审批流配置和工单状态跟踪。
- **应用发布**：Jenkins 发布工单与应用发布流程；支持执行、重试、中止和查看构建日志。
- **SQL 工单**：数据库 SQL 审核工单与审批执行流程。

### AIOps

- **智能助手**：结合主机、告警、数据库和 Kubernetes 上下文进行运维问答与分析。
- **分析与巡检**：主机及 Kubernetes 日志分析、资源分析、巡检和告警分析。
- **模型与知识**：模型提供商配置、运维知识库和声明式 Skill 管理。
- **MCP 工具**：MCP 服务和工具目录管理、工具配置及受控操作。
- **Agent 运行**：Agent Run 状态、流式结果和运行记录。

### 系统管理与审计

- **身份权限**：用户、角色、部门、岗位、菜单和按钮权限管理；支持 LDAP 配置与认证。
- **访问控制**：访客只读权限、堡垒机主机授权和连接前权限检查。
- **操作审计**：登录日志、操作日志、活动记录、终端会话录像、命令审计和敏感命令风险记录。

### v6.0 更新重点

- 优化 AI 助手分析流程及 Kubernetes AI 分析结果。
- 增加主机一键巡检、资产分析与资源治理视图。
- 支持上传离线 Agent 二进制并按目标主机架构部署。
- 增强日志采集器配置、异常提醒、日志预览和采集详情。
- 增强多厂商交换机监控和交换机 SSH 终端。
- 增加覆盖主机、Kubernetes、交换机和数据库中间件的统一监控大屏。

更详细的模块说明见[官网功能说明](https://autoops.xin/features.html)和[更新记录](https://autoops.xin/updates.html)。

## 项目组成

| 项目 | 说明 | 仓库 |
| --- | --- | --- |
| **AutoOps** | 部署脚本、Docker Compose 配置与项目介绍 | [GitHub](https://github.com/opsre/AutoOps) |
| **autoops-api** | Go / Gin 后端 API、调度与运维服务 | [Gitee](https://gitee.com/zhang_fan1024/autoops-api) |
| **autoops-web** | Vue 3 / Element Plus 管理工作台 | [Gitee](https://gitee.com/zhang_fan1024/autoops-web) |
| **autoops-site** | 产品官网、功能说明、部署指南与更新记录 | [Gitee](https://gitee.com/zhang_fan1024/autoops-site) |

部署包通过 Docker Compose 编排 Web、API、MySQL、Redis、VictoriaMetrics、VictoriaLogs、Kafka 和 Vector。前后端使用 Go 与 Vue 3 构建。

## 快速部署

需要一台 Linux 主机、Docker 20.10 或更高版本，以及 Docker Compose 2.0 或更高版本。4 核 CPU、4 GB 内存适用于体验和小规模部署；生产环境请根据资产规模、指标保留时间和日志吞吐规划资源。

```bash
git clone https://github.com/opsre/AutoOps.git
cd AutoOps/docker
chmod +x devops-start.sh devops-stop.sh
./devops-start.sh v6.0 <服务器IP或域名> 8080
```

启动脚本会拉取镜像、生成运行配置并启动服务。部署完成后访问 `http://<服务器IP或域名>:8080`，初始管理员账号为 `admin`，初始密码为 `123456`。首次登录后请立即修改密码。

```bash
# 查看容器状态
docker compose ps

# 查看实时日志
docker compose logs -f

# 停止服务并保留持久化数据
./devops-stop.sh
```

启动脚本还支持自定义 API、MySQL、Redis、VictoriaMetrics、VictoriaLogs 和 Kafka 端口。镜像预拉取、端口参数、服务清单及更新方式见[完整部署指南](https://autoops.xin/deploy.html)。

## 资源

- 在线测试环境：[https://autoops.xin/dashboard](https://autoops.xin/dashboard)（账号：`test`，密码：`123456`）
- 产品官网：[https://autoops.xin/](https://autoops.xin/)
- 功能说明：[https://autoops.xin/features.html](https://autoops.xin/features.html)
- 更新记录：[https://autoops.xin/updates.html](https://autoops.xin/updates.html)
- 安装与使用视频：[哔哩哔哩](https://www.bilibili.com/video/BV179Wxz1Ez6/)

## 联系与交流

- 技术交流微信：`zf5391621`
- 建议邮箱：`zfwh1024@163.com`

<details>
<summary>鸣谢项目支持者</summary>

<p align="center">
  <img src="assets/zanzhu/1.png" width="100" alt="项目支持者" />
  <img src="assets/zanzhu/4.png" width="100" alt="项目支持者" />
  <img src="assets/zanzhu/5.png" width="100" alt="项目支持者" />
  <img src="assets/zanzhu/6.png" width="100" alt="项目支持者" />
  <img src="assets/zanzhu/7.png" width="100" alt="项目支持者" />
  <img src="assets/zanzhu/8.png" width="100" alt="项目支持者" />
  <img src="assets/zanzhu/9.png" width="100" alt="项目支持者" />
  <img src="assets/zanzhu/10.png" width="100" alt="项目支持者" />
  <img src="assets/zanzhu/11.png" width="100" alt="项目支持者" />
  <img src="assets/zanzhu/12.png" width="100" alt="项目支持者" />
  <img src="assets/zanzhu/13.png" width="100" alt="项目支持者" />
  <img src="assets/zanzhu/14.png" width="100" alt="项目支持者" />
</p>
</details>
