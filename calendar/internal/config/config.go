package config

import (
	"encoding/json"
	"os"
	"time"
)

// Config is main config of app
type Config struct {
	LogFile string `json:"log_file"`
	LogLvl  string `json:"log_level"`

	HTTPListen   HTTPListenConfig `json:"http_listen"`
	GRPC         GRPCConfig       `json:"grpc"`
	DB           DBConfig         `json:"db"`
	RabbitConfig RabbitConfig     `json:"rabbitmq"`

	ScheduleDurationString string        `json:"schedule_duration"`
	ScheduleDuration       time.Duration `json:"-"`
}

// HTTPListenConfig is config of http server
type HTTPListenConfig struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

// GRPCConfig config
type GRPCConfig struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

// DBConfig is config for db connect
type DBConfig struct {
	DSN             string `json:"dsn"`
	MaxConnections  int    `json:"max_connections"`
	IdleConnections int    `json:"idle_connections"`
}

// RabbitConfig is config for rabbitmq connect
type RabbitConfig struct {
	DSN      string `json:"dsn"`
	Exchange string `json:"exchange"`
	Queue    string `json:"queue"`
}

// New - create new config from file path
func New(configPath string) (*Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	decoder := json.NewDecoder(file)

	cfg := &Config{}
	err = decoder.Decode(cfg)
	if err != nil {
		return nil, err
	}

	cfg.ScheduleDuration, err = time.ParseDuration(cfg.ScheduleDurationString)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
