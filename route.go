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

// AddGet registers a new request handle with the given path and Get-method.
func (rtr *Router) AddGet(path string, handler HandlerFuncWithParam) {
	rtr.Add(path, http.MethodGet, handler)

}

// AddPost registers a new request handle with the given path and Post-method.
func (rtr *Router) AddPost(path string, handler HandlerFuncWithParam) {
	rtr.Add(path, http.MethodPost, handler)
}

// AddPut registers a new request handle with the given path and Put-method.
func (rtr *Router) AddPut(path string, handler HandlerFuncWithParam) {
	rtr.Add(path, http.MethodPut, handler)
}

// AddDelete registers a new request handle with the given path and Delete-method.
func (rtr *Router) AddDelete(path string, handler HandlerFuncWithParam) {
	rtr.Add(path, http.MethodDelete, handler)
}

// AddOptions registers a new request handle with the given path and Options-method.
func (rtr *Router) AddOptions(path string, handler HandlerFuncWithParam) {
	rtr.Add(path, http.MethodOptions, handler)
}

// AddPatch registers a new request handle with the given path and Patch-method.
func (rtr *Router) AddPatch(path string, handler HandlerFuncWithParam) {
	rtr.Add(path, http.MethodPatch, handler)
}

// AddHead registers a new request handle with the given path and Head-method.
func (rtr *Router) AddHead(path string, handler HandlerFuncWithParam) {
	rtr.Add(path, http.MethodHead, handler)
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
