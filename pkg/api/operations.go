package api

import (
	"github.com/go-chi/render"
	"github.com/spf13/viper"
	"github.com/sujaykumarsuman/kubepod/pkg/types"
	"go.uber.org/zap"
	"net/http"
)

func GetNodes(w http.ResponseWriter, r *http.Request) {
	req := &types.GetNodesRequest{}
	if err := render.Bind(r, req); err != nil {
		_ = render.Render(w, r, types.ErrInvalidRequest(err))
		return
	}
	nodes, err := kpod.GetNodes()
	if err != nil {
		logger.Debug("unable to get nodes", zap.Error(err))
		_ = render.Render(w, r, types.ErrInternalServer(err))
		return
	}

	response := &types.GetNodesResponse{
		ClusterName: viper.GetString("eks.cluster.name"),
		Nodes:       nodes,
		CallerID:    req.CallerId,
	}
	if err := render.Render(w, r, response); err != nil {
		logger.Debug("unable to render response", zap.Error(err))
		_ = render.Render(w, r, types.ErrRender(err))
		return
	}
}

func GetNode(w http.ResponseWriter, r *http.Request) {
	req := &types.GetNodeRequest{}
	if err := render.Bind(r, req); err != nil {
		_ = render.Render(w, r, types.ErrInvalidRequest(err))
		return
	}
	node, err := kpod.GetNode(req.Name)
	if err != nil {
		logger.Debug("unable to get node", zap.Error(err))
		_ = render.Render(w, r, types.ErrInternalServer(err))
		return
	}

	response := &types.GetNodeResponse{
		ClusterName: viper.GetString("eks.cluster.name"),
		Node:        node,
		CallerID:    req.CallerId,
	}
	if err := render.Render(w, r, response); err != nil {
		logger.Debug("unable to render response", zap.Error(err))
		_ = render.Render(w, r, types.ErrRender(err))
		return
	}
}

func GetPods(w http.ResponseWriter, r *http.Request) {
	req := &types.GetPodsRequest{}
	if err := render.Bind(r, req); err != nil {
		_ = render.Render(w, r, types.ErrInvalidRequest(err))
		return
	}
	pods, err := kpod.GetPods(req.Namespace)
	if err != nil {
		logger.Debug("unable to get pods", zap.Error(err))
		_ = render.Render(w, r, types.ErrInternalServer(err))
		return
	}

	response := &types.GetPodsResponse{
		ClusterName: viper.GetString("eks.cluster.name"),
		Pods:        pods,
		CallerID:    req.CallerId,
	}
	if err := render.Render(w, r, response); err != nil {
		logger.Debug("unable to render response", zap.Error(err))
		_ = render.Render(w, r, types.ErrRender(err))
		return
	}
}

func GetPod(w http.ResponseWriter, r *http.Request) {
	req := &types.GetPodRequest{}
	if err := render.Bind(r, req); err != nil {
		_ = render.Render(w, r, types.ErrInvalidRequest(err))
		return
	}
	pod, err := kpod.GetPod(req.Name, req.Namespace)
	if err != nil {
		logger.Debug("unable to get pod", zap.Error(err))
		_ = render.Render(w, r, types.ErrInternalServer(err))
		return
	}

	response := &types.GetPodResponse{
		ClusterName: viper.GetString("eks.cluster.name"),
		Pod:         pod,
		CallerID:    req.CallerId,
	}
	if err := render.Render(w, r, response); err != nil {
		logger.Debug("unable to render response", zap.Error(err))
		_ = render.Render(w, r, types.ErrRender(err))
		return
	}
}
