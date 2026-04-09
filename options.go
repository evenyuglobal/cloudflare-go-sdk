package cfsdk

import (
	"errors"
	"net/http"
)

// Config 承载初始化 SDK 所需的所有内部配置
type Config struct {
	AppKey       string
	Secret       string
	AccountID    string
	PublicDomain string // 绑定的公开域名
	HTTPClient   *http.Client
}

// Option 定义了一个操作 Config 的函数类型
type Option func(*Config)

// WithCredentials 注入 R2 的 Access Key 和 Secret Key
func WithCredentials(appKey, secret string) Option {
	return func(c *Config) {
		c.AppKey = appKey
		c.Secret = secret
	}
}

// WithAccountID 注入 Cloudflare 账户 ID (R2 必需)
func WithAccountID(accountID string) Option {
	return func(c *Config) {
		c.AccountID = accountID
	}
}

// WithHTTPClient 允许外部传入自定义的 HTTP 客户端 (例如做代理或连接池优化)
func WithHTTPClient(client *http.Client) Option {
	return func(c *Config) {
		c.HTTPClient = client
	}
}

// WithPublicDomain 注入 R2 存储桶的公开访问域名，例如: "assets.yourcompany.com" 或 "pub-xxxx.r2.dev"
func WithPublicDomain(domain string) Option {
	return func(c *Config) {
		c.PublicDomain = domain
	}
}

// validate 检查核心参数是否齐全
func (c *Config) validate() error {
	if c.AccountID == "" {
		return errors.New("cloudflare account ID is required")
	}
	if c.AppKey == "" || c.Secret == "" {
		return errors.New("credentials (AppKey/Secret) are required")
	}
	return nil
}
