package email

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"time"
)

// Sender 邮件发送器
type Sender struct {
	config *Config
}

// NewSender 创建邮件发送器
func NewSender(config *Config) *Sender {
	return &Sender{
		config: config,
	}
}

// SendEmail 发送邮件
func (s *Sender) SendEmail(ctx context.Context, data *EmailData) error {
	// 设置超时
	ctx, cancel := context.WithTimeout(ctx, s.config.SMTP.Timeout)
	defer cancel()

	// 构建邮件内容
	message := s.buildMessage(data)

	// 配置SMTP认证
	auth := smtp.PlainAuth("", s.config.SMTP.Username, s.config.SMTP.Password, s.config.SMTP.Host)

	// 构建SMTP地址
	addr := fmt.Sprintf("%s:%d", s.config.SMTP.Host, s.config.SMTP.Port)

	// 发送邮件
	err := s.sendWithTLS(addr, auth, s.config.SMTP.From, []string{data.To}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// buildMessage 构建邮件消息
func (s *Sender) buildMessage(data *EmailData) string {
	message := fmt.Sprintf("From: %s\r\n", s.config.SMTP.From)
	message += fmt.Sprintf("To: %s\r\n", data.To)
	message += fmt.Sprintf("Subject: %s\r\n", data.Subject)
	message += "MIME-Version: 1.0\r\n"
	message += "Content-Type: text/html; charset=UTF-8\r\n"
	message += "\r\n"
	message += data.Body

	return message
}

// sendWithTLS 使用TLS发送邮件
func (s *Sender) sendWithTLS(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
	// 连接到SMTP服务器
	conn, err := tls.Dial("tcp", addr, &tls.Config{
		ServerName: s.config.SMTP.Host,
	})
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer conn.Close()

	// 创建SMTP客户端
	client, err := smtp.NewClient(conn, s.config.SMTP.Host)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Quit()

	// 认证
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP authentication failed: %w", err)
	}

	// 设置发件人
	if err = client.Mail(from); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	// 设置收件人
	for _, recipient := range to {
		if err = client.Rcpt(recipient); err != nil {
			return fmt.Errorf("failed to set recipient %s: %w", recipient, err)
		}
	}

	// 发送邮件内容
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}

	_, err = writer.Write(msg)
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close data writer: %w", err)
	}

	return nil
}

// SendTenantActivationEmail 发送租户激活邮件
func (s *Sender) SendTenantActivationEmail(ctx context.Context, to, userName, tenantName, activationLink, expireTime string) error {
	tm := NewTemplateManager()

	data := map[string]interface{}{
		"UserName":       userName,
		"TenantName":     tenantName,
		"ActivationLink": activationLink,
		"ExpireTime":     expireTime,
		"CurrentYear":    time.Now().Year(),
	}

	subject, body, err := tm.RenderTemplate(EmailTypeTenantActivation, data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	emailData := &EmailData{
		To:      to,
		Subject: subject,
		Body:    body,
	}

	return s.SendEmail(ctx, emailData)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// SendInvitationEmail 发送邀请邮件
func (s *Sender) SendInvitationEmail(ctx context.Context, to, userName, tenantName, departmentName, roleName, inviterName, inviteTime, acceptLink, declineLink, expireTime string) error {
	tm := NewTemplateManager()

	data := map[string]interface{}{
		"UserName":       userName,
		"TenantName":     tenantName,
		"DepartmentName": departmentName,
		"RoleName":       roleName,
		"InviterName":    inviterName,
		"InviteTime":     inviteTime,
		"AcceptLink":     acceptLink,
		"DeclineLink":    declineLink,
		"ExpireTime":     expireTime,
		"CurrentYear":    time.Now().Year(),
	}

	subject, body, err := tm.RenderTemplate(EmailTypeInvitation, data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	emailData := &EmailData{
		To:      to,
		Subject: subject,
		Body:    body,
	}

	return s.SendEmail(ctx, emailData)
}

// SendPasswordResetEmail 发送密码重置邮件
func (s *Sender) SendPasswordResetEmail(ctx context.Context, to, userName, resetLink, expireTime string) error {
	tm := NewTemplateManager()

	data := map[string]interface{}{
		"UserName":    userName,
		"ResetLink":   resetLink,
		"ExpireTime":  expireTime,
		"CurrentYear": time.Now().Year(),
	}

	subject, body, err := tm.RenderTemplate(EmailTypePasswordReset, data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	emailData := &EmailData{
		To:      to,
		Subject: subject,
		Body:    body,
	}

	return s.SendEmail(ctx, emailData)
}
