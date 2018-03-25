package cmd

import (
	"context"
	"fmt"
	"github.com/blevesearch/bleve/mapping"
	"github.com/mosuka/blast/index/config"
	"github.com/mosuka/blast/index/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type StartNodeCmdOpts struct {
	//configFile string

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
}

var defaultStartNodeCmdOpts = StartNodeCmdOpts{
	//configFile: "",

	logFormat: "text",
	logOutput: "",
	logLevel:  "info",

	grpcListenAddress: "0.0.0.0:5000",

	indexPath:        "./data/index",
	indexMappingFile: "",
	indexConfigFile:  "",

	httpListenAddress: "0.0.0.0:8000",

	restURI:    "/rest",
	metricsURI: "/metrics",
}

var startNodeCmdOpts = StartNodeCmdOpts{
	//configFile: defaultStartNodeCmdOpts.configFile,

	logFormat: defaultStartNodeCmdOpts.logFormat,
	logOutput: defaultStartNodeCmdOpts.logOutput,
	logLevel:  defaultStartNodeCmdOpts.logLevel,

	grpcListenAddress: defaultStartNodeCmdOpts.grpcListenAddress,

	indexPath:        defaultStartNodeCmdOpts.indexPath,
	indexMappingFile: defaultStartNodeCmdOpts.indexMappingFile,
	indexConfigFile:  defaultStartNodeCmdOpts.indexConfigFile,

	httpListenAddress: defaultStartNodeCmdOpts.httpListenAddress,

	restURI:    defaultStartNodeCmdOpts.restURI,
	metricsURI: defaultStartNodeCmdOpts.metricsURI,
}

