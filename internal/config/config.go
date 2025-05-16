package config

import (
	"log"
	"reflect"
	"time"

	"github.com/spf13/viper"
)

// use -ldflags flag to change value of these variables
//
// example: -ldflags="-X github.com/crazydw4rf/book-stock-manager/internal/config.APP_VERSION=1.0.4-beta"
var (
	APP_ENV          = "development"
	APP_VERSION      = "0.0.1-alpha"
	API_DOCS_ENABLED = true
)

const (
	BASE_API_HTTP_PATH            = "/api/v1"
	SESSION_COOKIE_NAME           = "__Host_session_"
	REFRESH_TOKEN_COOKIE_NAME     = "__Host_refreshtoken_"
	ACCESS_TOKEN_COOKIE_NAME      = "__Host_token_"
	ACCESS_TOKEN_HEADER_NAME      = "Authorization"
	CSRF_HEADER_NAME              = "X-Csrf-Token"
	CSRF_COOKIE_NAME              = "__Host_csrf_"
	ACCESS_TOKEN_EXPIRATION_TIME  = time.Minute * 15
	REFRESH_TOKEN_EXPIRATION_TIME = (time.Hour * 24) * 7
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
	v := viper.New()

	v.AddConfigPath(".")
	v.SetConfigType("env")
	v.SetConfigName(".env")
	v.AutomaticEnv()

	bindEnvStruct(v, cfg)

	err := v.ReadInConfig()
	if err != nil {
		log.Printf("Warning reading config file: %v\n", err)
	}

	// TODO: struct validation?
	err = v.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func bindEnvStruct(v *viper.Viper, s any) {
	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return
	}

	typ := val.Type()

	for i := range typ.NumField() {
		field := typ.Field(i)
		tagValue := field.Tag.Get("mapstructure")
		if tagValue != "" {
			v.BindEnv(field.Name, tagValue)
		}
	}
}
