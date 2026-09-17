package util

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"
	"time"
)

// CertInfo 证书信息
type CertInfo struct {
	Subject    string    // 证书主题
	Issuer     string    // 颁发者
	NotBefore  time.Time // 生效时间
	NotAfter   time.Time // 过期时间
	DaysLeft   int       // 剩余天数
	Algorithm  string    // 算法
	CommonName string    // 通用名称
}

// ParseCertificate 解析证书内容
func ParseCertificate(certContent string) (*CertInfo, error) {
	// 解析PEM格式证书
	block, _ := pem.Decode([]byte(certContent))
	if block == nil {
		return nil, fmt.Errorf("无效的PEM格式证书")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("证书解析失败: %v", err)
	}

	// 计算剩余天数
	daysLeft := int(time.Until(cert.NotAfter).Hours() / 24)

	// 提取算法信息
	algorithm := "RSA2048"
	if cert.PublicKeyAlgorithm == x509.ECDSA {
		algorithm = "EC256"
	}

	return &CertInfo{
		Subject:    cert.Subject.String(),
		Issuer:     cert.Issuer.Organization[0],
		NotBefore:  cert.NotBefore,
		NotAfter:   cert.NotAfter,
		DaysLeft:   daysLeft,
		Algorithm:  algorithm,
		CommonName: cert.Subject.CommonName,
	}, nil
}

// ExtractDomainFromFilename 从文件名提取域名
// 例如: lockin.com.key -> lockin.com
//       www.example.com.pem -> www.example.com
//       test.domain.com.crt -> test.domain.com
func ExtractDomainFromFilename(filename string) string {
	// 移除路径,只保留文件名
	if idx := strings.LastIndex(filename, "/"); idx != -1 {
		filename = filename[idx+1:]
	}
	if idx := strings.LastIndex(filename, "\\"); idx != -1 {
		filename = filename[idx+1:]
	}

	// 移除常见的证书文件扩展名
	extensions := []string{".pem", ".crt", ".cer", ".key", ".cert", ".certificate"}
	for _, ext := range extensions {
		filename = strings.TrimSuffix(filename, ext)
	}

	// 移除可能的后缀
	suffixes := []string{"_private", "_public", "-key", "-cert", "_key", "_cert"}
	for _, suffix := range suffixes {
		filename = strings.TrimSuffix(filename, suffix)
	}

	return strings.TrimSpace(filename)
}
