package resource

import (
	"context"
	"fmt"
	"time"

	v1 "codeup.aliyun.com/68ce48b215dfc6c8604f8fb0/go-heyin-common/api/gen/go/resource/v1"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/registry"
	kratosGrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"google.golang.org/grpc"
)

// Client 资源服务客户端
//
// 封装了资源服务的 gRPC 调用，提供简洁易用的接口
//
// 使用示例:
//
//	client, err := resource.NewClient(resource.DefaultConfig())
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer client.Close()
//
//	// 批量获取文件URL
//	urls, err := client.BatchGetFileUrls(ctx, fileIDs)
type Client struct {
	config *Config
	conn   *grpc.ClientConn // gRPC 连接
	client v1.FileServiceClient
	logger *log.Helper
}

// NewClient 创建资源服务客户端
//
// 参数:
//   - config: 客户端配置，可以使用 DefaultConfig() 获取默认配置
//
// 返回:
//   - *Client: 客户端实例
//   - error: 创建失败时的错误信息
//
// 使用示例:
//
//	// 使用默认配置
//	client, err := resource.NewClient(resource.DefaultConfig())
//
//	// 自定义配置
//	config := resource.DefaultConfig().
//	    WithAddress("resource-service:9000").
//	    WithTimeout(5 * time.Second)
//	client, err := resource.NewClient(config)
func NewClient(config *Config) (*Client, error) {
	if config == nil {
		config = DefaultConfig()
	}

	// 验证配置
	if config.Address == "" {
		return nil, fmt.Errorf("资源服务地址不能为空")
	}
	if config.Timeout <= 0 {
		config.Timeout = 10 * time.Second
	}

	// 创建日志helper
	logger := log.NewHelper(log.With(
		log.GetLogger(),
		"module", "resource-client",
	))

	// 创建 gRPC 连接
	conn, err := createGRPCConn(config, logger)
	if err != nil {
		return nil, fmt.Errorf("创建 gRPC 连接失败: %w", err)
	}

	// 创建 gRPC 客户端
	fileClient := v1.NewFileServiceClient(conn)

	return &Client{
		config: config,
		conn:   conn,
		client: fileClient,
		logger: logger,
	}, nil
}

// NewClientWithDiscovery 创建带服务发现的资源服务客户端
//
// 参数:
//   - config: 客户端配置
//   - discovery: 服务发现实例（如 Consul）
//
// 返回:
//   - *Client: 客户端实例
//   - error: 创建失败时的错误信息
//
// 使用示例:
//
//	// 创建带 Consul 服务发现的客户端
//	client, err := resource.NewClientWithDiscovery(
//	    resource.DefaultConfig(),
//	    consulDiscovery,
//	)
func NewClientWithDiscovery(config *Config, discovery registry.Discovery) (*Client, error) {
	if config == nil {
		config = DefaultConfig()
	}

	if discovery == nil {
		return nil, fmt.Errorf("服务发现实例不能为空")
	}

	// 验证配置
	if config.Address == "" {
		return nil, fmt.Errorf("资源服务地址不能为空")
	}
	if config.Timeout <= 0 {
		config.Timeout = 10 * time.Second
	}

	// 创建日志helper
	logger := log.NewHelper(log.With(
		log.GetLogger(),
		"module", "resource-client",
	))

	// 创建带服务发现的 gRPC 连接
	conn, err := createGRPCConnWithDiscovery(config, discovery, logger)
	if err != nil {
		return nil, fmt.Errorf("创建 gRPC 连接失败: %w", err)
	}

	// 创建 gRPC 客户端
	fileClient := v1.NewFileServiceClient(conn)

	return &Client{
		config: config,
		conn:   conn,
		client: fileClient,
		logger: logger,
	}, nil
}

// Close 关闭客户端连接
//
// 释放 gRPC 连接资源，应该在程序退出前调用
//
// 使用示例:
//
//	client, err := resource.NewClient(config)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer client.Close()
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// createGRPCConn 创建 gRPC 连接（无服务发现）
func createGRPCConn(config *Config, logger *log.Helper) (*grpc.ClientConn, error) {
	opts := []kratosGrpc.ClientOption{
		kratosGrpc.WithEndpoint(config.Address),
		kratosGrpc.WithTimeout(config.Timeout),
		kratosGrpc.WithMiddleware(
			recovery.Recovery(), // 恢复中间件
		),
	}

	conn, err := kratosGrpc.DialInsecure(
		context.Background(),
		opts...,
	)
	if err != nil {
		return nil, err
	}

	logger.Infof("资源服务客户端连接成功: address=%s, timeout=%v", config.Address, config.Timeout)
	return conn, nil
}

// createGRPCConnWithDiscovery 创建带服务发现的 gRPC 连接
func createGRPCConnWithDiscovery(config *Config, discovery registry.Discovery, logger *log.Helper) (*grpc.ClientConn, error) {
	opts := []kratosGrpc.ClientOption{
		kratosGrpc.WithEndpoint(config.Address),
		kratosGrpc.WithDiscovery(discovery), // 使用服务发现
		kratosGrpc.WithTimeout(config.Timeout),
		kratosGrpc.WithMiddleware(
			recovery.Recovery(), // 恢复中间件
		),
	}

	conn, err := kratosGrpc.DialInsecure(
		context.Background(),
		opts...,
	)
	if err != nil {
		return nil, err
	}

	logger.Infof("资源服务客户端连接成功 (服务发现): endpoint=%s, timeout=%v", config.Address, config.Timeout)
	return conn, nil
}
