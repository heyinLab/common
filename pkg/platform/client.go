package platform

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/registry"
	kratosGrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	v1 "github.com/heyinLab/common/api/gen/go/platform/v1"
	"google.golang.org/grpc"
)

// Client 平台服务客户端
//
// 聚合了所有平台相关的服务客户端，提供统一的访问入口
//
// 当前支持的服务：
// - IAM: 身份认证和权限管理服务
//
// 使用示例:
//
//	client, err := platform.NewClientWithDiscovery(
//	    platform.DefaultConfig(),
//	    consulDiscovery,
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer client.Close()
//
//	// 使用 IAM 服务
//	tree, err := client.IAM().GetTenantPermissionsTree(ctx, &platform.GetTenantPermissionsTreeOptions{
//	    Status: "GA",
//	})
type Client struct {
	config *Config
	conn   *grpc.ClientConn
	logger *log.Helper

	// 子服务客户端
	iamClient *IAMClient
}

// NewClient 创建平台服务客户端（直连方式）
//
// 参数:
//   - config: 客户端配置，可以使用 DefaultConfig() 获取默认配置
//
// 返回:
//   - *Client: 客户端实例
//   - error: 创建失败时的错误信息
func NewClient(config *Config) (*Client, error) {
	if config == nil {
		config = DefaultConfig()
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	logger := log.NewHelper(log.With(
		log.GetLogger(),
		"module", "platform-client",
	))

	conn, err := createGRPCConn(config, nil, logger)
	if err != nil {
		return nil, fmt.Errorf("创建 gRPC 连接失败: %w", err)
	}

	return &Client{
		config:    config,
		conn:      conn,
		logger:    logger,
		iamClient: newIAMClient(conn, logger),
	}, nil
}

// NewClientWithDiscovery 创建带服务发现的平台服务客户端
//
// 参数:
//   - config: 客户端配置
//   - discovery: 服务发现实例（如 Consul）
//
// 返回:
//   - *Client: 客户端实例
//   - error: 创建失败时的错误信息
func NewClientWithDiscovery(config *Config, discovery registry.Discovery) (*Client, error) {
	if config == nil {
		config = DefaultConfig()
	}

	if discovery == nil {
		return nil, fmt.Errorf("服务发现实例不能为空")
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	logger := log.NewHelper(log.With(
		log.GetLogger(),
		"module", "platform-client",
	))

	conn, err := createGRPCConn(config, discovery, logger)
	if err != nil {
		return nil, fmt.Errorf("创建 gRPC 连接失败: %w", err)
	}

	logger.Infof("平台服务客户端连接成功 (服务发现): endpoint=%s, timeout=%v", config.Endpoint, config.Timeout)

	return &Client{
		config:    config,
		conn:      conn,
		logger:    logger,
		iamClient: newIAMClient(conn, logger),
	}, nil
}

// Close 关闭客户端连接
//
// 释放 gRPC 连接资源，应该在程序退出前调用
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// ========== 服务访问器 ==========

// IAM 返回 IAM 服务客户端
//
// 用于访问身份认证和权限管理相关功能
//
// 使用示例:
//
//	tree, err := client.IAM().GetTenantPermissionsTree(ctx, &GetTenantPermissionsTreeOptions{
//	    Status: "GA",
//	})
func (c *Client) IAM() *IAMClient {
	return c.iamClient
}

// ========== 内部函数 ==========

// createGRPCConn 创建 gRPC 连接
func createGRPCConn(config *Config, discovery registry.Discovery, logger *log.Helper) (*grpc.ClientConn, error) {
	opts := []kratosGrpc.ClientOption{
		kratosGrpc.WithEndpoint(config.Endpoint),
		kratosGrpc.WithTimeout(config.Timeout),
		kratosGrpc.WithMiddleware(
			recovery.Recovery(),
		),
	}

	// 如果有服务发现，添加服务发现选项
	if discovery != nil {
		opts = append(opts, kratosGrpc.WithDiscovery(discovery))
	}

	conn, err := kratosGrpc.DialInsecure(
		context.Background(),
		opts...,
	)
	if err != nil {
		return nil, err
	}

	logger.Infof("平台服务客户端连接成功: endpoint=%s, timeout=%v", config.Endpoint, config.Timeout)

	return conn, nil
}

// ========== IAM 客户端 ==========

// IAMClient IAM 服务客户端
//
// 提供身份认证和权限管理相关功能
type IAMClient struct {
	client v1.PlatformIamServiceClient
	logger *log.Helper
}

// newIAMClient 创建 IAM 客户端
func newIAMClient(conn *grpc.ClientConn, logger *log.Helper) *IAMClient {
	return &IAMClient{
		client: v1.NewPlatformIamServiceClient(conn),
		logger: logger,
	}
}

// GetTenantPermissionsTreeOptions 获取租户权限树的选项
type GetTenantPermissionsTreeOptions struct {
	// Status 权限状态过滤：DEV, BETA, GA
	Status string
	// Timeout 自定义超时时间（可选）
	Timeout time.Duration
}

// GetTenantPermissionsTree 获取租户权限树
//
// 用于获取完整的租户权限树结构，包含嵌套的 children 节点
//
// 参数:
//   - ctx: 上下文
//   - opts: 查询选项（可选）
//
// 返回:
//   - []*v1.TenantPermissionTreeNode: 权限树节点列表
//   - uint32: 总数量
//   - error: 错误信息
//
// 使用场景：
//   - 前端菜单渲染
//   - 权限分配界面
//   - 角色权限配置
//
// 使用示例:
//
//	// 获取所有状态的权限
//	tree, total, err := client.IAM().GetTenantPermissionsTree(ctx, nil)
//
//	// 只获取正式发布的权限
//	tree, total, err := client.IAM().GetTenantPermissionsTree(ctx, &platform.GetTenantPermissionsTreeOptions{
//	    Status: "GA",
//	})
func (c *IAMClient) GetTenantPermissionsTree(ctx context.Context, opts *GetTenantPermissionsTreeOptions) ([]*v1.TenantPermissionTreeNode, uint32, error) {
	// 设置超时
	if opts != nil && opts.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, opts.Timeout)
		defer cancel()
	}

	// 构建请求
	req := &v1.GetTenantPermissionsTreeRequest{}
	if opts != nil && opts.Status != "" {
		req.Status = &opts.Status
	}

	// 执行请求
	resp, err := c.client.GetTenantPermissionsTree(ctx, req)
	if err != nil {
		c.logger.WithContext(ctx).Errorf("获取租户权限树失败: status=%s, error=%v",
			getStringValue(req.Status), err)
		return nil, 0, err
	}

	c.logger.WithContext(ctx).Infof("获取租户权限树成功: status=%s, total=%d",
		getStringValue(req.Status), resp.Total)

	return resp.Tree, resp.Total, nil
}

