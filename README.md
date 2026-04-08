cloudflare-go-sdk/
├── go.mod                  # 模块定义文件，声明模块路径和外部依赖 (如 aws-sdk-go-v2)
├── go.sum                  # 依赖版本的校验和
│
├── sdk.go                  # 【公共入口】定义 Client 接口，包含 NewClient() 构造工厂
├── options.go              # 【配置选项】定义 Config 和 Option 函数 (WithCredentials 等)
│                           # 【功能接口】定义 R2Module 接口 (Upload, Update, Download)
│
├── internal/               # 【核心隔离区】该目录下的代码外部服务无法 import
│   └── r2/                 # R2 存储的具体实现模块
│       ├── client.go       # 定义 Provider 结构体和 NewProvider()，初始化 S3 Client
│       └── server.go       # 实现 Provider 上的 Upload, Update, Download 方法
│
└── examples/               # 【示例代码】提供给调用者看的快速上手 Demo