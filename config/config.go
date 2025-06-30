package config

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

var config *viper.Viper
var BaseApiUrl string

// Init is an exported method that takes the environment starts the viper
// (external lib) and returns the configuration struct.
func init() {

	BaseApiUrl = os.Getenv("BASE_API_URL")
	if BaseApiUrl == "" {
		BaseApiUrl = "https://cms-ideas-app-c77ce0f9e466.herokuapp.com/"
	}

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
