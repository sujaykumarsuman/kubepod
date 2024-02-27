package kubepod

import (
	"context"
	"encoding/base64"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
	"sigs.k8s.io/aws-iam-authenticator/pkg/token"
)

type Kubepod struct {
	*kubernetes.Clientset
	cluster *types.Cluster
	logger  *zap.Logger
}

type Interface interface {
	GetNodes() (*v1.NodeList, error)
	GetNode(name string) (*v1.Node, error)
	GetPods(namespace string) (*v1.PodList, error)
	GetPod(name, namespace string) (*v1.Pod, error)
}

func newClientset(log *zap.Logger, cluster *types.Cluster) (*kubernetes.Clientset, error) {
	log.Info("Creating new clientSet with cluster", zap.Any("cluster", cluster))
	gen, err := token.NewGenerator(true, false)
	if err != nil {
		return nil, err
	}
	opts := &token.GetTokenOptions{
		ClusterID: *cluster.Name,
	}
	tok, err := gen.GetWithOptions(opts)
	if err != nil {
		return nil, err
	}
	ca, err := base64.StdEncoding.DecodeString(*cluster.CertificateAuthority.Data)
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(
		&rest.Config{
			Host:        *cluster.Endpoint,
			BearerToken: tok.Token,
			TLSClientConfig: rest.TLSClientConfig{
				CAData: ca,
			},
		},
	)
	if err != nil {
		return nil, err
	}
	return clientset, nil
}

func NewKubepod(ctx context.Context, logger *zap.Logger, eksClient *eks.Client, clusterName string) *Kubepod {
	clusterDescription, err := eksClient.DescribeCluster(ctx, &eks.DescribeClusterInput{Name: &clusterName})
	if err != nil {
		logger.Fatal("unable to describe cluster", zap.Error(err))
		return nil
	}

	clientset, err := newClientset(logger, clusterDescription.Cluster)
	if err != nil {
		logger.Fatal("unable to create clientset", zap.Error(err))
		return nil
	}

	return &Kubepod{clientset, clusterDescription.Cluster, logger}
}

// CreateEKSClient creates a new AWS session with the provided ARN and returns the credentials
func CreateEKSClient(ctx context.Context, logger *zap.Logger, arn string) (*eks.Client, error) {
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
		setAWSCredsEnv(logger, &creds)
	}
	cfg.Credentials = aws.NewCredentialsCache(roleProvider)

	eksClient := eks.NewFromConfig(cfg)
	return eksClient, nil
}

// setAWSCredsEnv sets the AWS credentials in the environment
func setAWSCredsEnv(logger *zap.Logger, creds *aws.Credentials) {
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
