package notification

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// SMSNotifier 短信通知器结构体
type SMSNotifier struct {
	AccountSID string
	AuthToken  string
	FromNumber string
	ToNumbers  []string
	APIURL     string
}

// NewSMSNotifier 创建并返回一个新的 SMSNotifier 实例
func NewSMSNotifier(accountSID, authToken, fromNumber string, toNumbers []string) *SMSNotifier {
	return &SMSNotifier{
		AccountSID: accountSID,
		AuthToken:  authToken,
		FromNumber: fromNumber,
		ToNumbers:  toNumbers,
		APIURL:     "https://api.twilio.com/2010-04-01/Accounts/" + accountSID + "/Messages.json",
	}
}

// Notify 发送证书到期通知短信
func (s *SMSNotifier) Notify(domain string, daysLeft int) error {
	message := fmt.Sprintf("Alert: The SSL/TLS certificate for domain %s is expiring in %d days. Please renew it soon to avoid downtime.", domain, daysLeft)

	// 循环给每个号码发送通知
	for _, toNumber := range s.ToNumbers {
		err := s.sendSMS(toNumber, message)
		if err != nil {
			return fmt.Errorf("failed to send SMS to %s: %v", toNumber, err)
		}
	}

	fmt.Printf("SMS notification sent successfully for domain %s\n", domain)
	return nil
}

// sendSMS 发送单条短信
func (s *SMSNotifier) sendSMS(to, message string) error {
	// 准备短信内容
	data := url.Values{}
	data.Set("To", to)
	data.Set("From", s.FromNumber)
	data.Set("Body", message)

	// 创建 HTTP 请求
	client := &http.Client{}
	req, err := http.NewRequest("POST", s.APIURL, strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}
	req.SetBasicAuth(s.AccountSID, s.AuthToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send SMS via Twilio: %v", err)
	}
	defer resp.Body.Close()

	// 检查 Twilio API 返回的状态码
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("non-OK response returned from Twilio: %s", resp.Status)
	}

	return nil
}
