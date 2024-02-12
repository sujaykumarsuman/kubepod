package kubepod

import (
	_ "github.com/aws/aws-sdk-go-v2/config"
	_ "github.com/aws/aws-sdk-go-v2/service/eks"
	_ "github.com/spf13/pflag"
	_ "github.com/spf13/viper"
	_ "go.uber.org/zap"
	_ "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/rest"
	_ "sigs.k8s.io/aws-iam-authenticator/pkg/token"
)
