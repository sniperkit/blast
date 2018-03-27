package config

import "github.com/spf13/viper"

const (
	DefaultConfigPath          = "./etc/blast_master.yaml"
	DefaultLogFormat           = "text"
	DefaultLogOutput           = ""
	DefaultLogLevel            = "info"
	DefaultGRPCListenAddress   = "0.0.0.0:5000"
	DefaultSuperviseConfigPath = "./etc/supervise_config.json"
	DefaultHTTPListenAddress   = "0.0.0.0:8000"
	DefaultRESTURI             = "/rest"
	DefaultMetricsURI          = "/metrics"
)

func NewConfig(configPath string) (*viper.Viper, error) {
	masterConfig := viper.New()
	masterConfig.SetDefault("log_format", DefaultLogFormat)
	masterConfig.SetDefault("log_output", DefaultLogOutput)
	masterConfig.SetDefault("log_level", DefaultLogLevel)
	masterConfig.SetDefault("grpc_listen_address", DefaultGRPCListenAddress)
	masterConfig.SetDefault("supervise_config_path", DefaultSuperviseConfigPath)
	masterConfig.SetDefault("http_listen_address", DefaultHTTPListenAddress)
	masterConfig.SetDefault("rest_uri", DefaultRESTURI)
	masterConfig.SetDefault("metrics_uri", DefaultMetricsURI)

	masterConfig.SetEnvPrefix("blast_master")
	masterConfig.AutomaticEnv()

	if configPath != "" {
		masterConfig.SetConfigFile(configPath)
	} else {
		masterConfig.SetConfigName("blast_master")
		masterConfig.SetConfigType("yaml")
		masterConfig.AddConfigPath("/etc")
		masterConfig.AddConfigPath("${HOME}/etc")
		masterConfig.AddConfigPath("./etc")
	}
	err := masterConfig.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return masterConfig, nil
}
