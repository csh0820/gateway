package config

import (
	"log"
	"sync"

	"github.com/spf13/viper"
)

type Config struct {
	GatewayAddress string      `mapstructure:"gateway_address"`
	GatewayMode    string      `mapstructure:"gateway_mode"`
	Etcd           *EtcdConfig `mapstructure:"etcd"`
}

type EtcdConfig struct {
	Endpoints []string `mapstructure:"endpoints"`
}

var (
	once   sync.Once
	config *Config
)

func GetConfig() *Config {
	once.Do(func() {
		config = loadConfig()
	})
	return config
}

func loadConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	setDefaults()

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to decode config into struct: %v", err)
	}

	return config
}

func setDefaults() {
	viper.SetDefault("gateway_address", ":8080")
	viper.SetDefault("gateway_mode", "release")
	// viper.SetDefault("server.read_timeout", 30)
	// viper.SetDefault("server.write_timeout", 30)
	// viper.SetDefault("server.idle_timeout", 120)

	viper.SetDefault("etcd.endpoints", []string{"localhost:2379"})

	// viper.SetDefault("etcd.dial_timeout", 5)
	// viper.SetDefault("etcd.lease_ttl", 10)
	// viper.SetDefault("etcd.namespace", "microservices")

	// viper.SetDefault("discovery.prefix", "/gateway")
	// viper.SetDefault("discovery.watch_interval", 10)
	// viper.SetDefault("discovery.cache_expiration", 30)

	// viper.SetDefault("http_pool.max_idle_conns", 100)
	// viper.SetDefault("http_pool.max_idle_conns_per_host", 10)
	// viper.SetDefault("http_pool.max_conns_per_host", 50)
	// viper.SetDefault("http_pool.idle_conn_timeout", 90)

	// viper.SetDefault("rate_limit.enabled", true)
	// viper.SetDefault("rate_limit.rps", 100)
	// viper.SetDefault("rate_limit.burst", 50)
}
