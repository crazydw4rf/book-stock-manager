package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

// use -ldflags flag to change value of these variables
//
// example: -ldflags="-X github.com/crazydw4rf/book-stock-manager/internal/config.APP_VERSION=1.0.4-beta"
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

type Config struct {
	APP_HOST                 string `mapstructure:"APP_HOST"`
	APP_PORT                 int    `mapstructure:"APP_PORT"`
	DB_HOST                  string `mapstructure:"DB_HOST"`
	DB_PORT                  uint16 `mapstructure:"DB_PORT"`
	DB_USER                  string `mapstructure:"DB_USER"`
	DB_PASSWORD              string `mapstructure:"DB_PASSWORD"`
	DB_NAME                  string `mapstructure:"DB_NAME"`
	JWT_ACCESS_TOKEN_SECRET  string `mapstructure:"JWT_ACCESS_TOKEN_SECRET"`
	JWT_REFRESH_TOKEN_SECRET string `mapstructure:"JWT_REFRESH_TOKEN_SECRET"`
}

func InitConfig() (*Config, error) {
	cfg := new(Config)

	// v := viper.NewWithOptions(
	// 	viper.WithLogger(
	// 		slog.New(
	// 			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
	// 				Level: slog.LevelDebug,
	// 			}))))
	v := viper.New()

	v.AddConfigPath(".")
	v.SetConfigType("env")
	v.SetConfigName(".env")
	v.AutomaticEnv()

	err := v.ReadInConfig()
	if err != nil {
		log.Printf("Error reading config file: %v\n", err)
	}

	// TODO: struct validation?
	err = v.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
