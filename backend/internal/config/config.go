package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBName     string `mapstructure:"DB_NAME"`

	JWTSecret    string `mapstructure:"JWT_SECRET"`
	JWTExpiredIn string `mapstructure:"JWT_EXPIRED_IN"`

	SMTPHost     string `mapstructure:"SMTP_HOST"`
	SMTPPort     int    `mapstructure:"SMTP_PORT"`
	SMTPEmail    string `mapstructure:"SMTP_EMAIL"`
	SMTPPassword string `mapstructure:"SMTP_PASSWORD"`

	Port    string `mapstructure:"PORT"`
	GinMode string `mapstructure:"GIN_MODE"`
}

func LoadConfig() (config *Config, err error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatal("Could not load config file:", err)
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatal("Could not unmarshal config:", err)
		return
	}

	return
}
