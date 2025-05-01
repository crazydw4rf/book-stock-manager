package config

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

// use -ldflags flag to change value of these variables
//
// example: -ldflags="-X github.com/crazydw4rf/magic-url/internal/config.APP_VERSION=1.0.4-beta"
var (
	APP_ENV     = "development"
	APP_VERSION = "0.0.1-alpha"
)

const (
	BaseApiHttpPath            = "/api/v1"
	SessionCookieName          = "__Host_session_"
	RefreshTokenCookieName     = "__Host_refreshtoken_"
	AccessTokenCookieName      = "__Host_token_"
	AccessTokenHeaderName      = "Authorization"
	CSRFHeaderName             = "X-Csrf-Token"
	CSRFCookieName             = "__Host_csrf_"
	AccessTokenExpirationTime  = time.Minute * 15
	RefreshTokenExpirationTime = (time.Hour * 24) * 7
)

type databaseConfig struct {
	Host     string
	Port     uint16
	User     string
	Password string
	DBName   string `mapstructure:"db_name"`
}

type redisConfig struct {
	Host string
	Port int
}

type jwtConfig struct {
	AccessTokenSecret  string `mapstructure:"access_token_secret"`
	RefreshTokenSecret string `mapstructure:"refresh_token_secret"`
}

type appConfig struct {
	Host string
	Port int
}

type Config struct {
	App      appConfig      `mapstructure:"app"`
	Database databaseConfig `mapstructure:"database"`
	JWT      jwtConfig      `mapstructure:"jwt"`
	Redis    redisConfig    `mapstructure:"redis"`
}

func InitConfig() (*Config, error) {
	cfg := new(Config)

	v := viper.New()

	v.AddConfigPath(".")
	v.SetConfigType("toml")
	v.SetConfigName("config")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	_ = v.ReadInConfig()

	// TODO: struct validation?
	err := v.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
