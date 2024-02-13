package kubepod

import (
	"context"
	"encoding/base64"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/aws-iam-authenticator/pkg/token"
)

type Kubepod struct {
	*kubernetes.Clientset
	logger *zap.Logger
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

func NewKubepod(ctx context.Context, logger *zap.Logger, clusterName string) *Kubepod {
	cfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(viper.GetString("aws.region")))
	if err != nil {
		logger.Fatal("unable to load SDK config", zap.Error(err))
		return nil
	}
	printConfig(logger, cfg)
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

	return &Kubepod{clientset, logger}
}

func (k *Kubepod) GetPods(namespace string) {
	pods, err := k.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		k.logger.Fatal("unable to list pods", zap.Error(err))
		return
	}
	for _, pod := range pods.Items {
		k.logger.Info("pod", zap.String("name", pod.Name))
	}
}

func (k *Kubepod) GetNodes() {
	nodes, err := k.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		k.logger.Fatal("unable to list nodes", zap.Error(err))
		return
	}
	for _, node := range nodes.Items {
		k.logger.Info("node", zap.String("name", node.Name))
	}
}
