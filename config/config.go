package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// LoadConfig 负责加载配置文件
func LoadConfig(path string) error {
	// 设置配置文件名称（不需要扩展名）
	viper.SetConfigName("config")

	// 设置配置文件路径
	viper.AddConfigPath(path)              // 指定配置文件路径
	viper.AddConfigPath(".")               // 当前目录
	viper.AddConfigPath("/etc/certwatch/") // 全局目录

	// 设置配置文件类型为 yaml
	viper.SetConfigType("yaml")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("Error reading config file: %w", err)
	}

	// 验证配置文件中的域名是否存在
	if len(viper.GetStringMap("domains")) == 0 {
		return fmt.Errorf("No domains found in config file")
	}

	log.Println("Configuration loaded successfully")
	return nil
}

// GetDomains 返回配置中的域名列表
func GetDomains() map[string]interface{} {
	return viper.GetStringMap("domains")
}

// GetNotificationSettings 返回通知配置
func GetNotificationSettings() map[string]interface{} {
	return viper.GetStringMap("notifications")
}
