package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
	Server   ServerConfig   `yaml:"server"`
	Payment  PaymentConfig  `yaml:"payment"`
	Captcha  CaptchaConfig  `yaml:"captcha"`
	Wechat   WechatConfig   `yaml:"wechat"`
	Redis    RedisConfig    `yaml:"redis"`
}

type DatabaseConfig struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type JWTConfig struct {
	WebSecret    string `yaml:"web_secret"`
	WebExpire    string `yaml:"web_expire"`
	PluginSecret string `yaml:"plugin_secret"`
	PluginExpire string `yaml:"plugin_expire"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type PaymentConfig struct {
	OrderExpire string `yaml:"order_expire"`
}

type CaptchaConfig struct {
	Expiration string `yaml:"expiration"`
	Length     int    `yaml:"length"`
}

type WechatConfig struct {
	AppID     string `yaml:"app_id"`
	AppSecret string `yaml:"app_secret"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
}

var (
	config      *Config
	configMutex sync.RWMutex
)

func init() {
	var err error
	config, err = Load()
	if err != nil {
		panic(fmt.Sprintf("加载配置文件失败: %v", err))
	}

	go watchConfigFile()
}

func GetConfig() *Config {
	configMutex.RLock()
	defer configMutex.RUnlock()
	return config
}

func Load() (*Config, error) {
	// 获取当前文件的路径
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("无法获取当前文件路径")
	}

	// 构建 config.yaml 的完整路径
	configPath := filepath.Join(filepath.Dir(filename), "config.yaml")

	// 读取 config.yaml 文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 解析 YAML 数据
	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	return &cfg, nil
}

func watchConfigFile() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Printf("创建文件监视器失败: %v\n", err)
		return
	} else {
		fmt.Println("文件监视器创建成功")
	}
	defer watcher.Close()

	// 获取配置文件路径
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Println("无法获取当前文件路径")
		return
	}
	configPath := filepath.Join(filepath.Dir(filename), "config.yaml")

	err = watcher.Add(configPath)
	if err != nil {
		fmt.Printf("添加文件到监视器失败: %v\n", err)
		return
	} else {
		fmt.Println("文件添加到监视器成功")
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				fmt.Println("配置文件已修改，正在重新加载...")
				newConfig, err := Load()
				if err != nil {
					fmt.Printf("重新加载配置文件失败: %v\n", err)
				} else {
					configMutex.Lock()
					config = newConfig
					configMutex.Unlock()
					fmt.Println("配置文件已成功重新加载")
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Printf("文件监视器错误: %v\n", err)
		}
	}
}
