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
		env = "development"
	}

	config = viper.New()
	config.SetConfigType("yaml")
	config.SetConfigName(env)
	config.AddConfigPath("config/") // Hạn chế dùng đường dẫn tương đối khó hiểu
	config.AddConfigPath(".")       // Cho phép tìm file tại thư mục gốc dự án

	// Đọc biến môi trường, chuyển đổi tên (VD: "db.host" → "DB_HOST")
	config.AutomaticEnv()
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Đọc file config nếu có, log phù hợp
	if err := config.ReadInConfig(); err != nil {
		log.Println("⚠️ Không tìm thấy file cấu hình YAML, đang dùng biến môi trường.")
	} else {
		log.Println("✅ Đã load file config:", config.ConfigFileUsed())
	}
}

// GetConfig trả về con trỏ tới đối tượng Viper
func GetConfig() *viper.Viper {
	return config
}
