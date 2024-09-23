package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"time"
)

func main() {
	// 要查询的域名
	domain := "geektool.org:443"

	// 建立与服务器的TLS连接
	conn, err := tls.Dial("tcp", domain, nil)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// 获取对方服务器的证书
	certs := conn.ConnectionState().PeerCertificates

	// 假设我们只关注第一个证书
	if len(certs) > 0 {
		cert := certs[0]

		// 获取证书的有效期
		fmt.Printf("Domain: %s\n", domain)
		fmt.Printf("Valid from: %s\n", cert.NotBefore.Format(time.RFC3339))
		fmt.Printf("Valid until: %s\n", cert.NotAfter.Format(time.RFC3339))

		// 计算剩余有效天数
		daysLeft := cert.NotAfter.Sub(time.Now()).Hours() / 24
		fmt.Printf("Days left: %.0f days\n", daysLeft)
	} else {
		log.Println("No certificates found")
	}

}
