# 系统模块设计

## 1.配置模块（Configuration Module）

功能: 读取用户配置文件，提取域名列表和过期时间阈值。配置文件可以是 YAML、JSON、TOML 等格式。
扩展性: 允许将配置文件替换为数据库读取（后续版本使用 SQLite）。
数据示例
```
domains:
  - name: "example.com"
    notify_before_days: 30
  - name: "anotherdomain.com"
    notify_before_days: 14
notifications:
  email:
    enabled: true
    smtp_server: "smtp.example.com"
    port: 587
    username: "user@example.com"
    password: "password"
    to: ["recipient@example.com"]
  slack:
    enabled: true
    webhook_url: "https://hooks.slack.com/services/..."
  sms:
    enabled: false
    provider: "twilio"
    account_sid: "ACXXXXXXXXXXXXXXXXX"
    auth_token: "your_auth_token"
    from: "+123456789"
    to: ["+987654321"]
```

## 2.证书检查模块（Certificate Checking Module）
功能: 使用 TLS 连接检查每个域名的证书有效期，并计算距离到期的天数。
扩展性: 后续可以引入并发检查、批量操作或外部 API 查询。

## 3.通知模块（Notification Module）
功能: 根据用户配置的通知方式发送提醒（例如通过电子邮件、Slack、短信等）。
接口设计: 定义统一的通知接口，使得不同的通知方式实现统一。
扩展性: 支持新的报警方式时，只需实现通知接口即可。

## 4.调度模块（Scheduler Module）
功能: 定期检查证书的状态。可以使用内置的定时器（如 time.Ticker），也可以在后续版本引入更复杂的任务调度器（如 cron-like 功能）。

## 5.日志模块（Logging Module）
功能: 记录每次检查结果、错误和报警操作，便于调试和排查问题。
扩展性: 后续可以集成日志的持久化存储或外部日志收集服务。


# 项目目录介绍
certwatch/
│
├── cmd/                   # 主程序入口
│   └── certwatch.go        # 项目的启动文件，负责启动应用、读取配置文件、初始化调度器等。
│
├── config/                # 配置文件及相关解析逻辑
│   └── config.go          # 解析配置文件的代码
│
├── pkg/                   # 核心业务逻辑
│   ├── certificate/       # 证书检查逻辑
│   │   └── checker.go     # 证书检查功能的实现。检查证书有效期。使用tls.Dial实现检查功能，并返回证书过期的时间。
│   │
│   ├── notification/      # 通知模块，一个通知接口包含多个实现，如邮件、Slack、短信等。通过实现一个通用接口来支持多种通知方式。
│   │   ├── email.go       # 邮件通知
│   │   ├── slack.go       # Slack 通知
│   │   ├── sms.go         # 短信通知
│   │   └── notifier.go    # 通知接口和相关逻辑
│   │
│   ├── scheduler/         # 任务调度模块
│   │   └── scheduler.go   # 定时器实现及相关逻辑
│   │
│   ├── logging/           # 日志模块
│   │   └── logger.go      # 日志记录的实现
│   │
│   └── db/                # 数据存储层，处理数据库的操作
│       └── sqlite.go      # SQLite 数据库相关操作
│
├── config.yaml            # 初始配置文件
├── LICENSE                # 项目许可证
├── README.md              # 项目说明文件
└── go.mod                 # Go 模块文件

# 模块功能说明
cmd/certwatch.go
这是项目的主入口。负责启动应用、读取配置文件、初始化调度器等。

config/config.go
负责解析配置文件（如 config.yaml）。读取用户输入的域名、通知方式以及到期提醒天数。

pkg/certificate/checker.go
负责检查证书有效期。使用 tls.Dial 实现检查功能，并返回证书过期的时间。

pkg/notification/
这是一个通知接口，包含多个实现，如邮件、Slack、短信等。通过实现一个通用接口来支持多种通知方式。
每个文件（如 email.go）实现对应的通知功能。

pkg/scheduler/scheduler.go
负责定期调度证书检查任务。可以使用简单的 time.Ticker 实现，也可以后续扩展为支持 cron 表达式的调度器。

pkg/logging/logger.go
日志记录功能，记录每次检查、错误信息、通知发送等重要操作。

pkg/db/sqlite.go
虽然第一版是从配置文件读取，但为后续数据库存储留一个模块。在后续版本中，可以将配置和证书检查历史记录存储到 SQLite 数据库中。

# 进一步扩展
持久化存储: 后续可以将域名信息和检查历史记录存储到 SQLite 数据库中。
Web UI: 如果需要一个前端界面，后续可以通过简单的 Golang Web 框架（如 Gin）开发 Web 界面，让用户在浏览器中管理和监控证书。
更多报警方式: 第一版完成后，逐步加入短信、Slack、Webhook 等更多报警方式的支持。