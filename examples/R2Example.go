package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	cfsdk "github.com/evenyuglobal/cloudflare-go-sdk"
)

func main() {
	ctx := context.Background()
	appKey := "app-key"
	secret := "secret-key"
	accountId := "account-id"
	bucketName := "bucket"
	//localFilePath := "/Users/evenyu/Downloads/file_3.jpg"
	objectKey := "path/file.jpg"
	publicDomain := "xxxx.com"

	// 1. 初始化客户端，带上你的公开域名
	cfClient, _ := cfsdk.NewClient(ctx,
		cfsdk.WithCredentials(appKey, secret),
		cfsdk.WithAccountID(accountId),
		cfsdk.WithPublicDomain(publicDomain), // 👈 这里填入你的域名
	)
	generatePresignedURL(cfClient, ctx, bucketName, objectKey)
	//uploadAndGetPublicURL(localFilePath, cfClient, ctx, bucketName, objectKey)
}

func uploadAndGetPublicURL(localFilePath string, cfClient cfsdk.Client, ctx context.Context, bucketName string, objectKey string) {
	// 2. 准备上传文件
	file, _ := os.Open(localFilePath)
	defer file.Close()

	// 3. 一键上传并获取永久链接
	publicURL, err := cfClient.R2().UploadAndGetPublicURL(ctx, bucketName, objectKey, file)

	if err != nil {
		log.Fatalf("上传失败: %v", err)
	}

	// 成功！将这个链接存入数据库，或者返回给前端展示
	fmt.Printf("文件上传成功！永久访问链接为:\n%s\n", publicURL)
	// 输出: https://assets.my-website.com/images/2026/logo.png
}

func generatePresignedURL(cfClient cfsdk.Client, ctx context.Context, bucketName string, objectKey string) string {
	// 设置有效期为 15 分钟
	temporaryURL, err := cfClient.R2().GeneratePresignedURL(ctx, bucketName, objectKey, 15*time.Minute)

	if err != nil {
		log.Fatalf("生成临时链接失败: %v", err)
	}

	fmt.Printf("生成的限时链接为 (15分钟内有效):\n%s\n", temporaryURL)
	return temporaryURL
}
