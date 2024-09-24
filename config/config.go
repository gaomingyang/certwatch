package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

// LoadConfig 读取配置文件
func LoadConfig(path string) error {
	// 设置配置文件名称（不需要扩展名）
	viper.SetConfigName("config")
	path = strings.TrimSpace(path)
	// 设置自定义配置文件路径
	if path != "" && path != "." && path != "/etc/certwatch/" {
		viper.AddConfigPath(path) // 指定配置文件路径
	}
	viper.AddConfigPath(".")               // 当前目录
	viper.AddConfigPath("/etc/certwatch/") // 全局目录

	viper.SetConfigType("yaml") // 设置配置文件类型为 yaml

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("error reading config file: %s", err.Error())
	}

	// fmt.Printf("Loaded Config: %v\n", viper.AllSettings()) // 调试信息

	log.Println("Configuration loaded successfully")
	return nil
}

// GetDomains 返回配置中的域名列表
// func GetDomains() map[string]interface{} {
// 	return viper.GetStringMap("domains")
// }
