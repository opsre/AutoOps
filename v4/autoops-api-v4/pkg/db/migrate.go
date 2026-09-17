// pkg/db/migrate.go
package db

import (
	"crypto/sha256"
	agentmodel "dodevops-api/api/agent/model"
	aimodel "dodevops-api/api/ai/model"
	appmodel "dodevops-api/api/app/model"
	service_release "dodevops-api/api/app/model"
	cmdbmodel "dodevops-api/api/cmdb/model"
	ccmodel "dodevops-api/api/configcenter/model"
	esmodel "dodevops-api/api/db/elasticsearch/model"
	mongomodel "dodevops-api/api/db/mongo/model"
	dbmodel "dodevops-api/api/db/mysql/model"
	rmodel "dodevops-api/api/db/redis/model"
	k8smodel "dodevops-api/api/k8s/model"
	monitormodel "dodevops-api/api/monitor/model"
	systemmodel "dodevops-api/api/system/model"
	taskmodel "dodevops-api/api/task/model"
	toolmodel "dodevops-api/api/tool/model"
	"encoding/hex"
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 注册所有需要自动建表的 model
var models = []interface{}{
	&aimodel.AiModel{},             // AI模型管理
	&agentmodel.AgentTask{},        // Agent任务
	&agentmodel.AgentChatHistory{}, // Agent对话历史
	&agentmodel.KnowledgeBase{},    // 知识库
	&agentmodel.PromptTemplate{},   // Prompt模板
	&cmdbmodel.CmdbGroup{},
	&ccmodel.EcsAuth{},
	&ccmodel.KeyManage{},
	&ccmodel.SyncSchedule{},
	&cmdbmodel.CmdbHost{},
	&cmdbmodel.CmdbSQLRecord{},
	&cmdbmodel.CmdbSQL{},
	&ccmodel.AccountAuth{},
	&taskmodel.TaskTemplate{},
	&taskmodel.Task{},
	&taskmodel.TaskWork{},
	&taskmodel.TaskAnsible{},
	&taskmodel.TaskAnsibleWork{},
	&monitormodel.Agent{},
	&monitormodel.Domain{},                // 域名监控
	&monitormodel.DomainSchedule{},        // 域名定时巡检配置
	&monitormodel.Incident{},              // 故障管理
	&monitormodel.SslCert{},               // SSL证书
	&monitormodel.AliyunConfig{},          // 阿里云配置
	&monitormodel.AlertSource{},           // 告警源配置
	&monitormodel.AlertRule{},             // 告警规则
	&monitormodel.AlertRuleTemplate{},     // 告警规则模板
	&monitormodel.AlertPolicy{},           // 告警策略
	&monitormodel.AlertPolicyRule{},       // 告警策略规则
	&monitormodel.AlertRecord{},           // 告警记录
	&monitormodel.AlertMute{},             // 告警静默
	&monitormodel.NotifyRobot{},           // 通知机器人（钉钉/飞书/企业微信）
	&monitormodel.WebhookLog{},            // Webhook接收日志
	&monitormodel.WebhookNotifyLog{},      // Webhook通知日志详情
	&monitormodel.APIBusinessFunction{},   // API业务功能分类
	&monitormodel.APIMonitor{},            // API监控配置
	&monitormodel.APIMonitorResult{},      // API监控检测结果
	&cmdbmodel.CmdbSnmpDevice{},           // SNMP设备资产(仅此一张表,其他数据存Prometheus)
	&cmdbmodel.BastionUserGroup{},         // 堡垒机用户组
	&cmdbmodel.BastionHostAuthorization{}, // 堡垒机主机授权
	&k8smodel.KubeCluster{},
	&appmodel.Application{},
	&appmodel.JenkinsEnv{},
	&appmodel.QuickDeployment{},
	&appmodel.QuickDeploymentTask{},
	&systemmodel.SysOperationLog{},
	&systemmodel.SysActivityLog{},      // 系统动态日志
	&systemmodel.SysConfig{},           // 系统配置
	&systemmodel.SysSessionRecording{}, // SSH会话录制
	&systemmodel.SysCommandAudit{},     // SSH命令审计
	&toolmodel.Tool{},
	&toolmodel.ServiceDeploy{},
	&toolmodel.FeishuUserToken{}, // 飞书 OAuth 用户令牌
	&service_release.ServiceRelease{},
	&service_release.ServiceReleaseItem{}, // 服务上线子表
	&appmodel.ShRelease{},                 // 脚本上线
	// 数据库管理模型
	&dbmodel.DbInstance{},
	&dbmodel.Db{},
	&dbmodel.DbSql{},
	&dbmodel.DbSqlExec{},
	&dbmodel.ExportTask{}, // 导出任务表
	// Redis 管理模型
	&rmodel.RedisInstance{},
	// Elasticsearch 管理模型
	&esmodel.EsInstance{},
	// MongoDB 管理模型
	&mongomodel.MongoInstance{},
	// 可以继续添加其他模型...
}

const autoMigrateSignatureKey = "auto_migrate_models_v1"
const autoMigrateBaselineTableThreshold = 10

// 自动迁移所有模型
func AutoMigrate(db *gorm.DB) error {
	signature, err := calculateModelSignature()
	if err != nil {
		return err
	}

	if err := ensureMigrationMetaTable(db); err != nil {
		return err
	}

	storedSignature, err := getStoredMigrationSignature(db, autoMigrateSignatureKey)
	if err != nil {
		return err
	}

	if storedSignature == "" {
		existingTableCount, err := countUserTables(db)
		if err != nil {
			return err
		}

		// 已有成熟业务库首次接入该优化版本时，直接写入当前模型签名作为基线，
		// 避免首次启动仍然执行全量 AutoMigrate 扫描全部 information_schema。
		if existingTableCount >= autoMigrateBaselineTableThreshold {
			log.Printf("AutoMigrate baseline initialized: detected %d existing tables, storing signature %s",
				existingTableCount, shortSignature(signature))
			if err := createUnifiedInstanceView(db); err != nil {
				return err
			}
			return setStoredMigrationSignature(db, autoMigrateSignatureKey, signature)
		}
	}

	if storedSignature == signature {
		log.Printf("AutoMigrate skipped: model signature unchanged (%s)", signature[:12])
		return createUnifiedInstanceView(db)
	}

	log.Printf("AutoMigrate running: model signature changed old=%s new=%s",
		shortSignature(storedSignature), shortSignature(signature))

	// 1. 先迁移所有表结构
	if err := db.AutoMigrate(models...); err != nil {
		return err
	}

	// 2. 创建统一数据库实例视图
	if err := createUnifiedInstanceView(db); err != nil {
		return err
	}

	if err := setStoredMigrationSignature(db, autoMigrateSignatureKey, signature); err != nil {
		return err
	}

	return nil
}

func calculateModelSignature() (string, error) {
	parts := make([]string, 0, len(models))

	for _, m := range models {
		t := reflect.TypeOf(m)
		if t == nil {
			continue
		}
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		if t.Kind() != reflect.Struct {
			return "", fmt.Errorf("unsupported model type: %T", m)
		}

		fieldParts := make([]string, 0, t.NumField())
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			fieldParts = append(fieldParts, strings.Join([]string{
				f.Name,
				f.Type.String(),
				string(f.Tag),
			}, "|"))
		}

		tableName := t.PkgPath() + "." + t.Name()
		if namer, ok := reflect.New(t).Interface().(interface{ TableName() string }); ok {
			tableName = namer.TableName()
		}

		parts = append(parts, tableName+"::"+strings.Join(fieldParts, "||"))
	}

	sort.Strings(parts)
	sum := sha256.Sum256([]byte(strings.Join(parts, "\n")))
	return hex.EncodeToString(sum[:]), nil
}

