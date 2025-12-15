package auth

import (
	"context"
	"strconv"

	businessErrors "github.com/heyinLab/common/pkg/errors"
	"github.com/heyinLab/common/pkg/middleware/common"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

func Server(needTenant bool) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			// 从 context 中获取 transport 信息 (HTTP/gRPC)
			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				return nil, errors.New(int(businessErrors.ErrSystemError.HttpCode), businessErrors.ErrSystemError.Type, businessErrors.ErrSystemError.Message)
			}

			// 信任上游传递来的header X-User-ID X-User-Type  X-Tenant-ID X-Region-Name
			userId := tr.RequestHeader().Get(common.USERID)
			regionName := tr.RequestHeader().Get(common.REGIONNAME)
			// 检查必需的 header
			if userId == "" {
				return nil, errors.New(
					int(businessErrors.ErrAuthHeaderMissing.HttpCode),
					businessErrors.ErrAuthHeaderMissing.Type,
					"X-User-ID header is missing",
				)
			}

			// 解析用户ID
			userIdUint, err := strconv.ParseUint(userId, 10, 32)
			if err != nil {
				return nil, errors.New(
					int(businessErrors.ErrAuthHeaderInvalid.HttpCode),
					businessErrors.ErrAuthHeaderInvalid.Type,
					"Invalid X-User-ID format",
				)
			}

			var tenantIdUint uint64 = 0

			if needTenant {
				tenantId := tr.RequestHeader().Get(common.TENANTID)
				if tenantId == "" {
					return nil, errors.New(
						int(businessErrors.ErrTenantMissing.HttpCode),
						businessErrors.ErrTenantMissing.Type,
						businessErrors.ErrTenantMissing.Message,
					)
				}
				t, err := strconv.ParseUint(tenantId, 10, 32)
				if err != nil {
					return nil, errors.New(
						int(businessErrors.ErrTenantInvalid.HttpCode),
						businessErrors.ErrTenantInvalid.Type,
						businessErrors.ErrTenantInvalid.Message,
					)
				}
				tenantIdUint = t
			}

			claims := &Claims{
				UserID:     uint32(userIdUint),
				TenantID:   uint32(tenantIdUint),
				RegionName: regionName,
			}
			newCtx := NewContext(ctx, claims)

			return handler(newCtx, req)
		}
	}
}
