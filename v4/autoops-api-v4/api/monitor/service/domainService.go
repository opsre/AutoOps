package service

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	configCenterDao "dodevops-api/api/configcenter/dao"
	"dodevops-api/api/monitor/dao"
	"dodevops-api/api/monitor/model"
	"dodevops-api/common/result"
	"dodevops-api/common/util/aliyun"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DomainService 域名监控服务接口
type DomainService interface {
	Create(c *gin.Context, dto *model.DomainDto)
	GetByID(c *gin.Context, id uint)
	Update(c *gin.Context, dto *model.DomainDto)
	Delete(c *gin.Context, id uint)
	GetList(c *gin.Context, query *model.DomainQuery)
	CheckDomain(c *gin.Context, id uint)
	CheckAllDomains(c *gin.Context)
	GetSSLExpiringSoon(c *gin.Context, days int)
	GetStatistics(c *gin.Context)
	BatchUploadDomains(c *gin.Context)
	ImportFromAliyunDNS(c *gin.Context, dto *model.ImportFromAliyunDNSDto)
}

// DomainServiceImpl 域名监控服务实现
type DomainServiceImpl struct {
	domainDao    dao.DomainDao
	keyManageDao *configCenterDao.KeyManageDao
}

// NewDomainService 创建域名监控服务实例
func NewDomainService() DomainService {
	return &DomainServiceImpl{
		domainDao:    dao.NewDomainDao(),
		keyManageDao: configCenterDao.NewKeyManageDao(),
	}
}

// Create 创建域名监控
func (s *DomainServiceImpl) Create(c *gin.Context, dto *model.DomainDto) {
	// 检查域名是否已存在
	existing, _ := s.domainDao.GetByDomain(dto.Domain)
	if existing != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "域名已存在")
		return
	}

	domain := &model.Domain{
		Domain: dto.Domain,
		Tags:   dto.Tags,
		Remark: dto.Remark,
		Status: 1, // 默认启用
	}

	if err := s.domainDao.Create(domain); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("创建失败: %v", err))
		return
	}

	// 创建成功后立即执行一次检查
	checkResult := s.performCheck(domain.Domain)
	checkResult.DomainID = domain.ID
	s.domainDao.UpdateCheckResult(domain.ID, checkResult)

	// 重新获取更新后的域名信息
	domain, _ = s.domainDao.GetByID(domain.ID)

	result.Success(c, domain)
}

// GetByID 根据ID获取域名监控
func (s *DomainServiceImpl) GetByID(c *gin.Context, id uint) {
	domain, err := s.domainDao.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			result.Failed(c, int(result.ApiCode.FAILED), "域名监控不存在")
		} else {
			result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("查询失败: %v", err))
		}
		return
	}

	result.Success(c, domain)
}

// Update 更新域名监控
func (s *DomainServiceImpl) Update(c *gin.Context, dto *model.DomainDto) {
	// 检查域名是否存在
	domain, err := s.domainDao.GetByID(dto.ID)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "域名监控不存在")
		return
	}

	// 如果修改了域名，检查新域名是否已存在
	if dto.Domain != domain.Domain {
		existing, _ := s.domainDao.GetByDomain(dto.Domain)
		if existing != nil && existing.ID != dto.ID {
			result.Failed(c, int(result.ApiCode.FAILED), "域名已存在")
			return
		}
	}

	// 更新字段
	domain.Domain = dto.Domain
	domain.Tags = dto.Tags
	domain.Remark = dto.Remark
	domain.Status = dto.Status

	if err := s.domainDao.Update(domain); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("更新失败: %v", err))
		return
	}

	result.Success(c, domain)
}

// Delete 删除域名监控
func (s *DomainServiceImpl) Delete(c *gin.Context, id uint) {
	// 检查是否存在
	_, err := s.domainDao.GetByID(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "域名监控不存在")
		return
	}

	if err := s.domainDao.Delete(id); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("删除失败: %v", err))
		return
	}

	result.Success(c, nil)
}

