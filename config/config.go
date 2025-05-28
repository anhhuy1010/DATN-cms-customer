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
	config.AddConfigPath("./config/")
	config.AddConfigPath("../config/")
	config.AutomaticEnv()
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := config.ReadInConfig(); err != nil {
		log.Fatal(err.Error())
	}
}
func GetConfig() *viper.Viper {
	return config
}
func CustomerServiceCheckTokenURL() string {
	baseURL := config.GetString("customer_service.base_url")
	path := config.GetString("customer_service.check_token_path")
	return baseURL + path
}
