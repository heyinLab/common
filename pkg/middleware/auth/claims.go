package auth

import "context"

type Claims struct {
	UserID     uint32
	TenantID   uint32
	RegionName string
}

// 定义用于在 context 中传递 Claims 的 key
type claimsKey struct{}

// NewContext 将 Claims 存入 context
func NewContext(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, claimsKey{}, claims)
}

// FromContext 从 context 中获取 Claims
func FromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(claimsKey{}).(*Claims)
	return claims, ok
}