// GetList 获取域名监控列表
func (s *DomainServiceImpl) GetList(c *gin.Context, query *model.DomainQuery) {
	domains, total, err := s.domainDao.GetList(query)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("查询失败: %v", err))
		return
	}

	data := map[string]interface{}{
		"list":     domains,
		"total":    total,
		"pageNum":  query.PageNum,
		"pageSize": query.PageSize,
	}

	result.Success(c, data)
}

// CheckDomain 检查单个域名
func (s *DomainServiceImpl) CheckDomain(c *gin.Context, id uint) {
	// 获取域名信息
	domain, err := s.domainDao.GetByID(id)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "域名监控不存在")
		return
	}

	// 执行检查
	checkResult := s.performCheck(domain.Domain)
	checkResult.DomainID = id

	// 更新检查结果
	if err := s.domainDao.UpdateCheckResult(id, checkResult); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("更新检查结果失败: %v", err))
		return
	}

	result.Success(c, checkResult)
}

// CheckAllDomains 检查所有启用的域名（异步执行）
func (s *DomainServiceImpl) CheckAllDomains(c *gin.Context) {
	domains, err := s.domainDao.GetEnabledDomains()
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("获取域名列表失败: %v", err))
		return
	}

	total := len(domains)
	if total == 0 {
		result.Success(c, map[string]interface{}{
			"total":   0,
			"message": "没有需要检查的域名",
		})
		return
	}

	// 立即返回，告知任务已启动
	result.Success(c, map[string]interface{}{
		"total":   total,
		"message": fmt.Sprintf("批量检查任务已启动，共 %d 个域名，正在后台执行", total),
		"status":  "running",
	})

	// 在后台异步执行检查任务
	go s.checkDomainsAsync(domains)
}

// checkDomainsAsync 异步检查域名（后台任务）
func (s *DomainServiceImpl) checkDomainsAsync(domains []model.Domain) {
	// 使用带缓冲的通道控制并发数，避免同时发起太多请求
	maxConcurrent := 10 // 最多同时检查10个域名
	semaphore := make(chan struct{}, maxConcurrent)

	// 使用 WaitGroup 等待所有检查完成
	var wg sync.WaitGroup

	for _, domain := range domains {
		wg.Add(1)

		// 启动 goroutine 检查单个域名
		go func(d model.Domain) {
			defer wg.Done()

			// 获取信号量
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// 执行检查
			checkResult := s.performCheck(d.Domain)
			checkResult.DomainID = d.ID

			// 更新检查结果
			s.domainDao.UpdateCheckResult(d.ID, checkResult)
		}(domain)
	}

	// 等待所有检查完成
	wg.Wait()
}

// GetSSLExpiringSoon 获取SSL证书即将过期的域名
func (s *DomainServiceImpl) GetSSLExpiringSoon(c *gin.Context, days int) {
	domains, err := s.domainDao.GetSSLExpiringSoon(days)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("查询失败: %v", err))
		return
	}

	result.Success(c, map[string]interface{}{
		"total":   len(domains),
		"days":    days,
		"domains": domains,
	})
}

// GetStatistics 获取统计信息
func (s *DomainServiceImpl) GetStatistics(c *gin.Context) {
	stats, err := s.domainDao.GetStatistics()
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("获取统计信息失败: %v", err))
		return
	}

	result.Success(c, stats)
}

