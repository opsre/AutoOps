package util

import "time"

// SSHConfig SSH连接配置
type SSHConfig struct {
	IP                 string        // 主机IP
	Port               int           // SSH端口
	Type               int           // 认证类型:1->密码,2->私钥,3->公钥(免认证),4->跳板机
	Username           string        // 用户名
	Password           string        // 密码(type=1时使用)
	PublicKey          string        // 私钥内容(type=2时使用) - 注意：字段名为PublicKey但实际存储私钥
	Timeout            time.Duration // 超时时间
	BastionHost        string        // 跳板机IP (type=4时使用)
	BastionPort        int           // 跳板机端口 (type=4时使用)
	BastionUsername    string        // 跳板机用户名 (type=4时使用)
	BastionAuthType    int           // 跳板机认证类型 (type=4时使用)
	BastionPassword    string        // 跳板机密码 (type=4时使用)
	BastionPrivateKey  string        // 跳板机私钥 (type=4时使用)
	TargetAuthType     int           // 目标机器认证类型 (type=4时使用)
	TargetPassword     string        // 目标机器密码 (type=4时使用)
	TargetPrivateKey   string        // 目标机器私钥 (type=4时使用)
}
