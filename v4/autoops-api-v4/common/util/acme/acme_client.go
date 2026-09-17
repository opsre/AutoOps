// ACME客户端封装 - 支持Let's Encrypt、ZeroSSL等CA
// 基于 go-acme/lego 库实现
package acme

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"
	"time"

	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/challenge/dns01"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
)

// AcmeClient ACME客户端
type AcmeClient struct {
	client     *lego.Client
	email      string
	caProvider string
}

// AcmeConfig ACME配置
type AcmeConfig struct {
	Email      string // 联系邮箱
	CAProvider string // CA提供商: LetsEncrypt/ZeroSSL/Google/SSLCOM/BuyPass
	Algorithm  string // 加密算法: RSA2048/RSA4096/EC256/EC384
	EABKid     string // EAB Key ID (ZeroSSL等需要)
	EABHmacKey string // EAB HMAC Key (ZeroSSL等需要)
}

// AcmeUser 实现ACME用户接口
type AcmeUser struct {
	Email        string
	Registration *registration.Resource
	key          crypto.PrivateKey
}

func (u *AcmeUser) GetEmail() string {
	return u.Email
}

func (u *AcmeUser) GetRegistration() *registration.Resource {
	return u.Registration
}

func (u *AcmeUser) GetPrivateKey() crypto.PrivateKey {
	return u.key
}

// CA服务器URL映射
var caURLMap = map[string]string{
	"LetsEncrypt": lego.LEDirectoryProduction,
	"ZeroSSL":     "https://acme.zerossl.com/v2/DV90",
	"Google":      "https://dv.acme-v02.api.pki.goog/directory",
	"SSLCOM":      "https://acme.ssl.com/sslcom-dv-rsa",
	"BuyPass":     "https://api.buypass.com/acme/directory",
}

// NewAcmeClient 创建ACME客户端
func NewAcmeClient(config AcmeConfig) (*AcmeClient, error) {
	// 1. 生成私钥
	var privateKey crypto.PrivateKey
	var err error

	algorithm := config.Algorithm
	if algorithm == "" {
		algorithm = "RSA2048" // 默认使用RSA2048
	}

	switch algorithm {
	case "RSA2048":
		privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	case "RSA4096":
		privateKey, err = rsa.GenerateKey(rand.Reader, 4096)
	case "EC256":
		privateKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case "EC384":
		privateKey, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	default:
		return nil, fmt.Errorf("不支持的加密算法: %s", algorithm)
	}

	if err != nil {
		return nil, fmt.Errorf("生成私钥失败: %v", err)
	}

	// 2. 创建ACME用户
	acmeUser := &AcmeUser{
		Email: config.Email,
		key:   privateKey,
	}

	// 3. 创建lego配置
	legoConfig := lego.NewConfig(acmeUser)

	// 设置CA服务器
	caProvider := config.CAProvider
	if caProvider == "" {
		caProvider = "LetsEncrypt" // 默认使用Let's Encrypt
	}

	caURL, ok := caURLMap[caProvider]
	if !ok {
		return nil, fmt.Errorf("不支持的CA提供商: %s", caProvider)
	}
	legoConfig.CADirURL = caURL

	// 设置证书类型
	switch algorithm {
	case "RSA2048", "RSA4096":
		legoConfig.Certificate.KeyType = certcrypto.RSA2048
		if algorithm == "RSA4096" {
			legoConfig.Certificate.KeyType = certcrypto.RSA4096
		}
	case "EC256":
		legoConfig.Certificate.KeyType = certcrypto.EC256
	case "EC384":
		legoConfig.Certificate.KeyType = certcrypto.EC384
	}

	// 4. 创建lego客户端
	client, err := lego.NewClient(legoConfig)
	if err != nil {
		return nil, fmt.Errorf("创建ACME客户端失败: %v", err)
	}

	// 5. 注册用户（支持EAB）
	var reg *registration.Resource

	// 如果提供了EAB凭证，使用EAB注册（ZeroSSL等需要）
	if config.EABKid != "" && config.EABHmacKey != "" {
		fmt.Printf("[INFO] 使用EAB注册ACME用户 - Kid: %s\n", config.EABKid)
		reg, err = client.Registration.RegisterWithExternalAccountBinding(registration.RegisterEABOptions{
			TermsOfServiceAgreed: true,
			Kid:                  config.EABKid,
			HmacEncoded:          config.EABHmacKey,
		})
	} else {
		// 普通注册（Let's Encrypt等）
		reg, err = client.Registration.Register(registration.RegisterOptions{
			TermsOfServiceAgreed: true,
		})
	}

	if err != nil {
		return nil, fmt.Errorf("注册ACME用户失败: %v", err)
	}
	acmeUser.Registration = reg

	return &AcmeClient{
		client:     client,
		email:      config.Email,
		caProvider: caProvider,
	}, nil
}

// DNSProviderFunc DNS提供商函数类型
type DNSProviderFunc func(fqdn, value string) error

// CustomDNSProvider 自定义DNS提供商
type CustomDNSProvider struct {
	addFunc    DNSProviderFunc
	cleanupFunc DNSProviderFunc
}

func (p *CustomDNSProvider) Present(domain, token, keyAuth string) error {
	fqdn, value := dns01.GetRecord(domain, keyAuth)
	return p.addFunc(fqdn, value)
}

func (p *CustomDNSProvider) CleanUp(domain, token, keyAuth string) error {
	fqdn, _ := dns01.GetRecord(domain, "")
	return p.cleanupFunc(fqdn, "")
}