// performCheck 执行域名检查（核心逻辑）
func (s *DomainServiceImpl) performCheck(domainName string) *model.DomainCheckResult {
	checkResult := &model.DomainCheckResult{}

	// 构建URL（优先尝试HTTPS）
	url := fmt.Sprintf("https://%s", domainName)

	// 创建HTTP客户端（允许不安全的SSL证书）
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // 允许不安全的SSL证书，以便获取证书信息
			},
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// 允许重定向，但最多10次
			if len(via) >= 10 {
				return fmt.Errorf("stopped after 10 redirects")
			}
			return nil
		},
	}

	// 记录开始时间
	startTime := time.Now()

	// 发起HTTPS请求
	resp, err := client.Get(url)
	httpsSuccess := false
	hasSSL := false

	if err != nil {
		// HTTPS失败，尝试HTTP
		url = fmt.Sprintf("http://%s", domainName)
		startTime = time.Now() // 重新计时
		resp, err = client.Get(url)
		if err != nil {
			// HTTP和HTTPS都失败
			checkResult.IsAlive = 0
			checkResult.ErrorMsg = fmt.Sprintf("连接失败: %v", err)
			return checkResult
		}
	} else {
		httpsSuccess = true
	}
	defer resp.Body.Close()

	// 计算响应时间
	checkResult.ResponseTime = int(time.Since(startTime).Milliseconds())
	checkResult.StatusCode = resp.StatusCode // 记录HTTP状态码

	// 检查SSL证书（仅HTTPS）
	if resp.TLS != nil && len(resp.TLS.PeerCertificates) > 0 {
		cert := resp.TLS.PeerCertificates[0]
		hasSSL = true

		// 证书过期时间
		checkResult.SSLExpireAt = cert.NotAfter.Format("2006-01-02 15:04:05")

		// 计算剩余天数
		daysLeft := int(time.Until(cert.NotAfter).Hours() / 24)
		checkResult.SSLDaysLeft = daysLeft

		// 证书颁发者
		checkResult.SSLIssuer = cert.Issuer.CommonName
	}

	// 判断域名状态:
	// 1. HTTPS成功且有SSL证书 -> 正常(IsAlive=1)
	// 2. 只能HTTP访问，没有SSL证书 -> 异常(IsAlive=0)，提示缺少SSL证书
	// 3. HTTPS成功但没有证书信息(异常情况) -> 异常
	if httpsSuccess && hasSSL {
		// 正常情况:HTTPS访问成功且有SSL证书
		checkResult.IsAlive = 1
	} else if !hasSSL {
		// 异常情况:没有SSL证书
		checkResult.IsAlive = 0
		if httpsSuccess {
			checkResult.ErrorMsg = "HTTPS访问成功但未检测到SSL证书"
		} else {
			checkResult.ErrorMsg = "域名缺少SSL证书，仅支持HTTP访问"
		}
	} else {
		// 其他情况视为正常
		checkResult.IsAlive = 1
	}

	return checkResult
}

// BatchUploadDomains 批量上传域名文件
func (s *DomainServiceImpl) BatchUploadDomains(c *gin.Context) {
	// 接收文件
	file, err := c.FormFile("file")
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "文件上传失败: "+err.Error())
		return
	}

	// 获取可选的统一标签和备注
	tags := c.PostForm("tags")
	remark := c.PostForm("remark")

	// 打开文件
	fileContent, err := file.Open()
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "文件打开失败: "+err.Error())
		return
	}
	defer fileContent.Close()

	// 解析文件内容，提取域名并去重
	domainSet := make(map[string]bool)
	scanner := bufio.NewScanner(fileContent)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// 跳过空行和注释行（以#开头）
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 简单处理：去除可能的协议前缀
		line = strings.TrimPrefix(line, "http://")
		line = strings.TrimPrefix(line, "https://")
		line = strings.TrimPrefix(line, "www.")

		// 去除可能的路径部分，只保留域名
		if idx := strings.Index(line, "/"); idx != -1 {
			line = line[:idx]
		}

		// 添加到去重集合
		if line != "" {
			domainSet[line] = true
		}
	}

	if err := scanner.Err(); err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "文件读取失败: "+err.Error())
		return
	}

	// 转换为域名列表
	var domains []string
	for domain := range domainSet {
		domains = append(domains, domain)
	}

	totalCount := len(domains)
	if totalCount == 0 {
		result.Failed(c, int(result.ApiCode.FAILED), "文件中没有有效的域名")
		return
	}

	// 立即返回响应，告知任务已启动
	result.Success(c, map[string]interface{}{
		"total":   totalCount,
		"message": fmt.Sprintf("批量上传任务已启动，共解析到 %d 个域名（已去重），正在后台异步创建", totalCount),
		"status":  "running",
	})

	// 在后台异步执行批量创建任务
	go s.batchCreateDomainsAsync(domains, tags, remark)
}

