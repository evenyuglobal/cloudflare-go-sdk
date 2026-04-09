package r2

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// Upload 执行纯粹的上传（或覆盖）逻辑
func (p *Provider) Upload(ctx context.Context, bucket string, key string, content io.Reader) error {
	o, err := p.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   content,
	})
	if err != nil {
		return fmt.Errorf("r2 upload failed [%s]: %w", key, err)
	}
	return nil
}

// Download 下载文件数据流
func (p *Provider) Download(ctx context.Context, bucket string, key string) (io.ReadCloser, error) {
	output, err := p.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		var noSuchKey *types.NoSuchKey
		if errors.As(err, &noSuchKey) {
			return nil, fmt.Errorf("file not found: %s", key)
		}
		return nil, fmt.Errorf("r2 download failed [%s]: %w", key, err)
	}

	// 返回 Body (io.ReadCloser)，将关闭流的责任交接给调用方
	return output.Body, nil
}

// GeneratePresignedURL 生成限时下载链接
func (p *Provider) GeneratePresignedURL(ctx context.Context, bucket string, key string, expireDuration time.Duration) (string, error) {
	// 1. 初始化预签名客户端 (复用我们已经建好的 S3 Client)
	presignClient := s3.NewPresignClient(p.s3Client)

	// 2. 生成预签名的 GetObject (获取/下载) 请求
	req, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(expireDuration)) // 注入过期时间

	if err != nil {
		return "", err
	}

	// 3. 返回生成的超长安全链接
	return req.URL, nil
}

// GetPublicURL 拼接永久链接
func (p *Provider) GetPublicURL(key string) string {
	// 去除 key 前面可能多余的 "/"，防止出现 https://domain//key
	cleanKey := strings.TrimPrefix(key, "/")
	return fmt.Sprintf("https://%s/%s", p.publicDomain, cleanKey)
}

// UploadAndGetPublicURL 上传并返回永久链接
func (p *Provider) UploadAndGetPublicURL(ctx context.Context, bucket string, key string, content io.Reader) (string, error) {
	if p.publicDomain == "" {
		return "", errors.New("public domain is not configured in SDK")
	}

	// 1. 复用之前写好的普通上传逻辑
	err := p.Upload(ctx, bucket, key, content)
	if err != nil {
		return "", err
	}

	// 2. 上传成功后，返回拼接好的永久 URL
	return p.GetPublicURL(key), nil
}

// UploadAndGetURL 上传并返回链接
func (p *Provider) UploadAndGetURL(ctx context.Context, bucket string, key string, content io.Reader, expireDuration time.Duration) (string, error) {
	// 1. 复用之前写好的普通上传逻辑
	err := p.Upload(ctx, bucket, key, content)
	if err != nil {
		return "", err
	}

	// 2. 上传成功后，返回 URL
	return p.GeneratePresignedURL(ctx, bucket, key, expireDuration)
}
