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

func NewMasterConfig(configPath string) (*viper.Viper, error) {
	masterConfig := viper.New()
	masterConfig.SetDefault("log_format", config.DefaultLogFormat)
	masterConfig.SetDefault("log_output", config.DefaultLogOutput)
	masterConfig.SetDefault("log_level", config.DefaultLogLevel)
	masterConfig.SetDefault("grpc_listen_address", config.DefaultGRPCListenAddress)
	masterConfig.SetDefault("http_listen_address", config.DefaultHTTPListenAddress)
	masterConfig.SetDefault("cluster_meta_path", config.DefaultClusterMetaPath)
	masterConfig.SetDefault("rest_uri", config.DefaultRESTURI)
	masterConfig.SetDefault("metrics_uri", config.DefaultMetricsURI)

	masterConfig.SetEnvPrefix("blast_master")
	masterConfig.AutomaticEnv()

	if configPath != "" {
		masterConfig.SetConfigFile(configPath)

		err := masterConfig.ReadInConfig()
		if err != nil {
			return nil, err
		}
	}

	return masterConfig, nil
}