func ensureMigrationMetaTable(db *gorm.DB) error {
	sql := `
CREATE TABLE IF NOT EXISTS system_migration_meta (
  meta_key VARCHAR(128) PRIMARY KEY,
  meta_value VARCHAR(128) NOT NULL,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
)`
	return db.Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)}).Exec(sql).Error
}

func getStoredMigrationSignature(db *gorm.DB, key string) (string, error) {
	type migrationMeta struct {
		MetaValue string `gorm:"column:meta_value"`
	}

	var meta migrationMeta
	err := db.Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)}).
		Table("system_migration_meta").
		Select("meta_value").
		Where("meta_key = ?", key).
		Take(&meta).Error
	if err == gorm.ErrRecordNotFound {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return meta.MetaValue, nil
}

func setStoredMigrationSignature(db *gorm.DB, key, value string) error {
	sql := `
INSERT INTO system_migration_meta (meta_key, meta_value)
VALUES (?, ?)
ON DUPLICATE KEY UPDATE
  meta_value = VALUES(meta_value),
  updated_at = CURRENT_TIMESTAMP`

	return db.Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)}).Exec(sql, key, value).Error
}

func countUserTables(db *gorm.DB) (int64, error) {
	type tableCount struct {
		Count int64 `gorm:"column:cnt"`
	}

	var result tableCount
	err := db.Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)}).
		Raw(`
SELECT COUNT(*) AS cnt
FROM information_schema.tables
WHERE table_schema = DATABASE()
  AND table_type = 'BASE TABLE'
  AND table_name <> 'system_migration_meta'`).
		Scan(&result).Error
	if err != nil {
		return 0, err
	}
	return result.Count, nil
}

