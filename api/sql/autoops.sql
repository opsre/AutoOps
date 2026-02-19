-- MySQL dump 10.13  Distrib 8.0.45, for Linux (x86_64)
--
-- Host: 8.130.14.34    Database: autoops
-- ------------------------------------------------------
-- Server version	8.0.45-0ubuntu0.24.04.1

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `ai_agent_chat_history`
--

DROP TABLE IF EXISTS `ai_agent_chat_history`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ai_agent_chat_history` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键''',
  `session_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''会话ID''',
  `user_id` bigint unsigned NOT NULL COMMENT '''用户ID''',
  `role` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''角色:user/assistant/system''',
  `message` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''消息内容''',
  `intent` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''识别的意图类型''',
  `intent_conf` decimal(3,2) DEFAULT NULL COMMENT '''意图识别置信度''',
  `entities` json DEFAULT NULL COMMENT '''提取的实体(JSON)''',
  `task_id` bigint unsigned DEFAULT NULL COMMENT '''关联的任务ID''',
  `task_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''任务类型''',
  `status` tinyint DEFAULT '1' COMMENT '''状态:1=处理中,2=成功,3=失败''',
  `error_msg` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''错误信息''',
  `create_time` datetime(3) NOT NULL COMMENT '''创建时间''',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_ai_agent_chat_history_session_id` (`session_id`) USING BTREE,
  KEY `idx_ai_agent_chat_history_user_id` (`user_id`) USING BTREE,
  KEY `idx_ai_agent_chat_history_task_id` (`task_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=143 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ai_agent_chat_history`
--

LOCK TABLES `ai_agent_chat_history` WRITE;
/*!40000 ALTER TABLE `ai_agent_chat_history` DISABLE KEYS */;
INSERT INTO `ai_agent_chat_history` VALUES (1,'d04b525e-d49c-41bd-a8e0-40752eca5a00',89,'user','检查 ip=180.76.231.65 的服务器压力','incident',0.70,'{\"checkType\": \"service\", \"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',1,'','2026-01-09 19:16:46.152'),(2,'889abebd-e1da-4cf6-bec8-f0a0594fb67a',89,'user','检查 ip=180.76.231.65 的服务器压力','incident',0.70,'{\"checkType\": \"service\", \"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',1,'','2026-01-09 19:19:01.533'),(3,'d15f101d-0a0a-4843-936b-21dcb1926f45',89,'user','检查 ip=180.76.231.65 的服务器压力','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',1,'','2026-01-09 19:24:18.448'),(4,'d15f101d-0a0a-4843-936b-21dcb1926f45',89,'assistant','# 服务器压力巡检报告\n\n**主机信息**\n- 主机名称: 百度云-ops\n- 主机IP: 180.76.231.65\n- 检查时间: 2026-01-09 19:24:18\n- 压力等级: **正常**\n\n## CPU信息\n- CPU核心数: 2\n- CPU使用率: 16.7%\n\n## 内存信息\n- 总内存: 1886 MB\n- 已使用: 1022 MB\n- 使用率: 54.19%\n\n## 磁盘信息\n- 根分区使用率: 51%\n```\n/dev/vda2        40G   20G   19G  51% /\n```\n\n## 系统负载\n- 负载平均值: 0.11, 0.04, 0.05\n\n## TOP 10 进程（按CPU使用率）\n| 用户 | PID | CPU% | MEM% | 命令 |\n|------|-----|------|------|------|\n| root | 94 | 11.7 | 0.0 | [kswapd0] |\n| root | 750164 | 3.1 | 1.8 | /opt/heyeAgent/bin/heyeAgent |\n| root | 4094944 | 1.2 | 0.4 | /lib/systemd/systemd --user |\n| root | 1 | 1.0 | 0.5 | /lib/systemd/systemd --switched-root --system --deserialize 31 nopti nospectre_v2 nospec_store_bypass_disable |\n| root | 1330 | 0.5 | 1.6 | /usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock |\n| root | 4083585 | 0.3 | 0.8 | /opt/agent/dodevops-agent |\n| root | 852 | 0.2 | 0.9 | /usr/bin/containerd |\n| systemd+ | 723 | 0.1 | 0.1 | /lib/systemd/systemd-networkd |\n| systemd+ | 725 | 0.1 | 0.4 | /lib/systemd/systemd-resolved |\n| fwupd-r+ | 920317 | 0.1 | 0.5 | /usr/bin/fwupdmgr refresh |\n\n','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',2,'','2026-01-09 19:24:25.613'),(5,'2cf25658-833a-4104-aaf1-de08cb5852cf',89,'user','检查 ip=180.76.231.65 的服务器压力','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',1,'','2026-01-09 19:25:06.014'),(6,'2cf25658-833a-4104-aaf1-de08cb5852cf',89,'assistant','# 服务器压力巡检报告\n\n**主机信息**\n- 主机名称: 百度云-ops\n- 主机IP: 180.76.231.65\n- 检查时间: 2026-01-09 19:25:06\n- 压力等级: **正常**\n\n## CPU信息\n- CPU核心数: 2\n- CPU使用率: 50.0%\n\n## 内存信息\n- 总内存: 1886 MB\n- 已使用: 1017 MB\n- 使用率: 53.92%\n\n## 磁盘信息\n- 根分区使用率: 51%\n```\n/dev/vda2        40G   20G   19G  51% /\n```\n\n## 系统负载\n- 负载平均值: 0.05, 0.03, 0.05\n\n## TOP 10 进程（按CPU使用率）\n| 用户 | PID | CPU% | MEM% | 命令 |\n|------|-----|------|------|------|\n| root | 94 | 11.7 | 0.0 | [kswapd0] |\n| root | 4095956 | 4.0 | 0.5 | sshd: root@notty |\n| root | 750164 | 3.1 | 1.8 | /opt/heyeAgent/bin/heyeAgent |\n| root | 1 | 1.0 | 0.5 | /lib/systemd/systemd --switched-root --system --deserialize 31 nopti nospectre_v2 nospec_store_bypass_disable |\n| root | 4095582 | 1.0 | 0.4 | /lib/systemd/systemd --user |\n| root | 1330 | 0.5 | 1.6 | /usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock |\n| root | 4095567 | 0.4 | 0.5 | sshd: root [priv] |\n| root | 4083585 | 0.3 | 0.8 | /opt/agent/dodevops-agent |\n| root | 852 | 0.2 | 0.9 | /usr/bin/containerd |\n| systemd+ | 723 | 0.1 | 0.1 | /lib/systemd/systemd-networkd |\n\n','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',2,'','2026-01-09 19:25:12.588'),(7,'867e3d9a-805a-4a5b-8493-36dbcfbacfc2',89,'user','帮我检查 ip=180.76.231.65 的服务器压力','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',1,'','2026-01-09 19:32:25.778'),(8,'867e3d9a-805a-4a5b-8493-36dbcfbacfc2',89,'assistant','# 服务器压力巡检报告\n\n**主机信息**\n- 主机名称: 百度云-ops\n- 主机IP: 180.76.231.65\n- 检查时间: 2026-01-09 19:32:26\n- 压力等级: **正常**\n\n## CPU信息\n- CPU核心数: 2\n- CPU使用率: 0.0%\n\n## 内存信息\n- 总内存: 1886 MB\n- 已使用: 1002 MB\n- 使用率: 53.13%\n\n## 磁盘信息\n- 根分区使用率: 51%\n```\n/dev/vda2        40G   20G   19G  51% /\n```\n\n## 系统负载\n- 负载平均值: 0.42, 0.10, 0.04\n\n## TOP 10 进程（按CPU使用率）\n| 用户 | PID | CPU% | MEM% | 命令 |\n|------|-----|------|------|------|\n| root | 94 | 11.7 | 0.0 | [kswapd0] |\n| root | 4096939 | 4.0 | 0.5 | sshd: root@notty |\n| root | 750164 | 3.1 | 1.8 | /opt/heyeAgent/bin/heyeAgent |\n| root | 4096550 | 1.2 | 0.5 | /lib/systemd/systemd --user |\n| root | 1 | 1.0 | 0.5 | /lib/systemd/systemd --switched-root --system --deserialize 31 nopti nospectre_v2 nospec_store_bypass_disable |\n| root | 4096886 | 1.0 | 0.4 | sshd: [accepted] |\n| root | 1330 | 0.5 | 1.6 | /usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock |\n| root | 4096848 | 0.5 | 0.4 | sshd: [accepted] |\n| root | 4083585 | 0.3 | 0.8 | /opt/agent/dodevops-agent |\n| root | 852 | 0.2 | 0.9 | /usr/bin/containerd |\n\n','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',2,'','2026-01-09 19:32:33.274'),(9,'f69033b5-8a94-4c03-8b0a-fa9ed056030e',89,'user','帮我检查 华为云-ops  这台机器压力如何','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"华为云-ops\", \"targetType\": \"host\"}',0,'',1,'','2026-01-09 21:47:00.398'),(10,'f69033b5-8a94-4c03-8b0a-fa9ed056030e',89,'assistant','# 服务器压力巡检报告\n\n**主机信息**\n- 主机名称: 华为云-ops\n- 主机IP: 139.9.205.38\n- 检查时间: 2026-01-09 21:47:00\n- 压力等级: **中等压力**\n\n## CPU信息\n- CPU核心数: 2\n- CPU使用率: 0.0%\n\n## 内存信息\n- 总内存: 1775 MB\n- 已使用: 1617 MB\n- 使用率: 91.10%\n\n## 磁盘信息\n- 根分区使用率: 43%\n```\n/dev/vda1        40G   17G   22G  43% /\noverlay          40G   17G   22G  43% /var/lib/docker/rootfs/overlayfs/40c5be161fc2ce582c3a2d1184fe03404db6dd46211fe3a79cf8931fc45f464c\noverlay          40G   17G   22G  43% /var/lib/docker/rootfs/overlayfs/c13a75becc246b4ece772df06bd03180f36529786290db2576d02cc7d7bda6c8\noverlay          40G   17G   22G  43% /var/lib/docker/rootfs/overlayfs/bf601832cdd47b0a419d55724d5bdb87e758e0061677058dba2defc6bcce1c11\noverlay          40G   17G   22G  43% /var/lib/docker/rootfs/overlayfs/9f487f841e39bdb68678c0f49d02f20618b497b71b5980631e62aaa58257e437\noverlay          40G   17G   22G  43% /var/lib/docker/rootfs/overlayfs/c13c84710575d8c918abfeda7a405b515d6b9042581017a6725dc10672608668\noverlay          40G   17G   22G  43% /var/lib/docker/rootfs/overlayfs/b60c6f281069d57bb2d0839e062446982829c387ace8d7595d1f559dcedff9a1\noverlay          40G   17G   22G  43% /var/lib/docker/rootfs/overlayfs/f3c0e0e68b86114b329da11998fd36528e02e705dae866c29033c225f387cb8e\noverlay          40G   17G   22G  43% /var/lib/docker/rootfs/overlayfs/69b46b1a5d278c65a2e736709bf22d1651142b1b996b5e0d149dbaff1706dd41\n```\n\n## 系统负载\n- 负载平均值: 0.05, 0.08, 0.08\n\n## TOP 10 进程（按CPU使用率）\n| 用户 | PID | CPU% | MEM% | 命令 |\n|------|-----|------|------|------|\n| root | 2257046 | 22.1 | 1.6 | /usr/bin/python3 /usr/lib/ubuntu-release-upgrader/check-new-release -q |\n| root | 2257429 | 4.0 | 0.5 | sshd: root@notty |\n| root | 2627393 | 0.9 | 0.4 | /usr/bin/containerd-shim-runc-v2 -namespace moby -id b60c6f281069d57bb2d0839e062446982829c387ace8d7595d1f559dcedff9a1 -address /run/containerd/containerd.sock |\n| dnsmasq | 1958583 | 0.3 | 22.1 | mysqld --default-authentication-plugin=caching_sha2_password --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci |\n| root | 2047742 | 0.3 | 0.4 | /opt/agent/dodevops-agent |\n| mysql | 1269514 | 0.3 | 20.9 | /usr/sbin/mysqld |\n| root | 379517 | 0.2 | 2.7 | /usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock |\n| root | 379334 | 0.1 | 1.8 | /usr/bin/containerd |\n| root | 2843811 | 0.1 | 2.0 | /app/devops |\n| root | 2609356 | 0.1 | 11.6 | /opt/netops/netops_linux/main |\n\n','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"华为云-ops\", \"targetType\": \"host\"}',0,'',2,'','2026-01-09 21:47:05.490'),(11,'2f248c63-fac5-462f-81ff-836681c49a4c',89,'user','帮我检查 华为云-ops 这台机器压力如何','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"华为云-ops\", \"targetType\": \"host\"}',0,'',1,'','2026-01-09 21:48:01.642'),(12,'2f248c63-fac5-462f-81ff-836681c49a4c',89,'assistant','# 服务器压力巡检报告\n\n**主机信息**\n- 主机名称: 华为云-ops\n- 主机IP: 139.9.205.38\n- 检查时间: 2026-01-09 21:48:02\n- 压力等级: **中等压力**\n\n## CPU信息\n- CPU核心数: 2\n- CPU使用率: 0.0%\n\n## 内存信息\n- 总内存: 1775 MB\n- 已使用: 1592 MB\n- 使用率: 89.69%\n\n## 磁盘信息\n- 根分区使用率: 43%\n```\n/dev/vda1        40G   17G   22G  43% /\noverlay          40G   17G   22G  43% /var/lib/docker/rootfs/overlayfs/40c5be161fc2ce582c3a2d1184fe03404db6dd46211fe3a79cf8931fc45f464c\noverlay          40G   17G   22G  43% /var/lib/docker/rootfs/overlayfs/c13a75becc246b4ece772df06bd03180f36529786290db2576d02cc7d7bda6c8\noverlay          40G   17G   22G  43% /var/lib/docker/rootfs/overlayfs/bf601832cdd47b0a419d55724d5bdb87e758e0061677058dba2defc6bcce1c11\noverlay          40G   17G   22G  43% /var/lib/docker/rootfs/overlayfs/9f487f841e39bdb68678c0f49d02f20618b497b71b5980631e62aaa58257e437\noverlay          40G   17G   22G  43% /var/lib/docker/rootfs/overlayfs/c13c84710575d8c918abfeda7a405b515d6b9042581017a6725dc10672608668\noverlay          40G   17G   22G  43% /var/lib/docker/rootfs/overlayfs/b60c6f281069d57bb2d0839e062446982829c387ace8d7595d1f559dcedff9a1\noverlay          40G   17G   22G  43% /var/lib/docker/rootfs/overlayfs/f3c0e0e68b86114b329da11998fd36528e02e705dae866c29033c225f387cb8e\noverlay          40G   17G   22G  43% /var/lib/docker/rootfs/overlayfs/69b46b1a5d278c65a2e736709bf22d1651142b1b996b5e0d149dbaff1706dd41\n```\n\n## 系统负载\n- 负载平均值: 0.16, 0.10, 0.08\n\n## TOP 10 进程（按CPU使用率）\n| 用户 | PID | CPU% | MEM% | 命令 |\n|------|-----|------|------|------|\n| root | 2258586 | 5.2 | 0.5 | sshd: root@notty |\n| root | 2627393 | 0.9 | 0.5 | /usr/bin/containerd-shim-runc-v2 -namespace moby -id b60c6f281069d57bb2d0839e062446982829c387ace8d7595d1f559dcedff9a1 -address /run/containerd/containerd.sock |\n| dnsmasq | 1958583 | 0.3 | 22.2 | mysqld --default-authentication-plugin=caching_sha2_password --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci |\n| root | 2047742 | 0.3 | 0.5 | /opt/agent/dodevops-agent |\n| mysql | 1269514 | 0.3 | 20.8 | /usr/sbin/mysqld |\n| root | 379517 | 0.2 | 3.0 | /usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock |\n| root | 379334 | 0.1 | 1.8 | /usr/bin/containerd |\n| root | 2843811 | 0.1 | 2.3 | /app/devops |\n| root | 2609356 | 0.1 | 11.6 | /opt/netops/netops_linux/main |\n| dnsmasq | 2626797 | 0.1 | 0.1 | redis-server *:6379 |\n\n','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"华为云-ops\", \"targetType\": \"host\"}',0,'',2,'','2026-01-09 21:48:05.711'),(13,'422b13f7-1254-4755-abbd-2b577a10d424',89,'user','帮我查询下  阿里云-ops 服务器的压力','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"阿里云-ops\", \"targetType\": \"host\"}',0,'',1,'','2026-01-09 22:05:03.025'),(14,'422b13f7-1254-4755-abbd-2b577a10d424',89,'assistant','# 服务器压力巡检报告\n\n**主机信息**\n- 主机名称: 阿里云-ops\n- 主机IP: 8.130.14.34\n- 检查时间: 2026-01-09 22:05:03\n- 压力等级: **中等压力**\n\n## CPU信息\n- CPU核心数: 2\n- CPU使用率: 4.8%\n\n## 内存信息\n- 总内存: 1613 MB\n- 已使用: 1386 MB\n- 使用率: 85.93%\n\n## 磁盘信息\n- 根分区使用率: 47%\n```\nefivarfs        256K  7.2K  244K   3% /sys/firmware/efi/efivars\n/dev/vda3        40G   18G   20G  47% /\n/dev/vda2       197M  6.2M  191M   4% /boot/efi\noverlay          40G   18G   20G  47% /var/lib/docker/overlay2/47bf984488f4fb9f9ddbc314c7139add6b96959daa1deaad119cf180af7d0a4b/merged\noverlay          40G   18G   20G  47% /var/lib/docker/overlay2/b369df5540ea748d3f4acda3aecc5e3b5c40a2722e791bb7846dac2133626079/merged\noverlay          40G   18G   20G  47% /var/lib/docker/overlay2/b5237f18229c37dd13f5f602aca464d85f255e6c35d40437db8e719e985d292b/merged\n```\n\n## 系统负载\n- 负载平均值: 0.03, 0.05, 0.07\n\n## TOP 10 进程（按CPU使用率）\n| 用户 | PID | CPU% | MEM% | 命令 |\n|------|-----|------|------|------|\n| root | 1452447 | 12.6 | 0.6 | sshd: root@notty |\n| root | 2461 | 1.1 | 1.4 | /usr/local/aegis/aegis_client/aegis_12_81/AliYunDunMonitor |\n| root | 1389472 | 0.7 | 0.7 | /opt/agent/dodevops-agent |\n| root | 1023 | 0.6 | 0.6 | /usr/local/cloudmonitor/bin/argusagent |\n| root | 855 | 0.6 | 1.9 | /opt/pushgateway |\n| root | 1062 | 0.4 | 1.5 | /usr/local/aegis/aegis_client/aegis_12_81/AliYunDun |\n| mysql | 1115 | 0.4 | 31.5 | /usr/sbin/mysqld |\n| root | 854 | 0.3 | 5.7 | /opt/prometheus/prometheus/prometheus --config.file=/opt/prometheus/prometheus/prometheus.yml --storage.tsdb.path=/opt/prometheus/prometheus/data --web.listen-address=0.0.0.0:9090 |\n| redis | 856 | 0.2 | 0.7 | /usr/bin/redis-server 0.0.0.0:6379 |\n| root | 1186 | 0.2 | 2.9 | /usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock |\n\n','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"阿里云-ops\", \"targetType\": \"host\"}',0,'',2,'','2026-01-09 22:05:13.034'),(15,'1bb6bc5f-2bc2-45f5-9812-f99cbf181c11',89,'user','帮我查询下 阿里云-ops 服务器的压力','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"阿里云-ops\", \"targetType\": \"host\"}',0,'',1,'','2026-01-09 22:06:36.120'),(16,'1bb6bc5f-2bc2-45f5-9812-f99cbf181c11',89,'assistant','# 服务器压力巡检报告\n\n**主机信息**\n- 主机名称: 阿里云-ops\n- 主机IP: 8.130.14.34\n- 检查时间: 2026-01-09 22:06:36\n- 压力等级: **中等压力**\n\n## CPU信息\n- CPU核心数: 2\n- CPU使用率: 4.8%\n\n## 内存信息\n- 总内存: 1613 MB\n- 已使用: 1407 MB\n- 使用率: 87.23%\n\n## 磁盘信息\n- 根分区使用率: 47%\n```\nefivarfs        256K  7.2K  244K   3% /sys/firmware/efi/efivars\n/dev/vda3        40G   18G   20G  47% /\n/dev/vda2       197M  6.2M  191M   4% /boot/efi\noverlay          40G   18G   20G  47% /var/lib/docker/overlay2/47bf984488f4fb9f9ddbc314c7139add6b96959daa1deaad119cf180af7d0a4b/merged\noverlay          40G   18G   20G  47% /var/lib/docker/overlay2/b369df5540ea748d3f4acda3aecc5e3b5c40a2722e791bb7846dac2133626079/merged\noverlay          40G   18G   20G  47% /var/lib/docker/overlay2/b5237f18229c37dd13f5f602aca464d85f255e6c35d40437db8e719e985d292b/merged\n```\n\n## 系统负载\n- 负载平均值: 0.41, 0.15, 0.10\n\n## TOP 10 进程（按CPU使用率）\n| 用户 | PID | CPU% | MEM% | 命令 |\n|------|-----|------|------|------|\n| root | 1453368 | 100 | 0.2 | ps aux --sort=-%cpu |\n| root | 1453318 | 13.0 | 0.6 | sshd: root@notty |\n| root | 2461 | 1.1 | 1.4 | /usr/local/aegis/aegis_client/aegis_12_81/AliYunDunMonitor |\n| root | 1389472 | 0.7 | 0.6 | /opt/agent/dodevops-agent |\n| root | 1023 | 0.6 | 0.6 | /usr/local/cloudmonitor/bin/argusagent |\n| root | 855 | 0.6 | 2.0 | /opt/pushgateway |\n| root | 1062 | 0.4 | 1.5 | /usr/local/aegis/aegis_client/aegis_12_81/AliYunDun |\n| mysql | 1115 | 0.4 | 31.5 | /usr/sbin/mysqld |\n| root | 854 | 0.3 | 5.7 | /opt/prometheus/prometheus/prometheus --config.file=/opt/prometheus/prometheus/prometheus.yml --storage.tsdb.path=/opt/prometheus/prometheus/data --web.listen-address=0.0.0.0:9090 |\n| redis | 856 | 0.2 | 0.7 | /usr/bin/redis-server 0.0.0.0:6379 |\n\n','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"阿里云-ops\", \"targetType\": \"host\"}',0,'',2,'','2026-01-09 22:06:46.080'),(17,'2e4b862f-82b0-4cf4-8781-cff149a574b1',89,'user','帮我检查下 华为云-ops  这个服务器压力怎么样','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"华为云-ops\", \"targetType\": \"host\"}',0,'',1,'','2026-01-09 22:15:59.122'),(18,'2e4b862f-82b0-4cf4-8781-cff149a574b1',89,'assistant','# 服务器压力巡检报告\n\n**主机信息**\n- 主机名称: 华为云-ops\n- 主机IP: 139.9.205.38\n- 检查时间: 2026-01-09 22:15:59\n- 压力等级: **中等压力**\n\n## CPU信息\n- CPU核心数: 2\n- CPU使用率: 0.0%\n\n## 内存信息\n- 总内存: 1775 MB\n- 已使用: 1594 MB\n- 使用率: 89.80%\n\n## 磁盘信息\n- 根分区使用率: 43%\n```\n/dev/vda1        40G   17G   22G  43% /\noverlay          40G   17G   22G  43% /var/lib/docker/rootfs/overlayfs/40c5be161fc2ce582c3a2d1184fe03404db6dd46211fe3a79cf8931fc45f464c\noverlay          40G   17G   22G  43% /var/lib/docker/rootfs/overlayfs/c13a75becc246b4ece772df06bd03180f36529786290db2576d02cc7d7bda6c8\noverlay          40G   17G   22G  43% /var/lib/docker/rootfs/overlayfs/bf601832cdd47b0a419d55724d5bdb87e758e0061677058dba2defc6bcce1c11\noverlay          40G   17G   22G  43% /var/lib/docker/rootfs/overlayfs/9f487f841e39bdb68678c0f49d02f20618b497b71b5980631e62aaa58257e437\noverlay          40G   17G   22G  43% /var/lib/docker/rootfs/overlayfs/c13c84710575d8c918abfeda7a405b515d6b9042581017a6725dc10672608668\noverlay          40G   17G   22G  43% /var/lib/docker/rootfs/overlayfs/b60c6f281069d57bb2d0839e062446982829c387ace8d7595d1f559dcedff9a1\noverlay          40G   17G   22G  43% /var/lib/docker/rootfs/overlayfs/f3c0e0e68b86114b329da11998fd36528e02e705dae866c29033c225f387cb8e\noverlay          40G   17G   22G  43% /var/lib/docker/rootfs/overlayfs/69b46b1a5d278c65a2e736709bf22d1651142b1b996b5e0d149dbaff1706dd41\n```\n\n## 系统负载\n- 负载平均值: 0.09, 0.08, 0.08\n\n## TOP 10 进程（按CPU使用率）\n| 用户 | PID | CPU% | MEM% | 命令 |\n|------|-----|------|------|------|\n| root | 2277913 | 5.0 | 0.5 | sshd: root@notty |\n| root | 2627393 | 0.9 | 0.4 | /usr/bin/containerd-shim-runc-v2 -namespace moby -id b60c6f281069d57bb2d0839e062446982829c387ace8d7595d1f559dcedff9a1 -address /run/containerd/containerd.sock |\n| dnsmasq | 1958583 | 0.3 | 22.3 | mysqld --default-authentication-plugin=caching_sha2_password --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci |\n| root | 2047742 | 0.3 | 0.6 | /opt/agent/dodevops-agent |\n| mysql | 1269514 | 0.3 | 21.2 | /usr/sbin/mysqld |\n| root | 379517 | 0.2 | 2.8 | /usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock |\n| root | 379334 | 0.1 | 1.7 | /usr/bin/containerd |\n| root | 2843811 | 0.1 | 2.2 | /app/devops |\n| root | 2609356 | 0.1 | 11.9 | /opt/netops/netops_linux/main |\n| dnsmasq | 2626797 | 0.1 | 0.2 | redis-server *:6379 |\n\n','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"华为云-ops\", \"targetType\": \"host\"}',0,'',2,'','2026-01-09 22:16:03.255'),(19,'a0c21ecc-787b-4743-a0b5-7efe7933a069',89,'user','帮忙查看 180.76.231.65 服务器的压力情况','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',1,'','2026-01-09 22:24:35.702'),(20,'a0c21ecc-787b-4743-a0b5-7efe7933a069',89,'assistant','# 服务器压力巡检报告\n\n**主机信息**\n- 主机名称: 百度云-ops\n- 主机IP: 180.76.231.65\n- 检查时间: 2026-01-09 22:24:36\n- 压力等级: **正常**\n\n## CPU信息\n- CPU核心数: 2\n- CPU使用率: 0.0%\n\n## 内存信息\n- 总内存: 1886 MB\n- 已使用: 1446 MB\n- 使用率: 76.67%\n\n## 磁盘信息\n- 根分区使用率: 54%\n```\n/dev/vda2        40G   21G   18G  54% /\noverlay          40G   21G   18G  54% /var/lib/docker/overlay2/36d3c5e2aa5b2e721d81c27c2923c1645e788af8c62af73ca78c562acb1dc6c6/merged\noverlay          40G   21G   18G  54% /var/lib/docker/overlay2/3a19d7d7ec48824a825359d54df4f9c59bde749d91c29a4c2bf7c39d5bedb135/merged\noverlay          40G   21G   18G  54% /var/lib/docker/overlay2/273c4fcfacbd0de3a76036142e72f3838211b24fc5753e7a39494964479738ce/merged\n```\n\n## 系统负载\n- 负载平均值: 0.13, 0.18, 0.13\n\n## TOP 10 进程（按CPU使用率）\n| 用户 | PID | CPU% | MEM% | 命令 |\n|------|-----|------|------|------|\n| root | 94 | 11.7 | 0.0 | [kswapd0] |\n| root | 4123568 | 4.0 | 0.5 | sshd: root@notty |\n| root | 750164 | 3.1 | 1.6 | /opt/heyeAgent/bin/heyeAgent |\n| root | 1 | 1.0 | 0.5 | /lib/systemd/systemd --switched-root --system --deserialize 31 nopti nospectre_v2 nospec_store_bypass_disable |\n| root | 4122377 | 1.0 | 2.1 | /usr/lib/snapd/snapd |\n| root | 4123165 | 0.8 | 0.4 | /lib/systemd/systemd --user |\n| root | 4083585 | 0.3 | 0.5 | /opt/agent/dodevops-agent |\n| root | 4101046 | 0.3 | 1.7 | /usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock |\n| root | 852 | 0.2 | 1.1 | /usr/bin/containerd |\n| systemd+ | 723 | 0.1 | 0.2 | /lib/systemd/systemd-networkd |\n\n','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',2,'','2026-01-09 22:24:43.976'),(21,'b4c12f12-da77-4efa-9438-ca14747d93df',89,'user','检查 180.76.231.65 服务器压力 ip=180.76.231.65','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',1,'','2026-01-09 22:34:20.694'),(22,'b4c12f12-da77-4efa-9438-ca14747d93df',89,'assistant','# 服务器压力巡检报告\n\n**主机信息**\n- 主机名称: 百度云-ops\n- 主机IP: 180.76.231.65\n- 检查时间: 2026-01-09 22:34:21\n- 压力等级: **正常**\n\n## CPU信息\n- CPU核心数: 2\n- CPU使用率: 21.2%\n\n## 内存信息\n- 总内存: 1886 MB\n- 已使用: 1463 MB\n- 使用率: 77.57%\n\n## 磁盘信息\n- 根分区使用率: 54%\n```\n/dev/vda2        40G   21G   18G  54% /\noverlay          40G   21G   18G  54% /var/lib/docker/overlay2/36d3c5e2aa5b2e721d81c27c2923c1645e788af8c62af73ca78c562acb1dc6c6/merged\noverlay          40G   21G   18G  54% /var/lib/docker/overlay2/3a19d7d7ec48824a825359d54df4f9c59bde749d91c29a4c2bf7c39d5bedb135/merged\noverlay          40G   21G   18G  54% /var/lib/docker/overlay2/273c4fcfacbd0de3a76036142e72f3838211b24fc5753e7a39494964479738ce/merged\n```\n\n## 系统负载\n- 负载平均值: 0.27, 0.14, 0.11\n\n## TOP 10 进程（按CPU使用率）\n| 用户 | PID | CPU% | MEM% | 命令 |\n|------|-----|------|------|------|\n| root | 94 | 11.7 | 0.0 | [kswapd0] |\n| root | 4126009 | 4.0 | 0.5 | sshd: root@notty |\n| root | 750164 | 3.1 | 1.6 | /opt/heyeAgent/bin/heyeAgent |\n| root | 1 | 1.0 | 0.5 | /lib/systemd/systemd --switched-root --system --deserialize 31 nopti nospectre_v2 nospec_store_bypass_disable |\n| root | 4125611 | 1.0 | 0.5 | /lib/systemd/systemd --user |\n| root | 4083585 | 0.3 | 0.5 | /opt/agent/dodevops-agent |\n| root | 852 | 0.2 | 1.0 | /usr/bin/containerd |\n| root | 4101046 | 0.2 | 1.7 | /usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock |\n| root | 4122377 | 0.2 | 1.8 | /usr/lib/snapd/snapd |\n| systemd+ | 723 | 0.1 | 0.2 | /lib/systemd/systemd-networkd |\n\n','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',2,'','2026-01-09 22:34:28.842'),(23,'a0dc851a-7311-4e2b-96e9-bd10739f5ff4',89,'user','检查 180.76.231.65 服务器压力 ip=180.76.231.65','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',1,'','2026-01-09 22:36:07.532'),(24,'a0dc851a-7311-4e2b-96e9-bd10739f5ff4',89,'assistant','# 服务器压力巡检报告\n\n**主机信息**\n- 主机名称: 百度云-ops\n- 主机IP: 180.76.231.65\n- 检查时间: 2026-01-09 22:36:08\n- 压力等级: **正常**\n\n## CPU信息\n- CPU核心数: 2\n- CPU使用率: 3.2%\n\n## 内存信息\n- 总内存: 1886 MB\n- 已使用: 1464 MB\n- 使用率: 77.62%\n\n## 磁盘信息\n- 根分区使用率: 54%\n```\n/dev/vda2        40G   21G   18G  54% /\noverlay          40G   21G   18G  54% /var/lib/docker/overlay2/36d3c5e2aa5b2e721d81c27c2923c1645e788af8c62af73ca78c562acb1dc6c6/merged\noverlay          40G   21G   18G  54% /var/lib/docker/overlay2/3a19d7d7ec48824a825359d54df4f9c59bde749d91c29a4c2bf7c39d5bedb135/merged\noverlay          40G   21G   18G  54% /var/lib/docker/overlay2/273c4fcfacbd0de3a76036142e72f3838211b24fc5753e7a39494964479738ce/merged\n```\n\n## 系统负载\n- 负载平均值: 0.27, 0.14, 0.11\n\n## TOP 10 进程（按CPU使用率）\n| 用户 | PID | CPU% | MEM% | 命令 |\n|------|-----|------|------|------|\n| root | 94 | 11.7 | 0.0 | [kswapd0] |\n| root | 4126894 | 4.0 | 0.5 | sshd: root@notty |\n| root | 750164 | 3.1 | 1.6 | /opt/heyeAgent/bin/heyeAgent |\n| root | 1 | 1.0 | 0.5 | /lib/systemd/systemd --switched-root --system --deserialize 31 nopti nospectre_v2 nospec_store_bypass_disable |\n| root | 4126511 | 1.0 | 0.4 | /lib/systemd/systemd --user |\n| root | 4083585 | 0.3 | 0.5 | /opt/agent/dodevops-agent |\n| root | 852 | 0.2 | 1.0 | /usr/bin/containerd |\n| root | 4101046 | 0.2 | 1.7 | /usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock |\n| root | 4122377 | 0.2 | 1.8 | /usr/lib/snapd/snapd |\n| systemd+ | 723 | 0.1 | 0.2 | /lib/systemd/systemd-networkd |\n\n','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',2,'','2026-01-09 22:36:15.772'),(25,'d1a02db9-5fc7-4ccf-b73e-bbb03b147866',89,'user','检查 180.76.231.65 服务器压力 ip=180.76.231.65','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',1,'','2026-01-09 22:44:16.227'),(26,'d1a02db9-5fc7-4ccf-b73e-bbb03b147866',89,'assistant','# 服务器压力巡检报告\n\n**主机信息**\n- 主机名称: 百度云-ops\n- 主机IP: 180.76.231.65\n- 检查时间: 2026-01-09 22:44:16\n- 压力等级: **正常**\n\n## CPU信息\n- CPU核心数: 2\n- CPU使用率: 3.2%\n\n## 内存信息\n- 总内存: 1886 MB\n- 已使用: 1452 MB\n- 使用率: 76.99%\n\n## 磁盘信息\n- 根分区使用率: 54%\n```\n/dev/vda2        40G   21G   18G  54% /\noverlay          40G   21G   18G  54% /var/lib/docker/overlay2/36d3c5e2aa5b2e721d81c27c2923c1645e788af8c62af73ca78c562acb1dc6c6/merged\noverlay          40G   21G   18G  54% /var/lib/docker/overlay2/3a19d7d7ec48824a825359d54df4f9c59bde749d91c29a4c2bf7c39d5bedb135/merged\noverlay          40G   21G   18G  54% /var/lib/docker/overlay2/273c4fcfacbd0de3a76036142e72f3838211b24fc5753e7a39494964479738ce/merged\n```\n\n## 系统负载\n- 负载平均值: 0.00, 0.06, 0.08\n\n## TOP 10 进程（按CPU使用率）\n| 用户 | PID | CPU% | MEM% | 命令 |\n|------|-----|------|------|------|\n| root | 94 | 11.7 | 0.0 | [kswapd0] |\n| root | 4129054 | 4.0 | 0.5 | sshd: root@notty |\n| root | 750164 | 3.1 | 1.6 | /opt/heyeAgent/bin/heyeAgent |\n| root | 1 | 1.0 | 0.5 | /lib/systemd/systemd --switched-root --system --deserialize 31 nopti nospectre_v2 nospec_store_bypass_disable |\n| root | 4128653 | 1.0 | 0.5 | /lib/systemd/systemd --user |\n| root | 4083585 | 0.3 | 0.5 | /opt/agent/dodevops-agent |\n| root | 4128644 | 0.3 | 0.5 | sshd: root [priv] |\n| root | 852 | 0.2 | 1.0 | /usr/bin/containerd |\n| root | 4101046 | 0.2 | 1.7 | /usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock |\n| systemd+ | 723 | 0.1 | 0.2 | /lib/systemd/systemd-networkd |\n\n','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',2,'','2026-01-09 22:44:24.422'),(27,'d39783c4-ee5f-4117-bd7b-f3afe52017e4',89,'user','检查 180.76.231.65 服务器压力 ip=180.76.231.65','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',1,'','2026-01-09 22:58:54.414'),(28,'d39783c4-ee5f-4117-bd7b-f3afe52017e4',89,'assistant','# 服务器压力巡检报告\n\n**主机信息**\n- 主机名称: 百度云-ops\n- 主机IP: 180.76.231.65\n- 检查时间: 2026-01-09 22:58:54\n- 压力等级: **正常**\n\n## CPU信息\n- CPU核心数: 2\n- CPU使用率: 0.0%\n\n## 内存信息\n- 总内存: 1886 MB\n- 已使用: 1443 MB\n- 使用率: 76.51%\n\n## 磁盘信息\n- 根分区使用率: 54%\n```\n/dev/vda2        40G   21G   18G  54% /\noverlay          40G   21G   18G  54% /var/lib/docker/overlay2/36d3c5e2aa5b2e721d81c27c2923c1645e788af8c62af73ca78c562acb1dc6c6/merged\noverlay          40G   21G   18G  54% /var/lib/docker/overlay2/3a19d7d7ec48824a825359d54df4f9c59bde749d91c29a4c2bf7c39d5bedb135/merged\noverlay          40G   21G   18G  54% /var/lib/docker/overlay2/273c4fcfacbd0de3a76036142e72f3838211b24fc5753e7a39494964479738ce/merged\n```\n\n## 系统负载\n- 负载平均值: 0.13, 0.08, 0.06\n\n## TOP 10 进程（按CPU使用率）\n| 用户 | PID | CPU% | MEM% | 命令 |\n|------|-----|------|------|------|\n| root | 94 | 11.7 | 0.0 | [kswapd0] |\n| root | 750164 | 3.1 | 1.6 | /opt/heyeAgent/bin/heyeAgent |\n| root | 4132006 | 1.2 | 0.5 | /lib/systemd/systemd --user |\n| root | 1 | 1.0 | 0.5 | /lib/systemd/systemd --switched-root --system --deserialize 31 nopti nospectre_v2 nospec_store_bypass_disable |\n| root | 4083585 | 0.3 | 0.5 | /opt/agent/dodevops-agent |\n| root | 852 | 0.2 | 1.0 | /usr/bin/containerd |\n| root | 4101046 | 0.2 | 1.8 | /usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock |\n| systemd+ | 723 | 0.1 | 0.2 | /lib/systemd/systemd-networkd |\n| systemd+ | 725 | 0.1 | 0.4 | /lib/systemd/systemd-resolved |\n| fwupd-r+ | 920317 | 0.1 | 0.5 | /usr/bin/fwupdmgr refresh |\n\n','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',2,'','2026-01-09 22:59:02.482'),(29,'7469c0e3-a402-4aec-acee-d565df0641f9',89,'user','帮我 检查 这个服务器的压力 8.130.14.34 ip=8.130.14.34','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"8.130.14.34\", \"targetType\": \"host\"}',0,'',1,'','2026-01-09 23:11:47.218'),(30,'7469c0e3-a402-4aec-acee-d565df0641f9',89,'assistant','# 服务器压力巡检报告\n\n**主机信息**\n- 主机名称: 阿里云-ops\n- 主机IP: 8.130.14.34\n- 检查时间: 2026-01-09 23:11:47\n- 压力等级: **中等压力**\n\n## CPU信息\n- CPU核心数: 2\n- CPU使用率: 9.1%\n\n## 内存信息\n- 总内存: 1613 MB\n- 已使用: 1383 MB\n- 使用率: 85.74%\n\n## 磁盘信息\n- 根分区使用率: 47%\n```\nefivarfs        256K  7.2K  244K   3% /sys/firmware/efi/efivars\n/dev/vda3        40G   18G   20G  47% /\n/dev/vda2       197M  6.2M  191M   4% /boot/efi\noverlay          40G   18G   20G  47% /var/lib/docker/overlay2/47bf984488f4fb9f9ddbc314c7139add6b96959daa1deaad119cf180af7d0a4b/merged\noverlay          40G   18G   20G  47% /var/lib/docker/overlay2/b369df5540ea748d3f4acda3aecc5e3b5c40a2722e791bb7846dac2133626079/merged\noverlay          40G   18G   20G  47% /var/lib/docker/overlay2/b5237f18229c37dd13f5f602aca464d85f255e6c35d40437db8e719e985d292b/merged\n```\n\n## 系统负载\n- 负载平均值: 0.00, 0.04, 0.08\n\n## TOP 10 进程（按CPU使用率）\n| 用户 | PID | CPU% | MEM% | 命令 |\n|------|-----|------|------|------|\n| root | 1465883 | 100 | 0.2 | ps aux --sort=-%cpu |\n| root | 1465833 | 12.8 | 0.6 | sshd: root@notty |\n| root | 2461 | 1.2 | 1.4 | /usr/local/aegis/aegis_client/aegis_12_81/AliYunDunMonitor |\n| root | 1389472 | 0.7 | 0.7 | /opt/agent/dodevops-agent |\n| root | 1023 | 0.6 | 0.7 | /usr/local/cloudmonitor/bin/argusagent |\n| root | 855 | 0.6 | 2.0 | /opt/pushgateway |\n| root | 1062 | 0.4 | 1.5 | /usr/local/aegis/aegis_client/aegis_12_81/AliYunDun |\n| mysql | 1115 | 0.4 | 31.8 | /usr/sbin/mysqld |\n| root | 854 | 0.3 | 5.2 | /opt/prometheus/prometheus/prometheus --config.file=/opt/prometheus/prometheus/prometheus.yml --storage.tsdb.path=/opt/prometheus/prometheus/data --web.listen-address=0.0.0.0:9090 |\n| redis | 856 | 0.2 | 0.7 | /usr/bin/redis-server 0.0.0.0:6379 |\n\n','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"8.130.14.34\", \"targetType\": \"host\"}',0,'',2,'','2026-01-09 23:11:57.658'),(31,'75aa6d6f-42d7-4bfc-b4f7-d87ce7f1662c',89,'user','180.76.231.65  帮忙检查下这个服务器的压力 ip=180.76.231.65','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',1,'','2026-01-09 23:31:43.933'),(32,'75aa6d6f-42d7-4bfc-b4f7-d87ce7f1662c',89,'assistant','# 服务器压力巡检报告\n\n**主机信息**\n- 主机名称: 百度云-ops\n- 主机IP: 180.76.231.65\n- 检查时间: 2026-01-09 23:31:45\n- 压力等级: **正常**\n\n## CPU信息\n- CPU核心数: 2\n- CPU使用率: 3.2%\n\n## 内存信息\n- 总内存: 1886 MB\n- 已使用: 1450 MB\n- 使用率: 76.88%\n\n## 磁盘信息\n- 根分区使用率: 54%\n```\n/dev/vda2        40G   21G   18G  54% /\noverlay          40G   21G   18G  54% /var/lib/docker/overlay2/36d3c5e2aa5b2e721d81c27c2923c1645e788af8c62af73ca78c562acb1dc6c6/merged\noverlay          40G   21G   18G  54% /var/lib/docker/overlay2/3a19d7d7ec48824a825359d54df4f9c59bde749d91c29a4c2bf7c39d5bedb135/merged\noverlay          40G   21G   18G  54% /var/lib/docker/overlay2/273c4fcfacbd0de3a76036142e72f3838211b24fc5753e7a39494964479738ce/merged\n```\n\n## 系统负载\n- 负载平均值: 0.00, 0.06, 0.03\n\n## TOP 10 进程（按CPU使用率）\n| 用户 | PID | CPU% | MEM% | 命令 |\n|------|-----|------|------|------|\n| root | 94 | 11.7 | 0.0 | [kswapd0] |\n| root | 750164 | 3.1 | 1.5 | /opt/heyeAgent/bin/heyeAgent |\n| root | 1 | 1.0 | 0.5 | /lib/systemd/systemd --switched-root --system --deserialize 31 nopti nospectre_v2 nospec_store_bypass_disable |\n| root | 4138275 | 1.0 | 0.5 | /lib/systemd/systemd --user |\n| root | 4083585 | 0.3 | 0.5 | /opt/agent/dodevops-agent |\n| root | 852 | 0.2 | 1.0 | /usr/bin/containerd |\n| root | 4101046 | 0.2 | 1.9 | /usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock |\n| systemd+ | 723 | 0.1 | 0.2 | /lib/systemd/systemd-networkd |\n| systemd+ | 725 | 0.1 | 0.4 | /lib/systemd/systemd-resolved |\n| fwupd-r+ | 920317 | 0.1 | 0.5 | /usr/bin/fwupdmgr refresh |\n\n','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',2,'','2026-01-09 23:32:10.398'),(33,'75aa6d6f-42d7-4bfc-b4f7-d87ce7f1662c',89,'user','180.76.231.65 帮忙检查下这个服务器的压力 ip=180.76.231.65','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',1,'','2026-01-09 23:32:58.950'),(34,'75aa6d6f-42d7-4bfc-b4f7-d87ce7f1662c',89,'assistant','# 服务器压力巡检报告\n\n**主机信息**\n- 主机名称: 百度云-ops\n- 主机IP: 180.76.231.65\n- 检查时间: 2026-01-09 23:33:01\n- 压力等级: **正常**\n\n## CPU信息\n- CPU核心数: 2\n- CPU使用率: 9.4%\n\n## 内存信息\n- 总内存: 1886 MB\n- 已使用: 1439 MB\n- 使用率: 76.30%\n\n## 磁盘信息\n- 根分区使用率: 54%\n```\n/dev/vda2        40G   21G   18G  54% /\noverlay          40G   21G   18G  54% /var/lib/docker/overlay2/36d3c5e2aa5b2e721d81c27c2923c1645e788af8c62af73ca78c562acb1dc6c6/merged\noverlay          40G   21G   18G  54% /var/lib/docker/overlay2/3a19d7d7ec48824a825359d54df4f9c59bde749d91c29a4c2bf7c39d5bedb135/merged\noverlay          40G   21G   18G  54% /var/lib/docker/overlay2/273c4fcfacbd0de3a76036142e72f3838211b24fc5753e7a39494964479738ce/merged\n```\n\n## 系统负载\n- 负载平均值: 0.02, 0.06, 0.03\n\n## TOP 10 进程（按CPU使用率）\n| 用户 | PID | CPU% | MEM% | 命令 |\n|------|-----|------|------|------|\n| root | 94 | 11.7 | 0.0 | [kswapd0] |\n| root | 750164 | 3.1 | 1.5 | /opt/heyeAgent/bin/heyeAgent |\n| root | 1 | 1.0 | 0.5 | /lib/systemd/systemd --switched-root --system --deserialize 31 nopti nospectre_v2 nospec_store_bypass_disable |\n| root | 4139027 | 1.0 | 0.4 | /lib/systemd/systemd --user |\n| root | 4083585 | 0.3 | 0.6 | /opt/agent/dodevops-agent |\n| root | 852 | 0.2 | 1.0 | /usr/bin/containerd |\n| root | 4101046 | 0.2 | 1.9 | /usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock |\n| systemd+ | 723 | 0.1 | 0.2 | /lib/systemd/systemd-networkd |\n| systemd+ | 725 | 0.1 | 0.4 | /lib/systemd/systemd-resolved |\n| fwupd-r+ | 920317 | 0.1 | 0.5 | /usr/bin/fwupdmgr refresh |\n\n','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',2,'','2026-01-09 23:33:09.524'),(35,'75aa6d6f-42d7-4bfc-b4f7-d87ce7f1662c',89,'user','180.76.231.65 帮忙检查下这个服务器的压力 ip=180.76.231.65','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',1,'','2026-01-09 23:33:59.759'),(36,'75aa6d6f-42d7-4bfc-b4f7-d87ce7f1662c',89,'assistant','# 服务器压力巡检报告\n\n**主机信息**\n- 主机名称: 百度云-ops\n- 主机IP: 180.76.231.65\n- 检查时间: 2026-01-09 23:34:00\n- 压力等级: **正常**\n\n## CPU信息\n- CPU核心数: 2\n- CPU使用率: 0.0%\n\n## 内存信息\n- 总内存: 1886 MB\n- 已使用: 1448 MB\n- 使用率: 76.78%\n\n## 磁盘信息\n- 根分区使用率: 54%\n```\n/dev/vda2        40G   21G   18G  54% /\noverlay          40G   21G   18G  54% /var/lib/docker/overlay2/36d3c5e2aa5b2e721d81c27c2923c1645e788af8c62af73ca78c562acb1dc6c6/merged\noverlay          40G   21G   18G  54% /var/lib/docker/overlay2/3a19d7d7ec48824a825359d54df4f9c59bde749d91c29a4c2bf7c39d5bedb135/merged\noverlay          40G   21G   18G  54% /var/lib/docker/overlay2/273c4fcfacbd0de3a76036142e72f3838211b24fc5753e7a39494964479738ce/merged\n```\n\n## 系统负载\n- 负载平均值: 0.01, 0.05, 0.03\n\n## TOP 10 进程（按CPU使用率）\n| 用户 | PID | CPU% | MEM% | 命令 |\n|------|-----|------|------|------|\n| root | 94 | 11.7 | 0.0 | [kswapd0] |\n| root | 4140179 | 4.0 | 0.5 | sshd: root@notty |\n| root | 750164 | 3.1 | 1.6 | /opt/heyeAgent/bin/heyeAgent |\n| root | 1 | 1.0 | 0.5 | /lib/systemd/systemd --switched-root --system --deserialize 31 nopti nospectre_v2 nospec_store_bypass_disable |\n| root | 4139746 | 0.8 | 0.4 | /lib/systemd/systemd --user |\n| root | 4083585 | 0.3 | 0.6 | /opt/agent/dodevops-agent |\n| root | 852 | 0.2 | 1.0 | /usr/bin/containerd |\n| root | 4101046 | 0.2 | 1.8 | /usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock |\n| systemd+ | 723 | 0.1 | 0.2 | /lib/systemd/systemd-networkd |\n| systemd+ | 725 | 0.1 | 0.4 | /lib/systemd/systemd-resolved |\n\n','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',2,'','2026-01-09 23:34:09.034'),(37,'afff8e56-7585-452d-8c66-fb08d097625e',89,'user','帮我 检查 这个服务器的压力 8.130.14.34 ip=8.130.14.34','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"8.130.14.34\", \"targetType\": \"host\"}',0,'',1,'','2026-01-09 23:48:22.039'),(38,'afff8e56-7585-452d-8c66-fb08d097625e',89,'assistant','# 服务器压力巡检报告\n\n**主机信息**\n- 主机名称: 阿里云-ops\n- 主机IP: 8.130.14.34\n- 检查时间: 2026-01-09 23:48:22\n- 压力等级: **中等压力**\n\n## CPU信息\n- CPU核心数: 2\n- CPU使用率: 4.2%\n\n## 内存信息\n- 总内存: 1613 MB\n- 已使用: 1393 MB\n- 使用率: 86.36%\n\n## 磁盘信息\n- 根分区使用率: 47%\n```\nefivarfs        256K  7.2K  244K   3% /sys/firmware/efi/efivars\n/dev/vda3        40G   18G   20G  47% /\n/dev/vda2       197M  6.2M  191M   4% /boot/efi\noverlay          40G   18G   20G  47% /var/lib/docker/overlay2/47bf984488f4fb9f9ddbc314c7139add6b96959daa1deaad119cf180af7d0a4b/merged\noverlay          40G   18G   20G  47% /var/lib/docker/overlay2/b369df5540ea748d3f4acda3aecc5e3b5c40a2722e791bb7846dac2133626079/merged\noverlay          40G   18G   20G  47% /var/lib/docker/overlay2/b5237f18229c37dd13f5f602aca464d85f255e6c35d40437db8e719e985d292b/merged\n```\n\n## 系统负载\n- 负载平均值: 0.10, 0.10, 0.09\n\n## TOP 10 进程（按CPU使用率）\n| 用户 | PID | CPU% | MEM% | 命令 |\n|------|-----|------|------|------|\n| root | 1473151 | 13.8 | 0.6 | sshd: root@notty |\n| root | 2461 | 1.2 | 1.3 | /usr/local/aegis/aegis_client/aegis_12_81/AliYunDunMonitor |\n| root | 1389472 | 0.7 | 0.7 | /opt/agent/dodevops-agent |\n| root | 1023 | 0.6 | 0.6 | /usr/local/cloudmonitor/bin/argusagent |\n| root | 855 | 0.6 | 1.9 | /opt/pushgateway |\n| root | 1062 | 0.4 | 1.5 | /usr/local/aegis/aegis_client/aegis_12_81/AliYunDun |\n| mysql | 1115 | 0.4 | 31.8 | /usr/sbin/mysqld |\n| root | 854 | 0.3 | 5.2 | /opt/prometheus/prometheus/prometheus --config.file=/opt/prometheus/prometheus/prometheus.yml --storage.tsdb.path=/opt/prometheus/prometheus/data --web.listen-address=0.0.0.0:9090 |\n| redis | 856 | 0.2 | 0.7 | /usr/bin/redis-server 0.0.0.0:6379 |\n| root | 1186 | 0.2 | 2.9 | /usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock |\n\n','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"8.130.14.34\", \"targetType\": \"host\"}',0,'',2,'','2026-01-09 23:48:31.976'),(39,'afff8e56-7585-452d-8c66-fb08d097625e',89,'user','帮我 检查 这个服务器的压力 8.130.14.34 ip=8.130.14.34','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"8.130.14.34\", \"targetType\": \"host\"}',0,'',1,'','2026-01-09 23:49:04.246'),(40,'afff8e56-7585-452d-8c66-fb08d097625e',89,'assistant','# 服务器压力巡检报告\n\n**主机信息**\n- 主机名称: 阿里云-ops\n- 主机IP: 8.130.14.34\n- 检查时间: 2026-01-09 23:49:04\n- 压力等级: **中等压力**\n\n## CPU信息\n- CPU核心数: 2\n- CPU使用率: 5.0%\n\n## 内存信息\n- 总内存: 1613 MB\n- 已使用: 1402 MB\n- 使用率: 86.92%\n\n## 磁盘信息\n- 根分区使用率: 47%\n```\nefivarfs        256K  7.2K  244K   3% /sys/firmware/efi/efivars\n/dev/vda3        40G   18G   20G  47% /\n/dev/vda2       197M  6.2M  191M   4% /boot/efi\noverlay          40G   18G   20G  47% /var/lib/docker/overlay2/47bf984488f4fb9f9ddbc314c7139add6b96959daa1deaad119cf180af7d0a4b/merged\noverlay          40G   18G   20G  47% /var/lib/docker/overlay2/b369df5540ea748d3f4acda3aecc5e3b5c40a2722e791bb7846dac2133626079/merged\noverlay          40G   18G   20G  47% /var/lib/docker/overlay2/b5237f18229c37dd13f5f602aca464d85f255e6c35d40437db8e719e985d292b/merged\n```\n\n## 系统负载\n- 负载平均值: 0.11, 0.11, 0.09\n\n## TOP 10 进程（按CPU使用率）\n| 用户 | PID | CPU% | MEM% | 命令 |\n|------|-----|------|------|------|\n| root | 1473688 | 12.9 | 0.6 | sshd: root@notty |\n| root | 2461 | 1.2 | 1.3 | /usr/local/aegis/aegis_client/aegis_12_81/AliYunDunMonitor |\n| root | 1389472 | 0.7 | 0.6 | /opt/agent/dodevops-agent |\n| root | 1023 | 0.6 | 0.6 | /usr/local/cloudmonitor/bin/argusagent |\n| root | 855 | 0.6 | 1.9 | /opt/pushgateway |\n| root | 1062 | 0.4 | 1.5 | /usr/local/aegis/aegis_client/aegis_12_81/AliYunDun |\n| mysql | 1115 | 0.4 | 31.8 | /usr/sbin/mysqld |\n| root | 854 | 0.3 | 4.9 | /opt/prometheus/prometheus/prometheus --config.file=/opt/prometheus/prometheus/prometheus.yml --storage.tsdb.path=/opt/prometheus/prometheus/data --web.listen-address=0.0.0.0:9090 |\n| redis | 856 | 0.2 | 0.7 | /usr/bin/redis-server 0.0.0.0:6379 |\n| root | 1186 | 0.2 | 2.8 | /usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock |\n\n','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"8.130.14.34\", \"targetType\": \"host\"}',0,'',2,'','2026-01-09 23:49:12.731'),(41,'afff8e56-7585-452d-8c66-fb08d097625e',89,'user','帮我 检查 这个服务器的压力 8.130.14.34 ip=8.130.14.34','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"8.130.14.34\", \"targetType\": \"host\"}',0,'',1,'','2026-01-09 23:49:55.498'),(42,'afff8e56-7585-452d-8c66-fb08d097625e',89,'assistant','# 服务器压力巡检报告\n\n**主机信息**\n- 主机名称: 阿里云-ops\n- 主机IP: 8.130.14.34\n- 检查时间: 2026-01-09 23:49:55\n- 压力等级: **中等压力**\n\n## CPU信息\n- CPU核心数: 2\n- CPU使用率: 8.7%\n\n## 内存信息\n- 总内存: 1613 MB\n- 已使用: 1387 MB\n- 使用率: 85.99%\n\n## 磁盘信息\n- 根分区使用率: 47%\n```\nefivarfs        256K  7.2K  244K   3% /sys/firmware/efi/efivars\n/dev/vda3        40G   18G   20G  47% /\n/dev/vda2       197M  6.2M  191M   4% /boot/efi\noverlay          40G   18G   20G  47% /var/lib/docker/overlay2/47bf984488f4fb9f9ddbc314c7139add6b96959daa1deaad119cf180af7d0a4b/merged\noverlay          40G   18G   20G  47% /var/lib/docker/overlay2/b369df5540ea748d3f4acda3aecc5e3b5c40a2722e791bb7846dac2133626079/merged\noverlay          40G   18G   20G  47% /var/lib/docker/overlay2/b5237f18229c37dd13f5f602aca464d85f255e6c35d40437db8e719e985d292b/merged\n```\n\n## 系统负载\n- 负载平均值: 0.60, 0.28, 0.15\n\n## TOP 10 进程（按CPU使用率）\n| 用户 | PID | CPU% | MEM% | 命令 |\n|------|-----|------|------|------|\n| root | 1474444 | 100 | 0.2 | ps aux --sort=-%cpu |\n| root | 1474394 | 12.9 | 0.6 | sshd: root@notty |\n| root | 2461 | 1.2 | 1.3 | /usr/local/aegis/aegis_client/aegis_12_81/AliYunDunMonitor |\n| root | 1389472 | 0.7 | 0.6 | /opt/agent/dodevops-agent |\n| root | 1023 | 0.6 | 0.6 | /usr/local/cloudmonitor/bin/argusagent |\n| root | 855 | 0.6 | 2.0 | /opt/pushgateway |\n| root | 1062 | 0.4 | 1.5 | /usr/local/aegis/aegis_client/aegis_12_81/AliYunDun |\n| mysql | 1115 | 0.4 | 31.6 | /usr/sbin/mysqld |\n| root | 854 | 0.3 | 5.0 | /opt/prometheus/prometheus/prometheus --config.file=/opt/prometheus/prometheus/prometheus.yml --storage.tsdb.path=/opt/prometheus/prometheus/data --web.listen-address=0.0.0.0:9090 |\n| redis | 856 | 0.2 | 0.7 | /usr/bin/redis-server 0.0.0.0:6379 |\n\n','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"8.130.14.34\", \"targetType\": \"host\"}',0,'',2,'','2026-01-09 23:50:05.366'),(43,'75aa6d6f-42d7-4bfc-b4f7-d87ce7f1662c',89,'user','180.76.231.65 帮忙检查下这个服务器的压力 ip=180.76.231.65','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',1,'','2026-01-10 00:02:41.561'),(44,'75aa6d6f-42d7-4bfc-b4f7-d87ce7f1662c',89,'assistant','# 服务器压力巡检报告\n\n**主机信息**\n- 主机名称: 百度云-ops\n- 主机IP: 180.76.231.65\n- 检查时间: 2026-01-10 00:02:42\n- 压力等级: **正常**\n\n## CPU信息\n- CPU核心数: 2\n- CPU使用率: 0.0%\n\n## 内存信息\n- 总内存: 1886 MB\n- 已使用: 1438 MB\n- 使用率: 76.25%\n\n## 磁盘信息\n- 根分区使用率: 54%\n```\n/dev/vda2        40G   21G   18G  54% /\noverlay          40G   21G   18G  54% /var/lib/docker/overlay2/36d3c5e2aa5b2e721d81c27c2923c1645e788af8c62af73ca78c562acb1dc6c6/merged\noverlay          40G   21G   18G  54% /var/lib/docker/overlay2/3a19d7d7ec48824a825359d54df4f9c59bde749d91c29a4c2bf7c39d5bedb135/merged\noverlay          40G   21G   18G  54% /var/lib/docker/overlay2/273c4fcfacbd0de3a76036142e72f3838211b24fc5753e7a39494964479738ce/merged\n```\n\n## 系统负载\n- 负载平均值: 0.02, 0.03, 0.00\n\n## TOP 10 进程（按CPU使用率）\n| 用户 | PID | CPU% | MEM% | 命令 |\n|------|-----|------|------|------|\n| root | 94 | 11.7 | 0.0 | [kswapd0] |\n| root | 750164 | 3.1 | 1.5 | /opt/heyeAgent/bin/heyeAgent |\n| root | 1 | 1.0 | 0.5 | /lib/systemd/systemd --switched-root --system --deserialize 31 nopti nospectre_v2 nospec_store_bypass_disable |\n| root | 4145227 | 1.0 | 0.4 | /lib/systemd/systemd --user |\n| root | 4083585 | 0.3 | 0.6 | /opt/agent/dodevops-agent |\n| root | 852 | 0.2 | 1.1 | /usr/bin/containerd |\n| systemd+ | 723 | 0.1 | 0.2 | /lib/systemd/systemd-networkd |\n| systemd+ | 725 | 0.1 | 0.4 | /lib/systemd/systemd-resolved |\n| fwupd-r+ | 920317 | 0.1 | 0.5 | /usr/bin/fwupdmgr refresh |\n| jenkins | 1997010 | 0.1 | 34.3 | /usr/bin/java -Djava.awt.headless=true -jar /usr/share/java/jenkins.war --webroot=/var/cache/jenkins/war --httpPort=8080 |\n\n','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',2,'','2026-01-10 00:02:49.514'),(45,'afff8e56-7585-452d-8c66-fb08d097625e',89,'user','帮我 检查 这个服务器的压力 8.130.14.34 ip=8.130.14.34','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"8.130.14.34\", \"targetType\": \"host\"}',0,'',1,'','2026-01-10 00:22:42.928'),(46,'afff8e56-7585-452d-8c66-fb08d097625e',89,'assistant','# 服务器压力巡检报告\n\n**主机信息**\n- 主机名称: 阿里云-ops\n- 主机IP: 8.130.14.34\n- 检查时间: 2026-01-10 00:22:43\n- 压力等级: **中等压力**\n\n## CPU信息\n- CPU核心数: 2\n- CPU使用率: 4.2%\n\n## 内存信息\n- 总内存: 1613 MB\n- 已使用: 1396 MB\n- 使用率: 86.55%\n\n## 磁盘信息\n- 根分区使用率: 47%\n```\nefivarfs        256K  7.2K  244K   3% /sys/firmware/efi/efivars\n/dev/vda3        40G   18G   20G  47% /\n/dev/vda2       197M  6.2M  191M   4% /boot/efi\noverlay          40G   18G   20G  47% /var/lib/docker/overlay2/47bf984488f4fb9f9ddbc314c7139add6b96959daa1deaad119cf180af7d0a4b/merged\noverlay          40G   18G   20G  47% /var/lib/docker/overlay2/b369df5540ea748d3f4acda3aecc5e3b5c40a2722e791bb7846dac2133626079/merged\noverlay          40G   18G   20G  47% /var/lib/docker/overlay2/b5237f18229c37dd13f5f602aca464d85f255e6c35d40437db8e719e985d292b/merged\n```\n\n## 系统负载\n- 负载平均值: 0.25, 0.13, 0.10\n\n## TOP 10 进程（按CPU使用率）\n| 用户 | PID | CPU% | MEM% | 命令 |\n|------|-----|------|------|------|\n| root | 1481918 | 100 | 0.2 | ps aux --sort=-%cpu |\n| root | 1481868 | 12.5 | 0.6 | sshd: root@notty |\n| root | 2461 | 1.2 | 1.3 | /usr/local/aegis/aegis_client/aegis_12_81/AliYunDunMonitor |\n| root | 1389472 | 0.7 | 0.7 | /opt/agent/dodevops-agent |\n| root | 1023 | 0.6 | 0.6 | /usr/local/cloudmonitor/bin/argusagent |\n| root | 855 | 0.6 | 1.9 | /opt/pushgateway |\n| root | 1062 | 0.4 | 1.5 | /usr/local/aegis/aegis_client/aegis_12_81/AliYunDun |\n| mysql | 1115 | 0.4 | 31.5 | /usr/sbin/mysqld |\n| root | 854 | 0.3 | 5.0 | /opt/prometheus/prometheus/prometheus --config.file=/opt/prometheus/prometheus/prometheus.yml --storage.tsdb.path=/opt/prometheus/prometheus/data --web.listen-address=0.0.0.0:9090 |\n| redis | 856 | 0.2 | 0.7 | /usr/bin/redis-server 0.0.0.0:6379 |\n\n','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"8.130.14.34\", \"targetType\": \"host\"}',0,'',2,'','2026-01-10 00:22:53.962'),(47,'75aa6d6f-42d7-4bfc-b4f7-d87ce7f1662c',89,'user','巡检下这个服务器  180.76.231.65 ip=180.76.231.65','inspection',0.70,'{\"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',1,'','2026-01-10 00:29:41.038'),(48,'75aa6d6f-42d7-4bfc-b4f7-d87ce7f1662c',89,'assistant','✅ 已为 1 台主机创建巡检任务\n\n目标主机：\n- 百度云-ops (180.76.231.65)\n\n巡检任务正在执行中，稍后可查看巡检报告...','inspection',0.70,'{\"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',2,'','2026-01-10 00:29:41.585'),(49,'75aa6d6f-42d7-4bfc-b4f7-d87ce7f1662c',89,'user','帮忙查看下 139.9.205.38 这个服务器都部署了什么服务 ip=139.9.205.38','incident',0.70,'{\"targetName\": \"139.9.205.38\", \"targetType\": \"host\"}',0,'',1,'','2026-01-10 01:05:43.057'),(50,'75aa6d6f-42d7-4bfc-b4f7-d87ce7f1662c',89,'assistant','请指定分析类型，例如：\n- 检查服务器压力\n- 分析内存问题\n- 查看磁盘使用\n- 检查端口情况\n- 查看进程状态','incident',0.70,'{\"targetName\": \"139.9.205.38\", \"targetType\": \"host\"}',0,'',2,'','2026-01-10 01:05:43.471'),(51,'75aa6d6f-42d7-4bfc-b4f7-d87ce7f1662c',89,'user','检查端口情况 ip=139.9.205.38','incident',0.70,'{\"checkType\": \"port\", \"targetName\": \"139.9.205.38\", \"targetType\": \"host\"}',0,'',1,'','2026-01-10 01:06:01.204'),(52,'75aa6d6f-42d7-4bfc-b4f7-d87ce7f1662c',89,'assistant','# 端口扫描报告\n\n**主机**: 华为云-ops (139.9.205.38)\n\n⚠️ 端口扫描功能开发中，敬请期待...\n\n建议手动执行：`ss -tuln | grep LISTEN`','incident',0.70,'{\"checkType\": \"port\", \"targetName\": \"139.9.205.38\", \"targetType\": \"host\"}',0,'',2,'','2026-01-10 01:06:01.732'),(53,'75aa6d6f-42d7-4bfc-b4f7-d87ce7f1662c',89,'user','查看磁盘使用 ip=139.9.205.38','incident',0.70,'{\"checkType\": \"disk\", \"targetName\": \"139.9.205.38\", \"targetType\": \"host\"}',0,'',1,'','2026-01-10 01:06:22.196'),(54,'75aa6d6f-42d7-4bfc-b4f7-d87ce7f1662c',89,'assistant','# 磁盘空间分析\n\n**主机**: 华为云-ops (139.9.205.38)\n\n## 磁盘使用情况\n\n⚠️ 详细磁盘分析功能开发中...\n\n','incident',0.70,'{\"checkType\": \"disk\", \"targetName\": \"139.9.205.38\", \"targetType\": \"host\"}',0,'',2,'','2026-01-10 01:06:26.922'),(55,'75aa6d6f-42d7-4bfc-b4f7-d87ce7f1662c',89,'user','帮忙查看下 139.9.205.38 这个服务器都部署了什么服务 ip=139.9.205.38','incident',0.70,'{\"checkType\": \"service\", \"targetName\": \"139.9.205.38\", \"targetType\": \"host\"}',0,'',1,'','2026-01-10 01:12:37.365'),(56,'75aa6d6f-42d7-4bfc-b4f7-d87ce7f1662c',89,'assistant','# 服务列表\n\n**主机**: 华为云-ops (139.9.205.38)\n\n## 检测到的服务\n\n- ✅ **MySQL 数据库** (PID: 0, CPU: 0.0%, MEM: 0.0%)\n- ✅ **Docker 容器** (PID: 0, CPU: 0.0%, MEM: 0.0%)\n- ✅ **Redis 缓存** (PID: 0, CPU: 0.0%, MEM: 0.0%)\n','incident',0.70,'{\"checkType\": \"service\", \"targetName\": \"139.9.205.38\", \"targetType\": \"host\"}',0,'',2,'','2026-01-10 01:12:42.277'),(57,'75aa6d6f-42d7-4bfc-b4f7-d87ce7f1662c',89,'user','帮忙查看下 139.9.205.38 这个服务器内存为什么这个高 ip=139.9.205.38','incident',0.70,'{\"checkType\": \"memory\", \"targetName\": \"139.9.205.38\", \"targetType\": \"host\", \"needAnalysis\": \"true\"}',0,'',1,'','2026-01-10 01:13:08.448'),(58,'75aa6d6f-42d7-4bfc-b4f7-d87ce7f1662c',89,'assistant','# 服务器内存分析报告\n\n**主机信息**\n- 主机名称: 华为云-ops\n- 主机IP: 139.9.205.38\n- 分析时间: 2026-01-10 01:13:09\n\n## 内存使用模式\n内存使用率较高(>80%); 可用内存较低(<20%)\n\n## 内存详细信息\n- 总内存: 1775 MB\n- 已使用: 1596 MB\n- 可用内存: 178 MB\n- 使用率: 89.92%\n\n## Swap信息\n- Swap总量: 2047 MB\n- Swap已使用: 0 MB\n- Swap使用率: 0.00%\n\n## TOP 15 进程（按内存使用率）\n| 用户 | PID | CPU% | MEM% | VSZ | RSS | 命令 |\n|------|-----|------|------|-----|-----|------|\n| dnsmasq | 1958583 | 0.3 | 22.3 | 1950420 | 405852 | mysqld --default-authentication-plugin=caching_sha2_password --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci |\n| mysql | 1269514 | 0.3 | 20.8 | 1803192 | 379328 | /usr/sbin/mysqld |\n| root | 2609356 | 0.1 | 11.6 | 857268 | 211064 | /opt/netops/netops_linux/main |\n| root | 4162698 | 0.0 | 3.2 | 1272612 | 59096 | /www/allinssl/allinssl start |\n| root | 2627163 | 0.1 | 3.2 | 2957452 | 58380 | /bin/prometheus --config.file=/etc/prometheus/prometheus.yml --storage.tsdb.path=/prometheus --web.console.libraries=/usr/share/prometheus/console_libraries --web.console.templates=/usr/share/prometheus/consoles --storage.tsdb.retention.time=15d --web.enable-lifecycle |\n| root | 379517 | 0.2 | 2.9 | 2646480 | 53676 | /usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock |\n| root | 2843811 | 0.1 | 2.3 | 1413280 | 42444 | /app/devops |\n| root | 379334 | 0.1 | 1.9 | 2090148 | 35476 | /usr/bin/containerd |\n| root | 2596894 | 0.0 | 1.4 | 223416 | 27008 | /sbin/multipathd -d -s |\n| dnsmasq | 880425 | 0.0 | 0.9 | 216332 | 17792 | postgres: checkpointer |\n| root | 2591913 | 0.0 | 0.9 | 177808 | 17696 | /usr/lib/systemd/systemd-journald |\n| root | 1958877 | 0.0 | 0.7 | 1118136 | 13812 | /usr/local/hostguard/bin/containerserver-1.0.0 |\n| dnsmasq | 880324 | 0.0 | 0.7 | 216056 | 13568 | postgres |\n| root | 1958776 | 0.1 | 0.7 | 1309392 | 13116 | /usr/local/hostguard/bin/hostguard --daemon |\n| root | 894 | 0.0 | 0.7 | 109632 | 12800 | /usr/bin/python3 /usr/share/unattended-upgrades/unattended-upgrade-shutdown --wait-for-signal |\n\n## 发现的问题\n1. 可用内存不足: 178 MB\n\n## AI建议\n- **可用内存不足**: 建议释放不必要的缓存，关闭不需要的服务，或增加物理内存\n\n**通用优化建议**:\n1. 定期监控内存使用情况，及时发现异常\n2. 使用 `free -h` 和 `top` 命令实时查看内存状态\n3. 检查系统日志，查找可能的内存相关错误\n4. 考虑使用内存分析工具（如 valgrind）检查内存泄漏\n','incident',0.70,'{\"checkType\": \"memory\", \"targetName\": \"139.9.205.38\", \"targetType\": \"host\", \"needAnalysis\": \"true\"}',0,'',2,'','2026-01-10 01:13:10.326'),(59,'75aa6d6f-42d7-4bfc-b4f7-d87ce7f1662c',89,'user','帮忙查看下 139.9.205.38 磁盘快满了，哪些目录占空间最多 ip=139.9.205.38','incident',0.70,'{\"checkType\": \"disk\", \"queryType\": \"directories\", \"targetName\": \"139.9.205.38\", \"targetType\": \"host\"}',0,'',1,'','2026-01-10 01:13:43.551'),(60,'75aa6d6f-42d7-4bfc-b4f7-d87ce7f1662c',89,'assistant','# 磁盘空间分析\n\n**主机**: 华为云-ops (139.9.205.38)\n\n## 磁盘使用情况\n\n⚠️ 详细磁盘分析功能开发中...\n\n## 占用空间最多的目录\n\n? 建议手动执行以下命令查看：\n\n```bash\ndu -h --max-depth=1 / 2>/dev/null | sort -hr | head -10\n```\n','incident',0.70,'{\"checkType\": \"disk\", \"queryType\": \"directories\", \"targetName\": \"139.9.205.38\", \"targetType\": \"host\"}',0,'',2,'','2026-01-10 01:13:48.052'),(61,'b8d98633-60f5-44c3-8872-7c8109dd3578',89,'user','47.252.17.169   帮我检查下这台服务器的压力 ip=47.252.17.169','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"47.252.17.169\", \"targetType\": \"host\"}',0,'',1,'','2026-01-10 16:49:14.846'),(62,'b8d98633-60f5-44c3-8872-7c8109dd3578',89,'assistant','❌ 分析失败: 未找到名称为 47.252.17.169 的服务器: record not found','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"47.252.17.169\", \"targetType\": \"host\"}',0,'',2,'','2026-01-10 16:49:15.631'),(63,'d6d18098-ac55-4168-b39f-a78d0f6bdb64',89,'user','47.244.194.219  帮我检查下这个服务器压力 ip=47.244.194.219','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"47.244.194.219\", \"targetType\": \"host\"}',0,'',1,'','2026-01-10 17:00:18.039'),(64,'d6d18098-ac55-4168-b39f-a78d0f6bdb64',89,'assistant','❌ 分析失败: 未找到名称为 47.244.194.219 的服务器: record not found','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"47.244.194.219\", \"targetType\": \"host\"}',0,'',2,'','2026-01-10 17:00:18.851'),(65,'4d481c48-cf6b-4c87-8262-7f6bdb1c3bd9',89,'user','8.130.14.34  帮我检查这个服务器的压力 ip=8.130.14.34','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"8.130.14.34\", \"targetType\": \"host\"}',0,'',1,'','2026-01-10 17:31:37.733'),(66,'4d481c48-cf6b-4c87-8262-7f6bdb1c3bd9',89,'assistant','# 服务器压力巡检报告\n\n**主机信息**\n- 主机名称: 阿里云-ops\n- 主机IP: 8.130.14.34\n- 检查时间: 2026-01-10 17:31:38\n- 压力等级: **中等压力**\n\n## CPU信息\n- CPU核心数: 2\n- CPU使用率: 4.2%\n\n## 内存信息\n- 总内存: 1613 MB\n- 已使用: 1391 MB\n- 使用率: 86.24%\n\n## 磁盘信息\n- 根分区使用率: 48%\n```\nefivarfs        256K  7.2K  244K   3% /sys/firmware/efi/efivars\n/dev/vda3        40G   18G   20G  48% /\n/dev/vda2       197M  6.2M  191M   4% /boot/efi\noverlay          40G   18G   20G  48% /var/lib/docker/overlay2/47bf984488f4fb9f9ddbc314c7139add6b96959daa1deaad119cf180af7d0a4b/merged\noverlay          40G   18G   20G  48% /var/lib/docker/overlay2/b369df5540ea748d3f4acda3aecc5e3b5c40a2722e791bb7846dac2133626079/merged\noverlay          40G   18G   20G  48% /var/lib/docker/overlay2/b5237f18229c37dd13f5f602aca464d85f255e6c35d40437db8e719e985d292b/merged\n```\n\n## 系统负载\n- 负载平均值: 0.58, 0.23, 0.23\n\n## TOP 10 进程（按CPU使用率）\n| 用户 | PID | CPU% | MEM% | 命令 |\n|------|-----|------|------|------|\n| root | 1676420 | 11.1 | 0.6 | sshd: root@notty |\n| root | 2461 | 1.2 | 1.4 | /usr/local/aegis/aegis_client/aegis_12_81/AliYunDunMonitor |\n| root | 1389472 | 0.7 | 1.1 | /opt/agent/dodevops-agent |\n| root | 1023 | 0.6 | 0.5 | /usr/local/cloudmonitor/bin/argusagent |\n| root | 855 | 0.6 | 2.2 | /opt/pushgateway |\n| root | 1062 | 0.4 | 1.6 | /usr/local/aegis/aegis_client/aegis_12_81/AliYunDun |\n| mysql | 1115 | 0.4 | 34.3 | /usr/sbin/mysqld |\n| root | 854 | 0.3 | 4.7 | /opt/prometheus/prometheus/prometheus --config.file=/opt/prometheus/prometheus/prometheus.yml --storage.tsdb.path=/opt/prometheus/prometheus/data --web.listen-address=0.0.0.0:9090 |\n| root | 1112559 | 0.2 | 1.1 | /usr/lib/snapd/snapd |\n| root | 1186 | 0.2 | 3.0 | /usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock |\n\n','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"8.130.14.34\", \"targetType\": \"host\"}',0,'',2,'','2026-01-10 17:31:47.678'),(67,'d3a86402-63a8-4891-9bd7-1144508c07c8',89,'user','8.130.14.34  检查这个服务器的压力 ip=8.130.14.34','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"8.130.14.34\", \"targetType\": \"host\"}',0,'',1,'','2026-01-10 17:34:06.305'),(68,'d3a86402-63a8-4891-9bd7-1144508c07c8',89,'assistant','# 服务器压力巡检报告\n\n**主机信息**\n- 主机名称: 阿里云-ops\n- 主机IP: 8.130.14.34\n- 检查时间: 2026-01-10 17:34:06\n- 压力等级: **中等压力**\n\n## CPU信息\n- CPU核心数: 2\n- CPU使用率: 4.3%\n\n## 内存信息\n- 总内存: 1613 MB\n- 已使用: 1398 MB\n- 使用率: 86.67%\n\n## 磁盘信息\n- 根分区使用率: 48%\n```\nefivarfs        256K  7.2K  244K   3% /sys/firmware/efi/efivars\n/dev/vda3        40G   18G   20G  48% /\n/dev/vda2       197M  6.2M  191M   4% /boot/efi\noverlay          40G   18G   20G  48% /var/lib/docker/overlay2/47bf984488f4fb9f9ddbc314c7139add6b96959daa1deaad119cf180af7d0a4b/merged\noverlay          40G   18G   20G  48% /var/lib/docker/overlay2/b369df5540ea748d3f4acda3aecc5e3b5c40a2722e791bb7846dac2133626079/merged\noverlay          40G   18G   20G  48% /var/lib/docker/overlay2/b5237f18229c37dd13f5f602aca464d85f255e6c35d40437db8e719e985d292b/merged\n```\n\n## 系统负载\n- 负载平均值: 0.05, 0.14, 0.19\n\n## TOP 10 进程（按CPU使用率）\n| 用户 | PID | CPU% | MEM% | 命令 |\n|------|-----|------|------|------|\n| root | 1677327 | 10.9 | 0.6 | sshd: root@notty |\n| root | 2461 | 1.2 | 1.4 | /usr/local/aegis/aegis_client/aegis_12_81/AliYunDunMonitor |\n| root | 1389472 | 0.7 | 1.0 | /opt/agent/dodevops-agent |\n| root | 1023 | 0.6 | 0.5 | /usr/local/cloudmonitor/bin/argusagent |\n| root | 855 | 0.6 | 2.2 | /opt/pushgateway |\n| root | 1062 | 0.4 | 1.5 | /usr/local/aegis/aegis_client/aegis_12_81/AliYunDun |\n| mysql | 1115 | 0.4 | 34.1 | /usr/sbin/mysqld |\n| root | 854 | 0.3 | 4.6 | /opt/prometheus/prometheus/prometheus --config.file=/opt/prometheus/prometheus/prometheus.yml --storage.tsdb.path=/opt/prometheus/prometheus/data --web.listen-address=0.0.0.0:9090 |\n| root | 1112559 | 0.2 | 1.5 | /usr/lib/snapd/snapd |\n| root | 1186 | 0.2 | 2.8 | /usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock |\n\n','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"8.130.14.34\", \"targetType\": \"host\"}',0,'',2,'','2026-01-10 17:34:16.047'),(69,'1fe65393-31f5-40e2-9839-3f0426eb1ff1',89,'user','8.130.14.34  帮忙检查这个服务器的压力 ip=8.130.14.34','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"8.130.14.34\", \"targetType\": \"host\"}',0,'',1,'','2026-01-10 17:35:59.639'),(70,'1fe65393-31f5-40e2-9839-3f0426eb1ff1',89,'assistant','# 服务器压力巡检报告\n\n**主机信息**\n- 主机名称: 阿里云-ops\n- 主机IP: 8.130.14.34\n- 检查时间: 2026-01-10 17:36:00\n- 压力等级: **中等压力**\n\n## CPU信息\n- CPU核心数: 2\n- CPU使用率: 4.8%\n\n## 内存信息\n- 总内存: 1613 MB\n- 已使用: 1405 MB\n- 使用率: 87.10%\n\n## 磁盘信息\n- 根分区使用率: 48%\n```\nefivarfs        256K  7.2K  244K   3% /sys/firmware/efi/efivars\n/dev/vda3        40G   18G   20G  48% /\n/dev/vda2       197M  6.2M  191M   4% /boot/efi\noverlay          40G   18G   20G  48% /var/lib/docker/overlay2/47bf984488f4fb9f9ddbc314c7139add6b96959daa1deaad119cf180af7d0a4b/merged\noverlay          40G   18G   20G  48% /var/lib/docker/overlay2/b369df5540ea748d3f4acda3aecc5e3b5c40a2722e791bb7846dac2133626079/merged\noverlay          40G   18G   20G  48% /var/lib/docker/overlay2/b5237f18229c37dd13f5f602aca464d85f255e6c35d40437db8e719e985d292b/merged\n```\n\n## 系统负载\n- 负载平均值: 0.19, 0.16, 0.18\n\n## TOP 10 进程（按CPU使用率）\n| 用户 | PID | CPU% | MEM% | 命令 |\n|------|-----|------|------|------|\n| root | 1678247 | 100 | 0.2 | ps aux --sort=-%cpu |\n| root | 1678197 | 12.3 | 0.6 | sshd: root@notty |\n| root | 2461 | 1.2 | 1.4 | /usr/local/aegis/aegis_client/aegis_12_81/AliYunDunMonitor |\n| root | 1389472 | 0.7 | 0.9 | /opt/agent/dodevops-agent |\n| root | 1023 | 0.6 | 0.7 | /usr/local/cloudmonitor/bin/argusagent |\n| root | 855 | 0.6 | 2.1 | /opt/pushgateway |\n| root | 1062 | 0.4 | 1.5 | /usr/local/aegis/aegis_client/aegis_12_81/AliYunDun |\n| mysql | 1115 | 0.4 | 34.1 | /usr/sbin/mysqld |\n| root | 854 | 0.3 | 4.5 | /opt/prometheus/prometheus/prometheus --config.file=/opt/prometheus/prometheus/prometheus.yml --storage.tsdb.path=/opt/prometheus/prometheus/data --web.listen-address=0.0.0.0:9090 |\n| root | 1112559 | 0.2 | 1.1 | /usr/lib/snapd/snapd |\n\n','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"8.130.14.34\", \"targetType\": \"host\"}',0,'',2,'','2026-01-10 17:36:08.875'),(71,'1fe65393-31f5-40e2-9839-3f0426eb1ff1',89,'user','查看 47.244.194.219  这个服务器的压力 ip=47.244.194.219','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"47.244.194.219\", \"targetType\": \"host\"}',0,'',1,'','2026-01-10 17:37:26.927'),(72,'1fe65393-31f5-40e2-9839-3f0426eb1ff1',89,'assistant','❌ 该机器 47.244.194.219 没有配置密钥无法访问，请增加密钥后重新提问','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"47.244.194.219\", \"targetType\": \"host\"}',0,'',2,'','2026-01-10 17:37:27.730'),(73,'e14c7db6-1578-4870-bb3f-742e9a2ff461',89,'user','检查 139.9.205.38 内存为什么这么高 ip=139.9.205.38','incident',0.70,'{\"checkType\": \"memory\", \"targetName\": \"139.9.205.38\", \"targetType\": \"host\", \"needAnalysis\": \"true\"}',0,'',1,'','2026-01-10 18:10:25.381'),(74,'e14c7db6-1578-4870-bb3f-742e9a2ff461',89,'assistant','# 服务器内存分析报告\n\n**主机信息**\n- 主机名称: 华为云-ops\n- 主机IP: 139.9.205.38\n- 分析时间: 2026-01-10 18:10:25\n\n## 内存使用模式\n内存使用率较高(>80%); 可用内存较低(<20%)\n\n## 内存详细信息\n- 总内存: 1775 MB\n- 已使用: 1587 MB\n- 可用内存: 188 MB\n- 使用率: 89.41%\n\n## Swap信息\n- Swap总量: 2047 MB\n- Swap已使用: 0 MB\n- Swap使用率: 0.00%\n\n## TOP 15 进程（按内存使用率）\n| 用户 | PID | CPU% | MEM% | VSZ | RSS | 命令 |\n|------|-----|------|------|-----|-----|------|\n| dnsmasq | 1958583 | 0.3 | 22.4 | 1954684 | 408284 | mysqld --default-authentication-plugin=caching_sha2_password --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci |\n| mysql | 1269514 | 0.3 | 21.1 | 1803192 | 384448 | /usr/sbin/mysqld |\n| root | 2609356 | 0.1 | 11.9 | 857268 | 217336 | /opt/netops/netops_linux/main |\n| root | 4162698 | 0.0 | 3.2 | 1272612 | 59104 | /www/allinssl/allinssl start |\n| root | 379517 | 0.2 | 2.9 | 2646480 | 53908 | /usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock |\n| root | 2627163 | 0.1 | 2.7 | 2815672 | 50548 | /bin/prometheus --config.file=/etc/prometheus/prometheus.yml --storage.tsdb.path=/prometheus --web.console.libraries=/usr/share/prometheus/console_libraries --web.console.templates=/usr/share/prometheus/consoles --storage.tsdb.retention.time=15d --web.enable-lifecycle |\n| root | 2843811 | 0.1 | 2.2 | 1413280 | 40704 | /app/devops |\n| root | 379334 | 0.1 | 1.7 | 2090148 | 31832 | /usr/bin/containerd |\n| root | 2596894 | 0.0 | 1.4 | 223416 | 27008 | /sbin/multipathd -d -s |\n| dnsmasq | 880425 | 0.0 | 0.9 | 216332 | 17792 | postgres: checkpointer |\n| root | 2591913 | 0.0 | 0.8 | 91544 | 16172 | /usr/lib/systemd/systemd-journald |\n| dnsmasq | 880324 | 0.0 | 0.7 | 216056 | 13440 | postgres |\n| root | 2626785 | 0.0 | 0.7 | 1240440 | 13336 | /bin/pushgateway |\n| root | 894 | 0.0 | 0.7 | 109632 | 12800 | /usr/bin/python3 /usr/share/unattended-upgrades/unattended-upgrade-shutdown --wait-for-signal |\n| test | 1269667 | 0.0 | 0.6 | 45684 | 11748 | /home/test/.dhpcd |\n\n## 发现的问题\n1. 可用内存不足: 188 MB\n\n## ? AI智能分析与优化建议\n\nAI分析服务暂时不可用，请稍后重试或联系管理员配置AI模型。\n','incident',0.70,'{\"checkType\": \"memory\", \"targetName\": \"139.9.205.38\", \"targetType\": \"host\", \"needAnalysis\": \"true\"}',0,'',2,'','2026-01-10 18:10:28.243'),(75,'e14c7db6-1578-4870-bb3f-742e9a2ff461',89,'user','检查 139.9.205.38 内存为什么这么高 ip=139.9.205.38','incident',0.70,'{\"checkType\": \"memory\", \"targetName\": \"139.9.205.38\", \"targetType\": \"host\", \"needAnalysis\": \"true\"}',0,'',1,'','2026-01-10 18:15:04.711'),(76,'e14c7db6-1578-4870-bb3f-742e9a2ff461',89,'assistant','# 服务器内存分析报告\n\n**主机信息**\n- 主机名称: 华为云-ops\n- 主机IP: 139.9.205.38\n- 分析时间: 2026-01-10 18:15:05\n\n## 内存使用模式\n内存使用率极高(>90%); 可用内存极低(<10%)\n\n## 内存详细信息\n- 总内存: 1775 MB\n- 已使用: 1600 MB\n- 可用内存: 175 MB\n- 使用率: 90.14%\n\n## Swap信息\n- Swap总量: 2047 MB\n- Swap已使用: 0 MB\n- Swap使用率: 0.00%\n\n## TOP 15 进程（按内存使用率）\n| 用户 | PID | CPU% | MEM% | VSZ | RSS | 命令 |\n|------|-----|------|------|-----|-----|------|\n| dnsmasq | 1958583 | 0.3 | 22.5 | 1954684 | 409308 | mysqld --default-authentication-plugin=caching_sha2_password --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci |\n| mysql | 1269514 | 0.3 | 20.8 | 1803192 | 379328 | /usr/sbin/mysqld |\n| root | 2609356 | 0.1 | 11.6 | 857268 | 211064 | /opt/netops/netops_linux/main |\n| root | 2627163 | 0.1 | 3.5 | 2815672 | 65432 | /bin/prometheus --config.file=/etc/prometheus/prometheus.yml --storage.tsdb.path=/prometheus --web.console.libraries=/usr/share/prometheus/console_libraries --web.console.templates=/usr/share/prometheus/consoles --storage.tsdb.retention.time=15d --web.enable-lifecycle |\n| root | 4162698 | 0.0 | 3.2 | 1272612 | 59032 | /www/allinssl/allinssl start |\n| root | 379517 | 0.2 | 3.1 | 2646480 | 56852 | /usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock |\n| root | 2843811 | 0.1 | 2.5 | 1413280 | 45748 | /app/devops |\n| root | 379334 | 0.1 | 1.7 | 2090148 | 32060 | /usr/bin/containerd |\n| root | 2596894 | 0.0 | 1.4 | 223416 | 27008 | /sbin/multipathd -d -s |\n| dnsmasq | 880425 | 0.0 | 0.9 | 216332 | 17792 | postgres: checkpointer |\n| dnsmasq | 880324 | 0.0 | 0.7 | 216056 | 13440 | postgres |\n| root | 3039604 | 0.0 | 0.7 | 1902572 | 12928 | runc --root /var/run/docker/runtime-runc/moby --log /run/containerd/io.containerd.runtime.v2.task/moby/c13c84710575d8c918abfeda7a405b515d6b9042581017a6725dc10672608668/log.json --log-format json --systemd-cgroup exec --process /tmp/runc-process2152991518 --detach --pid-file /run/containerd/io.containerd.runtime.v2.task/moby/c13c84710575d8c918abfeda7a405b515d6b9042581017a6725dc10672608668/20c4c9d090b9f2446e2ce69cfa241f22be2e118588eb4771932a6d30a73347e0.pid c13c84710575d8c918abfeda7a405b515d6b9042581017a6725dc10672608668 |\n| root | 894 | 0.0 | 0.7 | 109632 | 12800 | /usr/bin/python3 /usr/share/unattended-upgrades/unattended-upgrade-shutdown --wait-for-signal |\n| root | 3039605 | 0.0 | 0.6 | 1828584 | 12672 | runc --root /var/run/docker/runtime-runc/moby --log /run/containerd/io.containerd.runtime.v2.task/moby/9f487f841e39bdb68678c0f49d02f20618b497b71b5980631e62aaa58257e437/log.json --log-format json --systemd-cgroup exec --process /tmp/runc-process1507548198 --detach --pid-file /run/containerd/io.containerd.runtime.v2.task/moby/9f487f841e39bdb68678c0f49d02f20618b497b71b5980631e62aaa58257e437/88b6483c5d6112ea480b30d2b1897fd092693442fa18493c062fd267208c4ca9.pid 9f487f841e39bdb68678c0f49d02f20618b497b71b5980631e62aaa58257e437 |\n| root | 2626785 | 0.0 | 0.6 | 1240440 | 12312 | /bin/pushgateway |\n\n## 发现的问题\n1. 内存使用率过高: 90.14%\n2. 可用内存不足: 175 MB\n\n## ? AI智能分析与优化建议\n\nAI分析服务暂时不可用，请稍后重试或联系管理员配置AI模型。\n','incident',0.70,'{\"checkType\": \"memory\", \"targetName\": \"139.9.205.38\", \"targetType\": \"host\", \"needAnalysis\": \"true\"}',0,'',2,'','2026-01-10 18:15:07.616'),(77,'e14c7db6-1578-4870-bb3f-742e9a2ff461',89,'user','linux巡检哪些内容 ip=139.9.205.38','inspection',0.70,'{\"targetName\": \"139.9.205.38\", \"targetType\": \"host\"}',0,'',1,'','2026-01-10 18:18:18.124'),(78,'e14c7db6-1578-4870-bb3f-742e9a2ff461',89,'assistant','✅ 已为 1 台主机创建巡检任务\n\n目标主机：\n- 华为云-ops (139.9.205.38)\n\n巡检任务正在执行中，稍后可查看巡检报告...','inspection',0.70,'{\"targetName\": \"139.9.205.38\", \"targetType\": \"host\"}',0,'',2,'','2026-01-10 18:18:18.681'),(79,'e14c7db6-1578-4870-bb3f-742e9a2ff461',89,'user','检查 139.9.205.38 内存为什么这么高 ip=139.9.205.38','incident',0.70,'{\"checkType\": \"memory\", \"targetName\": \"139.9.205.38\", \"targetType\": \"host\", \"needAnalysis\": \"true\"}',0,'',1,'','2026-01-10 18:23:46.360'),(80,'e14c7db6-1578-4870-bb3f-742e9a2ff461',89,'assistant','# 服务器内存分析报告\n\n**主机信息**\n- 主机名称: 华为云-ops\n- 主机IP: 139.9.205.38\n- 分析时间: 2026-01-10 18:23:46\n\n## 内存使用模式\n内存使用率较高(>80%); 可用内存较低(<20%)\n\n## 内存详细信息\n- 总内存: 1775 MB\n- 已使用: 1576 MB\n- 可用内存: 198 MB\n- 使用率: 88.79%\n\n## Swap信息\n- Swap总量: 2047 MB\n- Swap已使用: 0 MB\n- Swap使用率: 0.00%\n\n## TOP 15 进程（按内存使用率）\n| 用户 | PID | CPU% | MEM% | VSZ | RSS | 命令 |\n|------|-----|------|------|-----|-----|------|\n| dnsmasq | 1958583 | 0.3 | 22.4 | 1954684 | 408540 | mysqld --default-authentication-plugin=caching_sha2_password --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci |\n| mysql | 1269514 | 0.3 | 21.1 | 1803192 | 384320 | /usr/sbin/mysqld |\n| root | 2609356 | 0.1 | 11.9 | 857268 | 217464 | /opt/netops/netops_linux/main |\n| root | 4162698 | 0.0 | 3.2 | 1272612 | 59080 | /www/allinssl/allinssl start |\n| root | 379517 | 0.2 | 2.9 | 2646480 | 54420 | /usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock |\n| root | 2627163 | 0.1 | 2.7 | 2815672 | 49844 | /bin/prometheus --config.file=/etc/prometheus/prometheus.yml --storage.tsdb.path=/prometheus --web.console.libraries=/usr/share/prometheus/console_libraries --web.console.templates=/usr/share/prometheus/consoles --storage.tsdb.retention.time=15d --web.enable-lifecycle |\n| root | 2843811 | 0.1 | 2.2 | 1413280 | 41372 | /app/devops |\n| root | 379334 | 0.1 | 1.9 | 2090148 | 34620 | /usr/bin/containerd |\n| root | 2596894 | 0.0 | 1.4 | 223416 | 27008 | /sbin/multipathd -d -s |\n| dnsmasq | 880425 | 0.0 | 0.9 | 216332 | 17792 | postgres: checkpointer |\n| dnsmasq | 880324 | 0.0 | 0.7 | 216056 | 13440 | postgres |\n| root | 2626785 | 0.0 | 0.7 | 1240440 | 13336 | /bin/pushgateway |\n| root | 894 | 0.0 | 0.7 | 109632 | 12800 | /usr/bin/python3 /usr/share/unattended-upgrades/unattended-upgrade-shutdown --wait-for-signal |\n| test | 1269667 | 0.0 | 0.6 | 45692 | 11748 | /home/test/.dhpcd |\n| root | 2842024 | 0.0 | 0.6 | 601448 | 11608 | /usr/libexec/fwupd/fwupd |\n\n## 发现的问题\n1. 可用内存不足: 198 MB\n\n## ? AI智能分析与优化建议\n\n### ? 基础优化建议\n\n#### ? 内存使用率较高 (>80%)\n\n建议监控内存趋势，考虑优化配置或扩容\n\n','incident',0.70,'{\"checkType\": \"memory\", \"targetName\": \"139.9.205.38\", \"targetType\": \"host\", \"needAnalysis\": \"true\"}',0,'',2,'','2026-01-10 18:23:49.186'),(81,'e14c7db6-1578-4870-bb3f-742e9a2ff461',89,'user','检查 139.9.205.38 内存为什么这么高 ip=139.9.205.38','incident',0.70,'{\"checkType\": \"memory\", \"targetName\": \"139.9.205.38\", \"targetType\": \"host\", \"needAnalysis\": \"true\"}',0,'',1,'','2026-01-10 18:28:27.760'),(82,'e14c7db6-1578-4870-bb3f-742e9a2ff461',89,'assistant','# 服务器内存分析报告\n\n**主机信息**\n- 主机名称: 华为云-ops\n- 主机IP: 139.9.205.38\n- 分析时间: 2026-01-10 18:28:28\n\n## 内存使用模式\n内存使用率较高(>80%); 可用内存较低(<20%)\n\n## 内存详细信息\n- 总内存: 1775 MB\n- 已使用: 1583 MB\n- 可用内存: 191 MB\n- 使用率: 89.18%\n\n## Swap信息\n- Swap总量: 2047 MB\n- Swap已使用: 0 MB\n- Swap使用率: 0.00%\n\n## TOP 15 进程（按内存使用率）\n| 用户 | PID | CPU% | MEM% | VSZ | RSS | 命令 |\n|------|-----|------|------|-----|-----|------|\n| dnsmasq | 1958583 | 0.3 | 22.4 | 1954684 | 408540 | mysqld --default-authentication-plugin=caching_sha2_password --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci |\n| mysql | 1269514 | 0.3 | 21.1 | 1803192 | 384320 | /usr/sbin/mysqld |\n| root | 2609356 | 0.1 | 11.8 | 857268 | 215416 | /opt/netops/netops_linux/main |\n| root | 4162698 | 0.0 | 3.2 | 1272612 | 59272 | /www/allinssl/allinssl start |\n| root | 2627163 | 0.1 | 3.2 | 2815672 | 58360 | /bin/prometheus --config.file=/etc/prometheus/prometheus.yml --storage.tsdb.path=/prometheus --web.console.libraries=/usr/share/prometheus/console_libraries --web.console.templates=/usr/share/prometheus/consoles --storage.tsdb.retention.time=15d --web.enable-lifecycle |\n| root | 379517 | 0.2 | 3.0 | 2646480 | 54932 | /usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock |\n| root | 2843811 | 0.1 | 2.3 | 1413280 | 42732 | /app/devops |\n| root | 379334 | 0.1 | 1.8 | 2090148 | 34108 | /usr/bin/containerd |\n| root | 2596894 | 0.0 | 1.4 | 223416 | 27008 | /sbin/multipathd -d -s |\n| dnsmasq | 880425 | 0.0 | 0.9 | 216332 | 17792 | postgres: checkpointer |\n| dnsmasq | 880324 | 0.0 | 0.7 | 216056 | 13440 | postgres |\n| root | 2626785 | 0.0 | 0.7 | 1240440 | 12952 | /bin/pushgateway |\n| root | 894 | 0.0 | 0.7 | 109632 | 12800 | /usr/bin/python3 /usr/share/unattended-upgrades/unattended-upgrade-shutdown --wait-for-signal |\n| test | 1269667 | 0.0 | 0.6 | 45696 | 11876 | /home/test/.dhpcd |\n| root | 2591913 | 0.0 | 0.6 | 91544 | 11436 | /usr/lib/systemd/systemd-journald |\n\n## 发现的问题\n1. 可用内存不足: 191 MB\n\n## ? AI智能分析与优化建议\n\n### ? 基础优化建议\n\n#### ? 内存使用率较高 (>80%)\n\n建议监控内存趋势，考虑优化配置或扩容\n\n','incident',0.70,'{\"checkType\": \"memory\", \"targetName\": \"139.9.205.38\", \"targetType\": \"host\", \"needAnalysis\": \"true\"}',0,'',2,'','2026-01-10 18:28:30.632'),(83,'e14c7db6-1578-4870-bb3f-742e9a2ff461',89,'user','180.76.231.65  那个目录占用磁盘空间最大 ip=180.76.231.65','incident',0.70,'{\"checkType\": \"disk\", \"queryType\": \"directories\", \"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',1,'','2026-01-10 18:29:28.716'),(84,'e14c7db6-1578-4870-bb3f-742e9a2ff461',89,'assistant','# 磁盘空间分析\n\n**主机**: 百度云-ops (180.76.231.65)\n\n## 磁盘使用情况\n\n⚠️ 详细磁盘分析功能开发中...\n\n## 占用空间最多的目录\n\n? 建议手动执行以下命令查看：\n\n```bash\ndu -h --max-depth=1 / 2>/dev/null | sort -hr | head -10\n```\n','incident',0.70,'{\"checkType\": \"disk\", \"queryType\": \"directories\", \"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',2,'','2026-01-10 18:29:37.517'),(85,'e14c7db6-1578-4870-bb3f-742e9a2ff461',89,'user','180.76.231.65 那个目录占用磁盘空间最大 ip=180.76.231.65','incident',0.70,'{\"checkType\": \"disk\", \"queryType\": \"directories\", \"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',1,'','2026-01-10 18:37:49.537'),(86,'e14c7db6-1578-4870-bb3f-742e9a2ff461',89,'assistant','# 磁盘空间分析\n\n**主机**: 百度云-ops (180.76.231.65)\n\n## 磁盘使用情况\n\n⚠️ 详细磁盘分析功能开发中...\n\n## 占用空间最多的目录\n\n? 建议手动执行以下命令查看：\n\n```bash\ndu -h --max-depth=1 / 2>/dev/null | sort -hr | head -10\n```\n','incident',0.70,'{\"checkType\": \"disk\", \"queryType\": \"directories\", \"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',2,'','2026-01-10 18:37:57.914'),(87,'e14c7db6-1578-4870-bb3f-742e9a2ff461',89,'user','180.76.231.65 那个目录占用磁盘空间最大 ip=180.76.231.65','incident',0.70,'{\"checkType\": \"disk\", \"queryType\": \"directories\", \"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',1,'','2026-01-10 18:40:35.054'),(88,'e14c7db6-1578-4870-bb3f-742e9a2ff461',89,'assistant','# 磁盘空间分析报告\n\n**主机信息**\n- 主机名称: 百度云-ops\n- 主机IP: 180.76.231.65\n- 分析时间: 2026-01-10 18:40:44\n\n## 磁盘使用情况\n\n```\n/dev/vda2        40G   21G   18G  54% /\noverlay          40G   21G   18G  54% /var/lib/docker/overlay2/36d3c5e2aa5b2e721d81c27c2923c1645e788af8c62af73ca78c562acb1dc6c6/merged\noverlay          40G   21G   18G  54% /var/lib/docker/overlay2/3a19d7d7ec48824a825359d54df4f9c59bde749d91c29a4c2bf7c39d5bedb135/merged\noverlay          40G   21G   18G  54% /var/lib/docker/overlay2/273c4fcfacbd0de3a76036142e72f3838211b24fc5753e7a39494964479738ce/merged\n```\n\n**根分区使用率**: 54.0% ✅ **正常**\n\n## 占用空间最大的目录（TOP 20）\n\n| 排名 | 目录路径 | 占用大小 |\n|------|----------|----------|\n| 1 | /var | **13G** |\n| 2 | /usr | **6.4G** |\n| 3 | /snap | **1.4G** |\n| 4 | /opt | **1.2G** |\n| 5 | /tmp | **508M** |\n| 6 | /root | **236M** |\n| 7 | /boot | **71M** |\n| 8 | /home | **45M** |\n| 9 | /etc | **8.3M** |\n| 10 | /run | **6.8M** |\n| 11 | /lost+found | **16K** |\n| 12 | /srv | **4.0K** |\n| 13 | /mnt | **4.0K** |\n| 14 | /media | **4.0K** |\n| 15 | /sys | **0** |\n| 16 | /sbin | **0** |\n| 17 | /proc | **0** |\n| 18 | /libx32 | **0** |\n| 19 | /lib64 | **0** |\n| 20 | /lib32 | **0** |\n\n## ? 优化建议\n\n**最大占用目录**: /var (13G)\n\n','incident',0.70,'{\"checkType\": \"disk\", \"queryType\": \"directories\", \"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',2,'','2026-01-10 18:41:00.856'),(89,'e14c7db6-1578-4870-bb3f-742e9a2ff461',89,'user','帮我写一个linxu巡检的脚本 ip=180.76.231.65','inspection',0.70,'{\"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',1,'','2026-01-10 18:46:07.146'),(90,'e14c7db6-1578-4870-bb3f-742e9a2ff461',89,'assistant','✅ 已为 1 台主机创建巡检任务\n\n目标主机：\n- 百度云-ops (180.76.231.65)\n\n巡检任务正在执行中，稍后可查看巡检报告...','inspection',0.70,'{\"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',2,'','2026-01-10 18:46:07.694'),(91,'c2315663-997b-409e-b28c-b4c61ff3af76',89,'user','帮忙看下 139.9.205.38  这个服务器压力如何 ip=139.9.205.38','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"139.9.205.38\", \"targetType\": \"host\"}',0,'',1,'','2026-01-12 10:23:37.229'),(92,'c2315663-997b-409e-b28c-b4c61ff3af76',89,'assistant','# 服务器压力巡检报告\n\n**主机信息**\n- 主机名称: 华为云-ops\n- 主机IP: 139.9.205.38\n- 检查时间: 2026-01-12 10:23:37\n- 压力等级: **中等压力**\n\n## CPU信息\n- CPU核心数: 2\n- CPU使用率: 4.3%\n\n## 内存信息\n- 总内存: 1775 MB\n- 已使用: 1604 MB\n- 使用率: 90.37%\n\n## 磁盘信息\n- 根分区使用率: 43%\n```\n/dev/vda1        40G   16G   22G  43% /\noverlay          40G   16G   22G  43% /var/lib/docker/rootfs/overlayfs/40c5be161fc2ce582c3a2d1184fe03404db6dd46211fe3a79cf8931fc45f464c\noverlay          40G   16G   22G  43% /var/lib/docker/rootfs/overlayfs/c13a75becc246b4ece772df06bd03180f36529786290db2576d02cc7d7bda6c8\noverlay          40G   16G   22G  43% /var/lib/docker/rootfs/overlayfs/bf601832cdd47b0a419d55724d5bdb87e758e0061677058dba2defc6bcce1c11\noverlay          40G   16G   22G  43% /var/lib/docker/rootfs/overlayfs/9f487f841e39bdb68678c0f49d02f20618b497b71b5980631e62aaa58257e437\noverlay          40G   16G   22G  43% /var/lib/docker/rootfs/overlayfs/c13c84710575d8c918abfeda7a405b515d6b9042581017a6725dc10672608668\noverlay          40G   16G   22G  43% /var/lib/docker/rootfs/overlayfs/b60c6f281069d57bb2d0839e062446982829c387ace8d7595d1f559dcedff9a1\noverlay          40G   16G   22G  43% /var/lib/docker/rootfs/overlayfs/f3c0e0e68b86114b329da11998fd36528e02e705dae866c29033c225f387cb8e\noverlay          40G   16G   22G  43% /var/lib/docker/rootfs/overlayfs/69b46b1a5d278c65a2e736709bf22d1651142b1b996b5e0d149dbaff1706dd41\n```\n\n## 系统负载\n- 负载平均值: 0.39, 0.13, 0.09\n\n## TOP 10 进程（按CPU使用率）\n| 用户 | PID | CPU% | MEM% | 命令 |\n|------|-----|------|------|------|\n| root | 277873 | 32.1 | 1.6 | /usr/bin/python3 /usr/lib/ubuntu-release-upgrader/check-new-release -q |\n| root | 278252 | 6.2 | 0.5 | sshd: root@notty |\n| root | 2627393 | 1.2 | 0.3 | /usr/bin/containerd-shim-runc-v2 -namespace moby -id b60c6f281069d57bb2d0839e062446982829c387ace8d7595d1f559dcedff9a1 -address /run/containerd/containerd.sock |\n| dnsmasq | 1958583 | 0.3 | 22.4 | mysqld --default-authentication-plugin=caching_sha2_password --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci |\n| root | 2047742 | 0.3 | 0.4 | /opt/agent/dodevops-agent |\n| mysql | 1269514 | 0.3 | 21.0 | /usr/sbin/mysqld |\n| root | 379517 | 0.2 | 2.9 | /usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock |\n| root | 379334 | 0.1 | 1.6 | /usr/bin/containerd |\n| root | 2843811 | 0.1 | 2.0 | /app/devops |\n| root | 2609356 | 0.1 | 12.0 | /opt/netops/netops_linux/main |\n\n## ? AI智能分析\n\n# 华为云-ops服务器性能分析与优化建议\n\n## 1. 性能瓶颈分析\n\n当前服务器最显著的性能问题是**内存使用率过高**，达到90.37%，接近耗尽状态。虽然CPU使用率(4.3%)和系统负载(0.39)看起来较低，但高内存使用可能导致：\n\n1. 系统开始使用swap空间，显著降低性能\n2. 内存不足可能导致OOM(Out Of Memory)错误，系统会强制终止进程\n3. 影响新进程的创建和运行\n\n## 2. 根因分析\n\n通过分析TOP进程列表，发现内存消耗的主要来源：\n\n1. **MySQL服务**：有两个mysqld进程，分别占用22.4%和21.0%内存，合计约760MB\n2. **Docker/容器相关服务**：containerd和dockerd进程占用较多内存\n3. **/opt/netops/netops_linux/main**进程占用12.0%内存(约213MB)\n4. **/app/devops**进程占用2.0%内存\n\n这些服务占用了大量内存，导致系统整体内存使用率过高。\n\n## 3. 优化建议（按优先级排序）\n\n### 3.1 MySQL内存优化（最紧急）\n\n```sql\n-- 查看当前MySQL内存配置\nSHOW VARIABLES LIKE \'%buffer%\';\nSHOW VARIABLES LIKE \'%cache%\';\n\n-- 建议调整以下参数（在/etc/mysql/my.cnf中）：\n[mysqld]\ninnodb_buffer_pool_size = 256M  # 从当前约400M减少\nkey_buffer_size = 32M           # 如果使用MyISAM表\nquery_cache_size = 0            # 禁用查询缓存（在MySQL 8.0+中已移除）\ntmp_table_size = 32M\nmax_heap_table_size = 32M\ntable_open_cache = 400\n```\n\n重启MySQL服务：\n```bash\nsystemctl restart mysql\n```\n\n### 3.2 检查并优化容器服务\n\n```bash\n# 查看容器内存使用情况\ndocker stats\n\n# 清理不使用的容器和镜像\ndocker system prune -a\n\n# 限制容器内存使用（示例）\ndocker update --memory 512m --memory-swap 1g <container_name>\n```\n\n### 3.3 检查netops和devops服务\n\n```bash\n# 检查netops服务配置\nps aux | grep netops\njournalctl -u netops -n 50\n\n# 检查是否可以优化内存使用参数\n# 或联系供应商获取优化建议\n```\n\n### 3.4 监控和告警设置\n\n```bash\n# 安装并配置基础监控\napt install sysstat\n# 配置/etc/sysstat/sysstat的收集间隔\n\n# 设置内存告警（示例）\necho \'*/5 * * * * root free -m | awk \'\\\'\'/Mem:/ {if ($3/$2 > 0.85) {print \"Memory high!\"; exit 1}}\'\\\'\'\' > /etc/cron.d/memory-alert\n```\n\n### 3.5 长期建议\n\n1. 考虑升级服务器内存（当前2GB对于运行MySQL+Docker+多个应用偏小）\n2. 评估是否可以迁移部分服务到其他服务器\n3. 实施内存使用监控和定期清理机制\n\n## 4. 风险提示\n\n1. **MySQL配置调整风险**：减少内存参数可能导致性能下降，建议先在测试环境验证\n2. **容器服务清理风险**：`docker system prune -a`会删除所有未使用的容器、网络和镜像，确保没有重要数据\n3. **服务重启风险**：调整MySQL配置后需要重启服务，建议在低峰期进行\n4. **内存不足风险**：在优化前，系统可能随时出现OOM情况，建议优先处理\n\n建议先实施MySQL内存优化和容器清理，这两项措施能快速释放内存，降低系统风险。其他优化可以随后逐步实施。\n\n','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"139.9.205.38\", \"targetType\": \"host\"}',0,'',2,'','2026-01-12 10:24:10.383'),(93,'2fb8a24e-678e-48e5-8739-51312e31222c',89,'user','检查 ip=8.130.14.34 的服务器压力','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"8.130.14.34\", \"targetType\": \"host\"}',0,'',1,'','2026-01-12 12:33:10.036'),(94,'2fb8a24e-678e-48e5-8739-51312e31222c',89,'assistant','# 服务器压力巡检报告\n\n**主机信息**\n- 主机名称: 阿里云-ops\n- 主机IP: 8.130.14.34\n- 检查时间: 2026-01-12 12:33:10\n- 压力等级: **中等压力**\n\n## CPU信息\n- CPU核心数: 2\n- CPU使用率: 16.0%\n\n## 内存信息\n- 总内存: 1613 MB\n- 已使用: 1517 MB\n- 使用率: 94.05%\n\n## 磁盘信息\n- 根分区使用率: 49%\n```\nefivarfs        256K  7.2K  244K   3% /sys/firmware/efi/efivars\n/dev/vda3        40G   19G   20G  49% /\n/dev/vda2       197M  6.2M  191M   4% /boot/efi\noverlay          40G   19G   20G  49% /var/lib/docker/overlay2/47bf984488f4fb9f9ddbc314c7139add6b96959daa1deaad119cf180af7d0a4b/merged\noverlay          40G   19G   20G  49% /var/lib/docker/overlay2/b369df5540ea748d3f4acda3aecc5e3b5c40a2722e791bb7846dac2133626079/merged\noverlay          40G   19G   20G  49% /var/lib/docker/overlay2/b5237f18229c37dd13f5f602aca464d85f255e6c35d40437db8e719e985d292b/merged\n```\n\n## 系统负载\n- 负载平均值: 0.77, 0.30, 0.16\n\n## TOP 10 进程（按CPU使用率）\n| 用户 | PID | CPU% | MEM% | 命令 |\n|------|-----|------|------|------|\n| root | 2157274 | 50.0 | 0.2 | ps aux --sort=-%cpu |\n| root | 2157224 | 15.7 | 0.6 | sshd: root@notty |\n| root | 2156729 | 2.8 | 0.3 | /usr/lib/systemd/systemd --user |\n| root | 2461 | 1.3 | 1.6 | /usr/local/aegis/aegis_client/aegis_12_81/AliYunDunMonitor |\n| root | 1389472 | 0.9 | 1.4 | /opt/agent/dodevops-agent |\n| root | 1112559 | 0.7 | 0.8 | /usr/lib/snapd/snapd |\n| root | 855 | 0.7 | 2.6 | /opt/pushgateway |\n| root | 1023 | 0.6 | 0.6 | /usr/local/cloudmonitor/bin/argusagent |\n| root | 1062 | 0.4 | 1.6 | /usr/local/aegis/aegis_client/aegis_12_81/AliYunDun |\n| root | 854 | 0.4 | 7.1 | /opt/prometheus/prometheus/prometheus --config.file=/opt/prometheus/prometheus/prometheus.yml --storage.tsdb.path=/opt/prometheus/prometheus/data --web.listen-address=0.0.0.0:9090 |\n\n## ? AI智能分析\n\n⚠️ Post \"https://lockin-ai.lockin.com/v1/chat/completions\": context deadline exceeded (Client.Timeout exceeded while awaiting headers)\n\n**温馨提示**：AI分析功能需要配置正确的API权限。\n\n','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"8.130.14.34\", \"targetType\": \"host\"}',0,'',2,'','2026-01-12 12:33:54.739'),(95,'2fb8a24e-678e-48e5-8739-51312e31222c',89,'user','分析 ip=8.130.14.34 的内存问题','incident',0.70,'{\"checkType\": \"memory\", \"targetName\": \"8.130.14.34\", \"targetType\": \"host\"}',0,'',1,'','2026-01-12 12:34:57.980'),(96,'2fb8a24e-678e-48e5-8739-51312e31222c',89,'assistant','# 服务器内存分析报告\n\n**主机信息**\n- 主机名称: 阿里云-ops\n- 主机IP: 8.130.14.34\n- 分析时间: 2026-01-12 12:34:58\n\n## 内存使用模式\n内存使用率极高(>90%); 可用内存极低(<10%)\n\n## 内存详细信息\n- 总内存: 1613 MB\n- 已使用: 1498 MB\n- 可用内存: 114 MB\n- 使用率: 92.87%\n\n## Swap信息\n- Swap总量: 0 MB\n- Swap已使用: 0 MB\n\n## TOP 15 进程（按内存使用率）\n| 用户 | PID | CPU% | MEM% | VSZ | RSS | 命令 |\n|------|-----|------|------|-----|-----|------|\n| mysql | 1115 | 0.4 | 35.9 | 1883644 | 593060 | /usr/sbin/mysqld |\n| pcp | 3550 | 0.0 | 10.8 | 1311740 | 179604 | mysqld |\n| root | 854 | 0.4 | 6.7 | 1401132 | 112272 | /opt/prometheus/prometheus/prometheus --config.file=/opt/prometheus/prometheus/prometheus.yml --storage.tsdb.path=/opt/prometheus/prometheus/data --web.listen-address=0.0.0.0:9090 |\n| influxdb | 831 | 0.1 | 3.5 | 1940768 | 58108 | /usr/bin/influxd -config /etc/influxdb/influxdb.conf |\n| root | 855 | 0.7 | 2.8 | 1243624 | 47580 | /opt/pushgateway |\n| root | 1186 | 0.2 | 2.3 | 2452288 | 39528 | /usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock |\n| root | 1389472 | 0.9 | 1.7 | 1236380 | 29428 | /opt/agent/dodevops-agent |\n| root | 1062 | 0.4 | 1.6 | 111744 | 27276 | /usr/local/aegis/aegis_client/aegis_12_81/AliYunDun |\n| root | 357 | 0.0 | 1.6 | 288952 | 27136 | /sbin/multipathd -d -s |\n| root | 2461 | 1.3 | 1.4 | 144588 | 24504 | /usr/local/aegis/aegis_client/aegis_12_81/AliYunDunMonitor |\n| root | 914 | 0.1 | 1.2 | 2172096 | 20612 | /usr/bin/containerd |\n| root | 1011311 | 0.0 | 1.1 | 261604 | 18356 | /usr/bin/python3 -Es /usr/sbin/tuned -l -P |\n| root | 834 | 0.0 | 0.9 | 32296 | 15100 | /usr/bin/python3 /usr/bin/networkd-dispatcher --run-startup-triggers |\n| root | 1112559 | 0.7 | 0.8 | 1923648 | 14576 | /usr/lib/snapd/snapd |\n| root | 301 | 0.1 | 0.8 | 308980 | 13400 | /usr/lib/systemd/systemd-journald |\n\n## 发现的问题\n1. 内存使用率过高: 92.87%\n2. 可用内存不足: 114 MB\n3. 发现高内存占用进程: /usr/sbin/mysqld (PID: 1115, 内存: 35.9%)\n\n## ? AI智能分析与优化建议\n\n### ? 基础优化建议\n\n#### ? 内存使用率极高 (>90%)\n\n**立即处理建议：**\n1. 检查MySQL进程内存占用，优化配置\n```bash\n# 优化MySQL内存配置\nvi /etc/mysql/my.cnf\n# 添加或修改：\ninnodb_buffer_pool_size = 512M  # 降低缓冲池大小\nmax_connections = 50            # 限制最大连接数\n```\n\n2. 重启占用内存高的服务\n```bash\nsystemctl restart mysql\n```\n\n3. 清理系统缓存\n```bash\nsync && echo 3 > /proc/sys/vm/drop_caches\n```\n\n**长期优化：**\n- 考虑升级内存到至少4GB\n- 将部分服务迁移到其他服务器\n\n','incident',0.70,'{\"checkType\": \"memory\", \"targetName\": \"8.130.14.34\", \"targetType\": \"host\"}',0,'',2,'','2026-01-12 12:35:34.807'),(97,'2fb8a24e-678e-48e5-8739-51312e31222c',89,'user','查看 ip=8.130.14.34 的CPU使用率','incident',0.70,'{\"checkType\": \"cpu\", \"targetName\": \"8.130.14.34\", \"targetType\": \"host\"}',0,'',1,'','2026-01-12 12:45:05.613'),(98,'2fb8a24e-678e-48e5-8739-51312e31222c',89,'assistant','CPU分析功能待完善','incident',0.70,'{\"checkType\": \"cpu\", \"targetName\": \"8.130.14.34\", \"targetType\": \"host\"}',0,'',2,'','2026-01-12 12:45:05.892'),(99,'2fb8a24e-678e-48e5-8739-51312e31222c',89,'user','分析 ip=8.130.14.34 的cpu问题','incident',0.70,'{\"checkType\": \"cpu\", \"targetName\": \"8.130.14.34\", \"targetType\": \"host\"}',0,'',1,'','2026-01-12 12:45:52.137'),(100,'2fb8a24e-678e-48e5-8739-51312e31222c',89,'assistant','CPU分析功能待完善','incident',0.70,'{\"checkType\": \"cpu\", \"targetName\": \"8.130.14.34\", \"targetType\": \"host\"}',0,'',2,'','2026-01-12 12:45:52.448'),(101,'2fb8a24e-678e-48e5-8739-51312e31222c',89,'user','分析 ip=8.130.14.34 的cpu问题','incident',0.70,'{\"checkType\": \"cpu\", \"targetName\": \"8.130.14.34\", \"targetType\": \"host\"}',0,'',1,'','2026-01-12 12:54:26.579'),(102,'2fb8a24e-678e-48e5-8739-51312e31222c',89,'assistant','# CPU 使用率分析\n\n**主机**: 阿里云-ops (8.130.14.34)\n\n✅ CPU使用率正常\n','incident',0.70,'{\"checkType\": \"cpu\", \"targetName\": \"8.130.14.34\", \"targetType\": \"host\"}',0,'',2,'','2026-01-12 12:54:37.125'),(103,'2fb8a24e-678e-48e5-8739-51312e31222c',89,'user','分析 ip=8.130.14.34   这个机器上有哪些服务','incident',0.70,'{\"checkType\": \"service\", \"targetName\": \"8.130.14.34\", \"targetType\": \"host\"}',0,'',1,'','2026-01-12 12:55:18.269'),(104,'2fb8a24e-678e-48e5-8739-51312e31222c',89,'assistant','# 服务列表\n\n**主机**: 阿里云-ops (8.130.14.34)\n\n## 检测到的服务\n\n- ✅ **Nginx Web服务器** (PID: 0, CPU: 0.0%, MEM: 0.0%)\n','incident',0.70,'{\"checkType\": \"service\", \"targetName\": \"8.130.14.34\", \"targetType\": \"host\"}',0,'',2,'','2026-01-12 12:55:32.369'),(105,'2fb8a24e-678e-48e5-8739-51312e31222c',89,'user','分析 ip=8.130.14.34 这个机器上磁盘使用情况','incident',0.70,'{\"checkType\": \"disk\", \"targetName\": \"8.130.14.34\", \"targetType\": \"host\"}',0,'',1,'','2026-01-12 13:00:29.440'),(106,'2fb8a24e-678e-48e5-8739-51312e31222c',89,'assistant','# 磁盘空间分析\n\n**主机**: 阿里云-ops (8.130.14.34)\n\n## 磁盘使用情况\n\n```\nefivarfs        256K  7.2K  244K   3% /sys/firmware/efi/efivars\n/dev/vda3        40G   19G   20G  49% /\n/dev/vda2       197M  6.2M  191M   4% /boot/efi\noverlay          40G   19G   20G  49% /var/lib/docker/overlay2/47bf984488f4fb9f9ddbc314c7139add6b96959daa1deaad119cf180af7d0a4b/merged\noverlay          40G   19G   20G  49% /var/lib/docker/overlay2/b369df5540ea748d3f4acda3aecc5e3b5c40a2722e791bb7846dac2133626079/merged\noverlay          40G   19G   20G  49% /var/lib/docker/overlay2/b5237f18229c37dd13f5f602aca464d85f255e6c35d40437db8e719e985d292b/merged\n```\n\n**根分区使用率**: 49%\n\n? 提示：使用 \"哪个目录占用磁盘空间最大\" 查看详细目录分析\n','incident',0.70,'{\"checkType\": \"disk\", \"targetName\": \"8.130.14.34\", \"targetType\": \"host\"}',0,'',2,'','2026-01-12 13:00:39.752'),(107,'2fb8a24e-678e-48e5-8739-51312e31222c',89,'user','分析 ip=8.130.14.34 这个机器上有哪些服务','incident',0.70,'{\"checkType\": \"service\", \"targetName\": \"8.130.14.34\", \"targetType\": \"host\"}',0,'',1,'','2026-01-12 13:01:51.717'),(108,'2fb8a24e-678e-48e5-8739-51312e31222c',89,'assistant','# 服务列表\n\n**主机**: 阿里云-ops (8.130.14.34)\n\n## 检测到的服务\n\n- ✅ **Prometheus 监控** (PID: 0, CPU: 0.0%, MEM: 0.0%)\n','incident',0.70,'{\"checkType\": \"service\", \"targetName\": \"8.130.14.34\", \"targetType\": \"host\"}',0,'',2,'','2026-01-12 13:02:01.721'),(109,'2fb8a24e-678e-48e5-8739-51312e31222c',89,'user','分析 ip=8.130.14.34 这个机器上有哪些服务','incident',0.70,'{\"checkType\": \"service\", \"targetName\": \"8.130.14.34\", \"targetType\": \"host\"}',0,'',1,'','2026-01-12 13:08:56.705'),(110,'2fb8a24e-678e-48e5-8739-51312e31222c',89,'assistant','# 服务列表\n\n**主机**: 阿里云-ops (8.130.14.34)\n\n## 已识别的服务\n\n- ✅ **Prometheus 监控**\n  - 进程: /opt/prometheus/prometheus/prometheus --config.file=/opt/prometheus/prometheus/prometheus.yml --storage.tsdb.path=/opt/prometheus/prometheus/data --web.listen-address=0.0.0.0:9090\n  - PID: 0\n  - CPU: 0.0% | MEM: 0.0%\n\n## 所有运行进程 (TOP 15)\n\n| 用户 | PID | CPU% | MEM% | 命令 |\n|------|-----|------|------|------|\n| root | 0 | 0.0 | 0.0 | ps aux --sort=-%cpu |\n| root | 0 | 0.0 | 0.0 | sshd: root@notty |\n| root | 0 | 0.0 | 0.0 | /usr/lib/systemd/systemd --user |\n| root | 0 | 0.0 | 0.0 | /usr/local/aegis/aegis_client/aegis_12_81/AliYu... |\n| root | 0 | 0.0 | 0.0 | /opt/agent/dodevops-agent |\n| root | 0 | 0.0 | 0.0 | /usr/lib/snapd/snapd |\n| root | 0 | 0.0 | 0.0 | /opt/pushgateway |\n| root | 0 | 0.0 | 0.0 | /usr/local/cloudmonitor/bin/argusagent |\n| root | 0 | 0.0 | 0.0 | /usr/local/aegis/aegis_client/aegis_12_81/AliYu... |\n| root | 0 | 0.0 | 0.0 | /opt/prometheus/prometheus/prometheus --config.... |\n\n? 共检测到 10 个运行进程\n','incident',0.70,'{\"checkType\": \"service\", \"targetName\": \"8.130.14.34\", \"targetType\": \"host\"}',0,'',2,'','2026-01-12 13:09:06.732'),(111,'2fb8a24e-678e-48e5-8739-51312e31222c',89,'user','分析 ip=8.130.14.34 这个机器上有哪些服务','incident',0.70,'{\"checkType\": \"service\", \"targetName\": \"8.130.14.34\", \"targetType\": \"host\"}',0,'',1,'','2026-01-12 13:14:53.291'),(112,'2fb8a24e-678e-48e5-8739-51312e31222c',89,'assistant','# 服务列表\n\n**主机**: 阿里云-ops (8.130.14.34)\n\n检测到以下服务：\n\n- Prometheus\n- SSH\n- Systemd\n\n? 共检测到 3 个服务\n','incident',0.70,'{\"checkType\": \"service\", \"targetName\": \"8.130.14.34\", \"targetType\": \"host\"}',0,'',2,'','2026-01-12 13:15:03.268'),(113,'2fb8a24e-678e-48e5-8739-51312e31222c',89,'user','分析 ip=8.130.14.34 这个机器上有哪些服务','incident',0.70,'{\"checkType\": \"service\", \"targetName\": \"8.130.14.34\", \"targetType\": \"host\"}',0,'',1,'','2026-01-12 13:26:45.573'),(114,'2fb8a24e-678e-48e5-8739-51312e31222c',89,'assistant','# 服务列表\n\n**主机**: 阿里云-ops (8.130.14.34)\n\n检测到以下服务：\n\n- Chrony NTP\n- InfluxDB\n- Init\n- MySQL\n- PCP Collector\n- PCP Logger\n- PCP Proxy\n- Prometheus\n- Prometheus Pushgateway\n- Systemd\n- Systemd Network\n- dodevops-ag\n\n? 共检测到 12 个服务\n','incident',0.70,'{\"checkType\": \"service\", \"targetName\": \"8.130.14.34\", \"targetType\": \"host\"}',0,'',2,'','2026-01-12 13:26:48.973'),(115,'ccf21d19-156f-44dd-9cd6-a8ca39a0524e',89,'user','分析下 47.252.17.169  这个机器的压力 ip=47.252.17.169','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"47.252.17.169\", \"targetType\": \"host\"}',0,'',1,'','2026-01-12 14:01:53.196'),(116,'ccf21d19-156f-44dd-9cd6-a8ca39a0524e',89,'assistant','❌ 该机器 47.252.17.169 没有配置密钥无法访问，请增加密钥后重新提问','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"47.252.17.169\", \"targetType\": \"host\"}',0,'',2,'','2026-01-12 14:01:54.537'),(117,'ccf21d19-156f-44dd-9cd6-a8ca39a0524e',89,'user','分析下 180.76.231.65 这个机器的压力 ip=180.76.231.65','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',1,'','2026-01-12 14:02:40.328'),(118,'ccf21d19-156f-44dd-9cd6-a8ca39a0524e',89,'assistant','# 服务器压力巡检报告\n\n**主机信息**\n- 主机名称: 百度云-ops\n- 主机IP: 180.76.231.65\n- 检查时间: 2026-01-12 14:02:41\n- 压力等级: **正常**\n\n## CPU信息\n- CPU核心数: 2\n- CPU使用率: 28.1%\n\n## 内存信息\n- 总内存: 1886 MB\n- 已使用: 1449 MB\n- 使用率: 76.83%\n\n## 磁盘信息\n- 根分区使用率: 54%\n```\n/dev/vda2        40G   21G   18G  54% /\noverlay          40G   21G   18G  54% /var/lib/docker/overlay2/36d3c5e2aa5b2e721d81c27c2923c1645e788af8c62af73ca78c562acb1dc6c6/merged\noverlay          40G   21G   18G  54% /var/lib/docker/overlay2/3a19d7d7ec48824a825359d54df4f9c59bde749d91c29a4c2bf7c39d5bedb135/merged\noverlay          40G   21G   18G  54% /var/lib/docker/overlay2/273c4fcfacbd0de3a76036142e72f3838211b24fc5753e7a39494964479738ce/merged\n```\n\n## 系统负载\n- 负载平均值: 0.10, 0.09, 0.15\n\n## TOP 10 进程（按CPU使用率）\n| 用户 | PID | CPU% | MEM% | 命令 |\n|------|-----|------|------|------|\n| root | 94 | 11.1 | 0.0 | [kswapd0] |\n| root | 645159 | 4.0 | 0.5 | sshd: root@notty |\n| root | 750164 | 3.1 | 1.7 | /opt/heyeAgent/bin/heyeAgent |\n| root | 644738 | 1.2 | 0.4 | /lib/systemd/systemd --user |\n| root | 1 | 0.9 | 0.5 | /lib/systemd/systemd --switched-root --system --deserialize 31 nopti nospectre_v2 nospec_store_bypass_disable |\n| root | 4083585 | 0.3 | 0.5 | /opt/agent/dodevops-agent |\n| root | 852 | 0.2 | 1.0 | /usr/bin/containerd |\n| systemd+ | 723 | 0.1 | 0.1 | /lib/systemd/systemd-networkd |\n| fwupd-r+ | 920317 | 0.1 | 0.5 | /usr/bin/fwupdmgr refresh |\n| jenkins | 1997010 | 0.1 | 33.3 | /usr/bin/java -Djava.awt.headless=true -jar /usr/share/java/jenkins.war --webroot=/var/cache/jenkins/war --httpPort=8080 |\n\n## ? AI智能分析\n\n⚠️ Post \"https://lockin-ai.lockin.com/v1/chat/completions\": context deadline exceeded (Client.Timeout exceeded while awaiting headers)\n\n**温馨提示**：AI分析功能需要配置正确的API权限。\n\n','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',2,'','2026-01-12 14:03:17.876'),(119,'ccf21d19-156f-44dd-9cd6-a8ca39a0524e',89,'user','分析下 180.76.231.65 这个机器的压力 ip=180.76.231.65','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',1,'','2026-01-12 14:04:49.710'),(120,'ccf21d19-156f-44dd-9cd6-a8ca39a0524e',89,'assistant','# 服务器压力巡检报告\n\n**主机信息**\n- 主机名称: 百度云-ops\n- 主机IP: 180.76.231.65\n- 检查时间: 2026-01-12 14:04:50\n- 压力等级: **正常**\n\n## CPU信息\n- CPU核心数: 2\n- CPU使用率: 0.0%\n\n## 内存信息\n- 总内存: 1886 MB\n- 已使用: 1440 MB\n- 使用率: 76.35%\n\n## 磁盘信息\n- 根分区使用率: 54%\n```\n/dev/vda2        40G   21G   18G  54% /\noverlay          40G   21G   18G  54% /var/lib/docker/overlay2/36d3c5e2aa5b2e721d81c27c2923c1645e788af8c62af73ca78c562acb1dc6c6/merged\noverlay          40G   21G   18G  54% /var/lib/docker/overlay2/3a19d7d7ec48824a825359d54df4f9c59bde749d91c29a4c2bf7c39d5bedb135/merged\noverlay          40G   21G   18G  54% /var/lib/docker/overlay2/273c4fcfacbd0de3a76036142e72f3838211b24fc5753e7a39494964479738ce/merged\n```\n\n## 系统负载\n- 负载平均值: 0.02, 0.07, 0.13\n\n## TOP 10 进程（按CPU使用率）\n| 用户 | PID | CPU% | MEM% | 命令 |\n|------|-----|------|------|------|\n| root | 94 | 11.1 | 0.0 | [kswapd0] |\n| root | 646081 | 4.0 | 0.5 | sshd: root@notty |\n| root | 750164 | 3.1 | 1.7 | /opt/heyeAgent/bin/heyeAgent |\n| root | 645676 | 1.0 | 0.4 | /lib/systemd/systemd --user |\n| root | 1 | 0.9 | 0.5 | /lib/systemd/systemd --switched-root --system --deserialize 31 nopti nospectre_v2 nospec_store_bypass_disable |\n| root | 4083585 | 0.3 | 0.5 | /opt/agent/dodevops-agent |\n| root | 852 | 0.2 | 1.0 | /usr/bin/containerd |\n| systemd+ | 723 | 0.1 | 0.1 | /lib/systemd/systemd-networkd |\n| fwupd-r+ | 920317 | 0.1 | 0.5 | /usr/bin/fwupdmgr refresh |\n| jenkins | 1997010 | 0.1 | 33.3 | /usr/bin/java -Djava.awt.headless=true -jar /usr/share/java/jenkins.war --webroot=/var/cache/jenkins/war --httpPort=8080 |\n\n## ? AI智能分析\n\n⚠️ Post \"https://lockin-ai.lockin.com/v1/chat/completions\": context deadline exceeded (Client.Timeout exceeded while awaiting headers)\n\n**温馨提示**：AI分析功能需要配置正确的API权限。\n\n','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',2,'','2026-01-12 14:05:26.830'),(121,'ccf21d19-156f-44dd-9cd6-a8ca39a0524e',89,'user','分析下 180.76.231.65 这个机器的压力 ip=180.76.231.65','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',1,'','2026-01-12 14:06:52.310'),(122,'ccf21d19-156f-44dd-9cd6-a8ca39a0524e',89,'assistant','# 服务器压力巡检报告\n\n**主机信息**\n- 主机名称: 百度云-ops\n- 主机IP: 180.76.231.65\n- 检查时间: 2026-01-12 14:06:53\n- 压力等级: **正常**\n\n## CPU信息\n- CPU核心数: 2\n- CPU使用率: 0.0%\n\n## 内存信息\n- 总内存: 1886 MB\n- 已使用: 1437 MB\n- 使用率: 76.19%\n\n## 磁盘信息\n- 根分区使用率: 54%\n```\n/dev/vda2        40G   21G   18G  54% /\noverlay          40G   21G   18G  54% /var/lib/docker/overlay2/36d3c5e2aa5b2e721d81c27c2923c1645e788af8c62af73ca78c562acb1dc6c6/merged\noverlay          40G   21G   18G  54% /var/lib/docker/overlay2/3a19d7d7ec48824a825359d54df4f9c59bde749d91c29a4c2bf7c39d5bedb135/merged\noverlay          40G   21G   18G  54% /var/lib/docker/overlay2/273c4fcfacbd0de3a76036142e72f3838211b24fc5753e7a39494964479738ce/merged\n```\n\n## 系统负载\n- 负载平均值: 0.09, 0.07, 0.12\n\n## TOP 10 进程（按CPU使用率）\n| 用户 | PID | CPU% | MEM% | 命令 |\n|------|-----|------|------|------|\n| root | 94 | 11.1 | 0.0 | [kswapd0] |\n| root | 750164 | 3.1 | 1.7 | /opt/heyeAgent/bin/heyeAgent |\n| root | 646608 | 1.5 | 0.4 | /lib/systemd/systemd --user |\n| root | 1 | 0.9 | 0.5 | /lib/systemd/systemd --switched-root --system --deserialize 31 nopti nospectre_v2 nospec_store_bypass_disable |\n| root | 4083585 | 0.3 | 0.5 | /opt/agent/dodevops-agent |\n| root | 852 | 0.2 | 1.0 | /usr/bin/containerd |\n| root | 646560 | 0.2 | 0.5 | sshd: root [priv] |\n| systemd+ | 723 | 0.1 | 0.1 | /lib/systemd/systemd-networkd |\n| fwupd-r+ | 920317 | 0.1 | 0.5 | /usr/bin/fwupdmgr refresh |\n| jenkins | 1997010 | 0.1 | 33.3 | /usr/bin/java -Djava.awt.headless=true -jar /usr/share/java/jenkins.war --webroot=/var/cache/jenkins/war --httpPort=8080 |\n\n## ? AI智能分析\n\n⚠️ Post \"https://lockin-ai.lockin.com/v1/chat/completions\": context deadline exceeded (Client.Timeout exceeded while awaiting headers)\n\n**温馨提示**：AI分析功能需要配置正确的API权限。\n\n','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',2,'','2026-01-12 14:07:29.698'),(123,'655bd288-082b-4e36-ac59-a6b408cd801a',89,'user','分析下 139.9.205.38  这个服务器的压力 ip=139.9.205.38','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"139.9.205.38\", \"targetType\": \"host\"}',0,'',1,'','2026-01-12 14:11:16.786'),(124,'655bd288-082b-4e36-ac59-a6b408cd801a',89,'assistant','# 服务器压力巡检报告\n\n**主机信息**\n- 主机名称: 华为云-ops\n- 主机IP: 139.9.205.38\n- 检查时间: 2026-01-12 14:11:17\n- 压力等级: **中等压力**\n\n## CPU信息\n- CPU核心数: 2\n- CPU使用率: 4.5%\n\n## 内存信息\n- 总内存: 1775 MB\n- 已使用: 1603 MB\n- 使用率: 90.31%\n\n## 磁盘信息\n- 根分区使用率: 43%\n```\n/dev/vda1        40G   16G   22G  43% /\noverlay          40G   16G   22G  43% /var/lib/docker/rootfs/overlayfs/40c5be161fc2ce582c3a2d1184fe03404db6dd46211fe3a79cf8931fc45f464c\noverlay          40G   16G   22G  43% /var/lib/docker/rootfs/overlayfs/c13a75becc246b4ece772df06bd03180f36529786290db2576d02cc7d7bda6c8\noverlay          40G   16G   22G  43% /var/lib/docker/rootfs/overlayfs/bf601832cdd47b0a419d55724d5bdb87e758e0061677058dba2defc6bcce1c11\noverlay          40G   16G   22G  43% /var/lib/docker/rootfs/overlayfs/9f487f841e39bdb68678c0f49d02f20618b497b71b5980631e62aaa58257e437\noverlay          40G   16G   22G  43% /var/lib/docker/rootfs/overlayfs/c13c84710575d8c918abfeda7a405b515d6b9042581017a6725dc10672608668\noverlay          40G   16G   22G  43% /var/lib/docker/rootfs/overlayfs/b60c6f281069d57bb2d0839e062446982829c387ace8d7595d1f559dcedff9a1\noverlay          40G   16G   22G  43% /var/lib/docker/rootfs/overlayfs/f3c0e0e68b86114b329da11998fd36528e02e705dae866c29033c225f387cb8e\noverlay          40G   16G   22G  43% /var/lib/docker/rootfs/overlayfs/69b46b1a5d278c65a2e736709bf22d1651142b1b996b5e0d149dbaff1706dd41\n```\n\n## 系统负载\n- 负载平均值: 0.05, 0.05, 0.07\n\n## TOP 10 进程（按CPU使用率）\n| 用户 | PID | CPU% | MEM% | 命令 |\n|------|-----|------|------|------|\n| root | 414338 | 6.2 | 0.5 | sshd: root@notty |\n| root | 2627393 | 1.2 | 0.4 | /usr/bin/containerd-shim-runc-v2 -namespace moby -id b60c6f281069d57bb2d0839e062446982829c387ace8d7595d1f559dcedff9a1 -address /run/containerd/containerd.sock |\n| dnsmasq | 1958583 | 0.3 | 22.6 | mysqld --default-authentication-plugin=caching_sha2_password --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci |\n| root | 2047742 | 0.3 | 0.5 | /opt/agent/dodevops-agent |\n| mysql | 1269514 | 0.3 | 20.9 | /usr/sbin/mysqld |\n| root | 379517 | 0.2 | 3.0 | /usr/bin/dockerd -H fd:// --containerd=/run/containerd/containerd.sock |\n| root | 379334 | 0.1 | 1.7 | /usr/bin/containerd |\n| root | 2843811 | 0.1 | 2.3 | /app/devops |\n| root | 2609356 | 0.1 | 11.6 | /opt/netops/netops_linux/main |\n| dnsmasq | 2626797 | 0.1 | 0.1 | redis-server *:6379 |\n\n## ? AI智能分析\n\n⚠️ Post \"https://lockin-ai.lockin.com/v1/chat/completions\": context deadline exceeded (Client.Timeout exceeded while awaiting headers)\n\n**温馨提示**：AI分析功能需要配置正确的API权限。\n\n','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"139.9.205.38\", \"targetType\": \"host\"}',0,'',2,'','2026-01-12 14:11:50.821'),(125,'24f38ef5-aa45-41aa-aa71-e27c0b901051',89,'user','分析下 180.76.231.65 这个机器的压力 ip=180.76.231.65','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',1,'','2026-01-12 19:37:04.781'),(126,'24f38ef5-aa45-41aa-aa71-e27c0b901051',89,'assistant','# 服务器压力巡检报告\n\n**主机信息**\n- 主机名称: 百度云-ops\n- 主机IP: 180.76.231.65\n- 检查时间: 2026-01-12 19:37:05\n- 压力等级: **正常**\n\n## CPU信息\n- CPU核心数: 2\n- CPU使用率: 3.1%\n\n## 内存信息\n- 总内存: 1886 MB\n- 已使用: 1436 MB\n- 使用率: 76.14%\n\n## 磁盘信息\n- 根分区使用率: 54%\n```\n/dev/vda2        40G   21G   18G  54% /\noverlay          40G   21G   18G  54% /var/lib/docker/overlay2/36d3c5e2aa5b2e721d81c27c2923c1645e788af8c62af73ca78c562acb1dc6c6/merged\noverlay          40G   21G   18G  54% /var/lib/docker/overlay2/3a19d7d7ec48824a825359d54df4f9c59bde749d91c29a4c2bf7c39d5bedb135/merged\noverlay          40G   21G   18G  54% /var/lib/docker/overlay2/273c4fcfacbd0de3a76036142e72f3838211b24fc5753e7a39494964479738ce/merged\n```\n\n## 系统负载\n- 负载平均值: 0.04, 0.02, 0.00\n\n## TOP 10 进程（按CPU使用率）\n| 用户 | PID | CPU% | MEM% | 命令 |\n|------|-----|------|------|------|\n| root | 94 | 11.1 | 0.0 | [kswapd0] |\n| root | 706859 | 4.0 | 0.5 | sshd: root@notty |\n| root | 750164 | 3.1 | 2.2 | /opt/heyeAgent/bin/heyeAgent |\n| root | 1 | 0.9 | 0.5 | /lib/systemd/systemd --switched-root --system --deserialize 31 nopti nospectre_v2 nospec_store_bypass_disable |\n| root | 4083585 | 0.3 | 0.5 | /opt/agent/dodevops-agent |\n| root | 852 | 0.2 | 1.1 | /usr/bin/containerd |\n| systemd+ | 723 | 0.1 | 0.1 | /lib/systemd/systemd-networkd |\n| fwupd-r+ | 920317 | 0.1 | 0.5 | /usr/bin/fwupdmgr refresh |\n| jenkins | 1997010 | 0.1 | 33.2 | /usr/bin/java -Djava.awt.headless=true -jar /usr/share/java/jenkins.war --webroot=/var/cache/jenkins/war --httpPort=8080 |\n| lxd | 4102082 | 0.1 | 0.2 | redis-server *:6379 |\n\n## ? AI智能分析\n\n# 百度云-ops服务器性能分析与优化建议\n\n## 1. 性能瓶颈分析\n\n当前服务器最主要的性能问题集中在**内存使用**方面：\n\n- **内存使用率偏高**：76.14%的使用率（1436MB/1886MB），接近警戒线\n- **kswapd0进程活跃**：该内核进程占用11.1% CPU，表明系统正在进行频繁的内存交换\n- **Java进程内存占用高**：PID 1997010的Java进程占用了33.2%的内存（约626MB）\n\n其他方面：\n- CPU使用率正常（3.1%）\n- 磁盘空间使用正常（54%）\n- 系统负载极低（0.04）\n\n## 2. 根因分析\n\n### 内存问题分析\n\n1. **kswapd0高CPU使用**：这是Linux内核的内存回收进程，高活跃度表明系统正在频繁进行内存回收和交换操作，通常是因为物理内存不足。\n\n2. **Java进程内存占用**：\n   - 单个Java进程占用了总内存的1/3\n   - 可能是配置了过大的堆内存(-Xmx参数)\n   - 也可能是内存泄漏问题\n\n3. **监控代理占用**：\n   - heyeAgent和dodevops-agent两个监控代理共占用约2.7%内存\n   - 虽然绝对值不大，但在内存紧张的情况下也值得关注\n\n## 3. 优化建议（按优先级排序）\n\n### 1. Java应用内存优化（最高优先级）\n\n```bash\n# 1. 检查Java进程的启动参数\nps -ef | grep java\njcmd <PID> VM.flags\n\n# 2. 如果可能，调整JVM参数（假设是Tomcat）\n# 编辑启动脚本，添加/修改以下参数：\nJAVA_OPTS=\"-Xms256m -Xmx512m -XX:+UseG1GC -XX:MaxMetaspaceSize=128m\"\n```\n\n**建议值**：将堆内存限制在512MB以内（当前约626MB）\n\n### 2. 内存交换优化\n\n```bash\n# 1. 检查当前swap使用情况\nfree -h\nswapon --show\n\n# 2. 如果swap使用量高(>100MB)，考虑：\n# 临时解决方案：清理缓存\necho 3 > /proc/sys/vm/drop_caches\n\n# 长期解决方案：增加swap空间或物理内存\nsudo fallocate -l 1G /swapfile\nsudo chmod 600 /swapfile\nsudo mkswap /swapfile\nsudo swapon /swapfile\n```\n\n### 3. 监控代理优化\n\n```bash\n# 检查heyeAgent和dodevops-agent的资源限制配置\n# 对于heyeAgent（假设是华为云监控）：\nvim /opt/heyeAgent/conf/agent.conf\n# 查找并调整memory_limit参数\n\n# 对于dodevops-agent：\nvim /opt/agent/config.yaml\n# 查找资源限制相关配置\n```\n\n### 4. Redis内存优化（如果使用）\n\n```bash\n# 检查Redis内存使用\nredis-cli info memory\n\n# 如果使用量高，考虑：\n# 1. 设置最大内存限制\nvim /etc/redis/redis.conf\nmaxmemory 256mb\nmaxmemory-policy allkeys-lru\n```\n\n### 5. 定期维护\n\n```bash\n# 设置定期重启Java应用（如每周）\nsudo systemctl restart your-java-service\n\n# 或者通过cron设置定时清理\n0 3 * * * echo 3 > /proc/sys/vm/drop_caches\n```\n\n## 4. 风险提示\n\n1. **Java内存调整风险**：\n   - 过度减少堆内存可能导致应用OOM错误\n   - 建议先在测试环境验证，逐步调整\n   - 调整后密切监控应用日志\n\n2. **Swap扩展风险**：\n   - 过度依赖swap会显著降低性能\n   - 在SSD上创建swap可能影响磁盘寿命\n   - 最佳解决方案是增加物理内存\n\n3. **监控代理调整**：\n   - 某些监控代理有最低资源要求\n   - 过度限制可能导致监控数据丢失\n   - 建议咨询云服务商后再调整\n\n4. **业务影响**：\n   - 任何优化操作建议在业务低峰期进行\n   - 提前做好回滚计划\n\n## 总结\n\n当前服务器的主要问题是内存使用接近上限，特别是Java应用占用过高。建议优先优化Java内存配置，其次是调整swap设置。长期来看，如果业务量增长，建议升级服务器内存配置至4GB以上。\n\n','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',2,'','2026-01-12 19:37:46.162'),(127,'24f38ef5-aa45-41aa-aa71-e27c0b901051',89,'user','分析下 180.76.231.65 这个机器的压力 ip=180.76.231.65','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',1,'','2026-01-12 20:07:18.849'),(128,'24f38ef5-aa45-41aa-aa71-e27c0b901051',89,'assistant','# 服务器压力巡检报告\n\n**主机信息**\n- 主机名称: 百度云-ops\n- 主机IP: 180.76.231.65\n- 检查时间: 2026-01-12 20:07:19\n- 压力等级: **正常**\n\n## CPU信息\n- CPU核心数: 2\n- CPU使用率: 0.0%\n\n## 内存信息\n- 总内存: 1886 MB\n- 已使用: 1456 MB\n- 使用率: 77.20%\n\n## 磁盘信息\n- 根分区使用率: 54%\n```\n/dev/vda2        40G   21G   18G  54% /\noverlay          40G   21G   18G  54% /var/lib/docker/overlay2/36d3c5e2aa5b2e721d81c27c2923c1645e788af8c62af73ca78c562acb1dc6c6/merged\noverlay          40G   21G   18G  54% /var/lib/docker/overlay2/3a19d7d7ec48824a825359d54df4f9c59bde749d91c29a4c2bf7c39d5bedb135/merged\noverlay          40G   21G   18G  54% /var/lib/docker/overlay2/273c4fcfacbd0de3a76036142e72f3838211b24fc5753e7a39494964479738ce/merged\n```\n\n## 系统负载\n- 负载平均值: 0.00, 0.01, 0.00\n\n## TOP 10 进程（按CPU使用率）\n| 用户 | PID | CPU% | MEM% | 命令 |\n|------|-----|------|------|------|\n| root | 94 | 11.1 | 0.0 | [kswapd0] |\n| root | 712663 | 4.0 | 0.5 | sshd: root@notty |\n| root | 750164 | 3.1 | 2.1 | /opt/heyeAgent/bin/heyeAgent |\n| root | 1 | 0.9 | 0.5 | /lib/systemd/systemd --switched-root --system --deserialize 31 nopti nospectre_v2 nospec_store_bypass_disable |\n| root | 712400 | 0.3 | 0.5 | sshd: unknown [priv] |\n| root | 4083585 | 0.3 | 0.5 | /opt/agent/dodevops-agent |\n| root | 852 | 0.2 | 1.0 | /usr/bin/containerd |\n| systemd+ | 723 | 0.1 | 0.1 | /lib/systemd/systemd-networkd |\n| fwupd-r+ | 920317 | 0.1 | 0.5 | /usr/bin/fwupdmgr refresh |\n| jenkins | 1997010 | 0.1 | 33.2 | /usr/bin/java -Djava.awt.headless=true -jar /usr/share/java/jenkins.war --webroot=/var/cache/jenkins/war --httpPort=8080 |\n\n## ? AI智能分析\n\n# 百度云-ops服务器性能分析与优化建议\n\n## 1. 性能瓶颈分析\n\n当前服务器最主要的性能问题集中在**内存使用**方面：\n- 内存使用率高达77.2%，接近警戒线(80%)\n- `kswapd0`进程(内核内存交换守护进程)占用CPU最高(11.1%)，表明存在内存压力\n- 有一个Java进程占用33.2%的内存，是最大的内存消耗者\n\n其他方面：\n- CPU负载极低(0.00-0.01)，2核CPU使用率为0%，说明CPU不是瓶颈\n- 磁盘使用率54%处于安全范围\n- 系统负载极低(0.00-0.01)，远低于CPU核心数(2)\n\n## 2. 根因分析\n\n内存高使用率的根本原因：\n1. **Java应用内存占用过高**：PID 1997010的Java进程占用33.2%内存(约626MB)，是最大的内存消耗者\n2. **监控/代理进程消耗资源**：`heyeAgent`(2.1%内存)和`dodevops-agent`(0.5%内存)等监控代理占用额外资源\n3. **内存交换活动**：`kswapd0`高CPU使用表明系统可能在进行内存交换，这会降低性能\n\n## 3. 优化建议 (按优先级排序)\n\n### 3.1 优化Java应用内存配置 (最高优先级)\n```bash\n# 1. 检查Java进程的启动参数\nps -ef | grep java\njcmd 1997010 VM.flags\n\n# 2. 建议调整JVM参数(根据实际应用需求调整)\n# 例如将-Xmx设置为适当值(当前可能设置过高)\n# 修改启动脚本中的JVM参数，如:\n# -Xms256m -Xmx512m -XX:MaxMetaspaceSize=256m\n```\n\n### 3.2 检查并优化监控代理 (次优先级)\n```bash\n# 1. 评估heyeAgent和dodevops-agent的必要性\nsystemctl status heyeAgent\nsystemctl status dodevops-agent\n\n# 2. 如果非必要，考虑停用或调整采集频率\n# 停用命令(谨慎操作):\n# systemctl stop heyeAgent\n# systemctl disable heyeAgent\n```\n\n### 3.3 增加系统交换空间 (临时缓解)\n```bash\n# 1. 检查当前交换空间\nfree -h\nswapon --show\n\n# 2. 创建额外交换文件(如果当前交换空间不足)\nsudo fallocate -l 1G /swapfile\nsudo chmod 600 /swapfile\nsudo mkswap /swapfile\nsudo swapon /swapfile\n\n# 3. 永久生效，添加到/etc/fstab\necho \'/swapfile none swap sw 0 0\' | sudo tee -a /etc/fstab\n```\n\n### 3.4 定期清理无用进程和连接\n```bash\n# 1. 检查异常的SSH连接(如notty连接)\nsudo netstat -tnpa | grep ESTABLISHED | grep sshd\n\n# 2. 清理闲置SSH会话\nsudo kill -9 712663  # 谨慎操作，确认会话可终止\n```\n\n### 3.5 长期监控和建议\n```bash\n# 安装并配置基础监控工具\nsudo apt install htop sysstat\n\n# 设置cron定期收集内存使用情况\n(crontab -l 2>/dev/null; echo \"0 * * * * free -m > /var/log/memory.log\") | crontab -\n```\n\n## 4. 风险提示\n\n1. **Java应用调整风险**：\n   - 直接修改Java内存参数可能导致应用崩溃\n   - 建议先在测试环境验证，或联系应用开发者确认合适的内存配置\n\n2. **服务中断风险**：\n   - 停止监控代理可能影响运维监控系统\n   - 建议先与运维团队确认这些代理的重要性\n\n3. **会话终止风险**：\n   - 强制终止SSH会话可能导致用户工作丢失\n   - 建议先确认这些会话是否确实闲置\n\n4. **内存不足风险**：\n   - 当前内存使用已接近警戒线，如果应用内存需求突增，可能导致OOM(Out Of Memory)错误\n   - 建议考虑升级内存(如增加到4GB)以获得更稳定的运行环境\n\n建议先实施3.1和3.3的优化措施，这些相对安全且能快速见效。其他措施应在业务低峰期实施，并做好回滚准备。\n\n','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',2,'','2026-01-12 20:07:59.545'),(129,'24f38ef5-aa45-41aa-aa71-e27c0b901051',89,'user','分析下 180.76.231.65 这个机器的压力 ip=180.76.231.65','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',1,'','2026-01-12 20:14:41.551'),(130,'24f38ef5-aa45-41aa-aa71-e27c0b901051',89,'assistant','# 服务器压力巡检报告\n\n**主机信息**\n- 主机名称: 百度云-ops\n- 主机IP: 180.76.231.65\n- 检查时间: 2026-01-12 20:14:42\n- 压力等级: **正常**\n\n## CPU信息\n- CPU核心数: 2\n- CPU使用率: 0.0%\n\n## 内存信息\n- 总内存: 1886 MB\n- 已使用: 1472 MB\n- 使用率: 78.05%\n\n## 磁盘信息\n- 根分区使用率: 54%\n```\n/dev/vda2        40G   21G   18G  54% /\noverlay          40G   21G   18G  54% /var/lib/docker/overlay2/36d3c5e2aa5b2e721d81c27c2923c1645e788af8c62af73ca78c562acb1dc6c6/merged\noverlay          40G   21G   18G  54% /var/lib/docker/overlay2/3a19d7d7ec48824a825359d54df4f9c59bde749d91c29a4c2bf7c39d5bedb135/merged\noverlay          40G   21G   18G  54% /var/lib/docker/overlay2/273c4fcfacbd0de3a76036142e72f3838211b24fc5753e7a39494964479738ce/merged\n```\n\n## 系统负载\n- 负载平均值: 0.23, 0.40, 0.22\n\n## TOP 10 进程（按CPU使用率）\n| 用户 | PID | CPU% | MEM% | 命令 |\n|------|-----|------|------|------|\n| root | 94 | 11.1 | 0.0 | [kswapd0] |\n| root | 750164 | 3.1 | 2.1 | /opt/heyeAgent/bin/heyeAgent |\n| root | 714736 | 2.5 | 0.5 | sshd: test [priv] |\n| root | 1 | 0.9 | 0.5 | /lib/systemd/systemd --switched-root --system --deserialize 31 nopti nospectre_v2 nospec_store_bypass_disable |\n| root | 714565 | 0.8 | 0.5 | sshd: test [priv] |\n| root | 4083585 | 0.3 | 0.5 | /opt/agent/dodevops-agent |\n| root | 852 | 0.2 | 1.0 | /usr/bin/containerd |\n| systemd+ | 723 | 0.1 | 0.1 | /lib/systemd/systemd-networkd |\n| fwupd-r+ | 920317 | 0.1 | 0.5 | /usr/bin/fwupdmgr refresh |\n| jenkins | 1997010 | 0.1 | 33.2 | /usr/bin/java -Djava.awt.headless=true -jar /usr/share/java/jenkins.war --webroot=/var/cache/jenkins/war --httpPort=8080 |\n\n## ? AI智能分析\n\n# 百度云-ops服务器性能分析与优化建议\n\n## 1. 性能瓶颈分析\n\n当前服务器最显著的性能问题是**内存使用率过高**（78.05%），接近警戒线（通常80%为警戒线）。具体分析如下：\n\n- **内存压力**：kswapd0进程（内核内存交换守护进程）占用11.1% CPU，表明系统正在频繁进行内存交换操作\n- **潜在内存泄漏**：Java进程占用33.2%内存（约625MB），是最大的单一内存用户\n- **监控代理负担**：heyeAgent和dodevops-agent两个监控代理共占用2.6% CPU和2.6%内存\n\n## 2. 根因分析\n\n### 内存问题根本原因\n1. **kswapd0活跃**：表明系统物理内存不足，开始使用交换空间，这会显著降低系统性能\n2. **Java应用内存配置不当**：单个Java进程占用过多内存（33.2%），可能未设置合理的JVM内存参数\n3. **监控代理资源消耗**：多个监控代理运行增加了系统负担\n\n### 其他潜在问题\n- **fwupdmgr进程**：正在运行固件更新检查，可能临时性占用资源\n- **容器运行时**：containerd进程运行，表明系统可能运行容器工作负载\n\n## 3. 优化建议（按优先级排序）\n\n### 1. Java应用内存优化（最高优先级）\n```bash\n# 检查当前Java进程的JVM参数\njcmd <PID> VM.flags\n\n# 建议为Java应用设置合理的堆内存限制，例如：\n# 在启动命令中添加（根据实际需求调整）：\n-Xms256m -Xmx512m -XX:MaxMetaspaceSize=128m\n```\n\n### 2. 监控代理整合与优化\n```bash\n# 检查监控代理的必要性和配置\n# 对于heyeAgent：\nps -ef | grep heyeAgent\n/opt/heyeAgent/bin/heyeAgent -v  # 查看版本和配置\n\n# 建议：\n# 1. 评估是否可以合并或减少监控代理数量\n# 2. 调整监控采集频率（如有配置选项）\n```\n\n### 3. 内存使用分析与优化\n```bash\n# 安装并运行smem进行内存使用分析\nsudo apt install smem\nsmem -t -k -u\n\n# 查看内存详细使用情况\ncat /proc/meminfo\nfree -h\n\n# 临时缓解（如需立即释放缓存）：\necho 3 > /proc/sys/vm/drop_caches\n```\n\n### 4. 交换空间优化\n```bash\n# 检查当前交换空间配置\nswapon --show\nfree -h\n\n# 如果尚未配置交换空间或太小，建议：\nsudo fallocate -l 1G /swapfile\nsudo chmod 600 /swapfile\nsudo mkswap /swapfile\nsudo swapon /swapfile\n\n# 永久生效，添加到/etc/fstab：\n/swapfile none swap sw 0 0\n```\n\n### 5. 定期维护任务调度\n```bash\n# 检查定时任务\nsystemctl list-timers\ncrontab -l\n\n# 建议将资源密集型任务（如fwupd）安排在低峰期\n```\n\n## 4. 风险提示\n\n1. **Java内存调整风险**：\n   - 过度减少Java堆内存可能导致应用OOM错误\n   - 建议先在测试环境验证新参数，并监控应用日志\n\n2. **监控代理调整风险**：\n   - 停用或合并监控代理可能影响监控系统数据收集\n   - 建议先与监控团队协调，确保关键指标仍被采集\n\n3. **交换空间调整风险**：\n   - 增加交换空间会占用磁盘空间\n   - 在SSD上频繁交换可能影响磁盘寿命\n\n4. **系统稳定性**：\n   - 当前内存压力已较高，任何增加内存使用的操作都可能触发OOM Killer\n   - 建议优先处理Java内存问题，再进行其他优化\n\n## 后续监控建议\n\n1. 设置内存使用警报（建议阈值：80%警告，90%严重）\n2. 持续监控kswapd0活动，理想情况下应很少看到它占用CPU\n3. 定期检查TOP进程列表，识别新的资源消耗大户\n\n','incident',0.70,'{\"checkType\": \"performance\", \"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',2,'','2026-01-12 20:15:25.196'),(131,'dcdc4483-6b75-46c8-b965-286b757b9185',89,'user','180.76.231.65\n139.9.205.38\n8.130.14.34\n检查这三台主机的 cpu是否正常 name=180.76.231.65','incident',0.70,'{\"checkType\": \"cpu\", \"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',1,'','2026-01-13 10:26:10.864'),(132,'dcdc4483-6b75-46c8-b965-286b757b9185',89,'assistant','# CPU 使用率分析\n\n**主机**: 百度云-ops (180.76.231.65)\n\n✅ CPU使用率正常\n','incident',0.70,'{\"checkType\": \"cpu\", \"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',2,'','2026-01-13 10:26:18.648'),(133,'dcdc4483-6b75-46c8-b965-286b757b9185',89,'user','180.76.231.65, 139.9.205.38, 8.130.14.34\n  检查这三台主机的 cpu是否正常 ip=180.76.231.65','incident',0.70,'{\"checkType\": \"cpu\", \"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',1,'','2026-01-13 11:10:00.140'),(134,'dcdc4483-6b75-46c8-b965-286b757b9185',89,'assistant','# CPU 使用率分析\n\n**主机**: 百度云-ops (180.76.231.65)\n\n✅ CPU使用率正常\n','incident',0.70,'{\"checkType\": \"cpu\", \"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',2,'','2026-01-13 11:10:07.936'),(135,'dcdc4483-6b75-46c8-b965-286b757b9185',89,'user','180.76.231.65, 139.9.205.38, 8.130.14.34 检查这三台主机的 cpu是否正常 ip=180.76.231.65','incident',0.70,'{\"checkType\": \"cpu\", \"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',1,'','2026-01-13 11:10:34.489'),(136,'dcdc4483-6b75-46c8-b965-286b757b9185',89,'assistant','# CPU 使用率分析\n\n**主机**: 百度云-ops (180.76.231.65)\n\n✅ CPU使用率正常\n','incident',0.70,'{\"checkType\": \"cpu\", \"targetName\": \"180.76.231.65\", \"targetType\": \"host\"}',0,'',2,'','2026-01-13 11:10:42.334'),(137,'dcdc4483-6b75-46c8-b965-286b757b9185',89,'user','180.76.231.65, 139.9.205.38, 8.130.14.34 检查这三台主机的 cpu是否正常 ip=180.76.231.65','incident',0.70,'{\"checkType\": \"cpu\", \"targetName\": \"180.76.231.65,139.9.205.38,8.130.14.34\", \"targetType\": \"host\"}',0,'',1,'','2026-01-13 11:14:45.166'),(138,'dcdc4483-6b75-46c8-b965-286b757b9185',89,'assistant','## 批量分析报告\n\n共分析 3 台主机，成功 3 台\n\n---\n\n### 百度云-ops (180.76.231.65)\n\n# CPU 使用率分析\n\n**主机**: 百度云-ops (180.76.231.65)\n\n✅ CPU使用率正常\n\n\n---\n\n### 华为云-ops (139.9.205.38)\n\n# CPU 使用率分析\n\n**主机**: 华为云-ops (139.9.205.38)\n\n✅ CPU使用率正常\n\n\n---\n\n### 阿里云-ops (8.130.14.34)\n\n# CPU 使用率分析\n\n**主机**: 阿里云-ops (8.130.14.34)\n\n✅ CPU使用率正常\n','incident',0.70,'{\"checkType\": \"cpu\", \"targetName\": \"180.76.231.65,139.9.205.38,8.130.14.34\", \"targetType\": \"host\"}',0,'',2,'','2026-01-13 11:15:05.203'),(139,'baa2ecd7-7693-4ae8-a6e4-20cd73cd9715',89,'user','8.130.14.34\n139.9.205.38\n180.76.231.65\n检查这三个机器的cpu name=8.130.14.34','incident',0.70,'{\"checkType\": \"cpu\", \"targetName\": \"8.130.14.34\", \"targetType\": \"host\"}',0,'',1,'','2026-01-15 10:52:10.869'),(140,'baa2ecd7-7693-4ae8-a6e4-20cd73cd9715',89,'assistant','# CPU 使用率分析\n\n**主机**: 阿里云-ops (8.130.14.34)\n\n✅ CPU使用率正常\n','incident',0.70,'{\"checkType\": \"cpu\", \"targetName\": \"8.130.14.34\", \"targetType\": \"host\"}',0,'',2,'','2026-01-15 10:52:19.884'),(141,'baa2ecd7-7693-4ae8-a6e4-20cd73cd9715',89,'user','8.130.14.34  139.9.205.38  180.76.231.65 \n检查这三个机器的cpu ip=8.130.14.34','incident',0.70,'{\"checkType\": \"cpu\", \"targetName\": \"8.130.14.34,139.9.205.38,180.76.231.65\", \"targetType\": \"host\"}',0,'',1,'','2026-01-15 10:54:37.242'),(142,'baa2ecd7-7693-4ae8-a6e4-20cd73cd9715',89,'assistant','## 批量分析报告\n\n共分析 3 台主机，成功 3 台\n\n---\n\n### 阿里云-ops (8.130.14.34)\n\n# CPU 使用率分析\n\n**主机**: 阿里云-ops (8.130.14.34)\n\n✅ CPU使用率正常\n\n\n---\n\n### 华为云-ops (139.9.205.38)\n\n# CPU 使用率分析\n\n**主机**: 华为云-ops (139.9.205.38)\n\n✅ CPU使用率正常\n\n\n---\n\n### 百度云-ops (180.76.231.65)\n\n# CPU 使用率分析\n\n**主机**: 百度云-ops (180.76.231.65)\n\n✅ CPU使用率正常\n','incident',0.70,'{\"checkType\": \"cpu\", \"targetName\": \"8.130.14.34,139.9.205.38,180.76.231.65\", \"targetType\": \"host\"}',0,'',2,'','2026-01-15 10:54:57.193');
/*!40000 ALTER TABLE `ai_agent_chat_history` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ai_agent_task`
--

DROP TABLE IF EXISTS `ai_agent_task`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ai_agent_task` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键''',
  `type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''任务类型:inspection/incident/prediction/custom''',
  `name` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''任务名称''',
  `description` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''任务描述''',
  `priority` bigint DEFAULT '5' COMMENT '''优先级:1-10''',
  `status` bigint NOT NULL COMMENT '''状态:1=待执行,2=执行中,3=成功,4=失败,5=取消''',
  `params` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''任务参数(JSON)''',
  `result` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''执行结果(JSON)''',
  `error_msg` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''错误信息''',
  `retry` bigint DEFAULT '0' COMMENT '''重试次数''',
  `max_retry` bigint DEFAULT '3' COMMENT '''最大重试次数''',
  `start_time` datetime(3) DEFAULT NULL COMMENT '''开始时间''',
  `end_time` datetime(3) DEFAULT NULL COMMENT '''结束时间''',
  `create_time` datetime(3) NOT NULL COMMENT '''创建时间''',
  `update_time` datetime(3) NOT NULL COMMENT '''更新时间''',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_ai_agent_task_type` (`type`) USING BTREE,
  KEY `idx_ai_agent_task_status` (`status`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ai_agent_task`
--

LOCK TABLES `ai_agent_task` WRITE;
/*!40000 ALTER TABLE `ai_agent_task` DISABLE KEYS */;
/*!40000 ALTER TABLE `ai_agent_task` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ai_model`
--

DROP TABLE IF EXISTS `ai_model`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `ai_model` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键''',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''模型名称''',
  `type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''模型类型(阿里百炼,OpenAI等)''',
  `url` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''API地址''',
  `api_key` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''API密钥''',
  `model` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''模型标识(qwen-plus,gpt-4等)''',
  `status` bigint NOT NULL DEFAULT '1' COMMENT '''状态:1->启用,2->禁用''',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''备注''',
  `create_time` datetime(3) NOT NULL COMMENT '''创建时间''',
  `update_time` datetime(3) NOT NULL COMMENT '''更新时间''',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ai_model`
--

LOCK TABLES `ai_model` WRITE;
/*!40000 ALTER TABLE `ai_model` DISABLE KEYS */;
/*!40000 ALTER TABLE `ai_model` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `app_application`
--

DROP TABLE IF EXISTS `app_application`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `app_application` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `code` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `business_group_id` bigint unsigned NOT NULL,
  `business_dept_id` bigint unsigned NOT NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
  `repo_url` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `dev_owners` json DEFAULT NULL,
  `test_owners` json DEFAULT NULL,
  `ops_owners` json DEFAULT NULL,
  `programming_lang` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `start_command` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
  `stop_command` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
  `health_api` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `domains` json DEFAULT NULL,
  `hosts` json DEFAULT NULL,
  `databases` json DEFAULT NULL,
  `other_res` json DEFAULT NULL,
  `status` tinyint DEFAULT '1',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `idx_app_application_code` (`code`) USING BTREE,
  KEY `idx_app_application_deleted_at` (`deleted_at`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=22 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `app_application`
--

LOCK TABLES `app_application` WRITE;
/*!40000 ALTER TABLE `app_application` DISABLE KEYS */;
INSERT INTO `app_application` VALUES (20,'s3-api','s3-api',60,2,'s3-api','https://github.com/zhang1024fan/deviops','[107]','[106]','[89]','Java','','','','[\"www.deviops.cn\"]','[1]','[2]','{}',2,'2026-01-19 15:55:00.906','2026-01-19 15:55:47.509',NULL),(21,'s3-web','s3-web',60,2,'s3-web','https://github.com/zhang1024fan/deviops','[106]','[106]','[89]','Go','','','',NULL,NULL,NULL,'{}',2,'2026-01-19 16:36:50.376','2026-01-19 17:00:54.643',NULL);
/*!40000 ALTER TABLE `app_application` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `app_jenkins_env`
--

DROP TABLE IF EXISTS `app_jenkins_env`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `app_jenkins_env` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `app_id` bigint unsigned NOT NULL,
  `env_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `jenkins_server_id` bigint unsigned DEFAULT NULL,
  `job_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '',
  `job_url` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `build_params` json DEFAULT NULL,
  `deploy_config` json DEFAULT NULL,
  `notification` json DEFAULT NULL,
  `is_active` tinyint DEFAULT '1',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_app_jenkins_env_app_id` (`app_id`) USING BTREE,
  KEY `idx_app_jenkins_env_deleted_at` (`deleted_at`) USING BTREE,
  CONSTRAINT `fk_app_application_jenkins_envs` FOREIGN KEY (`app_id`) REFERENCES `app_application` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=60 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `app_jenkins_env`
--

LOCK TABLES `app_jenkins_env` WRITE;
/*!40000 ALTER TABLE `app_jenkins_env` DISABLE KEYS */;
INSERT INTO `app_jenkins_env` VALUES (54,20,'prod',NULL,'',NULL,NULL,NULL,NULL,1,'2026-01-19 15:55:01.154','2026-01-19 15:55:01.154',NULL),(55,20,'test',16,'s3-api',NULL,NULL,NULL,NULL,1,'2026-01-19 15:55:01.389','2026-01-19 15:55:47.192',NULL),(56,20,'dev',NULL,'',NULL,NULL,NULL,NULL,1,'2026-01-19 15:55:01.620','2026-01-19 15:55:01.620',NULL),(57,21,'prod',NULL,'',NULL,NULL,NULL,NULL,1,'2026-01-19 16:36:50.614','2026-01-19 16:36:50.614',NULL),(58,21,'test',16,'s3-web',NULL,NULL,NULL,NULL,1,'2026-01-19 16:36:50.828','2026-01-19 17:00:54.205',NULL),(59,21,'dev',NULL,'',NULL,NULL,NULL,NULL,1,'2026-01-19 16:36:51.039','2026-01-19 16:36:51.039',NULL);
/*!40000 ALTER TABLE `app_jenkins_env` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `app_service_release`
--

DROP TABLE IF EXISTS `app_service_release`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `app_service_release` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '发布标题',
  `business_group_id` bigint unsigned NOT NULL COMMENT '业务组ID',
  `impact_feature` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '影响功能描述',
  `applicant_id` bigint unsigned NOT NULL COMMENT '申请人ID',
  `applicant_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '申请人姓名',
  `owner_approver_id` bigint unsigned DEFAULT NULL COMMENT '负责人审批人ID',
  `owner_approver_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '负责人审批人姓名',
  `security_approver_id` bigint unsigned DEFAULT NULL COMMENT '安全审批人ID',
  `security_approver_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '安全审批人姓名',
  `test_approver_id` bigint unsigned DEFAULT NULL COMMENT '测试审批人ID',
  `test_approver_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '测试审批人姓名',
  `owner_approval_status` bigint DEFAULT '1' COMMENT '负责人审批状态',
  `security_approval_status` bigint DEFAULT '1' COMMENT '安全审批状态',
  `test_approval_status` bigint DEFAULT '1' COMMENT '测试审批状态',
  `owner_approval_time` datetime(3) DEFAULT NULL COMMENT '负责人审批时间',
  `security_approval_time` datetime(3) DEFAULT NULL COMMENT '安全审批时间',
  `test_approval_time` datetime(3) DEFAULT NULL COMMENT '测试审批时间',
  `owner_approval_remark` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '负责人审批意见',
  `security_approval_remark` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '安全审批意见',
  `test_approval_remark` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '测试审批意见',
  `deploy_status` bigint DEFAULT '1' COMMENT '运维发布状态',
  `regression_test_status` bigint DEFAULT '1' COMMENT '回归测试状态',
  `status` bigint DEFAULT '1' COMMENT '流程状态',
  `start_time` datetime(3) DEFAULT NULL COMMENT '发布开始时间',
  `end_time` datetime(3) DEFAULT NULL COMMENT '发布结束时间',
  `duration` bigint DEFAULT '0' COMMENT '发布耗时(秒)',
  `service_count` bigint DEFAULT '0' COMMENT '关联服务数量',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_app_service_release_deleted_at` (`deleted_at`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `app_service_release`
--

LOCK TABLES `app_service_release` WRITE;
/*!40000 ALTER TABLE `app_service_release` DISABLE KEYS */;
/*!40000 ALTER TABLE `app_service_release` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `app_service_release_item`
--

DROP TABLE IF EXISTS `app_service_release_item`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `app_service_release_item` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `release_id` bigint unsigned NOT NULL COMMENT '上线工单ID',
  `app_id` bigint unsigned NOT NULL COMMENT '应用ID',
  `app_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '应用名称',
  `app_code` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '应用编码',
  `project_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '项目服务名称',
  `repo_url` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '项目地址',
  `branch` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT 'master' COMMENT '发布分支',
  `commit_id` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT 'Commit ID',
  `impact_feature` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '影响功能',
  `function_module` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '功能模块',
  `db_change` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '数据库变更',
  `config_change` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '配置变更',
  `remark` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '备注信息',
  `jenkins_env_id` bigint unsigned DEFAULT NULL COMMENT 'Jenkins环境配置ID',
  `jenkins_job_url` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT 'Jenkins任务URL',
  `parameters` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT 'Jenkins构建参数(JSON格式)',
  `build_number` bigint DEFAULT '0' COMMENT '构建编号',
  `log_url` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '构建日志URL',
  `status` bigint DEFAULT '1' COMMENT '发布状态',
  `start_time` datetime(3) DEFAULT NULL COMMENT '发布开始时间',
  `end_time` datetime(3) DEFAULT NULL COMMENT '发布结束时间',
  `duration` bigint DEFAULT '0' COMMENT '发布耗时(秒)',
  `error_message` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '错误信息',
  `execute_order` bigint DEFAULT '0' COMMENT '执行顺序',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_app_service_release_item_release_id` (`release_id`) USING BTREE,
  KEY `idx_app_service_release_item_deleted_at` (`deleted_at`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `app_service_release_item`
--

LOCK TABLES `app_service_release_item` WRITE;
/*!40000 ALTER TABLE `app_service_release_item` DISABLE KEYS */;
INSERT INTO `app_service_release_item` VALUES (1,2,14,'dev-lku-sass-pack','dev-lku-sass-pack','dev-lku-sass-pack','git@code.dding.net/lockin/saas-toc-server.git','master','083a2d6669effe63670b0d5ee899d29eeb187abf','影响用户登录','','没有','没有','没有',35,'','{\"commit_id\":\"083a2d6669effe63670b0d5ee899d29eeb187abf\",\"compile\":\"true\"}',408,'http://test-ops-jenkins-tc2.dding.net:8080/job/dev-lku-sass-pack/408/console',3,'2025-11-22 18:18:03.823','2025-11-22 18:21:44.766',220,'',1,'2025-11-22 18:13:37.484','2025-11-22 18:21:44.821',NULL),(2,3,15,'dev-lku-sass-deploy','dev-lku-sass-deploy','dev-lku-sass-deploy','git@code.dding.net/lockin/saas-toc-server.git','master','c7417cbe54a11d89c819b4541316980dc1634687','测试多任务工单发布','','无','无','无',39,'','',19,'http://test-ops-jenkins-tc2.dding.net:8080/job/test1-saas-deploy/19/console',4,'2025-11-22 19:36:10.184','2025-11-22 19:36:37.458',27,'Jenkins构建失败',1,'2025-11-22 19:31:58.231','2025-11-22 19:36:37.524',NULL),(3,3,14,'dev-lku-sass-pack','dev-lku-sass-pack','dev-lku-sass-pack','git@code.dding.net/lockin/saas-toc-server.git','master','c7417cbe54a11d89c819b4541316980dc1634687','测试多任务工单发布','','无','无','无',35,'','',409,'http://test-ops-jenkins-tc2.dding.net:8080/job/dev-lku-sass-pack/409/console',3,'2025-11-22 19:36:10.184','2025-11-22 19:39:59.417',229,'',2,'2025-11-22 19:31:58.710','2025-11-22 19:39:59.476',NULL),(4,3,18,'test1-paas-deploy','test1-paas-deploy','test1-paas-deploy','git@code.dding.net/lockin/cloud-platform.git','master','c7417cbe54a11d89c819b4541316980dc1634687','测试多任务工单发布','','无','无','无',48,'','',13,'http://test-ops-jenkins-tc2.dding.net:8080/job/test1-paas-deploy/13/console',4,'2025-11-22 19:36:10.184','2025-11-22 19:36:38.200',28,'Jenkins构建失败',3,'2025-11-22 19:31:59.196','2025-11-22 19:36:38.282',NULL),(5,3,17,'test1-paas-pack','test1-paas-pack','test1-paas-pack','git@code.dding.net/lockin/cloud-platform.git','master','c7417cbe54a11d89c819b4541316980dc1634687','测试多任务工单发布','','无','无','无',45,'','',14,'http://test-ops-jenkins-tc2.dding.net:8080/job/test1-paas-pack/14/console',3,'2025-11-22 19:36:10.184','2025-11-22 19:41:18.055',307,'',4,'2025-11-22 19:31:59.673','2025-11-22 19:41:18.106',NULL),(6,4,17,'test1-paas-pack','test1-paas-pack','test1-paas-pack','git@code.dding.net/lockin/cloud-platform.git','master','c7417cbe54a11d89c819b4541316980dc1634687','测试审批流程','','测试审批流程','测试审批流程','测试审批流程',45,'','',15,'http://test-ops-jenkins-tc2.dding.net:8080/job/test1-paas-pack/15/console',3,'2025-11-22 20:29:17.308','2025-11-22 20:34:08.653',291,'',1,'2025-11-22 20:00:19.453','2025-11-22 20:34:08.746',NULL),(7,4,18,'test1-paas-deploy','test1-paas-deploy','test1-paas-deploy','git@code.dding.net/lockin/cloud-platform.git','master','c7417cbe54a11d89c819b4541316980dc1634687','测试审批流程','','测试审批流程','测试审批流程','测试审批流程',48,'','',14,'http://test-ops-jenkins-tc2.dding.net:8080/job/test1-paas-deploy/14/console',4,'2025-11-22 20:34:09.058','2025-11-22 20:34:21.625',12,'Jenkins构建失败',2,'2025-11-22 20:00:19.908','2025-11-22 20:34:21.687',NULL),(8,5,14,'dev-lku-sass-pack','dev-lku-sass-pack','dev-lku-sass-pack','git@code.dding.net/lockin/saas-toc-server.git','master','083a2d6669effe63670b0d5ee899d29eeb187abf','无','','无','无','无',35,'','',410,'http://test-ops-jenkins-tc2.dding.net:8080/job/dev-lku-sass-pack/410/console',3,'2025-11-22 23:52:18.208','2025-11-22 23:55:51.178',212,'',1,'2025-11-22 23:13:40.348','2025-11-22 23:55:51.229',NULL),(9,6,17,'test1-paas-pack','test1-paas-pack','test1-paas-pack','git@code.dding.net/lockin/cloud-platform.git','master','c7417cbe54a11d89c819b4541316980dc1634687','123','','123','123','123',45,'','{\"commit_id\":\"c7417cbe54a11d89c819b4541316980dc1634687\",\"compile\":\"true\"}',0,'',1,NULL,NULL,0,'',1,'2025-11-23 11:32:20.185','2025-11-23 11:32:20.185',NULL),(10,6,14,'dev-lku-sass-pack','dev-lku-sass-pack','dev-lku-sass-pack','git@code.dding.net/lockin/saas-toc-server.git','master','083a2d6669effe63670b0d5ee899d29eeb187abf','123','','123','123','123',35,'','{\"commit_id\":\"083a2d6669effe63670b0d5ee899d29eeb187abf\",\"compile\":\"true\"}',0,'',1,NULL,NULL,0,'',2,'2025-11-23 11:32:20.664','2025-11-23 11:32:20.664',NULL),(11,7,14,'dev-lku-sass-pack','dev-lku-sass-pack','dev-lku-sass-pack','git@code.dding.net/lockin/saas-toc-server.git','master','083a2d6669effe63670b0d5ee899d29eeb187abf','影响用户登录','','无','无','123',35,'','{\"commit_id\":\"083a2d6669effe63670b0d5ee899d29eeb187abf\",\"compile\":\"true\"}',411,'http://test-ops-jenkins-tc2.dding.net:8080/job/dev-lku-sass-pack/411/console',3,'2025-11-24 00:46:53.314','2025-11-24 00:52:37.448',344,'',1,'2025-11-24 00:42:25.108','2025-11-24 00:52:37.512',NULL),(12,8,19,'prod_saas3.0_data-export','prod-saas30-data-export','prod_saas3.0_data-export','git@gitee.com:zhang_fan1024/zf-k8s.git','master','','影响用户登录','','无','无','无',51,'','{\"commit_id\":\"123456789\"}',10,'http://test-ops-jenkins-tc2.dding.net:8080/job/prod_saas3.0_data-export/10/console',3,'2025-12-05 20:47:40.229','2025-12-05 20:47:52.992',12,'',1,'2025-12-05 20:46:53.654','2025-12-05 20:47:53.043',NULL);
/*!40000 ALTER TABLE `app_service_release_item` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `app_sh_release`
--

DROP TABLE IF EXISTS `app_sh_release`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `app_sh_release` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '上线标题',
  `reason` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '上线原因',
  `business_group_id` bigint unsigned NOT NULL COMMENT '业务线ID',
  `app_id` bigint unsigned NOT NULL COMMENT '服务ID',
  `app_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '服务名称',
  `app_code` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '服务编码',
  `applicant_id` bigint unsigned NOT NULL COMMENT '申请人ID',
  `applicant_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '申请人姓名',
  `approver_id` bigint unsigned DEFAULT NULL COMMENT '审核人ID',
  `approver_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '审核人姓名',
  `executor_id` bigint unsigned DEFAULT NULL COMMENT '执行人ID',
  `executor_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '执行人姓名',
  `execute_dir` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '执行目录',
  `script_content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '脚本内容',
  `approval_status` bigint DEFAULT '1' COMMENT '审核状态',
  `approval_time` datetime(3) DEFAULT NULL COMMENT '审核时间',
  `approval_remark` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '审核意见',
  `execute_status` bigint DEFAULT '1' COMMENT '执行状态',
  `status` bigint DEFAULT '1' COMMENT '流程状态',
  `start_time` datetime(3) DEFAULT NULL COMMENT '脚本执行开始时间',
  `end_time` datetime(3) DEFAULT NULL COMMENT '脚本执行结束时间',
  `duration` bigint DEFAULT '0' COMMENT '执行耗时(秒)',
  `jenkins_env_id` bigint unsigned DEFAULT NULL COMMENT 'Jenkins环境配置ID',
  `build_number` bigint DEFAULT '0' COMMENT '构建编号',
  `log_url` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '构建日志URL',
  `error_message` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '错误信息',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `parameters` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT 'Jenkins构建参数(JSON格式)',
  `server_host_id` bigint unsigned NOT NULL COMMENT '执行服务器主机ID(关联cmdb_host)',
  `pull_code_start_time` datetime(3) DEFAULT NULL COMMENT '拉取代码开始时间',
  `pull_code_end_time` datetime(3) DEFAULT NULL COMMENT '拉取代码结束时间',
  `script_output` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '脚本执行输出',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_app_sh_release_deleted_at` (`deleted_at`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `app_sh_release`
--

LOCK TABLES `app_sh_release` WRITE;
/*!40000 ALTER TABLE `app_sh_release` DISABLE KEYS */;
INSERT INTO `app_sh_release` VALUES (1,'prod_saas3.0_data-export-数据导出','prod_saas3.0_data-export-数据导出',18,19,'prod_saas3.0_data-export','prod-saas30-data-export',1,'管理员',89,'admin',89,'admin','/home/dingding/saas3-data-export/','ls',2,'2025-11-28 12:42:48.493','111',3,2,'2025-11-28 12:43:08.943','2025-11-28 12:43:24.437',14,51,3,'http://test-ops-jenkins-tc2.dding.net:8080/job/prod_saas3.0_data-export/3/console','','2025-11-28 12:41:29.167','2025-11-28 12:43:24.488',NULL,'{\"commit_id\":\"123456789\"}',0,NULL,NULL,NULL),(2,'测试001','测试001',18,19,'prod_saas3.0_data-export','prod-saas30-data-export',1,'管理员',98,'李四',89,'admin','/home/','ls  /root/',2,'2025-11-28 14:18:58.866','同意',3,1,NULL,NULL,0,51,4,'http://test-ops-jenkins-tc2.dding.net:8080/job/prod_saas3.0_data-export/4/console','','2025-11-28 13:37:27.963','2025-11-28 14:47:38.777',NULL,'{\"commit_id\":\"123456789\"}',0,'2025-11-28 14:47:22.389','2025-11-28 14:47:38.725',NULL),(3,'测试脚本执行','测试脚本执行',18,19,'prod_saas3.0_data-export','prod-saas30-data-export',1,'管理员',89,'admin',89,'admin','/home/dingding/saas3-data-export/','pwd\nls\nhostname -I',2,'2025-11-28 15:19:15.004','111',6,2,'2025-11-28 15:24:27.057','2025-11-28 15:24:31.379',3,51,5,'http://test-ops-jenkins-tc2.dding.net:8080/job/prod_saas3.0_data-export/5/console','','2025-11-28 15:09:18.474','2025-11-28 15:24:31.432',NULL,'{\"commit_id\":\"123456789\"}',501,'2025-11-28 15:19:21.580','2025-11-28 15:19:34.271','bash: line 1: cd: /home/dingding/saas3-data-export/: No such file or directory\ndocker\nelk\njdk11\njdk17\njdk18\nluban-master\nnode\nprometheus\nsnap\n172.20.236.121 172.18.0.1 172.17.0.1 172.19.0.1 \n'),(4,'测试脚本002','测试脚本002',18,19,'prod_saas3.0_data-export','prod-saas30-data-export',1,'管理员',89,'admin',89,'admin','/home/dingding/saas3-data-export/','hostname\npwd\ndate',2,'2025-11-28 15:56:35.290','11',6,2,'2025-11-28 15:57:23.260','2025-11-28 15:57:27.569',3,51,6,'http://test-ops-jenkins-tc2.dding.net:8080/job/prod_saas3.0_data-export/6/console','','2025-11-28 15:56:13.417','2025-11-28 15:57:27.644',NULL,'{\"commit_id\":\"123456789\"}',501,'2025-11-28 15:56:44.717','2025-11-28 15:57:00.303','/root\nFri Nov 28 03:57:27 PM CST 2025\nbash: line 1: cd: /home/dingding/saas3-data-export/: No such file or directory\n'),(5,'测试002','测试002',18,19,'prod_saas3.0_data-export','prod-saas30-data-export',1,'管理员',89,'admin',89,'admin','/home/','pwd\nls\nhostname\n',2,'2025-11-28 16:12:04.059','ok',6,2,'2025-11-28 16:19:18.636','2025-11-28 16:19:22.516',3,51,7,'http://test-ops-jenkins-tc2.dding.net:8080/job/prod_saas3.0_data-export/7/console','','2025-11-28 16:10:22.217','2025-11-28 16:19:22.570',NULL,'{\"commit_id\":\"123456789\"}',501,'2025-11-28 16:12:11.076','2025-11-28 16:12:23.702','/home\ndevops\ndevops.tar.gz\ngo-ops\n'),(6,'test1111111111','test',18,19,'prod_saas3.0_data-export','prod-saas30-data-export',1,'管理员',89,'admin',89,'admin','/home/dingding/saas3-data-export/','pwd  \nls \nhostname',2,'2025-12-01 01:01:48.329','ok',6,2,'2025-12-01 01:02:36.996','2025-12-01 01:02:39.204',1,51,8,'http://test-ops-jenkins-tc2.dding.net:8080/job/prod_saas3.0_data-export/8/console','','2025-12-01 01:01:08.782','2025-12-01 01:02:39.260',NULL,'{\"commit_id\":\"123456789\"}',501,'2025-12-01 01:01:59.779','2025-12-01 01:02:16.154','/home/dingding/saas3-data-export\nxlsx\ngo-ops\n'),(7,'测试002','测试002',12,19,'prod_saas3.0_data-export','prod-saas30-data-export',1,'管理员',89,'管理员',89,'管理员','/home/dingding/saas3-data-export/','hostname -I',2,'2025-12-25 11:36:50.014','ok',6,2,'2025-12-25 11:39:15.664','2025-12-25 11:39:17.092',0,51,11,'http://test-ops-jenkins-tc2.dding.net:8080/job/prod_saas3.0_data-export/11/console','','2025-12-25 11:36:31.470','2025-12-25 11:39:17.147',NULL,'{\"commit_id\":\"123456789\"}',506,'2025-12-25 11:37:37.238','2025-12-25 11:37:49.711','172.31.6.35 172.17.0.1 172.20.0.1 \n');
/*!40000 ALTER TABLE `app_sh_release` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `bastion_host_authorization`
--

DROP TABLE IF EXISTS `bastion_host_authorization`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `bastion_host_authorization` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键''',
  `subject_type` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''授权对象类型:user/group''',
  `subject_id` bigint unsigned NOT NULL COMMENT '''授权对象ID(用户ID或用户组ID)''',
  `host_id` bigint unsigned NOT NULL COMMENT '''主机ID''',
  `create_time` datetime(3) NOT NULL COMMENT '''创建时间''',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_subject_host` (`subject_type`,`subject_id`,`host_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `bastion_host_authorization`
--

LOCK TABLES `bastion_host_authorization` WRITE;
/*!40000 ALTER TABLE `bastion_host_authorization` DISABLE KEYS */;
INSERT INTO `bastion_host_authorization` VALUES (5,'group',3,2,'2026-01-18 20:32:10.636');
/*!40000 ALTER TABLE `bastion_host_authorization` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `bastion_user_group`
--

DROP TABLE IF EXISTS `bastion_user_group`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `bastion_user_group` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键''',
  `name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''用户组名称''',
  `description` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''用户组描述''',
  `member_ids` json DEFAULT NULL COMMENT '''成员ID列表(JSON数组)''',
  `status` bigint NOT NULL DEFAULT '1' COMMENT '''状态:1->启用,2->禁用''',
  `create_time` datetime(3) NOT NULL COMMENT '''创建时间''',
  `update_time` datetime(3) DEFAULT NULL COMMENT '''更新时间''',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `idx_bastion_user_group_name` (`name`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `bastion_user_group`
--

LOCK TABLES `bastion_user_group` WRITE;
/*!40000 ALTER TABLE `bastion_user_group` DISABLE KEYS */;
INSERT INTO `bastion_user_group` VALUES (3,'运维组','ops','[105, 106]',1,'2026-01-18 20:30:59.559','2026-01-19 15:39:55.785');
/*!40000 ALTER TABLE `bastion_user_group` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `cmdb_group`
--

DROP TABLE IF EXISTS `cmdb_group`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `cmdb_group` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键''',
  `parent_id` bigint unsigned NOT NULL DEFAULT '0' COMMENT '''父级分组ID''',
  `name` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''分组名称''',
  `create_time` datetime(3) NOT NULL COMMENT '''创建时间''',
  `remark` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''备注''',
  `update_time` datetime(3) DEFAULT NULL COMMENT '''更新时间''',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=66 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cmdb_group`
--

LOCK TABLES `cmdb_group` WRITE;
/*!40000 ALTER TABLE `cmdb_group` DISABLE KEYS */;
INSERT INTO `cmdb_group` VALUES (1,0,'默认业务组','2025-07-10 11:02:07.226',NULL,NULL),(40,0,'支撑业务','2025-12-28 01:14:36.499',NULL,NULL),(41,0,'信息中心','2025-12-28 01:14:51.443',NULL,NULL),(44,43,'web','2025-12-28 01:15:42.077',NULL,NULL),(45,43,'wyze','2025-12-28 01:15:52.905',NULL,NULL),(46,43,'embed','2025-12-28 01:16:01.587',NULL,NULL),(54,41,'sun','2025-12-28 01:17:52.361',NULL,NULL),(55,40,'ops','2025-12-28 01:18:07.873',NULL,NULL),(56,40,'pingtai','2025-12-28 01:18:25.301',NULL,NULL),(57,39,'ddm','2025-12-28 01:18:42.607',NULL,NULL),(58,39,'toc','2025-12-28 01:18:49.421',NULL,NULL),(65,64,'test-ops','2025-12-28 15:53:44.096',NULL,NULL);
/*!40000 ALTER TABLE `cmdb_group` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `cmdb_host`
--

DROP TABLE IF EXISTS `cmdb_host`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `cmdb_host` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键''',
  `host_name` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''名称''',
  `group_id` bigint unsigned NOT NULL COMMENT '''分组ID''',
  `private_ip` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''私网IP''',
  `public_ip` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''公网IP''',
  `ssh_name` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''SSH用户名''',
  `ssh_key_id` bigint unsigned DEFAULT NULL COMMENT '''SSH凭据ID''',
  `ssh_port` bigint DEFAULT '22' COMMENT '''SSH端口''',
  `remark` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''备注''',
  `vendor` bigint DEFAULT NULL COMMENT '''1->自建,2->阿里云,3->腾讯云''',
  `region` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''区域''',
  `instance_id` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''实例ID''',
  `os` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''操作系统''',
  `status` bigint DEFAULT NULL COMMENT '''状态:1->认证成功,2->未认证,3->认证失败''',
  `cpu` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''CPU信息''',
  `memory` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''内存信息''',
  `disk` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''磁盘信息''',
  `billing_type` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''计费方式''',
  `create_time` datetime(3) NOT NULL COMMENT '''创建时间''',
  `expire_time` datetime(3) DEFAULT NULL COMMENT '''到期时间''',
  `update_time` datetime(3) DEFAULT NULL COMMENT '''更新时间''',
  `ssh_ip` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''SSH连接IP''',
  `name` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''ecs主机名称''',
  `ssh_gateway_id` bigint unsigned DEFAULT NULL COMMENT '''中转网关凭据ID''',
  `tag` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''标签(格式:key=value,key=value)''',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `fk_cmdb_group_hosts` (`group_id`) USING BTREE,
  CONSTRAINT `fk_cmdb_group_hosts` FOREIGN KEY (`group_id`) REFERENCES `cmdb_group` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=718 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cmdb_host`
--

LOCK TABLES `cmdb_host` WRITE;
/*!40000 ALTER TABLE `cmdb_host` DISABLE KEYS */;
/*!40000 ALTER TABLE `cmdb_host` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `cmdb_snmp_devices`
--

DROP TABLE IF EXISTS `cmdb_snmp_devices`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `cmdb_snmp_devices` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键''',
  `device_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''设备名称''',
  `device_type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''设备类型(switch/router/windows/linux)''',
  `ip_address` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''设备IP地址''',
  `snmp_version` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT 'v2c' COMMENT '''SNMP版本''',
  `snmp_community` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT 'public' COMMENT '''Community字符串''',
  `snmp_port` bigint DEFAULT '161' COMMENT '''SNMP端口''',
  `remote_port` bigint DEFAULT NULL COMMENT '''远程端口(SSH默认22,RDP默认3389)''',
  `remote_username` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''远程登录用户名''',
  `remote_password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''远程登录密码''',
  `remote_domain` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''域名(仅Windows RDP需要)''',
  `monitor_enabled` bigint DEFAULT '1' COMMENT '''是否启用监控(0=否,1=是)''',
  `location` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''物理位置''',
  `owner` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''负责人''',
  `remark` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''备注''',
  `status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT 'offline' COMMENT '''设备状态(online/offline)''',
  `last_seen` datetime(3) DEFAULT NULL COMMENT '''最后在线时间''',
  `create_time` datetime(3) NOT NULL COMMENT '''创建时间''',
  `update_time` datetime(3) NOT NULL COMMENT '''更新时间''',
  `ssh_port` bigint DEFAULT '22' COMMENT '''SSH端口''',
  `ssh_username` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''SSH用户名''',
  `ssh_password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''SSH密码''',
  `rdp_port` bigint DEFAULT '3389' COMMENT '''RDP端口''',
  `rdp_username` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''RDP用户名''',
  `rdp_password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''RDP密码''',
  `rdp_domain` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''RDP域''',
  `exporter_installed` bigint DEFAULT '0' COMMENT '''Exporter状态(0=未安装,1=已安装运行中,2=安装中,3=已安装未运行,4=安装失败,5=启动中)''',
  `exporter_port` bigint DEFAULT '9182' COMMENT '''Exporter监听端口''',
  `exporter_version` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''Exporter版本号''',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `idx_cmdb_snmp_devices_ip_address` (`ip_address`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cmdb_snmp_devices`
--

LOCK TABLES `cmdb_snmp_devices` WRITE;
/*!40000 ALTER TABLE `cmdb_snmp_devices` DISABLE KEYS */;
/*!40000 ALTER TABLE `cmdb_snmp_devices` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `cmdb_sql`
--

DROP TABLE IF EXISTS `cmdb_sql`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `cmdb_sql` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `type` int NOT NULL,
  `account_id` bigint unsigned NOT NULL,
  `group_id` bigint unsigned NOT NULL,
  `tags` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `description` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cmdb_sql`
--

LOCK TABLES `cmdb_sql` WRITE;
/*!40000 ALTER TABLE `cmdb_sql` DISABLE KEYS */;
INSERT INTO `cmdb_sql` VALUES (1,'saas3-mysql',1,1,1,'prod','1111','2025-07-29 21:23:17.309','2025-07-29 21:23:17.309'),(2,'saas3-redis-1',3,3,1,'prod','1111','2025-07-29 21:24:57.985','2025-09-06 15:12:22.605'),(4,'saas3-pgsql',2,1,1,'prod','1111','2025-07-29 21:36:11.147','2025-09-06 15:12:44.586');
/*!40000 ALTER TABLE `cmdb_sql` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `cmdb_sql_log`
--

DROP TABLE IF EXISTS `cmdb_sql_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `cmdb_sql_log` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `instance_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `database` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `operation_type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `sql_content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `exec_user` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `ip` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `scanned_rows` bigint DEFAULT '0',
  `affected_rows` bigint DEFAULT '0',
  `execution_time` bigint DEFAULT '0',
  `returned_rows` bigint DEFAULT '0',
  `result` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `query_time` datetime(3) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_cmdb_sql_log_query_time` (`query_time`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cmdb_sql_log`
--

LOCK TABLES `cmdb_sql_log` WRITE;
/*!40000 ALTER TABLE `cmdb_sql_log` DISABLE KEYS */;
/*!40000 ALTER TABLE `cmdb_sql_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `cmdb_sql_records`
--

DROP TABLE IF EXISTS `cmdb_sql_records`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `cmdb_sql_records` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `instance_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `database` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `operation_type` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `sql_content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `exec_user` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `scanned_rows` bigint DEFAULT '0',
  `affected_rows` bigint DEFAULT '0',
  `execution_time` bigint DEFAULT '0',
  `returned_rows` bigint DEFAULT '0',
  `result` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `query_time` datetime(3) NOT NULL,
  `name` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_cmdb_sql_records_query_time` (`query_time`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=19 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cmdb_sql_records`
--

LOCK TABLES `cmdb_sql_records` WRITE;
/*!40000 ALTER TABLE `cmdb_sql_records` DISABLE KEYS */;
INSERT INTO `cmdb_sql_records` VALUES (1,'8.130.14.34:3306','gin-api','SELECT','select * from cmdb_group;','anonymous',14,0,403,14,'SUCCESS','2025-07-29 11:14:13.686',''),(2,'8.130.14.30','saas3-mysql','SELECT','select * from cmdb_group;','',0,0,50,10,'SUCCESS','2025-07-30 13:29:24.409',''),(3,'8.130.14.34','gin-api','SELECT','select * from cmdb_group;','',0,0,54,14,'SUCCESS','2025-07-30 13:58:13.386',''),(4,'8.130.14.34','saas3-mysql','INSERT','UPDATE `cmdb_group`SET `name` = \'sql测试组0000\' WHERE `parent_id` = 17AND `name` = \'sql测试组\';','',0,1,80,0,'SUCCESS','2025-07-30 14:00:44.370',''),(5,'8.130.14.34','saas3-mysql','INSERT','UPDATE `cmdb_group` SET `name` = \'test123\' WHERE `id` = 22;','',0,1,80,0,'SUCCESS','2025-07-30 14:04:30.684',''),(6,'8.130.14.34','saas3-mysql','INSERT','UPDATE `cmdb_group` SET `name` = \'test123111\' WHERE `id` = 22;','',0,1,80,0,'SUCCESS','2025-07-30 14:06:34.692',''),(7,'8.130.14.34','gin-api','SELECT','UPDATE `cmdb_group` SET `name` = \'test123111\' WHERE `id` = 22;','',0,0,122,0,'SUCCESS','2025-07-30 14:07:44.151',''),(8,'8.130.14.34','gin-api','SELECT','UPDATE `cmdb_group` SET `name` = \'test001\' WHERE `id` = 22;','',0,0,55,0,'SUCCESS','2025-07-30 14:11:42.626',''),(9,'8.130.14.34','gin-api','SELECT','select * from cmdb_group;','',0,0,67,15,'SUCCESS','2025-07-30 15:24:57.109',''),(10,'8.130.14.34','gin-api','SELECT','select * from cmdb_host;','',0,0,114,75,'SUCCESS','2025-07-30 15:25:19.542',''),(11,'8.130.14.34','gin-api','SELECT','select * from cmdb_group;','',0,0,56,15,'SUCCESS','2025-07-30 15:29:16.507',''),(12,'8.130.14.34','saas3-mysql','EXECUTE','create databases  db;','',0,1,100,0,'SUCCESS','2025-07-30 15:34:31.246',''),(13,'8.130.14.34','gin-api','SELECT','select id,name from cmdb_group;','',0,0,51,15,'SUCCESS','2025-07-30 17:18:52.168',''),(14,'8.130.14.34','gin-api','SELECT','select id,name from cmdb_group;','',0,0,53,15,'SUCCESS','2025-07-30 17:29:36.465',''),(15,'8.130.14.34','gin-api','SELECT','select id,name from  cmdb_group;','',0,0,56,15,'SUCCESS','2025-07-30 21:06:07.136',''),(16,'8.130.14.34','gin-api','SELECT','select id,name from  cmdb_group;','',0,0,64,15,'SUCCESS','2025-07-30 21:12:04.886',''),(17,'8.130.14.34','gin-api','SELECT','select id,name from  cmdb_group;','admin',0,0,51,15,'SUCCESS','2025-07-30 21:23:55.532',''),(18,'8.130.14.34','gin-api','SELECT','select id,name from  cmdb_group;','zhangsan',0,0,54,15,'SUCCESS','2025-07-30 21:26:51.642','');
/*!40000 ALTER TABLE `cmdb_sql_records` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `config_account`
--

DROP TABLE IF EXISTS `config_account`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `config_account` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `alias` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `host` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `port` bigint NOT NULL,
  `name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `password` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `type` bigint NOT NULL,
  `remark` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `config_account`
--

LOCK TABLES `config_account` WRITE;
/*!40000 ALTER TABLE `config_account` DISABLE KEYS */;
/*!40000 ALTER TABLE `config_account` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `config_ecsauth`
--

DROP TABLE IF EXISTS `config_ecsauth`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `config_ecsauth` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键''',
  `name` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''凭证名称''',
  `type` bigint NOT NULL COMMENT '''认证类型:1->密码,2->私钥,3->公钥(免认证)''',
  `username` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''用户名''',
  `password` longtext COMMENT '''密码(type=1时使用)''',
  `public_key` text COMMENT '''私钥内容(type=2时使用，字段名历史原因)''',
  `create_time` datetime(3) NOT NULL COMMENT '''创建时间''',
  `remark` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''备注''',
  `bastion_host` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''跳板机IP''',
  `bastion_port` bigint DEFAULT '22' COMMENT '''跳板机端口''',
  `bastion_username` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''跳板机用户名''',
  `bastion_auth_type` bigint DEFAULT NULL COMMENT '''跳板机认证类型:1->密码,2->私钥,3->免密''',
  `bastion_password` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''跳板机密码''',
  `bastion_private_key` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''跳板机私钥''',
  `target_auth_type` bigint DEFAULT NULL COMMENT '''目标机器认证类型:1->密码,2->私钥,3->免密''',
  `target_password` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''目标机器密码''',
  `target_private_key` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''目标机器私钥''',
  `port` bigint DEFAULT '22' COMMENT '''端口号''',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=33 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `config_ecsauth`
--

LOCK TABLES `config_ecsauth` WRITE;
/*!40000 ALTER TABLE `config_ecsauth` DISABLE KEYS */;
INSERT INTO `config_ecsauth` VALUES (32,'免密认证',3,'root','','','2026-02-19 19:40:02.282','123',NULL,22,NULL,NULL,NULL,NULL,NULL,NULL,NULL,22);
/*!40000 ALTER TABLE `config_ecsauth` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `config_keymanage`
--

DROP TABLE IF EXISTS `config_keymanage`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `config_keymanage` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `key_type` bigint NOT NULL,
  `key_id` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `key_secret` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `remark` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `key_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `config_keymanage`
--

LOCK TABLES `config_keymanage` WRITE;
/*!40000 ALTER TABLE `config_keymanage` DISABLE KEYS */;
/*!40000 ALTER TABLE `config_keymanage` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `config_sync_schedule`
--

DROP TABLE IF EXISTS `config_sync_schedule`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `config_sync_schedule` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `cron_expr` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `key_types` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `status` bigint NOT NULL DEFAULT '1',
  `last_run_time` datetime(3) DEFAULT NULL,
  `next_run_time` datetime(3) DEFAULT NULL,
  `remark` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `sync_log` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `config_sync_schedule`
--

LOCK TABLES `config_sync_schedule` WRITE;
/*!40000 ALTER TABLE `config_sync_schedule` DISABLE KEYS */;
/*!40000 ALTER TABLE `config_sync_schedule` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `db`
--

DROP TABLE IF EXISTS `db`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `db` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键''',
  `code` varchar(36) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''数据库编码''',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''数据库名称''',
  `database` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''数据库名(多个用空格分隔)''',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''备注''',
  `instance_id` bigint unsigned NOT NULL COMMENT '''实例ID''',
  `instance_code` varchar(36) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''实例编码''',
  `status` bigint DEFAULT '1' COMMENT '''状态:1->启用,2->禁用''',
  `create_time` datetime(3) NOT NULL COMMENT '''创建时间''',
  `update_time` datetime(3) DEFAULT NULL COMMENT '''更新时间''',
  `creator` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''创建人''',
  `creator_id` bigint unsigned DEFAULT NULL COMMENT '''创建人ID''',
  `modifier` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''修改人''',
  `modifier_id` bigint unsigned DEFAULT NULL COMMENT '''修改人ID''',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `idx_db_code` (`code`) USING BTREE,
  KEY `idx_db_instance_id` (`instance_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `db`
--

LOCK TABLES `db` WRITE;
/*!40000 ALTER TABLE `db` DISABLE KEYS */;
/*!40000 ALTER TABLE `db` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `db_es_instance`
--

DROP TABLE IF EXISTS `db_es_instance`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `db_es_instance` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键''',
  `code` varchar(36) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''实例编码''',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''实例名称''',
  `protocol` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT 'http' COMMENT '''协议:http|https''',
  `host` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''主机''',
  `port` bigint DEFAULT '9200' COMMENT '''端口''',
  `username` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''用户名(可选)''',
  `password` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '''密码(加密)''',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''备注''',
  `ssh_tunnel_machine_id` bigint unsigned DEFAULT '0' COMMENT '''SSH隧道机器ID''',
  `status` bigint DEFAULT '1' COMMENT '''状态:1->启用,2->禁用''',
  `create_time` datetime(3) NOT NULL COMMENT '''创建时间''',
  `update_time` datetime(3) DEFAULT NULL COMMENT '''更新时间''',
  `creator` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''创建人''',
  `creator_id` bigint unsigned DEFAULT NULL COMMENT '''创建人ID''',
  `modifier` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''修改人''',
  `modifier_id` bigint unsigned DEFAULT NULL COMMENT '''修改人ID''',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `idx_db_es_instance_code` (`code`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `db_es_instance`
--

LOCK TABLES `db_es_instance` WRITE;
/*!40000 ALTER TABLE `db_es_instance` DISABLE KEYS */;
/*!40000 ALTER TABLE `db_es_instance` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `db_export_task`
--

DROP TABLE IF EXISTS `db_export_task`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `db_export_task` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键''',
  `task_id` varchar(36) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''任务ID''',
  `db_id` bigint unsigned NOT NULL COMMENT '''数据库ID''',
  `db_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''数据库名''',
  `export_type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''导出类型:structure/full''',
  `status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''状态''',
  `file_path` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''文件路径''',
  `file_size` bigint DEFAULT NULL COMMENT '''文件大小(字节)''',
  `error_message` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''错误信息''',
  `start_time` datetime(3) DEFAULT NULL COMMENT '''开始时间''',
  `end_time` datetime(3) DEFAULT NULL COMMENT '''结束时间''',
  `create_time` datetime(3) NOT NULL COMMENT '''创建时间''',
  `update_time` datetime(3) DEFAULT NULL COMMENT '''更新时间''',
  `creator` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''创建人''',
  `creator_id` bigint unsigned DEFAULT NULL COMMENT '''创建人ID''',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `idx_db_export_task_task_id` (`task_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `db_export_task`
--

LOCK TABLES `db_export_task` WRITE;
/*!40000 ALTER TABLE `db_export_task` DISABLE KEYS */;
INSERT INTO `db_export_task` VALUES (1,'91c691d3-8e72-40d3-9de7-4453057ca1ee',1,'devops','full','completed','data/exports/devops_full_1765299147.sql',689403,'','2025-12-10 00:52:27.431','2025-12-10 00:52:40.023','2025-12-10 00:52:27.082','2025-12-10 00:52:40.023','',0),(2,'e1c6f2c1-7a85-43fb-8255-3d9a9fd229cb',1,'mayfly-go','full','completed','data/exports/mayfly-go_full_1765333697.sql',100718,'','2025-12-10 10:28:17.221','2025-12-10 10:28:28.959','2025-12-10 10:28:16.888','2025-12-10 10:28:28.959','',0),(3,'d788c68a-bab0-4d4c-8b6f-467d672876b1',2,'database_name','full','completed','data/exports/database_name_full_1765342425.sql',1307,'','2025-12-10 12:53:45.551','2025-12-10 12:53:46.084','2025-12-10 12:53:45.207','2025-12-10 12:53:46.084','',0),(4,'c6c4abdf-5dff-4a43-ba4c-64fac0dfc6de',1,'gin-api','full','completed','data/exports/gin-api_full_1765636630.sql',652608,'','2025-12-13 22:37:10.521','2025-12-13 22:37:24.556','2025-12-13 22:37:10.153','2025-12-13 22:37:24.556','',0);
/*!40000 ALTER TABLE `db_export_task` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `db_instance`
--

DROP TABLE IF EXISTS `db_instance`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `db_instance` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键''',
  `code` varchar(36) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''实例编码''',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''实例名称''',
  `type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''数据库类型:mysql,postgres,oracle等''',
  `host` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''主机地址''',
  `port` bigint NOT NULL COMMENT '''端口''',
  `network` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT 'tcp' COMMENT '''网络类型:tcp,unix''',
  `params` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''连接参数''',
  `username` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''用户名''',
  `password` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''密码(加密)''',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''备注''',
  `ssh_tunnel_machine_id` bigint unsigned DEFAULT '0' COMMENT '''SSH隧道机器ID''',
  `status` bigint DEFAULT '1' COMMENT '''状态:1->启用,2->禁用''',
  `create_time` datetime(3) NOT NULL COMMENT '''创建时间''',
  `update_time` datetime(3) DEFAULT NULL COMMENT '''更新时间''',
  `creator` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''创建人''',
  `creator_id` bigint unsigned DEFAULT NULL COMMENT '''创建人ID''',
  `modifier` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''修改人''',
  `modifier_id` bigint unsigned DEFAULT NULL COMMENT '''修改人ID''',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `idx_db_instance_code` (`code`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `db_instance`
--

LOCK TABLES `db_instance` WRITE;
/*!40000 ALTER TABLE `db_instance` DISABLE KEYS */;
/*!40000 ALTER TABLE `db_instance` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Temporary view structure for view `db_instance_all`
--

DROP TABLE IF EXISTS `db_instance_all`;
/*!50001 DROP VIEW IF EXISTS `db_instance_all`*/;
SET @saved_cs_client     = @@character_set_client;
/*!50503 SET character_set_client = utf8mb4 */;
/*!50001 CREATE VIEW `db_instance_all` AS SELECT 
 1 AS `id`,
 1 AS `code`,
 1 AS `name`,
 1 AS `db_type`,
 1 AS `sub_type`,
 1 AS `host`,
 1 AS `port`,
 1 AS `username`,
 1 AS `password`,
 1 AS `remark`,
 1 AS `ssh_tunnel_machine_id`,
 1 AS `status`,
 1 AS `create_time`,
 1 AS `update_time`,
 1 AS `creator`,
 1 AS `creator_id`,
 1 AS `modifier`,
 1 AS `modifier_id`,
 1 AS `connection_config`*/;
SET character_set_client = @saved_cs_client;

--
-- Table structure for table `db_mongo_instance`
--

DROP TABLE IF EXISTS `db_mongo_instance`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `db_mongo_instance` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键''',
  `code` varchar(36) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''实例编码''',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''实例名称''',
  `uri` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''连接URI''',
  `ssh_tunnel_machine_id` bigint unsigned DEFAULT '0' COMMENT '''SSH隧道机器ID''',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''备注''',
  `status` bigint DEFAULT '1' COMMENT '''状态:1->启用,2->禁用''',
  `create_time` datetime(3) NOT NULL COMMENT '''创建时间''',
  `update_time` datetime(3) DEFAULT NULL COMMENT '''更新时间''',
  `creator` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''创建人''',
  `creator_id` bigint unsigned DEFAULT NULL COMMENT '''创建人ID''',
  `modifier` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''修改人''',
  `modifier_id` bigint unsigned DEFAULT NULL COMMENT '''修改人ID''',
  `type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT 'mongodb' COMMENT '''MongoDB类型:mongodb,mongodb-atlas等''',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `idx_db_mongo_instance_code` (`code`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `db_mongo_instance`
--

LOCK TABLES `db_mongo_instance` WRITE;
/*!40000 ALTER TABLE `db_mongo_instance` DISABLE KEYS */;
/*!40000 ALTER TABLE `db_mongo_instance` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `db_redis_instance`
--

DROP TABLE IF EXISTS `db_redis_instance`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `db_redis_instance` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键''',
  `code` varchar(36) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''实例编码''',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''实例名称''',
  `mode` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''模式:standalone,cluster,sentinel''',
  `host` varchar(300) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''主机: standalone为host:port, cluster为逗号分隔, sentinel为master=hosts''',
  `port` bigint DEFAULT '0' COMMENT '''端口(standalone可用)''',
  `db` bigint DEFAULT '0' COMMENT '''默认库号''',
  `username` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''用户名(可选)''',
  `password` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''密码(加密)''',
  `redis_node_password` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '''节点密码(仅sentinel)''',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''备注''',
  `ssh_tunnel_machine_id` bigint unsigned DEFAULT '0' COMMENT '''SSH隧道机器ID''',
  `status` bigint DEFAULT '1' COMMENT '''状态:1->启用,2->禁用''',
  `create_time` datetime(3) NOT NULL COMMENT '''创建时间''',
  `update_time` datetime(3) DEFAULT NULL COMMENT '''更新时间''',
  `creator` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''创建人''',
  `creator_id` bigint unsigned DEFAULT NULL COMMENT '''创建人ID''',
  `modifier` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''修改人''',
  `modifier_id` bigint unsigned DEFAULT NULL COMMENT '''修改人ID''',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `idx_db_redis_instance_code` (`code`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `db_redis_instance`
--

LOCK TABLES `db_redis_instance` WRITE;
/*!40000 ALTER TABLE `db_redis_instance` DISABLE KEYS */;
/*!40000 ALTER TABLE `db_redis_instance` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `db_sql`
--

DROP TABLE IF EXISTS `db_sql`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `db_sql` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键''',
  `db_id` bigint unsigned NOT NULL COMMENT '''数据库ID''',
  `db` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''数据库名''',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''SQL名称''',
  `type` bigint DEFAULT '1' COMMENT '''类型:1->查询,2->更新''',
  `sql` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''SQL语句''',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''备注''',
  `create_time` datetime(3) NOT NULL COMMENT '''创建时间''',
  `update_time` datetime(3) DEFAULT NULL COMMENT '''更新时间''',
  `creator` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''创建人''',
  `creator_id` bigint unsigned DEFAULT NULL COMMENT '''创建人ID''',
  `modifier` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''修改人''',
  `modifier_id` bigint unsigned DEFAULT NULL COMMENT '''修改人ID''',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_db_sql_db_id` (`db_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `db_sql`
--

LOCK TABLES `db_sql` WRITE;
/*!40000 ALTER TABLE `db_sql` DISABLE KEYS */;
/*!40000 ALTER TABLE `db_sql` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `db_sql_exec`
--

DROP TABLE IF EXISTS `db_sql_exec`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `db_sql_exec` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键''',
  `db_id` bigint unsigned NOT NULL COMMENT '''数据库ID''',
  `db_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''数据库名''',
  `table_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''表名''',
  `type` tinyint NOT NULL COMMENT '''类型:1->查询,2->插入,3->更新,4->删除,5->DDL''',
  `sql` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''SQL语句''',
  `old_value` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''旧值(用于回滚)''',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''备注''',
  `status` tinyint NOT NULL COMMENT '''状态:1->成功,2->失败''',
  `res` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''执行结果''',
  `exec_time` bigint DEFAULT NULL COMMENT '''执行时长(ms)''',
  `create_time` datetime(3) NOT NULL COMMENT '''创建时间''',
  `creator` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''创建人''',
  `creator_id` bigint unsigned DEFAULT NULL COMMENT '''创建人ID''',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_db_sql_exec_db_id` (`db_id`) USING BTREE,
  KEY `idx_db_sql_exec_create_time` (`create_time`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=524 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `db_sql_exec`
--

LOCK TABLES `db_sql_exec` WRITE;
/*!40000 ALTER TABLE `db_sql_exec` DISABLE KEYS */;
INSERT INTO `db_sql_exec` VALUES (1,1,'gin-api','APP_APPLICATION',0,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',2,'数据库不存在',110,'2025-11-29 18:05:04.883','',0),(2,1,'gin-api','APP_APPLICATION',0,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',2,'数据库不存在',112,'2025-11-29 18:05:05.487','',0),(3,1,'gin-api','CMDB_GROUP',0,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',2,'数据库不存在',119,'2025-11-29 18:05:10.402','',0),(4,1,'gin-api','CMDB_GROUP',0,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',2,'数据库不存在',116,'2025-11-29 18:05:10.598','',0),(5,1,'gin-api','CMDB_SQL',0,'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0','','',2,'数据库不存在',122,'2025-11-29 18:05:17.424','',0),(6,1,'gin-api','CMDB_SQL',0,'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0','','',2,'数据库不存在',136,'2025-11-29 18:05:17.580','',0),(7,1,'gin-api','APP_JENKINS_ENV',0,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',2,'数据库不存在',119,'2025-11-29 18:11:06.550','',0),(8,1,'gin-api','APP_JENKINS_ENV',0,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',2,'数据库不存在',493,'2025-11-29 18:11:06.802','',0),(9,1,'gin-api','APP_JENKINS_ENV',0,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',2,'数据库不存在',127,'2025-11-29 18:11:15.554','',0),(10,1,'gin-api','APP_JENKINS_ENV',0,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',2,'数据库不存在',110,'2025-11-29 18:11:15.687','',0),(11,1,'gin-api','APP_SERVICE_RELEASE_ITEM',1,'SELECT * FROM app_service_release_item LIMIT 20 OFFSET 0','','',1,'执行成功',548,'2025-11-29 18:28:12.401','',0),(12,1,'gin-api','APP_SERVICE_RELEASE_ITEM',1,'SELECT * FROM app_service_release_item LIMIT 20 OFFSET 0','','',1,'执行成功',1083,'2025-11-29 18:28:12.656','',0),(13,1,'gin-api','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',277,'2025-11-29 18:28:25.668','',0),(14,1,'gin-api','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',281,'2025-11-29 18:28:25.951','',0),(15,1,'devops','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',308,'2025-11-29 18:31:23.236','',0),(16,1,'devops','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',346,'2025-11-29 18:31:23.606','',0),(17,1,'gin-api','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',298,'2025-11-29 18:35:52.379','',0),(18,1,'gin-api','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',303,'2025-11-29 18:35:52.538','',0),(19,1,'gin-api','APP_APPLICATION',1,'SELECT COUNT(*) as total FROM app_application','','',1,'执行成功',266,'2025-11-29 18:35:53.254','',0),(20,1,'gin-api','APP_APPLICATION',1,'SELECT COUNT(*) as total FROM app_application','','',1,'执行成功',259,'2025-11-29 18:35:53.562','',0),(21,1,'devops','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',283,'2025-11-29 18:37:11.404','',0),(22,1,'devops','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',320,'2025-11-29 18:37:12.000','',0),(23,1,'devops','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',273,'2025-11-29 18:37:11.713','',0),(24,1,'devops','APP_APPLICATION',1,'SELECT COUNT(*) as total FROM app_application','','',1,'执行成功',299,'2025-11-29 18:37:12.382','',0),(25,1,'devops','APP_APPLICATION',1,'SELECT COUNT(*) as total FROM app_application','','',1,'执行成功',255,'2025-11-29 18:37:12.911','',0),(26,1,'devops','APP_APPLICATION',1,'SELECT COUNT(*) as total FROM app_application','','',1,'执行成功',282,'2025-11-29 18:37:13.296','',0),(27,1,'devops','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',248,'2025-11-29 18:37:58.462','',0),(28,1,'devops','APP_APPLICATION',1,'SELECT COUNT(*) as total FROM app_application','','',1,'执行成功',260,'2025-11-29 18:37:59.326','',0),(29,1,'devops','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',275,'2025-11-29 18:38:46.987','',0),(30,1,'devops','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',277,'2025-11-29 18:38:47.151','',0),(31,1,'devops','APP_APPLICATION',1,'SELECT COUNT(*) as total FROM app_application','','',1,'执行成功',265,'2025-11-29 18:38:47.925','',0),(32,1,'devops','APP_APPLICATION',1,'SELECT COUNT(*) as total FROM app_application','','',1,'执行成功',266,'2025-11-29 18:38:48.246','',0),(33,1,'gin-api','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',265,'2025-11-29 18:39:00.606','',0),(34,1,'gin-api','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',304,'2025-11-29 18:39:00.775','',0),(35,1,'gin-api','APP_APPLICATION',1,'SELECT COUNT(*) as total FROM app_application','','',1,'执行成功',287,'2025-11-29 18:39:01.496','',0),(36,1,'gin-api','APP_APPLICATION',1,'SELECT COUNT(*) as total FROM app_application','','',1,'执行成功',263,'2025-11-29 18:39:01.812','',0),(37,1,'devops','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',403,'2025-11-29 18:48:01.028','',0),(38,1,'devops','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',335,'2025-11-29 18:48:01.099','',0),(39,1,'devops','APP_APPLICATION',1,'SELECT COUNT(*) as total FROM app_application','','',1,'执行成功',10295,'2025-11-29 18:48:15.394','',0),(40,1,'devops','APP_APPLICATION',1,'SELECT COUNT(*) as total FROM app_application','','',1,'执行成功',2650,'2025-11-29 18:48:14.638','',0),(41,1,'devops','CMDB_HOST',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',6494,'2025-11-29 18:48:15.978','',0),(42,1,'devops','CMDB_HOST',1,'SELECT COUNT(*) as total FROM cmdb_host','','',1,'执行成功',293,'2025-11-29 18:48:16.898','',0),(43,1,'devops','CMDB_HOST',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',3692,'2025-11-29 18:48:12.891','',0),(44,1,'devops','CMDB_HOST',1,'SELECT COUNT(*) as total FROM cmdb_host','','',1,'执行成功',281,'2025-11-29 18:48:21.045','',0),(45,1,'devops','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',340,'2025-11-29 18:50:12.488','',0),(46,1,'devops','APP_APPLICATION',1,'SELECT COUNT(*) as total FROM app_application','','',1,'执行成功',284,'2025-11-29 18:50:13.429','',0),(47,1,'devops','APP_JENKINS_ENV',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',304,'2025-11-29 18:54:02.563','',0),(48,1,'devops','APP_JENKINS_ENV',1,'SELECT COUNT(*) as total FROM app_jenkins_env','','',1,'执行成功',294,'2025-11-29 18:54:03.463','',0),(49,1,'devops','CMDB_GROUP',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',328,'2025-11-29 18:54:07.320','',0),(50,1,'devops','CMDB_GROUP',1,'SELECT COUNT(*) as total FROM cmdb_group','','',1,'执行成功',313,'2025-11-29 18:54:08.231','',0),(51,1,'devops','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',282,'2025-11-29 18:55:28.076','',0),(52,1,'devops','APP_APPLICATION',1,'SELECT COUNT(*) as total FROM app_application','','',1,'执行成功',280,'2025-11-29 18:55:28.949','',0),(53,1,'devops','APP_JENKINS_ENV',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',267,'2025-11-29 18:55:29.997','',0),(54,1,'devops','APP_JENKINS_ENV',1,'SELECT COUNT(*) as total FROM app_jenkins_env','','',1,'执行成功',277,'2025-11-29 18:55:30.867','',0),(55,1,'devops','APP_JENKINS_ENV',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',282,'2025-11-29 18:56:33.893','',0),(56,1,'devops','APP_JENKINS_ENV',1,'SELECT COUNT(*) as total FROM app_jenkins_env','','',1,'执行成功',275,'2025-11-29 18:56:34.766','',0),(57,1,'devops','CMDB_GROUP',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',309,'2025-11-29 18:56:37.254','',0),(58,1,'devops','CMDB_GROUP',1,'SELECT COUNT(*) as total FROM cmdb_group','','',1,'执行成功',293,'2025-11-29 18:56:38.144','',0),(59,1,'devops','CMDB_HOST',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',276,'2025-11-29 18:57:16.850','',0),(60,1,'devops','CMDB_HOST',1,'SELECT COUNT(*) as total FROM cmdb_host','','',1,'执行成功',284,'2025-11-29 18:57:17.734','',0),(61,1,'devops','CMDB_SQL_LOG',1,'SELECT * FROM cmdb_sql_log LIMIT 20 OFFSET 0','','',1,'执行成功',402,'2025-11-29 18:57:49.376','',0),(62,1,'devops','CMDB_SQL_LOG',1,'SELECT COUNT(*) as total FROM cmdb_sql_log','','',1,'执行成功',274,'2025-11-29 18:57:50.233','',0),(63,1,'devops','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',315,'2025-11-29 19:00:49.626','',0),(64,1,'devops','APP_APPLICATION',1,'SELECT COUNT(*) as total FROM app_application','','',1,'执行成功',330,'2025-11-29 19:00:50.546','',0),(65,1,'devops','CMDB_HOST',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',484,'2025-11-29 19:00:57.485','',0),(66,1,'devops','CMDB_HOST',1,'SELECT COUNT(*) as total FROM cmdb_host','','',1,'执行成功',656,'2025-11-29 19:00:59.337','',0),(67,1,'devops','CMDB_SQL',1,'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0','','',1,'执行成功',432,'2025-11-29 19:01:33.694','',0),(68,1,'devops','CMDB_SQL',1,'SELECT COUNT(*) as total FROM cmdb_sql','','',1,'执行成功',285,'2025-11-29 19:01:34.585','',0),(69,1,'devops','CONFIG_ACCOUNT',1,'SELECT * FROM config_account LIMIT 20 OFFSET 0','','',1,'执行成功',300,'2025-11-29 19:01:52.297','',0),(70,1,'devops','CONFIG_ACCOUNT',1,'SELECT COUNT(*) as total FROM config_account','','',1,'执行成功',282,'2025-11-29 19:01:53.222','',0),(71,1,'devops','CMDB_SQL_LOG',1,'SELECT * FROM cmdb_sql_log LIMIT 20 OFFSET 0','','',1,'执行成功',442,'2025-11-29 19:02:12.743','',0),(72,1,'devops','CMDB_SQL_LOG',1,'SELECT COUNT(*) as total FROM cmdb_sql_log','','',1,'执行成功',388,'2025-11-29 19:02:13.732','',0),(73,1,'devops','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',321,'2025-11-29 19:04:24.978','',0),(74,1,'devops','APP_APPLICATION',1,'SELECT COUNT(*) as total FROM app_application','','',1,'执行成功',291,'2025-11-29 19:04:25.863','',0),(75,1,'devops','CMDB_HOST',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',286,'2025-11-29 19:05:02.134','',0),(76,1,'devops','CMDB_HOST',1,'SELECT COUNT(*) as total FROM cmdb_host','','',1,'执行成功',309,'2025-11-29 19:05:03.008','',0),(77,1,'gin-api','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',361,'2025-11-29 19:06:09.051','',0),(78,1,'gin-api','APP_APPLICATION',1,'SELECT COUNT(*) as total FROM app_application','','',1,'执行成功',291,'2025-11-29 19:06:10.002','',0),(79,1,'gin-api','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',1109,'2025-11-29 19:06:32.787','',0),(80,1,'gin-api','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',1586,'2025-11-29 19:06:33.149','',0),(81,1,'gin-api','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',303,'2025-11-29 19:06:33.300','',0),(82,1,'gin-api','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',552,'2025-11-29 19:06:33.428','',0),(83,1,'gin-api','APP_APPLICATION',1,'SELECT COUNT(*) as total FROM app_application','','',1,'执行成功',307,'2025-11-29 19:06:33.766','',0),(84,1,'gin-api','APP_APPLICATION',1,'SELECT COUNT(*) as total FROM app_application','','',1,'执行成功',254,'2025-11-29 19:06:34.069','',0),(85,1,'gin-api','APP_APPLICATION',1,'SELECT COUNT(*) as total FROM app_application','','',1,'执行成功',553,'2025-11-29 19:06:34.715','',0),(86,1,'gin-api','APP_APPLICATION',1,'SELECT COUNT(*) as total FROM app_application','','',1,'执行成功',302,'2025-11-29 19:06:35.225','',0),(87,1,'gin-api','APP_JENKINS_ENV',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',295,'2025-11-29 19:06:36.173','',0),(88,1,'gin-api','APP_JENKINS_ENV',1,'SELECT COUNT(*) as total FROM app_jenkins_env','','',1,'执行成功',276,'2025-11-29 19:06:37.028','',0),(89,1,'devops','CMDB_SQL',1,'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0','','',1,'执行成功',275,'2025-11-29 19:08:02.114','',0),(90,1,'devops','CMDB_SQL',1,'SELECT COUNT(*) as total FROM cmdb_sql','','',1,'执行成功',331,'2025-11-29 19:08:03.077','',0),(91,1,'devops','CONFIG_ACCOUNT',1,'SELECT * FROM config_account LIMIT 20 OFFSET 0','','',1,'执行成功',324,'2025-11-29 19:08:59.699','',0),(92,1,'devops','CMDB_SQL_RECORDS',1,'SELECT * FROM cmdb_sql_records LIMIT 20 OFFSET 0','','',1,'执行成功',299,'2025-11-29 19:08:59.964','',0),(93,1,'devops','CMDB_SQL_RECORDS',1,'SELECT * FROM cmdb_sql_records LIMIT 20 OFFSET 0','','',1,'执行成功',275,'2025-11-29 19:09:00.054','',0),(94,1,'devops','CMDB_SQL_RECORDS',1,'SELECT * FROM cmdb_sql_records LIMIT 20 OFFSET 0','','',1,'执行成功',259,'2025-11-29 19:09:00.169','',0),(95,1,'devops','CONFIG_ACCOUNT',1,'SELECT * FROM config_account LIMIT 20 OFFSET 0','','',1,'执行成功',360,'2025-11-29 19:08:59.813','',0),(96,1,'devops','CONFIG_ACCOUNT',1,'SELECT COUNT(*) as total FROM config_account','','',1,'执行成功',292,'2025-11-29 19:09:00.602','',0),(97,1,'devops','CMDB_SQL_RECORDS',1,'SELECT COUNT(*) as total FROM cmdb_sql_records','','',1,'执行成功',303,'2025-11-29 19:09:00.855','',0),(98,1,'devops','CONFIG_ACCOUNT',1,'SELECT COUNT(*) as total FROM config_account','','',1,'执行成功',326,'2025-11-29 19:09:01.089','',0),(99,1,'devops','CMDB_SQL_RECORDS',1,'SELECT COUNT(*) as total FROM cmdb_sql_records','','',1,'执行成功',295,'2025-11-29 19:09:01.217','',0),(100,1,'devops','CMDB_SQL_RECORDS',1,'SELECT COUNT(*) as total FROM cmdb_sql_records','','',1,'执行成功',334,'2025-11-29 19:09:01.611','',0),(101,1,'devops','CMDB_SQL',1,'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0','','',1,'执行成功',280,'2025-11-29 19:09:14.451','',0),(102,1,'devops','CMDB_SQL',1,'SELECT COUNT(*) as total FROM cmdb_sql','','',1,'执行成功',278,'2025-11-29 19:09:15.298','',0),(103,1,'devops','CMDB_HOST',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',1029,'2025-11-29 19:09:29.074','',0),(104,1,'devops','CMDB_HOST',1,'SELECT COUNT(*) as total FROM cmdb_host','','',1,'执行成功',6538,'2025-11-29 19:09:36.486','',0),(105,1,'devops','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',264,'2025-11-29 19:23:51.118','',0),(106,1,'devops','APP_APPLICATION',1,'SELECT COUNT(*) as total FROM app_application','','',1,'执行成功',296,'2025-11-29 19:23:51.990','',0),(107,1,'devops','CMDB_SQL',1,'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0','','',1,'执行成功',295,'2025-11-29 19:23:59.137','',0),(108,1,'devops','CMDB_SQL',1,'SELECT COUNT(*) as total FROM cmdb_sql','','',1,'执行成功',279,'2025-11-29 19:24:00.023','',0),(109,1,'devops','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 100;','','',1,'执行成功',280,'2025-11-29 19:25:10.430','',0),(110,1,'devops','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',284,'2025-11-29 19:49:42.427','',0),(111,1,'devops','APP_APPLICATION',1,'SELECT COUNT(*) as total FROM app_application','','',1,'执行成功',287,'2025-11-29 19:49:43.313','',0),(112,1,'devops','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',239,'2025-11-29 19:57:17.599','',0),(113,1,'devops','CMDB_SQL_LOG',1,'SELECT * FROM cmdb_sql_log LIMIT 20 OFFSET 0','','',1,'执行成功',213,'2025-11-29 19:57:40.935','',0),(114,1,'devops','CMDB_SQL',1,'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0','','',1,'执行成功',213,'2025-11-29 19:57:40.616','',0),(115,1,'devops','CMDB_SQL_LOG',1,'SELECT * FROM cmdb_sql_log LIMIT 20 OFFSET 0','','',1,'执行成功',212,'2025-11-29 19:57:40.797','',0),(116,1,'devops','CMDB_SQL_LOG',1,'SELECT * FROM cmdb_sql_log LIMIT 20 OFFSET 0','','',1,'执行成功',210,'2025-11-29 19:57:59.070','',0),(117,1,'devops','CMDB_SQL_RECORDS',1,'SELECT * FROM cmdb_sql_records LIMIT 20 OFFSET 0','','',1,'执行成功',248,'2025-11-29 19:58:15.888','',0),(118,1,'devops','CMDB_HOST',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',307,'2025-11-29 19:58:47.223','',0),(119,1,'devops','CMDB_HOST;',1,'SELECT * FROM cmdb_host;','','',1,'执行成功',212,'2025-11-29 20:00:09.885','',0),(120,1,'devops','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',263,'2025-11-29 20:03:26.517','',0),(121,1,'devops','CMDB_HOST',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',213,'2025-11-29 20:03:44.248','',0),(122,1,'devops','CMDB_HOST',1,'SELECT * FROM cmdb_host LIMIT 100;','','',1,'执行成功',220,'2025-11-29 20:04:04.538','',0),(123,1,'devops','TASK_WORK',1,'SELECT * FROM task_work LIMIT 20 OFFSET 0','','',1,'执行成功',234,'2025-11-29 20:04:29.401','',0),(124,1,'devops','K8S_CLUSTER',1,'SELECT * FROM k8s_cluster LIMIT 20 OFFSET 0','','',1,'执行成功',247,'2025-11-29 20:05:20.336','',0),(125,1,'devops','CONFIG_ECSAUTH',1,'SELECT * FROM config_ecsauth LIMIT 20 OFFSET 0','','',1,'执行成功',515,'2025-11-29 20:07:54.567','',0),(126,1,'devops','CONFIG_ACCOUNT',1,'SELECT * FROM config_account LIMIT 20 OFFSET 0','','',1,'执行成功',219,'2025-11-29 20:08:12.059','',0),(127,1,'devops','CONFIG_ECSAUTH',1,'SELECT * FROM config_ecsauth LIMIT 20 OFFSET 0','','',1,'执行成功',386,'2025-11-29 20:08:56.093','',0),(128,1,'devops','CONFIG_ACCOUNT',1,'SELECT * FROM config_account LIMIT 20 OFFSET 0','','',1,'执行成功',212,'2025-11-29 20:09:15.428','',0),(129,1,'devops','CONFIG_ACCOUNT',1,'SELECT * FROM config_account LIMIT 20 OFFSET 0','','',1,'执行成功',241,'2025-11-29 20:09:18.899','',0),(130,1,'devops','SYS_MENU',1,'SELECT * FROM sys_menu LIMIT 20 OFFSET 0','','',1,'执行成功',211,'2025-11-29 20:09:31.833','',0),(131,1,'devops','CONFIG_KEYMANAGE',1,'SELECT * FROM config_keymanage LIMIT 20 OFFSET 0','','',1,'执行成功',218,'2025-11-29 20:10:16.736','',0),(132,1,'gin-api','APP_SERVICE_RELEASE',1,'SELECT * FROM app_service_release LIMIT 20 OFFSET 0','','',1,'执行成功',229,'2025-11-29 20:10:41.975','',0),(133,1,'gin-api','APP_JENKINS_ENV',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',321,'2025-11-29 20:11:33.117','',0),(134,1,'devops','APP_JENKINS_ENV',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',213,'2025-11-29 20:14:46.781','',0),(135,1,'devops','APP_JENKINS_ENV',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',213,'2025-11-29 20:17:07.095','',0),(136,1,'devops','CMDB_HOST',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',211,'2025-11-29 20:17:20.352','',0),(137,1,'devops','CMDB_GROUP',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',240,'2025-11-29 20:27:29.297','',0),(138,1,'devops','CMDB_SQL',1,'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0','','',1,'执行成功',258,'2025-11-29 20:28:55.241','',0),(139,1,'devops','CMDB_SQL',1,'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0','','',1,'执行成功',579,'2025-11-29 20:28:55.466','',0),(140,1,'devops','CMDB_GROUP',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',205,'2025-11-29 20:30:47.480','',0),(141,1,'devops','CMDB_HOST',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',199,'2025-11-29 20:30:58.495','',0),(142,1,'devops','CMDB_GROUP',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',267,'2025-11-29 20:34:21.225','',0),(143,1,'devops','CMDB_GROUP',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',4246,'2025-11-29 20:45:39.121','',0),(144,1,'devops','CMDB_GROUP',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',582,'2025-11-29 20:49:42.415','',0),(145,1,'devops','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',212,'2025-11-29 20:59:18.272','',0),(146,1,'devops','APP_APPLICATION',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',265,'2025-11-29 21:02:11.486','',0),(147,1,'devops','APP_JENKINS_ENV',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',242,'2025-11-29 21:05:56.924','',0),(148,1,'devops','APP_JENKINS_ENV',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',269,'2025-11-29 21:06:38.995','',0),(149,1,'devops','APP_JENKINS_ENV',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',327,'2025-11-29 21:08:52.177','',0),(150,1,'devops','APP_JENKINS_ENV',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',646,'2025-11-29 21:10:07.252','',0),(151,1,'devops','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 100;','','',1,'执行成功',293,'2025-11-29 21:18:18.161','',0),(152,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',320,'2025-11-29 21:21:01.619','',0),(153,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',297,'2025-11-29 21:21:10.168','',0),(154,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',326,'2025-11-29 21:21:27.886','',0),(155,1,'devops','cmdb_sql',1,'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0','','',1,'执行成功',283,'2025-11-29 21:22:09.300','',0),(156,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',283,'2025-11-29 21:22:19.066','',0),(157,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',276,'2025-11-29 21:22:21.160','',0),(158,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',264,'2025-11-29 21:34:45.094','',0),(159,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',1520,'2025-11-29 21:35:20.063','',0),(160,1,'sys','sys_config',1,'SELECT * FROM sys_config LIMIT 20 OFFSET 0','','',1,'执行成功',293,'2025-11-29 21:36:51.862','',0),(161,1,'gin-api','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',265,'2025-11-29 21:37:07.653','',0),(162,1,'gin-api','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',292,'2025-11-29 21:37:26.201','',0),(163,1,'devops','cmdb_sql_log',1,'SELECT * FROM cmdb_sql_log LIMIT 20 OFFSET 0','','',1,'执行成功',298,'2025-11-29 21:37:48.390','',0),(164,1,'devops','cmdb_sql_log',1,'SELECT * FROM cmdb_sql_log LIMIT 20 OFFSET 0','','',1,'执行成功',277,'2025-11-29 21:38:13.100','',0),(165,1,'devops','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',282,'2025-11-29 21:40:52.663','',0),(166,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',283,'2025-11-29 21:42:31.577','',0),(167,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',283,'2025-11-29 21:45:10.707','',0),(168,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',285,'2025-11-29 21:45:23.126','',0),(169,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',270,'2025-11-29 21:45:42.521','',0),(170,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',491,'2025-11-29 21:46:03.897','',0),(171,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',295,'2025-11-29 21:50:35.944','',0),(172,1,'devops','cmdb_sql',1,'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0','','',1,'执行成功',270,'2025-11-29 21:50:58.532','',0),(173,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',258,'2025-11-29 21:51:00.879','',0),(174,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',266,'2025-11-29 21:54:42.129','',0),(175,1,'devops','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',299,'2025-11-29 21:57:35.040','',0),(176,1,'devops','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',408,'2025-11-29 21:57:38.152','',0),(177,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',279,'2025-11-29 21:57:43.563','',0),(178,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',341,'2025-11-29 21:57:53.234','',0),(179,1,'devops','config_ecsauth',1,'SELECT * FROM config_ecsauth LIMIT 20 OFFSET 0','','',1,'执行成功',309,'2025-11-29 21:58:20.086','',0),(180,1,'devops','config_ecsauth',1,'SELECT * FROM config_ecsauth LIMIT 20 OFFSET 0','','',1,'执行成功',579,'2025-11-29 21:58:20.467','',0),(181,1,'devops','config_ecsauth',1,'SELECT * FROM config_ecsauth LIMIT 20 OFFSET 0','','',1,'执行成功',299,'2025-11-29 21:58:28.291','',0),(182,1,'devops','monitor_agent',1,'SELECT * FROM monitor_agent LIMIT 20 OFFSET 0','','',1,'执行成功',277,'2025-11-29 21:58:36.381','',0),(183,1,'devops','cmdb_sql',1,'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0','','',1,'执行成功',330,'2025-11-29 21:59:23.190','',0),(184,1,'devops','cmdb_sql',1,'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0','','',1,'执行成功',450,'2025-11-29 21:59:44.784','',0),(185,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',295,'2025-11-29 22:10:51.726','',0),(186,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',266,'2025-11-29 22:10:54.896','',0),(187,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',323,'2025-11-29 22:11:01.210','',0),(188,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',288,'2025-11-29 22:11:01.374','',0),(189,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',319,'2025-11-29 22:14:15.745','',0),(190,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',508,'2025-11-29 22:14:16.040','',0),(191,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',261,'2025-11-29 22:14:22.801','',0),(192,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',278,'2025-11-29 22:14:29.071','',0),(193,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',587,'2025-11-29 22:17:49.851','',0),(194,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',273,'2025-11-29 22:21:41.845','',0),(195,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',297,'2025-11-29 22:21:41.996','',0),(196,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',338,'2025-11-29 22:21:41.638','',0),(197,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',935,'2025-11-29 22:21:42.359','',0),(198,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',280,'2025-11-29 22:21:45.552','',0),(199,1,'devops','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',290,'2025-11-29 22:25:38.566','',0),(200,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',262,'2025-11-29 22:26:39.755','',0),(201,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',299,'2025-11-29 22:29:18.342','',0),(202,1,'devops','cmdb_sql',1,'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0','','',1,'执行成功',288,'2025-11-29 22:30:54.805','',0),(203,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',282,'2025-11-29 22:31:17.619','',0),(204,1,'gin-api','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',289,'2025-11-29 22:32:34.668','',0),(205,1,'devops','cmdb_sql_log',1,'SELECT * FROM cmdb_sql_log LIMIT 20 OFFSET 0','','',1,'执行成功',272,'2025-11-29 22:32:43.903','',0),(206,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',312,'2025-11-29 22:32:59.875','',0),(207,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',268,'2025-11-29 22:38:19.169','',0),(208,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',285,'2025-11-29 22:42:46.982','',0),(209,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',283,'2025-11-29 22:48:01.180','',0),(210,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',279,'2025-11-29 22:53:17.671','',0),(211,1,'devops','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',271,'2025-11-29 22:57:40.841','',0),(212,1,'devops','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',262,'2025-11-29 22:59:48.099','',0),(213,1,'devops','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',278,'2025-11-29 23:00:34.141','',0),(214,1,'devops','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',569,'2025-11-29 23:04:04.638','',0),(215,1,'devops','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',542,'2025-11-29 23:10:59.541','',0),(216,1,'gin-api','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',309,'2025-11-29 23:13:35.399','',0),(217,1,'gin-api','app_sh_release',1,'SELECT * FROM app_sh_release LIMIT 20 OFFSET 0','','',1,'执行成功',316,'2025-11-29 23:17:36.816','',0),(218,1,'gin-api','app_sh_release',1,'SELECT * FROM app_sh_release LIMIT 20 OFFSET 0','','',1,'执行成功',713,'2025-11-29 23:17:37.046','',0),(219,1,'gin-api','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',311,'2025-11-29 23:17:41.461','',0),(220,1,'gin-api','config_ecsauth',1,'SELECT * FROM config_ecsauth LIMIT 20 OFFSET 0','','',1,'执行成功',451,'2025-11-29 23:18:01.612','',0),(221,1,'gin-api','config_keymanage',1,'SELECT * FROM config_keymanage LIMIT 20 OFFSET 0','','',1,'执行成功',300,'2025-11-29 23:18:03.029','',0),(222,1,'gin-api','config_sync_schedule',1,'SELECT * FROM config_sync_schedule LIMIT 20 OFFSET 0','','',1,'执行成功',305,'2025-11-29 23:18:04.708','',0),(223,1,'gin-api','sys_operation_log',1,'SELECT * FROM sys_operation_log LIMIT 20 OFFSET 0','','',1,'执行成功',303,'2025-11-29 23:18:07.975','',0),(224,1,'gin-api','sys_operation_log',1,'SELECT * FROM sys_operation_log LIMIT 20 OFFSET 0','','',1,'执行成功',334,'2025-11-29 23:18:13.728','',0),(225,1,'gin-api','app_service_release',1,'SELECT * FROM app_service_release LIMIT 20 OFFSET 0','','',1,'执行成功',304,'2025-11-29 23:26:31.908','',0),(226,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',420,'2025-11-29 23:30:01.262','',0),(227,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',331,'2025-11-29 23:35:23.442','',0),(228,1,'devops','cmdb_sql_log',1,'SELECT * FROM cmdb_sql_log LIMIT 20 OFFSET 0','','',1,'执行成功',306,'2025-11-29 23:36:11.800','',0),(229,1,'devops','cmdb_sql_log',1,'SELECT * FROM cmdb_sql_log LIMIT 20 OFFSET 0','','',1,'执行成功',293,'2025-11-29 23:36:35.095','',0),(230,1,'gin-api','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',297,'2025-11-29 23:39:22.251','',0),(231,1,'gin-api','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',340,'2025-11-29 23:39:24.638','',0),(232,1,'devops','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',289,'2025-11-29 23:39:32.922','',0),(233,1,'devops','',5,'RENAME TABLE app_jenkins_env_copy_20251129233940 TO app_jenkins_env_copy_20251129233940_new111','','',1,'执行成功',298,'2025-11-29 23:39:51.148','',0),(234,1,'devops','app_jenkins_env_copy_20251129233940_new111',1,'SELECT * FROM app_jenkins_env_copy_20251129233940_new111 LIMIT 20 OFFSET 0','','',1,'执行成功',563,'2025-11-29 23:40:20.699','',0),(235,1,'devops','',5,'DROP TABLE app_jenkins_env_copy_20251129233940_new111','','',1,'执行成功',291,'2025-11-29 23:40:30.953','',0),(236,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',434,'2025-11-29 23:44:31.265','',0),(237,1,'devops','cmdb_sql',1,'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0','','',1,'执行成功',466,'2025-11-29 23:45:21.575','',0),(238,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',288,'2025-11-29 23:52:27.278','',0),(239,1,'devops','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',302,'2025-11-29 23:55:28.566','',0),(240,1,'gin-api','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',319,'2025-11-29 23:58:06.595','',0),(241,2,'postgres','',5,'CREATE DATABASE database_name;','','',1,'执行成功',599,'2025-11-30 00:56:01.163','',0),(242,2,'database_name','',5,'CREATE TABLE users (\n    id SERIAL PRIMARY KEY,\n    name VARCHAR(100) NOT NULL,\n    email VARCHAR(150) UNIQUE NOT NULL\n);','','',1,'执行成功',3259,'2025-11-30 01:10:07.720','',0),(243,2,'database_name','users',1,'SELECT * FROM users LIMIT 20 OFFSET 0','','',1,'执行成功',4965,'2025-11-30 01:10:34.369','',0),(244,2,'database_name','users',1,'SELECT * FROM users LIMIT 20 OFFSET 0','','',1,'执行成功',198,'2025-11-30 01:10:34.867','',0),(245,2,'database_name','users',1,'SELECT * FROM users LIMIT 20 OFFSET 0','','',1,'执行成功',201,'2025-11-30 01:10:35.037','',0),(246,2,'database_name','users',1,'SELECT * FROM users LIMIT 20 OFFSET 0','','',1,'执行成功',208,'2025-11-30 01:10:35.181','',0),(247,2,'database_name','users',1,'SELECT * FROM users LIMIT 20 OFFSET 0','','',1,'执行成功',208,'2025-11-30 01:10:35.442','',0),(248,2,'database_name','users',1,'SELECT * FROM users LIMIT 20 OFFSET 0','','',1,'执行成功',192,'2025-11-30 01:10:35.755','',0),(249,2,'database_name','users',1,'SELECT * FROM users LIMIT 20 OFFSET 0','','',1,'执行成功',205,'2025-11-30 01:10:35.988','',0),(250,2,'database_name','users',1,'SELECT * FROM users LIMIT 20 OFFSET 0','','',1,'执行成功',202,'2025-11-30 01:10:39.921','',0),(251,2,'database_name','users',2,'INSERT INTO users (name, email) VALUES (\'张三\', \'zhangsan@123.com\')','','',1,'执行成功',376,'2025-11-30 01:11:20.911','',0),(252,2,'database_name','users',1,'SELECT * FROM users LIMIT 20 OFFSET 0','','',1,'执行成功',299,'2025-11-30 01:11:21.506','',0),(253,1,'gin-api','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',299,'2025-11-30 11:38:28.572','',0),(254,1,'gin-api','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',275,'2025-11-30 11:38:44.346','',0),(255,1,'gin-api','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',272,'2025-11-30 11:38:53.975','',0),(256,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',342,'2025-11-30 11:39:11.611','',0),(257,1,'devops','cmdb_sql',1,'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0','','',1,'执行成功',277,'2025-11-30 11:39:53.115','',0),(258,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',284,'2025-11-30 11:40:04.010','',0),(259,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',273,'2025-11-30 14:35:46.125','',0),(260,2,'database_name','users',1,'SELECT * FROM users LIMIT 20 OFFSET 0','','',1,'执行成功',219,'2025-11-30 15:43:25.630','',0),(261,1,'RECOVER_YOUR_DATA','',5,'CREATE DATABASE IF NOT EXISTS test \n  DEFAULT CHARACTER SET utf8mb4;','','',1,'执行成功',762,'2025-11-30 23:49:38.276','',0),(262,1,'devops','cmdb_sql',1,'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0','','',1,'执行成功',345,'2025-11-30 23:50:09.029','',0),(263,1,'devops','cmdb_sql_log',1,'SELECT * FROM cmdb_sql_log LIMIT 20 OFFSET 0','','',1,'执行成功',304,'2025-11-30 23:50:24.695','',0),(264,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',311,'2025-11-30 23:50:29.549','',0),(265,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',285,'2025-11-30 23:51:51.702','',0),(266,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',308,'2025-11-30 23:53:56.382','',0),(267,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',291,'2025-11-30 23:54:23.552','',0),(268,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',292,'2025-11-30 23:54:33.359','',0),(269,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',303,'2025-11-30 23:54:38.550','',0),(270,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',333,'2025-11-30 23:54:46.965','',0),(271,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',319,'2025-11-30 23:54:58.898','',0),(272,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',314,'2025-11-30 23:54:59.768','',0),(273,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',315,'2025-11-30 23:55:02.322','',0),(274,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',290,'2025-11-30 23:55:11.683','',0),(275,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',282,'2025-11-30 23:55:36.034','',0),(276,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',261,'2025-11-30 23:55:40.374','',0),(277,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 100;','','',1,'执行成功',288,'2025-11-30 23:56:00.161','',0),(278,1,'devops','cmdb_group;',1,'SELECT * FROM cmdb_group;','','',1,'执行成功',272,'2025-11-30 23:56:57.273','',0),(279,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',279,'2025-11-30 23:57:18.128','',0),(280,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',286,'2025-11-30 23:57:32.189','',0),(281,1,'devops','cmdb_group',2,'INSERT INTO cmdb_group (parent_id, name, create_time, remark, update_time) VALUES (\'0\', \'test\', \'\', \'123\', \'\')','','',2,'Error 1292 (22007): Incorrect datetime value: \'\' for column \'create_time\' at row 1',234,'2025-11-30 23:58:10.839','',0),(282,1,'devops','cmdb_host_copy_20251130235846',1,'SELECT * FROM cmdb_host_copy_20251130235846 LIMIT 20 OFFSET 0','','',1,'执行成功',286,'2025-11-30 23:58:49.791','',0),(283,1,'devops','',5,'DROP TABLE cmdb_host_copy_20251130235846','','',1,'执行成功',295,'2025-11-30 23:58:54.686','',0),(284,1,'devops','',5,'RENAME TABLE cmdb_host_copy_20251130235914 TO cmdb_host_copy_123','','',1,'执行成功',274,'2025-11-30 23:59:22.577','',0),(285,2,'database_name','users',1,'SELECT * FROM users LIMIT 20 OFFSET 0','','',1,'执行成功',207,'2025-11-30 23:59:58.611','',0),(286,2,'database_name','users',2,'INSERT INTO users (name, email) VALUES (\'李四\', \'lisi@123.com\')','','',1,'执行成功',285,'2025-12-01 00:00:30.571','',0),(287,2,'database_name','users',1,'SELECT * FROM users LIMIT 20 OFFSET 0','','',1,'执行成功',214,'2025-12-01 00:00:31.076','',0),(288,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',280,'2025-12-01 00:05:29.288','',0),(289,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',295,'2025-12-01 00:06:23.704','',0),(290,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',1184,'2025-12-01 00:06:54.903','',0),(291,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',611,'2025-12-01 00:15:50.948','',0),(292,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',291,'2025-12-01 00:17:04.651','',0),(293,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',303,'2025-12-01 00:19:18.036','',0),(294,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',417,'2025-12-01 00:19:22.279','',0),(295,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',280,'2025-12-01 00:20:01.970','',0),(296,1,'devops','cmdb_host',3,'UPDATE cmdb_host SET host_name = \'虚拟机001111\', group_id = 12, private_ip = \'172.16.226.16\', public_ip = \'120.231.244.158\', ssh_name = \'root\', ssh_key_id = 13, ssh_port = 22, remark = \'123\', vendor = 1, region = \'\', instance_id = \'\', os = \'CentOSLinux7(Core)\', status = 1, cpu = \'2\', memory = \'4\', disk = \'17\', billing_type = \'\', create_time = \'2025-11-23 23:44:35\', expire_time = NULL, update_time = \'2025-11-23 23:45:32\', ssh_ip = \'172.16.226.16\', name = \'jenkins\', ssh_gateway_id = NULL WHERE id = 511;','','',1,'执行成功',254,'2025-12-01 00:20:08.106','',0),(297,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',292,'2025-12-01 00:20:08.651','',0),(298,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',301,'2025-12-01 00:20:11.764','',0),(299,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 100;','','',1,'执行成功',293,'2025-12-01 00:20:31.016','',0),(300,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',302,'2025-12-01 00:21:51.502','',0),(301,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 100;','','',1,'执行成功',569,'2025-12-01 00:21:57.257','',0),(302,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',367,'2025-12-01 00:23:39.693','',0),(303,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 100;','','',1,'执行成功',540,'2025-12-01 00:23:47.277','',0),(304,2,'database_name','users',1,'SELECT * FROM users LIMIT 20 OFFSET 0','','',1,'执行成功',217,'2025-12-01 00:24:18.173','',0),(305,2,'database_name','users',1,'SELECT * FROM users LIMIT 20 OFFSET 0','','',1,'执行成功',289,'2025-12-01 00:25:05.373','',0),(306,2,'database_name','users',2,'INSERT INTO users (name, email) VALUES (\'test\', \'123@123.com\')','','',1,'执行成功',553,'2025-12-01 00:30:03.796','',0),(307,2,'database_name','users',1,'SELECT * FROM users LIMIT 20 OFFSET 0','','',1,'执行成功',205,'2025-12-01 00:30:04.267','',0),(308,2,'database_name','users',1,'SELECT * FROM users LIMIT 20 OFFSET 0','','',1,'执行成功',256,'2025-12-01 00:30:24.140','',0),(309,2,'database_name','users',1,'SELECT * FROM users LIMIT 20 OFFSET 0','','',1,'执行成功',452,'2025-12-01 00:37:20.181','',0),(310,2,'database_name','users',1,'SELECT * FROM users LIMIT 20 OFFSET 0','','',1,'执行成功',235,'2025-12-01 00:37:27.599','',0),(311,2,'database_name','users_copy_20251201003823',1,'SELECT * FROM users_copy_20251201003823 LIMIT 20 OFFSET 0','','',1,'执行成功',197,'2025-12-01 00:38:47.198','',0),(312,2,'database_name','users_copy_20251201003823',1,'SELECT * FROM users_copy_20251201003823 LIMIT 100;','','',1,'执行成功',227,'2025-12-01 00:38:56.092','',0),(313,1,'RECOVER_YOUR_DATA','',5,'CREATE DATABASE IF NOT EXISTS ops \n  DEFAULT CHARACTER SET utf8mb4;','','',1,'执行成功',729,'2025-12-01 00:43:24.905','',0),(314,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',296,'2025-12-01 00:43:58.889','',0),(315,1,'devops','cmdb_host_copy_123',1,'SELECT * FROM cmdb_host_copy_123 LIMIT 20 OFFSET 0','','',1,'执行成功',351,'2025-12-01 00:44:00.240','',0),(316,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',320,'2025-12-01 00:44:02.951','',0),(317,1,'devops','cmdb_host',3,'UPDATE cmdb_host SET host_name = \'虚拟机01\', group_id = 12, private_ip = \'172.16.226.16\', public_ip = \'120.231.244.158\', ssh_name = \'root\', ssh_key_id = 13, ssh_port = 22, remark = \'123\', vendor = 1, region = \'\', instance_id = \'\', os = \'CentOSLinux7(Core)\', status = 1, cpu = \'2\', memory = \'4\', disk = \'17\', billing_type = \'\', create_time = \'2025-11-23 23:44:35\', expire_time = NULL, update_time = \'2025-11-23 23:45:32\', ssh_ip = \'172.16.226.16\', name = \'jenkins\', ssh_gateway_id = NULL WHERE id = 511;','','',1,'执行成功',722,'2025-12-01 00:44:23.250','',0),(318,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',333,'2025-12-01 00:44:26.196','',0),(319,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',297,'2025-12-01 00:44:31.613','',0),(320,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',316,'2025-12-01 00:44:46.505','',0),(321,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',475,'2025-12-01 00:46:30.077','',0),(322,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',326,'2025-12-01 00:46:48.843','',0),(323,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',290,'2025-12-01 00:46:53.212','',0),(324,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',296,'2025-12-01 00:46:57.044','',0),(325,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',306,'2025-12-01 00:47:00.548','',0),(326,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',296,'2025-12-01 00:47:09.235','',0),(327,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',327,'2025-12-01 00:47:22.518','',0),(328,1,'devops','cmdb_host;',1,'SELECT * FROM  cmdb_host;','','',1,'执行成功',297,'2025-12-01 00:48:12.389','',0),(329,1,'devops','',5,'DROP TABLE cmdb_host_copy_123','','',1,'执行成功',578,'2025-12-01 00:49:01.428','',0),(330,1,'devops','cmdb_host_copy_20251201004911',1,'SELECT * FROM cmdb_host_copy_20251201004911 LIMIT 20 OFFSET 0','','',1,'执行成功',321,'2025-12-01 00:49:15.842','',0),(331,1,'devops','',5,'RENAME TABLE cmdb_host_copy_20251201004911 TO cmdb_host_copy_123','','',1,'执行成功',253,'2025-12-01 00:49:25.590','',0),(332,1,'devops','',5,'DROP TABLE cmdb_host_copy_123','','',1,'执行成功',259,'2025-12-01 00:49:43.125','',0),(333,2,'database_name','',5,'CREATE TABLE students (\n    id SERIAL PRIMARY KEY,                -- 学生ID，自增主键\n    student_id VARCHAR(20) UNIQUE NOT NULL, -- 学号，唯一且非空\n    name VARCHAR(50) NOT NULL,            -- 姓名\n    gender CHAR(1) CHECK (gender IN (\'M\', \'F\')), -- 性别：M 男，F 女\n    birth_date DATE,                      -- 出生日期\n    email VARCHAR(100) UNIQUE,            -- 邮箱，唯一\n    phone VARCHAR(20),                    -- 电话\n    enrollment_date DATE DEFAULT CURRENT_DATE, -- 入学日期，默认为当前日期\n    major VARCHAR(100),                   -- 专业\n    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- 记录创建时间\n);','','',1,'执行成功',239,'2025-12-01 00:50:33.682','',0),(334,2,'database_name','students',1,'SELECT * FROM students LIMIT 20 OFFSET 0','','',1,'执行成功',222,'2025-12-01 00:50:38.712','',0),(335,2,'database_name','students',2,'INSERT INTO students (student_id, name, gender, birth_date, email, phone, enrollment_date, major) VALUES (\'1\', \'张三\', \'m\', \'2004-05-15\', \'123@456.com\', \'12345678911\', \'2004-05-15\', \'软件工程\')','','',2,'pq: new row for relation \"students\" violates check constraint \"students_gender_check\"',186,'2025-12-01 00:52:04.216','',0),(336,2,'database_name','students',2,'INSERT INTO students (student_id, name, gender, birth_date, email, phone, enrollment_date, major) VALUES (\'20230001\', \'张三\', \'m\', \'2004-05-15\', \'123@456.com\', \'13800138001\', \'2004-05-15\', \'软件工程\')','','',2,'pq: new row for relation \"students\" violates check constraint \"students_gender_check\"',182,'2025-12-01 00:52:31.225','',0),(337,2,'database_name','students',2,'INSERT INTO students (student_id, name, gender, birth_date, email, phone, enrollment_date, major) VALUES (\'20230001\', \'张三\', \'\'\'M\'\'\', \'2004-05-15\', \'123@456.com\', \'13800138001\', \'2004-05-15\', \'软件工程\')','','',2,'pq: value too long for type character(1)',214,'2025-12-01 00:52:54.401','',0),(338,2,'database_name','users',1,'SELECT * FROM users LIMIT 20 OFFSET 0','','',1,'执行成功',542,'2025-12-01 00:52:58.529','',0),(339,2,'database_name','users',2,'INSERT INTO users (name, email) VALUES (\'test123\', \'123@123.com\')','','',1,'执行成功',183,'2025-12-01 00:53:22.992','',0),(340,2,'database_name','users',1,'SELECT * FROM users LIMIT 20 OFFSET 0','','',1,'执行成功',253,'2025-12-01 00:53:23.494','',0),(341,2,'database_name','users',1,'SELECT * FROM users LIMIT 20 OFFSET 0','','',1,'执行成功',207,'2025-12-01 00:53:33.443','',0),(342,2,'database_name','users',1,'SELECT * FROM users LIMIT 100;','','',1,'执行成功',196,'2025-12-01 00:53:54.378','',0),(343,2,'database_name','users_copy_20251201003823',1,'SELECT * FROM users_copy_20251201003823 LIMIT 20 OFFSET 0','','',1,'执行成功',231,'2025-12-01 00:53:59.552','',0),(344,2,'database_name','',5,'DROP TABLE users_copy_20251201003823','','',1,'执行成功',184,'2025-12-01 00:54:05.164','',0),(345,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',286,'2025-12-01 10:18:58.108','',0),(346,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 100;','','',1,'执行成功',257,'2025-12-01 10:19:34.501','',0),(347,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',265,'2025-12-01 12:05:45.106','',0),(348,1,'RECOVER_YOUR_DATA','RECOVER_YOUR_DATA',1,'SELECT * FROM RECOVER_YOUR_DATA LIMIT 20 OFFSET 0','','',1,'执行成功',288,'2025-12-01 14:07:15.387','',0),(349,1,'devops','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',256,'2025-12-01 14:07:23.511','',0),(350,1,'gin-api','app_sh_release',1,'SELECT * FROM app_sh_release LIMIT 20 OFFSET 0','','',1,'执行成功',254,'2025-12-01 17:10:51.225','',0),(351,1,'gin-api','app_service_release',1,'SELECT * FROM app_service_release LIMIT 20 OFFSET 0','','',1,'执行成功',245,'2025-12-01 17:11:07.319','',0),(352,1,'gin-api','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',253,'2025-12-01 17:11:39.058','',0),(353,1,'gin-api','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',248,'2025-12-01 17:11:57.884','',0),(354,1,'gin-api','quick_deployment_tasks',1,'SELECT * FROM quick_deployment_tasks LIMIT 20 OFFSET 0','','',1,'执行成功',280,'2025-12-01 17:12:31.887','',0),(355,1,'gin-api','quick_deployments',1,'SELECT * FROM quick_deployments LIMIT 20 OFFSET 0','','',1,'执行成功',254,'2025-12-01 17:12:34.351','',0),(356,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',248,'2025-12-01 19:50:46.392','',0),(357,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',294,'2025-12-02 11:13:02.970','',0),(358,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',265,'2025-12-02 11:34:17.783','',0),(359,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',256,'2025-12-02 11:35:41.192','',0),(360,1,'gin-api','sys_operation_log',1,'SELECT * FROM sys_operation_log LIMIT 20 OFFSET 0','','',1,'执行成功',292,'2025-12-02 12:01:05.583','',0),(361,1,'gin-api','sys_operation_log',1,'SELECT * FROM sys_operation_log LIMIT 100;','','',1,'执行成功',555,'2025-12-02 12:06:45.939','',0),(362,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',257,'2025-12-02 12:08:01.959','',0),(363,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',250,'2025-12-02 12:25:37.979','',0),(364,1,'devops','app_application',1,'SELECT *\nFROM  app_application\nLIMIT  100;','','',1,'执行成功',252,'2025-12-02 12:25:58.932','',0),(365,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',277,'2025-12-02 12:30:05.820','',0),(366,2,'database_name','users',1,'SELECT * FROM users LIMIT 20 OFFSET 0','','',1,'执行成功',253,'2025-12-02 12:35:09.666','',0),(367,2,'database_name','students',1,'SELECT * FROM students LIMIT 20 OFFSET 0','','',1,'执行成功',180,'2025-12-02 12:35:11.544','',0),(368,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',262,'2025-12-02 12:38:20.200','',0),(369,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',251,'2025-12-02 12:38:39.691','',0),(370,1,'gin-api','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',250,'2025-12-02 14:59:49.691','',0),(371,1,'gin-api','db',1,'SELECT * FROM db LIMIT 20 OFFSET 0','','',1,'执行成功',261,'2025-12-02 15:00:14.754','',0),(372,1,'gin-api','db_instance',1,'SELECT * FROM db_instance LIMIT 20 OFFSET 0','','',1,'执行成功',258,'2025-12-02 15:01:14.757','',0),(373,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',255,'2025-12-02 15:09:10.006','',0),(374,1,'gin-api','cmdb_sql',1,'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0','','',1,'执行成功',247,'2025-12-02 15:41:50.960','',0),(375,1,'gin-api','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',251,'2025-12-02 15:42:08.133','',0),(376,1,'gin-api','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',252,'2025-12-02 15:42:47.779','',0),(377,1,'RECOVER_YOUR_DATA','RECOVER_YOUR_DATA',1,'SELECT * FROM RECOVER_YOUR_DATA LIMIT 20 OFFSET 0','','',1,'执行成功',284,'2025-12-02 16:00:08.047','',0),(378,1,'mayfly-go','t_db_backup_history',1,'SELECT * FROM t_db_backup_history LIMIT 20 OFFSET 0','','',1,'执行成功',284,'2025-12-02 16:00:14.480','',0),(379,1,'mayfly-go','t_db_restore',1,'SELECT * FROM t_db_restore LIMIT 20 OFFSET 0','','',1,'执行成功',283,'2025-12-02 16:00:16.976','',0),(380,3,'gin-api','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',256,'2025-12-02 16:59:07.917','',0),(381,3,'gin-api','app_service_release_item',1,'SELECT * FROM app_service_release_item LIMIT 20 OFFSET 0','','',1,'执行成功',261,'2025-12-02 16:59:08.689','',0),(382,1,'gin-api','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',265,'2025-12-03 10:26:41.230','',0),(383,1,'gin-api','app_service_release_item',1,'SELECT * FROM app_service_release_item LIMIT 20 OFFSET 0','','',1,'执行成功',323,'2025-12-03 10:27:20.235','',0),(384,1,'gin-api','app_sh_release',1,'SELECT * FROM app_sh_release LIMIT 20 OFFSET 0','','',1,'执行成功',256,'2025-12-03 10:27:21.025','',0),(385,1,'gin-api','cmdb_sql_log',1,'SELECT * FROM cmdb_sql_log LIMIT 20 OFFSET 0','','',1,'执行成功',548,'2025-12-03 10:31:37.487','',0),(386,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',263,'2025-12-03 10:36:57.065','',0),(387,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',340,'2025-12-03 10:37:02.051','',0),(388,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',266,'2025-12-03 10:40:18.048','',0),(389,1,'devops','cmdb_host_1',2,'INSERT INTO cmdb_host_1 (id, host_name, group_id, private_ip, public_ip, ssh_name, ssh_key_id, ssh_port, remark, vendor, region, instance_id, os, status, cpu, memory, disk, billing_type, create_time, expire_time, update_time, ssh_ip, name, ssh_gateway_id) VALUES (506, \'华为云ops\', 4, \'172.31.6.35\', \'139.9.205.38\', \'root\', 22, 22, \'123\', 5, \'\', \'\', \'Ubuntu24.04.2\', 1, \'2\', \'2\', \'40\', \'\', \'2025-11-11 17:24:36\', NULL, \'2025-11-26 17:22:07\', \'139.9.205.38\', \'hw-ops\', NULL);','','',2,'Error 1146 (42S02): Table \'devops.cmdb_host_1\' doesn\'t exist',202,'2025-12-03 10:43:12.842','',0),(390,1,'devops','cmdb_host_1',2,'INSERT INTO cmdb_host_1 (id, host_name, group_id, private_ip, public_ip, ssh_name, ssh_key_id, ssh_port, remark, vendor, region, instance_id, os, status, cpu, memory, disk, billing_type, create_time, expire_time, update_time, ssh_ip, name, ssh_gateway_id) VALUES (506, \'华为云ops\', 4, \'172.31.6.35\', \'139.9.205.38\', \'root\', 22, 22, \'123\', 5, \'\', \'\', \'Ubuntu24.04.2\', 1, \'2\', \'2\', \'40\', \'\', \'2025-11-11 17:24:36\', NULL, \'2025-11-26 17:22:07\', \'139.9.205.38\', \'hw-ops\', NULL);','','',2,'Error 1146 (42S02): Table \'devops.cmdb_host_1\' doesn\'t exist',216,'2025-12-03 10:43:14.148','',0),(391,1,'devops','cmdb_sql',1,'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0','','',1,'执行成功',785,'2025-12-03 11:09:19.304','',0),(392,1,'devops','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',265,'2025-12-03 11:24:14.607','',0),(393,1,'gin-api','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',275,'2025-12-03 11:35:15.655','',0),(394,1,'mayfly-go','t_db',1,'SELECT * FROM t_db LIMIT 20 OFFSET 0','','',1,'执行成功',266,'2025-12-03 11:40:08.029','',0),(395,1,'mayfly-go','t_db_backup',1,'SELECT * FROM t_db_backup LIMIT 20 OFFSET 0','','',1,'执行成功',267,'2025-12-03 11:40:10.424','',0),(396,1,'mayfly-go','t_db_data_sync_log',1,'SELECT * FROM t_db_data_sync_log LIMIT 20 OFFSET 0','','',1,'执行成功',416,'2025-12-03 11:41:35.563','',0),(397,1,'mayfly-go','t_db_instance',1,'SELECT * FROM t_db_instance LIMIT 20 OFFSET 0','','',1,'执行成功',259,'2025-12-03 11:41:38.235','',0),(398,1,'gin-api','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',256,'2025-12-03 11:41:41.904','',0),(399,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',257,'2025-12-03 11:41:59.928','',0),(400,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',269,'2025-12-03 11:43:46.442','',0),(401,1,'devops','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',317,'2025-12-03 11:45:13.925','',0),(402,1,'devops','',5,'CREATE TABLE `app_jenkins_env123` ( `id` bigint(20) NOT NULL AUTO_INCREMENT, `app_id` bigint(20) NOT NULL, `env_name` varchar(50) NOT NULL, `jenkins_server_id` bigint(20), `job_name` varchar(255), `job_url` varchar(500), `build_params` json, `deploy_config` json, `notification` json, `is_active` tinyint(3) DEFAULT 1, `created_at` datetime, `updated_at` datetime, `deleted_at` datetime, PRIMARY KEY (id) ); ALTER TABLE `app_jenkins_env` ADD INDEX `idx_app_jenkins_env_app_id`(`app_id`) USING BTREE; ALTER TABLE `app_jenkins_env` ADD INDEX `idx_app_jenkins_env_deleted_at`(`deleted_at`) USING BTREE','','',2,'Error 1064 (42000): You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near \'ALTER TABLE `app_jenkins_env` ADD INDEX `idx_app_jenkins_env_app_id`(`app_id`) U\' at line 1',204,'2025-12-03 11:47:45.166','',0),(403,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',270,'2025-12-03 11:48:51.742','',0),(404,1,'devops','cmdb_sql',1,'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0','','',1,'执行成功',267,'2025-12-03 11:52:55.287','',0),(405,1,'devops','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',257,'2025-12-03 11:52:58.121','',0),(406,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',260,'2025-12-03 11:53:00.314','',0),(407,1,'gin-api','cmdb_sql',1,'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0','','',1,'执行成功',295,'2025-12-03 11:56:02.107','',0),(408,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',256,'2025-12-03 12:04:05.208','',0),(409,1,'devops','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',254,'2025-12-03 12:07:28.099','',0),(410,1,'gin-api','app_sh_release',1,'SELECT * FROM app_sh_release LIMIT 20 OFFSET 0','','',1,'执行成功',261,'2025-12-03 12:08:37.417','',0),(411,1,'devops','cmdb_sql',1,'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0','','',1,'执行成功',327,'2025-12-03 12:08:52.889','',0),(412,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',258,'2025-12-03 12:14:02.575','',0),(413,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',342,'2025-12-03 12:16:54.823','',0),(414,1,'gin-api','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',261,'2025-12-03 12:19:56.646','',0),(415,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',257,'2025-12-03 12:23:13.577','',0),(416,1,'devops','cmdb_sql',1,'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0','','',1,'执行成功',261,'2025-12-03 12:26:28.689','',0),(417,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',253,'2025-12-03 12:28:20.840','',0),(418,1,'gin-api','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',244,'2025-12-03 12:29:04.357','',0),(419,1,'devops','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',395,'2025-12-03 12:32:00.866','',0),(420,1,'devops','',5,'CREATE TABLE `cmdb_group123` (\n `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT \'\'\'主键\'\'\',\n `parent_id` bigint(20) NOT NULL DEFAULT 0 COMMENT \'\'\'父级分组ID\'\'\',\n `name` longtext NOT NULL COMMENT \'\'\'分组名称\'\'\',\n `create_time` datetime NOT NULL COMMENT \'\'\'创建时间\'\'\',\n `remark` longtext COMMENT \'\'\'备注\'\'\',\n `update_time` datetime COMMENT \'\'\'更新时间\'\'\', \nPRIMARY KEY (id)\n)','','',1,'执行成功',244,'2025-12-03 12:42:43.341','',0),(421,1,'devops','cmdb_group123',1,'SELECT * FROM cmdb_group123 LIMIT 20 OFFSET 0','','',1,'执行成功',262,'2025-12-03 12:42:51.081','',0),(422,1,'devops','',5,'DROP TABLE cmdb_group123','','',1,'执行成功',228,'2025-12-03 12:43:02.777','',0),(423,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',265,'2025-12-03 12:46:17.553','',0),(424,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',278,'2025-12-03 15:54:52.850','',0),(425,1,'devops','cmdb_sql_log',1,'SELECT * FROM cmdb_sql_log LIMIT 20 OFFSET 0','','',1,'执行成功',252,'2025-12-03 15:55:18.939','',0),(426,1,'gin-api','sys_operation_log',1,'SELECT * FROM sys_operation_log LIMIT 20 OFFSET 0','','',1,'执行成功',254,'2025-12-03 15:55:30.381','',0),(427,1,'gin-api','sys_operation_log',1,'SELECT *\nFROM  sys_operation_log\nLIMIT  100;','','',1,'执行成功',266,'2025-12-03 15:55:49.290','',0),(428,1,'gin-api','sys_admin',1,'SELECT *\nFROM  sys_admin\nLIMIT  100;','','',1,'执行成功',250,'2025-12-03 15:57:53.622','',0),(429,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',256,'2025-12-05 20:29:40.131','',0),(430,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',262,'2025-12-05 20:32:28.584','',0),(431,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',392,'2025-12-05 20:33:31.694','',0),(432,1,'devops','app_application',1,'SELECT *\nFROM  app_application\nLIMIT  100;','','',1,'执行成功',308,'2025-12-05 20:33:51.765','',0),(433,1,'devops','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',424,'2025-12-05 20:34:23.519','',0),(434,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',273,'2025-12-05 20:34:30.817','',0),(435,1,'devops','cmdb_host',3,'UPDATE cmdb_host SET host_name = \'虚拟机2025\', group_id = 4, private_ip = \'172.16.226.13\', public_ip = \'120.231.244.158\', ssh_name = \'root\', ssh_key_id = 13, ssh_port = 22, remark = \'123\', vendor = 1, region = \'\', instance_id = \'\', os = \'Ubuntu20.04.2\', status = 1, cpu = \'2\', memory = \'3\', disk = \'19\', billing_type = \'\', create_time = \'2025-11-23 23:46:03\', expire_time = NULL, update_time = \'2025-11-23 23:46:07\', ssh_ip = \'172.16.226.13\', name = \'k8s-node02\', ssh_gateway_id = NULL WHERE id = 512;','','',1,'执行成功',436,'2025-12-05 20:34:45.595','',0),(436,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',267,'2025-12-05 20:34:46.094','',0),(437,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',267,'2025-12-05 20:34:49.893','',0),(438,1,'devops','',5,'CREATE DATABASE IF NOT EXISTS test1\n  DEFAULT CHARACTER SET utf8mb4;','','',1,'执行成功',237,'2025-12-05 20:35:34.369','',0),(439,1,'test1','',5,'CREATE TABLE test1.your_table_name (\n  id BIGINT PRIMARY KEY AUTO_INCREMENT,\n  -- your columns here\n) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;','','',2,'Error 1064 (42000): You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near \') ENGINE=InnoDB DEFAULT CHARSET=utf8mb4\' at line 4',219,'2025-12-05 20:35:52.483','',0),(440,2,'database_name','users',1,'SELECT * FROM users LIMIT 20 OFFSET 0','','',1,'执行成功',217,'2025-12-05 20:36:57.655','',0),(441,2,'database_name','users',2,'INSERT INTO users (name, email) VALUES (\'王五\', \'xxxx@123.com\')','','',1,'执行成功',191,'2025-12-05 20:37:14.770','',0),(442,2,'database_name','users',1,'SELECT * FROM users LIMIT 20 OFFSET 0','','',1,'执行成功',205,'2025-12-05 20:37:15.207','',0),(443,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',267,'2025-12-08 10:52:11.033','',0),(444,1,'gin-api','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',288,'2025-12-08 23:01:38.354','',0),(445,1,'gin-api','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',267,'2025-12-08 23:01:38.638','',0),(446,1,'gin-api','sys_config',1,'SELECT * FROM sys_config LIMIT 20 OFFSET 0','','',1,'执行成功',259,'2025-12-09 13:21:51.379','',0),(447,1,'gin-api','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',280,'2025-12-09 23:45:28.566','',0),(448,1,'gin-api','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',268,'2025-12-09 23:48:06.048','',0),(449,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',267,'2025-12-10 11:22:14.627','',0),(450,1,'gin-api','sys_admin',1,'SELECT * FROM sys_admin LIMIT 20 OFFSET 0','','',1,'执行成功',268,'2025-12-10 17:17:50.042','',0),(451,1,'gin-api','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',276,'2025-12-11 18:04:01.397','',0),(452,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',314,'2025-12-13 00:03:02.341','',0),(453,2,'database_name','users',1,'SELECT * FROM users LIMIT 20 OFFSET 0','','',1,'执行成功',287,'2025-12-13 00:03:27.273','',0),(454,2,'database_name','students',1,'SELECT * FROM students LIMIT 20 OFFSET 0','','',1,'执行成功',257,'2025-12-13 00:03:29.263','',0),(455,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',303,'2025-12-13 16:15:43.928','',0),(456,1,'gin-api','monitor_aliyun_config',1,'SELECT * FROM monitor_aliyun_config LIMIT 20 OFFSET 0','','',1,'执行成功',278,'2025-12-15 15:44:32.136','',0),(457,1,'gin-api','',5,'ALTER TABLE `monitor_aliyun_config`\nADD COLUMN `eab_kid` VARCHAR(256) NULL COMMENT \'ACME EAB Key ID (ZeroSSL等需要)\' AFTER `username`,\nADD COLUMN `eab_hmac_key` VARCHAR(256) NULL COMMENT \'ACME EAB HMAC Key (ZeroSSL等需要)\' AFTER `eab_kid`;','','',2,'Error 1060 (42S21): Duplicate column name \'eab_kid\'',206,'2025-12-15 15:44:47.282','',0),(458,1,'gin-api','',5,'ALTER TABLE `monitor_aliyun_config` ADD COLUMN `eab_kid` VARCHAR(256) NULL COMMENT \'ACME EAB Key ID (ZeroSSL等需要)\' AFTER `username`, ADD COLUMN `eab_hmac_key` VARCHAR(256) NULL COMMENT \'ACME EAB HMAC Key (ZeroSSL等需要)\' AFTER `eab_kid`;','','',2,'Error 1060 (42S21): Duplicate column name \'eab_kid\'',222,'2025-12-15 15:45:18.947','',0),(459,1,'gin-api','monitor_aliyun_config',1,'SELECT * FROM monitor_aliyun_config LIMIT 20 OFFSET 0','','',1,'执行成功',271,'2025-12-15 15:45:25.985','',0),(460,1,'gin-api','monitor_aliyun_config',1,'SELECT * FROM monitor_aliyun_config LIMIT 20 OFFSET 0','','',1,'执行成功',271,'2025-12-15 15:48:58.000','',0),(461,1,'gin-api','monitor_ssl_cert',1,'SELECT * FROM monitor_ssl_cert LIMIT 20 OFFSET 0','','',1,'执行成功',313,'2025-12-17 11:33:48.528','',0),(462,1,'gin-api','monitor_ssl_cert_deploy_log',1,'SELECT * FROM monitor_ssl_cert_deploy_log LIMIT 20 OFFSET 0','','',1,'执行成功',257,'2025-12-17 12:07:06.756','',0),(463,1,'devops','app_service_release_item',1,'SELECT * FROM app_service_release_item LIMIT 20 OFFSET 0','','',1,'执行成功',278,'2025-12-22 11:33:29.093','',0),(464,1,'devops','cmdb_sql',1,'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0','','',1,'执行成功',261,'2025-12-24 11:21:35.417','',0),(465,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',268,'2025-12-24 11:21:45.959','',0),(466,1,'devops','db',1,'SELECT * FROM db LIMIT 20 OFFSET 0','','',1,'执行成功',255,'2025-12-24 11:21:52.729','',0),(467,1,'devops','',5,'DROP TABLE db','','',1,'执行成功',232,'2025-12-24 11:21:57.771','',0),(468,1,'devops','app_sh_release',1,'SELECT * FROM app_sh_release LIMIT 20 OFFSET 0','','',1,'执行成功',269,'2025-12-24 11:22:12.661','',0),(469,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',247,'2025-12-24 11:22:22.124','',0),(470,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',261,'2025-12-24 11:39:12.632','',0),(471,1,'devops','cmdb_group',1,'SELECT * FROM cmdb_group LIMIT 20 OFFSET 0','','',1,'执行成功',253,'2025-12-24 11:39:18.600','',0),(472,1,'devops','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',257,'2025-12-24 11:39:26.919','',0),(473,1,'devops','cmdb_sql',1,'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0','','',1,'执行成功',487,'2025-12-24 11:39:27.994','',0),(474,1,'devops','cmdb_sql',1,'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0','','',1,'执行成功',263,'2025-12-24 11:39:33.242','',0),(475,1,'devops','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',282,'2025-12-24 20:45:17.632','',0),(476,1,'gin-api','app_application',1,'SELECT * FROM app_application LIMIT 20 OFFSET 0','','',1,'执行成功',262,'2025-12-25 12:00:45.401','',0),(477,1,'gin-api','app_service_release',1,'SELECT * FROM app_service_release LIMIT 20 OFFSET 0','','',1,'执行成功',274,'2025-12-25 12:00:56.275','',0),(478,1,'gin-api','app_jenkins_env',1,'SELECT * FROM app_jenkins_env LIMIT 20 OFFSET 0','','',1,'执行成功',269,'2025-12-25 12:01:17.797','',0),(479,1,'gin-api','app_sh_release',1,'SELECT * FROM app_sh_release LIMIT 20 OFFSET 0','','',1,'执行成功',257,'2025-12-25 12:01:22.749','',0),(480,1,'gin-api','app_service_release_item',1,'SELECT * FROM app_service_release_item LIMIT 20 OFFSET 0','','',1,'执行成功',264,'2025-12-25 12:01:50.270','',0),(481,1,'gin-api','app_sh_release',5,'-- 上线趋势统计调试SQL -- 用于排查为什么最新的发布记录没有被统计到 -- =========================================== -- 1. 检查脚本上线记录 (近7天) -- =========================================== SELECT \'脚本上线-近7天\' AS \'统计类型\', DATE(created_at) as \'日期\', COUNT(*) as \'数量\'\nFROM  app_sh_release\nWHERE\n  created_at >= DATE_SUB(CURDATE(), INTERVAL 6 DAY)\n  AND created_at < DATE_ADD(CURDATE(), INTERVAL 1 DAY)\n  AND deleted_at IS NULL\n  AND status = 2\nGROUP BY  DATE(created_at)\nORDER BY  DATE(created_at) DESC; -- 查看所有脚本上线记录（包括不符合条件的） SELECT \'所有脚本上线记录\' AS \'类型\', id, title, status, CASE status\n    WHEN 1 THEN \'流程进行中\'\n    WHEN 2 THEN \'执行成功\'\n    WHEN 3 THEN \'执行失败\'\n    WHEN 4 THEN \'已驳回\'\n    WHEN 5 THEN \'已取消\'\n  END AS \'状态文本\', DATE(created_at) as \'日期\', created_at AS \'创建时间\', deleted_at AS \'删除时间\'\nFROM  app_sh_release\nWHERE  created_at >= DATE_SUB(CURDATE(), INTERVAL 6 DAY)\nORDER BY  created_at DESC\nLIMIT  10; -- =========================================== -- 2. 检查快速发布记录 (近7天) -- =========================================== SELECT \'快速发布-近7天\' AS \'统计类型\', DATE(created_at) as \'日期\', COUNT(*) as \'数量\'\nFROM  quick_deployments\nWHERE\n  created_at >= DATE_SUB(CURDATE(), INTERVAL 6 DAY)\n  AND created_at < DATE_ADD(CURDATE(), INTERVAL 1 DAY)\n  AND status = 3\nGROUP BY  DATE(created_at)\nORDER BY  DATE(created_at) DESC; -- 查看所有快速发布记录（包括不符合条件的） SELECT \'所有快速发布记录\' AS \'类型\', id, title, status, CASE status\n    WHEN 1 THEN \'待发布\'\n    WHEN 2 THEN \'发布中\'\n    WHEN 3 THEN \'发布成功\'\n    WHEN 4 THEN \'发布失败\'\n    WHEN 5 THEN \'已取消\'\n  END AS \'状态文本\', DATE(created_at) as \'日期\', created_at AS \'创建时间\'\nFROM  quick_deployments\nWHERE  created_at >= DATE_SUB(CURDATE(), INTERVAL 6 DAY)\nORDER BY  created_at DESC\nLIMIT  10; -- =========================================== -- 3. 检查服务上线记录 (近7天) -- =========================================== SELECT \'服务上线-近7天\' AS \'统计类型\', DATE(created_at) as \'日期\', COUNT(*) as \'数量\'\nFROM  app_service_release\nWHERE\n  created_at >= DATE_SUB(CURDATE(), INTERVAL 6 DAY)\n  AND created_at < DATE_ADD(CURDATE(), INTERVAL 1 DAY)\n  AND deleted_at IS NULL\n  AND status = 2\nGROUP BY  DATE(created_at)\nORDER BY  DATE(created_at) DESC; -- 查看所有服务上线记录（包括不符合条件的） SELECT \'所有服务上线记录\' AS \'类型\', id, title, status, CASE status\n    WHEN 1 THEN \'流程进行中\'\n    WHEN 2 THEN \'发布成功\'\n    WHEN 3 THEN \'发布失败\'\n    WHEN 4 THEN \'已驳回\'\n    WHEN 5 THEN \'已取消\'\n  END AS \'状态文本\', DATE(created_at) as \'日期\', created_at AS \'创建时间\', deleted_at AS \'删除时间\'\nFROM  app_service_release\nWHERE  created_at >= DATE_SUB(CURDATE(), INTERVAL 6 DAY)\nORDER BY  created_at DESC\nLIMIT  10; -- =========================================== -- 4. 查看当前时间和时区 -- =========================================== SELECT NOW() AS \'当前时间\', CURDATE() AS \'当前日期\', @@system_time_zone AS \'系统时区\', @@session.time_zone AS \'会话时区\'; -- =========================================== -- 5. 测试你提到的两条记录 -- =========================================== -- 检查 2025-12-25 11:36:31 的脚本发布 SELECT \'你的脚本发布记录\' AS \'类型\', id, title, status, DATE(created_at) as \'日期\', created_at AS \'创建时间\', deleted_at AS \'删除时间\', CASE WHEN status = 2 AND deleted_at IS NULL THEN \'✅ 应该被统计\'\n    WHEN status != 2 THEN \'❌ 状态不符合(需要status=2)\'\n    WHEN deleted_at IS NOT NULL THEN \'❌ 已被软删除\'\n  END AS \'是否统计\'\nFROM  app_sh_release\nWHERE\n  created_at >= \'2025-12-25 11:36:00\'\n  AND created_at <= \'2025-12-25 11:37:00\'\nORDER BY  created_at DESC; -- 检查 2025-12-25 11:57:45 的快速发布 SELECT \'你的快速发布记录\' AS \'类型\', id, title, status, DATE(created_at) as \'日期\', created_at AS \'创建时间\', CASE WHEN status = 3 THEN \'✅ 应该被统计\'\n    WHEN status != 3 THEN CONCAT(\'❌ 状态不符合(当前status=\', status, \', 需要status=3)\')\n  END AS \'是否统计\'\nFROM  quick_deployments\nWHERE\n  created_at >= \'2025-12-25 11:57:00\'\n  AND created_at <= \'2025-12-25 11:58:00\'\nORDER BY  created_at DESC;','','',2,'Error 1064 (42000): You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near \'FROM  app_sh_release\nWHERE\n  created_at >= DATE_SUB(CURDATE(), INTERVAL 6 DAY)\n \' at line 2',202,'2025-12-25 12:11:10.622','',0),(482,1,'gin-api','app_sh_release',5,'-- 上线趋势统计调试SQL\n-- 用于排查为什么最新的发布记录没有被统计到\n\n-- ===========================================\n-- 1. 检查脚本上线记录 (近7天)\n-- ===========================================\nSELECT\n    \'脚本上线-近7天\' AS `统计类型`,\n    DATE(created_at) as `日期`,\n    COUNT(*) as `数量`\nFROM app_sh_release\nWHERE created_at >= DATE_SUB(CURDATE(), INTERVAL 6 DAY)\n  AND created_at < DATE_ADD(CURDATE(), INTERVAL 1 DAY)\n  AND deleted_at IS NULL\n  AND status = 2\nGROUP BY DATE(created_at)\nORDER BY DATE(created_at) DESC;\n\n-- 查看所有脚本上线记录（包括不符合条件的）\nSELECT\n    \'所有脚本上线记录\' AS `类型`,\n    id,\n    title,\n    status,\n    CASE status\n        WHEN 1 THEN \'流程进行中\'\n        WHEN 2 THEN \'执行成功\'\n        WHEN 3 THEN \'执行失败\'\n        WHEN 4 THEN \'已驳回\'\n        WHEN 5 THEN \'已取消\'\n    END AS `状态文本`,\n    DATE(created_at) as `日期`,\n    created_at AS `创建时间`,\n    deleted_at AS `删除时间`\nFROM app_sh_release\nWHERE created_at >= DATE_SUB(CURDATE(), INTERVAL 6 DAY)\nORDER BY created_at DESC\nLIMIT 10;\n\n-- ===========================================\n-- 2. 检查快速发布记录 (近7天)\n-- ===========================================\nSELECT\n    \'快速发布-近7天\' AS `统计类型`,\n    DATE(created_at) as `日期`,\n    COUNT(*) as `数量`\nFROM quick_deployments\nWHERE created_at >= DATE_SUB(CURDATE(), INTERVAL 6 DAY)\n  AND created_at < DATE_ADD(CURDATE(), INTERVAL 1 DAY)\n  AND status = 3\nGROUP BY DATE(created_at)\nORDER BY DATE(created_at) DESC;\n\n-- 查看所有快速发布记录（包括不符合条件的）\nSELECT\n    \'所有快速发布记录\' AS `类型`,\n    id,\n    title,\n    status,\n    CASE status\n        WHEN 1 THEN \'待发布\'\n        WHEN 2 THEN \'发布中\'\n        WHEN 3 THEN \'发布成功\'\n        WHEN 4 THEN \'发布失败\'\n        WHEN 5 THEN \'已取消\'\n    END AS `状态文本`,\n    DATE(created_at) as `日期`,\n    created_at AS `创建时间`\nFROM quick_deployments\nWHERE created_at >= DATE_SUB(CURDATE(), INTERVAL 6 DAY)\nORDER BY created_at DESC\nLIMIT 10;\n\n-- ===========================================\n-- 3. 检查服务上线记录 (近7天)\n-- ===========================================\nSELECT\n    \'服务上线-近7天\' AS `统计类型`,\n    DATE(created_at) as `日期`,\n    COUNT(*) as `数量`\nFROM app_service_release\nWHERE created_at >= DATE_SUB(CURDATE(), INTERVAL 6 DAY)\n  AND created_at < DATE_ADD(CURDATE(), INTERVAL 1 DAY)\n  AND deleted_at IS NULL\n  AND status = 2\nGROUP BY DATE(created_at)\nORDER BY DATE(created_at) DESC;\n\n-- 查看所有服务上线记录（包括不符合条件的）\nSELECT\n    \'所有服务上线记录\' AS `类型`,\n    id,\n    title,\n    status,\n    CASE status\n        WHEN 1 THEN \'流程进行中\'\n        WHEN 2 THEN \'发布成功\'\n        WHEN 3 THEN \'发布失败\'\n        WHEN 4 THEN \'已驳回\'\n        WHEN 5 THEN \'已取消\'\n    END AS `状态文本`,\n    DATE(created_at) as `日期`,\n    created_at AS `创建时间`,\n    deleted_at AS `删除时间`\nFROM app_service_release\nWHERE created_at >= DATE_SUB(CURDATE(), INTERVAL 6 DAY)\nORDER BY created_at DESC\nLIMIT 10;\n\n-- ===========================================\n-- 4. 查看当前时间和时区\n-- ===========================================\nSELECT\n    NOW() AS `当前时间`,\n    CURDATE() AS `当前日期`,\n    @@system_time_zone AS `系统时区`,\n    @@session.time_zone AS `会话时区`;\n\n-- ===========================================\n-- 5. 测试你提到的两条记录\n-- ===========================================\n-- 检查 2025-12-25 11:36:31 的脚本发布\nSELECT\n    \'你的脚本发布记录\' AS `类型`,\n    id,\n    title,\n    status,\n    DATE(created_at) as `日期`,\n    created_at AS `创建时间`,\n    deleted_at AS `删除时间`,\n    CASE\n        WHEN status = 2 AND deleted_at IS NULL THEN \'✅ 应该被统计\'\n        WHEN status != 2 THEN \'❌ 状态不符合(需要status=2)\'\n        WHEN deleted_at IS NOT NULL THEN \'❌ 已被软删除\'\n    END AS `是否统计`\nFROM app_sh_release\nWHERE created_at >= \'2025-12-25 11:36:00\'\n  AND created_at <= \'2025-12-25 11:37:00\'\nORDER BY created_at DESC;\n\n-- 检查 2025-12-25 11:57:45 的快速发布\nSELECT\n    \'你的快速发布记录\' AS `类型`,\n    id,\n    title,\n    status,\n    DATE(created_at) as `日期`,\n    created_at AS `创建时间`,\n    CASE\n        WHEN status = 3 THEN \'✅ 应该被统计\'\n        WHEN status != 3 THEN CONCAT(\'❌ 状态不符合(当前status=\', status, \', 需要status=3)\')\n    END AS `是否统计`\nFROM quick_deployments\nWHERE created_at >= \'2025-12-25 11:57:00\'\n  AND created_at <= \'2025-12-25 11:58:00\'\nORDER BY created_at DESC;','','',2,'Error 1064 (42000): You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near \'SELECT\n    \'所有脚本上线记录\' AS `类型`,\n    id,\n    title,\n    status\' at line 20',213,'2025-12-25 12:13:32.454','',0),(483,1,'gin-api','app_sh_release',1,'-- 上线趋势统计调试SQL\n-- 用于排查为什么最新的发布记录没有被统计到\n\n-- ===========================================\n-- 1. 检查脚本上线记录 (近7天)\n-- ===========================================\nSELECT\n    \'脚本上线-近7天\' AS `统计类型`,\n    DATE(created_at) as `日期`,\n    COUNT(*) as `数量`\nFROM app_sh_release\nWHERE created_at >= DATE_SUB(CURDATE(), INTERVAL 6 DAY)\n  AND created_at < DATE_ADD(CURDATE(), INTERVAL 1 DAY)\n  AND deleted_at IS NULL\n  AND status = 2\nGROUP BY DATE(created_at)\nORDER BY DATE(created_at) DESC;\n\n-- 查看所有脚本上线记录（包括不符合条件的）\nSELECT\n    \'所有脚本上线记录\' AS `类型`,\n    id,\n    title,\n    status,\n    CASE status\n        WHEN 1 THEN \'流程进行中\'\n        WHEN 2 THEN \'执行成功\'\n        WHEN 3 THEN \'执行失败\'\n        WHEN 4 THEN \'已驳回\'\n        WHEN 5 THEN \'已取消\'\n    END AS `状态文本`,\n    DATE(created_at) as `日期`,\n    created_at AS `创建时间`,\n    deleted_at AS `删除时间`\nFROM app_sh_release\nWHERE created_at >= DATE_SUB(CURDATE(), INTERVAL 6 DAY)\nORDER BY created_at DESC\nLIMIT 10;\n\n-- ===========================================\n-- 2. 检查快速发布记录 (近7天)\n-- ===========================================\nSELECT\n    \'快速发布-近7天\' AS `统计类型`,\n    DATE(created_at) as `日期`,\n    COUNT(*) as `数量`\nFROM quick_deployments\nWHERE created_at >= DATE_SUB(CURDATE(), INTERVAL 6 DAY)\n  AND created_at < DATE_ADD(CURDATE(), INTERVAL 1 DAY)\n  AND status = 3\nGROUP BY DATE(created_at)\nORDER BY DATE(created_at) DESC;\n\n-- 查看所有快速发布记录（包括不符合条件的）\nSELECT\n    \'所有快速发布记录\' AS `类型`,\n    id,\n    title,\n    status,\n    CASE status\n        WHEN 1 THEN \'待发布\'\n        WHEN 2 THEN \'发布中\'\n        WHEN 3 THEN \'发布成功\'\n        WHEN 4 THEN \'发布失败\'\n        WHEN 5 THEN \'已取消\'\n    END AS `状态文本`,\n    DATE(created_at) as `日期`,\n    created_at AS `创建时间`\nFROM quick_deployments\nWHERE created_at >= DATE_SUB(CURDATE(), INTERVAL 6 DAY)\nORDER BY created_at DESC\nLIMIT 10;\n\n-- ===========================================\n-- 3. 检查服务上线记录 (近7天)\n-- ===========================================\nSELECT\n    \'服务上线-近7天\' AS `统计类型`,\n    DATE(created_at) as `日期`,\n    COUNT(*) as `数量`\nFROM app_service_release\nWHERE created_at >= DATE_SUB(CURDATE(), INTERVAL 6 DAY)\n  AND created_at < DATE_ADD(CURDATE(), INTERVAL 1 DAY)\n  AND deleted_at IS NULL\n  AND status = 2\nGROUP BY DATE(created_at)\nORDER BY DATE(created_at) DESC;\n\n-- 查看所有服务上线记录（包括不符合条件的）\nSELECT\n    \'所有服务上线记录\' AS `类型`,\n    id,\n    title,\n    status,\n    CASE status\n        WHEN 1 THEN \'流程进行中\'\n        WHEN 2 THEN \'发布成功\'\n        WHEN 3 THEN \'发布失败\'\n        WHEN 4 THEN \'已驳回\'\n        WHEN 5 THEN \'已取消\'\n    END AS `状态文本`,\n    DATE(created_at) as `日期`,\n    created_at AS `创建时间`,\n    deleted_at AS `删除时间`\nFROM app_service_release\nWHERE created_at >= DATE_SUB(CURDATE(), INTERVAL 6 DAY)\nORDER BY created_at DESC\nLIMIT 10;\n\n-- ===========================================\n-- 4. 查看当前时间和时区\n-- ===========================================\nSELECT\n    NOW() AS `当前时间`,\n    CURDATE() AS `当前日期`,\n    @@system_time_zone AS `系统时区`,\n    @@session.time_zone AS `会话时区`;\n\n-- ===========================================\n-- 5. 测试你提到的两条记录\n-- ===========================================\n-- 检查 2025-12-25 11:36:31 的脚本发布\nSELECT\n    \'你的脚本发布记录\' AS `类型`,\n    id,\n    title,\n    status,\n    DATE(created_at) as `日期`,\n    created_at AS `创建时间`,\n    deleted_at AS `删除时间`,\n    CASE\n        WHEN status = 2 AND deleted_at IS NULL THEN \'✅ 应该被统计\'\n        WHEN status != 2 THEN \'❌ 状态不符合(需要status=2)\'\n        WHEN deleted_at IS NOT NULL THEN \'❌ 已被软删除\'\n    END AS `是否统计`\nFROM app_sh_release\nWHERE created_at >= \'2025-12-25 11:36:00\'\n  AND created_at <= \'2025-12-25 11:37:00\'\nORDER BY created_at DESC;\n\n-- 检查 2025-12-25 11:57:45 的快速发布\nSELECT\n    \'你的快速发布记录\' AS `类型`,\n    id,\n    title,\n    status,\n    DATE(created_at) as `日期`,\n    created_at AS `创建时间`,\n    CASE\n        WHEN status = 3 THEN \'✅ 应该被统计\'\n        WHEN status != 3 THEN CONCAT(\'❌ 状态不符合(当前status=\', status, \', 需要status=3)\')\n    END AS `是否统计`\nFROM quick_deployments\nWHERE created_at >= \'2025-12-25 11:57:00\'\n  AND created_at <= \'2025-12-25 11:58:00\'\nORDER BY created_at DESC;','','',1,'执行成功',1403,'2025-12-25 12:19:38.925','',0),(484,1,'devops','app_service_release',1,'SELECT * FROM app_service_release LIMIT 20 OFFSET 0','','',1,'执行成功',311,'2025-12-25 22:54:47.702','',0),(485,1,'devops','app_service_release',1,'SELECT * FROM app_service_release LIMIT 20 OFFSET 0','','',1,'执行成功',642,'2025-12-25 22:54:47.900','',0),(486,1,'devops','app_service_release',1,'SELECT * FROM app_service_release LIMIT 20 OFFSET 0','','',1,'执行成功',278,'2025-12-25 22:54:52.507','',0),(487,1,'gin-api','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',257,'2025-12-26 14:00:52.509','',0),(488,1,'gin-api','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',255,'2025-12-26 15:12:14.551','',0),(489,1,'gin-api','cmdb_host',4,'DELETE FROM cmdb_host\nWHERE create_time >= \'2025-12-26 00:00:00\'\n  AND create_time < \'2025-12-27 00:00:00\';','','',1,'执行成功',226,'2025-12-26 15:12:27.345','',0),(490,1,'gin-api','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',252,'2025-12-26 15:12:35.241','',0),(491,1,'gin-api','cmdb_sql',1,'SELECT * FROM cmdb_sql LIMIT 20 OFFSET 0','','',1,'执行成功',271,'2025-12-27 23:04:13.781','',0),(492,1,'gin-api','',5,'ALTER TABLE monitor_domain_schedule \n  DROP COLUMN interval_mins,\n  ADD COLUMN cron_expr VARCHAR(100) COMMENT \'Cron表达式\' AFTER enabled;','','',1,'执行成功',254,'2025-12-27 23:04:24.413','',0),(493,1,'gin-api','monitor_domain_schedule',1,'SELECT * FROM monitor_domain_schedule LIMIT 20 OFFSET 0','','',1,'执行成功',396,'2025-12-27 23:04:37.906','',0),(494,1,'gin-api','monitor_domain_schedule',1,'SELECT * FROM monitor_domain_schedule LIMIT 20 OFFSET 0','','',1,'执行成功',283,'2025-12-27 23:04:47.881','',0),(495,1,'gin-api','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',281,'2025-12-28 13:05:35.747','',0),(496,1,'gin-api','cmdb_host',4,'DELETE FROM cmdb_host\nWHERE create_time >= NOW() - INTERVAL 1 DAY;','','',1,'执行成功',237,'2025-12-28 13:15:53.019','',0),(497,1,'gin-api','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',273,'2025-12-28 13:15:58.982','',0),(498,1,'gin-api','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',274,'2025-12-28 13:16:00.930','',0),(499,1,'gin-api','cmdb_host',4,'DELETE FROM cmdb_host\nWHERE create_time >= DATE_SUB(NOW(), INTERVAL 1 DAY);','','',1,'执行成功',215,'2025-12-28 13:16:38.008','',0),(500,1,'gin-api','cmdb_host',4,'DELETE FROM cmdb_host\nWHERE DATE(create_time) = CURDATE();','','',1,'执行成功',203,'2025-12-28 13:16:55.887','',0),(501,1,'gin-api','',5,'TRUNCATE TABLE cmdb_host;','','',1,'执行成功',281,'2025-12-28 13:22:10.901','',0),(502,1,'gin-api','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',271,'2025-12-28 13:22:17.241','',0),(503,1,'gin-api','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',295,'2025-12-28 13:47:17.865','',0),(504,1,'gin-api','',5,'TRUNCATE TABLE cmdb_host;','','',1,'执行成功',359,'2025-12-28 13:47:31.736','',0),(505,1,'gin-api','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',267,'2025-12-28 13:47:38.145','',0),(506,1,'gin-api','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',279,'2025-12-28 15:38:07.744','',0),(507,1,'gin-api','',5,'TRUNCATE TABLE cmdb_host;','','',1,'执行成功',322,'2025-12-28 15:41:49.483','',0),(508,1,'gin-api','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',295,'2025-12-28 15:41:55.663','',0),(509,1,'gin-api','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',284,'2026-01-07 15:00:19.261','',0),(510,1,'gin-api','',5,'TRUNCATE TABLE cmdb_host;','','',1,'执行成功',280,'2026-01-07 15:01:03.819','',0),(511,1,'gin-api','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',253,'2026-01-07 15:01:11.610','',0),(512,1,'gin-api','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',256,'2026-01-07 15:49:17.511','',0),(513,1,'gin-api','',5,'TRUNCATE TABLE cmdb_host;','','',1,'执行成功',264,'2026-01-07 15:50:33.865','',0),(514,1,'gin-api','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',284,'2026-01-07 15:50:38.659','',0),(515,1,'gin-api','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',273,'2026-01-08 20:04:08.959','',0),(516,1,'gin-api','',5,'TRUNCATE TABLE cmdb_host;','','',1,'执行成功',300,'2026-01-08 20:06:35.649','',0),(517,1,'gin-api','cmdb_host',1,'SELECT * FROM cmdb_host LIMIT 20 OFFSET 0','','',1,'执行成功',272,'2026-01-08 20:06:39.629','',0),(518,1,'devops','app_sh_release',1,'SELECT * FROM app_sh_release LIMIT 20 OFFSET 0','','',1,'执行成功',460,'2026-01-08 23:09:00.812','',0),(519,1,'autoops','quick_deployments',1,'SELECT * FROM quick_deployments LIMIT 20 OFFSET 0','','',1,'执行成功',549,'2026-01-19 17:04:49.659','',0),(520,1,'autoops','bastion_user_group',1,'SELECT * FROM bastion_user_group LIMIT 20 OFFSET 0','','',1,'执行成功',256,'2026-01-19 17:14:50.760','',0),(521,1,'autoops','bastion_host_authorization',1,'SELECT * FROM bastion_host_authorization LIMIT 20 OFFSET 0','','',1,'执行成功',324,'2026-01-19 17:15:25.235','',0),(522,1,'autoops','quick_deployment_tasks',1,'SELECT * FROM quick_deployment_tasks LIMIT 20 OFFSET 0','','',1,'执行成功',271,'2026-01-19 17:15:42.080','',0),(523,1,'autoops','quick_deployments',1,'SELECT * FROM quick_deployments LIMIT 20 OFFSET 0','','',1,'执行成功',348,'2026-01-19 17:15:43.521','',0);
/*!40000 ALTER TABLE `db_sql_exec` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `k8s_cluster`
--

DROP TABLE IF EXISTS `k8s_cluster`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `k8s_cluster` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键ID''',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''集群名称''',
  `version` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''集群版本''',
  `status` bigint NOT NULL DEFAULT '1' COMMENT '''集群状态:1-创建中,2-运行中,3-离线''',
  `credential` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''集群凭证(kubeconfig)''',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''集群描述''',
  `cluster_type` bigint NOT NULL DEFAULT '1' COMMENT '''集群类型:1-自建,2-导入''',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `node_count` bigint DEFAULT '0' COMMENT '''节点数量''',
  `ready_nodes` bigint DEFAULT '0' COMMENT '''就绪节点数''',
  `master_nodes` bigint DEFAULT '0' COMMENT '''Master节点数''',
  `worker_nodes` bigint DEFAULT '0' COMMENT '''Worker节点数''',
  `last_sync_at` datetime(3) DEFAULT NULL COMMENT '''最后同步时间''',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `idx_k8s_cluster_name` (`name`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=34 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `k8s_cluster`
--

LOCK TABLES `k8s_cluster` WRITE;
/*!40000 ALTER TABLE `k8s_cluster` DISABLE KEYS */;
/*!40000 ALTER TABLE `k8s_cluster` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `knowledge_base`
--

DROP TABLE IF EXISTS `knowledge_base`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `knowledge_base` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键''',
  `type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''类型:command/solution/best-practice/faq''',
  `category` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''分类''',
  `title` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''标题''',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''内容''',
  `keywords` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''关键词(JSON数组)''',
  `tags` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''标签(JSON数组)''',
  `score` decimal(3,2) DEFAULT '0.50' COMMENT '''评分''',
  `use_count` bigint DEFAULT '0' COMMENT '''使用次数''',
  `enabled` bigint DEFAULT '1' COMMENT '''是否启用:1=启用,2=禁用''',
  `create_time` datetime(3) NOT NULL COMMENT '''创建时间''',
  `update_time` datetime(3) NOT NULL COMMENT '''更新时间''',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_knowledge_base_type` (`type`) USING BTREE,
  KEY `idx_knowledge_base_category` (`category`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `knowledge_base`
--

LOCK TABLES `knowledge_base` WRITE;
/*!40000 ALTER TABLE `knowledge_base` DISABLE KEYS */;
/*!40000 ALTER TABLE `knowledge_base` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `monitor_agent`
--

DROP TABLE IF EXISTS `monitor_agent`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `monitor_agent` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键''',
  `host_id` bigint unsigned NOT NULL COMMENT '''主机ID''',
  `host_name` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''主机名称''',
  `version` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '1.0.0' COMMENT '''Agent版本''',
  `status` bigint DEFAULT NULL COMMENT '''状态:1->部署中,2->部署失败,3->运行中,4->已停止''',
  `install_path` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''安装路径''',
  `port` bigint DEFAULT '9100' COMMENT '''监听端口''',
  `pid` bigint DEFAULT NULL COMMENT '''进程ID''',
  `last_heartbeat` datetime(3) DEFAULT NULL COMMENT '''最后心跳时间''',
  `update_time` datetime(3) DEFAULT NULL COMMENT '''更新时间''',
  `create_time` datetime(3) NOT NULL COMMENT '''创建时间''',
  `error_msg` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''错误信息''',
  `install_progress` bigint DEFAULT '0' COMMENT '''安装进度(0-100)''',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_monitor_agent_host_id` (`host_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=120 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `monitor_agent`
--

LOCK TABLES `monitor_agent` WRITE;
/*!40000 ALTER TABLE `monitor_agent` DISABLE KEYS */;
/*!40000 ALTER TABLE `monitor_agent` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `monitor_alert_source`
--

DROP TABLE IF EXISTS `monitor_alert_source`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `monitor_alert_source` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键''',
  `name` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''告警源名称''',
  `type` bigint NOT NULL COMMENT '''告警源类型(1=flashduty,2=zabbix,3=n9e)''',
  `app_key` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''应用密钥''',
  `api_base_url` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''API基础URL''',
  `status` bigint DEFAULT '1' COMMENT '''状态(0=禁用,1=启用)''',
  `remark` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''备注''',
  `create_time` bigint DEFAULT NULL COMMENT '''创建时间''',
  `update_time` bigint DEFAULT NULL COMMENT '''更新时间''',
  `key_id` bigint unsigned DEFAULT '0' COMMENT '''云账号ID''',
  `host_id` bigint unsigned DEFAULT '0' COMMENT '''主机ID''',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `monitor_alert_source`
--

LOCK TABLES `monitor_alert_source` WRITE;
/*!40000 ALTER TABLE `monitor_alert_source` DISABLE KEYS */;
/*!40000 ALTER TABLE `monitor_alert_source` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `monitor_aliyun_config`
--

DROP TABLE IF EXISTS `monitor_aliyun_config`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `monitor_aliyun_config` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键''',
  `name` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''配置名称''',
  `access_key` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''AccessKey''',
  `access_secret` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''AccessSecret''',
  `region` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT 'cn-hangzhou' COMMENT '''区域''',
  `email` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''联系邮箱''',
  `phone` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''联系电话''',
  `username` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''联系人姓名''',
  `status` bigint DEFAULT '1' COMMENT '''状态:1->启用,0->禁用''',
  `create_time` datetime(3) NOT NULL COMMENT '''创建时间''',
  `update_time` datetime(3) DEFAULT NULL COMMENT '''更新时间''',
  `eab_kid` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''ACME EAB Key ID (ZeroSSL等需要)''',
  `eab_hmac_key` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''ACME EAB HMAC Key (ZeroSSL等需要)''',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `monitor_aliyun_config`
--

LOCK TABLES `monitor_aliyun_config` WRITE;
/*!40000 ALTER TABLE `monitor_aliyun_config` DISABLE KEYS */;
INSERT INTO `monitor_aliyun_config` VALUES (2,'阿里云DNS','LTAI5tHGT2uC14zHAxfWko1J','3pkc****Emv5','','zfwh1024@163.com','13826541566','zhangfan',1,'2026-01-17 14:19:54.333','2026-01-18 23:48:35.636','','');
/*!40000 ALTER TABLE `monitor_aliyun_config` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `monitor_domain`
--

DROP TABLE IF EXISTS `monitor_domain`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `monitor_domain` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `domain` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `tags` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `remark` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
  `status` bigint DEFAULT '1' COMMENT '状态(1:启用,0:禁用)',
  `is_alive` bigint DEFAULT '0' COMMENT '存活状态(1:正常,0:异常)',
  `status_code` bigint DEFAULT NULL COMMENT 'HTTP状态码',
  `response_time` bigint DEFAULT NULL COMMENT '响应时间(ms)',
  `ssl_expire_at` datetime DEFAULT NULL COMMENT 'SSL证书过期时间',
  `ssl_days_left` bigint DEFAULT NULL COMMENT 'SSL证书剩余天数',
  `ssl_issuer` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT 'SSL证书颁发者',
  `last_check_at` datetime DEFAULT NULL COMMENT '最后检查时间',
  `error_msg` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '错误信息',
  `create_time` datetime(3) DEFAULT NULL,
  `update_time` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `idx_monitor_domain_domain` (`domain`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=43 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `monitor_domain`
--

LOCK TABLES `monitor_domain` WRITE;
/*!40000 ALTER TABLE `monitor_domain` DISABLE KEYS */;
INSERT INTO `monitor_domain` VALUES (6,'gitee.com','','',1,1,200,501,'2026-02-17 23:59:59',29,'TrustAsia DV TLS RSA CA 2025','2026-01-19 13:29:01','','2025-12-04 10:21:13.921','2026-01-19 13:29:00.997');
/*!40000 ALTER TABLE `monitor_domain` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `monitor_domain_schedule`
--

DROP TABLE IF EXISTS `monitor_domain_schedule`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `monitor_domain_schedule` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `enabled` tinyint(1) DEFAULT '0' COMMENT '是否启用',
  `cron_expr` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT 'Cron表达式',
  `next_run_at` datetime DEFAULT NULL COMMENT '下次执行时间',
  `last_run_at` datetime DEFAULT NULL COMMENT '上次执行时间',
  `status` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '任务状态(running/paused/stopped)',
  `create_time` datetime(3) DEFAULT NULL,
  `update_time` datetime(3) DEFAULT NULL,
  `notify_enabled` tinyint(1) DEFAULT '0' COMMENT '是否开启通知',
  `notify_robot_id` bigint DEFAULT NULL COMMENT '通知机器人ID',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `monitor_domain_schedule`
--

LOCK TABLES `monitor_domain_schedule` WRITE;
/*!40000 ALTER TABLE `monitor_domain_schedule` DISABLE KEYS */;
INSERT INTO `monitor_domain_schedule` VALUES (1,1,'29 13 * * *',NULL,NULL,'running','2025-12-04 23:09:32.660','2026-01-30 22:27:59.819',1,1),(2,0,NULL,NULL,NULL,'stopped','2025-12-04 23:09:32.660','2025-12-04 23:09:32.660',0,NULL);
/*!40000 ALTER TABLE `monitor_domain_schedule` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `monitor_incident`
--

DROP TABLE IF EXISTS `monitor_incident`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `monitor_incident` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键''',
  `alert_time` datetime(3) NOT NULL COMMENT '''告警时间''',
  `business_line` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''业务线(saas3/saas4/saas5等)''',
  `frequency` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''频率(偶发/频繁)''',
  `alert_desc` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''告警描述''',
  `alert_level` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT 'P4' COMMENT '''告警级别(P1/P2/P3/P4)''',
  `incident_cause` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''故障原因''',
  `department` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''所属部门(研发部/运维部)''',
  `solution` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''解决方案''',
  `detail_url` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''故障详情URL链接''',
  `handler` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''处理人(用户名)''',
  `handler_id` bigint unsigned DEFAULT NULL COMMENT '''处理人ID''',
  `status` bigint DEFAULT '1' COMMENT '''处理状态:1->未处理,2->处理中,3->已完成''',
  `remark` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''备注信息''',
  `create_time` datetime(3) NOT NULL COMMENT '''创建时间''',
  `update_time` datetime(3) DEFAULT NULL COMMENT '''更新时间''',
  `business_line_id` bigint unsigned DEFAULT NULL COMMENT '''业务线ID(关联cmdb_group二级分组)''',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `monitor_incident`
--

LOCK TABLES `monitor_incident` WRITE;
/*!40000 ALTER TABLE `monitor_incident` DISABLE KEYS */;
INSERT INTO `monitor_incident` VALUES (4,'2025-12-05 00:00:00.000',NULL,'偶发','al-hz-s3-pgsql-prod01-slave1   事务执行耗时_warning  379','P4','','研发部','','https://us-west-2.console.aws.amazon.com/eks/clusters/aws-uw2-lku-k8s-test01?region=us-west-2&selectedResourceId=services&selectedTab=cluster-resources-tab','admin',89,1,'','2025-12-10 17:47:02.925','2025-12-10 17:52:44.598',16),(7,'2025-12-06 00:00:00.000',NULL,'偶发','al-hz-b-ops-monitor-prod01磁盘使用超过90%','P4','','运维部','','','admin',89,3,'','2025-12-10 17:48:39.865','2025-12-10 17:48:39.865',4),(8,'2025-12-05 00:00:00.000',NULL,'偶发','sun阳光雨露项目cpu使用率打满 ','P4','','运维部','','','admin',89,3,'','2025-12-10 17:49:14.614','2025-12-10 17:49:14.614',4);
/*!40000 ALTER TABLE `monitor_incident` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `monitor_notify_robot`
--

DROP TABLE IF EXISTS `monitor_notify_robot`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `monitor_notify_robot` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '机器人名称',
  `type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '机器人类型(dingtalk/feishu/wechat/email/webhook/teams)',
  `webhook` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT 'webhook地址或邮箱接收地址',
  `secret` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '加签密钥或邮箱密码',
  `status` tinyint DEFAULT '1' COMMENT '状态(1:启用,0:禁用)',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '备注说明',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `server` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT 'SMTP服务器地址',
  `port` bigint DEFAULT NULL COMMENT 'SMTP端口',
  `username` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '邮箱用户名',
  `password` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '邮箱密码',
  `nickname` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '发件人昵称',
  `headers` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '自定义HTTP头(JSON格式)',
  `method` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT 'POST' COMMENT 'HTTP方法(GET/POST/PUT)',
  `template` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '消息模板(支持变量替换)',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `monitor_notify_robot`
--

LOCK TABLES `monitor_notify_robot` WRITE;
/*!40000 ALTER TABLE `monitor_notify_robot` DISABLE KEYS */;
/*!40000 ALTER TABLE `monitor_notify_robot` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `monitor_ssl_cert`
--

DROP TABLE IF EXISTS `monitor_ssl_cert`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `monitor_ssl_cert` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键''',
  `domain` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''域名''',
  `aliyun_config_id` bigint unsigned DEFAULT NULL COMMENT '''阿里云配置ID''',
  `order_id` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''证书订单ID''',
  `cert_id` bigint DEFAULT NULL COMMENT '''证书ID''',
  `cert_name` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''证书名称''',
  `product_code` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT 'digicert-free-1-free' COMMENT '''产品代码''',
  `status` bigint DEFAULT NULL COMMENT '''状态:1->申请中,2->验证中,3->已签发,4->已过期,5->申请失败''',
  `validate_type` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT 'DNS' COMMENT '''验证方式:DNS/FILE''',
  `validate_info` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''验证信息(DNS记录或FILE内容)''',
  `cert` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''证书内容''',
  `private_key` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''私钥内容''',
  `issue_time` datetime(3) DEFAULT NULL COMMENT '''签发时间''',
  `expire_time` datetime(3) DEFAULT NULL COMMENT '''过期时间''',
  `days_left` bigint DEFAULT NULL COMMENT '''剩余天数''',
  `error_msg` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''错误信息''',
  `create_time` datetime(3) NOT NULL COMMENT '''创建时间''',
  `update_time` datetime(3) DEFAULT NULL COMMENT '''更新时间''',
  `cert_source` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT 'aliyun_cas' COMMENT '''证书来源:aliyun_cas/acme''',
  `ca_provider` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT 'DigiCert' COMMENT '''CA提供商:DigiCert/LetsEncrypt/ZeroSSL''',
  `issuer_cert` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''颁发者证书''',
  `algorithm` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT 'RSA2048' COMMENT '''加密算法:RSA2048/EC256''',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_monitor_ssl_cert_domain` (`domain`) USING BTREE,
  KEY `idx_monitor_ssl_cert_aliyun_config_id` (`aliyun_config_id`) USING BTREE,
  KEY `idx_monitor_ssl_cert_order_id` (`order_id`) USING BTREE,
  KEY `idx_monitor_ssl_cert_cert_source` (`cert_source`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `monitor_ssl_cert`
--

LOCK TABLES `monitor_ssl_cert` WRITE;
/*!40000 ALTER TABLE `monitor_ssl_cert` DISABLE KEYS */;
INSERT INTO `monitor_ssl_cert` VALUES (17,'deviops.cn',2,'',0,'cert-deviops.cn-1768630863','digicert-free-1-free',3,'DNS','','-----BEGIN CERTIFICATE-----\nMIIE9jCCA96gAwIBAgISBeRJ56UC4TS2k2LGe0nYeNJNMA0GCSqGSIb3DQEBCwUA\nMDMxCzAJBgNVBAYTAlVTMRYwFAYDVQQKEw1MZXQncyBFbmNyeXB0MQwwCgYDVQQD\nEwNSMTIwHhcNMjYwMTE3MDUyMjQ2WhcNMjYwNDE3MDUyMjQ1WjAVMRMwEQYDVQQD\nEwpkZXZpb3BzLmNuMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA28hp\nNXaieAE4JgnFEzepa79agjXCY6UyKU3xuGjqFwUl/CMGjTSJPE4vdZ1J8iIsFk24\n+u/cBZOGhMNCgY2u9jDGFlPcqauRFwCdOs+5jliLUdMPTUYEIIJQCT6Z7DqlfoZE\n2z/Gqap02fkH0I/TtLNx5aIyU5BkvaDvsEKsZQ9Vhc57f/Q8uX1zqhDYFb+c2JHL\nWb5vIW+VLThmPYBKhwmTwRKB6LT8zqJ14GvIJWNhkwxszUii+HBphR/yWv4Xqq7u\niv56LRmfdbm3ckvQBKckbgwcFv3km9UpJGEFiMGCX60DzV1vbDQiGpR2dp7H92Zk\nqbEuqNBfwS60jqy3bQIDAQABo4ICIDCCAhwwDgYDVR0PAQH/BAQDAgWgMB0GA1Ud\nJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjAMBgNVHRMBAf8EAjAAMB0GA1UdDgQW\nBBSdZZCQclf8FyXPwnhKp4gmQeE3QjAfBgNVHSMEGDAWgBQAtSnyLY5vMeibTK14\nPvrc6QzR0jAzBggrBgEFBQcBAQQnMCUwIwYIKwYBBQUHMAKGF2h0dHA6Ly9yMTIu\naS5sZW5jci5vcmcvMBUGA1UdEQQOMAyCCmRldmlvcHMuY24wEwYDVR0gBAwwCjAI\nBgZngQwBAgEwLgYDVR0fBCcwJTAjoCGgH4YdaHR0cDovL3IxMi5jLmxlbmNyLm9y\nZy83OC5jcmwwggEKBgorBgEEAdZ5AgQCBIH7BIH4APYAdQBJnJtp3h187Pw23s2H\nZKa4W68Kh4AZ0VVS++nrKd34wwAAAZvKnXt3AAAEAwBGMEQCIAEWeb/isWuGtYXK\nXn3K8yAOeViSK6v7nZ9zUVcNaQQlAiBEPvmUMpAmwJG9iPB44Syvucc3rOYTcg1u\nzkvkXXEfqAB9ABqLnWlKV5jImaDKiL30j8C0VmDMw2ANH3H0af/H0ayjAAABm8qd\nfiwACAAABQA4UT08BAMARjBEAiA1cbgRwcH02CdOBXz4pa+UzsvKRn45C5wkl377\nyV7fewIgKa2u7ktyCwk0Z2CaO3AwK3T4SZuX1vuShVL7STw3WrwwDQYJKoZIhvcN\nAQELBQADggEBAI72Qir6LxdJyoQykMJoHnm6L8GIbhIsxZIFeAjd1HYwUOUphgBB\nP0gEXOKGtTBA7TR34bDdKGe3PbLHJ1imXFWR5fG6l0bU48TaZpYgCxJvg4u8xfy4\n4nKSt4z76MVUiFdK7SbWUN0tcCyQ7RlsWzAKZ/Pd1GUqjhOm3yBoj+oOIYllswPS\nmBg1v3BQLNTTZr6+Uf53LM0XkQgjsr1lmJYJnICtYIXB803GO+kPYBaID7C/KdjT\nTV1MDZRo1juf0WYeo5MAoTLKExV1tg3eENrNqiz9EEuJ9JtE9B7usLH2h2/PmG7i\n5dlN2tnFTQD8meGj5DdVIPnNrGrFwSgCxvs=\n-----END CERTIFICATE-----\n\n-----BEGIN CERTIFICATE-----\nMIIFBjCCAu6gAwIBAgIRAMISMktwqbSRcdxA9+KFJjwwDQYJKoZIhvcNAQELBQAw\nTzELMAkGA1UEBhMCVVMxKTAnBgNVBAoTIEludGVybmV0IFNlY3VyaXR5IFJlc2Vh\ncmNoIEdyb3VwMRUwEwYDVQQDEwxJU1JHIFJvb3QgWDEwHhcNMjQwMzEzMDAwMDAw\nWhcNMjcwMzEyMjM1OTU5WjAzMQswCQYDVQQGEwJVUzEWMBQGA1UEChMNTGV0J3Mg\nRW5jcnlwdDEMMAoGA1UEAxMDUjEyMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIB\nCgKCAQEA2pgodK2+lP474B7i5Ut1qywSf+2nAzJ+Npfs6DGPpRONC5kuHs0BUT1M\n5ShuCVUxqqUiXXL0LQfCTUA83wEjuXg39RplMjTmhnGdBO+ECFu9AhqZ66YBAJpz\nkG2Pogeg0JfT2kVhgTU9FPnEwF9q3AuWGrCf4yrqvSrWmMebcas7dA8827JgvlpL\nThjp2ypzXIlhZZ7+7Tymy05v5J75AEaz/xlNKmOzjmbGGIVwx1Blbzt05UiDDwhY\nXS0jnV6j/ujbAKHS9OMZTfLuevYnnuXNnC2i8n+cF63vEzc50bTILEHWhsDp7CH4\nWRt/uTp8n1wBnWIEwii9Cq08yhDsGwIDAQABo4H4MIH1MA4GA1UdDwEB/wQEAwIB\nhjAdBgNVHSUEFjAUBggrBgEFBQcDAgYIKwYBBQUHAwEwEgYDVR0TAQH/BAgwBgEB\n/wIBADAdBgNVHQ4EFgQUALUp8i2ObzHom0yteD763OkM0dIwHwYDVR0jBBgwFoAU\nebRZ5nu25eQBc4AIiMgaWPbpm24wMgYIKwYBBQUHAQEEJjAkMCIGCCsGAQUFBzAC\nhhZodHRwOi8veDEuaS5sZW5jci5vcmcvMBMGA1UdIAQMMAowCAYGZ4EMAQIBMCcG\nA1UdHwQgMB4wHKAaoBiGFmh0dHA6Ly94MS5jLmxlbmNyLm9yZy8wDQYJKoZIhvcN\nAQELBQADggIBAI910AnPanZIZTKS3rVEyIV29BWEjAK/duuz8eL5boSoVpHhkkv3\n4eoAeEiPdZLj5EZ7G2ArIK+gzhTlRQ1q4FKGpPPaFBSpqV/xbUb5UlAXQOnkHn3m\nFVj+qYv87/WeY+Bm4sN3Ox8BhyaU7UAQ3LeZ7N1X01xxQe4wIAAE3JVLUCiHmZL+\nqoCUtgYIFPgcg350QMUIWgxPXNGEncT921ne7nluI02V8pLUmClqXOsCwULw+PVO\nZCB7qOMxxMBoCUeL2Ll4oMpOSr5pJCpLN3tRA2s6P1KLs9TSrVhOk+7LX28NMUlI\nusQ/nxLJID0RhAeFtPjyOCOscQBA53+NRjSCak7P4A5jX7ppmkcJECL+S0i3kXVU\ny5Me5BbrU8973jZNv/ax6+ZK6TM8jWmimL6of6OrX7ZU6E2WqazzsFrLG3o2kySb\nzlhSgJ81Cl4tv3SbYiYXnJExKQvzf83DYotox3f0fwv7xln1A2ZLplCb0O+l/AK0\nYE0DS2FPxSAHi0iwMfW2nNHJrXcY3LLHD77gRgje4Eveubi2xxa+Nmk/hmhLdIET\niVDFanoCrMVIpQ59XWHkzdFmoHXHBV7oibVjGSO7ULSQ7MJ1Nz51phuDJSgAIU7A\n0zrLnOrAj/dfrlEWRhCvAgbuwLZX1A2sjNjXoPOHbsPiy+lO1KF8/XY7\n-----END CERTIFICATE-----\n','-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEA28hpNXaieAE4JgnFEzepa79agjXCY6UyKU3xuGjqFwUl/CMG\njTSJPE4vdZ1J8iIsFk24+u/cBZOGhMNCgY2u9jDGFlPcqauRFwCdOs+5jliLUdMP\nTUYEIIJQCT6Z7DqlfoZE2z/Gqap02fkH0I/TtLNx5aIyU5BkvaDvsEKsZQ9Vhc57\nf/Q8uX1zqhDYFb+c2JHLWb5vIW+VLThmPYBKhwmTwRKB6LT8zqJ14GvIJWNhkwxs\nzUii+HBphR/yWv4Xqq7uiv56LRmfdbm3ckvQBKckbgwcFv3km9UpJGEFiMGCX60D\nzV1vbDQiGpR2dp7H92ZkqbEuqNBfwS60jqy3bQIDAQABAoIBACDEUpEuuH65OdCL\n/NaGYism+ny2vJ5ECwQ2826Jr7rSYO0kIK4tMrertW4sztE7cIS7vjQlThGDwtao\n0NHFWJwqPeImyxt9sXIzPH2/6ftzroO6uomMAoIBQ6L5MxQR00oXhQ0GB12UxPqZ\nYuWgh9DlPvocMjIcsTRHWezWp+Kh/aPum1KNMOkb8N4JdaB9eLLYqy2lpuy5xhKQ\ngtPZD3uFrpNFuXmjRBQwST/LhEvBg/ibF/i1O0BUPBqWsvw16XeuotltHy19E0dX\nT8eP7yUnv7uecqyT538xoO/p7iB749ZETVPViuYfL1bKmactZCLiBO8BoIU+lnmu\npR0jYDkCgYEA76dKIzhP80cGu65kp3yujVS0iNDvLYfv3h41d9MRLIx864ANwPsp\nD9gMuuhZAX21gE+tD/9zUYpjnoRArMuyxBpzOO5MrkJX0c3SCFmBJmwjv7VBUQ+b\n8rfzuZdBNW0FrvLRQAfoWLeTrHEdXFDUhqH9XgjnIMlmIt3CRvBMQuUCgYEA6sYm\ng35hj7dInNpX55P5fqZhba25ZX9pM25JtWnT7AgT6LPFAewXugRqctEefeozeJn7\nzfhHkSyqp1X6HnAf7UsS67VPj6Bd5KzjJ5cjgOJXyQ+Jq5SxmXXcXMwKxCS4Vt6p\n9G96D4q8fs5+sfrLK6pByeaedKKRVHtF1wk+MekCgYAXX8GjXcKuf6bwdznp/Acq\nUbCqB3bMgg2rJFcisLtkFO+6KGsu/1Ap83CKak8EDja9pQfEH64nrc+5jjygnd15\nsvq5+sQP2vilxeVOWyx1iuaRtCQCu/sB2V80IsSBaXi8b4JpNFw9CuqkTwQGFPwl\nhXcozMV+NSZ5EBf6w0Us4QKBgAQ1j0U2Q9uGxreBNqhBP+zkYMTNo5BNS36PkKS/\nA/ysXGcEqhiRug8suxfZ5fkN9ORJpnkX5XiV221qnjRjgD8J5InzUzzY/DX4X447\nKHjMFhTVwXshDFobAdtsERlZKzDrhxKJ2R9kNmLzotPgEjpzMADFsLuYmpBwMV1c\nng/ZAoGBAJrR2lv/EXuCYayvlRrROtiTvociyehleKCJ5inyll7/x1lZSjUnENMD\n3Kycy+G1ZYhdyXWMAB5sxINqLGB2gPZx43Wm6l9zcGf2ETxxsyr7BFrYrVBEJL3q\noKStvHvxB+x9YARsmsaR5rkUwmg4YF0d/erGGbrxv7zRfHy7twhK\n-----END RSA PRIVATE KEY-----\n','2026-01-17 08:00:00.000','2026-04-17 08:00:00.000',89,'','2026-01-17 14:21:03.143','2026-01-17 14:21:18.381','acme','LetsEncrypt','\n-----BEGIN CERTIFICATE-----\nMIIFBjCCAu6gAwIBAgIRAMISMktwqbSRcdxA9+KFJjwwDQYJKoZIhvcNAQELBQAw\nTzELMAkGA1UEBhMCVVMxKTAnBgNVBAoTIEludGVybmV0IFNlY3VyaXR5IFJlc2Vh\ncmNoIEdyb3VwMRUwEwYDVQQDEwxJU1JHIFJvb3QgWDEwHhcNMjQwMzEzMDAwMDAw\nWhcNMjcwMzEyMjM1OTU5WjAzMQswCQYDVQQGEwJVUzEWMBQGA1UEChMNTGV0J3Mg\nRW5jcnlwdDEMMAoGA1UEAxMDUjEyMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIB\nCgKCAQEA2pgodK2+lP474B7i5Ut1qywSf+2nAzJ+Npfs6DGPpRONC5kuHs0BUT1M\n5ShuCVUxqqUiXXL0LQfCTUA83wEjuXg39RplMjTmhnGdBO+ECFu9AhqZ66YBAJpz\nkG2Pogeg0JfT2kVhgTU9FPnEwF9q3AuWGrCf4yrqvSrWmMebcas7dA8827JgvlpL\nThjp2ypzXIlhZZ7+7Tymy05v5J75AEaz/xlNKmOzjmbGGIVwx1Blbzt05UiDDwhY\nXS0jnV6j/ujbAKHS9OMZTfLuevYnnuXNnC2i8n+cF63vEzc50bTILEHWhsDp7CH4\nWRt/uTp8n1wBnWIEwii9Cq08yhDsGwIDAQABo4H4MIH1MA4GA1UdDwEB/wQEAwIB\nhjAdBgNVHSUEFjAUBggrBgEFBQcDAgYIKwYBBQUHAwEwEgYDVR0TAQH/BAgwBgEB\n/wIBADAdBgNVHQ4EFgQUALUp8i2ObzHom0yteD763OkM0dIwHwYDVR0jBBgwFoAU\nebRZ5nu25eQBc4AIiMgaWPbpm24wMgYIKwYBBQUHAQEEJjAkMCIGCCsGAQUFBzAC\nhhZodHRwOi8veDEuaS5sZW5jci5vcmcvMBMGA1UdIAQMMAowCAYGZ4EMAQIBMCcG\nA1UdHwQgMB4wHKAaoBiGFmh0dHA6Ly94MS5jLmxlbmNyLm9yZy8wDQYJKoZIhvcN\nAQELBQADggIBAI910AnPanZIZTKS3rVEyIV29BWEjAK/duuz8eL5boSoVpHhkkv3\n4eoAeEiPdZLj5EZ7G2ArIK+gzhTlRQ1q4FKGpPPaFBSpqV/xbUb5UlAXQOnkHn3m\nFVj+qYv87/WeY+Bm4sN3Ox8BhyaU7UAQ3LeZ7N1X01xxQe4wIAAE3JVLUCiHmZL+\nqoCUtgYIFPgcg350QMUIWgxPXNGEncT921ne7nluI02V8pLUmClqXOsCwULw+PVO\nZCB7qOMxxMBoCUeL2Ll4oMpOSr5pJCpLN3tRA2s6P1KLs9TSrVhOk+7LX28NMUlI\nusQ/nxLJID0RhAeFtPjyOCOscQBA53+NRjSCak7P4A5jX7ppmkcJECL+S0i3kXVU\ny5Me5BbrU8973jZNv/ax6+ZK6TM8jWmimL6of6OrX7ZU6E2WqazzsFrLG3o2kySb\nzlhSgJ81Cl4tv3SbYiYXnJExKQvzf83DYotox3f0fwv7xln1A2ZLplCb0O+l/AK0\nYE0DS2FPxSAHi0iwMfW2nNHJrXcY3LLHD77gRgje4Eveubi2xxa+Nmk/hmhLdIET\niVDFanoCrMVIpQ59XWHkzdFmoHXHBV7oibVjGSO7ULSQ7MJ1Nz51phuDJSgAIU7A\n0zrLnOrAj/dfrlEWRhCvAgbuwLZX1A2sjNjXoPOHbsPiy+lO1KF8/XY7\n-----END CERTIFICATE-----\n','RSA2048'),(20,'autoops.com.cn',2,'15463131',0,'cert-autoops.com.cn-1768634453','digicert-free-1-free',3,'DNS','','-----BEGIN CERTIFICATE-----\nMIIGCzCCBPOgAwIBAgIQDy+JplihnBOQxsS7Ho42WjANBgkqhkiG9w0BAQsFADBu\nMQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3\nd3cuZGlnaWNlcnQuY29tMS0wKwYDVQQDEyRFbmNyeXB0aW9uIEV2ZXJ5d2hlcmUg\nRFYgVExTIENBIC0gRzIwHhcNMjYwMTE3MDAwMDAwWhcNMjYwNDE2MjM1OTU5WjAZ\nMRcwFQYDVQQDEw5hdXRvb3BzLmNvbS5jbjCCASIwDQYJKoZIhvcNAQEBBQADggEP\nADCCAQoCggEBALR9nqbHm5QkSNnmL60Iiqt3MQ3SvyY4XXMSDa4XAIOFdY4V9Xt7\nNVPfQVTy9viWAudYn21XP6wga8gGxdpYH0GZVdEXYh/Ddci7mPlzMvAd/P95IxPv\n4Rob9UFnqYOtDJHrfYmY0YecYsIQ2TIboINbKhFvLKnsIipF7una6Fg7DHu7OlRB\n+EJXe4syY2R1aMEbMwTYBAd0YJeUNni4QTGyxe3/3zSYYyPmVt4XX2zMQmm8bTbC\nH2azj8K99MoB5Aj38AskVDnCV8qYxwZ945ZXHdCUGnw2Y1RdzyUuFpXAtEOizUa0\nc2DMokMSlbDCmysBY5tMVPWsHmvyK5eCAtkCAwEAAaOCAvgwggL0MB8GA1UdIwQY\nMBaAFHjfkZBf7t6s9sV169VMVVPvJEq2MB0GA1UdDgQWBBSC9bs1/WjHbxftJulV\nizEXgD/2PjAtBgNVHREEJjAkgg5hdXRvb3BzLmNvbS5jboISd3d3LmF1dG9vcHMu\nY29tLmNuMD4GA1UdIAQ3MDUwMwYGZ4EMAQIBMCkwJwYIKwYBBQUHAgEWG2h0dHA6\nLy93d3cuZGlnaWNlcnQuY29tL0NQUzAOBgNVHQ8BAf8EBAMCBaAwHQYDVR0lBBYw\nFAYIKwYBBQUHAwEGCCsGAQUFBwMCMIGABggrBgEFBQcBAQR0MHIwJAYIKwYBBQUH\nMAGGGGh0dHA6Ly9vY3NwLmRpZ2ljZXJ0LmNvbTBKBggrBgEFBQcwAoY+aHR0cDov\nL2NhY2VydHMuZGlnaWNlcnQuY29tL0VuY3J5cHRpb25FdmVyeXdoZXJlRFZUTFND\nQS1HMi5jcnQwDAYDVR0TAQH/BAIwADCCAYEGCisGAQQB1nkCBAIEggFxBIIBbQFr\nAHcADleUvPOuqT4zGyyZB7P3kN+bwj1xMiXdIaklrGHFTiEAAAGbytj3awAABAMA\nSDBGAiEAzTnWkMQzgQM0uitWTlurPCJGUYDbhY8mO2KA+AXbgiYCIQDSDLE+tinp\nzCprJBWaxsvZgnk6p4DSi56Nplui9KbKjgB3ABaDLavwqSUPD/A6pUX/yL/II9CH\nS/YEKSf45x8zE/X6AAABm8rY944AAAQDAEgwRgIhAPxT4Lb3ttTm3IPl1pA4vGA2\nPTvMI2b7d6kWV2gVmel7AiEAngmuk88bY//x2+nClm4gI7n9BM2clrYQq5AyouCe\nQNkAdwBkEcRspBLsp4kcogIuALyrTygH1B41J6vq/tUDyX3N8AAAAZvK2PelAAAE\nAwBIMEYCIQD8n06EAhxjHOTqgYzI8gR+A/3SBYOR+zIdDlKVS9IYCgIhAJUlzSsy\nfoK0xvHWwPPCizksL+Uvt+w2RbrDztdIuquTMA0GCSqGSIb3DQEBCwUAA4IBAQBK\nD7tvUywWD0DFG9HSUnQ5aJ1cDdPO7rWqQSReISSyr8Vvx/vj+niW3BL/gG0Rzqvc\nm/ne4lQE+GNmOV0UVDGgmtTQUVp1MZlaaPNkebXUsGAxzz2+Y9fle5eSzGBFw0sD\nc1kEq+hikniCx2to8Vmi78HGWa51QBEpukz2Pm7GvlN4fe4k4Zc0tp2Ok8DR+1hh\nZAKO55D1IIfGS0cttffG+Xf5qWqZ9qSUSYBVCsSRmLSQ0VUYluOhxf4iCxE295WU\nU+sZu3P/CR+MxEk+F/uDiA0JoLnsaqRCahHLhMTAy7XeqQs7bE1ghOcFiXpGBNau\nAqSsMhm2CVsPVTPuwYBY\n-----END CERTIFICATE-----\n-----BEGIN CERTIFICATE-----\nMIIEqjCCA5KgAwIBAgIQDeD/te5iy2EQn2CMnO1e0zANBgkqhkiG9w0BAQsFADBh\nMQswCQYDVQQGEwJVUzEVMBMGA1UEChMMRGlnaUNlcnQgSW5jMRkwFwYDVQQLExB3\nd3cuZGlnaWNlcnQuY29tMSAwHgYDVQQDExdEaWdpQ2VydCBHbG9iYWwgUm9vdCBH\nMjAeFw0xNzExMjcxMjQ2NDBaFw0yNzExMjcxMjQ2NDBaMG4xCzAJBgNVBAYTAlVT\nMRUwEwYDVQQKEwxEaWdpQ2VydCBJbmMxGTAXBgNVBAsTEHd3dy5kaWdpY2VydC5j\nb20xLTArBgNVBAMTJEVuY3J5cHRpb24gRXZlcnl3aGVyZSBEViBUTFMgQ0EgLSBH\nMjCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAO8Uf46i/nr7pkgTDqnE\neSIfCFqvPnUq3aF1tMJ5hh9MnO6Lmt5UdHfBGwC9Si+XjK12cjZgxObsL6Rg1njv\nNhAMJ4JunN0JGGRJGSevbJsA3sc68nbPQzuKp5Jc8vpryp2mts38pSCXorPR+sch\nQisKA7OSQ1MjcFN0d7tbrceWFNbzgL2csJVQeogOBGSe/KZEIZw6gXLKeFe7mupn\nNYJROi2iC11+HuF79iAttMc32Cv6UOxixY/3ZV+LzpLnklFq98XORgwkIJL1HuvP\nha8yvb+W6JislZJL+HLFtidoxmI7Qm3ZyIV66W533DsGFimFJkz3y0GeHWuSVMbI\nlfsCAwEAAaOCAU8wggFLMB0GA1UdDgQWBBR435GQX+7erPbFdevVTFVT7yRKtjAf\nBgNVHSMEGDAWgBROIlQgGJXm427mD/r6uRLtBhePOTAOBgNVHQ8BAf8EBAMCAYYw\nHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMBIGA1UdEwEB/wQIMAYBAf8C\nAQAwNAYIKwYBBQUHAQEEKDAmMCQGCCsGAQUFBzABhhhodHRwOi8vb2NzcC5kaWdp\nY2VydC5jb20wQgYDVR0fBDswOTA3oDWgM4YxaHR0cDovL2NybDMuZGlnaWNlcnQu\nY29tL0RpZ2lDZXJ0R2xvYmFsUm9vdEcyLmNybDBMBgNVHSAERTBDMDcGCWCGSAGG\n/WwBAjAqMCgGCCsGAQUFBwIBFhxodHRwczovL3d3dy5kaWdpY2VydC5jb20vQ1BT\nMAgGBmeBDAECATANBgkqhkiG9w0BAQsFAAOCAQEAoBs1eCLKakLtVRPFRjBIJ9LJ\nL0s8ZWum8U8/1TMVkQMBn+CPb5xnCD0GSA6L/V0ZFrMNqBirrr5B241OesECvxIi\n98bZ90h9+q/X5eMyOD35f8YTaEMpdnQCnawIwiHx06/0BfiTj+b/XQih+mqt3ZXe\nxNCJqKexdiB2IWGSKcgahPacWkk/BAQFisKIFYEqHzV974S3FAz/8LIfD58xnsEN\nGfzyIDkH3JrwYZ8caPTf6ZX9M1GrISN8HnWTtdNCH2xEajRa/h9ZBXjUyFKQrGk2\nn2hcLrfZSbynEC/pSw/ET7H5nWwckjmAJ1l9fcnbqkU/pf6uMQmnfl0JQjJNSg==\n-----END CERTIFICATE-----\n','-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEAtH2epseblCRI2eYvrQiKq3cxDdK/JjhdcxINrhcAg4V1jhX1\ne3s1U99BVPL2+JYC51ifbVc/rCBryAbF2lgfQZlV0RdiH8N1yLuY+XMy8B38/3kj\nE+/hGhv1QWepg60Mket9iZjRh5xiwhDZMhugg1sqEW8sqewiKkXu6droWDsMe7s6\nVEH4Qld7izJjZHVowRszBNgEB3Rgl5Q2eLhBMbLF7f/fNJhjI+ZW3hdfbMxCabxt\nNsIfZrOPwr30ygHkCPfwCyRUOcJXypjHBn3jllcd0JQafDZjVF3PJS4WlcC0Q6LN\nRrRzYMyiQxKVsMKbKwFjm0xU9awea/Irl4IC2QIDAQABAoIBABXR9gMtsIZZWKFX\nLb5agQE3uKt2JfvF5q5JqXO5fitoUAzsZLAX9ht8qxZNBUSKRjGXyjZazOzPQph9\nESBTdPvmQLS1gfxbEftBir3limvzBCoSiAYaAvm4FY9hq+GqYxGlCZk/61zcsCfM\nB6YCJwaaV9yhtyGX1Y2MxNRqzZoU3ViV1MoALGT/ziepcjx8iqdQqsZAuqtQ/WMg\nYip7JGkvSnKKs8vV55qfeeaxl+CEmo287XhrFE8v1x7U929JU2vNgJTUB4pCEGob\nBaRBNsNmmpREAqowGqLZyQMIVWPAiBubreJiRTPvelVwNMmoBRSQjpV94owE5h8A\ndsivlt0CgYEA5X+xyxXNTkSfycEJtdctfrlI9ujBw8Y93P2kfquISNI5sTDaY0A0\nO5LQu9NMYl5iCr6t36PlKqQKS1QmoOZyLGR6G7HIBqrdaBLXEHKZj0b5zWY/Mzcp\nYCyIt8ceS1CNeR7xETR9RjxPFEbYO7VZF/aGSLmDPdjETYR3zlKUsR8CgYEAyVUt\nWi7A+EZpLH3ishPdBl82TI8AN6iWOpMgTeP06MGUK2TuvQs+qNu4XxmLBs+FrfKc\nMptoYwd6kWhJpLkMQNz5Pr0Nc8XV9GW8vtciwOtycjeFDIzac4jhiRNhPIR/+WLN\nidua3n2t2dJmapgotIZmpbiCXDqujnGWrPnCdQcCgYEAzN7ZUHn99Mi9+9UfXfzg\n2Cd/Qys1DheCUiz7Ezj9SESiy9LYoZFHD4dYXtOZdJbc3sqpXONqSE54Ahb1BN2z\nPn3To7Vo4srVK3rmEgC7pHOvu/Rypr2vafNXePd/y6v1HxLIEJvS2UgtL+SuU17o\nUcP4zCMPTb6cPrbqXHzLz/ECgYAA5BGZTo+oCZe5jK0OIKxftNgmqpNTq7lOdQw1\njAavaPvg3Ud5Kq2STbaZJSHKW+y5glNSngFuPsA6GH4K7o6xUsj/kBKZKozx0elt\n5sHOnLwa6Uf7cDwfOsHUu3M4s1VloxxJueELjv7P8fBmFFpJwjZ+ETzyb4gKNVJ3\nXXuQzQKBgBF7e4I3d4Z/3vQSWX7Sv68cFFbdE5y3GBwoggYijm8dVKJRl8hOiiDe\nR9A4HE653WMLstgKn2r5H7BcQKq1NeXRZbUcOGsINGqRH6/46YM0ksEwT6yLB3lR\nmf0RJGlRDN9+DXM54tBMGGEhDS+loAPJuqgnrsAhBnJdDJKS+r/N\n-----END RSA PRIVATE KEY-----','2026-01-17 08:00:00.000','2026-04-16 08:00:00.000',89,'','2026-01-17 15:20:53.046','2026-01-17 15:28:00.659','aliyun_cas','DigiCert','','RSA2048');
/*!40000 ALTER TABLE `monitor_ssl_cert` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `monitor_ssl_cert_deploy_log`
--

DROP TABLE IF EXISTS `monitor_ssl_cert_deploy_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `monitor_ssl_cert_deploy_log` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `cert_id` bigint unsigned NOT NULL COMMENT '证书ID',
  `domain` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '域名',
  `host_id` bigint unsigned NOT NULL COMMENT '目标主机ID',
  `host_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '主机名称',
  `deploy_path` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '部署路径',
  `status` int DEFAULT '1' COMMENT '部署状态:1->部署中,2->成功,3->失败',
  `backup_files` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '备份的文件列表(JSON)',
  `deploy_files` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '部署的文件列表(JSON)',
  `logs` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '执行日志(JSON)',
  `error_msg` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '错误信息',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `update_time` datetime DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_cert_id` (`cert_id`) USING BTREE,
  KEY `idx_host_id` (`host_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='SSL证书部署日志表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `monitor_ssl_cert_deploy_log`
--

LOCK TABLES `monitor_ssl_cert_deploy_log` WRITE;
/*!40000 ALTER TABLE `monitor_ssl_cert_deploy_log` DISABLE KEYS */;
INSERT INTO `monitor_ssl_cert_deploy_log` VALUES (1,6,'deviops.cn',525,'go-ops','/etc/nginx/ssl',3,'[]','[]','[]','获取SSH认证信息失败: Error 1054 (42S22): Unknown column \'auth\' in \'field list\'','2025-12-17 12:20:35','2025-12-17 12:20:35'),(2,6,'deviops.cn',525,'go-ops','/etc/nginx/ssl',2,'[]','[\"/etc/nginx/ssl/deviops.cn.pem\",\"/etc/nginx/ssl/deviops.cn.key\"]','[\"[2025-12-17 12:26:56] 开始部署证书到主机: go-ops (8.130.14.34)\",\"[2025-12-17 12:26:56] 检查部署目录: /etc/nginx/ssl\",\"[2025-12-17 12:26:59] 目录已就绪\",\"[2025-12-17 12:26:59] 检查旧证书文件\",\"[2025-12-17 12:27:02] 开始部署新证书\",\"[2025-12-17 12:27:03] 证书文件已部署: /etc/nginx/ssl/deviops.cn.pem\",\"[2025-12-17 12:27:04] 私钥文件已部署: /etc/nginx/ssl/deviops.cn.key\",\"[2025-12-17 12:27:05] 已设置文件权限为 600\",\"[2025-12-17 12:27:05] 开始重载Nginx\",\"[2025-12-17 12:27:07] Nginx重载成功\",\"[2025-12-17 12:27:07] Nginx输出: 2025/12/17 12:27:07 [notice] 356917#356917: signal process started\",\"[2025-12-17 12:27:07] 证书部署完成\"]','','2025-12-17 12:26:57','2025-12-17 12:27:07'),(3,1,'deviops.cn',525,'go-ops','/etc/nginx/ssl',2,'[\"/etc/nginx/ssl/deviops.cn.pem.20251217130203.bak\",\"/etc/nginx/ssl/deviops.cn.key.20251217130203.bak\"]','[\"/etc/nginx/ssl/deviops.cn.pem\",\"/etc/nginx/ssl/deviops.cn.key\"]','[\"[2025-12-17 13:02:00] 开始部署证书到主机: go-ops (8.130.14.34)\",\"[2025-12-17 13:02:00] 检查部署目录: /etc/nginx/ssl\",\"[2025-12-17 13:02:03] 目录已就绪\",\"[2025-12-17 13:02:03] 检查旧证书文件\",\"[2025-12-17 13:02:05] 已备份证书文件: /etc/nginx/ssl/deviops.cn.pem.20251217130203.bak\",\"[2025-12-17 13:02:08] 已备份密钥文件: /etc/nginx/ssl/deviops.cn.key.20251217130203.bak\",\"[2025-12-17 13:02:08] 开始部署新证书\",\"[2025-12-17 13:02:09] 证书文件已部署: /etc/nginx/ssl/deviops.cn.pem\",\"[2025-12-17 13:02:10] 私钥文件已部署: /etc/nginx/ssl/deviops.cn.key\",\"[2025-12-17 13:02:11] 已设置文件权限为 600\",\"[2025-12-17 13:02:11] 开始重载Nginx\",\"[2025-12-17 13:02:13] Nginx重载成功\",\"[2025-12-17 13:02:13] Nginx输出: 2025/12/17 13:02:13 [notice] 359465#359465: signal process started\",\"[2025-12-17 13:02:13] 证书部署完成\"]','','2025-12-17 13:02:00','2025-12-17 13:02:13'),(4,7,'deviops.cn',525,'go-ops','/etc/nginx/ssl',2,'[\"/etc/nginx/ssl/deviops.cn.pem.20251217140432.bak\",\"/etc/nginx/ssl/deviops.cn.key.20251217140432.bak\"]','[\"/etc/nginx/ssl/deviops.cn.pem\",\"/etc/nginx/ssl/deviops.cn.key\"]','[\"[2025-12-17 14:04:29] 开始部署证书到主机: go-ops (8.130.14.34)\",\"[2025-12-17 14:04:29] 检查部署目录: /etc/nginx/ssl\",\"[2025-12-17 14:04:32] 目录已就绪\",\"[2025-12-17 14:04:32] 检查旧证书文件\",\"[2025-12-17 14:04:34] 已备份证书文件: /etc/nginx/ssl/deviops.cn.pem.20251217140432.bak\",\"[2025-12-17 14:04:37] 已备份密钥文件: /etc/nginx/ssl/deviops.cn.key.20251217140432.bak\",\"[2025-12-17 14:04:37] 开始部署新证书\",\"[2025-12-17 14:04:37] 检测到证书已包含完整链\",\"[2025-12-17 14:04:38] 证书文件已部署: /etc/nginx/ssl/deviops.cn.pem\",\"[2025-12-17 14:04:39] 私钥文件已部署: /etc/nginx/ssl/deviops.cn.key\",\"[2025-12-17 14:04:41] 已设置文件权限为 600\",\"[2025-12-17 14:04:41] 开始重载Nginx\",\"[2025-12-17 14:04:42] Nginx重载成功\",\"[2025-12-17 14:04:42] Nginx输出: 2025/12/17 14:04:42 [notice] 363589#363589: signal process started\",\"[2025-12-17 14:04:42] 证书部署完成\"]','','2025-12-17 14:04:29','2025-12-17 14:04:42'),(5,8,'deviops.cn',525,'go-ops','/etc/nginx/ssl',2,'[\"/etc/nginx/ssl/deviops.cn.pem.20251217154631.bak\",\"/etc/nginx/ssl/deviops.cn.key.20251217154631.bak\"]','[\"/etc/nginx/ssl/deviops.cn.pem\",\"/etc/nginx/ssl/deviops.cn.key\"]','[\"[2025-12-17 15:46:28] 开始部署证书到主机: go-ops (8.130.14.34)\",\"[2025-12-17 15:46:28] 检查部署目录: /etc/nginx/ssl\",\"[2025-12-17 15:46:31] 目录已就绪\",\"[2025-12-17 15:46:31] 检查旧证书文件\",\"[2025-12-17 15:46:34] 已备份证书文件: /etc/nginx/ssl/deviops.cn.pem.20251217154631.bak\",\"[2025-12-17 15:46:36] 已备份密钥文件: /etc/nginx/ssl/deviops.cn.key.20251217154631.bak\",\"[2025-12-17 15:46:36] 开始部署新证书\",\"[2025-12-17 15:46:36] 检测到证书已包含完整链\",\"[2025-12-17 15:46:37] 证书文件已部署: /etc/nginx/ssl/deviops.cn.pem\",\"[2025-12-17 15:46:38] 私钥文件已部署: /etc/nginx/ssl/deviops.cn.key\",\"[2025-12-17 15:46:40] 已设置文件权限为 600\",\"[2025-12-17 15:46:40] 开始重载Nginx\",\"[2025-12-17 15:46:41] Nginx重载成功\",\"[2025-12-17 15:46:41] Nginx输出: 2025/12/17 15:46:41 [notice] 369478#369478: signal process started\",\"[2025-12-17 15:46:41] 证书部署完成\"]','','2025-12-17 15:46:29','2025-12-17 15:46:41'),(6,10,'deviops.cn',525,'go-ops','/etc/nginx/ssl',2,'[\"/etc/nginx/ssl/deviops.cn.pem.20251220153616.bak\",\"/etc/nginx/ssl/deviops.cn.key.20251220153616.bak\"]','[\"/etc/nginx/ssl/deviops.cn.pem\",\"/etc/nginx/ssl/deviops.cn.key\"]','[\"[2025-12-20 15:36:15] 开始部署证书到主机: go-ops (8.130.14.34)\",\"[2025-12-20 15:36:15] 检查部署目录: /etc/nginx/ssl\",\"[2025-12-20 15:36:16] 目录已就绪\",\"[2025-12-20 15:36:16] 检查旧证书文件\",\"[2025-12-20 15:36:19] 已备份证书文件: /etc/nginx/ssl/deviops.cn.pem.20251220153616.bak\",\"[2025-12-20 15:36:21] 已备份密钥文件: /etc/nginx/ssl/deviops.cn.key.20251220153616.bak\",\"[2025-12-20 15:36:21] 开始部署新证书\",\"[2025-12-20 15:36:21] 检测到证书已包含完整链\",\"[2025-12-20 15:36:22] 证书文件已部署: /etc/nginx/ssl/deviops.cn.pem\",\"[2025-12-20 15:36:23] 私钥文件已部署: /etc/nginx/ssl/deviops.cn.key\",\"[2025-12-20 15:36:25] 已设置文件权限为 600\",\"[2025-12-20 15:36:25] 开始重载Nginx\",\"[2025-12-20 15:36:26] Nginx重载成功\",\"[2025-12-20 15:36:26] Nginx输出: 2025/12/20 15:36:26 [notice] 451483#451483: signal process started\",\"[2025-12-20 15:36:26] 证书部署完成\"]','','2025-12-20 15:36:15','2025-12-20 15:36:26'),(7,17,'deviops.cn',1,'阿里云-ops','/etc/nginx/ssl',2,'[\"/etc/nginx/ssl/deviops.cn.pem.20260117143010.bak\",\"/etc/nginx/ssl/deviops.cn.key.20260117143010.bak\"]','[\"/etc/nginx/ssl/deviops.cn.pem\",\"/etc/nginx/ssl/deviops.cn.key\"]','[\"[2026-01-17 14:30:08] 开始部署证书到主机: 阿里云-ops (8.130.14.34)\",\"[2026-01-17 14:30:08] 检查部署目录: /etc/nginx/ssl\",\"[2026-01-17 14:30:10] 目录已就绪\",\"[2026-01-17 14:30:10] 检查旧证书文件\",\"[2026-01-17 14:30:13] 已备份证书文件: /etc/nginx/ssl/deviops.cn.pem.20260117143010.bak\",\"[2026-01-17 14:30:15] 已备份密钥文件: /etc/nginx/ssl/deviops.cn.key.20260117143010.bak\",\"[2026-01-17 14:30:15] 开始部署新证书\",\"[2026-01-17 14:30:15] 检测到证书已包含完整链\",\"[2026-01-17 14:30:16] 证书文件已部署: /etc/nginx/ssl/deviops.cn.pem\",\"[2026-01-17 14:30:17] 私钥文件已部署: /etc/nginx/ssl/deviops.cn.key\",\"[2026-01-17 14:30:18] 已设置文件权限为 600\",\"[2026-01-17 14:30:18] 开始重载Nginx\",\"[2026-01-17 14:30:19] Nginx重载成功\",\"[2026-01-17 14:30:19] Nginx输出: 2026/01/17 14:30:19 [notice] 486589#486589: signal process started\",\"[2026-01-17 14:30:19] 证书部署完成\"]','','2026-01-17 14:30:08','2026-01-17 14:30:19');
/*!40000 ALTER TABLE `monitor_ssl_cert_deploy_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `monitor_webhook_log`
--

DROP TABLE IF EXISTS `monitor_webhook_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `monitor_webhook_log` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `source` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '来源标识',
  `title` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '告警标题',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '告警内容',
  `level` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '告警级别',
  `tags` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '标签(JSON数组字符串)',
  `extra` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '额外数据(JSON字符串)',
  `notify_robot_ids` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '通知的机器人ID列表(逗号分隔)',
  `status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT 'success' COMMENT '处理状态(success/failed/partial)',
  `error_msg` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '错误信息',
  `notify_count` bigint DEFAULT '0' COMMENT '通知数量',
  `success_count` bigint DEFAULT '0' COMMENT '成功数量',
  `failed_count` bigint DEFAULT '0' COMMENT '失败数量',
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `monitor_webhook_log`
--

LOCK TABLES `monitor_webhook_log` WRITE;
/*!40000 ALTER TABLE `monitor_webhook_log` DISABLE KEYS */;
INSERT INTO `monitor_webhook_log` VALUES (1,'zabbix','CPU告警','CPU使用率过高','critical','','','','success','',1,0,0,'2026-01-08 18:08:01.934'),(2,'prometheus','服务宕机','API服务无响应','critical','','','','success','',1,0,0,'2026-01-08 18:08:28.929'),(3,'deployment','发布完成: test','发布标题: test\n发布环境: test\n服务名称: s3-web\n发布人员: 管理员\n发布时间: 2026-01-19 17:16:37\n运行时长: 22秒\n发布状态: 成功','info','','','[1]','success','',1,1,0,'2026-01-19 17:18:19.556'),(4,'deployment','发布完成: s3-登录功能迭代上线2.0','发布标题: s3-登录功能迭代上线2.0\n发布环境: test\n服务名称: s3-api, s3-web\n发布人员: 管理员\n发布时间: 2026-01-19 17:19:49\n运行时长: 132秒\n发布状态: 成功','info','','','[1]','success','',1,1,0,'2026-01-19 17:22:06.897');
/*!40000 ALTER TABLE `monitor_webhook_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `monitor_webhook_notify_log`
--

DROP TABLE IF EXISTS `monitor_webhook_notify_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `monitor_webhook_notify_log` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `webhook_log_id` bigint NOT NULL COMMENT 'Webhook日志ID',
  `robot_id` bigint NOT NULL COMMENT '机器人ID',
  `robot_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '机器人名称',
  `robot_type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '机器人类型',
  `status` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '发送状态(success/failed)',
  `error_msg` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '错误信息',
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_monitor_webhook_notify_log_webhook_log_id` (`webhook_log_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `monitor_webhook_notify_log`
--

LOCK TABLES `monitor_webhook_notify_log` WRITE;
/*!40000 ALTER TABLE `monitor_webhook_notify_log` DISABLE KEYS */;
INSERT INTO `monitor_webhook_notify_log` VALUES (1,1,1,'devops飞书机器人','feishu','success','','2026-01-08 18:08:02.660'),(2,2,1,'devops飞书机器人','feishu','success','','2026-01-08 18:08:29.570'),(3,3,1,'devops飞书机器人','feishu','success','','2026-01-19 17:18:20.363'),(4,4,1,'devops飞书机器人','feishu','success','','2026-01-19 17:22:07.907');
/*!40000 ALTER TABLE `monitor_webhook_notify_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `notify_robot`
--

DROP TABLE IF EXISTS `notify_robot`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `notify_robot` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '机器人名称',
  `type` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '机器人类型(dingtalk/feishu/wechat)',
  `webhook` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT 'webhook地址',
  `secret` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '加签密钥',
  `status` tinyint DEFAULT '1' COMMENT '状态(1:启用,0:禁用)',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '备注说明',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `notify_robot`
--

LOCK TABLES `notify_robot` WRITE;
/*!40000 ALTER TABLE `notify_robot` DISABLE KEYS */;
/*!40000 ALTER TABLE `notify_robot` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `prompt_template`
--

DROP TABLE IF EXISTS `prompt_template`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `prompt_template` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键''',
  `name` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''模板名称''',
  `category` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''分类:inspection/incident/prediction/general''',
  `scene` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''场景:summary/analysis/suggestion/fix''',
  `template` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''模板内容(支持变量替换)''',
  `variables` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''变量列表(JSON数组)''',
  `system_prompt` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''系统提示词''',
  `temperature` decimal(3,2) DEFAULT '0.70' COMMENT '''温度参数''',
  `max_tokens` bigint DEFAULT '2000' COMMENT '''最大Token数''',
  `model_id` bigint unsigned DEFAULT NULL COMMENT '''默认模型ID''',
  `enabled` bigint DEFAULT '1' COMMENT '''是否启用:1=启用,2=禁用''',
  `create_time` datetime(3) NOT NULL COMMENT '''创建时间''',
  `update_time` datetime(3) NOT NULL COMMENT '''更新时间''',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `idx_prompt_template_name` (`name`) USING BTREE,
  KEY `idx_prompt_template_category` (`category`) USING BTREE,
  KEY `idx_prompt_template_scene` (`scene`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `prompt_template`
--

LOCK TABLES `prompt_template` WRITE;
/*!40000 ALTER TABLE `prompt_template` DISABLE KEYS */;
INSERT INTO `prompt_template` VALUES (1,'inspection_summary','inspection','summary','你是一个专业的运维工程师，请根据以下巡检结果生成总结报告：\n\n## 巡检信息\n- 主机: {{hostName}} ({{hostIp}})\n- 巡检时间: {{inspectionTime}}\n- 总检查项: {{totalChecks}}\n- 通过项: {{passedChecks}}\n- 失败项: {{failedChecks}}\n- 警告项: {{warnChecks}}\n\n## 检查结果\n{{checkResults}}\n\n请生成一份简洁的总结，包括：\n1. 整体健康状态评估\n2. 主要问题汇总\n3. 风险等级评估\n4. 关键建议（不超过5条）\n\n要求：\n- 简洁专业，重点突出\n- 按优先级排序\n- 给出可执行的建议','[\"hostName\",\"hostIp\",\"inspectionTime\",\"totalChecks\",\"passedChecks\",\"failedChecks\",\"warnChecks\",\"checkResults\"]','你是一个专业的运维工程师，擅长分析系统健康状况和提供运维建议。',0.70,1000,0,1,'2026-01-09 19:16:12.751','2026-01-09 19:16:12.751'),(2,'incident_analysis','incident','analysis','你是一个专业的运维故障分析专家，请根据以下信息进行根因分析：\n\n## 故障基本信息\n- 主机: {{hostName}} ({{hostIp}})\n- 故障类别: {{category}}\n- 检测时间: {{detectTime}}\n- 症状: {{symptoms}}\n\n## 证据数据\n{{evidenceData}}\n\n## 分析要求\n1. 分析故障的根本原因\n2. 列出支持该结论的关键证据\n3. 评估故障的影响范围\n4. 给出详细的解决方案\n5. 提供可执行的修复脚本（如适用）\n6. 评估自动修复的风险等级（0-3）\n\n请按照以下JSON格式返回分析结果：\n{\n    \"rootCause\": \"根因描述\",\n    \"keyEvidence\": [\"证据1\", \"证据2\"],\n    \"impact\": \"影响范围描述\",\n    \"solution\": \"解决方案描述\",\n    \"fixScript\": \"#!/bin/bash\\n修复脚本\",\n    \"riskLevel\": 0-3,\n    \"riskReason\": \"风险评估说明\"\n}','[\"hostName\",\"hostIp\",\"category\",\"detectTime\",\"symptoms\",\"evidenceData\"]','你是一个专业的运维故障分析专家，擅长快速定位问题根因并提供解决方案。',0.50,2000,0,1,'2026-01-09 19:16:13.063','2026-01-09 19:16:13.063'),(3,'prediction_suggestion','prediction','suggestion','你是一个专业的运维容量规划专家，请根据以下信息给出优化建议：\n\n## 资源信息\n- 目标: {{targetName}}\n- 指标: {{metricName}} ({{metricType}})\n- 当前值: {{currentValue}}%\n- 阈值: {{threshold}}%\n\n## 趋势分析\n- 趋势类型: {{trendType}}\n- 变化率: {{changeRate}}%/天\n- 预测剩余天数: {{daysLeft}}天\n- 置信度: {{confidence}}%\n\n## 预测数据\n未来7天预测值: {{predictions}}\n\n## 要求\n请给出以下建议：\n1. 当前资源状态评估\n2. 风险分析和影响评估\n3. 具体的优化措施（优先级排序）\n4. 扩容建议（如需要）\n5. 成本优化建议\n6. 长期优化策略\n\n请用简洁、专业的语言回答，重点突出可执行的建议。','[\"targetName\",\"metricName\",\"metricType\",\"currentValue\",\"threshold\",\"trendType\",\"changeRate\",\"daysLeft\",\"confidence\",\"predictions\"]','你是一个专业的运维容量规划专家，擅长资源趋势分析和容量规划。',0.70,1500,0,1,'2026-01-09 19:16:13.383','2026-01-09 19:16:13.383');
/*!40000 ALTER TABLE `prompt_template` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `quick_deployment_tasks`
--

DROP TABLE IF EXISTS `quick_deployment_tasks`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `quick_deployment_tasks` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `deployment_id` bigint unsigned NOT NULL,
  `app_id` bigint unsigned NOT NULL,
  `app_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `app_code` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `environment` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `jenkins_env_id` bigint unsigned NOT NULL,
  `jenkins_job_url` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `build_number` bigint DEFAULT NULL,
  `status` tinyint DEFAULT '1',
  `execute_order` bigint NOT NULL,
  `start_time` datetime(3) DEFAULT NULL,
  `end_time` datetime(3) DEFAULT NULL,
  `duration` bigint DEFAULT NULL,
  `error_message` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
  `log_url` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `parameters` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_quick_deployment_tasks_deployment_id` (`deployment_id`) USING BTREE,
  KEY `fk_quick_deployment_tasks_application` (`app_id`) USING BTREE,
  KEY `fk_quick_deployment_tasks_jenkins_env` (`jenkins_env_id`) USING BTREE,
  CONSTRAINT `fk_quick_deployment_tasks_application` FOREIGN KEY (`app_id`) REFERENCES `app_application` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT,
  CONSTRAINT `fk_quick_deployment_tasks_jenkins_env` FOREIGN KEY (`jenkins_env_id`) REFERENCES `app_jenkins_env` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT,
  CONSTRAINT `fk_quick_deployments_tasks` FOREIGN KEY (`deployment_id`) REFERENCES `quick_deployments` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=52 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `quick_deployment_tasks`
--

LOCK TABLES `quick_deployment_tasks` WRITE;
/*!40000 ALTER TABLE `quick_deployment_tasks` DISABLE KEYS */;
INSERT INTO `quick_deployment_tasks` VALUES (50,39,20,'s3-api','s3-api','test',55,'',22,3,1,'2026-01-19 17:20:12.789','2026-01-19 17:22:05.410',112,'','http://180.76.231.65:8080/job/s3-api/22/console','2026-01-19 17:19:50.088','2026-01-19 17:22:05.460','{\"commit_id\":\"123456\",\"env\":\"test\"}'),(51,39,21,'s3-web','s3-web','test',58,'',13,3,2,'2026-01-19 17:20:22.426','2026-01-19 17:20:42.601',20,'','http://180.76.231.65:8080/job/s3-web/13/console','2026-01-19 17:19:50.510','2026-01-19 17:20:42.652','{\"commit_id\":\"123456\"}');
/*!40000 ALTER TABLE `quick_deployment_tasks` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `quick_deployments`
--

DROP TABLE IF EXISTS `quick_deployments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `quick_deployments` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `business_group_id` bigint unsigned NOT NULL,
  `business_dept_id` bigint unsigned NOT NULL,
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
  `status` tinyint DEFAULT '1',
  `task_count` bigint NOT NULL DEFAULT '0',
  `creator_id` bigint unsigned NOT NULL,
  `creator_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `start_time` datetime(3) DEFAULT NULL,
  `end_time` datetime(3) DEFAULT NULL,
  `duration` bigint DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `execution_mode` tinyint DEFAULT '1',
  `notice` tinyint DEFAULT '0',
  `robot_id` bigint unsigned DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=40 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `quick_deployments`
--

LOCK TABLES `quick_deployments` WRITE;
/*!40000 ALTER TABLE `quick_deployments` DISABLE KEYS */;
INSERT INTO `quick_deployments` VALUES (39,'s3-登录功能迭代上线2.0',60,2,'s3-登录功能迭代上线2.0',3,2,1,'管理员','2026-01-19 17:20:12.265','2026-01-19 17:22:06.168',132,'2026-01-19 17:19:49.672','2026-01-19 17:22:06.332',1,1,1);
/*!40000 ALTER TABLE `quick_deployments` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `redis_instance`
--

DROP TABLE IF EXISTS `redis_instance`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `redis_instance` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键''',
  `code` varchar(36) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''实例编码''',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''实例名称''',
  `mode` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''模式:standalone,cluster,sentinel''',
  `host` varchar(300) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''主机: standalone为host:port, cluster为逗号分隔, sentinel为master=hosts''',
  `port` bigint DEFAULT '0' COMMENT '''端口(standalone可用)''',
  `db` bigint DEFAULT '0' COMMENT '''默认库号''',
  `username` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''用户名(可选)''',
  `password` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''密码(加密)''',
  `redis_node_password` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT '' COMMENT '''节点密码(仅sentinel)''',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''备注''',
  `ssh_tunnel_machine_id` bigint unsigned DEFAULT '0' COMMENT '''SSH隧道机器ID''',
  `status` bigint DEFAULT '1' COMMENT '''状态:1->启用,2->禁用''',
  `create_time` datetime(3) NOT NULL COMMENT '''创建时间''',
  `update_time` datetime(3) DEFAULT NULL COMMENT '''更新时间''',
  `creator` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''创建人''',
  `creator_id` bigint unsigned DEFAULT NULL COMMENT '''创建人ID''',
  `modifier` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''修改人''',
  `modifier_id` bigint unsigned DEFAULT NULL COMMENT '''修改人ID''',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `idx_redis_instance_code` (`code`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `redis_instance`
--

LOCK TABLES `redis_instance` WRITE;
/*!40000 ALTER TABLE `redis_instance` DISABLE KEYS */;
/*!40000 ALTER TABLE `redis_instance` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_activity_log`
--

DROP TABLE IF EXISTS `sys_activity_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_activity_log` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `activity_type` bigint NOT NULL COMMENT '''动态类型：1=密钥同步，2=域名检查，3=服务器巡检，4=定时任务，5=其他''',
  `title` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''动态标题''',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''动态详细内容''',
  `status` bigint NOT NULL DEFAULT '1' COMMENT '''状态：1=成功，2=失败，3=部分成功''',
  `related_id` bigint unsigned DEFAULT NULL COMMENT '''关联ID（如同步任务ID、域名ID等）''',
  `summary` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''摘要信息''',
  `duration` bigint DEFAULT NULL COMMENT '''执行耗时（秒）''',
  `create_time` datetime(3) NOT NULL COMMENT '''创建时间''',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=61 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_activity_log`
--

LOCK TABLES `sys_activity_log` WRITE;
/*!40000 ALTER TABLE `sys_activity_log` DISABLE KEYS */;
INSERT INTO `sys_activity_log` VALUES (1,1,'阿里云主机同步','[2025-12-26 23:31:00] 开始同步\n\n未知云厂商(13):\n未配置密钥\n\n[2025-12-26 23:31:00] 同步完成，耗时: 183.355564ms\n',1,6,'同步1个云厂商类型的密钥',183,'2025-12-26 23:31:00.555'),(2,1,'cmdb云资产同步','[2025-12-27 00:09:00] 开始同步\n\n华为云:\n  - 华为云: 暂未实现\n\n百度云:\n  - 百度云: 暂未实现\n\n阿里云:\n  - 张帆-阿里云-密钥[阿里云]: 获取1台，新增0台，跳过1台，失败0台\n\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n本次同步 0 台新增主机，3 个密钥，共计耗时 1m0.2642405s\n[2025-12-27 00:10:00] 同步完成\n',1,7,'同步3个密钥，新增0台主机',60264,'2025-12-27 00:10:07.534'),(3,1,'cmdb云资产同步','[2025-12-27 10:20:00] 开始同步\n\n阿里云:\n  - 张帆-阿里云-密钥[阿里云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n\n华为云:\n  - 华为云: 暂未实现\n\n百度云:\n  - 百度云: 暂未实现\n\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n本次同步 0 台新增主机，3 个密钥，共计耗时 26.124756443s\n[2025-12-27 10:20:26] 同步完成\n',1,7,'同步3个密钥，新增0台主机',26124,'2025-12-27 10:20:26.466'),(4,1,'cmdb云资产同步','[2025-12-27 10:34:00] 开始同步\n\n阿里云:\n  - 张帆-阿里云-密钥[阿里云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n\n华为云:\n  - 张帆-华为云-密钥[华为云]: 获取1台，新增1台，跳过0台，失败0台 [cn-south-1:1台 ]\n\n百度云:\n  - 张帆-百度云-密钥[百度云]: 获取1台，新增1台，跳过0台，失败0台 [bj:1台 ]\n\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n本次同步 2 台新增主机，3 个密钥，共计耗时 54.232414686s\n[2025-12-27 10:34:54] 同步完成\n',1,7,'同步3个密钥，新增2台主机',54232,'2025-12-27 10:34:54.566'),(5,2,'域名定时巡检完成','巡检结果：\n- 巡检域名总数：17\n- 即将过期域名数：1（SSL证书剩余天数<30天）\n- 执行耗时：1.98秒',1,0,'本次巡检 17 个域名，1 个域名即将过期',1,'2025-12-27 22:19:42.694'),(6,1,'云资产同步','[2025-12-27 22:35:00] 开始同步\n\n华为云:\n  - 张帆-华为云-密钥[华为云]: 获取1台，新增0台，跳过1台，失败0台 [cn-south-1:1台 ]\n\n百度云:\n  - 张帆-百度云-密钥[百度云]: 获取1台，新增0台，跳过1台，失败0台 [bj:1台 ]\n\n阿里云:\n  - 张帆-阿里云-密钥[阿里云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n本次同步 0 台新增主机，3 个密钥，共计耗时 1m14.130908455s\n[2025-12-27 22:36:14] 同步完成\n',1,7,'同步3个密钥，共计新增0台主机，耗时74.1秒',74,'2025-12-27 22:36:14.484'),(7,2,'域名定时巡检完成','巡检结果：\n- 巡检域名总数：17\n- 即将过期域名数：1（SSL证书剩余天数<30天）\n- 执行耗时：4.88秒',1,0,'本次巡检 17 个域名，1 个域名即将过期',4,'2025-12-27 22:37:01.761'),(8,2,'域名定时巡检完成','巡检结果：\n- 巡检域名总数：17\n- 即将过期域名数：1（SSL证书剩余天数<30天）\n- 执行耗时：1.68秒',1,0,'本次巡检 17 个域名，1 个域名即将过期',1,'2025-12-27 22:39:58.552'),(9,2,'域名定时巡检完成','巡检结果：\n- 巡检域名总数：17\n- 即将过期域名数：1（SSL证书剩余天数<30天）\n- 执行耗时：3.01秒',1,0,'本次巡检 17 个域名，1 个域名即将过期',3,'2025-12-27 22:42:59.886'),(10,2,'域名定时巡检完成','巡检结果：\n- 巡检域名总数：17\n- 即将过期域名数：1（SSL证书剩余天数<30天）\n- 执行耗时：3.00秒',1,0,'本次巡检 17 个域名，1 个域名即将过期',2,'2025-12-27 23:24:02.998'),(11,1,'阿里云主机同步','[2025-12-27 23:31:00] 开始同步\n\n阿里云:\n  - 张帆-阿里云-密钥[阿里云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n本次同步 0 台新增主机，1 个密钥，共计耗时 44.119109395s\n[2025-12-27 23:31:44] 同步完成\n',1,6,'同步1个密钥，共计新增0台主机，耗时44.1秒',44,'2025-12-27 23:31:44.443'),(12,1,'devops云主机批量同步','[2025-12-28 01:53:00] 开始同步\n\nAWS云:\n  - AWS云: 暂未实现\n\n阿里云:\n  - test[阿里云]: 获取379台，新增379台，跳过0台，失败0台 [cn-shenzhen:2台 cn-hongkong:3台 us-east-1:2台 cn-qingdao:25台 cn-beijing:6台 cn-hangzhou:341台 ]\n  - 张三-阿里云-密钥[阿里云]: 获取1台，新增1台，跳过0台，失败0台 [cn-wulanchabu:1台 ]\n\n腾讯云:\n  - 腾讯云批量主机[腾讯云]: 获取82台，新增82台，跳过0台，失败0台 [ap-singapore:20台 na-siliconvalley:17台 ap-guangzhou:43台 ap-hongkong:1台 ap-tokyo:1台 ]\n\n华为云:\n  - 张三-华为云-密钥[华为云]: 获取1台，新增1台，跳过0台，失败0台 [cn-south-1:1台 ]\n\n百度云:\n  - 张三-百度云-密钥[百度云]: 获取1台，新增1台，跳过0台，失败0台 [bj:1台 ]\n本次同步 464 台新增主机，6 个密钥，共计耗时 4m9.844549966s\n[2025-12-28 01:57:09] 同步完成\n',1,8,'同步6个密钥，共计新增464台主机，耗时249.8秒',249,'2025-12-28 01:57:10.160'),(13,1,'手动同步任务: 阿里云主机同步','[2025-12-28 02:22:56] 手动触发同步\n\nAWS云:\n  - AWS云: 暂未实现\n\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n本次同步 0 台新增主机，1 个密钥，共计耗时 169.653907ms\n[2025-12-28 02:22:56] 同步完成\n',1,6,'手动触发同步1个密钥，新增0台主机',169,'2025-12-28 02:22:57.273'),(14,1,'手动同步任务: 阿里云主机同步','[2025-12-28 02:48:06] 手动触发同步\n\nAWS云:\n  - AWS云: 暂未实现\n\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n本次同步 0 台新增主机，1 个密钥，共计耗时 187.527522ms\n[2025-12-28 02:48:06] 同步完成\n',1,6,'手动触发同步1个密钥，新增0台主机',187,'2025-12-28 02:48:07.259'),(15,1,'手动同步任务: 阿里云主机同步','[2025-12-28 02:51:52] 手动触发同步\n\nAWS云:\n  - aws主机同步密钥[AWS云]: 获取31台，新增31台，跳过0台，失败0台\n\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n本次同步 31 台新增主机，1 个密钥，共计耗时 1m3.790532236s\n[2025-12-28 02:52:56] 同步完成\n',1,6,'手动触发同步1个密钥，新增31台主机',63790,'2025-12-28 02:52:56.442'),(16,1,'云主机同步','[2025-12-28 13:40:00] 开始同步\n\n阿里云:\n  - test[阿里云]: 获取377台，新增377台，跳过0台，失败0台 [cn-hongkong:3台 us-east-1:2台 cn-qingdao:25台 cn-beijing:6台 cn-hangzhou:341台 ]\n  - 张三-阿里云-密钥[阿里云]: 获取1台，新增1台，跳过0台，失败0台 [cn-wulanchabu:1台 ]\n\n腾讯云:\n  - 腾讯云批量主机[腾讯云]: 获取82台，新增82台，跳过0台，失败0台 [ap-guangzhou:43台 ap-hongkong:1台 ap-tokyo:1台 ap-singapore:20台 na-siliconvalley:17台 ]\n\nAWS云:\n  - aws主机同步密钥[AWS云]: 获取31台，新增31台，跳过0台，失败0台 [ap-southeast-1:5台 us-west-2:26台 ]\n\n百度云:\n  - 张三-百度云-密钥[百度云]: 获取1台，新增1台，跳过0台，失败0台 [bj:1台 ]\n\n华为云:\n  - 张三-华为云-密钥[华为云]: 获取1台，新增1台，跳过0台，失败0台 [cn-south-1:1台 ]\n本次同步 493 台新增主机，6 个密钥，共计耗时 5m12.033347678s\n[2025-12-28 13:45:12] 同步完成\n',1,9,'同步6个密钥，共计新增493台主机，耗时312.0秒',312,'2025-12-28 13:45:12.403'),(17,1,'云主机同步','[2025-12-28 13:51:00] 开始同步\n\n阿里云:\n  - test[阿里云]: 获取379台，新增379台，跳过0台，失败0台 [cn-qingdao:25台 cn-beijing:6台 cn-hangzhou:341台 cn-shenzhen:2台 cn-hongkong:3台 us-east-1:2台 ]\n  - 张三-阿里云-密钥[阿里云]: 获取1台，新增1台，跳过0台，失败0台 [cn-wulanchabu:1台 ]\n\n腾讯云:\n  - 腾讯云批量主机[腾讯云]: 获取82台，新增82台，跳过0台，失败0台 [ap-singapore:20台 na-siliconvalley:17台 ap-guangzhou:43台 ap-hongkong:1台 ap-tokyo:1台 ]\n\nAWS云:\n  - aws主机同步密钥[AWS云]: 获取31台，新增31台，跳过0台，失败0台 [ap-southeast-1:5台 us-west-2:26台 ]\n\n百度云:\n  - 张三-百度云-密钥[百度云]: 获取1台，新增1台，跳过0台，失败0台 [bj:1台 ]\n\n华为云:\n  - 张三-华为云-密钥[华为云]: 获取1台，新增1台，跳过0台，失败0台 [cn-south-1:1台 ]\n本次同步 495 台新增主机，6 个密钥，共计耗时 5m15.417975024s\n[2025-12-28 13:56:15] 同步完成\n',1,9,'同步6个密钥，共计新增495台主机，耗时315.4秒',315,'2025-12-28 13:56:15.714'),(18,1,'手动同步任务: test','[2025-12-28 14:36:20] 手动触发同步\n\n华为云:\n  - 张三-华为云-密钥[华为云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n\n阿里云:\n  - 张三-阿里云-密钥[阿里云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n本次同步 0 台新增主机，2 个密钥，共计耗时 35.24539009s\n[2025-12-28 14:36:55] 同步完成\n',1,10,'手动触发同步2个密钥，新增0台主机',35245,'2025-12-28 14:36:55.821'),(19,1,'手动同步任务: test','[2025-12-28 14:44:30] 手动触发同步\n\n阿里云:\n  - 张三-阿里云-密钥[阿里云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n\n华为云:\n  - 张三-华为云-密钥[华为云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n本次同步 0 台新增主机，2 个密钥，共计耗时 27.574998714s\n[2025-12-28 14:44:58] 同步完成\n',1,10,'手动触发同步2个密钥，新增0台主机',27574,'2025-12-28 14:44:58.614'),(20,1,'手动同步任务: 云主机同步','[2025-12-28 14:45:49] 手动触发同步\n\n华为云:\n  - 张三-华为云-密钥[华为云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n\n阿里云:\n  - test[阿里云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n  - 张三-阿里云-密钥[阿里云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n\n腾讯云:\n  - 腾讯云批量主机[腾讯云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n\nAWS云:\n  - aws主机同步密钥[AWS云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n\n百度云:\n  - 张三-百度云-密钥[百度云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n本次同步 0 台新增主机，6 个密钥，共计耗时 3m17.106106546s\n[2025-12-28 14:49:06] 同步完成\n',1,9,'手动触发同步6个密钥，新增0台主机',197106,'2025-12-28 14:49:07.063'),(21,1,'手动同步任务: test','[2025-12-28 15:21:24] 手动触发同步\n\n华为云:\n  - 张三-华为云-密钥[华为云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n\n百度云:\n  - 张三-百度云-密钥[百度云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n本次同步 0 台新增主机，2 个密钥，共计耗时 29.645202763s\n[2025-12-28 15:21:54] 同步完成\n',1,10,'手动触发同步2个密钥，新增0台主机',29645,'2025-12-28 15:21:54.399'),(22,1,'手动同步任务: test','[2025-12-28 15:22:23] 手动触发同步\n\n百度云:\n  - 张三-百度云-密钥[百度云]: 获取1台，新增1台，跳过0台，失败0台\n\n华为云:\n  - 张三-华为云-密钥[华为云]: 获取1台，新增1台，跳过0台，失败0台\n\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n本次同步 2 台新增主机，2 个密钥，共计耗时 25.748850164s\n[2025-12-28 15:22:49] 同步完成\n',1,10,'手动触发同步2个密钥，新增2台主机',25748,'2025-12-28 15:22:49.678'),(23,1,'手动同步任务: test1','[2025-12-28 15:29:54] 手动触发同步\n\nAWS云:\n  - aws主机同步密钥[AWS云]: 获取28台，新增28台，跳过0台，失败0台\n\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n本次同步 28 台新增主机，1 个密钥，共计耗时 55.546285392s\n[2025-12-28 15:30:50] 同步完成\n',1,12,'手动触发同步1个密钥，新增28台主机',55546,'2025-12-28 15:30:50.702'),(24,1,'云主机同步','[2025-12-28 15:46:00] 开始同步\n\n阿里云:\n  - test[阿里云]: 获取379台，新增379台，跳过0台，失败0台 [us-east-1:2台 cn-qingdao:25台 cn-beijing:6台 cn-hangzhou:341台 cn-shenzhen:2台 cn-hongkong:3台 ]\n  - 张三-阿里云-密钥[阿里云]: 获取1台，新增1台，跳过0台，失败0台 [cn-wulanchabu:1台 ]\n\n腾讯云:\n  - 腾讯云批量主机[腾讯云]: 获取82台，新增82台，跳过0台，失败0台 [ap-singapore:20台 na-siliconvalley:17台 ap-guangzhou:43台 ap-hongkong:1台 ap-tokyo:1台 ]\n\nAWS云:\n  - aws主机同步密钥[AWS云]: 获取31台，新增31台，跳过0台，失败0台 [ap-southeast-1:5台 us-west-2:26台 ]\n\n百度云:\n  - 张三-百度云-密钥[百度云]: 获取1台，新增1台，跳过0台，失败0台 [bj:1台 ]\n\n华为云:\n  - 张三-华为云-密钥[华为云]: 获取1台，新增1台，跳过0台，失败0台 [cn-south-1:1台 ]\n本次同步 495 台新增主机，6 个密钥，共计耗时 5m1.682881025s\n[2025-12-28 15:51:01] 同步完成\n',1,9,'同步6个密钥，共计新增495台主机，耗时301.7秒',301,'2025-12-28 15:51:02.024'),(25,1,'手动同步任务: test1','[2025-12-28 22:02:54] 手动触发同步\n\nAWS云:\n  - aws主机同步密钥[AWS云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n本次同步 0 台新增主机，1 个密钥，共计耗时 1m10.84188003s\n[2025-12-28 22:04:05] 同步完成\n',1,12,'手动触发同步1个密钥，新增0台主机',70841,'2025-12-28 22:04:05.522'),(26,1,'手动同步任务: test','[2025-12-28 22:59:41] 手动触发同步\n\n华为云:\n  - 张三-华为云-密钥[华为云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例（已更新 1 台主机配置）\n\n百度云:\n  - 张三-百度云-密钥[百度云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例（已更新 1 台主机配置）\n\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n本次同步 0 台新增主机，2 个密钥，共计耗时 26.44710232s\n[2025-12-28 23:00:08] 同步完成\n',1,10,'手动触发同步2个密钥，新增0台主机',26447,'2025-12-28 23:00:08.483'),(27,1,'手动同步任务: test','[2025-12-28 23:02:18] 手动触发同步\n\n华为云:\n  - 张三-华为云-密钥[华为云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例（已更新 1 台主机配置）\n\n百度云:\n  - 张三-百度云-密钥[百度云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例（已更新 1 台主机配置）\n\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n本次同步 0 台新增主机，2 个密钥，共计耗时 24.331334233s\n[2025-12-28 23:02:42] 同步完成\n',1,10,'手动触发同步2个密钥，新增0台主机',24331,'2025-12-28 23:02:43.216'),(28,1,'手动同步任务: test1','[2025-12-28 23:03:58] 手动触发同步\n\nAWS云:\n  - aws主机同步密钥[AWS云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例（已更新 2 台主机配置）\n\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n本次同步 0 台新增主机，1 个密钥，共计耗时 1m25.105810079s\n[2025-12-28 23:05:23] 同步完成\n',1,12,'手动触发同步1个密钥，新增0台主机',85105,'2025-12-28 23:05:23.419'),(29,2,'域名定时巡检完成','巡检结果：\n- 巡检域名总数：17\n- 即将过期域名数：1（SSL证书剩余天数<30天）\n- 执行耗时：2.89秒',1,0,'本次巡检 17 个域名，1 个域名即将过期',2,'2025-12-28 23:24:02.894'),(30,1,'云主机同步','[2026-01-01 15:46:00] 开始同步\n\n阿里云:\n  - test[阿里云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n  - 张三-阿里云-密钥[阿里云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n\n腾讯云:\n  - 腾讯云批量主机[腾讯云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n\nAWS云:\n  - aws主机同步密钥[AWS云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n\n百度云:\n  - 张三-百度云-密钥[百度云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n\n华为云:\n  - 张三-华为云-密钥[华为云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n本次同步 0 台新增主机，6 个密钥，共计耗时 3m10.777401569s\n[2026-01-01 15:49:10] 同步完成\n',1,9,'同步6个密钥，共计新增0台主机，耗时190.8秒',190,'2026-01-01 15:49:11.120'),(31,2,'域名定时巡检完成','巡检结果：\n- 巡检域名总数：17\n- 即将过期域名数：1（SSL证书剩余天数<30天）\n- 执行耗时：2.59秒',1,0,'本次巡检 17 个域名，1 个域名即将过期',2,'2026-01-01 23:24:02.593'),(32,1,'test','[2026-01-02 00:00:00] 开始同步\n\n百度云:\n  - 张三-百度云-密钥[百度云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n\n华为云:\n  - 张三-华为云-密钥[华为云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n本次同步 0 台新增主机，2 个密钥，共计耗时 33.593165805s\n[2026-01-02 00:00:33] 同步完成\n',1,10,'同步2个密钥，共计新增0台主机，耗时33.6秒',33,'2026-01-02 00:00:34.012'),(33,1,'test1','[2026-01-02 00:00:00] 开始同步\n\nAWS云:\n  - aws主机同步密钥[AWS云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n本次同步 0 台新增主机，1 个密钥，共计耗时 1m24.80736483s\n[2026-01-02 00:01:24] 同步完成\n',1,12,'同步1个密钥，共计新增0台主机，耗时84.8秒',84,'2026-01-02 00:01:25.168'),(34,1,'云主机同步','[2026-01-02 15:46:00] 开始同步\n\nAWS云:\n  - aws主机同步密钥[AWS云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n\n百度云:\n  - 张三-百度云-密钥[百度云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n\n华为云:\n  - 张三-华为云-密钥[华为云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n\n阿里云:\n  - test[阿里云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n  - 张三-阿里云-密钥[阿里云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n\n腾讯云:\n  - 腾讯云批量主机[腾讯云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n本次同步 0 台新增主机，6 个密钥，共计耗时 2m58.90052237s\n[2026-01-02 15:48:58] 同步完成\n',1,9,'同步6个密钥，共计新增0台主机，耗时178.9秒',178,'2026-01-02 15:48:59.313'),(35,1,'云主机同步','[2026-01-02 15:46:00] 开始同步\n\n华为云:\n  - 张三-华为云-密钥[华为云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n\n阿里云:\n  - test[阿里云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n  - 张三-阿里云-密钥[阿里云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n\n腾讯云:\n  - 腾讯云批量主机[腾讯云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n\nAWS云:\n  - aws主机同步密钥[AWS云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n\n百度云:\n  - 张三-百度云-密钥[百度云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n本次同步 0 台新增主机，6 个密钥，共计耗时 3m13.839925933s\n[2026-01-02 15:49:13] 同步完成\n',1,9,'同步6个密钥，共计新增0台主机，耗时193.8秒',193,'2026-01-02 15:49:14.209'),(36,2,'域名定时巡检完成','巡检结果：\n- 巡检域名总数：17\n- 即将过期域名数：1（SSL证书剩余天数<30天）\n- 执行耗时：3.96秒',1,0,'本次巡检 17 个域名，1 个域名即将过期',3,'2026-01-02 23:24:04.039'),(37,1,'云主机同步','[2026-01-03 15:46:00] 开始同步\n\nAWS云:\n  - aws主机同步密钥[AWS云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n\n百度云:\n  - 张三-百度云-密钥[百度云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n\n华为云:\n  - 张三-华为云-密钥[华为云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n\n阿里云:\n  - test[阿里云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n  - 张三-阿里云-密钥[阿里云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n\n腾讯云:\n  - 腾讯云批量主机[腾讯云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n本次同步 0 台新增主机，6 个密钥，共计耗时 3m23.459038453s\n[2026-01-03 15:49:23] 同步完成\n',1,9,'同步6个密钥，共计新增0台主机，耗时203.5秒',203,'2026-01-03 15:49:23.901'),(38,2,'域名定时巡检完成','巡检结果：\n- 巡检域名总数：17\n- 即将过期域名数：1（SSL证书剩余天数<30天）\n- 执行耗时：2.13秒',1,0,'本次巡检 17 个域名，1 个域名即将过期',2,'2026-01-03 23:24:02.138'),(39,2,'域名定时巡检完成','巡检结果：\n- 巡检域名总数：17\n- 即将过期域名数：1（SSL证书剩余天数<30天）\n- 执行耗时：1.36秒',1,0,'本次巡检 17 个域名，1 个域名即将过期',1,'2026-01-05 11:35:08.872'),(40,1,'云主机同步','[2026-01-05 15:46:00] 开始同步\n\n阿里云:\n  - test[阿里云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n  - 张三-阿里云-密钥[阿里云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n\n腾讯云:\n  - 腾讯云批量主机[腾讯云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n\nAWS云:\n  - aws主机同步密钥[AWS云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n\n百度云:\n  - 张三-百度云-密钥[百度云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n\n华为云:\n  - 张三-华为云-密钥[华为云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n本次同步 0 台新增主机，6 个密钥，共计耗时 2m37.627749811s\n[2026-01-05 15:48:37] 同步完成\n',1,9,'同步6个密钥，共计新增0台主机，耗时157.6秒',157,'2026-01-05 15:48:37.958'),(41,2,'域名定时巡检完成','巡检结果：\n- 巡检域名总数：17\n- 即将过期域名数：1（SSL证书剩余天数<30天）\n- 执行耗时：4.31秒',1,0,'本次巡检 17 个域名，1 个域名即将过期',4,'2026-01-05 23:24:04.313'),(42,1,'云主机同步','[2026-01-06 15:46:00] 开始同步\n\nAWS云:\n  - aws主机同步密钥[AWS云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n\n百度云:\n  - 张三-百度云-密钥[百度云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n\n华为云:\n  - 张三-华为云-密钥[华为云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n\n阿里云:\n  - test[阿里云]: 获取1台，新增1台，跳过0台，失败0台 [cn-hangzhou:1台 ]\n  - 张三-阿里云-密钥[阿里云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n\n腾讯云:\n  - 腾讯云批量主机[腾讯云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n本次同步 1 台新增主机，6 个密钥，共计耗时 2m31.321175594s\n[2026-01-06 15:48:31] 同步完成\n',1,9,'同步6个密钥，共计新增1台主机，耗时151.3秒',151,'2026-01-06 15:48:31.614'),(43,2,'域名定时巡检完成','巡检结果：\n- 巡检域名总数：17\n- 即将过期域名数：1（SSL证书剩余天数<30天）\n- 执行耗时：1.55秒',1,0,'本次巡检 17 个域名，1 个域名即将过期',1,'2026-01-06 23:24:01.550'),(44,2,'域名定时巡检完成','巡检结果：\n- 巡检域名总数：17\n- 即将过期域名数：1（SSL证书剩余天数<30天）\n- 执行耗时：2.53秒',1,0,'本次巡检 17 个域名，1 个域名即将过期',2,'2026-01-07 23:24:02.566'),(45,1,'云主机同步','[2026-01-08 16:21:01] 开始同步\n\n腾讯云:\n  - 腾讯云批量主机[腾讯云]: 获取82台，新增82台，跳过0台，失败0台 [ap-guangzhou:43台 ap-hongkong:1台 ap-tokyo:1台 ap-singapore:20台 na-siliconvalley:17台 ]\n\nAWS云:\n  - aws主机同步密钥[AWS云]: 获取31台，新增31台，跳过0台，失败0台 [ap-southeast-1:5台 us-west-2:26台 ]\n\n百度云:\n  - 张三-百度云-密钥[百度云]: 获取1台，新增1台，跳过0台，失败0台 [bj:1台 ]\n\n华为云:\n  - 张三-华为云-密钥[华为云]: 获取1台，新增1台，跳过0台，失败0台 [cn-south-1:1台 ]\n\n阿里云:\n  - test[阿里云]: 获取25台，新增25台，跳过0台，失败0台 [cn-qingdao:25台 ]\n  - 张三-阿里云-密钥[阿里云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n本次同步 140 台新增主机，6 个密钥，共计耗时 3m14.372654465s\n[2026-01-08 16:24:16] 同步完成\n',1,9,'同步6个密钥，共计新增140台主机，耗时194.4秒',194,'2026-01-08 16:24:16.622'),(46,1,'云主机同步','[2026-01-08 16:21:01] 开始同步\n\n阿里云:\n  - test[阿里云]: 获取380台，新增380台，跳过0台，失败0台 [cn-hongkong:3台 us-east-1:2台 cn-qingdao:25台 cn-beijing:6台 cn-hangzhou:342台 cn-shenzhen:2台 ]\n  - 张三-阿里云-密钥[阿里云]: 获取1台，新增1台，跳过0台，失败0台 [cn-wulanchabu:1台 ]\n\n腾讯云:\n  - 腾讯云批量主机[腾讯云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n\nAWS云:\n  - aws主机同步密钥[AWS云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n\n百度云:\n  - 张三-百度云-密钥[百度云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n\n华为云:\n  - 张三-华为云-密钥[华为云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例\n本次同步 381 台新增主机，6 个密钥，共计耗时 3m46.770498547s\n[2026-01-08 16:24:48] 同步完成\n',1,9,'同步6个密钥，共计新增381台主机，耗时226.8秒',226,'2026-01-08 16:24:48.744'),(47,1,'手动同步任务: test','[2026-01-08 20:09:56] 手动触发同步\n\n华为云:\n  - 张三-华为云-密钥[华为云]: 获取1台，新增1台，跳过0台，失败0台\n\n百度云:\n  - 张三-百度云-密钥[百度云]: 获取1台，新增1台，跳过0台，失败0台\n\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n本次同步 2 台新增主机，2 个密钥，共计耗时 29.615005407s\n[2026-01-08 20:10:26] 同步完成\n',1,10,'手动触发同步2个密钥，新增2台主机',29615,'2026-01-08 20:10:26.826'),(48,2,'域名定时巡检完成','巡检结果：\n- 巡检域名总数：17\n- 即将过期域名数：1（SSL证书剩余天数<30天）\n- 执行耗时：3.73秒',1,0,'本次巡检 17 个域名，1 个域名即将过期',3,'2026-01-09 23:24:03.729'),(49,1,'云主机同步','[2026-01-10 15:46:00] 开始同步\n\n腾讯云:\n  - 腾讯云批量主机[腾讯云]: 获取82台，新增82台，跳过0台，失败0台 [ap-guangzhou:43台 ap-hongkong:1台 ap-tokyo:1台 ap-singapore:20台 na-siliconvalley:17台 ]\n\nAWS云:\n  - aws主机同步密钥[AWS云]: 获取31台，新增31台，跳过0台，失败0台 [ap-southeast-1:5台 us-west-2:26台 ]\n\n百度云:\n  - 张三-百度云-密钥[百度云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例（已更新 1 台主机配置）\n\n华为云:\n  - 张三-华为云-密钥[华为云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例（已更新 1 台主机配置）\n\n阿里云:\n  - test[阿里云]: 获取380台，新增380台，跳过0台，失败0台 [cn-qingdao:25台 cn-beijing:6台 cn-hangzhou:342台 cn-shenzhen:2台 cn-hongkong:3台 us-east-1:2台 ]\n  - 张三-阿里云-密钥[阿里云]: 同步失败 - 同步失败: 所有实例均已导入，无新增实例（已更新 1 台主机配置）\n本次同步 493 台新增主机，6 个密钥，共计耗时 3m20.771246s\n[2026-01-10 15:49:20] 同步完成\n',1,9,'同步6个密钥，共计新增493台主机，耗时200.8秒',200,'2026-01-10 15:49:20.959'),(50,2,'域名定时巡检完成','巡检结果：\n- 巡检域名总数：17\n- 即将过期域名数：1（SSL证书剩余天数<30天）\n- 执行耗时：4.91秒',1,0,'本次巡检 17 个域名，1 个域名即将过期',4,'2026-01-13 23:24:04.913'),(51,2,'域名定时巡检完成','巡检结果：\n- 巡检域名总数：17\n- 即将过期域名数：1（SSL证书剩余天数<30天）\n- 执行耗时：3.33秒',1,0,'本次巡检 17 个域名，1 个域名即将过期',3,'2026-01-14 23:24:03.463'),(52,2,'域名定时巡检完成','巡检结果：\n- 巡检域名总数：17\n- 即将过期域名数：1（SSL证书剩余天数<30天）\n- 执行耗时：1.82秒',1,0,'本次巡检 17 个域名，1 个域名即将过期',1,'2026-01-15 12:53:01.820'),(53,2,'域名定时巡检完成','巡检结果：\n- 巡检域名总数：17\n- 即将过期域名数：1（SSL证书剩余天数<30天）\n- 执行耗时：1.36秒',1,0,'本次巡检 17 个域名，1 个域名即将过期',1,'2026-01-15 13:05:01.359'),(54,2,'域名定时巡检完成','巡检结果：\n- 巡检域名总数：17\n- 即将过期域名数：1（SSL证书剩余天数<30天）\n- 执行耗时：1.31秒',1,0,'本次巡检 17 个域名，1 个域名即将过期',1,'2026-01-15 13:14:01.313'),(55,2,'域名定时巡检完成','巡检结果：\n- 巡检域名总数：17\n- 即将过期域名数：1（SSL证书剩余天数<30天）\n- 执行耗时：1.57秒',1,0,'本次巡检 17 个域名，1 个域名即将过期',1,'2026-01-15 13:21:01.574'),(56,2,'域名定时巡检完成','巡检结果：\n- 巡检域名总数：17\n- 即将过期域名数：1（SSL证书剩余天数<30天）\n- 执行耗时：1.37秒',1,0,'本次巡检 17 个域名，1 个域名即将过期',1,'2026-01-15 13:28:01.372'),(57,2,'域名定时巡检完成','巡检结果：\n- 巡检域名总数：17\n- 即将过期域名数：1（SSL证书剩余天数<30天）\n- 执行耗时：1.20秒',1,0,'本次巡检 17 个域名，1 个域名即将过期',1,'2026-01-15 13:29:01.203'),(58,2,'域名定时巡检完成','巡检结果：\n- 巡检域名总数：17\n- 即将过期域名数：1（SSL证书剩余天数<30天）\n- 执行耗时：1.60秒',1,0,'本次巡检 17 个域名，1 个域名即将过期',1,'2026-01-16 13:29:01.657'),(59,2,'域名定时巡检完成','巡检结果：\n- 巡检域名总数：17\n- 即将过期域名数：0（SSL证书剩余天数<30天）\n- 执行耗时：1.83秒',1,0,'本次巡检 17 个域名，0 个域名即将过期',1,'2026-01-17 13:29:03.650'),(60,2,'域名定时巡检完成','巡检结果：\n- 巡检域名总数：17\n- 即将过期域名数：1（SSL证书剩余天数<30天）\n- 执行耗时：1.33秒',1,0,'本次巡检 17 个域名，1 个域名即将过期',1,'2026-01-19 13:29:01.360');
/*!40000 ALTER TABLE `sys_activity_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_admin`
--

DROP TABLE IF EXISTS `sys_admin`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_admin` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键',
  `post_id` int DEFAULT NULL COMMENT '岗位id',
  `dept_id` int DEFAULT NULL COMMENT '部门id',
  `username` varchar(64) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL COMMENT '账号',
  `password` varchar(64) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL COMMENT '密码',
  `nickname` varchar(64) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT '昵称',
  `icon` varchar(500) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT '头像',
  `email` varchar(64) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT '邮箱',
  `phone` varchar(64) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT '手机',
  `note` varchar(500) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT '备注信息',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  `status` int NOT NULL DEFAULT '1' COMMENT '帐号启用状态：1->启用,2->禁用',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `username` (`username`) USING BTREE,
  KEY `create_time` (`create_time`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=109 DEFAULT CHARSET=utf8mb3 ROW_FORMAT=DYNAMIC COMMENT='后台管理员表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_admin`
--

LOCK TABLES `sys_admin` WRITE;
/*!40000 ALTER TABLE `sys_admin` DISABLE KEYS */;
INSERT INTO `sys_admin` VALUES (89,1,15,'admin','e10adc3949ba59abbe56e057f20f883e','管理员','http://10.7.16.100:8080/api/v1/upload/avatar/20260119/109120000.svg','123456789@qq.com','13754354536','后端研发','2023-05-23 22:15:50',1),(108,1,2,'root','e10adc3949ba59abbe56e057f20f883e','root','','zfwh1024@163.com','123','123','2026-02-19 19:37:42',1);
/*!40000 ALTER TABLE `sys_admin` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_admin_role`
--

DROP TABLE IF EXISTS `sys_admin_role`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_admin_role` (
  `admin_id` int NOT NULL COMMENT '管理员id',
  `role_id` int NOT NULL COMMENT '角色id'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 ROW_FORMAT=DYNAMIC COMMENT='管理员和角色关系表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_admin_role`
--

LOCK TABLES `sys_admin_role` WRITE;
/*!40000 ALTER TABLE `sys_admin_role` DISABLE KEYS */;
INSERT INTO `sys_admin_role` VALUES (89,1),(108,1);
/*!40000 ALTER TABLE `sys_admin_role` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_blocking_policy`
--

DROP TABLE IF EXISTS `sys_blocking_policy`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_blocking_policy` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `description` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `enabled` tinyint(1) DEFAULT '1',
  `priority` bigint DEFAULT '0',
  `block_mode` bigint DEFAULT '1',
  `enable_alias_resolution` tinyint(1) DEFAULT '1',
  `enable_script_scanning` tinyint(1) DEFAULT '1',
  `custom_rules` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
  `whitelist_cmds` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_blocking_policy`
--

LOCK TABLES `sys_blocking_policy` WRITE;
/*!40000 ALTER TABLE `sys_blocking_policy` DISABLE KEYS */;
/*!40000 ALTER TABLE `sys_blocking_policy` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_command_audit`
--

DROP TABLE IF EXISTS `sys_command_audit`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_command_audit` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键ID''',
  `recording_id` bigint unsigned NOT NULL COMMENT '''录制记录ID''',
  `session_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''会话ID''',
  `command` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''执行的命令''',
  `timestamp` double NOT NULL COMMENT '''相对时间戳(秒)''',
  `sequence` bigint NOT NULL COMMENT '''命令序号''',
  `is_sensitive` tinyint(1) DEFAULT '0' COMMENT '''是否为敏感命令''',
  `risk_level` bigint DEFAULT '0' COMMENT '''风险等级:0-正常 1-可疑 2-高危''',
  `risk_reason` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''风险原因''',
  `execute_time` datetime(3) NOT NULL COMMENT '''执行时间(绝对时间)''',
  `create_time` datetime(3) NOT NULL COMMENT '''创建时间''',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_sys_command_audit_recording_id` (`recording_id`) USING BTREE,
  KEY `idx_sys_command_audit_session_id` (`session_id`) USING BTREE,
  KEY `idx_sys_command_audit_is_sensitive` (`is_sensitive`) USING BTREE,
  KEY `idx_sys_command_audit_risk_level` (`risk_level`) USING BTREE,
  KEY `idx_sys_command_audit_execute_time` (`execute_time`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=181 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_command_audit`
--

LOCK TABLES `sys_command_audit` WRITE;
/*!40000 ALTER TABLE `sys_command_audit` DISABLE KEYS */;
INSERT INTO `sys_command_audit` VALUES (1,1,'1767874274_3891b5c24c9f03d3','ls',2.80614476,1,0,0,'','2026-01-08 20:11:17.703','2026-01-08 20:11:23.162'),(2,1,'1767874274_3891b5c24c9f03d3','cat 1.sql',6.878329228,2,0,0,'','2026-01-08 20:11:21.775','2026-01-08 20:11:23.162'),(3,2,'1767874295_52026bc3c4ad4d16','ls',1.203503921,1,0,0,'','2026-01-08 20:11:36.248','2026-01-08 20:12:04.280'),(4,2,'1767874295_52026bc3c4ad4d16','vim 1.sql',6.002750905,2,0,0,'','2026-01-08 20:11:41.047','2026-01-08 20:12:04.280'),(5,2,'1767874295_52026bc3c4ad4d16','ls',15.174233342,3,0,0,'','2026-01-08 20:11:50.218','2026-01-08 20:12:04.280'),(6,2,'1767874295_52026bc3c4ad4d16','cat 1.sql',28.79382693,4,0,0,'','2026-01-08 20:12:03.838','2026-01-08 20:12:04.280'),(7,3,'1767874373_3f89dfb692f1aeb7','ls',77.005688101,1,0,0,'','2026-01-08 20:14:10.461','2026-01-08 20:15:59.683'),(8,3,'1767874373_3f89dfb692f1aeb7','ls elk/',80.187585552,2,0,0,'','2026-01-08 20:14:13.643','2026-01-08 20:15:59.683'),(9,3,'1767874373_3f89dfb692f1aeb7','ls luban-master/',83.943315437,3,0,0,'','2026-01-08 20:14:17.398','2026-01-08 20:15:59.684'),(10,3,'1767874373_3f89dfb692f1aeb7','cd  /opt/',141.774468394,4,0,0,'','2026-01-08 20:15:15.229','2026-01-08 20:15:59.684'),(11,3,'1767874373_3f89dfb692f1aeb7','ls',142.5387535,5,0,0,'','2026-01-08 20:15:15.994','2026-01-08 20:15:59.684'),(12,3,'1767874373_3f89dfb692f1aeb7','systemctl status   [200~prometheus[201~',159.812634471,6,0,0,'','2026-01-08 20:15:33.268','2026-01-08 20:15:59.684'),(13,3,'1767874373_3f89dfb692f1aeb7','ls',179.189554404,7,0,0,'','2026-01-08 20:15:52.645','2026-01-08 20:15:59.684'),(14,4,'1767874590_2cbaa6c8d286619a','ls',2.665041017,1,0,0,'','2026-01-08 20:16:33.113','2026-01-08 20:17:25.114'),(15,4,'1767874590_2cbaa6c8d286619a','ls elk',6.091859938,2,0,0,'','2026-01-08 20:16:36.540','2026-01-08 20:17:25.114'),(16,4,'1767874590_2cbaa6c8d286619a','[A/logs/',8.895688495,3,0,0,'','2026-01-08 20:16:39.344','2026-01-08 20:17:25.114'),(17,4,'1767874590_2cbaa6c8d286619a','rm -rf elk/',14.107609288,4,1,1,'包含敏感操作关键词','2026-01-08 20:16:44.556','2026-01-08 20:17:25.115'),(18,4,'1767874590_2cbaa6c8d286619a','ls',15.059702342,5,0,0,'','2026-01-08 20:16:45.508','2026-01-08 20:17:25.115'),(19,4,'1767874590_2cbaa6c8d286619a','ls luban-master/',18.74447007,6,0,0,'','2026-01-08 20:16:49.192','2026-01-08 20:17:25.116'),(20,4,'1767874590_2cbaa6c8d286619a','rm -rf luban-master/',30.774242247,7,1,1,'包含敏感操作关键词','2026-01-08 20:17:01.222','2026-01-08 20:17:25.116'),(21,4,'1767874590_2cbaa6c8d286619a','ls',31.790972634,8,0,0,'','2026-01-08 20:17:02.239','2026-01-08 20:17:25.116'),(22,4,'1767874590_2cbaa6c8d286619a','ls docker/',36.857322249,9,0,0,'','2026-01-08 20:17:07.305','2026-01-08 20:17:25.117'),(23,4,'1767874590_2cbaa6c8d286619a','[Afilebeat.yml/',40.616391193,10,0,0,'','2026-01-08 20:17:11.064','2026-01-08 20:17:25.117'),(24,5,'1767875841_1747f4f17da86ec3','ls',3.465370439,1,0,0,'','2026-01-08 20:37:24.992','2026-01-08 20:38:58.544'),(25,5,'1767875841_1747f4f17da86ec3','cat 1.sql',7.937137411,2,0,0,'','2026-01-08 20:37:29.464','2026-01-08 20:38:58.544'),(26,5,'1767875841_1747f4f17da86ec3','ls',11.052666884,3,0,0,'','2026-01-08 20:37:32.580','2026-01-08 20:38:58.544'),(27,5,'1767875841_1747f4f17da86ec3','ls /',88.56829096,4,0,0,'','2026-01-08 20:38:50.095','2026-01-08 20:38:58.544'),(28,6,'1767885992_b10b66647164e108','[Ils',2.82115313,1,0,0,'','2026-01-08 23:26:35.665','2026-01-08 23:27:51.720'),(29,6,'1767885992_b10b66647164e108','vim 1.sh',7.008717398,2,0,0,'','2026-01-08 23:26:39.853','2026-01-08 23:27:51.720'),(30,6,'1767885992_b10b66647164e108','[A[D[D[2',12.98857792,3,0,0,'','2026-01-08 23:26:45.833','2026-01-08 23:27:51.720'),(31,6,'1767885992_b10b66647164e108','[I[2;2R[3;1R[>0;276;0c]10;rgb:d4d4/d4d4/d4d4\\]11;rgb:1e1e/1e1e/1e1e\\[?12;2$[200~#!/bin/bash',15.418203422,4,0,0,'','2026-01-08 23:26:48.262','2026-01-08 23:27:51.720'),(32,6,'1767885992_b10b66647164e108','ZZlls',20.045342626,5,0,0,'','2026-01-08 23:26:52.889','2026-01-08 23:27:51.720'),(33,6,'1767885992_b10b66647164e108','ls',22.701426176,6,0,0,'','2026-01-08 23:26:55.545','2026-01-08 23:27:51.721'),(34,6,'1767885992_b10b66647164e108','chmod 2.sh',28.502691815,7,0,0,'','2026-01-08 23:27:01.347','2026-01-08 23:27:51.721'),(35,6,'1767885992_b10b66647164e108','[A[D[D[D[D[D[D[C +x',32.657462694,8,0,0,'','2026-01-08 23:27:05.501','2026-01-08 23:27:51.721'),(36,6,'1767885992_b10b66647164e108','ls',34.257368173,9,0,0,'','2026-01-08 23:27:07.101','2026-01-08 23:27:51.721'),(37,6,'1767885992_b10b66647164e108','sh 2.sh',39.350872312999996,10,0,0,'','2026-01-08 23:27:12.195','2026-01-08 23:27:51.721'),(38,9,'1768278417_2aedd83a62053660','ls',14.400572401,1,0,0,'','2026-01-13 12:27:11.685','2026-01-13 12:28:34.709'),(39,9,'1768278417_2aedd83a62053660','pwd',15.831482235,2,0,0,'','2026-01-13 12:27:13.115','2026-01-13 12:28:34.709'),(40,13,'1768284202_bd282732818a32e7','l sls',16.310576145,1,0,0,'','2026-01-13 14:03:38.358','2026-01-13 14:05:03.487'),(41,13,'1768284202_bd282732818a32e7','ls',17.022053776,2,0,0,'','2026-01-13 14:03:39.069','2026-01-13 14:05:03.488'),(42,14,'1768284411_74526a07af1c60d3','ls',7.406490787,1,0,0,'','2026-01-13 14:06:58.771','2026-01-13 14:10:36.152'),(43,14,'1768284411_74526a07af1c60d3','catt 1.sql',12.982747926,2,0,0,'','2026-01-13 14:07:04.348','2026-01-13 14:10:36.152'),(44,14,'1768284411_74526a07af1c60d3','cat 1.sql',16.715096885,3,0,0,'','2026-01-13 14:07:08.080','2026-01-13 14:10:36.153'),(45,14,'1768284411_74526a07af1c60d3','ls',19.4900968,4,0,0,'','2026-01-13 14:07:10.855','2026-01-13 14:10:36.153'),(46,16,'1768285870_4ec148e24e293986','ls',7.342037644,1,0,0,'','2026-01-13 14:31:18.037','2026-01-13 14:34:16.402'),(47,23,'1768287202_a6abbb4e78b6a9de','ls',392.323654318,1,0,0,'','2026-01-13 14:59:54.746','2026-01-13 14:59:58.113'),(48,24,'1768292563_785e25cd14af7ba5','ls',3.005419186,1,0,0,'','2026-01-13 16:22:46.070','2026-01-13 16:23:00.823'),(49,24,'1768292563_785e25cd14af7ba5','ls',5.088186679,2,0,0,'','2026-01-13 16:22:48.153','2026-01-13 16:23:00.823'),(50,27,'1768316313_7e36e18cf1e36db8','ls',6.015181758,1,0,0,'','2026-01-13 22:58:39.972','2026-01-14 00:02:33.600'),(51,27,'1768316313_7e36e18cf1e36db8','netstat   -nutpl',36.648042621,2,0,0,'','2026-01-13 22:59:10.605','2026-01-14 00:02:33.601'),(52,27,'1768316313_7e36e18cf1e36db8','docker ps',42.232433523,3,0,0,'','2026-01-13 22:59:16.189','2026-01-14 00:02:33.601'),(53,27,'1768316313_7e36e18cf1e36db8','curl  http://ifconfig.io',61.513548719,4,1,1,'包含敏感操作关键词','2026-01-13 22:59:35.471','2026-01-14 00:02:33.601'),(54,27,'1768316313_7e36e18cf1e36db8','ls',332.732776459,5,0,0,'','2026-01-13 23:04:06.690','2026-01-14 00:02:33.601'),(55,27,'1768316313_7e36e18cf1e36db8','ls',535.893509441,6,0,0,'','2026-01-13 23:07:29.851','2026-01-14 00:02:33.601'),(56,27,'1768316313_7e36e18cf1e36db8','docker ps',1785.791008184,7,0,0,'','2026-01-13 23:28:19.748','2026-01-14 00:02:33.601'),(57,27,'1768316313_7e36e18cf1e36db8','[200~curl http://localhost:3000/api/user_group_map[201~',1998.480037555,8,1,1,'包含敏感操作关键词','2026-01-13 23:31:52.437','2026-01-14 00:02:33.601'),(58,27,'1768316313_7e36e18cf1e36db8','docker ps',3672.961247759,9,0,0,'','2026-01-13 23:59:46.918','2026-01-14 00:02:33.601'),(59,27,'1768316313_7e36e18cf1e36db8','free -h',3688.351055737,10,0,0,'','2026-01-14 00:00:02.308','2026-01-14 00:02:33.601'),(60,27,'1768316313_7e36e18cf1e36db8','df -h',3690.35872877,11,0,0,'','2026-01-14 00:00:04.316','2026-01-14 00:02:33.602'),(61,27,'1768316313_7e36e18cf1e36db8','[200~sudo fallocate -l 2G /swapfile[201~',3748.320698335,12,1,1,'包含敏感操作关键词','2026-01-14 00:01:02.278','2026-01-14 00:02:33.602'),(62,27,'1768316313_7e36e18cf1e36db8','[200~sudo chmod 600 /swapfile',3754.861026839,13,1,1,'包含敏感操作关键词','2026-01-14 00:01:08.818','2026-01-14 00:02:33.602'),(63,27,'1768316313_7e36e18cf1e36db8','[200~sudo mkswap /swapfile',3759.251763057,14,1,1,'包含敏感操作关键词','2026-01-14 00:01:13.209','2026-01-14 00:02:33.602'),(64,27,'1768316313_7e36e18cf1e36db8','[200~sudo swapon /swapfile',3767.363378288,15,1,1,'包含敏感操作关键词','2026-01-14 00:01:21.320','2026-01-14 00:02:33.602'),(65,27,'1768316313_7e36e18cf1e36db8','[200~free -h',3771.71440026,16,0,0,'','2026-01-14 00:01:25.671','2026-01-14 00:02:33.602'),(66,27,'1768316313_7e36e18cf1e36db8','[200~echo \'/swapfile none swap sw 0 0\' | sudo tee -a /etc/fstab',3779.354333544,17,1,1,'包含敏感操作关键词','2026-01-14 00:01:33.311','2026-01-14 00:02:33.602'),(67,27,'1768316313_7e36e18cf1e36db8','[200~cat /etc/fstab | grep swapfile',3784.418108789,18,0,0,'','2026-01-14 00:01:38.375','2026-01-14 00:02:33.602'),(68,27,'1768316313_7e36e18cf1e36db8','[200~free -h',3800.08413335,19,0,0,'','2026-01-14 00:01:54.041','2026-01-14 00:02:33.602'),(69,27,'1768316313_7e36e18cf1e36db8','[200~swapon --show',3804.619038778,20,0,0,'','2026-01-14 00:01:58.576','2026-01-14 00:02:33.602'),(70,28,'1768316317_e6f89da89f84cd77','netstat   -nutpl',11.729624523,1,0,0,'','2026-01-13 22:58:49.520','2026-01-14 00:17:42.785'),(71,28,'1768316317_e6f89da89f84cd77','[A  | grep 3000',17.475192977,2,0,0,'','2026-01-13 22:58:55.265','2026-01-14 00:17:42.785'),(72,28,'1768316317_e6f89da89f84cd77','cd /etc/nginx/',66.852047002,3,0,0,'','2026-01-13 22:59:44.642','2026-01-14 00:17:42.786'),(73,28,'1768316317_e6f89da89f84cd77','ls',67.949395164,4,0,0,'','2026-01-13 22:59:45.739','2026-01-14 00:17:42.786'),(74,28,'1768316317_e6f89da89f84cd77','cd conf.d/',70.234091905,5,0,0,'','2026-01-13 22:59:48.024','2026-01-14 00:17:42.787'),(75,28,'1768316317_e6f89da89f84cd77','ls',70.783639971,6,0,0,'','2026-01-13 22:59:48.574','2026-01-14 00:17:42.787'),(76,28,'1768316317_e6f89da89f84cd77','cat AutoOps.conf',77.955560388,7,0,0,'','2026-01-13 22:59:55.746','2026-01-14 00:17:42.787'),(77,28,'1768316317_e6f89da89f84cd77','[Als',94.30170778,8,0,0,'','2026-01-13 23:00:12.092','2026-01-14 00:17:42.787'),(78,28,'1768316317_e6f89da89f84cd77','ls',96.337848812,9,0,0,'','2026-01-13 23:00:14.128','2026-01-14 00:17:42.788'),(79,28,'1768316317_e6f89da89f84cd77','cat devops.conf',99.026360098,10,0,0,'','2026-01-13 23:00:16.816','2026-01-14 00:17:42.788'),(80,28,'1768316317_e6f89da89f84cd77','ls',110.657067443,11,0,0,'','2026-01-13 23:00:28.447','2026-01-14 00:17:42.788'),(81,28,'1768316317_e6f89da89f84cd77','vim devops.conf',114.794213189,12,0,0,'','2026-01-13 23:00:32.584','2026-01-14 00:17:42.788'),(82,28,'1768316317_e6f89da89f84cd77','vim AutoOps.conf',128.469270883,13,0,0,'','2026-01-13 23:00:46.259','2026-01-14 00:17:42.788'),(83,28,'1768316317_e6f89da89f84cd77','nginx -t',148.627838984,14,0,0,'','2026-01-13 23:01:06.418','2026-01-14 00:17:42.788'),(84,28,'1768316317_e6f89da89f84cd77','ls',150.659797416,15,0,0,'','2026-01-13 23:01:08.450','2026-01-14 00:17:42.788'),(85,28,'1768316317_e6f89da89f84cd77','cat AutoOps.conf',158.531535477,16,0,0,'','2026-01-13 23:01:16.322','2026-01-14 00:17:42.789'),(86,28,'1768316317_e6f89da89f84cd77','ls',333.801062232,17,0,0,'','2026-01-13 23:04:11.591','2026-01-14 00:17:42.789'),(87,28,'1768316317_e6f89da89f84cd77','cp AutoOps.conf deviops.cn.conf[H[mv',361.454342189,18,0,0,'','2026-01-13 23:04:39.244','2026-01-14 00:17:42.789'),(88,28,'1768316317_e6f89da89f84cd77','ls',362.220644415,19,0,0,'','2026-01-13 23:04:40.011','2026-01-14 00:17:42.789'),(89,28,'1768316317_e6f89da89f84cd77','cp deviops.cn.conf deviops.cn.conf [D[D[D[D[Dai',377.850451932,20,0,0,'','2026-01-13 23:04:55.640','2026-01-14 00:17:42.789'),(90,28,'1768316317_e6f89da89f84cd77','vim deviops.ai.conf',383.89935745,21,0,0,'','2026-01-13 23:05:01.689','2026-01-14 00:17:42.789'),(91,28,'1768316317_e6f89da89f84cd77','ls',392.906639551,22,0,0,'','2026-01-13 23:05:10.697','2026-01-14 00:17:42.789'),(92,28,'1768316317_e6f89da89f84cd77','rm -rf deviops.ai.conf',399.996343799,23,1,1,'包含敏感操作关键词','2026-01-13 23:05:17.786','2026-01-14 00:17:42.789'),(93,28,'1768316317_e6f89da89f84cd77','vim deviops.cn.conf',405.561509876,24,0,0,'','2026-01-13 23:05:23.352','2026-01-14 00:17:42.790'),(94,28,'1768316317_e6f89da89f84cd77','nginx -t',419.971756289,25,0,0,'','2026-01-13 23:05:37.762','2026-01-14 00:17:42.790'),(95,28,'1768316317_e6f89da89f84cd77','systemctl restart nginx',426.217848884,26,0,0,'','2026-01-13 23:05:44.008','2026-01-14 00:17:42.790'),(96,28,'1768316317_e6f89da89f84cd77','ls',427.683605111,27,0,0,'','2026-01-13 23:05:45.474','2026-01-14 00:17:42.790'),(97,28,'1768316317_e6f89da89f84cd77','cat devops.conf',431.97275415,28,0,0,'','2026-01-13 23:05:49.763','2026-01-14 00:17:42.790'),(98,28,'1768316317_e6f89da89f84cd77','ls',433.997158071,29,0,0,'','2026-01-13 23:05:51.787','2026-01-14 00:17:42.790'),(99,28,'1768316317_e6f89da89f84cd77','ls',534.526195141,30,0,0,'','2026-01-13 23:07:32.316','2026-01-14 00:17:42.790'),(100,28,'1768316317_e6f89da89f84cd77','vim deviops.cn.conf',537.99637449,31,0,0,'','2026-01-13 23:07:35.786','2026-01-14 00:17:42.790'),(101,28,'1768316317_e6f89da89f84cd77','systemctl restart nginx',580.706194209,32,0,0,'','2026-01-13 23:08:18.496','2026-01-14 00:17:42.791'),(102,28,'1768316317_e6f89da89f84cd77','curl [200~http://180.76.231.65:3000/[201~[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D-I',589.133205987,33,1,1,'包含敏感操作关键词','2026-01-13 23:08:26.923','2026-01-14 00:17:42.791'),(103,28,'1768316317_e6f89da89f84cd77','[200~deviops.cn[201~[D[D[D[D[D[D[D[D[D[Dhttps://[C[C[C[C[C[C[C[C[C[C[C/ai',627.358534052,34,0,0,'','2026-01-13 23:09:05.149','2026-01-14 00:17:42.791'),(104,28,'1768316317_e6f89da89f84cd77','[A[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D[D',633.206113729,35,0,0,'','2026-01-13 23:09:10.996','2026-01-14 00:17:42.791'),(105,28,'1768316317_e6f89da89f84cd77','[A   ls',1083.204426974,36,0,0,'','2026-01-13 23:16:40.994','2026-01-14 00:17:42.791'),(106,28,'1768316317_e6f89da89f84cd77','ls',1087.954310531,37,0,0,'','2026-01-13 23:16:45.744','2026-01-14 00:17:42.792'),(107,28,'1768316317_e6f89da89f84cd77','cat deviops.cn.conf',1092.844349049,38,0,0,'','2026-01-13 23:16:50.634','2026-01-14 00:17:42.792'),(108,28,'1768316317_e6f89da89f84cd77','ls',1193.449721604,39,0,0,'','2026-01-13 23:18:31.240','2026-01-14 00:17:42.792'),(109,28,'1768316317_e6f89da89f84cd77','vim deviops.cn.conf',1200.220537001,40,0,0,'','2026-01-13 23:18:38.011','2026-01-14 00:17:42.792'),(110,28,'1768316317_e6f89da89f84cd77','[A[A[A[A[A[A[A[A[A',1213.869209923,41,0,0,'','2026-01-13 23:18:51.659','2026-01-14 00:17:42.793'),(111,28,'1768316317_e6f89da89f84cd77','[A[A',1623.441773403,42,0,0,'','2026-01-13 23:25:41.232','2026-01-14 00:17:42.793'),(112,28,'1768316317_e6f89da89f84cd77','[I[2;2R[3;1R[>0;276;0c]10;rgb:d4d4/d4d4/d4d4\\]11;rgb:1e1e/1e1e/1e1e\\[?12;2$ygg1000dda[200~# HTTPS  for deviops.cn',1629.034815126,43,1,1,'包含敏感操作关键词','2026-01-13 23:25:46.825','2026-01-14 00:17:42.793'),(113,28,'1768316317_e6f89da89f84cd77','ZZ[A[A[A[A[B[B',1634.082976584,44,0,0,'','2026-01-13 23:25:51.873','2026-01-14 00:17:42.793'),(114,28,'1768316317_e6f89da89f84cd77','ls',2707.11907249,45,0,0,'','2026-01-13 23:43:44.909','2026-01-14 00:17:42.793'),(115,28,'1768316317_e6f89da89f84cd77','curl  [200~curl -v http://180.76.231.65:3000/api/user_group_map[201~[200~curl -v http://180.76.231.65:3000/api/user_group_map[201~',2712.849749575,46,1,1,'包含敏感操作关键词','2026-01-13 23:43:50.640','2026-01-14 00:17:42.793'),(116,28,'1768316317_e6f89da89f84cd77','[200~#  nginx.conf',2765.231420295,47,0,0,'','2026-01-13 23:44:43.021','2026-01-14 00:17:42.793'),(117,28,'1768316317_e6f89da89f84cd77','[A',2768.91113746,48,0,0,'','2026-01-13 23:44:46.701','2026-01-14 00:17:42.793'),(118,28,'1768316317_e6f89da89f84cd77','[200~curl -H \"Host: deviops.cn\" http://127.0.0.1/api/user_group_map[201~',2790.05621414,49,1,1,'包含敏感操作关键词','2026-01-13 23:45:07.846','2026-01-14 00:17:42.794'),(119,28,'1768316317_e6f89da89f84cd77','[200~curl -k -H \"Host: deviops.cn\" https://127.0.0.1/api/user_group_map[201~',2808.365654357,50,1,1,'包含敏感操作关键词','2026-01-13 23:45:26.156','2026-01-14 00:17:42.794'),(120,28,'1768316317_e6f89da89f84cd77','[A[A[A[A[Als',2948.70384186,51,0,0,'','2026-01-13 23:47:46.494','2026-01-14 00:17:42.794'),(121,28,'1768316317_e6f89da89f84cd77','cat deviops.cn.conf',2952.770165855,52,0,0,'','2026-01-13 23:47:50.560','2026-01-14 00:17:42.794'),(122,28,'1768316317_e6f89da89f84cd77','ls',3009.264640652,53,0,0,'','2026-01-13 23:48:47.055','2026-01-14 00:17:42.794'),(123,28,'1768316317_e6f89da89f84cd77','vim  deviops.cn.conf',3013.770884208,54,0,0,'','2026-01-13 23:48:51.561','2026-01-14 00:17:42.794'),(124,28,'1768316317_e6f89da89f84cd77','[200~grep \"proxy_set_header Host\" /etc/nginx/conf.d/deviops.cn.conf[201~',3114.316680585,55,0,0,'','2026-01-13 23:50:32.107','2026-01-14 00:17:42.794'),(125,28,'1768316317_e6f89da89f84cd77','[200~sudo nginx -t && sudo systemctl reload nginx[201~',3125.352087523,56,1,1,'包含敏感操作关键词','2026-01-13 23:50:43.142','2026-01-14 00:17:42.794'),(126,28,'1768316317_e6f89da89f84cd77','ls',3551.37585761,57,0,0,'','2026-01-13 23:57:49.166','2026-01-14 00:17:42.795'),(127,28,'1768316317_e6f89da89f84cd77','cp deviops.cn.conf deviops.cn.bak',3573.075438903,58,0,0,'','2026-01-13 23:58:10.865','2026-01-14 00:17:42.795'),(128,28,'1768316317_e6f89da89f84cd77','ls',3574.296810284,59,0,0,'','2026-01-13 23:58:12.087','2026-01-14 00:17:42.795'),(129,28,'1768316317_e6f89da89f84cd77','vim deviops.cn.conf',3581.144246876,60,0,0,'','2026-01-13 23:58:18.934','2026-01-14 00:17:42.795'),(130,28,'1768316317_e6f89da89f84cd77','systemctl restart nginx',3615.018404581,61,0,0,'','2026-01-13 23:58:52.808','2026-01-14 00:17:42.795'),(131,29,'1768322330_5de93c439ce770af','df -h',12.6774142,1,0,0,'','2026-01-14 00:39:03.299','2026-01-14 00:39:32.750'),(132,29,'1768322330_5de93c439ce770af','free -h',16.806410586,2,0,0,'','2026-01-14 00:39:07.428','2026-01-14 00:39:32.752'),(133,29,'1768322330_5de93c439ce770af','docker ps',24.715458326,3,0,0,'','2026-01-14 00:39:15.337','2026-01-14 00:39:32.752'),(134,33,'1768575070_9ef963acef296c40','ls',6.00229082,1,0,0,'','2026-01-16 22:51:16.205','2026-01-16 22:52:47.651'),(135,33,'1768575070_9ef963acef296c40','cd /home/di',13.30212642,2,0,0,'','2026-01-16 22:51:23.505','2026-01-16 22:52:47.651'),(136,33,'1768575070_9ef963acef296c40','ls',13.924314325,3,0,0,'','2026-01-16 22:51:24.127','2026-01-16 22:52:47.651'),(137,33,'1768575070_9ef963acef296c40','ls -la',19.875913466,4,0,0,'','2026-01-16 22:51:30.079','2026-01-16 22:52:47.652'),(138,33,'1768575070_9ef963acef296c40','ls -la',42.188594392,5,0,0,'','2026-01-16 22:51:52.392','2026-01-16 22:52:47.652'),(139,33,'1768575070_9ef963acef296c40','[A test/',46.614850144,6,0,0,'','2026-01-16 22:51:56.818','2026-01-16 22:52:47.652'),(140,33,'1768575070_9ef963acef296c40','ps -ef | grep test',56.241941332,7,0,0,'','2026-01-16 22:52:06.445','2026-01-16 22:52:47.652'),(141,34,'1768575184_3c8cb16db4030013','ls -la  1.txt                               deploy_netops.sh                    mayfly-go/.ansible/                           deviops/                            .profile/home/test/dhpcd',12.30153889,1,0,0,'','2026-01-16 22:53:16.438','2026-01-16 22:59:01.927'),(142,34,'1768575184_3c8cb16db4030013','[200~kill -9 1646087[201~',55.45777552,2,1,1,'包含敏感操作关键词','2026-01-16 22:53:59.594','2026-01-16 22:59:01.927'),(143,34,'1768575184_3c8cb16db4030013','[200~crontab -u test -l',75.960351625,3,0,0,'','2026-01-16 22:54:20.096','2026-01-16 22:59:01.927'),(144,34,'1768575184_3c8cb16db4030013','[200~crontab -u test -r',123.908809969,4,0,0,'','2026-01-16 22:55:08.045','2026-01-16 22:59:01.927'),(145,34,'1768575184_3c8cb16db4030013','[200~pkill -f .dhpcd',191.606641686,5,1,1,'包含敏感操作关键词','2026-01-16 22:56:15.743','2026-01-16 22:59:01.927'),(146,34,'1768575184_3c8cb16db4030013','[200~rm -f /home/test/.dhpcd',197.95466153,6,1,1,'包含敏感操作关键词','2026-01-16 22:56:22.091','2026-01-16 22:59:01.927'),(147,34,'1768575184_3c8cb16db4030013','[200~grep \"Failed password for test\" /var/log/auth.log | tail -n 20[201~',217.780452627,7,0,0,'','2026-01-16 22:56:41.916','2026-01-16 22:59:01.927'),(148,34,'1768575184_3c8cb16db4030013','[200~usermod -s /usr/sbin/nologin test',242.445283775,8,1,1,'包含敏感操作关键词','2026-01-16 22:57:06.581','2026-01-16 22:59:01.927'),(149,34,'1768575184_3c8cb16db4030013','[200~userdel -r test',251.027024357,9,1,1,'包含敏感操作关键词','2026-01-16 22:57:15.163','2026-01-16 22:59:01.927'),(150,34,'1768575184_3c8cb16db4030013','ls',279.319796472,10,0,0,'','2026-01-16 22:57:43.456','2026-01-16 22:59:01.927'),(151,34,'1768575184_3c8cb16db4030013','cat /etc/papam.conf  pam.d/    passwd    passwd-   root@hw-ops:~# cat /etc/papam.conf  pam.d/    passwd    passwd-   root@hw-ops:~# cat /etccd /home/',293.506342841,11,1,1,'包含敏感操作关键词','2026-01-16 22:57:57.642','2026-01-16 22:59:01.928'),(152,34,'1768575184_3c8cb16db4030013','ls',294.363672269,12,0,0,'','2026-01-16 22:57:58.500','2026-01-16 22:59:01.928'),(153,34,'1768575184_3c8cb16db4030013','rm -rf test/',299.207112965,13,1,1,'包含敏感操作关键词','2026-01-16 22:58:03.343','2026-01-16 22:59:01.928'),(154,34,'1768575184_3c8cb16db4030013','ls',300.308497156,14,0,0,'','2026-01-16 22:58:04.444','2026-01-16 22:59:01.928'),(155,34,'1768575184_3c8cb16db4030013','rm -rf dingding/',311.641402161,15,1,1,'包含敏感操作关键词','2026-01-16 22:58:15.777','2026-01-16 22:59:01.928'),(156,34,'1768575184_3c8cb16db4030013','ls',312.7431588,16,0,0,'','2026-01-16 22:58:16.879','2026-01-16 22:59:01.928'),(157,34,'1768575184_3c8cb16db4030013','ls -la',315.143676122,17,0,0,'','2026-01-16 22:58:19.280','2026-01-16 22:59:01.928'),(158,34,'1768575184_3c8cb16db4030013','passwd root',323.538993301,18,1,1,'包含敏感操作关键词','2026-01-16 22:58:27.675','2026-01-16 22:59:01.928'),(159,34,'1768575184_3c8cb16db4030013','zhangfan@5391621',328.544721453,19,0,0,'','2026-01-16 22:58:32.681','2026-01-16 22:59:01.929'),(160,34,'1768575184_3c8cb16db4030013','zhangfan@123',334.618660387,20,0,0,'','2026-01-16 22:58:38.755','2026-01-16 22:59:01.929'),(161,34,'1768575184_3c8cb16db4030013','[A',337.641327698,21,0,0,'','2026-01-16 22:58:41.777','2026-01-16 22:59:01.929'),(162,34,'1768575184_3c8cb16db4030013','zhangfan@5391621',342.472208022,22,0,0,'','2026-01-16 22:58:46.608','2026-01-16 22:59:01.929'),(163,34,'1768575184_3c8cb16db4030013','zhangfan@5391621',347.419455368,23,0,0,'','2026-01-16 22:58:51.555','2026-01-16 22:59:01.929'),(164,34,'1768575184_3c8cb16db4030013','ls',352.080550706,24,0,0,'','2026-01-16 22:58:56.217','2026-01-16 22:59:01.929'),(165,35,'1768732425_c3519e27ba191582','ls',4.139905579,1,0,0,'','2026-01-18 18:33:50.134','2026-01-18 18:34:29.818'),(166,35,'1768732425_c3519e27ba191582','ls',34.805802896,2,0,0,'','2026-01-18 18:34:20.800','2026-01-18 18:34:29.818'),(167,35,'1768732425_c3519e27ba191582','pwd',36.338863367,3,0,0,'','2026-01-18 18:34:22.333','2026-01-18 18:34:29.818'),(168,36,'1768735284_a8cf4481c1cf228c','ls',2.92740767,1,0,0,'','2026-01-18 19:21:27.532','2026-01-18 19:23:34.426'),(169,36,'1768735284_a8cf4481c1cf228c','ls 1',5.865612885,2,0,0,'','2026-01-18 19:21:30.471','2026-01-18 19:23:34.427'),(170,36,'1768735284_a8cf4481c1cf228c','ls',9.327750285,3,0,0,'','2026-01-18 19:21:33.933','2026-01-18 19:23:34.428'),(171,36,'1768735284_a8cf4481c1cf228c','pwd',14.910184389,4,0,0,'','2026-01-18 19:21:39.515','2026-01-18 19:23:34.428'),(172,37,'1768739653_310088f9dc605532','ls',3.962174598,1,0,0,'','2026-01-18 20:34:17.136','2026-01-18 20:38:08.307'),(173,38,'1768740482_c8981eadeff6eb51','ls',4.468267509,1,0,0,'','2026-01-18 20:48:07.099','2026-01-18 20:49:06.698'),(174,38,'1768740482_c8981eadeff6eb51','vim 2.sh',11.85908694,2,0,0,'','2026-01-18 20:48:14.490','2026-01-18 20:49:06.699'),(175,38,'1768740482_c8981eadeff6eb51','ls',18.083710692,3,0,0,'','2026-01-18 20:48:20.715','2026-01-18 20:49:06.699'),(176,38,'1768740482_c8981eadeff6eb51','sh 2.sh',22.35491147,4,0,0,'','2026-01-18 20:48:24.986','2026-01-18 20:49:06.699'),(177,38,'1768740482_c8981eadeff6eb51','ls',25.302946835,5,0,0,'','2026-01-18 20:48:27.934','2026-01-18 20:49:06.699'),(178,38,'1768740482_c8981eadeff6eb51','rm -rf 1.sql',29.950304359,6,1,1,'包含敏感操作关键词','2026-01-18 20:48:32.581','2026-01-18 20:49:06.699'),(179,38,'1768740482_c8981eadeff6eb51','ls',31.391764929,7,0,0,'','2026-01-18 20:48:34.023','2026-01-18 20:49:06.700'),(180,38,'1768740482_c8981eadeff6eb51','systemctl status agent',43.751342835,8,0,0,'','2026-01-18 20:48:46.382','2026-01-18 20:49:06.700');
/*!40000 ALTER TABLE `sys_command_audit` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_command_blocking`
--

DROP TABLE IF EXISTS `sys_command_blocking`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_command_blocking` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `session_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `recording_id` bigint unsigned DEFAULT NULL,
  `original_cmd` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
  `resolved_cmd` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
  `blocked` tinyint(1) DEFAULT NULL,
  `block_reason` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
  `risk_level` bigint DEFAULT NULL,
  `is_alias` tinyint(1) DEFAULT '0',
  `is_script` tinyint(1) DEFAULT '0',
  `detected_issues` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
  `policy_id` bigint unsigned DEFAULT NULL,
  `policy_name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `user_id` bigint unsigned DEFAULT NULL,
  `username` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `host_id` bigint unsigned DEFAULT NULL,
  `host_ip` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_sys_command_blocking_session_id` (`session_id`) USING BTREE,
  KEY `idx_sys_command_blocking_recording_id` (`recording_id`) USING BTREE,
  KEY `idx_sys_command_blocking_blocked` (`blocked`) USING BTREE,
  KEY `idx_sys_command_blocking_risk_level` (`risk_level`) USING BTREE,
  KEY `idx_sys_command_blocking_user_id` (`user_id`) USING BTREE,
  KEY `idx_sys_command_blocking_host_id` (`host_id`) USING BTREE,
  KEY `idx_sys_command_blocking_created_at` (`created_at`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_command_blocking`
--

LOCK TABLES `sys_command_blocking` WRITE;
/*!40000 ALTER TABLE `sys_command_blocking` DISABLE KEYS */;
/*!40000 ALTER TABLE `sys_command_blocking` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_config`
--

DROP TABLE IF EXISTS `sys_config`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_config` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键''',
  `config_key` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''配置键''',
  `config_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''配置类型(ldap,email,sms等)''',
  `config_data` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''配置数据(JSON格式)''',
  `status` bigint NOT NULL DEFAULT '1' COMMENT '''状态:1->启用,2->禁用''',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''备注''',
  `create_time` datetime(3) NOT NULL COMMENT '''创建时间''',
  `update_time` datetime(3) NOT NULL COMMENT '''更新时间''',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `idx_sys_config_config_key` (`config_key`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_config`
--

LOCK TABLES `sys_config` WRITE;
/*!40000 ALTER TABLE `sys_config` DISABLE KEYS */;
INSERT INTO `sys_config` VALUES (1,'ldap','ldap','{\"enable\":true,\"host\":\"1111\",\"port\":389,\"baseDn\":\"123\",\"bindUser\":\"cn=reader,dc=dding,dc=cn\",\"bindPass\":\"123\",\"authFilter\":\"(\\u0026(cn=%s))\",\"coverAttributes\":true,\"tls\":false,\"startTLS\":false,\"defaultRoles\":null,\"defaultRoleId\":13,\"attributes\":{\"nickname\":\"cn\",\"phone\":\"mobile\",\"email\":\"mail\"},\"remark\":\"\"}',1,'','2025-12-09 13:19:56.671','2026-01-19 17:44:52.150');
/*!40000 ALTER TABLE `sys_config` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_dept`
--

DROP TABLE IF EXISTS `sys_dept`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_dept` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键',
  `parent_id` int NOT NULL COMMENT '父id',
  `dept_type` int NOT NULL COMMENT '部门类型（1->公司, 2->中心，3->部门）',
  `dept_name` varchar(30) CHARACTER SET utf8mb3 COLLATE utf8mb3_bin NOT NULL COMMENT '部门名称',
  `dept_status` int NOT NULL DEFAULT '1' COMMENT '部门状态（1->正常 2->停用）',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `dept_name` (`dept_name`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=26 DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_bin ROW_FORMAT=DYNAMIC COMMENT='部门表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_dept`
--

LOCK TABLES `sys_dept` WRITE;
/*!40000 ALTER TABLE `sys_dept` DISABLE KEYS */;
INSERT INTO `sys_dept` VALUES (1,0,1,'神舟科技有限公司',1,'2023-06-14 17:53:23'),(2,1,2,'深圳研发中心1',1,'2023-06-14 17:53:55');
/*!40000 ALTER TABLE `sys_dept` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_login_info`
--

DROP TABLE IF EXISTS `sys_login_info`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_login_info` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT 'ID',
  `username` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT '' COMMENT '用户账号',
  `ip_address` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT '' COMMENT '登录IP地址',
  `login_location` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT '' COMMENT '登录地点',
  `browser` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT '' COMMENT '浏览器类型',
  `os` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT '' COMMENT '操作系统',
  `login_status` int DEFAULT NULL COMMENT '登录状态（1-成功 2-失败）',
  `message` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin DEFAULT '' COMMENT '提示消息',
  `login_time` datetime DEFAULT NULL COMMENT '访问时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin ROW_FORMAT=DYNAMIC COMMENT='登录日志记录';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_login_info`
--

LOCK TABLES `sys_login_info` WRITE;
/*!40000 ALTER TABLE `sys_login_info` DISABLE KEYS */;
INSERT INTO `sys_login_info` VALUES (1,'admin','192.168.10.12','局域网','Chrome/145.0.0.0','Mac OS X 10_15_7',1,'登录成功','2026-02-19 15:32:48'),(2,'admin','192.168.10.12','局域网','Chrome/145.0.0.0','Mac OS X 10_15_7',1,'登录成功','2026-02-19 15:37:05'),(3,'admin','192.168.10.12','局域网','Chrome/145.0.0.0','Mac OS X 10_15_7',1,'登录成功','2026-02-19 15:45:15'),(4,'admin','192.168.10.12','局域网','Chrome/145.0.0.0','Mac OS X 10_15_7',1,'登录成功','2026-02-19 15:46:45'),(5,'admin','192.168.10.12','局域网','Chrome/145.0.0.0','Mac OS X 10_15_7',1,'登录成功','2026-02-19 15:50:25'),(6,'admin','192.168.10.12','局域网','Chrome/145.0.0.0','Mac OS X 10_15_7',1,'登录成功','2026-02-19 15:53:54'),(7,'admin','192.168.10.12','局域网','Chrome/145.0.0.0','Mac OS X 10_15_7',1,'登录成功','2026-02-19 15:56:25'),(8,'admin','192.168.10.12','局域网','Chrome/145.0.0.0','Mac OS X 10_15_7',1,'登录成功','2026-02-19 22:56:41');
/*!40000 ALTER TABLE `sys_login_info` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_menu`
--

DROP TABLE IF EXISTS `sys_menu`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_menu` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键',
  `parent_id` int DEFAULT NULL COMMENT '父级菜单id',
  `menu_name` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT '菜单名称',
  `icon` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT '图标',
  `value` varchar(100) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT '接口权限值',
  `menu_type` int DEFAULT NULL COMMENT '菜单类型：1->目录；2->菜单；3->按钮（接口绑定权限）',
  `url` varchar(200) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT '菜单url',
  `menu_status` int DEFAULT '2' COMMENT '启用状态；1->禁用；2->启用',
  `sort` int DEFAULT NULL COMMENT '排序',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=254 DEFAULT CHARSET=utf8mb3 ROW_FORMAT=DYNAMIC COMMENT='菜单表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_menu`
--

LOCK TABLES `sys_menu` WRITE;
/*!40000 ALTER TABLE `sys_menu` DISABLE KEYS */;
INSERT INTO `sys_menu` VALUES (4,0,'系统管理','StarFilled','',1,'',2,7,'2022-09-04 13:57:39'),(6,4,'用户信息','Avatar','base:admin:list',2,'system/admin',2,1,'2022-09-04 13:59:39'),(7,4,'角色信息','InfoFilled','base:role:list',2,'system/role',2,2,'2022-09-04 14:00:12'),(8,4,'菜单信息','Histogram','base:menu:list',2,'system/menu',2,3,'2022-09-04 14:00:17'),(10,4,'岗位信息','Promotion','base:post:list',2,'system/post',2,5,'2022-09-04 14:02:06'),(16,6,'新增用户','','base:admin:add',3,'',2,1,'2022-09-04 18:32:55'),(17,6,'修改用户','','base:admin:edit',3,'',2,2,'2022-09-04 18:33:29'),(18,6,'删除用户','','base:admin:delete',3,'',2,3,'2022-09-04 18:33:51'),(21,7,'新增角色','','base:role:add',3,'',2,1,'2022-09-04 18:44:47'),(22,7,'修改角色','','base:role:edit',3,'',2,2,'2022-09-04 18:45:10'),(23,7,'删除角色','','base:role:delete',3,'',2,3,'2022-09-04 18:45:46'),(24,7,'分配权限','','base:role:assign',3,'',2,4,'2022-09-04 18:46:20'),(26,8,'新增菜单','','base:menu:add',3,'',2,1,'2022-09-04 18:49:51'),(27,8,'修改菜单','','base:menu:edit',3,'',2,2,'2022-09-04 18:50:24'),(28,8,'删除菜单','','base:menu:delete',3,'',2,3,'2022-09-04 18:50:53'),(32,10,'新增岗位','','base:post:add',3,'',2,1,'2022-09-04 18:53:28'),(33,10,'修改岗位','','base:post:edit',3,'',2,2,'2022-09-04 18:53:48'),(34,10,'删除岗位','','base:post:delete',3,'',2,3,'2022-09-04 18:54:00'),(44,0,'操作审计','BellFilled','',1,'',2,9,'2022-09-05 11:06:57'),(45,44,'操作日志','User','monitor:operator:list',2,'monitor/operator',2,1,'2022-09-05 11:10:54'),(46,44,'登录日志','DocumentRemove','monitor:loginLog:list',2,'monitor/loginlog',2,2,'2022-09-05 11:11:31'),(47,45,'清空操作日志','','monitor:operator:clean',3,'',2,1,'2022-09-05 11:12:36'),(49,46,'清空登录日志','','monitor:loginLog:clean',3,'',2,1,'2022-09-05 11:16:01'),(60,6,'重置密码',NULL,'base:admin:reset',3,NULL,2,6,'2022-12-01 16:33:34'),(62,46,'删除登录日志','','monitor:loginLog:delete',3,'',2,2,'2022-12-02 15:41:56'),(72,0,'仪表盘','HomeFilled','',1,'dashboard',2,1,'2023-05-24 22:11:13'),(73,45,'删除操作日志','','monitor:operator:delete',3,'',2,3,'2023-06-02 10:09:38'),(78,80,'主机管理','Platform','cmdb:ecs:list',2,'cmdb/ecs',2,1,'2025-06-29 00:30:35'),(80,0,'资产管理','TrendCharts','',1,'',2,2,'2025-07-03 11:47:07'),(81,0,'容器管理','UploadFilled','',1,'',2,3,'2025-07-03 11:50:47'),(82,81,'集群管理','Menu','cloud:k8s:list',2,'k8s/list',2,1,'2025-07-03 11:56:44'),(83,81,'节点管理','Help','cloud:k8s:node',2,'k8s/node',2,2,'2025-07-03 12:04:59'),(84,0,'配置中心','Tools','',1,'',2,8,'2025-07-04 17:00:01'),(85,84,'主机凭据','Setting','config:ecs:key',2,'config/ecskey',2,1,'2025-07-04 17:03:10'),(86,84,'通用凭据','User','config:accountauth:key',2,'config/accountauth',2,2,'2025-07-04 17:08:20'),(88,80,'资产分组','Shop','cmdb:group',2,'cmdb/group',2,2,'2025-07-16 15:17:14'),(89,88,'创建资产','','cmdb:group:add',3,'',2,1,'2025-07-18 15:24:31'),(90,88,'修改资产','','cmdb:group:update',3,'',2,2,'2025-07-18 15:25:49'),(91,88,'删除资产','','cmdb:group:delete',3,'',2,3,'2025-07-18 15:26:21'),(93,81,'工作负载','Star','cloud:k8s:workload',2,'k8s/workload',2,4,'2025-07-24 14:38:31'),(95,80,'数据管理','Coin','cmdb:db',2,'cmdb/db',2,3,'2025-07-29 15:27:50'),(96,44,'数据日志','Coin','monitor:dblog:list',2,'monitor/dblog',2,3,'2025-07-31 12:44:07'),(97,0,'任务中心','User','',1,'',2,5,'2025-08-06 13:33:47'),(98,97,'任务模版','connection','task:template',2,'task/template',2,2,'2025-08-06 13:35:19'),(99,97,'任务作业','key','task:job',2,'task/job',2,1,'2025-08-06 13:36:06'),(100,97,'Ansible任务','Eleme','task:ansible',2,'task/ansible',2,3,'2025-08-23 18:35:24'),(101,0,'运维工具','Search','',1,'',2,6,'2025-08-29 10:59:35'),(102,101,'agent列表','price-tag','ops:agent',2,'ops/agent',2,2,'2025-08-29 11:22:20'),(103,101,'工具列表','present','ops:tools',2,'ops/tools',2,1,'2025-08-29 11:29:02'),(104,84,'密钥管理','Phone','config:keymanage:key',2,'config/keymanage',2,3,'2025-09-08 13:24:40'),(105,81,'命名空间','discount','k8s:namespace',2,'k8s/namespace',2,3,'2025-09-11 15:02:14'),(106,81,'网络管理','guide','k8s:network',2,'k8s/network',2,5,'2025-09-11 15:04:14'),(107,81,'配置管理','connection','k8s:config',2,'k8s/config',2,7,'2025-09-11 15:04:52'),(108,81,'存储管理','Coin','k8s:storage',2,'k8s/storage',2,6,'2025-09-11 15:05:40'),(109,0,'服务管理','ElemeFilled','',1,'',2,4,'2025-09-16 09:49:55'),(110,109,'应用列表','Menu','app:application',2,'app/application',2,2,'2025-09-16 09:52:58'),(111,109,'快速发布','TrendCharts','app:quick-release',2,'app/quick-release',2,3,'2025-09-16 17:12:11'),(113,45,'批量删除','','monitor:operator:delete',3,'',2,2,'2025-09-17 20:55:13'),(115,104,'删除密钥','','config:keymanage:delete',3,'',2,2,'2025-09-18 10:53:44'),(117,104,'同步主机','','config:keymanage:rsync',3,'',2,3,'2025-09-18 10:57:25'),(118,104,'创建密钥','','config:keymanage:create',3,'',2,4,'2025-09-18 11:01:12'),(119,86,'修改账号','','config:common:edit',3,'',2,2,'2025-09-18 11:47:33'),(120,86,'删除账号','','config:common:delete',3,'',2,3,'2025-09-18 11:48:17'),(121,86,'解密账号','','config:common:decrypt',3,'',2,4,'2025-09-18 11:48:57'),(122,86,'创建账号','','config:common:add',3,'',2,1,'2025-09-18 11:49:30'),(123,85,'修改凭据','','config:ecs:edit',3,'',2,1,'2025-09-18 11:54:16'),(124,85,'删除凭据','','config:ecs:delete',3,'',2,2,'2025-09-18 11:54:51'),(125,85,'创建凭据','','config:ecs:create',3,'',2,3,'2025-09-18 11:55:21'),(126,102,'卸载agent','','ops:agent:delete',3,'',2,1,'2025-09-18 12:47:48'),(127,102,'查看agent','','ops:agent:get',3,'',2,2,'2025-09-18 12:49:02'),(128,102,'部署agent','','ops:agent:create',3,'',2,3,'2025-09-18 12:49:56'),(129,102,'批量卸载agent','','ops:agent:deleteall',3,'',2,4,'2025-09-18 12:50:52'),(130,100,'启动ansible任务流程','','task:ansible:start',3,'',2,1,'2025-09-18 12:59:30'),(131,100,'删除ansible任务','','task:ansible:delete',3,'',2,2,'2025-09-18 13:00:03'),(132,100,'新增ansible任务','','task:ansible:create',3,'',2,3,'2025-09-18 13:00:45'),(133,98,'新增模版','','task:template:add',3,'',2,1,'2025-09-18 13:16:38'),(134,98,'修改模版','','task:template:edit',3,'',2,2,'2025-09-18 13:17:04'),(135,98,'删除模版','','task:template:delete',3,'',2,3,'2025-09-18 13:18:25'),(136,99,'新增任务','','task:job:add',3,'',2,1,'2025-09-18 13:24:19'),(137,99,'启动任务','','task:job:start',3,'',2,2,'2025-09-18 13:24:59'),(138,99,'删除任务','','task:job:delete',3,'',2,3,'2025-09-18 13:25:41'),(139,111,'新建发布','','app:quick-release:add',3,'',2,1,'2025-09-18 13:30:53'),(140,111,'启动发布','','app:quick-release:start',3,'',2,2,'2025-09-18 13:32:11'),(141,111,'删除发布','','app:quick-release:delete',3,'',2,3,'2025-09-18 13:32:32'),(142,110,'创建应用','','app:application:add',3,'',2,1,'2025-09-18 14:28:07'),(143,110,'修改应用','','app:application:edit',3,'',2,2,'2025-09-18 14:28:59'),(144,110,'配置应用环境','','app:application:env',3,'',2,3,'2025-09-18 14:29:34'),(145,110,'删除应用','','app:application:delete',3,'',2,4,'2025-09-18 14:30:11'),(146,95,'创建数据库账号','','cmdb:db:add',3,'',2,1,'2025-09-18 14:41:32'),(147,95,'修改数据库配置','','cmdb:db:edit',3,'',2,2,'2025-09-18 14:42:47'),(148,95,'删除数据库账号','','cmdb:db:delete',3,'',2,3,'2025-09-18 14:43:57'),(149,78,'创建主机','','cmdb:ecs:add',3,'',2,1,'2025-09-18 14:47:42'),(150,78,'主机终端','','cmdb:ecs:terminal',3,'',2,2,'2025-09-18 14:48:36'),(151,78,'修改主机信息','','cmdb:ecs:edit',3,'',2,3,'2025-09-18 14:49:43'),(152,78,'上传文件到主机','','cmdb:ecs:upload',3,'',2,4,'2025-09-18 14:50:38'),(153,78,'执行主机命令','','cmdb:ecs:shell',3,'',2,5,'2025-09-18 14:51:10'),(154,78,'监控主机','','cmdb:ecs:monitor',3,'',2,6,'2025-09-18 14:51:52'),(155,78,'删除主机','','cmdb:ecs:delete',3,'',2,7,'2025-09-18 14:52:20'),(156,99,'启动脚本','','task:job:jobstart',3,'',2,4,'2025-09-18 18:36:38'),(157,99,'停止脚本','','task:job:jobstop',3,'',2,5,'2025-09-18 18:39:23'),(159,100,'开始ansible任务作业','','task:ansible:jobstart',3,'',2,4,'2025-09-18 18:43:40'),(160,111,'启动jenkins任务','','app:quick-release:jobstart',3,'',2,4,'2025-09-18 18:47:39'),(161,111,'停止jenkins任务','','app:quick-release:jobstop',3,'',2,5,'2025-09-18 18:48:16'),(162,110,'新增环境','','app:application:envadd',3,'',2,5,'2025-09-18 21:02:28'),(163,110,'修改环境','','app:application:envedit',3,'',2,6,'2025-09-18 21:03:08'),(164,110,'删除环境','','app:application:envdelete',3,'',2,7,'2025-09-18 21:04:43'),(165,78,'连接主机终端','','cmdb:ecs:connecthost',3,'',2,8,'2025-09-18 21:11:43'),(166,78,'同步主机信息','','cmdb:ecs:rsync',3,'',2,9,'2025-09-19 21:35:06'),(167,82,'注册集群','','cloud:k8s:register',3,'',2,1,'2025-09-19 21:57:54'),(168,82,'创建集群','','cloud:k8s:add',3,'',2,2,'2025-09-19 21:58:19'),(169,82,'修改集群','','cloud:k8s:edit',3,'',2,3,'2025-09-19 21:59:06'),(170,82,'同步集群','','cloud:k8s:rsync',3,'',2,4,'2025-09-19 21:59:31'),(171,82,'删除集群','','cloud:k8s:delete',3,'',2,5,'2025-09-19 21:59:56'),(172,83,'查看监控仪表盘','','k8s:node:monitor',3,'',2,1,'2025-09-20 00:19:49'),(173,83,'查看节点资源详情','','k8s:node:details',3,'',2,2,'2025-09-20 00:21:20'),(174,83,'配置节点污点','','k8s:node:stain',3,'',2,3,'2025-09-20 00:22:17'),(175,83,'增加标签','','k8s:node:label',3,'',2,4,'2025-09-20 00:23:15'),(176,83,'封锁节点','','k8s:node:close',3,'',2,5,'2025-09-20 00:24:13'),(177,83,'驱逐节点','','k8s:node:expel',3,'',2,6,'2025-09-20 00:25:04'),(178,105,'创建命名空间','','k8s:namespace:add',3,'',2,1,'2025-09-20 00:36:14'),(179,105,'查看命名空间详情','','k8s:namespace:details',3,'',2,2,'2025-09-20 00:37:22'),(180,105,'查看命名空间资源配置','','k8s:namespace:setup',3,'',2,3,'2025-09-20 00:39:05'),(181,105,'添加命名空间资源配置','','k8s:namespace:setupadd',3,'',2,4,'2025-09-20 00:40:06'),(182,105,'查看限制命名空间资源','','k8s:namespace:restriction',3,'',2,5,'2025-09-20 00:41:25'),(183,105,'添加限制命名空间资源','','k8s:namespace:restrictionadd',3,'',2,6,'2025-09-20 00:42:23'),(184,105,'删除命名空间','','k8s:namespace:delete',3,'',2,7,'2025-09-20 00:43:03'),(185,93,'新增工作负载','','k8s:workload:add',3,'',2,1,'2025-09-20 01:05:08'),(186,93,'伸缩pod','','k8s:workload:expandable',3,'',2,2,'2025-09-20 01:06:18'),(187,93,'重启pod','','k8s:workload:restart',3,'',2,3,'2025-09-20 01:07:13'),(188,93,'更新pod资源限制','','k8s:workload:resource',3,'',2,4,'2025-09-20 01:08:39'),(189,93,'更新pod资调度','','k8s:workload:dispatch',3,'',2,5,'2025-09-20 01:09:36'),(190,93,'更新yaml文件','','k8s:workload:edityaml',3,'',2,6,'2025-09-20 01:10:57'),(191,93,'删除工作负载','','k8s:workload:delete',3,'',2,7,'2025-09-20 01:11:52'),(192,93,'回滚工作负载版本','','k8s:workload:rollback_version',3,'',2,8,'2025-09-20 01:39:38'),(193,93,'查看pod日志','','k8s:workload:podlog',3,'',2,9,'2025-09-20 01:44:37'),(194,93,'删除pod','','k8s:workload:poddelete',3,'',2,10,'2025-09-20 01:45:16'),(195,93,'登陆pod终端','','k8s:workload:terminal',3,'',2,11,'2025-09-20 01:46:07'),(196,93,'编辑pod yaml文件','','k8s:workload:edityaml',3,'',2,12,'2025-09-20 01:47:22'),(197,106,'新增service','','k8s:network:addservice',3,'',2,1,'2025-09-20 02:14:21'),(198,106,'编辑 Service','','k8s:network:editservice',3,'',2,2,'2025-09-20 02:15:19'),(199,106,'编辑service YAML','','k8s:network:edit_service_yaml',3,'',2,3,'2025-09-20 02:16:15'),(200,106,'查看Service 事件','','k8s:network:service_event',3,'',2,4,'2025-09-20 02:18:11'),(201,106,'删除Service','','k8s:network:deleteservice',3,'',2,5,'2025-09-20 02:18:59'),(202,106,'编辑ingress','','k8s:network:editingress',3,'',2,6,'2025-09-20 02:26:59'),(203,106,'新增ingress','','k8s:network:addingress',3,'',2,7,'2025-09-20 02:27:29'),(204,106,'编辑ingress_yaml','','k8s:network:edit_ingress_yaml',3,'',2,8,'2025-09-20 02:28:23'),(205,106,'查看ingress 事件','','k8s:network:ingress_event',3,'',2,9,'2025-09-20 02:29:24'),(206,106,'删除ingress','','k8s:network:delete_ingress',3,'',2,10,'2025-09-20 02:30:04'),(215,4,'部门信息','List','system:dept',2,'system/dept',2,6,'2025-12-09 11:03:54'),(217,95,'查看密码','','cmdb:db:passwd',3,'',2,4,'2025-12-13 14:27:22'),(218,95,'测试链接','','cmdb:db:test',3,'',2,5,'2025-12-13 14:28:38'),(219,95,'DBMS终端','','cmdb:db:dbms',3,'',2,6,'2025-12-13 14:32:59'),(220,95,'Redis终端','','cmdb:db:redis',3,'',2,7,'2025-12-13 14:33:27'),(221,95,'ES终端','','cmdb:db:es',3,'',2,8,'2025-12-13 14:33:53'),(222,95,'Mongo终端','','cmdb:db:mongo',3,'',2,9,'2025-12-13 14:34:14'),(241,78,'下载文件','','cmdb:ecs:download',3,'',2,10,'2025-12-13 16:13:17'),(243,78,'批量操作','','cmdb:ecs:batch',3,'',2,10,'2025-12-26 19:04:56'),(244,104,'查看密钥','','config:keymanage:decrypt',3,'',2,5,'2025-12-27 00:20:50'),(245,88,'监控大盘','','cmdb:group:collect',3,'',2,4,'2026-01-01 15:48:05'),(246,88,'终端按钮','','cmdb:snmp:ssh',3,'',2,5,'2026-01-01 22:18:54'),(247,44,'会话录制','Coin','monitor:recording',2,'monitor/recording',2,4,'2026-01-08 19:57:24'),(250,86,'测试连接','','config:common:test',3,'',2,5,'2026-01-18 23:27:29'),(251,215,'新增部门','','base:dept:add',3,'',2,1,'2026-02-19 15:52:17'),(252,215,'修改部门','','base:dept:edit',3,'',2,2,'2026-02-19 15:52:58'),(253,215,'删除部门','','base:dept:delete',3,'',2,3,'2026-02-19 15:53:27');
/*!40000 ALTER TABLE `sys_menu` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_operation_log`
--

DROP TABLE IF EXISTS `sys_operation_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_operation_log` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键',
  `admin_id` bigint unsigned NOT NULL COMMENT '''管理员id''',
  `username` varchar(64) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL COMMENT '''管理员账号''',
  `method` varchar(64) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL COMMENT '''请求方式''',
  `ip` varchar(64) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT '''IP''',
  `url` varchar(500) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT '''URL''',
  `create_time` datetime(3) NOT NULL COMMENT '''创建时间''',
  `description` varchar(255) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT '''操作描述''',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `username` (`username`) USING BTREE,
  KEY `create_time` (`create_time`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb3 ROW_FORMAT=DYNAMIC COMMENT='操作日志记录表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_operation_log`
--

LOCK TABLES `sys_operation_log` WRITE;
/*!40000 ALTER TABLE `sys_operation_log` DISABLE KEYS */;
INSERT INTO `sys_operation_log` VALUES (1,89,'admin','delete','192.168.10.12','/api/v1/sysOperationLog/clean','2026-02-19 22:59:04.635','清空操作日志'),(2,89,'admin','delete','192.168.10.12','/api/v1/monitor/agent/delete/118','2026-02-19 22:59:11.679','删除Agent'),(3,89,'admin','post','192.168.10.12','/api/v1/cmdb/hostcreate','2026-02-19 22:59:43.755','创建主机'),(4,89,'admin','post','192.168.10.12','/api/v1/monitor/agent/deploy','2026-02-19 23:01:12.494','部署Agent'),(5,89,'admin','delete','192.168.10.12','/api/v1/cmdb/hostdelete','2026-02-19 23:20:06.164','删除主机'),(6,89,'admin','delete','192.168.10.12','/api/v1/monitor/agent/delete/119','2026-02-19 23:20:15.587','删除Agent');
/*!40000 ALTER TABLE `sys_operation_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_post`
--

DROP TABLE IF EXISTS `sys_post`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_post` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '岗位ID',
  `post_code` varchar(64) CHARACTER SET utf8mb3 COLLATE utf8mb3_bin NOT NULL COMMENT '岗位编码',
  `post_name` varchar(50) CHARACTER SET utf8mb3 COLLATE utf8mb3_bin NOT NULL COMMENT '岗位名称',
  `post_status` int NOT NULL DEFAULT '1' COMMENT '状态（1->正常 2->停用）',
  `create_time` datetime DEFAULT NULL COMMENT '创建时间',
  `remark` varchar(500) CHARACTER SET utf8mb3 COLLATE utf8mb3_bin DEFAULT NULL COMMENT '备注',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_bin ROW_FORMAT=DYNAMIC COMMENT='岗位信息表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_post`
--

LOCK TABLES `sys_post` WRITE;
/*!40000 ALTER TABLE `sys_post` DISABLE KEYS */;
INSERT INTO `sys_post` VALUES (1,'AAA','研发总监',1,'2023-06-14 20:08:22','主管各个部门'),(10,'ops','运维工程师',1,'2025-06-28 22:46:33','运维工程师'),(11,'dev','研发工程师',1,'2025-06-28 22:50:29','研发工程师'),(12,'test','测试工程师',1,'2025-06-28 22:52:57','测试工程师'),(14,'cccc','运维总监',1,'2026-01-14 11:12:04',''),(15,'bbb','测试总监',1,'2026-01-14 11:12:44','');
/*!40000 ALTER TABLE `sys_post` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_role`
--

DROP TABLE IF EXISTS `sys_role`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_role` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '主键',
  `role_name` varchar(64) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL COMMENT '角色名称',
  `role_key` varchar(64) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci NOT NULL COMMENT '角色权限字符串',
  `status` int NOT NULL DEFAULT '1' COMMENT '启用状态：1->启用；2->禁用',
  `description` varchar(500) CHARACTER SET utf8mb3 COLLATE utf8mb3_general_ci DEFAULT NULL COMMENT '描述',
  `create_time` datetime NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `name` (`role_name`) USING BTREE,
  KEY `create_time` (`create_time`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb3 ROW_FORMAT=DYNAMIC COMMENT='后台角色表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_role`
--

LOCK TABLES `sys_role` WRITE;
/*!40000 ALTER TABLE `sys_role` DISABLE KEYS */;
INSERT INTO `sys_role` VALUES (1,'超级管理员','admin',1,'最大权限','2023-06-12 20:04:53'),(13,'游客','test',1,'test1','2025-07-03 18:47:25');
/*!40000 ALTER TABLE `sys_role` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_role_menu`
--

DROP TABLE IF EXISTS `sys_role_menu`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_role_menu` (
  `role_id` int DEFAULT NULL COMMENT '角色ID',
  `menu_id` int DEFAULT NULL COMMENT '菜单ID'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb3 ROW_FORMAT=DYNAMIC COMMENT='角色和菜单关系表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_role_menu`
--

LOCK TABLES `sys_role_menu` WRITE;
/*!40000 ALTER TABLE `sys_role_menu` DISABLE KEYS */;
INSERT INTO `sys_role_menu` VALUES (1,72),(1,80),(1,78),(1,149),(1,150),(1,151),(1,152),(1,153),(1,154),(1,155),(1,165),(1,166),(1,241),(1,243),(1,88),(1,89),(1,90),(1,91),(1,245),(1,246),(1,95),(1,146),(1,147),(1,148),(1,217),(1,218),(1,219),(1,220),(1,221),(1,222),(1,81),(1,82),(1,167),(1,168),(1,169),(1,170),(1,171),(1,83),(1,172),(1,173),(1,174),(1,175),(1,176),(1,177),(1,105),(1,178),(1,179),(1,180),(1,181),(1,182),(1,183),(1,184),(1,93),(1,185),(1,186),(1,187),(1,188),(1,189),(1,190),(1,191),(1,192),(1,193),(1,194),(1,195),(1,196),(1,106),(1,197),(1,198),(1,199),(1,200),(1,201),(1,202),(1,203),(1,204),(1,205),(1,206),(1,108),(1,107),(1,109),(1,110),(1,142),(1,143),(1,144),(1,145),(1,162),(1,163),(1,164),(1,111),(1,139),(1,140),(1,141),(1,160),(1,161),(1,97),(1,99),(1,136),(1,137),(1,138),(1,156),(1,157),(1,98),(1,133),(1,134),(1,135),(1,100),(1,130),(1,131),(1,132),(1,159),(1,101),(1,103),(1,102),(1,126),(1,127),(1,128),(1,129),(1,4),(1,6),(1,16),(1,17),(1,18),(1,60),(1,7),(1,21),(1,22),(1,23),(1,24),(1,8),(1,26),(1,27),(1,28),(1,10),(1,32),(1,33),(1,34),(1,215),(1,251),(1,252),(1,253),(1,84),(1,85),(1,123),(1,124),(1,125),(1,86),(1,122),(1,119),(1,120),(1,121),(1,250),(1,104),(1,115),(1,117),(1,118),(1,244),(1,45),(1,47),(1,113),(1,73),(1,46),(1,49),(1,62),(1,44),(13,72),(13,149),(13,150),(13,154),(13,165),(13,89),(13,146),(13,218),(13,219),(13,167),(13,168),(13,170),(13,172),(13,173),(13,174),(13,175),(13,178),(13,179),(13,180),(13,181),(13,182),(13,185),(13,186),(13,190),(13,193),(13,195),(13,197),(13,200),(13,203),(13,205),(13,108),(13,107),(13,142),(13,144),(13,162),(13,139),(13,140),(13,160),(13,136),(13,137),(13,156),(13,133),(13,130),(13,132),(13,159),(13,103),(13,127),(13,16),(13,125),(13,122),(13,118),(13,73),(13,62),(13,80),(13,78),(13,88),(13,95),(13,81),(13,82),(13,83),(13,105),(13,93),(13,106),(13,109),(13,110),(13,111),(13,97),(13,99),(13,98),(13,100),(13,101),(13,102),(13,4),(13,6),(13,84),(13,85),(13,86),(13,104),(13,44),(13,45),(13,46);
/*!40000 ALTER TABLE `sys_role_menu` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `sys_session_recording`
--

DROP TABLE IF EXISTS `sys_session_recording`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sys_session_recording` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键ID''',
  `session_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''会话唯一标识''',
  `admin_id` bigint unsigned NOT NULL COMMENT '''操作用户ID''',
  `username` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''用户名''',
  `host_id` bigint unsigned NOT NULL COMMENT '''目标主机ID''',
  `host_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''主机名称''',
  `host_ip` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''主机IP''',
  `ssh_user` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''SSH登录用户''',
  `start_time` datetime(3) NOT NULL COMMENT '''会话开始时间''',
  `end_time` datetime(3) DEFAULT NULL COMMENT '''会话结束时间''',
  `duration` bigint DEFAULT NULL COMMENT '''会话时长(秒)''',
  `terminal_width` bigint DEFAULT '80' COMMENT '''终端宽度''',
  `terminal_height` bigint DEFAULT '24' COMMENT '''终端高度''',
  `file_path` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''录制文件路径''',
  `file_size` bigint DEFAULT NULL COMMENT '''文件大小(字节)''',
  `storage_type` bigint DEFAULT '1' COMMENT '''存储类型:1-本地 2-OSS''',
  `oss_key` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''OSS对象key''',
  `input_count` bigint DEFAULT '0' COMMENT '''输入事件数量''',
  `output_count` bigint DEFAULT '0' COMMENT '''输出事件数量''',
  `resize_count` bigint DEFAULT '0' COMMENT '''窗口调整次数''',
  `command_count` bigint DEFAULT '0' COMMENT '''命令执行次数(估算)''',
  `client_ip` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''客户端IP''',
  `user_agent` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''用户代理''',
  `risk_level` bigint DEFAULT '0' COMMENT '''风险等级:0-正常 1-可疑 2-高危''',
  `has_sensitive_cmd` tinyint(1) DEFAULT '0' COMMENT '''是否包含敏感命令''',
  `status` bigint DEFAULT '1' COMMENT '''状态:1-录制中 2-已完成 3-异常终止''',
  `error_msg` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''错误信息''',
  `create_time` datetime(3) NOT NULL COMMENT '''创建时间''',
  `update_time` datetime(3) DEFAULT NULL COMMENT '''更新时间''',
  `delete_time` datetime(3) DEFAULT NULL COMMENT '''删除时间''',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `idx_sys_session_recording_session_id` (`session_id`) USING BTREE,
  KEY `idx_sys_session_recording_host_id` (`host_id`) USING BTREE,
  KEY `idx_sys_session_recording_host_ip` (`host_ip`) USING BTREE,
  KEY `idx_sys_session_recording_start_time` (`start_time`) USING BTREE,
  KEY `idx_sys_session_recording_risk_level` (`risk_level`) USING BTREE,
  KEY `idx_sys_session_recording_status` (`status`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=39 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `sys_session_recording`
--

LOCK TABLES `sys_session_recording` WRITE;
/*!40000 ALTER TABLE `sys_session_recording` DISABLE KEYS */;
/*!40000 ALTER TABLE `sys_session_recording` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `task_ansible`
--

DROP TABLE IF EXISTS `task_ansible`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `task_ansible` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键ID''',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''任务名称''',
  `description` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''任务描述''',
  `type` bigint NOT NULL DEFAULT '1' COMMENT '''任务类型:1-手动,2-Git,3-K8s''',
  `git_repo` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''Git仓库地址''',
  `host_groups` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''主机分组JSON''',
  `all_host_ids` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''所有主机ID JSON数组''',
  `global_vars` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''全局变量JSON''',
  `status` bigint NOT NULL DEFAULT '1' COMMENT '''任务状态:1-等待中,2-运行中,3-成功,4-异常''',
  `created_at` datetime(3) NOT NULL COMMENT '''创建时间''',
  `updated_at` datetime(3) NOT NULL COMMENT '''更新时间''',
  `error_msg` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''错误信息''',
  `task_count` bigint NOT NULL DEFAULT '0' COMMENT '''任务数量(Type=1时为上传文件数,Type=2时为解析的playbook数,Type=3时固定为1)''',
  `total_duration` bigint NOT NULL DEFAULT '0' COMMENT '''任务执行总耗时(秒,所有子任务耗时总和)''',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `idx_task_ansible_name` (`name`) USING BTREE,
  KEY `idx_task_status` (`status`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=52 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `task_ansible`
--

LOCK TABLES `task_ansible` WRITE;
/*!40000 ALTER TABLE `task_ansible` DISABLE KEYS */;
INSERT INTO `task_ansible` VALUES (49,'test2','',1,'','{\"web\":[511]}','[511]','',3,'2025-11-24 00:12:22.770','2025-11-24 00:12:47.140','',1,17);
/*!40000 ALTER TABLE `task_ansible` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `task_ansiblework`
--

DROP TABLE IF EXISTS `task_ansiblework`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `task_ansiblework` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键ID''',
  `task_id` bigint unsigned NOT NULL COMMENT '''父任务ID''',
  `entry_file_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''入口文件名''',
  `entry_file_path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''入口文件路径''',
  `log_path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '''日志路径''',
  `status` bigint NOT NULL DEFAULT '1' COMMENT '''子任务状态:1-等待中,2-运行中,3-成功,4-异常''',
  `start_time` datetime(3) DEFAULT NULL COMMENT '''开始时间''',
  `end_time` datetime(3) DEFAULT NULL COMMENT '''结束时间''',
  `duration` bigint DEFAULT NULL COMMENT '''执行耗时(秒)''',
  `exit_code` bigint DEFAULT NULL COMMENT '''退出代码''',
  `error_msg` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''错误信息''',
  `log` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''日志内容''',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_task_ansiblework_task_id` (`task_id`) USING BTREE,
  KEY `idx_task_id` (`task_id`) USING BTREE,
  KEY `idx_task_work_composite` (`task_id`,`status`) USING BTREE,
  CONSTRAINT `fk_task_ansible_works` FOREIGN KEY (`task_id`) REFERENCES `task_ansible` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE=InnoDB AUTO_INCREMENT=53 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `task_ansiblework`
--

LOCK TABLES `task_ansiblework` WRITE;
/*!40000 ALTER TABLE `task_ansiblework` DISABLE KEYS */;
INSERT INTO `task_ansiblework` VALUES (49,49,'01-linux-os.yaml.yml','task/49/test2/01-linux-os.yaml.yml','logs/ansible/49/49/01-linux-os.yaml.yml.log',3,'2025-11-24 00:12:28.894','2025-11-24 00:12:46.777',17,0,'',NULL);
/*!40000 ALTER TABLE `task_ansiblework` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `task_job`
--

DROP TABLE IF EXISTS `task_job`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `task_job` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '任务标题',
  `type` bigint DEFAULT NULL COMMENT '任务类型 1=普通任务,2=定时任务,3=ansible任务',
  `shell` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '任务内容(任务模板ID,多个用逗号分隔)',
  `host_ids` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '主机ID(多个用逗号分隔)',
  `cron_expr` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '定时表达式(* * * * *)',
  `tasklog` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '任务执行日志',
  `status` bigint DEFAULT NULL COMMENT '任务状态 1=等待中,2=运行中,3=成功,4=异常,5=已暂停',
  `duration` bigint DEFAULT NULL COMMENT '执行耗时(秒)',
  `remark` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '任务备注',
  `start_time` datetime(3) DEFAULT NULL COMMENT '任务开始时间',
  `end_time` datetime(3) DEFAULT NULL COMMENT '任务结束时间',
  `created_at` datetime(3) DEFAULT NULL COMMENT '任务创建时间',
  `task_count` bigint DEFAULT NULL COMMENT '任务数量',
  `is_recurring` tinyint(1) DEFAULT NULL COMMENT '是否周期性任务',
  `scheduled_time` datetime(3) DEFAULT NULL COMMENT '计划执行时间',
  `log_path` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '日志文件路径',
  `execute_count` bigint DEFAULT '0' COMMENT '执行次数',
  `next_run_time` datetime(3) DEFAULT NULL COMMENT '下次执行时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=62 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `task_job`
--

LOCK TABLES `task_job` WRITE;
/*!40000 ALTER TABLE `task_job` DISABLE KEYS */;
INSERT INTO `task_job` VALUES (49,'test4',1,'12','1','','',3,5,'',NULL,'2025-09-29 21:15:32.496','2025-09-29 21:01:51.867',1,NULL,NULL,NULL,1,NULL);
/*!40000 ALTER TABLE `task_job` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `task_template`
--

DROP TABLE IF EXISTS `task_template`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `task_template` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `type` bigint NOT NULL,
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `created_by` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `updated_by` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `task_template`
--

LOCK TABLES `task_template` WRITE;
/*!40000 ALTER TABLE `task_template` DISABLE KEYS */;
INSERT INTO `task_template` VALUES (2,'1-数字 1-100',1,'#!/bin/bash\n\n# 从 1 到 100，每秒打印一个数字\nfor ((i = 1; i <= 100; i++)); do\n    echo \"[$(date +%H:%M:%S)] $i\"\n    sleep 1\ndone\n\necho \"完成：所有数字 1-100 已打印完毕。\"\n','2025-08-06 12:47:57.073','2025-08-12 16:14:49.394','admin','admin','测试脚本');
/*!40000 ALTER TABLE `task_template` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `task_work`
--

DROP TABLE IF EXISTS `task_work`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `task_work` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `task_id` bigint unsigned DEFAULT NULL COMMENT '关联的任务ID',
  `template_id` bigint unsigned DEFAULT NULL COMMENT '任务模板ID',
  `host_id` bigint unsigned DEFAULT NULL COMMENT '执行主机ID',
  `status` bigint DEFAULT NULL COMMENT '任务状态 1=等待中,2=运行中,3=成功,4=异常',
  `log` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '任务日志',
  `log_path` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '日志文件路径',
  `start_time` datetime(3) DEFAULT NULL COMMENT '任务开始时间',
  `end_time` datetime(3) DEFAULT NULL COMMENT '任务结束时间',
  `duration` bigint DEFAULT NULL COMMENT '执行耗时(秒)',
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `type` bigint DEFAULT NULL COMMENT '任务类型 1=普通任务,2=定时任务',
  `scheduled_time` datetime(3) DEFAULT NULL COMMENT '定时任务执行时间',
  `cron_expr` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT 'cron表达式',
  `is_recurring` tinyint(1) DEFAULT NULL COMMENT '是否周期性任务',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `idx_task_work_task_id` (`task_id`) USING BTREE,
  KEY `idx_task_work_template_id` (`template_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=110 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `task_work`
--

LOCK TABLES `task_work` WRITE;
/*!40000 ALTER TABLE `task_work` DISABLE KEYS */;
INSERT INTO `task_work` VALUES (90,48,11,1,2,'[2025-09-29-21:42:02] 任务开始\n进程统计信息如下\n总进程数量为:79\nRunning 进程数为:1\nStoped 进程数为:0\nSleeping 进程数为:55\nZombie 进程数为:0\n[2025-09-29-21:42:02] 任务完成\n','logs/task_48/task_48_template_11.log','2025-09-29 21:42:00.299','2025-09-29 21:42:06.276',5,'2025-09-29 20:58:17.259',0,'2025-09-29 21:45:00.000',NULL,NULL),(91,48,0,0,4,'','',NULL,NULL,0,'2025-09-29 20:58:17.515',2,'2025-09-29 21:00:00.000',NULL,NULL),(92,49,12,1,3,'[2025-09-29-21:15:27] 任务开始\n协议:tcp 端口号:6379\n协议:tcp 端口号:4330\n协议:tcp 端口号:3306\n协议:tcp 端口号:80\n协议:tcp 端口号:33060\n协议:tcp 端口号:8088\n协议:tcp 端口号:8080\n协议:tcp 端口号:4330\n协议:tcp 端口号:80\n协议:tcp 端口号:8080\n协议:tcp 端口号:8086\n协议:tcp 端口号:9091\n协议:tcp 端口号:9090\n[2025-09-29-21:15:27] 任务完成\n','logs/task_49/task_49_template_12.log','2025-09-29 21:15:25.761','2025-09-29 21:15:31.262',5,'2025-09-29 21:01:51.970',0,NULL,NULL,NULL);
/*!40000 ALTER TABLE `task_work` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `tool_link`
--

DROP TABLE IF EXISTS `tool_link`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `tool_link` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键''',
  `title` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''导航标题''',
  `icon` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''导航图标''',
  `link` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''链接地址''',
  `sort` bigint DEFAULT '0' COMMENT '''排序''',
  `create_time` datetime(3) NOT NULL COMMENT '''创建时间''',
  `update_time` datetime(3) DEFAULT NULL COMMENT '''更新时间''',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `tool_link`
--

LOCK TABLES `tool_link` WRITE;
/*!40000 ALTER TABLE `tool_link` DISABLE KEYS */;
/*!40000 ALTER TABLE `tool_link` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `tool_service_deploy`
--

DROP TABLE IF EXISTS `tool_service_deploy`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `tool_service_deploy` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '''主键''',
  `service_name` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''服务名称''',
  `service_id` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''服务ID''',
  `version` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''服务版本''',
  `host_id` bigint unsigned NOT NULL COMMENT '''主机ID''',
  `host_ip` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''主机IP''',
  `install_dir` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '''安装目录''',
  `container_name` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''容器名称''',
  `ports` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''端口映射(JSON)''',
  `env_vars` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''环境变量(JSON)''',
  `status` bigint DEFAULT '0' COMMENT '''状态:0->部署中,1->运行中,2->已停止,3->部署失败''',
  `deploy_log` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '''部署日志''',
  `create_time` datetime(3) NOT NULL COMMENT '''创建时间''',
  `update_time` datetime(3) DEFAULT NULL COMMENT '''更新时间''',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `tool_service_deploy`
--

LOCK TABLES `tool_service_deploy` WRITE;
/*!40000 ALTER TABLE `tool_service_deploy` DISABLE KEYS */;
INSERT INTO `tool_service_deploy` VALUES (7,'Redis','redis','redis-7.2',501,'8.130.14.34','/opt/data/redis','','','{\"REDIS_MAXMEMORY\":\"2gb\",\"REDIS_PASSWORD\":\"redis123456\",\"REDIS_PORT\":\"6370\"}',1,'[2025-10-30 17:20:01] 开始部署 Redis Redis 7.2\n[2025-10-30 17:20:01] 连接主机 8.130.14.34...\n[2025-10-30 17:20:02] SSH连接成功\n[2025-10-30 17:20:03] 创建安装目录 /opt/data/redis...\n[2025-10-30 17:20:03] 读取模板文件 common/templates/05-redis/versions/redis-7.2-docker-compose.yml...\n[2025-10-30 17:20:03] 生成环境变量配置...\n[2025-10-30 17:20:03] 上传 docker-compose.yml...\n[2025-10-30 17:20:04] 上传 .env...\n[2025-10-30 17:20:04] 检查Docker环境...\n[2025-10-30 17:20:04] 启动服务容器...\n[2025-10-30 17:20:05] 容器启动输出:\n\n[2025-10-30 17:20:05] 验证容器状态...\n[2025-10-30 17:20:08] 容器状态:\nNAME      IMAGE                                                                                      COMMAND                  SERVICE   CREATED         STATUS                            PORTS\nredis72   crpi-aj3vgoxp9kzh2jx4.cn-hangzhou.personal.cr.aliyuncs.com/zhangfan_k8s/redis:7.2-alpine   \"docker-entrypoint.s…\"   redis     4 seconds ago   Up 3 seconds (health: starting)   0.0.0.0:6370->6379/tcp, [::]:6370->6379/tcp\n\n[2025-10-30 17:20:08] 部署完成！\n','2025-10-30 17:20:01.477','2025-10-30 17:20:01.477'),(11,'Java','java','java-17',501,'8.130.14.34','/opt/data/java','','','{\"APP_PORT\":\"8080\",\"JAVA_OPTS\":\"-Xmx512m -Xms256m\"}',1,'[2025-10-31 12:53:21] 开始部署 Java Java 17 LTS\n[2025-10-31 12:53:21] 连接主机 8.130.14.34...\n[2025-10-31 12:53:22] SSH连接成功\n[2025-10-31 12:53:22] 检查Docker环境...\n[2025-10-31 12:53:22] 使用镜像: crpi-aj3vgoxp9kzh2jx4.cn-hangzhou.personal.cr.aliyuncs.com/zhangfan_k8s/openjdk:17-jdk\n[2025-10-31 12:53:22] 拉取镜像...\n[2025-10-31 12:53:23] 镜像拉取成功\n[2025-10-31 12:53:23] 创建临时容器...\n[2025-10-31 12:53:23] 提取文件 /usr/local/openjdk-17 -> /usr/local/java17...\n[2025-10-31 12:53:29] 清理临时容器...\n[2025-10-31 12:53:29] 读取安装脚本 common/templates/02-java/versions/java-17-install.sh...\n[2025-10-31 12:53:29] 上传安装脚本...\n[2025-10-31 12:53:29] 执行安装脚本...\n[2025-10-31 12:53:30] 安装脚本输出:\n===== Java 17 LTS 安装配置 =====\n安装路径: /usr/local/java17\n环境变量文件: /etc/profile.d/java17.sh\n环境变量已配置: /etc/profile.d/java17.sh\n===== 安装成功 =====\n\n[2025-10-31 12:53:30] 验证安装...\n[2025-10-31 12:53:30] 验证结果:\nopenjdk version \"17.0.0.1\" 2024-07-02\nOpenJDK Runtime Environment (build 17.0.0.1+2-3)\nOpenJDK 64-Bit Server VM (build 17.0.0.1+2-3, mixed mode, sharing)\n\n[2025-10-31 12:53:30] 部署完成！\n','2025-10-31 12:53:21.350','2025-10-31 12:53:21.350'),(12,'Elasticsearch','elasticsearch','elasticsearch-8.x',506,'139.9.205.38','/opt/data/elasticsearch','','','{\"ES_HEAP_SIZE\":\"1g\",\"ES_HTTP_PORT\":\"9200\"}',3,'[2025-11-30 22:33:06] 开始部署 Elasticsearch Elasticsearch 8.x\n[2025-11-30 22:33:06] 连接主机 139.9.205.38...\n[2025-11-30 22:33:07] SSH连接成功\n[2025-11-30 22:33:07] 创建安装目录 /opt/data/elasticsearch...\n[2025-11-30 22:33:07] 读取模板文件 common/templates/06-elasticsearch/versions/elasticsearch-8.x-docker-compose.yml...\n[2025-11-30 22:33:07] 生成环境变量配置...\n[2025-11-30 22:33:07] 上传 docker-compose.yml...\n[2025-11-30 22:33:07] 上传 .env...\n[2025-11-30 22:33:07] 检查Docker环境...\n[2025-11-30 22:33:08] 启动服务容器...\n[2025-11-30 22:33:08] 启动失败: Process exited with status 127\n输出: \nSTDERR:\nbash: line 1: docker-compose: command not found\n\n','2025-11-30 22:33:06.149','2025-11-30 22:33:06.149');
/*!40000 ALTER TABLE `tool_service_deploy` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Dumping events for database 'autoops'
--

--
-- Dumping routines for database 'autoops'
--

--
-- Final view structure for view `db_instance_all`
--

/*!50001 DROP VIEW IF EXISTS `db_instance_all`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_general_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`root`@`%` SQL SECURITY DEFINER */
/*!50001 VIEW `db_instance_all` AS select `db_instance`.`id` AS `id`,`db_instance`.`code` AS `code`,`db_instance`.`name` AS `name`,(case when (`db_instance`.`type` = 'postgres') then 'postgres' else 'mysql' end) AS `db_type`,`db_instance`.`type` AS `sub_type`,`db_instance`.`host` AS `host`,`db_instance`.`port` AS `port`,`db_instance`.`username` AS `username`,`db_instance`.`password` AS `password`,`db_instance`.`remark` AS `remark`,`db_instance`.`ssh_tunnel_machine_id` AS `ssh_tunnel_machine_id`,`db_instance`.`status` AS `status`,`db_instance`.`create_time` AS `create_time`,`db_instance`.`update_time` AS `update_time`,`db_instance`.`creator` AS `creator`,`db_instance`.`creator_id` AS `creator_id`,`db_instance`.`modifier` AS `modifier`,`db_instance`.`modifier_id` AS `modifier_id`,json_object('network',`db_instance`.`network`,'params',`db_instance`.`params`) AS `connection_config` from `db_instance` union all select `db_redis_instance`.`id` AS `id`,`db_redis_instance`.`code` AS `code`,`db_redis_instance`.`name` AS `name`,'redis' AS `db_type`,`db_redis_instance`.`mode` AS `sub_type`,`db_redis_instance`.`host` AS `host`,`db_redis_instance`.`port` AS `port`,`db_redis_instance`.`username` AS `username`,`db_redis_instance`.`password` AS `password`,`db_redis_instance`.`remark` AS `remark`,`db_redis_instance`.`ssh_tunnel_machine_id` AS `ssh_tunnel_machine_id`,`db_redis_instance`.`status` AS `status`,`db_redis_instance`.`create_time` AS `create_time`,`db_redis_instance`.`update_time` AS `update_time`,`db_redis_instance`.`creator` AS `creator`,`db_redis_instance`.`creator_id` AS `creator_id`,`db_redis_instance`.`modifier` AS `modifier`,`db_redis_instance`.`modifier_id` AS `modifier_id`,json_object('mode',`db_redis_instance`.`mode`,'db',`db_redis_instance`.`db`) AS `connection_config` from `db_redis_instance` union all select `db_mongo_instance`.`id` AS `id`,`db_mongo_instance`.`code` AS `code`,`db_mongo_instance`.`name` AS `name`,'mongodb' AS `db_type`,`db_mongo_instance`.`type` AS `sub_type`,NULL AS `host`,NULL AS `port`,NULL AS `username`,NULL AS `password`,`db_mongo_instance`.`remark` AS `remark`,`db_mongo_instance`.`ssh_tunnel_machine_id` AS `ssh_tunnel_machine_id`,`db_mongo_instance`.`status` AS `status`,`db_mongo_instance`.`create_time` AS `create_time`,`db_mongo_instance`.`update_time` AS `update_time`,`db_mongo_instance`.`creator` AS `creator`,`db_mongo_instance`.`creator_id` AS `creator_id`,`db_mongo_instance`.`modifier` AS `modifier`,`db_mongo_instance`.`modifier_id` AS `modifier_id`,json_object('uri',`db_mongo_instance`.`uri`) AS `connection_config` from `db_mongo_instance` union all select `db_es_instance`.`id` AS `id`,`db_es_instance`.`code` AS `code`,`db_es_instance`.`name` AS `name`,'elasticsearch' AS `db_type`,`db_es_instance`.`protocol` AS `sub_type`,`db_es_instance`.`host` AS `host`,`db_es_instance`.`port` AS `port`,`db_es_instance`.`username` AS `username`,`db_es_instance`.`password` AS `password`,`db_es_instance`.`remark` AS `remark`,`db_es_instance`.`ssh_tunnel_machine_id` AS `ssh_tunnel_machine_id`,`db_es_instance`.`status` AS `status`,`db_es_instance`.`create_time` AS `create_time`,`db_es_instance`.`update_time` AS `update_time`,`db_es_instance`.`creator` AS `creator`,`db_es_instance`.`creator_id` AS `creator_id`,`db_es_instance`.`modifier` AS `modifier`,`db_es_instance`.`modifier_id` AS `modifier_id`,json_object('protocol',`db_es_instance`.`protocol`) AS `connection_config` from `db_es_instance` */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2026-02-19 23:34:37