// SetDNSProvider 设置DNS提供商
func (c *AcmeClient) SetDNSProvider(addFunc, cleanupFunc DNSProviderFunc) error {
	provider := &CustomDNSProvider{
		addFunc:    addFunc,
		cleanupFunc: cleanupFunc,
	}
	return c.client.Challenge.SetDNS01Provider(provider)
}

// ApplyCertRequest 申请证书请求
type ApplyCertRequest struct {
	Domain    string // 域名(支持通配符 *.example.com)
	Email     string // 联系邮箱
}

// ApplyCertResponse 申请证书响应
type ApplyCertResponse struct {
	Domain      string    // 域名
	Cert        string    // 证书内容
	PrivateKey  string    // 私钥内容
	IssuerCert  string    // 颁发者证书
	CertURL     string    // 证书URL
	IssueDate   string    // 签发日期
	ExpireDate  string    // 过期日期
	DaysLeft    int       // 剩余天数
}

// ApplyCert 申请证书
func (c *AcmeClient) ApplyCert(req ApplyCertRequest) (*ApplyCertResponse, error) {
	// 申请证书
	request := certificate.ObtainRequest{
		Domains: []string{req.Domain},
		Bundle:  true, // Bundle=true 表示返回完整证书链(域名证书+中间证书)
	}

	certificates, err := c.client.Certificate.Obtain(request)
	if err != nil {
		return nil, fmt.Errorf("申请证书失败: %v", err)
	}

	// 注意: 当Bundle=true时
	// - certificates.Certificate 包含完整证书链(域名证书 + 中间证书)
	// - certificates.IssuerCertificate 仅包含中间证书
	// 为了兼容Nginx等Web服务器,我们需要确保证书文件包含完整链

	// 解析证书获取有效期(只解析第一个证书块,即域名证书)
	block, _ := pem.Decode(certificates.Certificate)
	if block == nil {
		return nil, fmt.Errorf("解析证书失败")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("解析证书失败: %v", err)
	}

	now := time.Now()
	daysLeft := int(cert.NotAfter.Sub(now).Hours() / 24)

	// 确保证书链完整性
	// certificates.Certificate 理论上已包含完整链,但为了保险,我们显式合并
	fullChainCert := string(certificates.Certificate)

	// 如果IssuerCertificate不为空,并且Certificate中不包含它,则追加
	if len(certificates.IssuerCertificate) > 0 {
		issuerCertStr := string(certificates.IssuerCertificate)
		// 简单检查: 如果Certificate不包含IssuerCertificate的内容,则追加
		if !strings.Contains(fullChainCert, strings.TrimSpace(issuerCertStr)) {
			if !strings.HasSuffix(fullChainCert, "\n") {
				fullChainCert += "\n"
			}
			fullChainCert += issuerCertStr
		}
	}

	return &ApplyCertResponse{
		Domain:     req.Domain,
		Cert:       fullChainCert, // 完整证书链
		PrivateKey: string(certificates.PrivateKey),
		IssuerCert: string(certificates.IssuerCertificate), // 单独的中间证书
		CertURL:    certificates.CertURL,
		IssueDate:  cert.NotBefore.Format("2006-01-02"),
		ExpireDate: cert.NotAfter.Format("2006-01-02"),
		DaysLeft:   daysLeft,
	}, nil
}

// RenewCert 续期证书
func (c *AcmeClient) RenewCert(certPEM, privateKeyPEM []byte) (*ApplyCertResponse, error) {
	certResource := &certificate.Resource{
		Certificate:   certPEM,
		PrivateKey:    privateKeyPEM,
	}

	certificates, err := c.client.Certificate.Renew(*certResource, true, false, "")
	if err != nil {
		return nil, fmt.Errorf("续期证书失败: %v", err)
	}

	// 解析证书获取有效期
	block, _ := pem.Decode(certificates.Certificate)
	if block == nil {
		return nil, fmt.Errorf("解析证书失败")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("解析证书失败: %v", err)
	}

	now := time.Now()
	daysLeft := int(cert.NotAfter.Sub(now).Hours() / 24)

	// 提取域名
	domain := cert.Subject.CommonName
	if len(cert.DNSNames) > 0 {
		domain = cert.DNSNames[0]
	}

	// 确保证书链完整性
	fullChainCert := string(certificates.Certificate)
	if len(certificates.IssuerCertificate) > 0 {
		issuerCertStr := string(certificates.IssuerCertificate)
		if !strings.Contains(fullChainCert, strings.TrimSpace(issuerCertStr)) {
			if !strings.HasSuffix(fullChainCert, "\n") {
				fullChainCert += "\n"
			}
			fullChainCert += issuerCertStr
		}
	}

	return &ApplyCertResponse{
		Domain:     domain,
		Cert:       fullChainCert, // 完整证书链
		PrivateKey: string(certificates.PrivateKey),
		IssuerCert: string(certificates.IssuerCertificate),
		CertURL:    certificates.CertURL,
		IssueDate:  cert.NotBefore.Format("2006-01-02"),
		ExpireDate: cert.NotAfter.Format("2006-01-02"),
		DaysLeft:   daysLeft,
	}, nil
}

// GetDNSRecord 获取DNS验证记录信息
func GetDNSRecord(domain string) (recordName string, recordType string) {
	// 移除通配符前缀
	domain = strings.TrimPrefix(domain, "*.")

	// DNS-01挑战记录名称格式: _acme-challenge.domain.com
	recordName = "_acme-challenge." + domain
	recordType = "TXT"

	return recordName, recordType
}
