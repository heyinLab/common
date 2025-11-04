package user_type

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	businessErrors "go-heyin/pkg/errors"
	"go-heyin/pkg/middleware/auth"
)

const (
	UserTypePlatform = "platform"
	UserTypeMerchant = "merchant"
)

// "/api/v1/auth",
// "/api/v1/permissions",
// "/api/v1/tenants",
// "/api/v1/users",
// "/api/v1/departments",
// "/api/v1/invitations",
// "/api/v1/roles",
// "/api/v1/members",
// "/api/v1/groups"
var (
	platformNotAllowedPaths = []string{
		"/api/v1/tenants",
	}
	merchantNotAllowedPaths = []string{
		"/api/v1/groups",
	}
)

// Server 用户类型中间件
func Server() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			// 从JWT token中获取用户类型
			claims, ok := auth.FromContext(ctx)
			if !ok {
				return nil, errors.New(int(businessErrors.ErrTokenInvalid.HttpCode), businessErrors.ErrTokenInvalid.Type, businessErrors.ErrTokenInvalid.Message)
			}

			userType := claims.UserType
			if userType == "" {
				return nil, errors.New(int(businessErrors.ErrUserTypeUndefined.HttpCode), businessErrors.ErrUserTypeUndefined.Type, businessErrors.ErrUserTypeUndefined.Message)
			}

			// 获取当前请求的路径
			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				return nil, errors.New(int(businessErrors.ErrSystemError.HttpCode), businessErrors.ErrSystemError.Type, businessErrors.ErrSystemError.Message)
			}

			httpTr, ok := tr.(http.Transporter)
			if !ok {
				return nil, errors.New(int(businessErrors.ErrSystemError.HttpCode), businessErrors.ErrSystemError.Type, businessErrors.ErrSystemError.Message)
			}

			path := httpTr.Request().URL.Path

			// 检查用户类型权限
			if !hasPermission(userType, path) {
				return nil, errors.New(int(businessErrors.ErrAccessForbidden.HttpCode), businessErrors.ErrAccessForbidden.Type, businessErrors.ErrAccessForbidden.Message)
			}

			return handler(ctx, req)
		}
	}
}

// hasPermission 检查用户类型是否有权限访问指定路径
func hasPermission(userType, path string) bool {
	switch userType {
	case UserTypePlatform:
		// 平台用户只能访问平台管理API
		return notAllowedPath(path, platformNotAllowedPaths)
	case UserTypeMerchant:
		// 商户用户只能访问商户相关API
		return notAllowedPath(path, merchantNotAllowedPaths)
	default:
		return false
	}
}

func notAllowedPath(path string, paths []string) bool {
	for _, excludePath := range paths {
		if len(path) >= len(excludePath) && path[:len(excludePath)] == excludePath {
			return false
		}
	}

	// 其他路径允许访问
	return true
}
