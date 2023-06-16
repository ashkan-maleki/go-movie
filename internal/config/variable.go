package config

const (
	InfraConfigFile          = "INFRA_CONFIG_FILE"
	AppConfigFile            = "APP_CONFIG_FILE"
	VarServiceDiscoveryUrl   = "SERVICE_DISCOVERY_URL"
	VarJaegerUrl             = "JAEGER_URL"
	VarHttpServerPort        = "HTTP_SERVER_PORT"
	VarPrometheusMetricsPort = "PROMETHEUS_METRICS_PORT"
)

var (
	EnvVarFiles = []string{
		InfraConfigFile,
		AppConfigFile,
	}

	EnvVars = []string{
		VarServiceDiscoveryUrl,
		VarJaegerUrl,
		VarHttpServerPort,
		VarPrometheusMetricsPort,
	}
)
