package config

import (
	"errors"
	"github.com/spf13/viper"
)

var filePathIsEmpty = errors.New("the specified yaml file is empty")

func JaegerUrl() string {
	return viper.GetString(VarJaegerUrl)
}

func ServiceDiscoveryUrl() string {
	return viper.GetString(VarServiceDiscoveryUrl)
}

func HttpServerPort() int {
	return viper.GetInt(VarHttpServerPort)
}

func PrometheusMetricsPort() int {
	return viper.GetInt(VarPrometheusMetricsPort)
}


