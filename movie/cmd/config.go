package main

type jaegerConfig struct {
	URL string `yaml:"url"`
}

type apiConfig struct {
	Port int `yaml:"port"`
}

type prometheusConfig struct {
	MetricsPort int `yaml:"metricsPort"`
}

type config struct {
	API        apiConfig        `yaml:"api"`
	Jaeger     jaegerConfig     `yaml:"jaeger"`
	Prometheus prometheusConfig `yaml:"prometheus"`
}
