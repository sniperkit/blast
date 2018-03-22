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
	"github.com/mosuka/blast/index/config"
	"github.com/mosuka/blast/index/server"
	"github.com/mosuka/blast/version"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type RootCommandOptions struct {
	configFile string

	logFormat string
	logOutput string
	logLevel  string

	grpcListenAddress string

	indexPath        string
	indexMappingFile string
	indexConfigFile  string

	httpListenAddress string

	restURI    string
	metricsURI string

	versionFlag bool
}

var rootCmdOpts = RootCommandOptions{
	configFile: "",

	logFormat: "json",
	logOutput: "",
	logLevel:  "info",

	grpcListenAddress: "0.0.0.0:5000",

	indexPath:        "./data/index",
	indexMappingFile: "",
	indexConfigFile:  "",

	httpListenAddress: "0.0.0.0:8000",

	restURI:    "/rest",
	metricsURI: "/metrics",

	versionFlag: false,
}

var logOutput *os.File

var RootCmd = &cobra.Command{
	Use:                "blast",
	Short:              "Blast",
	Long:               `The Command Line Interface for the Blast.`,
	PersistentPreRunE:  persistentPreRunERootCmd,
	RunE:               runERootCmd,
	PersistentPostRunE: persistentPostRunERootCmd,
}

func persistentPreRunERootCmd(cmd *cobra.Command, args []string) error {
	if rootCmdOpts.versionFlag {
		fmt.Printf("%s\n", version.Version)
		os.Exit(0)
	}

	switch viper.GetString("log_format") {
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

	switch viper.GetString("log_level") {
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

	if viper.GetString("log_output") == "" {
		log.SetOutput(os.Stdout)
	} else {
		var err error
		logOutput, err = os.OpenFile(viper.GetString("log_output"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		} else {
			log.SetOutput(logOutput)
		}
	}

	return nil
}

func runERootCmd(cmd *cobra.Command, args []string) error {
	// create index mapping
	indexMapping := mapping.NewIndexMapping()
	if viper.GetString("index_mapping_file") != "" {
		file, err := os.Open(viper.GetString("index_mapping_file"))
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Fatal(fmt.Sprintf("failed to open index mapping file."))
			return err
		}
		defer file.Close()

		indexMapping, err = config.LoadIndexMapping(file)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Fatal(fmt.Sprintf("failed to load index mapping file."))
			return err
		}

		log.Info(fmt.Sprintf("index mapping file was loaded."))
	}

	indexConfig := config.NewIndexConfig()
	if viper.GetString("index_config_file") != "" {
		file, err := os.Open(viper.GetString("index_config_file"))
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Fatal(fmt.Sprintf("failed to open index config file."))
			return err
		}
		defer file.Close()

		indexConfig, err = config.LoadIndexConfig(file)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Fatal(fmt.Sprintf("failed to load index config file."))
			return err
		}

		log.Info(fmt.Sprintf("index config file was loaded."))
	}

	// create gRPC Server
	indexServer, err := server.NewGRPCServer(
		viper.GetString("grpc_listen_address"),
		viper.GetString("index_path"),
		indexMapping,
		indexConfig,
	)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Fatal("failed to create index server.")
		return err
	}
	log.Info(fmt.Sprintf("index server was created."))

	// start gRPC Server
	err = indexServer.Start()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Fatal("failed to start index server.")
		return err
	}
	log.Info(fmt.Sprintf("index server was started."))

	// create HTTP Server
	httpServer, err := server.NewHTTPServer(
		viper.GetString("http_listen_address"),
		viper.GetString("rest_uri"),
		viper.GetString("metrics_uri"),
		context.Background(),
		viper.GetString("grpc_listen_address"),
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

		err = indexServer.Stop()
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Fatal("failed to stop index server.")
		}
		log.Info("index server was stopped.")

		return nil
	}

	return nil
}

