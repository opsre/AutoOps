package aliyun

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"

	cas "github.com/alibabacloud-go/cas-20200407/v4/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
)

// CasClient 阿里云CAS客户端封装
type CasClient struct {
	client *cas.Client
	config CasConfig
}

// CasConfig 阿里云CAS配置
type CasConfig struct {
	AccessKey    string
	AccessSecret string
	Endpoint     string // 默认: cas.aliyuncs.com
}

// NewCasClient 创建阿里云CAS客户端
func NewCasClient(config CasConfig) (*CasClient, error) {
	if config.Endpoint == "" {
		config.Endpoint = "cas.aliyuncs.com"
	}

	openApiConfig := &openapi.Config{
		AccessKeyId:     tea.String(config.AccessKey),
		AccessKeySecret: tea.String(config.AccessSecret),
		Endpoint:        tea.String(config.Endpoint),
	}

	client, err := cas.NewClient(openApiConfig)
	if err != nil {
		return nil, fmt.Errorf("创建CAS客户端失败: %w", err)
	}

	return &CasClient{
		client: client,
		config: config,
	}, nil
}

// ApplyFreeCertRequest 申请免费证书请求参数
type ApplyFreeCertRequest struct {
	Domain       string // 域名(必填)
	Email        string // 联系邮箱(必填)
	Phone        string // 联系电话(可选)
	Username     string // 联系人姓名(可选)
	ValidateType string // 验证方式: DNS 或 FILE (默认DNS)
	ProductCode  string // 产品代码(默认: digicert-free-1-free)
}

// ApplyFreeCertResponse 申请免费证书响应
type ApplyFreeCertResponse struct {
	OrderID int64  // 订单ID
	Message string // 响应消息
}

// ApplyFreeCert 申请阿里云免费SSL证书
func (c *CasClient) ApplyFreeCert(req ApplyFreeCertRequest) (*ApplyFreeCertResponse, error) {
	// 设置默认值
	if req.ValidateType == "" {
		req.ValidateType = "DNS"
	}
	if req.ProductCode == "" {
		req.ProductCode = "digicert-free-1-free"
	}

	// 参数验证
	if req.Domain == "" {
		return nil, fmt.Errorf("域名不能为空")
	}
	if req.Email == "" {
		return nil, fmt.Errorf("邮箱不能为空")
	}

	// 创建证书申请请求
	request := &cas.CreateCertificateForPackageRequestRequest{
		ProductCode:  tea.String(req.ProductCode),
		Domain:       tea.String(req.Domain),
		Email:        tea.String(req.Email),
		ValidateType: tea.String(req.ValidateType),
	}

	// 可选参数
	if req.Phone != "" {
		request.Phone = tea.String(req.Phone)
	}
	if req.Username != "" {
		request.Username = tea.String(req.Username)
	}

	runtime := &util.RuntimeOptions{}

	// 提交证书申请
	response, err := c.client.CreateCertificateForPackageRequestWithOptions(request, runtime)
	if err != nil {
		return nil, fmt.Errorf("申请证书失败: %w", err)
	}

	if response.Body == nil {
		return nil, fmt.Errorf("证书申请响应为空")
	}

	return &ApplyFreeCertResponse{
		OrderID: tea.Int64Value(response.Body.OrderId),
		Message: "证书申请已提交",
	}, nil
}

// CertStateResponse 证书状态响应
type CertStateResponse struct {
	OrderID      int64  // 订单ID
	Type         string // 证书类型
	Status       string // 证书状态
	Domain       string // 域名
	ValidateType string // 验证方式
	ValidateInfo string // 验证信息(DNS记录或FILE内容)
}

// GetCertState 查询证书状态和验证信息
func (c *CasClient) GetCertState(orderID int64) (*CertStateResponse, error) {
	request := &cas.DescribeCertificateStateRequest{
		OrderId: tea.Int64(orderID),
	}

	runtime := &util.RuntimeOptions{}
	response, err := c.client.DescribeCertificateStateWithOptions(request, runtime)
	if err != nil {
		return nil, fmt.Errorf("查询证书状态失败: %w", err)
	}

	if response.Body == nil {
		return nil, fmt.Errorf("证书状态响应为空")
	}

	return &CertStateResponse{
		OrderID:      orderID,
		Type:         tea.StringValue(response.Body.Type),
		Status:       tea.StringValue(response.Body.Type), // Type字段表示状态
		Domain:       tea.StringValue(response.Body.Domain),
		ValidateType: tea.StringValue(response.Body.ValidateType),
		ValidateInfo: tea.StringValue(response.Body.Content),
	}, nil
}

// DownloadCertResponse 下载证书响应
type DownloadCertResponse struct {
	CertID      int64  // 证书ID
	Cert        string // 证书内容(PEM格式)
	PrivateKey  string // 私钥内容(PEM格式)
	CertChain   string // 证书链
	Domain      string // 域名
	IssueDate   string // 签发日期
	ExpireDate  string // 过期日期
	DaysLeft    int    // 剩余天数
	Fingerprint string // 证书指纹
}

