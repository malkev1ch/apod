package config

import "github.com/caarlos0/env/v6"

// Config is struct that contains information about service configuration options.
type Config struct {
	ServerListenAddress string `env:"SERVER_LISTEN_ADDRESS,notEmpty"`
	PGConnectionString  string `env:"PG_CONN_STRING,notEmpty"`
	NASAAPIKey          string `env:"NASA_API_KEY,notEmpty"`
	NASAPIAddress       string `env:"NASA_API_ADDRESS,notEmpty"`
	LogLevel            string `env:"LOG_LEVEL" envDefault:"info"`
	DevMode             bool   `env:"DEV_MODE" envDefault:"false"`
	S3Address           string `env:"S3_ADDRESS,notEmpty"`
	S3AccessKey         string `env:"S3_ACCESS_KEY,notEmpty"`
	S3SecretKey         string `env:"S3_SECRET_KEY,notEmpty"`
	S3Endpoint          string `env:"S3_ENDPOINT,notEmpty"`
	S3Secured           bool   `env:"S3_SECURED,notEmpty"`
	S3Bucket            string `env:"S3_BUCKET,notEmpty"`
}

// New creates a new Config.
func New() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)

	return cfg, err
}
