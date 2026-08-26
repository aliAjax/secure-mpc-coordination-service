package config

import (
	"errors"
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTPAddr        string
	StateFile       string
	ShutdownTimeout time.Duration
	MaxBodyBytes    int64
	Environment     string
}

func Load() Config {
	c := Config{HTTPAddr: env("MPC_HTTP_ADDR", ":8080"), StateFile: os.Getenv("MPC_STATE_FILE"), ShutdownTimeout: 10 * time.Second, MaxBodyBytes: 1 << 20, Environment: env("MPC_ENV", "development")}
	if v := os.Getenv("MPC_SHUTDOWN_TIMEOUT"); v != "" {
		if n, e := strconv.Atoi(v); e == nil {
			c.ShutdownTimeout = time.Duration(n) * time.Second
		}
	}
	return c
}
func Validate(c Config) error {
	if c.HTTPAddr == "" || c.MaxBodyBytes <= 0 || c.ShutdownTimeout <= 0 {
		return errors.New("invalid configuration")
	}
	return nil
}
func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
