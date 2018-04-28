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

package cmd

import (
	"context"
	"fmt"
	"github.com/blevesearch/bleve/mapping"
	"github.com/mosuka/blast/config"
	"github.com/mosuka/blast/index"
	nodeconfig "github.com/mosuka/blast/node/config"
	blastgrpc "github.com/mosuka/blast/node/server/grpc"
	blasthttp "github.com/mosuka/blast/node/server/http"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type StartNodeCmdOpts struct {
	configPath string

	logFormat string
	logOutput string
	logLevel  string

	grpcListenAddress string

	indexPath        string
	indexMappingPath string
	indexConfigPath  string

	httpListenAddress string

	restURI    string
	metricsURI string
}

var startNodeCmdOpts = StartNodeCmdOpts{
	configPath: config.DefaultConfigPath,

	logFormat: config.DefaultLogFormat,
	logOutput: config.DefaultLogOutput,
	logLevel:  config.DefaultLogLevel,

	grpcListenAddress: config.DefaultNodeGRPCListenAddress,

	indexPath:        config.DefaultIndexPath,
	indexMappingPath: config.DefaultIndexMappingPath,
	indexConfigPath:  config.DefaultIndexConfigPath,

	httpListenAddress: config.DefaultNodeHTTPListenAddress,

	restURI:    config.DefaultRESTURI,
	metricsURI: config.DefaultMetricsURI,
}

var startNodeCmd = &cobra.Command{
	Use:   "node",
	Short: "start node",
	Long:  `The start node command starts the Blast node.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		nodeConfig, err := nodeconfig.NewNodeConfig(startNodeCmdOpts.configPath)
		if err != nil {
			return err
		}

		nodeConfig.BindPFlag("log_format", cmd.Flags().Lookup("log-format"))
		nodeConfig.BindPFlag("log_output", cmd.Flags().Lookup("log-output"))
		nodeConfig.BindPFlag("log_level", cmd.Flags().Lookup("log-level"))
		nodeConfig.BindPFlag("grpc_listen_address", cmd.Flags().Lookup("grpc-listen-address"))
		nodeConfig.BindPFlag("index_path", cmd.Flags().Lookup("index-path"))
		nodeConfig.BindPFlag("index_mapping_path", cmd.Flags().Lookup("index-mapping-path"))
		nodeConfig.BindPFlag("index_config_path", cmd.Flags().Lookup("index-config-path"))
		nodeConfig.BindPFlag("http_listen_address", cmd.Flags().Lookup("http-listen-address"))
		nodeConfig.BindPFlag("rest_uri", cmd.Flags().Lookup("rest-uri"))
		nodeConfig.BindPFlag("metrics_uri", cmd.Flags().Lookup("metrics-uri"))

		switch nodeConfig.GetString("log_format") {
		case "text":
			log.SetFormatter(&log.TextFormatter{
				ForceColors:      false,
				DisableColors:    true,
				DisableTimestamp: false,
				FullTimestamp:    true,
				TimestampFormat:  time.RFC3339,
				DisableSorting:   false,
				QuoteEmptyFields: true,
			})
		case "color":
			log.SetFormatter(&log.TextFormatter{
				ForceColors:      true,
				DisableColors:    false,
				DisableTimestamp: false,
				FullTimestamp:    true,
				TimestampFormat:  time.RFC3339,
				DisableSorting:   false,
				QuoteEmptyFields: true,
			})
		case "json":
			log.SetFormatter(&log.JSONFormatter{
				TimestampFormat:  time.RFC3339,
				DisableTimestamp: false,
				FieldMap: log.FieldMap{
					log.FieldKeyTime:  "@timestamp",
					log.FieldKeyLevel: "@level",
					log.FieldKeyMsg:   "@message",
				},
			})
		default:
			log.SetFormatter(&log.TextFormatter{
				ForceColors:      false,
				DisableColors:    true,
				DisableTimestamp: false,
				FullTimestamp:    true,
				TimestampFormat:  time.RFC3339,
				DisableSorting:   false,
				QuoteEmptyFields: true,
			})
		}

		switch nodeConfig.GetString("log_level") {
		case "debug":
			log.SetLevel(log.DebugLevel)
		case "info":
			log.SetLevel(log.InfoLevel)
		case "warn":
			log.SetLevel(log.WarnLevel)
		case "error":
			log.SetLevel(log.ErrorLevel)
		case "fatal":
			log.SetLevel(log.FatalLevel)
		case "panic":
			log.SetLevel(log.PanicLevel)
		default:
			log.SetLevel(log.InfoLevel)
		}

		if nodeConfig.GetString("log_output") == "" {
			log.SetOutput(os.Stdout)
		} else {
			var err error
			logOutput, err := os.OpenFile(nodeConfig.GetString("log_output"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
			if err != nil {
				return err
			} else {
				log.SetOutput(logOutput)
			}
			defer logOutput.Close()
		}

		// create index mapping
		indexMapping := mapping.NewIndexMapping()
		if nodeConfig.GetString("index_mapping_path") != "" {
			file, err := os.Open(nodeConfig.GetString("index_mapping_path"))
			if err != nil {
				log.WithFields(log.Fields{
					"error": err.Error(),
				}).Fatal(fmt.Sprintf("failed to open index mapping file."))
				return err
			}
			defer file.Close()

			indexMapping, err = index.LoadIndexMapping(file)
			if err != nil {
				log.WithFields(log.Fields{
					"error": err.Error(),
				}).Fatal(fmt.Sprintf("failed to load index mapping file."))
				return err
			}

			log.Info(fmt.Sprintf("index mapping file was loaded."))
		}

		indexConfig := index.NewIndexMeta()
		if nodeConfig.GetString("index_config_path") != "" {
			file, err := os.Open(nodeConfig.GetString("index_config_path"))
			if err != nil {
				log.WithFields(log.Fields{
					"error": err.Error(),
				}).Fatal(fmt.Sprintf("failed to open index config file."))
				return err
			}
			defer file.Close()

			indexConfig, err = index.LoadIndexMeta(file)
			if err != nil {
				log.WithFields(log.Fields{
					"error": err.Error(),
				}).Fatal(fmt.Sprintf("failed to load index config file."))
				return err
			}

			log.Info(fmt.Sprintf("index config file was loaded."))
		}

		// create gRPC Server
		gRPCServer, err := blastgrpc.NewGRPCServer(
			nodeConfig.GetString("grpc_listen_address"),
			nodeConfig.GetString("index_path"),
			indexMapping,
			indexConfig,
		)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Fatal("failed to create gRPC server.")
			return err
		}
		log.Info(fmt.Sprintf("gRPC server was created."))

		// start gRPC Server
		err = gRPCServer.Start()
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Fatal("failed to start index server.")
			return err
		}
		log.Info(fmt.Sprintf("gRPC server was started."))

		// create HTTP Server
		httpServer, err := blasthttp.NewHTTPServer(
			nodeConfig.GetString("http_listen_address"),
			nodeConfig.GetString("rest_uri"),
			nodeConfig.GetString("metrics_uri"),
			context.Background(),
			nodeConfig.GetString("grpc_listen_address"),
			grpc.WithInsecure(),
		)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Fatal("failed to create HTTP server.")
			return err
		}
		log.Info(fmt.Sprintf("HTTP server was created."))

		// start HTTP Server
		err = httpServer.Start()
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Fatal("failed to start HTTP server.")
			return err
		}
		log.Info(fmt.Sprintf("HTTP server was started."))

		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan,
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT)
		for {
			sig := <-signalChan

			log.WithFields(log.Fields{
				"signal": sig,
			}).Info("trap signal.")

			err = httpServer.Stop()
			if err != nil {
				log.WithFields(log.Fields{
					"error": err.Error(),
				}).Fatal("failed to stop HTTP server.")
			}
			log.Info("HTTP server was stopped.")

			err = gRPCServer.Stop()
			if err != nil {
				log.WithFields(log.Fields{
					"error": err.Error(),
				}).Fatal("failed to stop index server.")
			}
			log.Info("index server was stopped.")

			return nil
		}

		return nil
	},
}

func init() {
	startNodeCmd.Flags().SortFlags = false

	startNodeCmd.Flags().StringVar(&startNodeCmdOpts.configPath, "config-path", config.DefaultConfigPath, "config path")

	startNodeCmd.Flags().StringVar(&startNodeCmdOpts.logFormat, "log-format", config.DefaultLogFormat, "log format")
	startNodeCmd.Flags().StringVar(&startNodeCmdOpts.logOutput, "log-output", config.DefaultLogOutput, "log output")
	startNodeCmd.Flags().StringVar(&startNodeCmdOpts.logLevel, "log-level", config.DefaultLogLevel, "log level")
	startNodeCmd.Flags().StringVar(&startNodeCmdOpts.grpcListenAddress, "grpc-listen-address", config.DefaultNodeGRPCListenAddress, "address to listen for the gRPC")
	startNodeCmd.Flags().StringVar(&startNodeCmdOpts.indexPath, "index-path", config.DefaultIndexPath, "index directory path")
	startNodeCmd.Flags().StringVar(&startNodeCmdOpts.indexMappingPath, "index-mapping-path", config.DefaultIndexMappingPath, "index mapping path")
	startNodeCmd.Flags().StringVar(&startNodeCmdOpts.indexConfigPath, "index-config-path", config.DefaultIndexConfigPath, "index config path")
	startNodeCmd.Flags().StringVar(&startNodeCmdOpts.httpListenAddress, "http-listen-address", config.DefaultNodeHTTPListenAddress, "address to listen for the HTTP")
	startNodeCmd.Flags().StringVar(&startNodeCmdOpts.restURI, "rest-uri", config.DefaultRESTURI, "base URI for REST endpoint")
	startNodeCmd.Flags().StringVar(&startNodeCmdOpts.metricsURI, "metrics-uri", config.DefaultMetricsURI, "base URI for metrics endpoint")

	startCmd.AddCommand(startNodeCmd)
}
