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

func printConfig(logger *zap.Logger, cfg aws.Config) {
	logger.Debug("cfg.region", zap.String("region", cfg.Region))
	logger.Debug("cfg.credentials", zap.Any("credentials", cfg.Credentials))
}

func NewKubepod(ctx context.Context, logger *zap.Logger, arn, clusterName string) *Kubepod {
	cfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(viper.GetString("aws.region")))
	if err != nil {
		logger.Fatal("unable to load SDK config", zap.Error(err))
		return nil
	}
	printConfig(logger, cfg)
	stsSvc := sts.NewFromConfig(cfg)
	creds := stscreds.NewAssumeRoleProvider(stsSvc, arn)
	if err != nil {
		logger.Fatal("unable to assume role", zap.Error(err))
		return nil
	}
	cfg.Credentials = aws.NewCredentialsCache(creds)
	logger.Debug("cfg.credentials", zap.Any("credentials", cfg.Credentials))
	eksClient := eks.NewFromConfig(cfg)
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
