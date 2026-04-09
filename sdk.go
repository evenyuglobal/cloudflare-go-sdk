package cfsdk

import (
	"context"
	"io"
	"time"
)

// Client 顶层接口定义
type Client interface {
	R2() R2Module
}

// R2Module R2 功能接口
type R2Module interface {
	Upload(ctx context.Context, bucket string, key string, content io.Reader) error
	Download(ctx context.Context, bucket string, key string) (io.ReadCloser, error)
	GetPublicURL(key string) string
	UploadAndGetPublicURL(ctx context.Context, bucket string, key string, content io.Reader) (string, error)
	UploadAndGetURL(ctx context.Context, bucket string, key string, content io.Reader, expireDuration time.Duration) (string, error)
	GeneratePresignedURL(ctx context.Context, bucket string, key string, expireDuration time.Duration) (string, error)
}
