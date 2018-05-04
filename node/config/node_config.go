//  Copyright (c) 2018 Minoru Osuka
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 		http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"github.com/mosuka/blast/config"
	"github.com/spf13/viper"
)

func NewNodeConfig(configPath string) (*viper.Viper, error) {
	nodeConfig := viper.New()
	nodeConfig.SetDefault("log_format", config.DefaultLogFormat)
	nodeConfig.SetDefault("log_output", config.DefaultLogOutput)
	nodeConfig.SetDefault("log_level", config.DefaultLogLevel)
	nodeConfig.SetDefault("grpc_listen_address", config.DefaultGRPCListenAddress)
	nodeConfig.SetDefault("index_path", config.DefaultIndexPath)
	nodeConfig.SetDefault("index_mapping_path", config.DefaultIndexMappingPath)
	nodeConfig.SetDefault("index_config_path", config.DefaultIndexConfigPath)
	nodeConfig.SetDefault("http_listen_address", config.DefaultHTTPListenAddress)
	nodeConfig.SetDefault("rest_uri", config.DefaultRESTURI)
	nodeConfig.SetDefault("metrics_uri", config.DefaultMetricsURI)

	nodeConfig.SetEnvPrefix("blast_node")
	nodeConfig.AutomaticEnv()

	if configPath != "" {
		nodeConfig.SetConfigFile(configPath)

		err := nodeConfig.ReadInConfig()
		if err != nil {
			return nil, err
		}
	}

	return nodeConfig, nil
}
