package config

import (
	"flag"
	"strings"

	"log"
)

type Config struct {
	Port         int
	MaxKeySize   int
	MaxValueSize int
	LogLevel     string
}

func LoadConfig() Config {
	cfg := Config{}

	flag.IntVar(&cfg.Port, "port", 8080, "server port")
	flag.IntVar(&cfg.MaxKeySize, "max-key-size", 256, "maximum key size as bytes")
	flag.IntVar(&cfg.MaxValueSize, "max-value-size", 1048576, "maximum value size as bytes")
	flag.StringVar(&cfg.LogLevel, "log-level", "debug", "log level: debug/info/warn/error")

	flag.Parse()

	validLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}

	cfg.LogLevel = strings.ToLower(cfg.LogLevel)

	if !validLevels[cfg.LogLevel] {
		log.Fatalf("invalid log level: %s", cfg.LogLevel)
	}

	return cfg
}