func persistentPostRunERootCmd(cmd *cobra.Command, args []string) error {
	if viper.GetString("log_output") != "" {
		logOutput.Close()
	}

	return nil
}

func LoadConfig() {
	viper.SetDefault("log_format", rootCmdOpts.logFormat)
	viper.SetDefault("log_output", rootCmdOpts.logOutput)
	viper.SetDefault("log_level", rootCmdOpts.logLevel)
	viper.SetDefault("grpc_listen_address", rootCmdOpts.grpcListenAddress)
	viper.SetDefault("index_path", rootCmdOpts.indexPath)
	viper.SetDefault("index_mapping_file", rootCmdOpts.indexMappingFile)
	viper.SetDefault("index_config_file", rootCmdOpts.indexConfigFile)
	viper.SetDefault("http_listen_address", rootCmdOpts.httpListenAddress)
	viper.SetDefault("rest_uri", rootCmdOpts.restURI)
	viper.SetDefault("metrics_uri", rootCmdOpts.metricsURI)

	if viper.GetString("config_file") != "" {
		viper.SetConfigFile(viper.GetString("config"))
	} else {
		viper.SetConfigName("blast")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("/etc")
		viper.AddConfigPath("${HOME}/etc")
		viper.AddConfigPath("./etc")
	}
	viper.SetEnvPrefix("blast")
	viper.AutomaticEnv()

	viper.ReadInConfig()
}

func init() {
	cobra.OnInitialize(LoadConfig)

	RootCmd.Flags().SortFlags = false

	RootCmd.Flags().String("config-file", rootCmdOpts.configFile, "config file path")
	RootCmd.Flags().String("log-format", rootCmdOpts.logFormat, "log format")
	RootCmd.Flags().String("log-output", rootCmdOpts.logOutput, "log output path")
	RootCmd.Flags().String("log-level", rootCmdOpts.logLevel, "log level")
	RootCmd.Flags().String("grpc-listen-address", rootCmdOpts.grpcListenAddress, "address to listen for the gRPC")
	RootCmd.Flags().String("index-path", rootCmdOpts.indexPath, "index directory path")
	RootCmd.Flags().String("index-mapping-file", rootCmdOpts.indexMappingFile, "index mapping file path")
	RootCmd.Flags().String("index-config-file", rootCmdOpts.indexConfigFile, "index config file path")
	RootCmd.Flags().String("http-listen-address", rootCmdOpts.httpListenAddress, "address to listen for the HTTP")
	RootCmd.Flags().String("rest-uri", rootCmdOpts.restURI, "base URI for REST endpoint")
	RootCmd.Flags().String("metrics-uri", rootCmdOpts.metricsURI, "base URI for metrics endpoint")
	RootCmd.Flags().BoolVarP(&rootCmdOpts.versionFlag, "version", "v", rootCmdOpts.versionFlag, "show version number")

	viper.BindPFlag("config_file", RootCmd.Flags().Lookup("config-file"))
	viper.BindPFlag("log_format", RootCmd.Flags().Lookup("log-format"))
	viper.BindPFlag("log_output", RootCmd.Flags().Lookup("log-output"))
	viper.BindPFlag("log_level", RootCmd.Flags().Lookup("log-level"))
	viper.BindPFlag("grpc_listen_address", RootCmd.Flags().Lookup("grpc-listen-address"))
	viper.BindPFlag("index_path", RootCmd.Flags().Lookup("index-path"))
	viper.BindPFlag("index_mapping_file", RootCmd.Flags().Lookup("index-mapping-file"))
	viper.BindPFlag("index_config_file", RootCmd.Flags().Lookup("index-config-file"))
	viper.BindPFlag("http_listen_address", RootCmd.Flags().Lookup("http-listen-address"))
	viper.BindPFlag("rest_uri", RootCmd.Flags().Lookup("rest-uri"))
	viper.BindPFlag("metrics_uri", RootCmd.Flags().Lookup("metrics-uri"))
}
