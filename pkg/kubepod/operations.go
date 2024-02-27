package kubepod

import (
	"context"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

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

func (k *Kubepod) GetPods(namespace string) (*v1.PodList, error) {
	pods, err := k.CoreV1().Pods(namespace).List(context.Background(), metaV1.ListOptions{})
	if err != nil {
		k.logger.Debug("error listing pods with kubepod",
			zap.Any("kubepod", k.Clientset),
			zap.Any("cluster", k.cluster))
		k.logger.Error("unable to list pods", zap.Error(err))
		return nil, err
	}
	for _, pod := range pods.Items {
		k.logger.Info("pod", zap.String("name", pod.Name))
	}
	return pods, nil
}

func (k *Kubepod) GetPod(name, namespace string) (*v1.Pod, error) {
	pod, err := k.CoreV1().Pods(namespace).Get(context.Background(), name, metaV1.GetOptions{})
	if err != nil {
		k.logger.Debug("error getting pod with kubepod",
			zap.Any("kubepod", k.Clientset),
			zap.Any("cluster", k.cluster))
		k.logger.Error("unable to get pod", zap.Error(err))
		return nil, err
	}
	k.logger.Info("pod", zap.String("name", pod.Name))
	return pod, nil
}

func (k *Kubepod) CreateObject(obj interface{}) error {

	return nil
}
