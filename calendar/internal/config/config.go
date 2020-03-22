package config

// Config is main config of app
type Config struct {
	LogFile string `json:"log_file"`
	LogLvl  string `json:"log_level"`

	HTTPListen HTTPListenConfig `json:"http_listen"`
	GRPC       GRPCConfig       `json:"grpc"`
	DB         DBConfig         `json:"db"`
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
