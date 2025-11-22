package resource

import (
	"context"
	"fmt"
	"time"

	v1 "codeup.aliyun.com/68ce48b215dfc6c8604f8fb0/go-heyin-common/api/gen/go/resource/v1"
)

// FileUrlInfo 文件URL信息（简化版）
//
// 包含文件的访问URL和基本元数据
type FileUrlInfo struct {
	// URL 文件访问URL（原图）
	URL string

	// VariantUrls 图片变体URL（缩略图、裁剪图等）
	// key: 变体ID, value: 变体URL
	VariantUrls map[string]string

	// IsPublic 是否为公开URL
	// true: 永久有效的CDN URL
	// false: 临时有效的预签名URL
	IsPublic bool

	// ExpiresIn URL过期时间（秒）
	// IsPublic=true时为0（永久有效）
	// IsPublic=false时为预签名URL的剩余有效期
	ExpiresIn int64

	// Filename 文件名
	Filename string

	// Size 文件大小（字节）
	Size int64

	// ContentType MIME类型
	// 示例: "image/jpeg", "image/png", "video/mp4"
	ContentType string

	// Success 是否成功获取URL
	Success bool

	// Error 错误信息（Success=false时）
	Error string
}

// BatchGetFileUrlsRequest 批量获取文件URL请求
type BatchGetFileUrlsRequest struct {
	// FileIDs 文件ID列表（必填）
	// 限制：建议单次不超过100个
	FileIDs []string

	// IncludeVariants 是否包含图片变体URL（可选，默认true）
	// true: 返回原图 + 所有缩略图/裁剪图
	// false: 只返回原图URL
	IncludeVariants bool

	// ExpiresIn URL有效期（秒，可选）
	// 默认: 3600（1小时）
	// 仅对私有文件有效，公开文件忽略此参数
	ExpiresIn int64
}

// BatchGetFileUrlsResponse 批量获取文件URL响应
type BatchGetFileUrlsResponse struct {
	// Results 文件URL信息映射
	// key: 文件ID, value: URL信息
	// 注意：即使某个文件失败，其他文件仍正常返回
	Results map[string]*FileUrlInfo

	// ExpiresIn URL有效期（秒）
	// 私有文件的URL过期时间
	// 公开文件忽略此值（永久有效）
	ExpiresIn int64
}

// BatchGetFileUrls 批量获取文件URL
//
// 最常用的方法，用于获取多个文件的访问URL，支持图片变体
//
// 参数:
//   - ctx: 上下文，用于超时控制和取消
//   - fileIDs: 文件ID列表，建议不超过100个
//
// 返回:
//   - map[string]*FileUrlInfo: 文件URL信息映射，key为文件ID
//   - error: 错误信息
//
// 使用示例:
//
//	// 获取商品图片URL
//	urls, err := client.BatchGetFileUrls(ctx, []string{
//	    "file_id_1",
//	    "file_id_2",
//	    "file_id_3",
//	})
//	if err != nil {
//	    return err
//	}
//
//	for fileID, info := range urls {
//	    if !info.Success {
//	        log.Errorf("文件 %s 获取失败: %s", fileID, info.Error)
//	        continue
//	    }
//	    fmt.Printf("原图URL: %s\n", info.URL)
//	    fmt.Printf("缩略图URL: %s\n", info.VariantUrls["thumbnail"])
//	}
func (c *Client) BatchGetFileUrls(ctx context.Context, fileIDs []string) (map[string]*FileUrlInfo, error) {
	return c.BatchGetFileUrlsWithOptions(ctx, &BatchGetFileUrlsRequest{
		FileIDs:         fileIDs,
		IncludeVariants: true, // 默认包含变体
		ExpiresIn:       3600, // 默认1小时
	})
}

// BatchGetFileUrlsWithOptions 批量获取文件URL（带选项）
//
// 提供更多自定义选项的批量获取方法
//
// 参数:
//   - ctx: 上下文
//   - req: 请求参数
//
// 返回:
//   - map[string]*FileUrlInfo: 文件URL信息映射
//   - error: 错误信息
//
// 使用示例:
//
//	// 只获取原图URL，不要变体
//	urls, err := client.BatchGetFileUrlsWithOptions(ctx, &resource.BatchGetFileUrlsRequest{
//	    FileIDs:         fileIDs,
//	    IncludeVariants: false,  // 不包含变体
//	    ExpiresIn:       7200,    // 2小时有效期
//	})
func (c *Client) BatchGetFileUrlsWithOptions(ctx context.Context, req *BatchGetFileUrlsRequest) (map[string]*FileUrlInfo, error) {
	if req == nil || len(req.FileIDs) == 0 {
		return make(map[string]*FileUrlInfo), nil
	}

	// 参数验证
	if len(req.FileIDs) > 100 {
		c.logger.Warnf("批量获取文件URL: 文件数量过多(%d)，建议不超过100个", len(req.FileIDs))
	}

	// 设置默认值
	if req.ExpiresIn <= 0 {
		req.ExpiresIn = 3600 // 默认1小时
	}

	// 创建超时上下文
	timeoutCtx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	startTime := time.Now()

	// 调用 gRPC 方法
	resp, err := c.client.BatchGetFileUrls(timeoutCtx, &v1.BatchGetFileUrlsRequest{
		FileIds:         req.FileIDs,
		IncludeVariants: req.IncludeVariants,
		ExpiresIn:       req.ExpiresIn,
	})
	if err != nil {
		c.logger.Errorf("批量获取文件URL失败: count=%d, error=%v, elapsed=%v",
			len(req.FileIDs), err, time.Since(startTime))
		return nil, fmt.Errorf("批量获取文件URL失败: %w", err)
	}

	elapsed := time.Since(startTime)
	c.logger.Infof("批量获取文件URL成功: count=%d, elapsed=%v", len(req.FileIDs), elapsed)

	// 转换响应
	results := make(map[string]*FileUrlInfo)
	for fileID, info := range resp.Results {
		results[fileID] = &FileUrlInfo{
			URL:         info.Url,
			VariantUrls: info.VariantUrls,
			IsPublic:    info.IsPublic,
			ExpiresIn:   info.ExpiresIn,
			Filename:    info.Filename,
			Size:        info.Size,
			ContentType: info.ContentType,
			Success:     info.Success,
			Error:       info.Error,
		}
	}

	return results, nil
}

