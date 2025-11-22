package resource

import "time"

// Config 资源服务客户端配置
type Config struct {
	// Address 资源服务的 gRPC 地址
	// 支持直接地址: "resource-service:9000"
	// 支持服务发现: "discovery:///resource-service"
	Address string

	// Timeout 调用超时时间，默认: 10秒
	Timeout time.Duration

	// EnableTrace 是否启用链路追踪，默认: true
	EnableTrace bool

	// EnableLog 是否启用日志，默认: true
	EnableLog bool
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Address:     "discovery:///resource-service", // 默认使用服务发现
		Timeout:     10 * time.Second,                // 默认10秒超时
		EnableTrace: true,                            // 默认启用追踪
		EnableLog:   true,                            // 默认启用日志
	}
}

// WithAddress 设置服务地址
func (c *Config) WithAddress(address string) *Config {
	c.Address = address
	return c
}

// WithTimeout 设置超时时间
func (c *Config) WithTimeout(timeout time.Duration) *Config {
	c.Timeout = timeout
	return c
}

// WithTrace 设置是否启用追踪
func (c *Config) WithTrace(enable bool) *Config {
	c.EnableTrace = enable
	return c
}

// WithLog 设置是否启用日志
func (c *Config) WithLog(enable bool) *Config {
	c.EnableLog = enable
	return c
}
