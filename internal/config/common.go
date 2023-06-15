package config

import (
	"log"
	"os"
)

type CommonConfig struct {
	HttpServerPort        int `mapstructure:"HTTP_SERVER_PORT" validate:"required"`
	PrometheusMetricsPort int `mapstructure:"PROMETHEUS_METRICS_PORT" validate:"required"`
}

func NewCommonConfig() (*CommonConfig, error) {
	log.Println("beginning the creation of the base config object.")
	var config CommonConfig
	err := Configure(os.Getenv(AppConfigFile), &config)
	if err != nil {
		return nil, err
	}
	log.Println("ending the creation of the base config object.")
	return &config, nil
}
