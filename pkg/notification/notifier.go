package notification

// Notifier 通用通知接口
type Notifier interface {
	// Notify 用于发送通知，接收域名和证书到期剩余天数
	Notify(domain string, daysLeft int) error
}
