package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Env      string     `yaml:"env"`
	Database Database   `mapstructure:"database"`
	Server   HTTPServer `mapstructure:"http_server"`
}

type Database struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSL      string `mapstructure:"ssl"`
}

type HTTPServer struct {
	Address     string        `mapstructure:"address"`
	Timeout     time.Duration `mapstructure:"timeout"`
	IdleTimeout time.Duration `mapstructure:"idle_timeout"`
}

func MustLoad() *Config {
	viper.SetConfigName("local")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./internal/config/")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()
	setupEnvBinds()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Printf("config file not found, using environment variables\n")
		}
		log.Printf("error openning config file: %s\n", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("error reading config file: %s\n", err)
	}

	return &cfg
}

func setupEnvBinds() {
	viper.BindEnv("env", "APP_ENV")
	viper.BindEnv("database.host", "DB_HOST")
	viper.BindEnv("database.port", "DB_PORT")
	viper.BindEnv("database.user", "DB_USER")
	viper.BindEnv("database.password", "DB_PASSWORD")
	viper.BindEnv("database.name", "DB_NAME")
	viper.BindEnv("database.ssl", "DB_SSL")
	viper.BindEnv("http_server.address", "SERVER_ADDRESS")
	viper.BindEnv("http_server.timeout", "SERVER_TIMEOUT")
	viper.BindEnv("http_server.idle_timeout", "SERVER_IDLE_TIMEOUT")
}
