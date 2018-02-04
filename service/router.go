package service

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string           `json:"name"`
	Method      string           `json:"method"`
	Pattern     string           `json:"pattern"`
	HandlerFunc http.HandlerFunc `json:"handler"`
}

type Routes []Route

type Router struct {
	service   *Service
	muxRouter *mux.Router
}

func NewRouter(service *Service) *Router {
	r := &Router{
		service:   service,
		muxRouter: mux.NewRouter().StrictSlash(true),
	}
	return r
}

func (r *Router) InitRouter() {
	var routes = Routes{
		Route{
			Name:        "HandleGet",
			Method:      "GET",
			Pattern:     "/rastore/{key}",
			HandlerFunc: r.service.HandleGet,
		},
		Route{
			Name:        "HandleSet",
			Method:      "POST",
			Pattern:     "/rastore",
			HandlerFunc: r.service.HandleSet,
		},
		Route{
			Name:        "HandleDelete",
			Method:      "DELETE",
			Pattern:     "/rastore/{key}",
			HandlerFunc: r.service.HandleDelete,
		},
	}

	for _, route := range routes {
		// var handler http.Handler
		handler := route.HandlerFunc
		// handler = Logger(handler, route.Name)
		r.muxRouter.Methods(route.Method).Path(route.Pattern).Name(route.Name).Handler(handler)
	}
}
