package api

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/sujaykumarsuman/kubepod/pkg/kubepod"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
	"net/http"
)

var (
	kpod   kubepod.Interface
	logger *zap.Logger
)

func GetRouter(log *zap.Logger, k8sClient kubepod.Interface) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	kpod = k8sClient
	logger = log
	buildTree(r)
	return r
}

func buildTree(r *chi.Mux) {
	r.HandleFunc("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, r.RequestURI+"/", http.StatusMovedPermanently)
	})
	r.Get("/swagger*", httpSwagger.Handler())

	r.Route("/nodes", func(r chi.Router) {
		r.Get("/", GetNodes)
		r.Get("/{node}", GetNode)
	})

	r.Route("/pods", func(r chi.Router) {
		r.Get("/", GetPods)
		r.Get("/{pod}", GetPod)
	})
}
