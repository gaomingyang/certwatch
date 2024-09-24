package notification

import (
	"fmt"
	"net/smtp"
	"strings"
)

// EmailNotifier 邮件通知器结构体
type EmailNotifier struct {
	SMTPServer string
	Port       int
	Username   string
	Password   string
	Recipients []string
}

// NewEmailNotifier 创建并返回一个新的 EmailNotifier 实例
func NewEmailNotifier(smtpServer string, port int, username, password string, recipients []string) *EmailNotifier {
	return &EmailNotifier{
		SMTPServer: smtpServer,
		Port:       port,
		Username:   username,
		Password:   password,
		Recipients: recipients,
	}
}

// Notify 发送证书到期通知邮件
func (e *EmailNotifier) Notify(domain string, daysLeft int) error {
	// 构造邮件内容
	subject := fmt.Sprintf("Certificate Expiration Warning for %s", domain)
	body := fmt.Sprintf("The SSL/TLS certificate for domain %s is expiring in %d days. Please renew it soon to avoid downtime.", domain, daysLeft)
	message := e.buildMessage(subject, body)

	// 连接到 SMTP 服务器并发送邮件
	auth := smtp.PlainAuth("", e.Username, e.Password, e.SMTPServer)
	smtpAddr := fmt.Sprintf("%s:%d", e.SMTPServer, e.Port)

	err := smtp.SendMail(smtpAddr, auth, e.Username, e.Recipients, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email notification: %v", err)
	}

	fmt.Printf("Email notification sent to %v for domain %s\n", e.Recipients, domain)
	return nil
}

// buildMessage 构建邮件的 MIME 消息格式
func (e *EmailNotifier) buildMessage(subject, body string) string {
	headers := make(map[string]string)
	headers["From"] = e.Username
	headers["To"] = strings.Join(e.Recipients, ",")
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/plain; charset=UTF-8"

	message := ""
	for key, value := range headers {
		message += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	message += "\r\n" + body

	return message
}
