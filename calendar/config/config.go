package config

// Config is main config of app
type Config struct {
	LogFile string `json:"log_file"`
	LogLvl  string `json:"log_level"`

	HTTPListen HTTPListenConfig `json:"http_listen"`
}

// HTTPListenConfig is config of http server
type HTTPListenConfig struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}
