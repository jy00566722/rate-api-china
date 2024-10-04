package config

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Name      string `mapstructure:"name"`
		Port      int    `mapstructure:"port"`
		JWTSecret string `mapstructure:"jwt_secret"`
	} `mapstructure:"app"`

	MySQL struct {
		Host            string        `mapstructure:"host"`
		Port            int           `mapstructure:"port"`
		Database        string        `mapstructure:"database"`
		Username        string        `mapstructure:"username"`
		Password        string        `mapstructure:"password"`
		MaxIdleConns    int           `mapstructure:"max_idle_conns"`
		MaxOpenConns    int           `mapstructure:"max_open_conns"`
		ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	} `mapstructure:"mysql"`
}

var (
	config Config
	once   sync.Once
)

// Load loads configuration from file and environment variables
func Load() (*Config, error) {
	var err error
	once.Do(func() {
		err = loadConfig()
	})
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func loadConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	// 读取环境变量覆盖
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return fmt.Errorf("config file not found: %w", err)
		}
		return fmt.Errorf("failed to read config file: %w", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 从环境变量覆盖敏感信息
	if envDBPass := os.Getenv("MYSQL_PASSWORD"); envDBPass != "" {
		config.MySQL.Password = envDBPass
	}
	if envJWTSecret := os.Getenv("JWT_SECRET"); envJWTSecret != "" {
		config.App.JWTSecret = envJWTSecret
	}

	return nil
}