// DownloadCert 下载已签发的证书（通过CertID）
func (c *CasClient) DownloadCert(certID int64) (*DownloadCertResponse, error) {
	request := &cas.GetUserCertificateDetailRequest{
		CertId: tea.Int64(certID),
	}

	runtime := &util.RuntimeOptions{}
	response, err := c.client.GetUserCertificateDetailWithOptions(request, runtime)
	if err != nil {
		return nil, fmt.Errorf("下载证书失败: %w", err)
	}

	if response.Body == nil {
		return nil, fmt.Errorf("证书下载响应为空")
	}

	// 计算剩余天数
	daysLeft := 0
	if response.Body.EndDate != nil {
		endDate := tea.StringValue(response.Body.EndDate)
		if expireTime, err := time.Parse("2006-01-02", endDate); err == nil {
			daysLeft = int(time.Until(expireTime).Hours() / 24)
		}
	}

	// 处理证书链信息(仅保存发行者信息)
	certChain := ""
	if len(response.Body.CertChain) > 0 && response.Body.CertChain[0] != nil {
		certChain = tea.StringValue(response.Body.CertChain[0].IssuerCommonName)
	}

	return &DownloadCertResponse{
		CertID:      certID,
		Cert:        tea.StringValue(response.Body.Cert),
		PrivateKey:  tea.StringValue(response.Body.Key),
		CertChain:   certChain,
		Domain:      tea.StringValue(response.Body.Name), // 使用Name字段
		IssueDate:   tea.StringValue(response.Body.StartDate),
		ExpireDate:  tea.StringValue(response.Body.EndDate),
		DaysLeft:    daysLeft,
		Fingerprint: tea.StringValue(response.Body.Fingerprint),
	}, nil
}

// CertMetadata 证书元信息（不包含证书内容）
type CertMetadata struct {
	CertID     int64  // 证书ID
	Domain     string // 域名
	IssueDate  string // 签发日期
	ExpireDate string // 过期日期
	DaysLeft   int    // 剩余天数
}

// GetCertMetadata 获取证书元信息（仅有效期，不下载证书内容）
func (c *CasClient) GetCertMetadata(orderID int64) (*CertMetadata, error) {
	// 1. 先通过DescribeCertificateState获取证书信息（轻量级API）
	stateReq := &cas.DescribeCertificateStateRequest{
		OrderId: tea.Int64(orderID),
	}

	runtime := &util.RuntimeOptions{}
	stateResp, err := c.client.DescribeCertificateStateWithOptions(stateReq, runtime)
	if err != nil {
		return nil, fmt.Errorf("查询证书状态失败: %w", err)
	}

	if stateResp.Body == nil {
		return nil, fmt.Errorf("证书状态响应为空")
	}

	// 2. 检查是否有证书内容（已签发）
	if stateResp.Body.Certificate == nil || tea.StringValue(stateResp.Body.Certificate) == "" {
		return nil, fmt.Errorf("证书内容为空，可能尚未签发")
	}

	// 3. 从证书内容解析有效期（临时解析，不保存证书内容）
	cert := tea.StringValue(stateResp.Body.Certificate)
	issueDate, expireDate, daysLeft := parseCertificateDates(cert)

	return &CertMetadata{
		CertID:     0, // DescribeCertificateState不返回CertID
		Domain:     tea.StringValue(stateResp.Body.Domain),
		IssueDate:  issueDate,
		ExpireDate: expireDate,
		DaysLeft:   daysLeft,
	}, nil
}

// DownloadCertByOrderID 通过OrderID下载已签发的证书
func (c *CasClient) DownloadCertByOrderID(orderID int64) (*DownloadCertResponse, error) {
	// 使用DescribeCertificateStateRequest API来获取证书
	request := &cas.DescribeCertificateStateRequest{
		OrderId: tea.Int64(orderID),
	}

	runtime := &util.RuntimeOptions{}
	stateResp, err := c.client.DescribeCertificateStateWithOptions(request, runtime)
	if err != nil {
		return nil, fmt.Errorf("查询证书状态失败: %w", err)
	}

	if stateResp.Body == nil {
		return nil, fmt.Errorf("证书状态响应为空")
	}

	// 获取证书内容
	if stateResp.Body.Certificate == nil || tea.StringValue(stateResp.Body.Certificate) == "" {
		return nil, fmt.Errorf("证书内容为空，可能尚未签发")
	}

	// 从证书内容解析有效期
	cert := tea.StringValue(stateResp.Body.Certificate)
	issueDate, expireDate, daysLeft := parseCertificateDates(cert)

	return &DownloadCertResponse{
		CertID:     0, // OrderID方式可能不返回CertID
		Cert:       cert,
		PrivateKey: tea.StringValue(stateResp.Body.PrivateKey),
		CertChain:  "",
		Domain:     tea.StringValue(stateResp.Body.Domain),
		IssueDate:  issueDate,
		ExpireDate: expireDate,
		DaysLeft:   daysLeft,
	}, nil
}

