package config

import (
	"errors"
	"fmt"
	"log"
	"os"
)

type CommonConfig struct {
	HttpServerPort        int `mapstructure:"HTTP_SERVER_PORT" validate:"required"`
	PrometheusMetricsPort int `mapstructure:"PROMETHEUS_METRICS_PORT" validate:"required"`
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

	return &config, nil
}
