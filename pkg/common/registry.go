package common

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/registry"
	consulAPI "github.com/hashicorp/consul/api"

	consul "github.com/go-kratos/kratos/contrib/registry/consul/v2"
)

type customRegistrar struct {
	registry.Registrar
	serviceHost string
	portMap     map[string]string // 容器内端口 -> 宿主机端口 (e.g. "8000" -> "25678")
}

func (r *customRegistrar) Register(ctx context.Context, service *registry.ServiceInstance) error {
	if r.serviceHost == "" {
		return r.Registrar.Register(ctx, service)
	}

	var httpEndpoints []string
	var otherEndpoints []string

	for _, endpoint := range service.Endpoints {
		// 1. 解析协议和地址
		parts := strings.SplitN(endpoint, "://", 2)
		if len(parts) != 2 {
			otherEndpoints = append(otherEndpoints, endpoint)
			continue
		}
		protocol, addr := parts[0], parts[1]

		// 2. 提取容器内端口
		_, port, err := net.SplitHostPort(addr)
		if err != nil {
			otherEndpoints = append(otherEndpoints, endpoint)
			continue
		}

		// 3. 映射宿主机端口
		finalPort := port
		if mappedPort, ok := r.portMap[port]; ok && mappedPort != "" {
			finalPort = mappedPort
		}

		// 4. 组装宿主机真实的访问地址
		finalAddr := fmt.Sprintf("%s://%s:%s", protocol, r.serviceHost, finalPort)

		if protocol == "http" {
			httpEndpoints = append(httpEndpoints, finalAddr)
		} else {
			otherEndpoints = append(otherEndpoints, finalAddr)
		}
	}

	service.Endpoints = append(httpEndpoints, otherEndpoints...)

	log.Context(ctx).Infof("服务注册中: Name=%s, 最终地址=%v", service.Name, service.Endpoints)
	return r.Registrar.Register(ctx, service)
}

func (r *customRegistrar) Deregister(ctx context.Context, service *registry.ServiceInstance) error {
	return r.Registrar.Deregister(ctx, service)
}

// NewConsulRegistrar 初始化通用的 Consul 注册器
func NewConsulRegistrar(consulAddr string, tags []string) registry.Registrar {
	// 初始化 Consul Client
	c := consulAPI.DefaultConfig()
	c.Address = consulAddr
	cli, err := consulAPI.NewClient(c)
	if err != nil {
		panic(fmt.Sprintf("Consul 客户端初始化失败: %v", err))
	}

	// 1. 获取宿主机 IP
	host := os.Getenv("SERVICE_HOST")
	if host == "" {
		// 如果没传，尝试自动获取（保底）
		host = getLocalIP()
	}

	// 2. 解析端口映射环境变量
	portMap := make(map[string]string)
	parseEnvToMap(os.Getenv("HTTP_PORT_MAP"), portMap)
	parseEnvToMap(os.Getenv("GRPC_PORT_MAP"), portMap)

	// 3. 创建基础注册器
	baseRegistrar := consul.New(cli,
		consul.WithHealthCheck(true),
		consul.WithTags(tags))

	return &customRegistrar{
		Registrar:   baseRegistrar,
		serviceHost: host,
		portMap:     portMap,
	}
}

// 辅助工具：解析 "8000:12345" 格式
func parseEnvToMap(envVal string, m map[string]string) {
	if envVal == "" {
		return
	}
	kv := strings.Split(envVal, ":")
	if len(kv) == 2 {
		m[kv[0]] = kv[1]
	}
}

// 辅助工具：获取本地 IP
func getLocalIP() string {
	address, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}
	for _, address := range address {
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			return ipNet.IP.String()
		}
	}
	return "127.0.0.1"
}
