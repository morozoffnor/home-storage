package config

import "os"

type Config struct {
	ListenAddr   string
	DatabaseAddr string
}

func New() *Config {
	c := &Config{}
	c.loadFromEnv()
	return c
}

func (c *Config) loadFromEnv() {
	sa := os.Getenv("HOME_STORAGE_LISTEN_ADDR")
	if sa != "" {
		c.ListenAddr = sa
	}

	db := os.Getenv("POSTGRES_STRING")
	if db != "" {
		c.DatabaseAddr = db
	}
}
