package email

import (
	"context"
	"fmt"
	"time"
)

// Service 邮件服务
type Service struct {
	sender *Sender
}

// NewService 创建邮件服务
func NewService(config *Config) Service {
	return Service{
		sender: NewSender(config),
	}
}

// SendTenantActivationEmail 发送租户激活邮件
func (s *Service) SendTenantActivationEmail(ctx context.Context, req *TenantActivationEmailRequest) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	if req.To == "" || req.UserName == "" || req.TenantName == "" || req.ActivationLink == "" {
		return fmt.Errorf("required fields cannot be empty")
	}

	// 设置默认过期时间（24小时）
	expireTime := "24小时"
	if req.ExpireTime != "" {
		expireTime = req.ExpireTime
	}

	return s.sender.SendTenantActivationEmail(
		ctx,
		req.To,
		req.UserName,
		req.TenantName,
		req.ActivationLink,
		expireTime,
	)
}

// SendInvitationEmail 发送邀请邮件
func (s *Service) SendInvitationEmail(ctx context.Context, req *InvitationEmailRequest) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	if req.To == "" || req.UserName == "" || req.TenantName == "" || req.DepartmentName == "" || req.AcceptLink == "" {
		return fmt.Errorf("required fields cannot be empty")
	}

	// 设置默认过期时间（7天）
	expireTime := "7天"
	if req.ExpireTime != "" {
		expireTime = req.ExpireTime
	}

	// 设置默认邀请时间
	inviteTime := time.Now().Format("2006-01-02 15:04:05")
	if req.InviteTime != "" {
		inviteTime = req.InviteTime
	}

	return s.sender.SendInvitationEmail(
		ctx,
		req.To,
		req.UserName,
		req.TenantName,
		req.DepartmentName,
		req.RoleName,
		req.InviterName,
		inviteTime,
		req.AcceptLink,
		req.DeclineLink,
		expireTime,
	)
}

// SendPasswordResetEmail 发送密码重置邮件
func (s *Service) SendPasswordResetEmail(ctx context.Context, req *PasswordResetEmailRequest) error {
	if req == nil {
		return fmt.Errorf("request cannot be nil")
	}

	if req.To == "" || req.UserName == "" || req.ResetLink == "" {
		return fmt.Errorf("required fields cannot be empty")
	}

	// 设置默认过期时间（1小时）
	expireTime := "1小时"
	if req.ExpireTime != "" {
		expireTime = req.ExpireTime
	}

	return s.sender.SendPasswordResetEmail(
		ctx,
		req.To,
		req.UserName,
		req.ResetLink,
		expireTime,
	)
}

// TenantActivationEmailRequest 租户激活邮件请求
type TenantActivationEmailRequest struct {
	To             string `json:"to"`              // 收件人邮箱
	UserName       string `json:"user_name"`       // 用户名
	TenantName     string `json:"tenant_name"`     // 租户名称
	ActivationLink string `json:"activation_link"` // 激活链接
	ExpireTime     string `json:"expire_time"`     // 过期时间（可选）
}

// InvitationEmailRequest 邀请邮件请求
type InvitationEmailRequest struct {
	To             string `json:"to"`              // 收件人邮箱
	UserName       string `json:"user_name"`       // 用户名
	TenantName     string `json:"tenant_name"`     // 租户名称
	DepartmentName string `json:"department_name"` // 部门名称
	RoleName       string `json:"role_name"`       // 角色名称
	InviterName    string `json:"inviter_name"`    // 邀请人姓名
	InviteTime     string `json:"invite_time"`     // 邀请时间（可选）
	AcceptLink     string `json:"accept_link"`     // 接受链接
	DeclineLink    string `json:"decline_link"`    // 拒绝链接
	ExpireTime     string `json:"expire_time"`     // 过期时间（可选）
}

// PasswordResetEmailRequest 密码重置邮件请求
type PasswordResetEmailRequest struct {
	To         string `json:"to"`          // 收件人邮箱
	UserName   string `json:"user_name"`   // 用户名
	ResetLink  string `json:"reset_link"`  // 重置链接
	ExpireTime string `json:"expire_time"` // 过期时间（可选）
}
