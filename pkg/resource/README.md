# Resource Service 客户端

资源服务的 Go 客户端封装，提供简洁易用的接口用于获取文件信息和URL。

## 功能特性

- ✅ **批量获取文件URL** - 一次请求获取多个文件的URL，包括图片变体
- ✅ **服务发现支持** - 支持 Consul 等服务发现机制
- ✅ **自动重试和恢复** - 内置错误恢复机制
- ✅ **详细日志** - 完整的调用日志，便于问题排查
- ✅ **超时控制** - 可配置的超时时间
- ✅ **中文注释** - 所有代码都有详细的中文注释

## 快速开始

### 安装

```bash
go get codeup.aliyun.com/68ce48b215dfc6c8604f8fb0/go-heyin-common
```

### 基本使用

#### 1. 创建客户端（直连方式）

```go
package main

import (
    "context"
    "fmt"
    "log"

    "codeup.aliyun.com/68ce48b215dfc6c8604f8fb0/go-heyin-common/pkg/resource"
)

func main() {
    // 使用默认配置
    client, err := resource.NewClient(resource.DefaultConfig())
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // 或者自定义配置
    config := resource.DefaultConfig().
        WithAddress("resource-service:9000").
        WithTimeout(5 * time.Second)
    client, err := resource.NewClient(config)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()
}
```

#### 2. 创建客户端（服务发现方式 - 推荐）

```go
package main

import (
    "context"
    "log"

    "codeup.aliyun.com/68ce48b215dfc6c8604f8fb0/go-heyin-common/pkg/resource"
    "github.com/go-kratos/kratos/v2/registry"
    // 导入你的服务发现实现，如 Consul
)

func main() {
    // 创建服务发现实例（以 Consul 为例）
    consulRegistry := NewConsulRegistry(consulConfig)

    // 创建带服务发现的客户端
    config := resource.DefaultConfig().
        WithAddress("discovery:///resource-service")  // 使用服务发现地址

    client, err := resource.NewClientWithDiscovery(config, consulRegistry)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()
}
```

## 核心功能

### 1. 批量获取文件URL（最常用）

适用场景：商品列表页、用户头像、相册等需要展示多个图片的场景

```go
// 场景：商品列表页获取图片
func GetProductImages(client *resource.Client, productIDs []string) error {
    ctx := context.Background()

    // 假设每个商品有一个主图文件ID
    fileIDs := []string{
        "file_id_001",
        "file_id_002",
        "file_id_003",
    }

    // 批量获取文件URL（包含缩略图）
    urls, err := client.BatchGetFileUrls(ctx, fileIDs)
    if err != nil {
        return fmt.Errorf("获取图片URL失败: %w", err)
    }

    // 处理结果
    for fileID, info := range urls {
        if !info.Success {
            log.Printf("文件 %s 获取失败: %s", fileID, info.Error)
            continue
        }

        fmt.Printf("文件ID: %s\n", fileID)
        fmt.Printf("原图URL: %s\n", info.URL)
        fmt.Printf("文件名: %s\n", info.Filename)
        fmt.Printf("大小: %d 字节\n", info.Size)
        fmt.Printf("类型: %s\n", info.ContentType)

        // 获取不同尺寸的缩略图
        if thumbnailURL, ok := info.VariantUrls["thumbnail_200x200"]; ok {
            fmt.Printf("小缩略图: %s\n", thumbnailURL)
        }
        if mediumURL, ok := info.VariantUrls["thumbnail_800x800"]; ok {
            fmt.Printf("中等尺寸: %s\n", mediumURL)
        }

        // 判断URL类型
        if info.IsPublic {
            fmt.Println("这是永久有效的公开URL（可长期缓存）")
        } else {
            fmt.Printf("这是临时URL，%d秒后过期（需定期刷新）\n", info.ExpiresIn)
        }
    }

    return nil
}
```

### 2. 批量获取文件URL（自定义选项）

```go
// 只获取原图URL，不要缩略图（节省响应大小）
func GetOriginalImagesOnly(client *resource.Client, fileIDs []string) error {
    ctx := context.Background()

    urls, err := client.BatchGetFileUrlsWithOptions(ctx, &resource.BatchGetFileUrlsRequest{
        FileIDs:         fileIDs,
        IncludeVariants: false,  // 不包含变体
        ExpiresIn:       7200,    // 2小时有效期
    })
    if err != nil {
        return err
    }

    // 处理结果...
    return nil
}
```

