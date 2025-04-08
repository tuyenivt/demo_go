package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	RedisShards int      `mapstructure:"redis_shards"`
	RedisHosts  []string `mapstructure:"redis_hosts"`
	APIPort     int      `mapstructure:"api_port"`
	RateLimit   int      `mapstructure:"rate_limit"`
	LogLevel    string   `mapstructure:"log_level"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Set defaults
	viper.SetDefault("redis_shards", 3)
	viper.SetDefault("redis_hosts", "localhost:6379,localhost:6380,localhost:6381")
	viper.SetDefault("api_port", 8080)
	viper.SetDefault("rate_limit", 100)
	viper.SetDefault("log_level", "info")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	config := &Config{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Validate configuration
	if err := config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) validate() error {
	if c.RedisShards <= 0 {
		return fmt.Errorf("redis_shards must be greater than 0")
	}

	if len(c.RedisHosts) != c.RedisShards {
		return fmt.Errorf("number of redis_hosts must match redis_shards")
	}

	if c.APIPort <= 0 {
		return fmt.Errorf("api_port must be greater than 0")
	}

	if c.RateLimit <= 0 {
		return fmt.Errorf("rate_limit must be greater than 0")
	}

	return nil
}

// GetRedisHosts returns the Redis hosts as a slice of strings
func GetRedisHosts() []string {
	hosts := os.Getenv("REDIS_HOSTS")
	if hosts == "" {
		return []string{"localhost:6379", "localhost:6380", "localhost:6381"}
	}
	return strings.Split(hosts, ",")
}

// GetRedisShards returns the number of Redis shards
func GetRedisShards() int {
	shards, err := strconv.Atoi(os.Getenv("REDIS_SHARDS"))
	if err != nil {
		return 3
	}
	return shards
}

// GetAPIPort returns the API port
func GetAPIPort() int {
	port, err := strconv.Atoi(os.Getenv("API_PORT"))
	if err != nil {
		return 8080
	}
	return port
}

// GetRateLimit returns the rate limit per second
func GetRateLimit() int {
	limit, err := strconv.Atoi(os.Getenv("RATE_LIMIT"))
	if err != nil {
		return 100
	}
	return limit
}

// GetLogLevel returns the logging level
func GetLogLevel() string {
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		return "info"
	}
	return level
}
