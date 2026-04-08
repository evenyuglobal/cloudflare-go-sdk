package main

import (
	"context"
	"io"
	"log"
	"os"
	"strings"

	"cloudflare-go-sdk"
)

func main() {
	ctx := context.Background()
	bucketName := "your bucket"
	objectKey := "your path/file.json"

	// 1. 初始化客户端
	cfClient, _ := cfsdk.NewClient(ctx,
		cfsdk.WithCredentials("your appKey", "your secret"),
		cfsdk.WithAccountID("your accountID"),
	)

	// ==========================================
	// 场景 A：下载文件 (Download)
	// ==========================================
	reader, err := cfClient.R2().Download(ctx, bucketName, objectKey)
	if err != nil {
		log.Printf("下载失败: %v", err)
	} else {
		// 【极其重要】调用方必须 defer Close 释放网络连接！
		defer reader.Close()

		// 将云端文件直接流式写入到本地磁盘
		localFile, _ := os.Create("local_settings.json")
		defer localFile.Close()
		io.Copy(localFile, reader)

		log.Println("文件下载并保存成功")
	}

	// ==========================================
	// 场景 B：更新文件内容 (Update)
	// ==========================================
	newContent := strings.NewReader(`{"status": "updated", "version": 2}`)

	err = cfClient.R2().Upload(ctx, bucketName, objectKey, newContent)
	if err != nil {
		// 如果我们实现了上面的拦截逻辑，当文件不存在时，这里会报错
		log.Printf("更新失败: %v", err)
	} else {
		log.Println("文件内容更新成功 (已覆盖旧版本)")
	}
}
