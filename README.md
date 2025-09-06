# NoWordsServer

# 服务配置方式

通过 **命令行参数** 或 **环境变量** 来配置服务：

| 参数 | 环境变量 | 说明 | 默认值 |
|------|----------|------|--------|
| `--addr` | `NOWORDS_SERVER_ADDRESS` | 服务器监听地址 | `0.0.0.0:8080` |
| `--url-prefix` | `NOWORDS_SERVER_URL_PREFIX` | 服务 URL 前缀 | `/` |
| `--dir` | `NOWORDS_SERVER_DATA_DIR` | 数据存储目录 | `./data` |
| `--dsn` | `NOWORDS_SERVER_DSN` | MySQL DSN，格式：`user:pass@tcp(host)/dbname` | 空 |
| `--serverless` | `NOWORDS_SERVER_SERVERLESS` | 无服务模式 | `false` |
| `--cert` | `NOWORDS_SERVER_CERT` | TLS 证书路径 | 空 |
| `--key` | `NOWORDS_SERVER_KEY` | TLS 证书私钥路径 | 空 |
| `--case-sensitive` | `NOWORDS_SERVER_CASE_SENSITIVE` | 启用 HTTP URL 大小写敏感 | `false` |
| `--strict-routing` | `NOWORDS_SERVER_STRICT_ROUTING` | 启用严格路由区分 | `false` |
| `--reduce-memory-usage` | `NOWORDS_SERVER_REDUCE_MEMORY_USAGE` | 降低内存占用（增加 CPU 消耗） | `false` |
| `--user, -u` | `NOWORDS_SERVER_BASIC_AUTH_USER` | 基础认证用户名 | 空 |
| `--password, -p` | `NOWORDS_SERVER_BASIC_AUTH_PASSWORD` | 基础认证密码 | 空 |
| `--proxy-header` | `NOWORDS_SERVER_PROXY_HEADER` | HTTP 头中远程 IP 地址来源 | 空 |
| `--max-batch-push-count` | `NOWORDS_SERVER_MAX_BATCH_PUSH_COUNT` | 批量推送最大数量，`-1` 表示无限制 | `-1` |
| `--max-apns-client-count` | `NOWORDS_SERVER_MAX_APNS_CLIENT_COUNT` | 最大 APNs 客户端连接数 | `1` |
| `--admins` | `NOWORDS_SERVER_ADMINS` | 管理员 ID 列表 | 空 |
| `--debug` | `NOWORDS_DEBUG` | 启用调试模式 | `false` |
| `--apns-private-key` | `NOWORDS_APPLE_APNS_PRIVATE_KEY` | APNs 私钥路径 | 空 |
| `--topic` | `NOWORDS_APPLE_TOPIC` | APNs Topic | 空 |
| `--key-id` | `NOWORDS_APPLE_KEY_ID` | APNs Key ID | 空 |
| `--team-id` | `NOWORDS_APPLE_TEAM_ID` | APNs Team ID | 空 |
| `--develop, --dev` | `NOWORDS_APPLE_DEVELOP` | 启用 APNs 开发环境 | `false` |
| `--Expired, --ex` | `NOWORDS_EXPIRED_TIME` | 语音过期时间（秒） | `120` |
| `--help, -h` | - | 显示帮助信息 | - |
