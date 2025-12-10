package resource

import (
	"fmt"
	"time"
)

const (
	// DefaultServiceName 默认的资源服务名称（用于服务发现）
	DefaultServiceName = "resourceServer"

	// DefaultTimeout 默认超时时间
	DefaultTimeout = 10 * time.Second

	// DefaultURLExpiresIn 默认URL过期时间（秒）
	DefaultURLExpiresIn = 3600
)

// InternalConfig 资源内部服务客户端配置
type InternalConfig struct {
	// Endpoint 服务端点
	// 直连方式: "localhost:9000" 或 "192.168.1.100:9000"
	// 服务发现方式: "discovery:///resourceServer"
	Endpoint string

	// ServiceName 服务名称（用于服务发现）
	ServiceName string

	// Timeout 请求超时时间
	Timeout time.Duration
}

// DefaultInternalConfig 返回默认的内部服务客户端配置
//
// 默认配置:
//   - Endpoint: "discovery:///resource-service"
//   - ServiceName: "resource-service"
//   - Timeout: 10s
func DefaultInternalConfig() *InternalConfig {
	return &InternalConfig{
		Endpoint:    fmt.Sprintf("discovery:///%s", DefaultServiceName),
		ServiceName: DefaultServiceName,
		Timeout:     DefaultTimeout,
	}
}

// Validate 验证配置
func (c *InternalConfig) Validate() error {
	if c.Endpoint == "" {
		return fmt.Errorf("服务端点不能为空")
	}
	if c.Timeout <= 0 {
		c.Timeout = DefaultTimeout
	}
	return nil
}

// WithEndpoint 设置服务端点
//
// 参数:
//   - endpoint: 服务端点地址
//
// 示例:
//   - 直连: "localhost:9000"
//   - 服务发现: "discovery:///resource-service"
func (c *InternalConfig) WithEndpoint(endpoint string) *InternalConfig {
	c.Endpoint = endpoint
	return c
}

// WithServiceName 设置服务名称
func (c *InternalConfig) WithServiceName(name string) *InternalConfig {
	c.ServiceName = name
	c.Endpoint = fmt.Sprintf("discovery:///%s", name)
	return c
}

// WithTimeout 设置请求超时时间
func (c *InternalConfig) WithTimeout(timeout time.Duration) *InternalConfig {
	c.Timeout = timeout
	return c
}
