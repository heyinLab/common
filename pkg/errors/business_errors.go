package errors

import (
	commonV1 "go-heyin/api/gen/go/common"
	"strings"
)

// 业务错误类型
type BusinessError struct {
	Code     int32  `json:"code"`      // 业务错误码，使用生成的枚举
	Message  string `json:"message"`   // 错误消息
	Type     string `json:"type"`      // 错误类型
	HttpCode int32  `json:"http_code"` // 对应的HTTP状态码
}

func (e *BusinessError) Error() string {
	return e.Message
}

// 预定义的业务错误
var (
	// 用户相关错误 (10001-10099)
	ErrUserNotFound      = &BusinessError{Code: convertToInt32(commonV1.ErrorCode_USER_NOT_FOUND), Message: "用户不存在", Type: "USER_NOT_FOUND", HttpCode: 404}
	ErrUserAlreadyExists = &BusinessError{Code: convertToInt32(commonV1.ErrorCode_USER_ALREADY_EXISTS), Message: "用户已存在", Type: "USER_ALREADY_EXISTS", HttpCode: 409}
	ErrInvalidPassword   = &BusinessError{Code: convertToInt32(commonV1.ErrorCode_INVALID_PASSWORD), Message: "密码格式不正确", Type: "INVALID_PASSWORD", HttpCode: 400}
	ErrUserDisabled      = &BusinessError{Code: convertToInt32(commonV1.ErrorCode_USER_DISABLED), Message: "用户已被禁用", Type: "USER_DISABLED", HttpCode: 403}
	ErrUserDeleted       = &BusinessError{Code: convertToInt32(commonV1.ErrorCode_USER_DELETED), Message: "用户已被删除", Type: "USER_DELETED", HttpCode: 404}

	// 租户相关错误 (10100-10199)
	ErrTenantNotFound      = &BusinessError{Code: convertToInt32(commonV1.ErrorCode_TENANT_NOT_FOUND), Message: "租户不存在", Type: "TENANT_NOT_FOUND", HttpCode: 404}
	ErrTenantAlreadyExists = &BusinessError{Code: convertToInt32(commonV1.ErrorCode_TENANT_ALREADY_EXISTS), Message: "租户已存在", Type: "TENANT_ALREADY_EXISTS", HttpCode: 409}
	ErrTenantDisabled      = &BusinessError{Code: convertToInt32(commonV1.ErrorCode_TENANT_DISABLED), Message: "租户已被禁用", Type: "TENANT_DISABLED", HttpCode: 403}
	ErrTenantPending       = &BusinessError{Code: convertToInt32(commonV1.ErrorCode_TENANT_PENDING), Message: "租户待审核", Type: "TENANT_PENDING", HttpCode: 403}
	ErrTenantRejected      = &BusinessError{Code: convertToInt32(commonV1.ErrorCode_TENANT_REJECTED), Message: "租户申请被拒绝", Type: "TENANT_REJECTED", HttpCode: 403}

	// 权限相关错误 (10200-10299)
	ErrPermissionDenied   = &BusinessError{Code: convertToInt32(commonV1.ErrorCode_PERMISSION_DENIED), Message: "权限不足", Type: "PERMISSION_DENIED", HttpCode: 403}
	ErrRoleNotFound       = &BusinessError{Code: convertToInt32(commonV1.ErrorCode_ROLE_NOT_FOUND), Message: "角色不存在", Type: "ROLE_NOT_FOUND", HttpCode: 404}
	ErrRoleDisabled       = &BusinessError{Code: convertToInt32(commonV1.ErrorCode_ROLE_DISABLED), Message: "角色已被禁用", Type: "ROLE_DISABLED", HttpCode: 403}
	ErrPermissionNotFound = &BusinessError{Code: convertToInt32(commonV1.ErrorCode_PERMISSION_NOT_FOUND), Message: "权限不存在", Type: "PERMISSION_NOT_FOUND", HttpCode: 404}

	// 认证相关错误 (10300-10399)
	ErrInvalidCredentials = &BusinessError{Code: convertToInt32(commonV1.ErrorCode_INVALID_CREDENTIALS), Message: "用户名或密码错误", Type: "INVALID_CREDENTIALS", HttpCode: 401}
	ErrTokenExpired       = &BusinessError{Code: convertToInt32(commonV1.ErrorCode_TOKEN_EXPIRED), Message: "Token已过期", Type: "TOKEN_EXPIRED", HttpCode: 401}
	ErrTokenInvalid       = &BusinessError{Code: convertToInt32(commonV1.ErrorCode_TOKEN_INVALID), Message: "Token无效", Type: "TOKEN_INVALID", HttpCode: 401}
	ErrTokenRevoked       = &BusinessError{Code: convertToInt32(commonV1.ErrorCode_TOKEN_REVOKED), Message: "Token已被撤销", Type: "TOKEN_REVOKED", HttpCode: 401}
	ErrAccountLocked      = &BusinessError{Code: convertToInt32(commonV1.ErrorCode_ACCOUNT_LOCKED), Message: "账户已被锁定", Type: "ACCOUNT_LOCKED", HttpCode: 403}
	ErrAuthHeaderMissing  = &BusinessError{Code: convertToInt32(commonV1.ErrorCode_AUTH_HEADER_MISSING), Message: "缺少Authorization头", Type: "AUTH_HEADER_MISSING", HttpCode: 401}
	ErrAuthHeaderInvalid  = &BusinessError{Code: convertToInt32(commonV1.ErrorCode_AUTH_HEADER_INVALID), Message: "Authorization头格式错误", Type: "AUTH_HEADER_INVALID", HttpCode: 401}
	ErrAuthServiceError   = &BusinessError{Code: convertToInt32(commonV1.ErrorCode_AUTH_SERVICE_ERROR), Message: "认证服务错误", Type: "AUTH_SERVICE_ERROR", HttpCode: 500}
	ErrUserTypeUndefined  = &BusinessError{Code: convertToInt32(commonV1.ErrorCode_USER_TYPE_UNDEFINED), Message: "用户类型未定义", Type: "USER_TYPE_UNDEFINED", HttpCode: 401}
	ErrAccessForbidden    = &BusinessError{Code: convertToInt32(commonV1.ErrorCode_ACCESS_FORBIDDEN), Message: "访问被禁止", Type: "ACCESS_FORBIDDEN", HttpCode: 403}
	ErrTenantMissing      = &BusinessError{Code: convertToInt32(commonV1.ErrorCode_TENANT_MISSING), Message: "缺少租户ID", Type: "TENANT_MISSING", HttpCode: 400}
	ErrTenantInvalid      = &BusinessError{Code: convertToInt32(commonV1.ErrorCode_TENANT_INVALID), Message: "租户ID格式错误", Type: "TENANT_INVALID", HttpCode: 400}
	ErrRegisterFailed     = &BusinessError{Code: convertToInt32(commonV1.ErrorCode_REGISTER_FAILED), Message: "注册失败", Type: "REGISTER_FAILED", HttpCode: 400}
	// 参数验证错误 (10400-10499)
	ErrInvalidParameter = &BusinessError{Code: convertToInt32(commonV1.ErrorCode_INVALID_PARAMETER), Message: "参数错误", Type: "INVALID_PARAMETER", HttpCode: 400}
	ErrMissingParameter = &BusinessError{Code: convertToInt32(commonV1.ErrorCode_MISSING_PARAMETER), Message: "缺少必要参数", Type: "MISSING_PARAMETER", HttpCode: 400}
	ErrInvalidFormat    = &BusinessError{Code: convertToInt32(commonV1.ErrorCode_INVALID_FORMAT), Message: "数据格式错误", Type: "INVALID_FORMAT", HttpCode: 400}
	ErrInvalidEmail     = &BusinessError{Code: convertToInt32(commonV1.ErrorCode_INVALID_EMAIL), Message: "邮箱格式错误", Type: "INVALID_EMAIL", HttpCode: 400}
	ErrInvalidPhone     = &BusinessError{Code: convertToInt32(commonV1.ErrorCode_INVALID_PHONE), Message: "手机号格式错误", Type: "INVALID_PHONE", HttpCode: 400}

	// 数据相关错误 (10500-10599)
	ErrDataNotFound   = &BusinessError{Code: convertToInt32(commonV1.ErrorCode_DATA_NOT_FOUND), Message: "数据不存在", Type: "DATA_NOT_FOUND", HttpCode: 404}
	ErrDataConflict   = &BusinessError{Code: convertToInt32(commonV1.ErrorCode_DATA_CONFLICT), Message: "数据冲突", Type: "DATA_CONFLICT", HttpCode: 409}
	ErrDataInvalid    = &BusinessError{Code: convertToInt32(commonV1.ErrorCode_DATA_INVALID), Message: "数据无效", Type: "DATA_INVALID", HttpCode: 400}
	ErrDataDuplicate  = &BusinessError{Code: convertToInt32(commonV1.ErrorCode_DATA_DUPLICATE), Message: "数据重复", Type: "DATA_DUPLICATE", HttpCode: 409}
	ErrDataConstraint = &BusinessError{Code: convertToInt32(commonV1.ErrorCode_DATA_CONSTRAINT), Message: "数据约束错误", Type: "DATA_CONSTRAINT", HttpCode: 400}

	// 系统相关错误 (19900-19999)
	ErrSystemError        = &BusinessError{Code: convertToInt32(commonV1.ErrorCode_SYSTEM_ERROR), Message: "系统错误", Type: "SYSTEM_ERROR", HttpCode: 500}
	ErrServiceUnavailable = &BusinessError{Code: convertToInt32(commonV1.ErrorCode_SERVICE_UNAVAILABLE), Message: "服务不可用", Type: "SERVICE_UNAVAILABLE", HttpCode: 503}
	ErrDatabaseError      = &BusinessError{Code: convertToInt32(commonV1.ErrorCode_DATABASE_ERROR), Message: "数据库错误", Type: "DATABASE_ERROR", HttpCode: 500}
	ErrNetworkError       = &BusinessError{Code: convertToInt32(commonV1.ErrorCode_NETWORK_ERROR), Message: "网络错误", Type: "NETWORK_ERROR", HttpCode: 500}
)

