package config

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"log"
	"os"
)

type CommonConfig struct {
	HttpServerPort        int `mapstructure:"HTTP_SERVER_PORT" structs:"HTTP_SERVER_PORT" env:"HTTP_SERVER_PORT"`
	PrometheusMetricsPort int `mapstructure:"PROMETHEUS_METRICS_PORT"`
}

func NewCommonConfig() (*CommonConfig, error) {
	log.Println("beginning the creation of the base config object.")
	viper, err := SetViperConfig(os.Getenv(AppConfigFile))
	if err != nil {
		return nil, err
	}

	var config CommonConfig
	if err := viper.Unmarshal(&config); err != nil {
		return nil, errors.New(fmt.Sprintf("unable to unmarshall the config %v", err))
	}
	log.Println(config)
	validate := validator.New()
	if err := validate.Struct(&config); err != nil {
		return nil, errors.New(fmt.Sprintf("Missing required attributes %v\n", err))
	}

	return &config, nil
}
