package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// SlackNotifier Slack通知器结构体
type SlackNotifier struct {
	WebhookURL string
}

// NewSlackNotifier 创建并返回一个新的 SlackNotifier 实例
func NewSlackNotifier(webhookURL string) *SlackNotifier {
	return &SlackNotifier{
		WebhookURL: webhookURL,
	}
}

// Notify 发送证书到期通知到 Slack
func (s *SlackNotifier) Notify(domain string, daysLeft int) error {
	// 构建Slack消息内容
	message := fmt.Sprintf("🔔 *Certificate Expiration Warning* 🔔\n\nThe SSL/TLS certificate for domain `%s` is expiring in *%d days*.\nPlease take action to renew the certificate to avoid any downtime.", domain, daysLeft)

	payload := map[string]string{
		"text": message,
	}

	// 将消息转换为 JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal Slack payload: %v", err)
	}

	// 发送POST请求到Slack Webhook URL
	resp, err := http.Post(s.WebhookURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to send notification to Slack: %v", err)
	}
	defer resp.Body.Close()

	// 检查Slack返回的响应状态
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("non-OK response returned from Slack: %s", resp.Status)
	}

	fmt.Printf("Slack notification sent successfully for domain %s\n", domain)
	return nil
}
