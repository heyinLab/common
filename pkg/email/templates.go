package email

import (
	"fmt"
	"html/template"
	"strings"
)

// TemplateManager 模板管理器
type TemplateManager struct {
	templates map[EmailType]*template.Template
}

// NewTemplateManager 创建模板管理器
func NewTemplateManager() *TemplateManager {
	tm := &TemplateManager{
		templates: make(map[EmailType]*template.Template),
	}
	tm.initTemplates()
	return tm
}

// initTemplates 初始化模板
func (tm *TemplateManager) initTemplates() {
	// 租户激活邮件模板
	tm.templates[EmailTypeTenantActivation] = template.Must(template.New("tenant_activation").Parse(tenantActivationTemplate))

	// 邀请加入邮件模板
	tm.templates[EmailTypeInvitation] = template.Must(template.New("invitation").Parse(invitationTemplate))

	// 密码重置邮件模板
	tm.templates[EmailTypePasswordReset] = template.Must(template.New("password_reset").Parse(passwordResetTemplate))

	// 验证模板是否正确解析
	for emailType, t := range tm.templates {
		if t.Lookup("subject") == nil {
			panic(fmt.Sprintf("subject template not found for %s", emailType))
		}
		if t.Lookup("body") == nil {
			panic(fmt.Sprintf("body template not found for %s", emailType))
		}
	}
}

// RenderTemplate 渲染模板
func (tm *TemplateManager) RenderTemplate(emailType EmailType, data map[string]interface{}) (string, string, error) {
	t, exists := tm.templates[emailType]
	if !exists {
		return "", "", fmt.Errorf("template not found for type: %s", emailType)
	}

	// 渲染主题
	var subjectBuilder strings.Builder
	subjectTemplate := t.Lookup("subject")
	if subjectTemplate == nil {
		return "", "", fmt.Errorf("subject template not found")
	}
	err := subjectTemplate.Execute(&subjectBuilder, data)
	if err != nil {
		return "", "", fmt.Errorf("failed to render subject: %w", err)
	}

	// 渲染正文
	var bodyBuilder strings.Builder
	bodyTemplate := t.Lookup("body")
	if bodyTemplate == nil {
		return "", "", fmt.Errorf("body template not found")
	}
	err = bodyTemplate.Execute(&bodyBuilder, data)
	if err != nil {
		return "", "", fmt.Errorf("failed to render body: %w", err)
	}

	return subjectBuilder.String(), bodyBuilder.String(), nil
}

// 这是一个 Go 代码文件，包含三个优化后的邮件模板常量。
// 这些模板具有更好的邮件客户端兼容性。

// 1. 租户激活邮件模板 (优化版)
const tenantActivationTemplate = `{{define "subject"}}欢迎加入 {{.TenantName}} - 请激活您的账户{{end}}
{{define "body"}}
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>账户激活</title>
    <style>
        body { 
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Helvetica Neue', Arial, sans-serif; 
            line-height: 1.6; 
            color: #333333; 
            font-size: 16px;
            margin: 0;
            padding: 0;
            background-color: #f4f4f7; /* 浅灰色背景 */
        }
        .container { 
            max-width: 600px; 
            margin: 20px auto; 
            padding: 0; 
            background-color: #ffffff; /* 白色卡片 */
            border: 1px solid #e0e0e0;
            border-radius: 8px;
            overflow: hidden; 
        }
        .header { 
            background-color: #ffffff; 
            padding: 30px 20px; 
            text-align: center; 
            border-bottom: 1px solid #e0e0e0;
        }
        .header h1 { margin: 0; color: #222222; font-size: 24px; }
        .content { background: #ffffff; padding: 32px; }
        .content p, .content ul { margin-bottom: 20px; }
        .footer { 
            background: #f9f9f9; 
            padding: 20px; 
            text-align: center; 
            font-size: 13px; 
            color: #777777; 
        }
        
        /* --- 基础按钮样式 (重要) --- */
        .button-base {
            display: inline-block; 
            padding: 14px 28px; 
            text-decoration: none !important; /* 强制无下划线 */
            border-radius: 8px; 
            margin: 20px 0; 
            font-size: 16px; 
            font-weight: 600; 
            text-align: center; 
            border: none;
            cursor: pointer;
            color: #ffffff !important; /* 强制白色文字 */
        }
        .button-primary { 
            background-color: #007bff; /* 纯蓝色 */
        }
        
        /* --- 辅助样式 --- */
        .highlight { color: #007bff; font-weight: bold; }
        .link-box {
            word-break: break-all; 
            background: #f8f9fa; 
            padding: 12px; 
            border-radius: 4px;
            font-family: 'Courier New', Courier, monospace;
        }
        .text-center { text-align: center; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>欢迎加入 {{.TenantName}}</h1>
        </div>
        <div class="content">
            <h2>亲爱的 {{.UserName}}，</h2>
            <p>欢迎加入 <span class="highlight">{{.TenantName}}</span>！您的账户已成功创建。</p>
            
            <p>请点击下面的按钮激活您的账户：</p>
            <div class="text-center">
             	<a href="{{.ActivationLink}}" class="button-base button-primary">激活账户</a>

            </div>
            
            <p>如果按钮无法点击，请复制以下链接到浏览器中打开：</p>
            <p class="link-box">{{.ActivationLink}}</p>
            
            <p><strong>注意事项：</strong></p>
            <ul>
                <li>此激活链接将在 {{.ExpireTime}} 后过期</li>
                <li>如果链接已过期，请联系管理员重新发送激活邮件</li>
                <li>请妥善保管您的登录凭据</li>
            </ul>
            
            <p>如有任何问题，请联系我们的技术支持团队。</p>
        </div>
        <div class="footer">
            <p>此邮件由系统自动发送，请勿回复。</p>
            <p>&copy; {{.CurrentYear}} {{.TenantName}}. 保留所有权利。</p>
        </div>
    </div>
</body>
</html>
{{end}}
`

