// LDAP认证工具
package util

import (
	"crypto/tls"
	"fmt"
	"strings"
	"time"

	"github.com/go-ldap/ldap/v3"
)

// LDAPAuthConfig LDAP认证配置
type LDAPAuthConfig struct {
	Host       string
	Port       int
	BaseDn     string
	BindUser   string
	BindPass   string
	AuthFilter string
	TLS        bool
	StartTLS   bool
}

// LDAPUserAttributes LDAP用户属性
type LDAPUserAttributes struct {
	Username string
	CN       string
	Mail     string
	Mobile   string
	Attrs    map[string]string
}

// LDAPAuth LDAP认证
func LDAPAuth(config *LDAPAuthConfig, username, password string) (*LDAPUserAttributes, error) {
	// 1. 连接LDAP服务器
	// 清理Host中的空格
	host := strings.TrimSpace(config.Host)
	ldapURL := fmt.Sprintf("%s:%d", host, config.Port)

	var l *ldap.Conn
	var err error

	if config.TLS {
		// 使用TLS加密连接
		l, err = ldap.DialTLS("tcp", ldapURL, &tls.Config{InsecureSkipVerify: true})
	} else {
		// 使用DialURL代替Dial
		l, err = ldap.DialURL(fmt.Sprintf("ldap://%s", ldapURL))
	}

	if err != nil {
		return nil, fmt.Errorf("LDAP连接失败: %v", err)
	}
	// 注意：不使用defer，手动在用户认证前关闭

	// 设置操作超时
	l.SetTimeout(10 * time.Second)

	// 2. 如果需要StartTLS，执行升级
	if config.StartTLS && !config.TLS {
		err = l.StartTLS(&tls.Config{InsecureSkipVerify: true})
		if err != nil {
			return nil, fmt.Errorf("StartTLS失败: %v", err)
		}
	}

	// 3. 使用管理员账号绑定
	err = l.Bind(config.BindUser, config.BindPass)
	if err != nil {
		return nil, fmt.Errorf("LDAP绑定失败: %v", err)
	}

	// 4. 搜索用户
	searchFilter := strings.ReplaceAll(config.AuthFilter, "%s", username)
	searchRequest := ldap.NewSearchRequest(
		config.BaseDn,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		searchFilter,
		[]string{"cn", "mail", "mobile", "telephoneNumber", "displayName"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("LDAP搜索失败: %v", err)
	}

	if len(sr.Entries) == 0 {
		return nil, fmt.Errorf("用户不存在")
	}

	if len(sr.Entries) > 1 {
		return nil, fmt.Errorf("找到多个用户，请联系管理员")
	}

	// 5. 获取用户DN
	userDN := sr.Entries[0].DN

	// 6. 提取用户属性（在关闭连接前）
	entry := sr.Entries[0]
	userAttrs := &LDAPUserAttributes{
		Username: username,
		CN:       entry.GetAttributeValue("cn"),
		Mail:     entry.GetAttributeValue("mail"),
		Mobile:   entry.GetAttributeValue("mobile"),
		Attrs:    make(map[string]string),
	}

	// 如果mobile为空，尝试获取telephoneNumber
	if userAttrs.Mobile == "" {
		userAttrs.Mobile = entry.GetAttributeValue("telephoneNumber")
	}

	// 如果CN为空，使用username
	if userAttrs.CN == "" {
		userAttrs.CN = username
	}

	// 保存所有属性
	for _, attr := range entry.Attributes {
		if len(attr.Values) > 0 {
			userAttrs.Attrs[attr.Name] = attr.Values[0]
		}
	}

	// 7. 关闭管理员连接
	l.Close()

	// 8. 重新建立新连接用于用户认证（避免在同一连接上多次Bind）
	var l2 *ldap.Conn
	if config.TLS {
		l2, err = ldap.DialTLS("tcp", ldapURL, &tls.Config{InsecureSkipVerify: true})
	} else {
		l2, err = ldap.DialURL(fmt.Sprintf("ldap://%s", ldapURL))
	}
	if err != nil {
		return nil, fmt.Errorf("LDAP重连失败: %v", err)
	}
	defer l2.Close()

	l2.SetTimeout(10 * time.Second)

	if config.StartTLS && !config.TLS {
		err = l2.StartTLS(&tls.Config{InsecureSkipVerify: true})
		if err != nil {
			return nil, fmt.Errorf("StartTLS失败: %v", err)
		}
	}

	// 9. 使用用户密码进行认证
	err = l2.Bind(userDN, password)
	if err != nil {
		return nil, fmt.Errorf("用户名或密码错误")
	}

	return userAttrs, nil
}

// TestLDAPConnection 测试LDAP连接
func TestLDAPConnection(config *LDAPAuthConfig) error {
	// 清理Host中的空格
	host := strings.TrimSpace(config.Host)
	ldapURL := fmt.Sprintf("%s:%d", host, config.Port)

	var l *ldap.Conn
	var err error

	// 使用带超时的Dial
	if config.TLS {
		tlsConfig := &tls.Config{InsecureSkipVerify: true}
		l, err = ldap.DialTLS("tcp", ldapURL, tlsConfig)
	} else {
		// 使用DialURL代替Dial，支持更好的超时控制
		l, err = ldap.DialURL(fmt.Sprintf("ldap://%s", ldapURL))
	}

	if err != nil {
		return fmt.Errorf("LDAP连接失败: %v", err)
	}
	defer l.Close()

	// 设置后续操作超时
	l.SetTimeout(10 * time.Second)

	if config.StartTLS && !config.TLS {
		err = l.StartTLS(&tls.Config{InsecureSkipVerify: true})
		if err != nil {
			return fmt.Errorf("StartTLS失败: %v", err)
		}
	}

	// 测试绑定
	err = l.Bind(config.BindUser, config.BindPass)
	if err != nil {
		return fmt.Errorf("LDAP绑定失败，请检查用户名和密码: %v", err)
	}

	return nil
}
