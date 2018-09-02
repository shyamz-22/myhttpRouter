package router

import (
	"net/http"
)

type HandlerFuncWithParam func(w http.ResponseWriter, request *http.Request, param PathParams)

type Router struct {
	routes map[string]*node
}

func New() *Router {
	return &Router{}
}

func (rtr *Router) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	path := request.URL.Path
	method := request.Method

	routes := rtr.routes[method]

	if routes == nil {
		handleError(rtr, w, path, method)
		return
	}

	handle, params := routes.findRoute(path)

	if handle == nil {
		handleError(rtr, w, path, method)
		return
	}

	handle(w, request, params)
}

func (rtr *Router) Add(path string, method string, handler HandlerFuncWithParam) {
	if rtr.routes == nil {
		rtr.routes = make(map[string]*node)
	}

	root := rtr.routes[method]

	if root == nil {
		root = new(node)
		rtr.routes[method] = root
	}

	root.addRoute(path, handler)
}

func handleError(router *Router, writer http.ResponseWriter, path, requestMethod string) {
	status := http.StatusNotFound

search:
	for method := range router.routes {
		// skip search as we know request method is already searched by normal flow
		if method == requestMethod {
			continue search
		}

		root := router.routes[method]

		handle, _ := root.findRoute(path)

		if handle != nil {
			status = http.StatusMethodNotAllowed
			break search
		}
	}

	writer.WriteHeader(status)
}
