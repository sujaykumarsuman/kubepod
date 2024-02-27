package kubepod

import (
	"context"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetPods returns a list of pods from the cluster in a namespace using the kubepod client
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

// GetPod returns a pod from the cluster in a namespace using the kubepod client
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