// 错误分类函数
func ClassifyError(err error) *BusinessError {
	if err == nil {
		return nil
	}

	// 检查是否已经是业务错误
	if businessErr, ok := err.(*BusinessError); ok {
		return businessErr
	}

	errMsg := strings.ToLower(err.Error())

	// 根据错误消息分类
	switch {
	// 数据库相关错误
	case strings.Contains(errMsg, "duplicate") || strings.Contains(errMsg, "unique constraint"):
		return &BusinessError{Code: convertToInt32(commonV1.ErrorCode_DATA_DUPLICATE), Message: "数据重复", Type: "DATA_DUPLICATE", HttpCode: 409}

	case strings.Contains(errMsg, "not found") || strings.Contains(errMsg, "no rows"):
		return &BusinessError{Code: convertToInt32(commonV1.ErrorCode_DATA_NOT_FOUND), Message: "数据不存在", Type: "DATA_NOT_FOUND", HttpCode: 404}

	case strings.Contains(errMsg, "foreign key") || strings.Contains(errMsg, "constraint"):
		return &BusinessError{Code: convertToInt32(commonV1.ErrorCode_DATA_CONSTRAINT), Message: "数据约束错误", Type: "DATA_CONSTRAINT", HttpCode: 400}

	case strings.Contains(errMsg, "permission") || strings.Contains(errMsg, "access denied"):
		return &BusinessError{Code: convertToInt32(commonV1.ErrorCode_PERMISSION_DENIED), Message: "权限不足", Type: "PERMISSION_DENIED", HttpCode: 403}

	case strings.Contains(errMsg, "invalid") || strings.Contains(errMsg, "malformed"):
		return &BusinessError{Code: convertToInt32(commonV1.ErrorCode_INVALID_PARAMETER), Message: "参数无效", Type: "INVALID_PARAMETER", HttpCode: 400}

	case strings.Contains(errMsg, "timeout") || strings.Contains(errMsg, "connection"):
		return &BusinessError{Code: convertToInt32(commonV1.ErrorCode_SERVICE_UNAVAILABLE), Message: "服务不可用", Type: "SERVICE_UNAVAILABLE", HttpCode: 503}

	// 默认业务错误
	default:
		return &BusinessError{Code: convertToInt32(commonV1.ErrorCode_SYSTEM_ERROR), Message: "系统错误", Type: "SYSTEM_ERROR", HttpCode: 500}
	}
}