// UploadCertRequest 上传证书请求
type UploadCertRequest struct {
	Name       string // 证书名称(必填)
	Cert       string // 证书内容(必填)
	PrivateKey string // 私钥内容(必填)
}

// UploadCertResponse 上传证书响应
type UploadCertResponse struct {
	CertID  int64  // 证书ID
	Message string // 响应消息
}

// UploadCert 上传证书到阿里云CAS
func (c *CasClient) UploadCert(req UploadCertRequest) (*UploadCertResponse, error) {
	// 参数验证
	if req.Name == "" {
		return nil, fmt.Errorf("证书名称不能为空")
	}
	if req.Cert == "" {
		return nil, fmt.Errorf("证书内容不能为空")
	}
	if req.PrivateKey == "" {
		return nil, fmt.Errorf("私钥内容不能为空")
	}

	request := &cas.UploadUserCertificateRequest{
		Name: tea.String(req.Name),
		Cert: tea.String(req.Cert),
		Key:  tea.String(req.PrivateKey),
	}

	runtime := &util.RuntimeOptions{}
	response, err := c.client.UploadUserCertificateWithOptions(request, runtime)
	if err != nil {
		return nil, fmt.Errorf("上传证书失败: %w", err)
	}

	if response.Body == nil {
		return nil, fmt.Errorf("上传证书响应为空")
	}

	return &UploadCertResponse{
		CertID:  tea.Int64Value(response.Body.CertId),
		Message: "证书上传成功",
	}, nil
}

// ListCertsResponse 证书列表响应
type ListCertsResponse struct {
	TotalCount int        // 总数量
	CertList   []CertInfo // 证书列表
}

// CertInfo 证书信息
type CertInfo struct {
	CertID      int64  // 证书ID
	Name        string // 证书名称
	Domain      string // 域名
	IssueDate   string // 签发日期
	ExpireDate  string // 过期日期
	Status      string // 状态
	Source      string // 来源
	Fingerprint string // 证书指纹
}

// ListCerts 查询证书列表
func (c *CasClient) ListCerts(currentPage, pageSize int) (*ListCertsResponse, error) {
	request := &cas.ListUserCertificateOrderRequest{
		CurrentPage: tea.Int64(int64(currentPage)),
		ShowSize:    tea.Int64(int64(pageSize)),
		OrderType:   tea.String("CERT"), // 证书类型
		Status:      tea.String("ISSUED"), // 只查询已签发的证书
	}

	runtime := &util.RuntimeOptions{}
	response, err := c.client.ListUserCertificateOrderWithOptions(request, runtime)
	if err != nil {
		return nil, fmt.Errorf("查询证书列表失败: %w", err)
	}

	if response.Body == nil {
		return nil, fmt.Errorf("证书列表响应为空")
	}

	result := &ListCertsResponse{
		TotalCount: int(tea.Int64Value(response.Body.TotalCount)),
		CertList:   make([]CertInfo, 0),
	}

	if response.Body.CertificateOrderList != nil {
		for _, cert := range response.Body.CertificateOrderList {
			// 时间戳转换为字符串
			issueDate := ""
			expireDate := ""
			if cert.CertStartTime != nil {
				t := time.Unix(tea.Int64Value(cert.CertStartTime)/1000, 0)
				issueDate = t.Format("2006-01-02")
			}
			if cert.CertEndTime != nil {
				t := time.Unix(tea.Int64Value(cert.CertEndTime)/1000, 0)
				expireDate = t.Format("2006-01-02")
			}

			result.CertList = append(result.CertList, CertInfo{
				CertID:      tea.Int64Value(cert.CertificateId),
				Name:        tea.StringValue(cert.Name),
				Domain:      tea.StringValue(cert.Domain),
				IssueDate:   issueDate,
				ExpireDate:  expireDate,
				Status:      tea.StringValue(cert.Status),
				Source:      tea.StringValue(cert.SourceType), // 使用SourceType字段
				Fingerprint: tea.StringValue(cert.Fingerprint),
			})
		}
	}

	return result, nil
}

// DeleteCert 删除证书
func (c *CasClient) DeleteCert(certID int64) error {
	request := &cas.DeleteUserCertificateRequest{
		CertId: tea.Int64(certID),
	}

	runtime := &util.RuntimeOptions{}
	_, err := c.client.DeleteUserCertificateWithOptions(request, runtime)
	if err != nil {
		return fmt.Errorf("删除证书失败: %w", err)
	}

	return nil
}

// parseCertificateDates 从证书PEM内容中解析有效期
func parseCertificateDates(certPEM string) (issueDate, expireDate string, daysLeft int) {
	// 导入必要的包来解析证书
	block, _ := pem.Decode([]byte(certPEM))
	if block == nil {
		return "", "", 0
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return "", "", 0
	}

	// 格式化日期
	issueDate = cert.NotBefore.Format("2006-01-02")
	expireDate = cert.NotAfter.Format("2006-01-02")

	// 计算剩余天数
	daysLeft = int(time.Until(cert.NotAfter).Hours() / 24)
	if daysLeft < 0 {
		daysLeft = 0
	}

	return issueDate, expireDate, daysLeft
}
