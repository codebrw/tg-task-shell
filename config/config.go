package config

type Config struct {
	TG_API_TOKEN string
}

func NewConfig(token string) *Config {
	return &Config{
		TG_API_TOKEN: token,
	}
}