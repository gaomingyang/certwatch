package main

import (
	"fmt"
	"log"
	"time"

	"certwatch/config"
	"certwatch/pkg/certificate"
	"certwatch/pkg/logging"
	"certwatch/pkg/notification"
	"certwatch/pkg/scheduler"

	"github.com/spf13/viper"
)

func main() {
	// 初始化日志模块
	logger := logging.NewLogger()

	// // 记录一般信息
	// logger.Info("Starting CertWatch application...")
	// // 记录警告
	// logger.Warn("SSL certificate for example.com is expiring soon")
	// // 记录错误
	// logger.Error("Failed to check certificate for example.com")

	// 读取配置文件
	if err := config.LoadConfig("."); err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// 解析配置中的域名列表
	domains := viper.GetStringMap("domains")
	if len(domains) == 0 {
		log.Fatal("No domains found in configuration")
	}

	// 初始化调度器
	checkInterval := viper.GetDuration("check_interval") * time.Minute
	s := scheduler.NewScheduler(checkInterval)

	// 为每个域名添加检查任务
	for domain, settings := range domains {
		domainName := fmt.Sprintf("%v", domain)
		notifyBeforeDays := settings.(map[string]interface{})["notify_before_days"].(int)

		// 定义每个域名的证书检查任务
		task := func() {
			logger.Info(fmt.Sprintf("Checking certificate for domain: %s", domainName))

			// 检查证书
			expirationDate, err := certificate.CheckCertificate(domainName)
			if err != nil {
				logger.Error(fmt.Sprintf("Error checking certificate for domain %s: %v", domainName, err))
				return
			}

			// 计算证书剩余天数
			daysLeft := int(expirationDate.Sub(time.Now()).Hours() / 24)
			logger.Info(fmt.Sprintf("Domain %s certificate expires in %d days", domainName, daysLeft))

			// 如果证书剩余天数小于通知天数，则发送报警
			if daysLeft <= notifyBeforeDays {
				logger.Warn(fmt.Sprintf("Domain %s certificate is expiring in %d days. Sending notifications.", domainName, daysLeft))
				sendNotifications(domainName, daysLeft)
			}
		}

		// 将任务加入调度器
		s.AddTask(domainName, task)
	}

	// 启动调度器
	logger.Info("Starting certificate monitoring...")
	s.Start()

	// 阻塞主进程，保持调度器运行
	select {}
}

// sendNotifications 发送通知
func sendNotifications(domain string, daysLeft int) {
	notifiers := []notification.Notifier{}

	// 检查是否启用了邮件通知
	if viper.GetBool("notifications.email.enabled") {
		emailNotifier := notification.NewEmailNotifier(
			viper.GetString("notifications.email.smtp_server"),
			viper.GetInt("notifications.email.port"),
			viper.GetString("notifications.email.username"),
			viper.GetString("notifications.email.password"),
			viper.GetStringSlice("notifications.email.to"),
		)
		notifiers = append(notifiers, emailNotifier)
	}

	// 检查是否启用了 Slack 通知
	if viper.GetBool("notifications.slack.enabled") {
		slackNotifier := notification.NewSlackNotifier(
			viper.GetString("notifications.slack.webhook_url"),
		)
		notifiers = append(notifiers, slackNotifier)
	}

	// 逐个发送通知
	for _, notifier := range notifiers {
		if err := notifier.Notify(domain, daysLeft); err != nil {
			log.Printf("Error sending notification for domain %s: %v", domain, err)
		}
	}
}