// 创建业务错误
func NewBusinessError(code int32, message, errorType string, httpCode int32) *BusinessError {
	return &BusinessError{
		Code:     code,
		Message:  message,
		Type:     errorType,
		HttpCode: httpCode,
	}
}

// 包装错误
func WrapError(err error, message string) *BusinessError {
	if businessErr, ok := err.(*BusinessError); ok {
		return &BusinessError{
			Code:     businessErr.Code,
			Message:  message + ": " + businessErr.Message,
			Type:     businessErr.Type,
			HttpCode: businessErr.HttpCode,
		}
	}
	return &BusinessError{
		Code:     convertToInt32(commonV1.ErrorCode_SYSTEM_ERROR),
		Message:  message + ": " + err.Error(),
		Type:     "WRAPPED_ERROR",
		HttpCode: 500,
	}
}

func convertToInt32(error commonV1.ErrorCode) int32 {
	return int32(error)
}

// 获取HTTP状态码
func (e *BusinessError) GetHttpCode() int32 {
	return e.HttpCode
}

// 判断是否为系统错误
func (e *BusinessError) IsSystemError() bool {
	return int32(e.Code) >= 19900
}

// 判断是否为业务错误
func (e *BusinessError) IsBusinessError() bool {
	return int32(e.Code) < 19900
}
