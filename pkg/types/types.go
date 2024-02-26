package types

import (
	"github.com/go-chi/render"
	v1 "k8s.io/api/core/v1"
	"net/http"
)

type GetNodeRequest struct {
	CallerId string `json:"callerId"`
	Name     string `json:"name"`
}

func (g *GetNodeRequest) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
func (g *GetNodeRequest) Bind(r *http.Request) error {
	return nil
}

type GetNodeResponse struct {
	CallerID    string   `json:"callerId"`
	ClusterName string   `json:"clusterName"`
	Node        *v1.Node `json:"node"`
}

func (g *GetNodeResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
func (g *GetNodeResponse) Bind(r *http.Request) error {
	return nil
}

type GetNodesRequest struct {
	CallerId string `json:"callerId"`
}

func (g *GetNodesRequest) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
func (g *GetNodesRequest) Bind(r *http.Request) error {
	return nil
}

type GetNodesResponse struct {
	CallerID    string       `json:"callerId"`
	ClusterName string       `json:"clusterName"`
	Nodes       *v1.NodeList `json:"nodes"`
}

func (g *GetNodesResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
func (g *GetNodesResponse) Bind(r *http.Request) error {
	return nil
}

type GetPodRequest struct {
	CallerId  string `json:"callerId"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

func (g *GetPodRequest) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
func (g *GetPodRequest) Bind(r *http.Request) error {
	return nil
}

type GetPodResponse struct {
	CallerID    string  `json:"callerId"`
	ClusterName string  `json:"clusterName"`
	Pod         *v1.Pod `json:"pod"`
}

func (g *GetPodResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
func (g *GetPodResponse) Bind(r *http.Request) error {
	return nil
}

type GetPodsRequest struct {
	CallerId  string `json:"callerId"`
	Namespace string `json:"namespace"`
}

func (g *GetPodsRequest) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
func (g *GetPodsRequest) Bind(r *http.Request) error {
	return nil
}

type GetPodsResponse struct {
	CallerID    string      `json:"callerId"`
	ClusterName string      `json:"clusterName"`
	Pods        *v1.PodList `json:"pods"`
}

func (g *GetPodsResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
func (g *GetPodsResponse) Bind(r *http.Request) error {
	return nil
}

// ErrResponse renderer type for handling all sorts of errors.
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status" example:"Resource not found."`                                         // user-level status message
	AppCode    int64  `json:"code,omitempty" example:"404"`                                                 // application-specific error code
	ErrorText  string `json:"error,omitempty" example:"The requested resource was not found on the server"` // application-level error message, for debugging
} // @name ErrorResponse

// Render implements the github.com/go-chi/render.Renderer interface for ErrResponse
func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// ErrInvalidRequest returns a structured http response for invalid requests
func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

// ErrRender returns a structured http response in case of rendering errors
func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusUnprocessableEntity,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}

// ErrInternalServer returns a structured http response for internal server errors
func ErrInternalServer(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusInternalServerError,
		StatusText:     "Internal server error.",
		ErrorText:      err.Error(),
	}
}
