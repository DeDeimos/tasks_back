package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config Структура конфигурации;
// Содержит все конфигурационные данные о сервисе;
// автоподгружается при изменении исходного файла
type Config struct {
	ServiceHost string
	ServicePort int

	Redis RedisConfig
}

type RedisConfig struct {
	Host        string
	Password    string
	Port        int
	User        string
	DialTimeout time.Duration
	ReadTimeout time.Duration
}

const (
	envRedisHost = "REDIS_HOST"
	envRedisPort = "REDIS_PORT"
	envRedisUser = "REDIS_USER"
	envRedisPass = "REDIS_PASSWORD"
)

type TokenConfig struct {
	JWT struct {
		SigningMethod jwt.SigningMethod `json:"signing_method"`
		ExpiresIn     time.Duration     `json:"expires_in"`
		Token         string            `json:"token"`
		Role          string            `json:"role"`
	} `json:"jwt"`
	Redis RedisConfig
}

func New() (*TokenConfig, error) {

	return &TokenConfig{
		JWT: struct {
			SigningMethod jwt.SigningMethod `json:"signing_method"`
			ExpiresIn     time.Duration     `json:"expires_in"`
			Token         string            `json:"token"`
			Role          string            `json:"role"`
		}{
			SigningMethod: jwt.SigningMethodHS256,
			ExpiresIn:     time.Hour * 24,
			Token:         "",
			Role:          "",
		},
	}, nil
}

// NewConfig Создаёт новый объект конфигурации, загружая данные из файла конфигурации
func NewConfig() (*Config, error) {
	var err error

	configName := "config"
	_ = godotenv.Load()
	if os.Getenv("CONFIG_NAME") != "" {
		configName = os.Getenv("CONFIG_NAME")
	}

	viper.SetConfigName(configName)
	viper.SetConfigType("toml")
	viper.AddConfigPath("config")
	viper.AddConfigPath(".")
	viper.WatchConfig()

	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = viper.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}

	cfg.Redis.Host = os.Getenv(envRedisHost)
	cfg.Redis.Port, err = strconv.Atoi(os.Getenv(envRedisPort))
	log.Println("Host")
	log.Println(cfg.Redis.Host)
	log.Println("Port")
	log.Println(cfg.Redis.Port)
	if err != nil {
		return nil, fmt.Errorf("redis must be")
	}

	cfg.Redis.Password = os.Getenv(envRedisPass)
	cfg.Redis.User = os.Getenv(envRedisUser)

	log.Info("config parsed")

	return cfg, nil
}
