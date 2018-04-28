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

import "github.com/spf13/viper"

const (
	DefaultConfigPath        = ""
	DefaultLogFormat         = "text"
	DefaultLogOutput         = ""
	DefaultLogLevel          = "info"
	DefaultGRPCListenAddress = "0.0.0.0:5000"
	DefaultClusterMetaPath   = ""
	DefaultHTTPListenAddress = "0.0.0.0:8000"
	DefaultRESTURI           = "/rest"
	DefaultMetricsURI        = "/metrics"
)

func NewConfig(configPath string) (*viper.Viper, error) {
	masterConfig := viper.New()
	masterConfig.SetDefault("log_format", DefaultLogFormat)
	masterConfig.SetDefault("log_output", DefaultLogOutput)
	masterConfig.SetDefault("log_level", DefaultLogLevel)
	masterConfig.SetDefault("grpc_listen_address", DefaultGRPCListenAddress)
	masterConfig.SetDefault("supervise_config_path", DefaultClusterMetaPath)
	masterConfig.SetDefault("http_listen_address", DefaultHTTPListenAddress)
	masterConfig.SetDefault("rest_uri", DefaultRESTURI)
	masterConfig.SetDefault("metrics_uri", DefaultMetricsURI)

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
