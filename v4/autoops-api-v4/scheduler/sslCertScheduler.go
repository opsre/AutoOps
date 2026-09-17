package scheduler

import (
	"fmt"
	"log"
	"time"

	"dodevops-api/api/monitor/model"
	"dodevops-api/common"
	"dodevops-api/common/util/aliyun"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

// SslCertScheduler SSL证书状态查询调度器
type SslCertScheduler struct {
	cron *cron.Cron
}

// NewSslCertScheduler 创建SSL证书调度器
func NewSslCertScheduler() *SslCertScheduler {
	return &SslCertScheduler{
		cron: cron.New(),
	}
}

// Start 启动调度器
func (s *SslCertScheduler) Start() error {
	log.Println("启动SSL证书状态查询调度器...")

	// 每2分钟查询一次申请中和验证中的证书状态
	_, err := s.cron.AddFunc("*/2 * * * *", func() {
		s.checkPendingCerts()
	})
	if err != nil {
		return fmt.Errorf("添加SSL证书查询任务失败: %v", err)
	}

	s.cron.Start()
	log.Println("SSL证书状态查询调度器启动成功")
	return nil
}

// Stop 停止调度器
func (s *SslCertScheduler) Stop() {
	log.Println("停止SSL证书状态查询调度器...")
	s.cron.Stop()
	log.Println("SSL证书状态查询调度器已停止")
}

// checkPendingCerts 检查待处理的证书
func (s *SslCertScheduler) checkPendingCerts() {
	db := common.GetDB()

	// 查询申请中(1)和验证中(2)的阿里云CAS证书
	var certs []model.SslCert
	err := db.Where("cert_source = ? AND status IN (?, ?)", model.CertSourceAliyunCAS,
		model.SslCertStatusApplying, model.SslCertStatusValidating).
		Find(&certs).Error

	if err != nil {
		log.Printf("[ERROR] 查询待处理证书失败: %v\n", err)
		return
	}

	if len(certs) == 0 {
		return
	}

	log.Printf("[INFO] 发现 %d 个待处理证书，开始查询状态...\n", len(certs))

	for _, cert := range certs {
		s.updateCertState(db, &cert)
	}
}

// updateCertState 更新证书状态
func (s *SslCertScheduler) updateCertState(db *gorm.DB, cert *model.SslCert) {
	// 1. 查询阿里云配置
	var aliyunConfig model.AliyunConfig
	if err := db.First(&aliyunConfig, cert.AliyunConfigID).Error; err != nil {
		log.Printf("[ERROR] 查询阿里云配置失败 - CertID: %d, Error: %v\n", cert.ID, err)
		return
	}

	// 2. 创建CAS客户端
	casClient, err := aliyun.NewCasClient(aliyun.CasConfig{
		AccessKey:    aliyunConfig.AccessKey,
		AccessSecret: aliyunConfig.AccessSecret,
	})
	if err != nil {
		log.Printf("[ERROR] 创建CAS客户端失败 - CertID: %d, Error: %v\n", cert.ID, err)
		return
	}

	// 3. 解析订单ID
	var orderID int64
	fmt.Sscanf(cert.OrderID, "%d", &orderID)
	if orderID == 0 {
		log.Printf("[WARN] 证书订单ID为空 - CertID: %d\n", cert.ID)
		return
	}

	// 4. 查询证书状态
	stateResp, err := casClient.GetCertState(orderID)
	if err != nil {
		log.Printf("[ERROR] 查询证书状态失败 - CertID: %d, OrderID: %d, Error: %v\n", cert.ID, orderID, err)
		// 更新错误信息
		db.Model(cert).Updates(map[string]interface{}{
			"error_msg":   fmt.Sprintf("查询证书状态失败: %v", err),
			"update_time": time.Now(),
		})
		return
	}

	log.Printf("[INFO] 证书状态查询 - CertID: %d, Domain: %s, Status: %s\n", cert.ID, cert.Domain, stateResp.Status)

	// 5. 更新证书状态
	updates := map[string]interface{}{
		"validate_info": stateResp.ValidateInfo,
		"update_time":   time.Now(),
	}

	// 根据阿里云返回的状态更新本地状态
	switch stateResp.Status {
	case "checking", "domain_verify":
		// 验证中 - 自动添加DNS验证记录
		if cert.Status == model.SslCertStatusApplying {
			// 第一次进入验证阶段，尝试自动添加DNS验证记录
			log.Printf("[INFO] 证书进入验证阶段，尝试自动添加DNS验证 - CertID: %d\n", cert.ID)
			s.autoAddDnsValidation(db, casClient, &aliyunConfig, cert, stateResp)
			return // autoAddDnsValidation内部会更新状态
		} else if cert.Status == model.SslCertStatusValidating {
			// 已在验证中，继续等待
			log.Printf("[INFO] 证书验证中，等待阿里云完成验证 - CertID: %d\n", cert.ID)
		}
	case "certificate":
		// 已签发，尝试下载证书
		if cert.Status != model.SslCertStatusIssued {
			log.Printf("[INFO] 证书已签发，开始下载 - CertID: %d, OrderID: %d\n", cert.ID, orderID)
			s.downloadAndSaveCert(db, casClient, cert, orderID)
			return // downloadAndSaveCert内部会更新状态
		}
	case "revoked":
		updates["status"] = model.SslCertStatusExpired
		log.Printf("[WARN] 证书已被撤销 - CertID: %d\n", cert.ID)
	case "expired":
		updates["status"] = model.SslCertStatusExpired
		log.Printf("[WARN] 证书已过期 - CertID: %d\n", cert.ID)
	}

	// 更新数据库
	if err := db.Model(cert).Updates(updates).Error; err != nil {
		log.Printf("[ERROR] 更新证书状态失败 - CertID: %d, Error: %v\n", cert.ID, err)
	}
}

// downloadAndSaveCert 下载并保存证书
func (s *SslCertScheduler) downloadAndSaveCert(db *gorm.DB, casClient *aliyun.CasClient, cert *model.SslCert, orderID int64) {
	// 1. 下载证书
	downloadResp, err := casClient.DownloadCertByOrderID(orderID)
	if err != nil {
		log.Printf("[ERROR] 下载证书失败 - CertID: %d, OrderID: %d, Error: %v\n", cert.ID, orderID, err)
		db.Model(cert).Updates(map[string]interface{}{
			"error_msg":   fmt.Sprintf("下载证书失败: %v", err),
			"update_time": time.Now(),
		})
		return
	}

	// 2. 解析证书有效期（使用 DownloadCertResponse 中已有的字段）
	var issueTime, expireTime time.Time
	var daysLeft int

	// 尝试解析日期字符串
	if downloadResp.IssueDate != "" {
		issueTime, _ = time.Parse("2006-01-02", downloadResp.IssueDate)
	}
	if downloadResp.ExpireDate != "" {
		expireTime, _ = time.Parse("2006-01-02", downloadResp.ExpireDate)
	}
	daysLeft = downloadResp.DaysLeft

	// 如果解析失败，使用默认值
	if issueTime.IsZero() {
		issueTime = time.Now()
	}
	if expireTime.IsZero() {
		expireTime = time.Now().AddDate(0, 3, 0) // 默认3个月
	}

	// 3. 更新证书记录
	updates := map[string]interface{}{
		"status":       model.SslCertStatusIssued,
		"cert":         downloadResp.Cert,
		"private_key":  downloadResp.PrivateKey,
		"issuer_cert":  downloadResp.CertChain,
		"issue_time":   issueTime,
		"expire_time":  expireTime,
		"days_left":    daysLeft,
		"error_msg":    "",
		"update_time":  time.Now(),
	}

	if err := db.Model(cert).Updates(updates).Error; err != nil {
		log.Printf("[ERROR] 保存证书失败 - CertID: %d, Error: %v\n", cert.ID, err)
		return
	}

	log.Printf("[SUCCESS] ✅ 证书下载并保存成功 - CertID: %d, Domain: %s, 有效期: %s ~ %s\n",
		cert.ID, cert.Domain, downloadResp.IssueDate, downloadResp.ExpireDate)
}

// autoAddDnsValidation 自动添加DNS验证记录
func (s *SslCertScheduler) autoAddDnsValidation(db *gorm.DB, casClient *aliyun.CasClient,
	aliyunConfig *model.AliyunConfig, cert *model.SslCert, stateResp *aliyun.CertStateResponse) {

	log.Printf("[INFO] 开始自动添加DNS验证记录 - CertID: %d, Domain: %s\n", cert.ID, cert.Domain)

	// 1. 检查验证信息是否存在
	if stateResp.ValidateInfo == "" {
		log.Printf("[WARN] 验证信息为空，等待阿里云生成 - CertID: %d\n", cert.ID)
		db.Model(cert).Updates(map[string]interface{}{
			"error_msg":   "等待阿里云生成DNS验证信息",
			"update_time": time.Now(),
		})
		return
	}

	// 2. 创建DNS客户端
	dnsClient, err := aliyun.NewDnsClient(aliyun.DnsConfig{
		AccessKey:    aliyunConfig.AccessKey,
		AccessSecret: aliyunConfig.AccessSecret,
	})
	if err != nil {
		log.Printf("[ERROR] 创建DNS客户端失败 - CertID: %d, Error: %v\n", cert.ID, err)
		db.Model(cert).Updates(map[string]interface{}{
			"error_msg":   fmt.Sprintf("创建DNS客户端失败: %v", err),
			"update_time": time.Now(),
		})
		return
	}

	// 3. 添加DNS验证记录
	addResp, err := dnsClient.AddCertValidationRecord(cert.Domain, stateResp.ValidateInfo)
	if err != nil {
		log.Printf("[ERROR] 添加DNS验证记录失败 - CertID: %d, Domain: %s, Error: %v\n", cert.ID, cert.Domain, err)
		db.Model(cert).Updates(map[string]interface{}{
			"error_msg":   fmt.Sprintf("添加DNS验证记录失败: %v", err),
			"update_time": time.Now(),
		})
		return
	}

	log.Printf("[SUCCESS] ✅ DNS验证记录添加成功 - CertID: %d, RecordID: %s\n", cert.ID, addResp.RecordID)

	// 4. 更新证书状态为验证中
	updates := map[string]interface{}{
		"status":        model.SslCertStatusValidating,
		"validate_info": stateResp.ValidateInfo,
		"error_msg":     "", // 清空错误信息
		"update_time":   time.Now(),
	}

	if err := db.Model(cert).Updates(updates).Error; err != nil {
		log.Printf("[ERROR] 更新证书状态失败 - CertID: %d, Error: %v\n", cert.ID, err)
		return
	}

	log.Printf("[INFO] 证书状态已更新为验证中 - CertID: %d\n", cert.ID)
}

// GetStats 获取调度器统计信息
func (s *SslCertScheduler) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"total_jobs":   len(s.cron.Entries()),
		"cron_entries": s.cron.Entries(),
	}
}
