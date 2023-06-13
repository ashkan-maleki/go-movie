package config

import (
	"errors"
	"github.com/spf13/viper"
	"log"
	"os"
)

type CommonConfig struct {
	HttpServerPort        int `yaml:"http_server_port"`
	PrometheusMetricsPort int `yaml:"prometheus_metrics_port"`
}

func NewCommonConfig() (*CommonConfig, error) {
	log.Println("beginning the creation of the base config object.")
	err := SetViperConfig(os.Getenv(AppConfigFile))
	if err != nil {
		return nil, err
	}
	httpServerPort := viper.GetInt("http_server_port")
	if httpServerPort == 0 {
		return nil, errors.New("http_server_port is not provided")
	}

	prometheusMetricsPort := viper.GetInt("prometheus_metrics_port")
	if prometheusMetricsPort == 0 {
		return nil, errors.New("prometheus_metrics_port is not provided")
	}

	return &CommonConfig{
		HttpServerPort:        httpServerPort,
		PrometheusMetricsPort: prometheusMetricsPort,
	}, nil
}
