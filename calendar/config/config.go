package config

// Config is main config of app
type Config struct {
	LogFile string `json:"log_file"`
	LogLvl  string `json:"log_level"`

	HttpListen HttpListenConfig `json:"http_listen"`
}

// HttpListenConfig is config of http server
type HttpListenConfig struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}
