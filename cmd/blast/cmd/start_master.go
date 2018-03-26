package cmd

import (
	"fmt"
	"github.com/mosuka/blast/master/config"
	"github.com/mosuka/blast/master/server"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type StartMasterCmdOpts struct {
	//configFile string

	logFormat string
	logOutput string
	logLevel  string

	grpcListenAddress string

	supervisorConfigFile string
}

var defaultStartMasterCmdOpts = StartMasterCmdOpts{
	//configFile: "",

	logFormat: "text",
	logOutput: "",
	logLevel:  "info",

	grpcListenAddress: "0.0.0.0:6000",

	supervisorConfigFile: "",
}

var startSupervisorCmdOpts = StartMasterCmdOpts{
	//configFile: defaultStartMasterCmdOpts.configFile,

	logFormat: defaultStartMasterCmdOpts.logFormat,
	logOutput: defaultStartMasterCmdOpts.logOutput,
	logLevel:  defaultStartMasterCmdOpts.logLevel,

	grpcListenAddress: defaultStartMasterCmdOpts.grpcListenAddress,

	supervisorConfigFile: defaultStartMasterCmdOpts.supervisorConfigFile,
}

var startSupervisorCmd = &cobra.Command{
	Use:   "master",
	Short: "start master",
	Long:  `The start master command starts the Blast master.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		switch startSupervisorCmdOpts.logFormat {
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

		switch startSupervisorCmdOpts.logLevel {
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

		if startSupervisorCmdOpts.logOutput == "" {
			log.SetOutput(os.Stdout)
		} else {
			var err error
			logOutput, err := os.OpenFile(startSupervisorCmdOpts.logOutput, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
			if err != nil {
				return err
			} else {
				log.SetOutput(logOutput)
			}
			defer logOutput.Close()
		}

		supervisorConfig := config.NewSupervisorConfig()
		if startSupervisorCmdOpts.supervisorConfigFile != "" {
			file, err := os.Open(startSupervisorCmdOpts.supervisorConfigFile)
			if err != nil {
				log.WithFields(log.Fields{
					"error": err.Error(),
				}).Fatal(fmt.Sprintf("failed to open supervisor config file."))
				return err
			}
			defer file.Close()

			supervisorConfig, err = config.LoadSupervisorConfig(file)
			if err != nil {
				log.WithFields(log.Fields{
					"error": err.Error(),
				}).Fatal(fmt.Sprintf("failed to load supervisor config file."))
				return err
			}

			log.Info(fmt.Sprintf("supervisor config file was loaded."))
		}

		// create gRPC Server
		grpcServer, err := server.NewGRPCServer(
			startSupervisorCmdOpts.grpcListenAddress,
			supervisorConfig,
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

//func loadSupervisorConfig() {
//	viper.SetDefault("log_format", startSupervisorCmdOpts.logFormat)
//	viper.SetDefault("log_output", startSupervisorCmdOpts.logOutput)
//	viper.SetDefault("log_level", startSupervisorCmdOpts.logLevel)
//	viper.SetDefault("grpc_listen_address", startSupervisorCmdOpts.grpcListenAddress)
//	viper.SetDefault("supervisor_config_file", startSupervisorCmdOpts.supervisorConfigFile)
//
//	if viper.GetString("config_file") != "" {
//		viper.SetConfigFile(viper.GetString("config"))
//	} else {
//		viper.SetConfigName("blastsv")
//		viper.SetConfigType("yaml")
//		viper.AddConfigPath("/etc")
//		viper.AddConfigPath("${HOME}/etc")
//		viper.AddConfigPath("./etc")
//	}
//	viper.SetEnvPrefix("blastsv")
//	viper.AutomaticEnv()
//
//	viper.ReadInConfig()
//}

func init() {
	//cobra.OnInitialize(loadSupervisorConfig)

	startSupervisorCmd.Flags().SortFlags = false

	//startSupervisorCmd.Flags().String("config-file", startSupervisorCmdOpts.configFile, "config file path")
	startSupervisorCmd.Flags().String("log-format", startSupervisorCmdOpts.logFormat, "log format")
	startSupervisorCmd.Flags().String("log-output", startSupervisorCmdOpts.logOutput, "log output path")
	startSupervisorCmd.Flags().String("log-level", startSupervisorCmdOpts.logLevel, "log level")
	startSupervisorCmd.Flags().String("grpc-listen-address", startSupervisorCmdOpts.grpcListenAddress, "address to listen for the gRPC")
	startSupervisorCmd.Flags().String("supervisor-config-file", startSupervisorCmdOpts.supervisorConfigFile, "supervisor config file path")

	//viper.BindPFlag("config_file", startSupervisorCmd.Flags().Lookup("config-file"))
	//viper.BindPFlag("log_format", startSupervisorCmd.Flags().Lookup("log-format"))
	//viper.BindPFlag("log_output", startSupervisorCmd.Flags().Lookup("log-output"))
	//viper.BindPFlag("log_level", startSupervisorCmd.Flags().Lookup("log-level"))
	//viper.BindPFlag("grpc_listen_address", startSupervisorCmd.Flags().Lookup("grpc-listen-address"))
	//viper.BindPFlag("supervisor_config_file", startSupervisorCmd.Flags().Lookup("supervisor-config-file"))

	startCmd.AddCommand(startSupervisorCmd)
}