// GetPermissionCodesByProductOptions 根据产品ID获取权限codes的选项
type GetPermissionCodesByProductOptions struct {
	// Status 权限状态过滤：DEV, BETA, GA（可选）
	Status string
	// Timeout 自定义超时时间（可选）
	Timeout time.Duration
}

// GetPermissionCodesByProduct 根据产品ID获取权限codes
//
// 用于获取指定产品下的所有权限代码（扁平列表），用于权限校验
//
// 参数:
//   - ctx: 上下文
//   - productID: 产品ID（必填）
//   - opts: 查询选项（可选）
//
// 返回:
//   - []string: 权限代码列表
//   - uint32: 总数量
//   - error: 错误信息
//
// 使用场景：
//   - 权限校验（检查用户是否拥有某个产品的权限）
//   - 权限初始化（批量授权）
//   - 权限统计
//
// 使用示例:
//
//	// 获取产品的所有权限codes
//	codes, total, err := client.IAM().GetPermissionCodesByProduct(ctx, 1, nil)
//
//	// 只获取正式发布的权限codes
//	codes, total, err := client.IAM().GetPermissionCodesByProduct(ctx, 1, &platform.GetPermissionCodesByProductOptions{
//	    Status: "GA",
//	})
func (c *IAMClient) GetPermissionCodesByProduct(ctx context.Context, productCode string, opts *GetPermissionCodesByProductOptions) ([]string, uint32, error) {
	// 设置超时
	if opts != nil && opts.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, opts.Timeout)
		defer cancel()
	}

	// 构建请求
	req := &v1.GetPermissionCodesByProductRequest{
		ProductCode: productCode,
	}
	if opts != nil && opts.Status != "" {
		req.Status = &opts.Status
	}

	// 执行请求
	resp, err := c.client.GetPermissionCodesByProduct(ctx, req)
	if err != nil {
		c.logger.WithContext(ctx).Errorf("获取产品权限codes失败: product_id=%d, status=%s, error=%v",
			productCode, getStringValue(req.Status), err)
		return nil, 0, err
	}

	c.logger.WithContext(ctx).Infof("获取产品权限codes成功: product_id=%d, status=%s, total=%d",
		productCode, getStringValue(req.Status), resp.Total)

	return resp.Codes, resp.Total, nil
}

// ========== 辅助函数 ==========

// getStringValue 获取指针字符串的值
func getStringValue(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}