func shortSignature(signature string) string {
	if signature == "" {
		return "none"
	}
	if len(signature) <= 12 {
		return signature
	}
	return signature[:12]
}

// createUnifiedInstanceView 创建统一数据库实例视图
func createUnifiedInstanceView(db *gorm.DB) error {
	// 先删除旧视图（如果存在）
	dropViewSQL := `DROP VIEW IF EXISTS db_instance_all`
	if err := db.Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)}).Exec(dropViewSQL).Error; err != nil {
		return err
	}

	// 创建新视图
	createViewSQL := `
CREATE VIEW db_instance_all AS
-- MySQL/PostgreSQL 实例
SELECT
    id,
    code,
    name,
    CASE
        WHEN type = 'postgres' THEN 'postgres'
        ELSE 'mysql'
    END AS db_type,
    type AS sub_type,
    host,
    port,
    username,
    password,
    remark,
    ssh_tunnel_machine_id,
    status,
    create_time,
    update_time,
    creator,
    creator_id,
    modifier,
    modifier_id,
    JSON_OBJECT(
        'network', network,
        'params', params
    ) AS connection_config
FROM db_instance

UNION ALL

-- Redis 实例
SELECT
    id,
    code,
    name,
    'redis' AS db_type,
    mode AS sub_type,
    host,
    port,
    username,
    password,
    remark,
    ssh_tunnel_machine_id,
    status,
    create_time,
    update_time,
    creator,
    creator_id,
    modifier,
    modifier_id,
    JSON_OBJECT(
        'mode', mode,
        'db', db
    ) AS connection_config
FROM db_redis_instance

UNION ALL

-- MongoDB 实例
SELECT
    id,
    code,
    name,
    'mongodb' AS db_type,
    type AS sub_type,
    NULL AS host,
    NULL AS port,
    NULL AS username,
    NULL AS password,
    remark,
    ssh_tunnel_machine_id,
    status,
    create_time,
    update_time,
    creator,
    creator_id,
    modifier,
    modifier_id,
    JSON_OBJECT(
        'uri', uri
    ) AS connection_config
FROM db_mongo_instance

UNION ALL

-- Elasticsearch 实例
SELECT
    id,
    code,
    name,
    'elasticsearch' AS db_type,
    protocol AS sub_type,
    host,
    port,
    username,
    password,
    remark,
    ssh_tunnel_machine_id,
    status,
    create_time,
    update_time,
    creator,
    creator_id,
    modifier,
    modifier_id,
    JSON_OBJECT(
        'protocol', protocol
    ) AS connection_config
FROM db_es_instance
`
	return db.Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Silent)}).Exec(createViewSQL).Error
}