// 2. 邀请加入邮件模板 (优化版)
const invitationTemplate = `
{{define "subject"}}邀请您加入 {{.TenantName}} 的 {{.DepartmentName}} 部门{{end}}
{{define "body"}}
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>部门邀请</title>
    <style>
        body { 
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Helvetica Neue', Arial, sans-serif; 
            line-height: 1.6; color: #333333; font-size: 16px;
            margin: 0; padding: 0; background-color: #f4f4f7;
        }
        .container { 
            max-width: 600px; margin: 20px auto; padding: 0; 
            background-color: #ffffff; border: 1px solid #e0e0e0;
            border-radius: 8px; overflow: hidden; 
        }
        .header { 
            background-color: #ffffff; padding: 30px 20px; 
            text-align: center; border-bottom: 1px solid #e0e0e0;
        }
        .header h1 { margin: 0; color: #222222; font-size: 24px; }
        .content { background: #ffffff; padding: 32px; }
        .content p, .content ul { margin-bottom: 20px; }
        .footer { 
            background: #f9f9f9; padding: 20px; text-align: center; 
            font-size: 13px; color: #777777; 
        }
        
        /* --- 基础按钮样式 (重要) --- */
        .button-base {
            display: inline-block; 
            padding: 14px 28px; 
            text-decoration: none !important; 
            border-radius: 8px; 
            margin: 10px 8px; /* 调整间距 */
            font-size: 16px; 
            font-weight: 600; 
            text-align: center; 
            border: none;
            cursor: pointer;
            color: #ffffff !important; 
        }
        .button-success { 
            background-color: #28a745; /* 纯绿色 */
        }
        .button-secondary { 
            background-color: #6c757d; /* 纯灰色 */
        }
        
        /* --- 辅助样式 --- */
        .highlight { color: #007bff; font-weight: bold; }
        .link-box {
            word-break: break-all; background: #f8f9fa; 
            padding: 12px; border-radius: 4px;
            font-family: 'Courier New', Courier, monospace;
        }
        .role-info { 
            background: #f8f9fa; padding: 15px; 
            border-radius: 4px; margin: 15px 0; 
        }
        .text-center { text-align: center; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>邀请</h1>
        </div>
        <div class="content">
            <h2>亲爱的 {{.UserName}}，</h2>
            <p><span class="highlight">{{.InviterName}}</span> 邀请您加入 <span class="highlight">{{.TenantName}}</span> 的 <span class="highlight">{{.DepartmentName}}</span> 部门。</p>
            
            <div class="role-info">
                <h3>邀请详情：</h3>
                <p><strong>组织：</strong>{{.TenantName}}</p>
                <p><strong>部门：</strong>{{.DepartmentName}}</p>
                <p><strong>角色：</strong>{{.RoleName}}</p>
                <p><strong>邀请人：</strong>{{.InviterName}}</p>
                <p><strong>邀请时间：</strong>{{.InviteTime}}</p>
            </div>
            
            <div class="text-center">
                <a href="{{.AcceptLink}}" class="button-base button-success">接受邀请</a>
            </div>
            
            <p>如果按钮无法点击，请复制以下链接到浏览器中打开：</p>
            <p><strong>接受邀请：</strong></p>
            <p class="link-box">{{.AcceptLink}}</p>
            
            <p><strong>注意事项：</strong></p>
            <ul>
                <li>此邀请将在 {{.ExpireTime}} 后过期</li>
                <li>接受邀请后，您将获得相应的部门权限</li>
                <li>如有疑问，请联系邀请人或技术支持团队</li>
            </ul>
        </div>
        <div class="footer">
            <p>此邮件由系统自动发送，请勿回复。</p>
            <p>&copy; {{.CurrentYear}} {{.TenantName}}. 保留所有权利。</p>
        </div>
    </div>
</body>
</html>
{{end}}
`

