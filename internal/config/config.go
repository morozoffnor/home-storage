package config

type Config struct {
	ListenAddr string
}

func New(listenAddr string) *Config {
	return &Config{
		ListenAddr: listenAddr,
	}
}
