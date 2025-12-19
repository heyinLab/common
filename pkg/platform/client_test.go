package platform

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	consulapi "github.com/hashicorp/consul/api"
)

func TestGetTenantPermissionsTree(t *testing.T) {
	// 创建 Consul 客户端
	consulClient, err := consulapi.NewClient(consulapi.DefaultConfig())
	if err != nil {
		t.Skipf("无法连接到 Consul: %v", err)
		return
	}

	// 创建 Consul 服务发现
	discovery := consul.New(consulClient)

	// 创建平台服务客户端
	client, err := NewClientWithDiscovery(DefaultConfig(), discovery)
	if err != nil {
		t.Fatalf("创建客户端失败: %v", err)
	}
	defer client.Close()

	// 测试获取权限树
	ctx := context.Background()
	tree, total, err := client.IAM().GetTenantPermissionsTree(ctx, nil)
	if err != nil {
		t.Logf("获取权限树失败（可能服务未启动）: %v", err)
		t.Skip("跳过测试，服务可能未启动")
		return
	}

	t.Logf("成功获取权限树，总数: %d", total)
	for i, node := range tree {
		if i < 3 { // 只打印前3个节点避免输出过多
			t.Logf("权限节点: ID=%d, Name=%s, Code=%s", node.Id, node.Name, getCodeValue(node.Code))
		}
	}
}

func TestGetTenantPermissionsTreeWithStatus(t *testing.T) {
	consulClient, err := consulapi.NewClient(consulapi.DefaultConfig())
	if err != nil {
		t.Skipf("无法连接到 Consul: %v", err)
		return
	}

	discovery := consul.New(consulClient)
	client, err := NewClientWithDiscovery(DefaultConfig(), discovery)
	if err != nil {
		t.Fatalf("创建客户端失败: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	tree, total, err := client.IAM().GetTenantPermissionsTree(ctx, &GetTenantPermissionsTreeOptions{})
	if err != nil {
		t.Logf("获取权限树失败（可能服务未启动）: %v", err)
		t.Skip("跳过测试，服务可能未启动")
		return
	}
	fmt.Printf("%s", tree)
	t.Logf("成功获取 权限树，总数: %d", total)
}

func TestGetPermissionCodesByProduct(t *testing.T) {
	consulClient, err := consulapi.NewClient(consulapi.DefaultConfig())
	if err != nil {
		t.Skipf("无法连接到 Consul: %v", err)
		return
	}

	discovery := consul.New(consulClient)
	client, err := NewClientWithDiscovery(DefaultConfig(), discovery)
	if err != nil {
		t.Fatalf("创建客户端失败: %v", err)
	}
	defer client.Close()

	// 测试获取产品权限codes
	ctx := context.Background()
	productCode := "111"

	codes, total, err := client.IAM().GetPermissionCodesByProduct(ctx, productCode, nil)
	if err != nil {
		t.Logf("获取产品权限codes失败（可能服务未启动）: %v", err)
		t.Skip("跳过测试，服务可能未启动")
		return
	}

	t.Logf("成功获取产品 %d 的权限codes，总数: %d", productCode, total)
	if len(codes) > 0 {
		t.Logf("权限codes: %v", codes)
	}
}

func TestGetPermissionCodesByProductWithStatus(t *testing.T) {
	consulClient, err := consulapi.NewClient(consulapi.DefaultConfig())
	if err != nil {
		t.Skipf("无法连接到 Consul: %v", err)
		return
	}

	discovery := consul.New(consulClient)
	client, err := NewClientWithDiscovery(DefaultConfig(), discovery)
	if err != nil {
		t.Fatalf("创建客户端失败: %v", err)
	}
	defer client.Close()

	// 测试只获取 GA 状态的产品权限codes
	ctx := context.Background()
	productCode := "111"

	codes, total, err := client.IAM().GetPermissionCodesByProduct(ctx, productCode, &GetPermissionCodesByProductOptions{
		Status: "GA",
	})

	if err != nil {
		t.Logf("获取产品权限codes失败（可能服务未启动）: %v", err)
		t.Skip("跳过测试，服务可能未启动")
		return
	}

	t.Logf("成功获取产品 %d 的 GA 状态权限codes，总数: %d, codes: %v", productCode, total, codes)
}

// 辅助函数
func getCodeValue(code *string) string {
	if code == nil {
		return ""
	}
	return *code
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