### 3. 获取单个文件元数据

适用场景：需要获取文件的详细信息（上传时间、状态等）

```go
func GetFileInfo(client *resource.Client, fileID string) error {
    ctx := context.Background()

    file, err := client.GetFile(ctx, fileID)
    if err != nil {
        return fmt.Errorf("获取文件信息失败: %w", err)
    }

    fmt.Printf("文件ID: %s\n", file.Id)
    fmt.Printf("文件名: %s\n", file.Filename)
    fmt.Printf("大小: %d 字节\n", file.Size)
    fmt.Printf("类型: %s\n", file.ContentType)
    fmt.Printf("状态: %s\n", file.Status)
    fmt.Printf("上传时间: %s\n", file.CreatedAt.AsTime())
    fmt.Printf("上传者ID: %s\n", file.UploaderId)
    fmt.Printf("分类: %s\n", file.FileCategory)  // image/video/document等

    // 检查自定义元数据
    if file.Metadata != nil {
        fmt.Printf("自定义元数据: %v\n", file.Metadata)
    }

    return nil
}
```

### 4. 获取单个文件下载URL

```go
func DownloadFile(client *resource.Client, fileID string) error {
    ctx := context.Background()

    url, variantUrls, err := client.GetDownloadUrl(ctx, fileID)
    if err != nil {
        return fmt.Errorf("获取下载URL失败: %w", err)
    }

    fmt.Printf("下载URL: %s\n", url)

    // 如果是图片，还会返回变体URL
    for variantID, variantURL := range variantUrls {
        fmt.Printf("变体 %s URL: %s\n", variantID, variantURL)
    }

    return nil
}
```

### 5. 列出文件

适用场景：文件管理后台、用户文件列表

```go
func ListUserFiles(client *resource.Client) error {
    ctx := context.Background()

    // 第1页，每页20条
    files, total, err := client.ListFiles(ctx, 1, 20)
    if err != nil {
        return fmt.Errorf("列出文件失败: %w", err)
    }

    fmt.Printf("总共 %d 个文件\n", total)

    for i, file := range files {
        fmt.Printf("%d. %s (%d 字节) - %s\n",
            i+1, file.Filename, file.Size, file.Status)
    }

    return nil
}
```

## 在 Kratos 服务中集成

### 1. 在 Wire 中配置

```go
// internal/server/server.go
package server

import (
    "codeup.aliyun.com/68ce48b215dfc6c8604f8fb0/go-heyin-common/pkg/resource"
    "github.com/go-kratos/kratos/v2/registry"
    "github.com/google/wire"
)

// ProviderSet 服务提供者集合
var ProviderSet = wire.NewSet(
    NewGRPCServer,
    NewHTTPServer,
    NewConsulRegistry,
    NewConsulDiscovery,
    NewResourceClient,  // 添加资源服务客户端
)

// NewResourceClient 创建资源服务客户端
func NewResourceClient(discovery registry.Discovery) (*resource.Client, error) {
    config := resource.DefaultConfig().
        WithAddress("discovery:///resource-service").  // 使用服务发现
        WithTimeout(10 * time.Second)

    return resource.NewClientWithDiscovery(config, discovery)
}
```

### 2. 在 Biz 层注入使用

```go
// internal/biz/product_usecase.go
package biz

import (
    "context"

    "codeup.aliyun.com/68ce48b215dfc6c8604f8fb0/go-heyin-common/pkg/resource"
)

type ProductUsecase struct {
    resourceClient *resource.Client
    // ... 其他依赖
}

func NewProductUsecase(
    resourceClient *resource.Client,
    // ... 其他依赖
) *ProductUsecase {
    return &ProductUsecase{
        resourceClient: resourceClient,
    }
}

// GetProductDetail 获取商品详情（包含图片）
func (uc *ProductUsecase) GetProductDetail(ctx context.Context, productID string) (*Product, error) {
    // 1. 查询商品信息（包含图片文件ID列表）
    product, err := uc.repo.GetProduct(ctx, productID)
    if err != nil {
        return nil, err
    }

    // 2. 批量获取图片URL
    if len(product.ImageFileIDs) > 0 {
        imageURLs, err := uc.resourceClient.BatchGetFileUrls(ctx, product.ImageFileIDs)
        if err != nil {
            uc.log.Errorf("获取商品图片URL失败: %v", err)
            // 降级处理：图片获取失败不影响商品信息返回
        } else {
            // 填充图片URL到商品对象
            for _, fileID := range product.ImageFileIDs {
                if urlInfo, ok := imageURLs[fileID]; ok && urlInfo.Success {
                    product.Images = append(product.Images, &ProductImage{
                        FileID:       fileID,
                        OriginalURL:  urlInfo.URL,
                        ThumbnailURL: urlInfo.VariantUrls["thumbnail_200x200"],
                        MediumURL:    urlInfo.VariantUrls["thumbnail_800x800"],
                    })
                }
            }
        }
    }

    return product, nil
}
```

