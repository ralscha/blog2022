package config

import (
	"github.com/spf13/viper"
	"time"
)

type Environment string

const (
	Production  Environment = "production"
	Development Environment = "development"
)

type Config struct {
	Environment Environment
	Session     struct {
		SecureCookie bool
		CookieDomain string
		Lifetime     time.Duration
	}
	HTTP struct {
		Addr                  string
		ReadTimeoutInSeconds  int64
		WriteTimeoutInSeconds int64
		IdleTimeoutInSeconds  int64
	}
	Passwordless struct {
		ApiUrl       string
		SecretApiKey string
	}
}

func applyDefaults() {
	viper.SetDefault("environment", Production)
	viper.SetDefault("http.readTimeoutInSeconds", 30)
	viper.SetDefault("http.writeTimeoutInSeconds", 30)
	viper.SetDefault("http.idleTimeoutInSeconds", 120)
	viper.SetDefault("db.maxOpenConns", 4)
	viper.SetDefault("db.maxIdleConns", 2)
	viper.SetDefault("db.maxIdleTime", "15m")
	viper.SetDefault("db.maxLifetime", "2h")
	viper.SetDefault("session.secureCookie", true)
	viper.SetDefault("session.lifetime", "24h")
}

func LoadConfig() (Config, error) {
	var cfg Config

	applyDefaults()
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		return cfg, err
	}

	viper.SetEnvPrefix("WEBAUTHN")
	viper.AutomaticEnv()

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
