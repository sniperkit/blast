package config

import "github.com/spf13/viper"

const (
	DefaultConfigPath        = "./etc/blast_node.yml"
	DefaultLogFormat         = "text"
	DefaultLogOutput         = ""
	DefaultLogLevel          = "info"
	DefaultGRPCListenAddress = "0.0.0.0:5000"
	DefaultIndexPath         = "./data/index"
	DefaultIndexMappingPath  = "./etc/index_mapping.json"
	DefaultIndexConfigPath   = "./etc/index_config.json"
	DefaultHTTPListenAddress = "0.0.0.0:8000"
	DefaultRESTURI           = "/rest"
	DefaultMetricsURI        = "/metrics"
)

func NewConfig(configPath string) (*viper.Viper, error) {
	nodeConfig := viper.New()
	nodeConfig.SetDefault("log_format", DefaultLogFormat)
	nodeConfig.SetDefault("log_output", DefaultLogOutput)
	nodeConfig.SetDefault("log_level", DefaultLogLevel)
	nodeConfig.SetDefault("grpc_listen_address", DefaultGRPCListenAddress)
	nodeConfig.SetDefault("index_path", DefaultIndexPath)
	nodeConfig.SetDefault("index_mapping_path", DefaultIndexMappingPath)
	nodeConfig.SetDefault("index_config_path", DefaultIndexConfigPath)
	nodeConfig.SetDefault("http_listen_address", DefaultHTTPListenAddress)
	nodeConfig.SetDefault("rest_uri", DefaultRESTURI)
	nodeConfig.SetDefault("metrics_uri", DefaultMetricsURI)

	nodeConfig.SetEnvPrefix("blast_node")
	nodeConfig.AutomaticEnv()

	if configPath != "" {
		nodeConfig.SetConfigFile(configPath)
	} else {
		nodeConfig.SetConfigName("blast_node")
		nodeConfig.SetConfigType("yaml")
		nodeConfig.AddConfigPath("/etc")
		nodeConfig.AddConfigPath("${HOME}/etc")
		nodeConfig.AddConfigPath("./etc")
	}
	err := nodeConfig.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return nodeConfig, nil
}
