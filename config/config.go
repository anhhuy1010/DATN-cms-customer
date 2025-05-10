package config

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

var config *viper.Viper

func init() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "production" // Heroku default
	}

	config = viper.New()
	config.SetConfigType("yaml")
	config.SetConfigName(env)
	config.AddConfigPath("./config/")
	config.AddConfigPath("../config/")

	// Tự động đọc từ biến môi trường
	config.AutomaticEnv()
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Load file YAML nếu có
	if err := config.ReadInConfig(); err != nil {
		log.Println("⚠️ Không tìm thấy file cấu hình YAML. Sử dụng biến môi trường.")
	} else {
		log.Println("✅ Đã load file config:", config.ConfigFileUsed())
	}
}

func GetConfig() *viper.Viper {
	return config
}