// batchCreateDomainsAsync 异步批量创建域名监控（后台任务）
func (s *DomainServiceImpl) batchCreateDomainsAsync(domains []string, tags, remark string) {
	// 使用带缓冲的通道控制并发数
	maxConcurrent := 5 // 最多同时创建5个域名
	semaphore := make(chan struct{}, maxConcurrent)

	// 使用 WaitGroup 等待所有任务完成
	var wg sync.WaitGroup

	// 统计结果
	var mu sync.Mutex
	successCount := 0
	skipCount := 0
	failCount := 0

	for _, domainName := range domains {
		wg.Add(1)

		// 启动 goroutine 创建单个域名
		go func(name string) {
			defer wg.Done()

			// 获取信号量
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// 检查域名是否已存在
			existing, _ := s.domainDao.GetByDomain(name)
			if existing != nil {
				mu.Lock()
				skipCount++
				mu.Unlock()
				return
			}

			// 创建域名记录
			domain := &model.Domain{
				Domain: name,
				Tags:   tags,
				Remark: remark,
				Status: 1, // 默认启用
			}

			if err := s.domainDao.Create(domain); err != nil {
				mu.Lock()
				failCount++
				mu.Unlock()
				return
			}

			// 创建成功后执行一次检查
			checkResult := s.performCheck(domain.Domain)
			checkResult.DomainID = domain.ID
			s.domainDao.UpdateCheckResult(domain.ID, checkResult)

			mu.Lock()
			successCount++
			mu.Unlock()
		}(domainName)
	}

	// 等待所有任务完成
	wg.Wait()

	// 这里可以添加日志记录或通知
	// 例如：记录批量创建的结果
	fmt.Printf("批量创建域名任务完成: 成功 %d, 跳过(已存在) %d, 失败 %d\n", successCount, skipCount, failCount)
}

// ImportFromAliyunDNS 从阿里云解析导入域名
func (s *DomainServiceImpl) ImportFromAliyunDNS(c *gin.Context, dto *model.ImportFromAliyunDNSDto) {
	// 1. 获取密钥信息
	keyManage, err := s.keyManageDao.GetByID(dto.KeyID)
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), "密钥不存在或已删除")
		return
	}

	// 验证是否为阿里云密钥
	if keyManage.KeyType != 1 {
		result.Failed(c, int(result.ApiCode.FAILED), "请选择阿里云类型的密钥")
		return
	}

	// 2. 解密密钥
	keyID, keySecret, err := keyManage.DecryptKeys()
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("密钥解密失败: %v", err))
		return
	}

	// 3. 创建阿里云DNS客户端
	dnsClient, err := aliyun.NewDnsClient(aliyun.DnsConfig{
		AccessKey:    keyID,
		AccessSecret: keySecret,
	})
	if err != nil {
		result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("创建阿里云DNS客户端失败: %v", err))
		return
	}

	// 4. 根据选项获取域名列表
	var domainNames []string
	var err2 error

	if dto.IncludeSubdomain {
		// 获取所有域名(包含子域名)
		domainNames, err2 = dnsClient.ListAllDomainsWithSubdomains()
		if err2 != nil {
			result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("获取阿里云域名和子域名列表失败: %v", err2))
			return
		}
	} else {
		// 只获取主域名
		domainInfos, err2 := dnsClient.ListAllDomains()
		if err2 != nil {
			result.Failed(c, int(result.ApiCode.FAILED), fmt.Sprintf("获取阿里云域名列表失败: %v", err2))
			return
		}

		// 转换为域名字符串列表
		domainNames = make([]string, len(domainInfos))
		for i, info := range domainInfos {
			domainNames[i] = info.DomainName
		}
	}

	if len(domainNames) == 0 {
		result.Success(c, model.AliyunDNSImportResult{
			Total:    0,
			Imported: 0,
			Skipped:  0,
			Failed:   0,
		})
		return
	}

	// 5. 立即返回响应，告知任务已启动
	importType := "主域名"
	if dto.IncludeSubdomain {
		importType = "域名(含子域名)"
	}
	result.Success(c, map[string]interface{}{
		"total":   len(domainNames),
		"message": fmt.Sprintf("从阿里云解析导入任务已启动，共 %d 个%s，正在后台异步导入", len(domainNames), importType),
		"status":  "running",
	})

	// 6. 在后台异步执行导入任务
	go s.importDomainNamesFromAliyunAsync(domainNames, dto.Tags, dto.Remark)
}

