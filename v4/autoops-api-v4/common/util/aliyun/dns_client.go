package aliyun

import (
	"fmt"
	"strings"

	alidns "github.com/alibabacloud-go/alidns-20150109/v4/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

// DnsClient 阿里云DNS客户端封装
type DnsClient struct {
	client *alidns.Client
	config DnsConfig
}

// DnsConfig 阿里云DNS配置
type DnsConfig struct {
	AccessKey    string
	AccessSecret string
	Endpoint     string // 默认: alidns.cn-hangzhou.aliyuncs.com
}

// NewDnsClient 创建阿里云DNS客户端
func NewDnsClient(config DnsConfig) (*DnsClient, error) {
	if config.Endpoint == "" {
		config.Endpoint = "alidns.cn-hangzhou.aliyuncs.com"
	}

	// 清理 AccessKey 和 AccessSecret 的空格
	config.AccessKey = strings.TrimSpace(config.AccessKey)
	config.AccessSecret = strings.TrimSpace(config.AccessSecret)

	// 调试日志：打印密钥长度（不打印实际内容）
	fmt.Printf("[DEBUG] DNS客户端配置 - AccessKey长度: %d, AccessSecret长度: %d, Endpoint: %s\n",
		len(config.AccessKey), len(config.AccessSecret), config.Endpoint)

	if config.AccessKey == "" || config.AccessSecret == "" {
		return nil, fmt.Errorf("AccessKey 或 AccessSecret 不能为空")
	}

	openApiConfig := &openapi.Config{
		AccessKeyId:     tea.String(config.AccessKey),
		AccessKeySecret: tea.String(config.AccessSecret),
		Endpoint:        tea.String(config.Endpoint),
	}

	client, err := alidns.NewClient(openApiConfig)
	if err != nil {
		return nil, fmt.Errorf("创建DNS客户端失败: %w", err)
	}

	return &DnsClient{
		client: client,
		config: config,
	}, nil
}

// AddDnsTxtRecordRequest 添加DNS TXT记录请求
type AddDnsTxtRecordRequest struct {
	Domain string // 主域名(如: example.com)
	RR     string // 主机记录(如: _acme-challenge)
	Value  string // 记录值
	TTL    int64  // TTL值(可选,默认600)
}

// AddDnsTxtRecordResponse 添加DNS TXT记录响应
type AddDnsTxtRecordResponse struct {
	RecordID string // 记录ID
	Message  string // 响应消息
}

// AddDnsTxtRecord 添加DNS TXT记录(用于证书验证)
// 支持两种调用方式:
// 1. 传入Domain+RR: Domain="example.com", RR="_acme-challenge"
// 2. 传入完整FQDN作为Domain: Domain="_acme-challenge.example.com", RR="" (自动解析)
func (d *DnsClient) AddDnsTxtRecord(req AddDnsTxtRecordRequest) (*AddDnsTxtRecordResponse, error) {
	var domainName, rr string

	fmt.Printf("[DEBUG] AddDnsTxtRecord 输入 - Domain: %s, RR: %s, Value: %s\n", req.Domain, req.RR, req.Value)

	// 参数验证
	if req.Domain == "" {
		return nil, fmt.Errorf("域名不能为空")
	}
	if req.Value == "" {
		return nil, fmt.Errorf("记录值不能为空")
	}

	// 如果RR为空，说明Domain是完整的FQDN，需要解析
	if req.RR == "" {
		// 先去掉尾随点，统一处理
		cleanDomain := strings.TrimSuffix(req.Domain, ".")

		// 从FQDN中提取根域名和主机记录
		rootDomain, err := GetRootDomain(cleanDomain)
		if err != nil {
			return nil, fmt.Errorf("解析根域名失败: %w", err)
		}

		domainName = rootDomain

		// 计算RR：cleanDomain去掉根域名部分
		// 例如: _acme-challenge.deviops.cn -> _acme-challenge
		//      _acme-challenge.sub.example.com -> _acme-challenge.sub
		rr = strings.TrimSuffix(cleanDomain, "."+rootDomain)
		if rr == cleanDomain {
			// 如果没有后缀可以去除，说明cleanDomain格式不正确
			return nil, fmt.Errorf("无法从FQDN中提取主机记录: %s", req.Domain)
		}

		fmt.Printf("[DEBUG] AddDnsTxtRecord FQDN解析 - 根域名: %s, RR: %s\n", rootDomain, rr)
	} else {
		// 如果RR已提供，直接使用
		domainName = req.Domain
		rr = req.RR
	}

	if rr == "" {
		return nil, fmt.Errorf("主机记录不能为空")
	}

	fmt.Printf("[DEBUG] AddDnsTxtRecord 最终参数 - DomainName: %s, RR: %s\n", domainName, rr)

	// 设置默认TTL
	if req.TTL == 0 {
		req.TTL = 600
	}

	request := &alidns.AddDomainRecordRequest{
		DomainName: tea.String(domainName),
		RR:         tea.String(rr),
		Type:       tea.String("TXT"),
		Value:      tea.String(req.Value),
		TTL:        tea.Int64(req.TTL),
	}

	runtime := &util.RuntimeOptions{}
	response, err := d.client.AddDomainRecordWithOptions(request, runtime)
	if err != nil {
		return nil, fmt.Errorf("添加DNS记录失败: %w", err)
	}

	if response.Body == nil {
		return nil, fmt.Errorf("添加DNS记录响应为空")
	}

	return &AddDnsTxtRecordResponse{
		RecordID: tea.StringValue(response.Body.RecordId),
		Message:  "DNS TXT记录添加成功",
	}, nil
}

// DeleteDnsRecordRequest 删除DNS记录请求
type DeleteDnsRecordRequest struct {
	RecordID string // 记录ID(必填)
}

// DeleteDnsRecord 删除DNS记录
func (d *DnsClient) DeleteDnsRecord(recordID string) error {
	if recordID == "" {
		return fmt.Errorf("记录ID不能为空")
	}

	request := &alidns.DeleteDomainRecordRequest{
		RecordId: tea.String(recordID),
	}

	runtime := &util.RuntimeOptions{}
	_, err := d.client.DeleteDomainRecordWithOptions(request, runtime)
	if err != nil {
		return fmt.Errorf("删除DNS记录失败: %w", err)
	}

	return nil
}

// DnsRecord DNS记录信息
type DnsRecord struct {
	RecordID   string // 记录ID
	RR         string // 主机记录
	Type       string // 记录类型
	Value      string // 记录值
	TTL        int64  // TTL
	Status     string // 状态
	DomainName string // 域名
}

// ListDnsRecordsRequest 查询DNS记录列表请求
type ListDnsRecordsRequest struct {
	Domain     string // 域名(必填)
	RR         string // 主机记录(可选,用于过滤)
	Type       string // 记录类型(可选,如: TXT)
	PageNumber int64  // 页码(可选,默认1)
	PageSize   int64  // 每页数量(可选,默认20)
}

// ListDnsRecordsResponse 查询DNS记录列表响应
type ListDnsRecordsResponse struct {
	TotalCount int64       // 总记录数
	Records    []DnsRecord // DNS记录列表
}

// ListDnsRecords 查询DNS记录列表
func (d *DnsClient) ListDnsRecords(req ListDnsRecordsRequest) (*ListDnsRecordsResponse, error) {
	if req.Domain == "" {
		return nil, fmt.Errorf("域名不能为空")
	}

	// 设置默认值
	if req.PageNumber == 0 {
		req.PageNumber = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	request := &alidns.DescribeDomainRecordsRequest{
		DomainName: tea.String(req.Domain),
		PageNumber: tea.Int64(req.PageNumber),
		PageSize:   tea.Int64(req.PageSize),
	}

	// 可选过滤条件
	if req.RR != "" {
		request.RRKeyWord = tea.String(req.RR)
	}
	if req.Type != "" {
		request.TypeKeyWord = tea.String(req.Type)
	}

	runtime := &util.RuntimeOptions{}
	response, err := d.client.DescribeDomainRecordsWithOptions(request, runtime)
	if err != nil {
		return nil, fmt.Errorf("查询DNS记录失败: %w", err)
	}

	if response.Body == nil {
		return nil, fmt.Errorf("查询DNS记录响应为空")
	}

	result := &ListDnsRecordsResponse{
		TotalCount: tea.Int64Value(response.Body.TotalCount),
		Records:    make([]DnsRecord, 0),
	}

	if response.Body.DomainRecords != nil && response.Body.DomainRecords.Record != nil {
		for _, record := range response.Body.DomainRecords.Record {
			result.Records = append(result.Records, DnsRecord{
				RecordID:   tea.StringValue(record.RecordId),
				RR:         tea.StringValue(record.RR),
				Type:       tea.StringValue(record.Type),
				Value:      tea.StringValue(record.Value),
				TTL:        tea.Int64Value(record.TTL),
				Status:     tea.StringValue(record.Status),
				DomainName: tea.StringValue(record.DomainName),
			})
		}
	}

	return result, nil
}

// GetRootDomain 从完整域名中提取根域名
// 例如: www.example.com -> example.com
//      sub.test.example.com -> example.com
//      *.example.com -> example.com (泛域名)
func GetRootDomain(fullDomain string) (string, error) {
	// 去掉尾随的点（FQDN格式）
	domain := strings.TrimSuffix(fullDomain, ".")

	// 处理泛域名：去掉 *. 前缀
	domain = strings.TrimPrefix(domain, "*.")

	parts := strings.Split(domain, ".")
	if len(parts) < 2 {
		return "", fmt.Errorf("无效的域名格式: %s", fullDomain)
	}

	// 取最后两部分作为根域名
	// 注意: 对于.com.cn等特殊后缀,可能需要更复杂的逻辑
	rootDomain := strings.Join(parts[len(parts)-2:], ".")
	return rootDomain, nil
}

// GetSubDomain 从完整域名中提取子域名部分
// 例如: www.example.com -> www
//      sub.test.example.com -> sub.test
//      *.example.com -> "" (泛域名视为根域名)
func GetSubDomain(fullDomain string) (string, error) {
	// 去掉尾随的点（FQDN格式）
	domain := strings.TrimSuffix(fullDomain, ".")

	// 处理泛域名：去掉 *. 前缀
	domain = strings.TrimPrefix(domain, "*.")

	parts := strings.Split(domain, ".")
	if len(parts) < 2 {
		return "", fmt.Errorf("无效的域名格式: %s", fullDomain)
	}

	// 如果只有两部分,说明没有子域名
	if len(parts) == 2 {
		return "", nil
	}

	// 取除了最后两部分的所有部分
	subDomain := strings.Join(parts[:len(parts)-2], ".")
	return subDomain, nil
}

// AddCertValidationRecord 添加证书验证DNS记录
// fullDomain: 完整域名(如: www.example.com)
// validationValue: 验证值
func (d *DnsClient) AddCertValidationRecord(fullDomain, validationValue string) (*AddDnsTxtRecordResponse, error) {
	// 提取根域名
	rootDomain, err := GetRootDomain(fullDomain)
	if err != nil {
		return nil, err
	}

	// 提取子域名
	subDomain, err := GetSubDomain(fullDomain)
	if err != nil {
		return nil, err
	}

	// 构造验证记录的RR
	var rr string
	if subDomain == "" {
		rr = "_acme-challenge"
	} else {
		rr = "_acme-challenge." + subDomain
	}

	// 添加TXT记录
	return d.AddDnsTxtRecord(AddDnsTxtRecordRequest{
		Domain: rootDomain,
		RR:     rr,
		Value:  validationValue,
		TTL:    600,
	})
}

// CleanupCertValidationRecord 清理证书验证DNS记录
// fullDomain: ACME验证域名，如 _acme-challenge.example.com
func (d *DnsClient) CleanupCertValidationRecord(fullDomain string) error {
	// 先去掉尾随点，统一处理
	cleanDomain := strings.TrimSuffix(fullDomain, ".")

	// 从验证域名中提取根域名和RR
	// 例如: _acme-challenge.example.com -> 根域名: example.com, RR: _acme-challenge
	rootDomain, err := GetRootDomain(cleanDomain)
	if err != nil {
		return fmt.Errorf("解析根域名失败: %w", err)
	}

	// 计算RR：cleanDomain去掉根域名部分
	// 例如: _acme-challenge.example.com -> _acme-challenge
	//      _acme-challenge.sub.example.com -> _acme-challenge.sub
	rr := strings.TrimSuffix(cleanDomain, "."+rootDomain)
	if rr == cleanDomain {
		// 如果没有后缀可以去除，说明cleanDomain格式不正确
		return fmt.Errorf("无法从验证域名中提取主机记录: %s", fullDomain)
	}

	fmt.Printf("[DEBUG] 清理DNS记录 - 根域名: %s, RR: %s\n", rootDomain, rr)

	// 查询所有匹配的TXT记录
	records, err := d.ListDnsRecords(ListDnsRecordsRequest{
		Domain:     rootDomain,
		RR:         rr,
		Type:       "TXT",
		PageNumber: 1,
		PageSize:   100,
	})
	if err != nil {
		// 如果查询失败，打印日志但不返回错误（可能是记录不存在）
		fmt.Printf("[WARN] 查询DNS记录失败（可能记录不存在）: %v\n", err)
		return nil
	}

	// 如果没有找到匹配的记录，直接返回
	if len(records.Records) == 0 {
		fmt.Printf("[INFO] 未找到需要清理的DNS记录 - 根域名: %s, RR: %s\n", rootDomain, rr)
		return nil
	}

	// 删除所有匹配的记录
	for _, record := range records.Records {
		fmt.Printf("[INFO] 删除DNS记录 - RecordID: %s, RR: %s, Value: %s\n", record.RecordID, record.RR, record.Value)
		if err := d.DeleteDnsRecord(record.RecordID); err != nil {
			// 记录错误但继续删除其他记录
			fmt.Printf("[ERROR] 删除DNS记录失败 RecordID=%s: %v\n", record.RecordID, err)
		}
	}

	return nil
}

// DomainInfo 域名信息
type DomainInfo struct {
	DomainName string // 域名
	RecordCount int64  // 记录数
}

// ListDomainsRequest 查询域名列表请求
type ListDomainsRequest struct {
	PageNumber int64 // 页码(可选,默认1)
	PageSize   int64 // 每页数量(可选,默认20,最大100)
}

// ListDomainsResponse 查询域名列表响应
type ListDomainsResponse struct {
	TotalCount int64        // 总域名数
	Domains    []DomainInfo // 域名列表
}

// ListDomains 查询域名列表
func (d *DnsClient) ListDomains(req ListDomainsRequest) (*ListDomainsResponse, error) {
	// 设置默认值
	if req.PageNumber == 0 {
		req.PageNumber = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	// 阿里云限制最大100
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	request := &alidns.DescribeDomainsRequest{
		PageNumber: tea.Int64(req.PageNumber),
		PageSize:   tea.Int64(req.PageSize),
	}

	runtime := &util.RuntimeOptions{}
	response, err := d.client.DescribeDomainsWithOptions(request, runtime)
	if err != nil {
		return nil, fmt.Errorf("查询域名列表失败: %w", err)
	}

	if response.Body == nil {
		return nil, fmt.Errorf("查询域名列表响应为空")
	}

	result := &ListDomainsResponse{
		TotalCount: tea.Int64Value(response.Body.TotalCount),
		Domains:    make([]DomainInfo, 0),
	}

	if response.Body.Domains != nil && response.Body.Domains.Domain != nil {
		for _, domain := range response.Body.Domains.Domain {
			result.Domains = append(result.Domains, DomainInfo{
				DomainName:  tea.StringValue(domain.DomainName),
				RecordCount: tea.Int64Value(domain.RecordCount),
			})
		}
	}

	return result, nil
}

// ListAllDomains 查询所有域名(自动分页)
func (d *DnsClient) ListAllDomains() ([]DomainInfo, error) {
	allDomains := make([]DomainInfo, 0)
	pageNumber := int64(1)
	pageSize := int64(100) // 使用最大页大小减少请求次数

	for {
		response, err := d.ListDomains(ListDomainsRequest{
			PageNumber: pageNumber,
			PageSize:   pageSize,
		})
		if err != nil {
			return nil, err
		}

		allDomains = append(allDomains, response.Domains...)

		// 如果已经获取所有域名，退出循环
		if int64(len(allDomains)) >= response.TotalCount {
			break
		}

		pageNumber++
	}

	return allDomains, nil
}

// ListAllSubdomains 查询域名下的所有子域名(自动分页)
// 返回完整的域名列表,如: www.example.com, api.example.com等
func (d *DnsClient) ListAllSubdomains(rootDomain string) ([]string, error) {
	allSubdomains := make([]string, 0)
	pageNumber := int64(1)
	pageSize := int64(100) // 使用最大页大小

	// 记录已添加的子域名,避免重复
	subdomainSet := make(map[string]bool)

	for {
		response, err := d.ListDnsRecords(ListDnsRecordsRequest{
			Domain:     rootDomain,
			PageNumber: pageNumber,
			PageSize:   pageSize,
		})
		if err != nil {
			return nil, err
		}

		// 遍历所有DNS记录,构建完整域名
		for _, record := range response.Records {
			// 只处理A记录和CNAME记录(有实际访问意义的记录)
			if record.Type != "A" && record.Type != "CNAME" {
				continue
			}

			var fullDomain string
			if record.RR == "@" {
				// @ 表示根域名本身
				fullDomain = rootDomain
			} else {
				// 构建完整域名: RR.rootDomain
				fullDomain = record.RR + "." + rootDomain
			}

			// 去重添加
			if !subdomainSet[fullDomain] {
				subdomainSet[fullDomain] = true
				allSubdomains = append(allSubdomains, fullDomain)
			}
		}

		// 如果已经获取所有记录,退出循环
		if int64(pageNumber*pageSize) >= response.TotalCount {
			break
		}

		pageNumber++
	}

	return allSubdomains, nil
}

// ListAllDomainsWithSubdomains 查询所有域名及其子域名
// 返回所有完整域名列表
func (d *DnsClient) ListAllDomainsWithSubdomains() ([]string, error) {
	// 第一步: 获取所有主域名
	domainInfos, err := d.ListAllDomains()
	if err != nil {
		return nil, fmt.Errorf("获取主域名列表失败: %w", err)
	}

	if len(domainInfos) == 0 {
		return []string{}, nil
	}

	allDomains := make([]string, 0)

	// 第二步: 遍历每个主域名,获取其子域名
	for _, domainInfo := range domainInfos {
		rootDomain := domainInfo.DomainName

		// 获取该主域名下的所有子域名
		subdomains, err := d.ListAllSubdomains(rootDomain)
		if err != nil {
			// 如果获取子域名失败,只记录错误,继续处理其他域名
			fmt.Printf("[WARN] 获取域名 %s 的子域名失败: %v\n", rootDomain, err)
			// 至少添加主域名本身
			allDomains = append(allDomains, rootDomain)
			continue
		}

		// 如果该域名没有任何DNS记录,至少添加主域名本身
		if len(subdomains) == 0 {
			allDomains = append(allDomains, rootDomain)
		} else {
			allDomains = append(allDomains, subdomains...)
		}
	}

	return allDomains, nil
}