## 配置选项

### Config 结构

```go
type Config struct {
    Address     string        // 服务地址
    Timeout     time.Duration // 超时时间
    EnableTrace bool          // 是否启用追踪
    EnableLog   bool          // 是否启用日志
}
```

### 配置示例

```go
// 默认配置
config := resource.DefaultConfig()

// 链式配置
config := resource.DefaultConfig().
    WithAddress("resource-service:9000").
    WithTimeout(5 * time.Second).
    WithTrace(true).
    WithLog(true)

// 或者直接创建
config := &resource.Config{
    Address:     "discovery:///resource-service",
    Timeout:     10 * time.Second,
    EnableTrace: true,
    EnableLog:   true,
}
```

## 最佳实践

### 1. 使用服务发现

推荐使用服务发现而不是硬编码IP地址：

```go
// ✅ 推荐：使用服务发现
config := resource.DefaultConfig().
    WithAddress("discovery:///resource-service")

client, err := resource.NewClientWithDiscovery(config, consulDiscovery)

// ❌ 不推荐：硬编码IP
config := resource.DefaultConfig().
    WithAddress("192.168.1.100:9000")
```

### 2. 合理设置超时时间

根据业务场景设置合适的超时时间：

```go
// 快速响应场景（如用户头像）
config.WithTimeout(3 * time.Second)

// 普通场景（如商品列表）
config.WithTimeout(5 * time.Second)

// 批量处理场景（如后台导出）
config.WithTimeout(30 * time.Second)
```

### 3. 错误处理和降级

图片获取失败不应该影响核心业务：

```go
func GetProductList(ctx context.Context) ([]*Product, error) {
    // 1. 查询商品信息
    products, err := repo.ListProducts(ctx)
    if err != nil {
        return nil, err
    }

    // 2. 收集所有图片ID
    var fileIDs []string
    for _, p := range products {
        if p.ImageFileID != "" {
            fileIDs = append(fileIDs, p.ImageFileID)
        }
    }

    // 3. 批量获取图片URL（失败时降级）
    imageURLs := make(map[string]*resource.FileUrlInfo)
    if len(fileIDs) > 0 {
        urls, err := resourceClient.BatchGetFileUrls(ctx, fileIDs)
        if err != nil {
            log.Errorf("获取图片URL失败（降级处理）: %v", err)
            // 降级：使用默认占位图
        } else {
            imageURLs = urls
        }
    }

    // 4. 填充图片URL
    for _, p := range products {
        if urlInfo, ok := imageURLs[p.ImageFileID]; ok && urlInfo.Success {
            p.ImageURL = urlInfo.URL
            p.ThumbnailURL = urlInfo.VariantUrls["thumbnail_200x200"]
        } else {
            p.ImageURL = "https://cdn.example.com/placeholder.png"  // 占位图
        }
    }

    return products, nil
}
```

### 4. 缓存优化

对于公开URL，可以长期缓存：

```go
func GetImageURL(ctx context.Context, fileID string) (string, error) {
    // 1. 先从缓存获取
    if cachedURL, ok := cache.Get(fileID); ok {
        return cachedURL, nil
    }

    // 2. 调用资源服务
    urls, err := resourceClient.BatchGetFileUrls(ctx, []string{fileID})
    if err != nil {
        return "", err
    }

    urlInfo := urls[fileID]
    if !urlInfo.Success {
        return "", fmt.Errorf(urlInfo.Error)
    }

    // 3. 根据URL类型设置缓存时间
    if urlInfo.IsPublic {
        // 公开URL永久有效，缓存1天
        cache.Set(fileID, urlInfo.URL, 24*time.Hour)
    } else {
        // 私有URL有过期时间，提前5分钟刷新
        cacheDuration := time.Duration(urlInfo.ExpiresIn-300) * time.Second
        cache.Set(fileID, urlInfo.URL, cacheDuration)
    }

    return urlInfo.URL, nil
}
```

