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
	"fmt"
	"github.com/mosuka/blast/cluster"
	"github.com/mosuka/blast/config"
	masterconfig "github.com/mosuka/blast/master/config"
	"github.com/mosuka/blast/master/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type StartMasterCmdOpts struct {
	configPath string

	logFormat string
	logOutput string
	logLevel  string

	grpcListenAddress string

	clusterMetaPath string

	httpListenAddress string

	restURI    string
	metricsURI string
}

var startMasterCmdOpts = StartMasterCmdOpts{
	configPath: config.DefaultConfigPath,

	logFormat: config.DefaultLogFormat,
	logOutput: config.DefaultLogOutput,
	logLevel:  config.DefaultLogLevel,

	grpcListenAddress: config.DefaultMasterGRPCListenAddress,

	clusterMetaPath: config.DefaultClusterMetaPath,

	httpListenAddress: config.DefaultMasterHTTPListenAddress,

	restURI:    config.DefaultRESTURI,
	metricsURI: config.DefaultMetricsURI,
}

var startMasterCmd = &cobra.Command{
	Use:   "master",
	Short: "start master",
	Long:  `The start master command starts the Blast master.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		masterConfig, err := masterconfig.NewMasterConfig(startMasterCmdOpts.configPath)
		if err != nil {
			return err
		}

		masterConfig.BindPFlag("log_format", cmd.Flags().Lookup("log-format"))
		masterConfig.BindPFlag("log_output", cmd.Flags().Lookup("log-output"))
		masterConfig.BindPFlag("log_level", cmd.Flags().Lookup("log-level"))
		masterConfig.BindPFlag("grpc_listen_address", cmd.Flags().Lookup("grpc-listen-address"))
		masterConfig.BindPFlag("cluster_meta_path", cmd.Flags().Lookup("cluster-meta-path"))
		masterConfig.BindPFlag("http_listen_address", cmd.Flags().Lookup("http-listen-address"))
		masterConfig.BindPFlag("rest_uri", cmd.Flags().Lookup("rest-uri"))
		masterConfig.BindPFlag("metrics_uri", cmd.Flags().Lookup("metrics-uri"))

		switch masterConfig.GetString("log_format") {
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

		switch masterConfig.GetString("log_level") {
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

		if masterConfig.GetString("log_output") == "" {
			log.SetOutput(os.Stdout)
		} else {
			var err error
			logOutput, err := os.OpenFile(startMasterCmdOpts.logOutput, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
			if err != nil {
				return err
			} else {
				log.SetOutput(logOutput)
			}
			defer logOutput.Close()
		}

		clusterMeta := cluster.NewClusterMeta()
		if masterConfig.GetString("cluster_meta_path") != "" {
			file, err := os.Open(masterConfig.GetString("cluster_meta_path"))
			if err != nil {
				log.WithFields(log.Fields{
					"error": err.Error(),
				}).Fatal(fmt.Sprintf("failed to open cluster meta file."))
				return err
			}
			defer file.Close()

			clusterMeta, err = cluster.LoadClusterMeta(file)
			if err != nil {
				log.WithFields(log.Fields{
					"error": err.Error(),
				}).Fatal(fmt.Sprintf("failed to load cluster meta file."))
				return err
			}

			log.Info(fmt.Sprintf("cluster meta file was loaded."))
		}

		// create gRPC Server
		grpcServer, err := server.NewGRPCServer(
			masterConfig.GetString("grpc_listen_address"),
			clusterMeta,
		)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Fatal(fmt.Sprintf("failed to create supervisor."))
			return err
		}
		log.Info(fmt.Sprintf("supervisor was created."))

		// start gRPC Server
		err = grpcServer.Start()
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Fatal(fmt.Sprintf("failed to start supervisor."))
			return err
		}
		log.Info(fmt.Sprintf("supervisor was started."))

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
			}).Info("trap signal")

			err = grpcServer.Stop()
			if err != nil {
				log.WithFields(log.Fields{
					"error": err.Error(),
				}).Fatal("failed to stop supervisor.")
			}
			log.Info("supervisor was stopped.")

			return nil
		}

		return nil
	},
}

func init() {
	startMasterCmd.Flags().SortFlags = false

	startMasterCmd.Flags().StringVar(&startMasterCmdOpts.configPath, "config-path", config.DefaultConfigPath, "config path")

	startMasterCmd.Flags().StringVar(&startMasterCmdOpts.logFormat, "log-format", config.DefaultLogFormat, "log format")
	startMasterCmd.Flags().StringVar(&startMasterCmdOpts.logOutput, "log-output", config.DefaultLogOutput, "log output")
	startMasterCmd.Flags().StringVar(&startMasterCmdOpts.logLevel, "log-level", config.DefaultLogLevel, "log level")
	startMasterCmd.Flags().StringVar(&startMasterCmdOpts.grpcListenAddress, "grpc-listen-address", config.DefaultMasterGRPCListenAddress, "address to listen for the gRPC")
	startMasterCmd.Flags().StringVar(&startMasterCmdOpts.clusterMetaPath, "cluster-meta-path", config.DefaultClusterMetaPath, "cluster meta path")
	startMasterCmd.Flags().StringVar(&startMasterCmdOpts.httpListenAddress, "http-listen-address", config.DefaultMasterHTTPListenAddress, "address to listen for the HTTP")
	startMasterCmd.Flags().StringVar(&startMasterCmdOpts.restURI, "rest-uri", config.DefaultRESTURI, "base URI for REST endpoint")
	startMasterCmd.Flags().StringVar(&startMasterCmdOpts.metricsURI, "metrics-uri", config.DefaultMetricsURI, "base URI for metrics endpoint")

	startCmd.AddCommand(startMasterCmd)
}
