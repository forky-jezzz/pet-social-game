package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	JWT      JWTConfig      `yaml:"jwt"`
	MySQL    MySQLConfig    `yaml:"mysql"`
	Redis    RedisConfig    `yaml:"redis"`
	Security SecurityConfig `yaml:"security"`
	Auth     AuthConfig     `yaml:"auth"`
}

type ServerConfig struct {
	Addr string `yaml:"addr"`
	Mode string `yaml:"mode"`
}

type JWTConfig struct {
	Secret      string `yaml:"secret"`
	ExpireHours int    `yaml:"expire_hours"`
}

type MySQLConfig struct {
	DSN string `yaml:"dsn"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type SecurityConfig struct {
	Provider string       `yaml:"provider"`
	Alipay   AlipayConfig `yaml:"alipay"`
}

type AlipayConfig struct {
	AppID           string `yaml:"app_id"`
	PrivateKey      string `yaml:"private_key"`
	AlipayPublicKey string `yaml:"alipay_public_key"`
	Gateway         string `yaml:"gateway"`
}

type AuthConfig struct {
	AllowDevLogin bool `yaml:"allow_dev_login"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	if cfg.Server.Addr == "" {
		cfg.Server.Addr = ":8080"
	}
	if cfg.JWT.ExpireHours <= 0 {
		cfg.JWT.ExpireHours = 168
	}
	if cfg.Security.Provider == "" {
		cfg.Security.Provider = "mock"
	}
	return &cfg, nil
}
