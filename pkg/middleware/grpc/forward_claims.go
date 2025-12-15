package middleware

import (
	"context"
	"strconv"

	"github.com/go-kratos/kratos/v2/middleware"
	authWare "github.com/heyinLab/common/pkg/middleware/auth"
	"github.com/heyinLab/common/pkg/middleware/common"
	"google.golang.org/grpc/metadata"
)

func ForwardClaims() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			// 1. 从当前上下文中获取认证信息 (通常是 HTTP 侧解析 token 后放进去的)
			claims, ok := authWare.FromContext(ctx)
			if ok && claims != nil && claims.UserID != 0 {
				// 2. 将关键字段转换为字符串并放入 gRPC Metadata
				// 使用 AppendToOutgoingContext 可以保留已有的 metadata (如 trace_id)
				ctx = metadata.AppendToOutgoingContext(ctx,
					common.USERID, strconv.FormatUint(uint64(claims.UserID), 10),
					common.TENANTID, strconv.FormatUint(uint64(claims.TenantID), 10),
					common.REGIONNAME, claims.RegionName,
				)
			}
			return handler(ctx, req)
		}
	}
}
