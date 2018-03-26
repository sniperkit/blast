package config

import "github.com/spf13/viper"

const (
	DefaultLogFormat         = "text"
	DefaultLogOutput         = ""
	DefaultLogLevel          = "info"
	DefaultGRPCListenAddress = "0.0.0.0:5000"
	DefaultIndexPath         = "./data/index"
	DefaultIndexMappingPath  = "./data/index"
	DefaultIndexConfigPath   = "./data/index"
	DefaultHTTPListenAddress = "0.0.0.0:5000"
	DefaultRESTURI           = "/rest"
	DefaultMetricsURI        = "/metrics"
)

func NewNodeConfig(nodeConfigPath string) (*viper.Viper, error) {
	nodeConfig := viper.New()
	nodeConfig.SetDefault("log_format", DefaultLogFormat)
	nodeConfig.SetDefault("log_output", DefaultLogOutput)
	nodeConfig.SetDefault("log_level", DefaultLogLevel)
	nodeConfig.SetDefault("grpc_listen_address", DefaultGRPCListenAddress)
	nodeConfig.SetDefault("index_path", DefaultIndexPath)
	nodeConfig.SetDefault("index_mapping_path", DefaultIndexMappingPath)
	nodeConfig.SetDefault("index_config_file", DefaultIndexConfigPath)
	nodeConfig.SetDefault("http_listen_address", DefaultHTTPListenAddress)
	nodeConfig.SetDefault("rest_uri", DefaultRESTURI)
	nodeConfig.SetDefault("metrics_uri", DefaultMetricsURI)

	nodeConfig.SetEnvPrefix("blast_node")
	nodeConfig.AutomaticEnv()

	if nodeConfigPath != "" {
		nodeConfig.SetConfigFile(nodeConfigPath)
	} else {
		nodeConfig.SetConfigName("blast")
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