// importDomainNamesFromAliyunAsync 异步导入阿里云域名（后台任务）
// 优化版本：使用批量插入，性能提升约10-20倍
func (s *DomainServiceImpl) importDomainNamesFromAliyunAsync(domainNames []string, tags, remark string) {
	startTime := time.Now()

	// 统计结果
	imported := 0
	skipped := 0
	failed := 0
	var failedList []string

	// 第一步：批量查询已存在的域名
	fmt.Printf("[INFO] 开始批量查询已存在的域名，总数: %d\n", len(domainNames))

	// 分批查询已存在的域名（避免IN子句过长）
	existingMap := make(map[string]*model.Domain)
	batchSize := 500
	for i := 0; i < len(domainNames); i += batchSize {
		end := i + batchSize
		if end > len(domainNames) {
			end = len(domainNames)
		}

		batch, err := s.domainDao.GetByDomains(domainNames[i:end])
		if err != nil {
			fmt.Printf("[ERROR] 批量查询域名失败: %v\n", err)
			continue
		}

		for k, v := range batch {
			existingMap[k] = v
		}
	}

	skipped = len(existingMap)
	fmt.Printf("[INFO] 查询完成，已存在: %d, 需要导入: %d\n", skipped, len(domainNames)-skipped)

	// 第二步：过滤出需要导入的域名
	toImport := make([]*model.Domain, 0, len(domainNames)-skipped)
	for _, domainName := range domainNames {
		if _, exists := existingMap[domainName]; exists {
			continue
		}

		toImport = append(toImport, &model.Domain{
			Domain: domainName,
			Tags:   tags,
			Remark: remark,
			Status: 1, // 默认启用
		})
	}

	// 第三步：批量插入域名
	if len(toImport) > 0 {
		fmt.Printf("[INFO] 开始批量插入域名，数量: %d\n", len(toImport))
		if err := s.domainDao.BatchCreate(toImport); err != nil {
			fmt.Printf("[ERROR] 批量插入域名失败: %v\n", err)
			failed = len(toImport)
			for _, d := range toImport {
				failedList = append(failedList, d.Domain)
			}
		} else {
			imported = len(toImport)
			fmt.Printf("[INFO] 批量插入成功: %d 个域名\n", imported)
		}
	}

	// 第四步：并发执行健康检查（使用并发控制）
	if imported > 0 {
		fmt.Printf("[INFO] 开始批量健康检查，域名数: %d\n", imported)
		s.batchCheckDomainsAsync(toImport)
	}

	// 记录导入结果
	elapsed := time.Since(startTime)
	fmt.Printf("阿里云域名导入任务完成: 总计 %d, 成功导入 %d, 跳过(已存在) %d, 失败 %d, 耗时: %v\n",
		len(domainNames), imported, skipped, failed, elapsed)
	if len(failedList) > 0 {
		fmt.Printf("失败的域名列表: %v\n", failedList)
	}
}

// batchCheckDomainsAsync 批量异步健康检查
func (s *DomainServiceImpl) batchCheckDomainsAsync(domains []*model.Domain) {
	// 使用带缓冲的通道控制并发数
	maxConcurrent := 20 // 健康检查可以有更高的并发数
	semaphore := make(chan struct{}, maxConcurrent)

	// 使用 WaitGroup 等待所有检查完成
	var wg sync.WaitGroup

	for _, domain := range domains {
		wg.Add(1)

		// 启动 goroutine 检查单个域名
		go func(d *model.Domain) {
			defer wg.Done()

			// 获取信号量
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// 执行检查
			checkResult := s.performCheck(d.Domain)
			checkResult.DomainID = d.ID

			// 更新检查结果（忽略错误，不影响其他域名）
			if err := s.domainDao.UpdateCheckResult(d.ID, checkResult); err != nil {
				fmt.Printf("[WARN] 更新域名 %s 检查结果失败: %v\n", d.Domain, err)
			}
		}(domain)
	}

	// 等待所有检查完成
	wg.Wait()
	fmt.Printf("[INFO] 批量健康检查完成\n")
}
