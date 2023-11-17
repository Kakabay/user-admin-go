package config

import "time"

type Config struct {
	Env string `yaml:"env"`
	Database `yaml:"database"`
	HTTPServer `yaml:"http_server"`
}

type Database struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Sslmode  string `yaml:"sslmode"`
}

type HTTPServer struct {
	Address string `yaml:"address"`
	Timeout time.Duration `yaml:"timeout"`
	Idle_Timeout time.Duration `yaml:"idle_timeout"`
}