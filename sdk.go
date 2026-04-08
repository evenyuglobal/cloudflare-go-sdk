package cfsdk

import (
	"cloudflare-go-sdk/internal/r2"
	"context"
	"fmt"
	"io"
)

// Client 顶层接口定义
type Client interface {
	R2() R2Module
}

// R2Module R2 功能接口
type R2Module interface {
	Upload(ctx context.Context, bucket string, key string, content io.Reader) error
	Download(ctx context.Context, bucket string, key string) (io.ReadCloser, error)
}

// cfClient 是 Client 接口的内部实现
type cfClient struct {
	r2 R2Module
}

func (c *cfClient) R2() R2Module { return c.r2 }

// NewClient 是创建 Cloudflare SDK 客户端的工厂函数
func NewClient(ctx context.Context, opts ...Option) (Client, error) {
	// 1. 初始化默认配置
	cfg := &Config{}

	// 2. 遍历并应用外部传入的选项
	for _, opt := range opts {
		opt(cfg)
	}

	// 3. 校验必填项
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	// 4. 初始化内部的 R2 驱动
	r2Provider, err := r2.NewProvider(ctx, cfg.AppKey, cfg.Secret, cfg.AccountID, cfg.HTTPClient)
	if err != nil {
		return nil, fmt.Errorf("failed to init R2 provider: %w", err)
	}

	// 5. 返回组装好的客户端
	return &cfClient{
		r2: r2Provider,
	}, nil
}
