package config

import (
	"errors"
	"fmt"
	"log"
	"os"
)

type CommonConfig struct {
	HttpServerPort        int `yaml:"HTTP_SERVER_PORT"`
	PrometheusMetricsPort int `yaml:"PROMETHEUS_METRICS_PORT"`
}

func NewCommonConfig() (*CommonConfig, error) {
	log.Println("beginning the creation of the base config object.")
	viper, err := SetViperConfig(os.Getenv(AppConfigFile))
	if err != nil {
		return nil, err
	}
	httpServerPort := viper.GetInt(VarHttpServerPort)
	if httpServerPort == 0 {
		return nil, errors.New(fmt.Sprintf("%s is not provided", VarHttpServerPort))
	}

	prometheusMetricsPort := viper.GetInt("prometheus_metrics_port")
	if prometheusMetricsPort == 0 {
		return nil, errors.New(fmt.Sprintf("%s is not provided", VarPrometheusMetricsPort))
	}

	return &CommonConfig{
		HttpServerPort:        httpServerPort,
		PrometheusMetricsPort: prometheusMetricsPort,
	}, nil
}
