package config

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
)

var (
	cfgCache map[string]*Config
	once     sync.Once
)

// FTPConfig defines the configuration for FTP connections.
type FTPConfig struct {
	Host             string `yaml:"host"`
	User             string `yaml:"user"`
	Password         string `yaml:"password"`
	RemotePathSend   string `yaml:"remote_path_send"`
	RemotePathResult string `yaml:"remote_path_result"`
	LocalPath        string `yaml:"local_path"`
}

// APIEndpoints defines the endpoints for external APIs.
type APIEndpoints struct {
	GetAlertSetting  string        `yaml:"get_alert_setting"`
	SendNotification string        `yaml:"send_notification"`
	Timeout          time.Duration `yaml:"timeout"`
}

// ScheduleConfig defines the schedule for batch jobs.
type ScheduleConfig struct {
	SendTime   string `yaml:"send_time"`
	ResultTime string `yaml:"result_time"`
}

// BatchConfig defines the configuration for a specific batch (Spending Alert or e-NCB).
type BatchConfig struct {
	FTP          FTPConfig      `yaml:"ftp"`
	Schedule     ScheduleConfig `yaml:"schedule"`
	ResultPrefix string         `yaml:"result_file_prefix"`
}

// Config holds the entire application configuration.
type Config struct {
	Environment   string       `yaml:"environment"`
	APIEndpoints  APIEndpoints `yaml:"api_endpoints"`
	SpendingAlert BatchConfig  `yaml:"spending_alert"`
	ENCB          BatchConfig  `yaml:"e_ncb"`
	LogPath       string       `yaml:"log_path"`
	APILogPrefix  string       `yaml:"api_log_prefix"`
}

// LoadConfig loads the configuration from the specified environment's YAML file.
func LoadConfig(env string) map[string]*Config {
	once.Do(func() {
		cfgCache = make(map[string]*Config)

		loadConfigFromFile(env)
	})
	return cfgCache
}

func loadConfigFromFile(env string) {
	cfgCache = make(map[string]*Config)
	filename := fmt.Sprintf("config/%s.yaml", env)
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to read config file '%s': %v", filename, err)
		return
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Failed to unmarshal config file '%s': %v", filename, err)
		return
	}

	cfgCache["default"] = &config
	cfgCache["spending_alert"] = &Config{
		Environment:   config.Environment,
		APIEndpoints:  config.APIEndpoints,
		SpendingAlert: config.SpendingAlert,
		LogPath:       config.LogPath,
		APILogPrefix:  config.APILogPrefix,
	}
	cfgCache["encb"] = &Config{
		Environment:  config.Environment,
		APIEndpoints: config.APIEndpoints,
		ENCB:         config.ENCB,
		LogPath:      config.LogPath,
		APILogPrefix: config.APILogPrefix,
	}
}

// GetConfig returns the configuration for the specified business from the cache.
func GetConfig(business string) (*Config, bool) {
	configMap := LoadConfig(GetEnv())
	conf, ok := configMap[business]
	return conf, ok
}

// GetEnv returns the current application environment (defaults to "dev" if APP_ENV is not set).
func GetEnv() string {
	env := os.Getenv("APP_ENV")
	if env == "" {
		return "dev"
	}
	return env
}