var startNodeCmd = &cobra.Command{
	Use:   "node",
	Short: "start node",
	Long:  `The start node command starts the Blast node.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		switch startNodeCmdOpts.logFormat {
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

		switch startNodeCmdOpts.logLevel {
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

		if startNodeCmdOpts.logOutput == "" {
			log.SetOutput(os.Stdout)
		} else {
			var err error
			logOutput, err := os.OpenFile(startNodeCmdOpts.logOutput, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
			if err != nil {
				return err
			} else {
				log.SetOutput(logOutput)
			}
			defer logOutput.Close()
		}

		// create index mapping
		indexMapping := mapping.NewIndexMapping()
		if startNodeCmdOpts.indexMappingFile != "" {
			file, err := os.Open(startNodeCmdOpts.indexMappingFile)
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
		if startNodeCmdOpts.indexConfigFile != "" {
			file, err := os.Open(startNodeCmdOpts.indexConfigFile)
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
		gRPCServer, err := server.NewGRPCServer(
			startNodeCmdOpts.grpcListenAddress,
			startNodeCmdOpts.indexPath,
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
		httpServer, err := server.NewHTTPServer(
			startNodeCmdOpts.httpListenAddress,
			startNodeCmdOpts.restURI,
			startNodeCmdOpts.metricsURI,
			context.Background(),
			startNodeCmdOpts.grpcListenAddress,
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

//func loadConfig() {
//	viper.SetDefault("config_file", defaultStartNodeCmdOpts.configFile)
//	viper.SetDefault("log_format", defaultStartNodeCmdOpts.logFormat)
//	viper.SetDefault("log_output", defaultStartNodeCmdOpts.logOutput)
//	viper.SetDefault("log_level", defaultStartNodeCmdOpts.logLevel)
//	viper.SetDefault("grpc_listen_address", defaultStartNodeCmdOpts.grpcListenAddress)
//	viper.SetDefault("index_path", defaultStartNodeCmdOpts.indexPath)
//	viper.SetDefault("index_mapping_file", defaultStartNodeCmdOpts.indexMappingFile)
//	viper.SetDefault("index_config_file", defaultStartNodeCmdOpts.indexConfigFile)
//	viper.SetDefault("http_listen_address", defaultStartNodeCmdOpts.httpListenAddress)
//	viper.SetDefault("rest_uri", defaultStartNodeCmdOpts.restURI)
//	viper.SetDefault("metrics_uri", defaultStartNodeCmdOpts.metricsURI)
//
//	viper.SetEnvPrefix("blast_node")
//	viper.AutomaticEnv()
//
//	if startNodeCmdOpts.configFile != "" {
//		viper.SetConfigFile(startNodeCmdOpts.configFile)
//	} else {
//		viper.SetConfigName("blast")
//		viper.SetConfigType("yaml")
//		viper.AddConfigPath("/etc")
//		viper.AddConfigPath("${HOME}/etc")
//		viper.AddConfigPath("./etc")
//	}
//	viper.ReadInConfig()
//}

func init() {
	//cobra.OnInitialize(loadConfig)

	startNodeCmd.Flags().SortFlags = false

	//startNodeCmd.Flags().StringVar(&startNodeCmdOpts.configFile, "config-file", defaultStartNodeCmdOpts.configFile, "config file path")
	startNodeCmd.Flags().StringVar(&startNodeCmdOpts.logFormat, "log-format", defaultStartNodeCmdOpts.logFormat, "log format")
	startNodeCmd.Flags().StringVar(&startNodeCmdOpts.logOutput, "log-output", defaultStartNodeCmdOpts.logOutput, "log output path")
	startNodeCmd.Flags().StringVar(&startNodeCmdOpts.logLevel, "log-level", defaultStartNodeCmdOpts.logLevel, "log level")
	startNodeCmd.Flags().StringVar(&startNodeCmdOpts.grpcListenAddress, "grpc-listen-address", defaultStartNodeCmdOpts.grpcListenAddress, "address to listen for the gRPC")
	startNodeCmd.Flags().StringVar(&startNodeCmdOpts.indexPath, "index-path", defaultStartNodeCmdOpts.indexPath, "index directory path")
	startNodeCmd.Flags().StringVar(&startNodeCmdOpts.indexMappingFile, "index-mapping-file", defaultStartNodeCmdOpts.indexMappingFile, "index mapping file path")
	startNodeCmd.Flags().StringVar(&startNodeCmdOpts.indexConfigFile, "index-config-file", defaultStartNodeCmdOpts.indexConfigFile, "index config file path")
	startNodeCmd.Flags().StringVar(&startNodeCmdOpts.httpListenAddress, "http-listen-address", defaultStartNodeCmdOpts.httpListenAddress, "address to listen for the HTTP")
	startNodeCmd.Flags().StringVar(&startNodeCmdOpts.restURI, "rest-uri", defaultStartNodeCmdOpts.restURI, "base URI for REST endpoint")
	startNodeCmd.Flags().StringVar(&startNodeCmdOpts.metricsURI, "metrics-uri", defaultStartNodeCmdOpts.metricsURI, "base URI for metrics endpoint")

	//viper.BindPFlag("config_file", startNodeCmd.Flags().Lookup("config-file"))
	//viper.BindPFlag("log_format", startNodeCmd.Flags().Lookup("log-format"))
	//viper.BindPFlag("log_output", startNodeCmd.Flags().Lookup("log-output"))
	//viper.BindPFlag("log_level", startNodeCmd.Flags().Lookup("log-level"))
	//viper.BindPFlag("grpc_listen_address", startNodeCmd.Flags().Lookup("grpc-listen-address"))
	//viper.BindPFlag("index_path", startNodeCmd.Flags().Lookup("index-path"))
	//viper.BindPFlag("index_mapping_file", startNodeCmd.Flags().Lookup("index-mapping-file"))
	//viper.BindPFlag("index_config_file", startNodeCmd.Flags().Lookup("index-config-file"))
	//viper.BindPFlag("http_listen_address", startNodeCmd.Flags().Lookup("http-listen-address"))
	//viper.BindPFlag("rest_uri", startNodeCmd.Flags().Lookup("rest-uri"))
	//viper.BindPFlag("metrics_uri", startNodeCmd.Flags().Lookup("metrics-uri"))

	startCmd.AddCommand(startNodeCmd)
}
