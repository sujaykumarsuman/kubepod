package main

import (
	"context"
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/sujaykumarsuman/kubepod/pkg/kubepod"
	"go.uber.org/zap"
	"strings"
)

const (
	// App
	defaultConfigDirectory = "deploy/"
	defaultConfigFile      = ""
	defaultSecretFile      = ""
	defaultApplicationID   = "kubepod"
	defaultLogLevel        = "info"

	// AWS
	defaultAWSRegion   = "us-east-1"
	defaultClusterName = "ddi-dev-use1"
)

var (
	logger *zap.Logger
	err    error
)

func initializeLogger() (*zap.Logger, error) {
	conf := zap.NewProductionConfig()
	err = conf.Level.UnmarshalText([]byte(viper.GetString("log-level")))
	if err != nil {
		return nil, err
	}

	return conf.Build()
}

func initializeFlags() {
	// App
	_ = pflag.String("config.source", defaultConfigDirectory, "config source")
	_ = pflag.String("config.file", defaultConfigFile, "directory of the configuration file")
	_ = pflag.String("config.secret.file", defaultSecretFile, "directory of the secrets configuration file")
	_ = pflag.String("app.id", defaultApplicationID, "identifier for the application")
	_ = pflag.String("log-level", defaultLogLevel, "log level (debug, info, warn, error, dpanic, panic, fatal)")

	// AWS
	_ = pflag.String("aws.region", defaultAWSRegion, "AWS region")
	_ = pflag.String("eks.cluster.name", defaultClusterName, "name of the EKS cluster")

	pflag.Parse()
}

func init() {
	initializeFlags()
	err = viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		fmt.Printf("[ERROR] could not bind pflags %v", err)
		return
	}
	logger, err = initializeLogger()
	if err != nil {
		logger.Fatal(fmt.Sprintf("cannot initialize logger: %v", err), zap.Error(err))
		return
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AddConfigPath(viper.GetString("config.source"))

	if viper.GetString("config.file") != "" {
		logger.Info("Serving from configuration file", zap.String("file", viper.GetString("config.file")))
		viper.SetConfigName(viper.GetString("config.file"))
		if err := viper.ReadInConfig(); err != nil {
			logger.Fatal("cannot load configuration", zap.Error(err))
		}
	} else {
		logger.Info("Serving from default values, environment variables, and/or flags")
	}
}

func main() {
	ctx := context.Background()
	// print all the configuration
	for _, key := range viper.AllKeys() {
		logger.Debug("Configuration", zap.String("key", key), zap.Any("value", viper.Get(key)))
	}

	kp := kubepod.NewKubepod(ctx, logger, viper.GetString("eks.cluster.name"))
	kp.GetNodes()
}
