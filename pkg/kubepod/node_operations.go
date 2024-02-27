package kubepod

import (
	"context"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetNodes returns a list of nodes from the cluster using the kubepod client
func (k *Kubepod) GetNodes() (*v1.NodeList, error) {
	nodes, err := k.CoreV1().Nodes().List(context.Background(), metaV1.ListOptions{})
	if err != nil {
		k.logger.Debug("error listing nodes with kubepod",
			zap.Any("kubepod", k.Clientset),
			zap.Any("cluster", k.cluster))
		k.logger.Error("unable to list nodes", zap.Error(err))
		return nil, err
	}
	for _, node := range nodes.Items {
		k.logger.Info("node", zap.String("name", node.Name))
	}
	return nodes, nil
}

// GetNode returns a node from the cluster using the kubepod client
func (k *Kubepod) GetNode(name string) (*v1.Node, error) {
	node, err := k.CoreV1().Nodes().Get(context.Background(), name, metaV1.GetOptions{})
	if err != nil {
		k.logger.Debug("error getting node with kubepod",
			zap.Any("kubepod", k.Clientset),
			zap.Any("cluster", k.cluster))
		k.logger.Error("unable to get node", zap.Error(err))
		return nil, err
	}
	k.logger.Info("node", zap.String("name", node.Name))
	return node, nil
}

// CreateObject creates a kubernetes object in the cluster using the kubepod client
func (k *Kubepod) CreateObject(obj interface{}) error {

	return nil
}
