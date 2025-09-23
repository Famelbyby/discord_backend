package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env string `yaml:"env" env-default:"local"`

	GRPC           GRPCConfig `yaml:"grpc"`
	MigrationsPath string
	TokenTTL       time.Duration `yaml:"token_ttl" env-default:"1h"`
	Clients        ClientsConfig `yaml:"clients"`
	AppSecret      string        `yaml:"app_secret" env-required:"true" env:"APP_SECRET"`
}
type Client struct {
	Address      string        `yaml:"address"`
	Timeout      time.Duration `yaml:"timeout"`
	RetriesCount int           `yaml:"retriesCount"`
	Insecure     string        `yaml:"incesure"`
}
type ClientsConfig struct {
	Auth      Client `yaml:"auth"`
	Relations Client `yaml:"relations"`
}
type GRPCConfig struct {
	AuthMS      MicroserviceGRPCConfig `yaml:"auth"`
	RelationsMS MicroserviceGRPCConfig `yaml:"relations"`
}
type MicroserviceGRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	configPath, _ := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadPath(configPath)
}

func MustLoadPath(configPath string) *Config {
	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	return &cfg
}

// fetchConfigPath fetches config path from command line flag or environment variable.
// Priority: flag > env > default.
// Default value is empty string.
func fetchConfigPath() (string, error) {
	const configFile = ".env"

	err := godotenv.Load(configFile)
	if err != nil {
		return "", err
	}

	cfgPath, exists := os.LookupEnv("CONFIG_PATH")
	if !exists {
		return "", fmt.Errorf("config wasn't found")
	}

	return cfgPath, nil
}
