package email

import (
	"time"
)

// Config 邮件配置
type Config struct {
	SMTP SMTPConfig `yaml:"smtp"`
}

// SMTPConfig SMTP配置
type SMTPConfig struct {
	Host     string        `yaml:"host"`     // SMTP服务器地址
	Port     int           `yaml:"port"`     // SMTP端口
	Username string        `yaml:"username"` // 用户名
	Password string        `yaml:"password"` // 密码
	From     string        `yaml:"from"`     // 发件人邮箱
	Timeout  time.Duration `yaml:"timeout"` // 超时时间
}

// EmailTemplate 邮件模板
type EmailTemplate struct {
	Subject string            `yaml:"subject"` // 邮件主题
	Body    string            `yaml:"body"`    // 邮件正文
	Params  map[string]string `yaml:"params"` // 模板参数
}

// EmailData 邮件数据
type EmailData struct {
	To      string            `json:"to"`      // 收件人
	Subject string            `json:"subject"` // 主题
	Body    string            `json:"body"`    // 正文
	Params  map[string]string `json:"params"` // 参数
}

// EmailType 邮件类型
type EmailType string

const (
	EmailTypeTenantActivation EmailType = "tenant_activation" // 租户激活邮件
	EmailTypeInvitation       EmailType = "invitation"        // 邀请加入邮件
	EmailTypePasswordReset    EmailType = "password_reset"    // 密码重置邮件
)