// GetFile 获取文件元数据
//
// 获取文件的完整元数据信息，不包含下载URL
//
// 参数:
//   - ctx: 上下文
//   - fileID: 文件ID
//
// 返回:
//   - *v1.FileObject: 文件对象（包含所有元数据）
//   - error: 错误信息
//
// 使用示例:
//
//	file, err := client.GetFile(ctx, "file_id_123")
//	if err != nil {
//	    return err
//	}
//	fmt.Printf("文件名: %s\n", file.Filename)
//	fmt.Printf("大小: %d 字节\n", file.Size)
//	fmt.Printf("上传时间: %s\n", file.CreatedAt.AsTime())
func (c *Client) GetFile(ctx context.Context, fileID string) (*v1.FileObject, error) {
	if fileID == "" {
		return nil, fmt.Errorf("文件ID不能为空")
	}

	// 创建超时上下文
	timeoutCtx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	startTime := time.Now()

	// 调用 gRPC 方法
	resp, err := c.client.GetFile(timeoutCtx, &v1.GetFileRequest{
		FileId: fileID,
	})
	if err != nil {
		c.logger.Errorf("获取文件元数据失败: file_id=%s, error=%v, elapsed=%v",
			fileID, err, time.Since(startTime))
		return nil, fmt.Errorf("获取文件元数据失败: %w", err)
	}

	elapsed := time.Since(startTime)
	c.logger.Infof("获取文件元数据成功: file_id=%s, filename=%s, elapsed=%v",
		fileID, resp.Filename, elapsed)

	return resp, nil
}

// GetDownloadUrl 获取文件下载URL
//
// 获取单个文件的下载URL，支持获取图片变体URL
//
// 参数:
//   - ctx: 上下文
//   - fileID: 文件ID
//
// 返回:
//   - string: 下载URL
//   - map[string]string: 变体URL（图片文件时返回）
//   - error: 错误信息
//
// 使用示例:
//
//	url, variantUrls, err := client.GetDownloadUrl(ctx, "file_id_123")
//	if err != nil {
//	    return err
//	}
//	fmt.Printf("原图URL: %s\n", url)
//	fmt.Printf("缩略图URL: %s\n", variantUrls["thumbnail_200x200"])
func (c *Client) GetDownloadUrl(ctx context.Context, fileID string) (string, map[string]string, error) {
	if fileID == "" {
		return "", nil, fmt.Errorf("文件ID不能为空")
	}

	// 创建超时上下文
	timeoutCtx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	startTime := time.Now()

	// 调用 gRPC 方法
	resp, err := c.client.GetDownloadUrl(timeoutCtx, &v1.GetDownloadUrlRequest{
		FileId:    fileID,
		ExpiresIn: 3600, // 默认1小时
	})
	if err != nil {
		c.logger.Errorf("获取下载URL失败: file_id=%s, error=%v, elapsed=%v",
			fileID, err, time.Since(startTime))
		return "", nil, fmt.Errorf("获取下载URL失败: %w", err)
	}

	elapsed := time.Since(startTime)
	c.logger.Infof("获取下载URL成功: file_id=%s, filename=%s, elapsed=%v",
		fileID, resp.Filename, elapsed)

	return resp.DownloadUrl, resp.VariantUrls, nil
}

// ListFiles 列出文件
//
// 分页查询文件列表
//
// 参数:
//   - ctx: 上下文
//   - page: 页码（从1开始）
//   - pageSize: 每页数量（建议不超过100）
//
// 返回:
//   - []*v1.FileObject: 文件列表
//   - int32: 总记录数
//   - error: 错误信息
//
// 使用示例:
//
//	files, total, err := client.ListFiles(ctx, 1, 20)
//	if err != nil {
//	    return err
//	}
//	fmt.Printf("总共 %d 个文件，当前页 %d 个\n", total, len(files))
//	for _, file := range files {
//	    fmt.Printf("- %s (%d 字节)\n", file.Filename, file.Size)
//	}
func (c *Client) ListFiles(ctx context.Context, page, pageSize int32) ([]*v1.FileObject, int32, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		c.logger.Warnf("列出文件: 每页数量过大(%d)，建议不超过100", pageSize)
	}

	// 创建超时上下文
	timeoutCtx, cancel := context.WithTimeout(ctx, c.config.Timeout)
	defer cancel()

	startTime := time.Now()

	// 调用 gRPC 方法
	resp, err := c.client.ListFiles(timeoutCtx, &v1.ListFilesRequest{
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		c.logger.Errorf("列出文件失败: page=%d, pageSize=%d, error=%v, elapsed=%v",
			page, pageSize, err, time.Since(startTime))
		return nil, 0, fmt.Errorf("列出文件失败: %w", err)
	}

	elapsed := time.Since(startTime)
	c.logger.Infof("列出文件成功: page=%d, pageSize=%d, count=%d, total=%d, elapsed=%v",
		page, pageSize, len(resp.Files), resp.Total, elapsed)

	return resp.Files, resp.Total, nil
}
