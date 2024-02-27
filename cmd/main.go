package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/sujaykumarsuman/kubepod/pkg/api"
	"github.com/sujaykumarsuman/kubepod/pkg/kubepod"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const (
	// App
	defaultConfigDirectory = "deploy/"
	defaultConfigFile      = ""
	defaultApplicationID   = "kubepod"
	defaultLogLevel        = "info"
	defaultHost            = ""
	defaultPort            = "8080"

	// AWS
	defaultAWSRegion   = "us-east-1"
	defaultClusterName = "ddi-dev-use1"
	defaultBiabARN     = "arn:aws:iam::123456789012:role/biab-eks-access"
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
	_ = pflag.String("app.id", defaultApplicationID, "identifier for the application")
	_ = pflag.String("log-level", defaultLogLevel, "log level (debug, info, warn, error, dpanic, panic, fatal)")
	_ = pflag.String("host", defaultHost, "address to serve requests")
	_ = pflag.String("port", defaultPort, "port to serve requests")

	// AWS
	_ = pflag.String("aws.region", defaultAWSRegion, "AWS region")
	_ = pflag.String("eks.cluster.name", defaultClusterName, "name of the EKS cluster")
	_ = pflag.String("aws.arn", defaultBiabARN, "ARN of the BIAB role")

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

	// gracefully exit on keyboard interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	// create EKS Client
	eksClient, err := createEKSClient(ctx, viper.GetString("aws.arn"))
	if err != nil {
		logger.Fatal("unable to create EKS client", zap.Error(err))
		return
	}

	// create kubepod
	kp := kubepod.NewKubepod(ctx, logger, eksClient, viper.GetString("eks.cluster.name"))
	if kp == nil {
		logger.Fatal("unable to create kubepod")
		return
	}

	// start the api server
	addr := viper.GetString("host") + ":" + viper.GetString("port")
	r := api.GetRouter(logger, kp)
	go func() {
		if err := http.ListenAndServe(addr, r); err != nil {
			logger.Error("failed to start server", zap.Error(err))
			os.Exit(1)
		}
	}()

	logger.Info("ready to serve requests on " + addr)
	<-c
	logger.Info("gracefully shutting down")
	os.Exit(0)
}

// createEKSClient creates a new AWS session with the provided ARN and returns the credentials
func createEKSClient(ctx context.Context, arn string) (*eks.Client, error) {
	cfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(viper.GetString("aws.region")))
	if err != nil {
		return nil, err
	}
	stsSvc := sts.NewFromConfig(cfg)
	roleProvider := stscreds.NewAssumeRoleProvider(stsSvc, arn)
	if err != nil {
		logger.Fatal("unable to assume role", zap.Error(err))
		return nil, err
	} else {
		creds, err := roleProvider.Retrieve(ctx)
		if err != nil {
			logger.Fatal("unable to retrieve credentials", zap.Error(err))
			return nil, err
		}
		setAWSCredsEnv(&creds)
	}
	cfg.Credentials = aws.NewCredentialsCache(roleProvider)

	eksClient := eks.NewFromConfig(cfg)
	return eksClient, nil
}

// setAWSCredsEnv sets the AWS credentials in the environment
func setAWSCredsEnv(creds *aws.Credentials) {
	if err := os.Setenv("AWS_ACCESS_KEY_ID", creds.AccessKeyID); err != nil {
		logger.Error("unable to set AWS_ACCESS_KEY_ID", zap.Error(err))
	} else {
		logger.Debug("env AWS_ACCESS_KEY_ID set", zap.String("AWS_ACCESS_KEY_ID", creds.AccessKeyID))
	}
	if err := os.Setenv("AWS_SECRET_ACCESS_KEY", creds.SecretAccessKey); err != nil {
		logger.Error("unable to set AWS_SECRET_ACCESS_KEY", zap.Error(err))
	} else {
		logger.Debug("env AWS_SECRET_ACCESS_KEY set", zap.String("AWS_SECRET_ACCESS_KEY", creds.SecretAccessKey))
	}
	if err := os.Setenv("AWS_SESSION_TOKEN", creds.SessionToken); err != nil {
		logger.Error("unable to set AWS_SESSION_TOKEN", zap.Error(err))
	} else {
		logger.Debug("env AWS_SESSION_TOKEN set", zap.String("AWS_SESSION_TOKEN", creds.SessionToken))
	}
}
