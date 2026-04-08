package r2

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Provider struct {
	s3Client *s3.Client
}

// NewProvider 创建底层的 R2 客户端引擎
func NewProvider(ctx context.Context, appKey, secret, accountID string, httpClient *http.Client) (*Provider, error) {
	// 1. 构建 Cloudflare R2 专属的 Endpoint 路由
	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID),
		}, nil
	})

	// 2. 配置加载选项
	loadOptions := []func(*config.LoadOptions) error{
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(appKey, secret, "")),
		config.WithRegion("auto"), // R2 的 Region 必须是 auto
	}

	// 3. 如果上层传入了自定义的 HTTP Client，则应用它
	if httpClient != nil {
		loadOptions = append(loadOptions, config.WithHTTPClient(httpClient))
	}

	// 4. 生成配置并创建客户端
	cfg, err := config.LoadDefaultConfig(ctx, loadOptions...)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)

	return &Provider{
		s3Client: client,
	}, nil
}

// Upload 和 Download 方法实现... (同上一条回复)
