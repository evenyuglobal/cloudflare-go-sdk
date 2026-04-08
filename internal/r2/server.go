package r2

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// Upload 执行纯粹的上传（或覆盖）逻辑
func (p *Provider) Upload(ctx context.Context, bucket string, key string, content io.Reader) error {
	_, err := p.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   content,
	})
	if err != nil {
		return fmt.Errorf("r2 upload failed [%s]: %w", key, err)
	}
	return nil
}

// Update 语义上的更新。如果你的业务要求“必须存在才能更新”，可以先加一个判断逻辑。
// 如果不需要严格判断，直接调用 Upload 覆盖即可。
func (p *Provider) Update(ctx context.Context, bucket string, key string, content io.Reader) error {
	// 进阶做法：如果你想确保文件之前是存在的，才允许更新（防止误创建）
	// 可以先发起一个 HeadObject 请求验证。如果不需要，直接 return p.Upload(...)
	_, err := p.s3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		// 判断是否是 404 文件不存在错误
		var notFound *types.NotFound
		if errors.As(err, &notFound) {
			return fmt.Errorf("cannot update: file %s does not exist", key)
		}
		// 其他网络或权限错误
		return fmt.Errorf("failed to check file existence before update: %w", err)
	}

	// 确认存在后，执行覆盖上传
	return p.Upload(ctx, bucket, key, content)
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