## 完整示例

### PWA 服务集成示例

```go
// cmd/pwa/wire.go
//go:build wireinject

package main

import (
    "codeup.aliyun.com/68ce48b215dfc6c8604f8fb0/go-heyin-common/pkg/resource"
    "github.com/google/wire"
)

func wireApp(*conf.Server, *conf.Data, log.Logger) (*kratos.App, func(), error) {
    panic(wire.Build(
        server.ProviderSet,
        data.ProviderSet,
        biz.ProviderSet,
        service.ProviderSet,
        newApp,
    ))
}

// internal/server/server.go
package server

import (
    "codeup.aliyun.com/68ce48b215dfc6c8604f8fb0/go-heyin-common/pkg/resource"
)

var ProviderSet = wire.NewSet(
    NewResourceClient,
    // ... 其他
)

func NewResourceClient(discovery registry.Discovery, logger log.Logger) *resource.Client {
    config := resource.DefaultConfig().
        WithAddress("discovery:///resource-service")

    client, err := resource.NewClientWithDiscovery(config, discovery)
    if err != nil {
        logger.Fatal("创建资源服务客户端失败", err)
    }

    return client
}

// internal/biz/app_usecase.go
package biz

type AppUsecase struct {
    resourceClient *resource.Client
    log            *log.Helper
}

func NewAppUsecase(resourceClient *resource.Client, logger log.Logger) *AppUsecase {
    return &AppUsecase{
        resourceClient: resourceClient,
        log:            log.NewHelper(logger),
    }
}

func (uc *AppUsecase) GetAppWithIcon(ctx context.Context, appID string) (*App, error) {
    // 查询应用信息
    app, err := uc.repo.GetApp(ctx, appID)
    if err != nil {
        return nil, err
    }

    // 获取图标URL
    if app.IconFileID != "" {
        urls, err := uc.resourceClient.BatchGetFileUrls(ctx, []string{app.IconFileID})
        if err != nil {
            uc.log.Warnf("获取应用图标失败: %v", err)
        } else if urlInfo, ok := urls[app.IconFileID]; ok && urlInfo.Success {
            app.IconURL = urlInfo.URL
            app.IconThumbnail = urlInfo.VariantUrls["thumbnail_200x200"]
        }
    }

    return app, nil
}
```

## 常见问题

### Q1: 如何处理批量获取时部分文件失败？

A: `BatchGetFileUrls` 返回的是 map，每个文件都有 `Success` 字段，失败的文件不会影响其他文件：

```go
urls, err := client.BatchGetFileUrls(ctx, fileIDs)
if err != nil {
    return err  // 整个请求失败
}

for fileID, info := range urls {
    if !info.Success {
        log.Errorf("文件 %s 失败: %s", fileID, info.Error)
        continue  // 跳过失败的文件
    }
    // 处理成功的文件
}
```

### Q2: URL什么时候会过期？

A: 取决于文件的公开设置：
- **公开文件**（`IsPublic=true`）：返回CDN URL，永久有效
- **私有文件**（`IsPublic=false`）：返回预签名URL，有过期时间（默认1小时）

### Q3: 如何获取不同尺寸的缩略图？

A: 使用 `VariantUrls` 字段：

```go
urlInfo := urls[fileID]
thumbnailURL := urlInfo.VariantUrls["thumbnail_200x200"]  // 小图
mediumURL := urlInfo.VariantUrls["thumbnail_800x800"]      // 中图
largeURL := urlInfo.VariantUrls["thumbnail_1600x1600"]     // 大图
```

具体有哪些变体取决于上传时的策略配置。

### Q4: 性能怎么样？

A: 内网 gRPC 调用通常在 10ms 以内：
- 单个文件查询：~5-10ms
- 批量查询（50个文件）：~10-20ms
- 建议一次不超过100个文件

## API 文档

详细的 API 文档请参考：
- [资源服务 API 文档](../../docs/integration/service-integration.md)
- [Proto 定义](../../api/protos/resource/v1/)

## 版本历史

- v1.0.0 (2025-11-21): 初始版本
  - 支持批量获取文件URL
  - 支持服务发现
  - 完整的中文注释

## 许可证

内部项目，仅供公司内部使用。
