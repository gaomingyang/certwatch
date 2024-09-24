package certificate

import (
	"crypto/tls"
	"fmt"
	"time"
)

// CheckCertificate 检查给定域名的证书到期时间
func CheckCertificate(domain string) (time.Time, error) {
	// 连接到指定域名的443端口（HTTPS 默认端口）
	conn, err := tls.Dial("tcp", domain+":443", nil)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to connect to %s: %v", domain, err)
	}
	defer conn.Close()

	// 获取对方服务器的证书链
	certs := conn.ConnectionState().PeerCertificates

	// 假设我们只关心第一个证书
	if len(certs) > 0 {
		cert := certs[0]
		return cert.NotAfter, nil
	}

	// 如果没有找到证书
	return time.Time{}, fmt.Errorf("no certificates found for domain %s", domain)
}