// 3. 密码重置邮件模板 (优化版)
const passwordResetTemplate = `
{{define "subject"}}密码重置请求 - {{.TenantName}}{{end}}
{{define "body"}}
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>密码重置</title>
    <style>
        body { 
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Helvetica Neue', Arial, sans-serif; 
            line-height: 1.6; color: #333333; font-size: 16px;
            margin: 0; padding: 0; background-color: #f4f4f7;
        }
        .container { 
            max-width: 600px; margin: 20px auto; padding: 0; 
            background-color: #ffffff; border: 1px solid #e0e0e0;
            border-radius: 8px; overflow: hidden; 
        }
        .header { 
            background-color: #ffffff; padding: 30px 20px; 
            text-align: center; border-bottom: 1px solid #e0e0e0;
        }
        .header h1 { margin: 0; color: #222222; font-size: 24px; }
        .content { background: #ffffff; padding: 32px; }
        .content p, .content ul { margin-bottom: 20px; }
        .footer { 
            background: #f9f9f9; padding: 20px; text-align: center; 
            font-size: 13px; color: #777777; 
        }
        
        /* --- 基础按钮样式 (重要) --- */
        .button-base {
            display: inline-block; 
            padding: 14px 28px; 
            text-decoration: none !important; 
            border-radius: 8px; 
            margin: 20px 0; 
            font-size: 16px; 
            font-weight: 600; 
            text-align: center; 
            border: none;
            cursor: pointer;
            color: #ffffff !important; 
        }
        .button-danger { 
            background-color: #dc3545; /* 纯红色 */
        }
        
        /* --- 辅助样式 --- */
        .highlight { color: #dc3545; font-weight: bold; }
        .link-box {
            word-break: break-all; background: #f8f9fa; 
            padding: 12px; border-radius: 4px;
            font-family: 'Courier New', Courier, monospace;
        }
        .warning { 
            background: #fff3cd; 
            border: 1px solid #ffeeba; 
            padding: 15px; 
            border-radius: 4px; 
            margin: 15px 0; 
            color: #856404; /* 确保文字可读 */
        }
        .text-center { text-align: center; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>密码重置请求</h1>
        </div>
        <div class="content">
            <h2>亲爱的 {{.UserName}}，</h2>
            <p>我们收到了您对该账户的密码重置请求。</p>
            
            <div class="warning">
                <h3>⚠️ 安全提醒</h3>
                <p>如果您没有请求密码重置，请忽略此邮件。您的账户仍然是安全的。</p>
            </div>
            
            <p>要重置您的密码，请点击下面的按钮：</p>
            <div class="text-center">
                <a href="{{.ResetLink}}" class="button-base button-danger">重置密码</a>
            </div>
            
            <p>如果按钮无法点击，请复制以下链接到浏览器中打开：</p>
            <p class="link-box">{{.ResetLink}}</p>
            
            <p><strong>重要信息：</strong></p>
            <ul>
                <li>此重置链接将在 {{.ExpireTime}} 后过期</li>
                <li>链接只能使用一次，使用后立即失效</li>
                <li>为了账户安全，请设置一个强密码</li>
            </ul>
        </div>
        <div class="footer">
            <p>此邮件由系统自动发送，请勿回复。</p>
            <p>&copy; {{.CurrentYear}} {{.TenantName}}. 保留所有权利。</p>
        </div>
    </div>
</body>
</html>
{{end}}
`
