package router

import (
	"fmt"
	"strings"
)

const (
	sep              = "/"
	pathParamSepChar = ':'
	sepChar          = '/'
)

type node struct {
	path     string
	handle   HandlerFuncWithParam
	children []*node
}

func (n *node) addRoute(path string, handler HandlerFuncWithParam) {
	var parts []string
	var child *node

	if path[0] != sepChar {
		panic(fmt.Sprintf("Invalid Path: %s. Path must begin with '/'\n", path))
	}

	if isIndex(path) {
		parts = []string{path}
	} else {
		path = path[1:] // remove leading slash
		parts = strings.Split(path, sep)
	}

	for i := range parts {
		child = addPath(n, child, parts[i], i)
	}

	child.handle = handler
}

// Index vars
//		paramSize: path param current size
//		nextSepIndex : Index Byte of next / found in path

func (n *node) findRoute(path string) (HandlerFuncWithParam, []Param) {
	var (
		child      *node
		p          Param
		params     []Param
		paramsSize int
	)

	// handle Index
	if isIndex(path) {
		child, p = findPath(n, child, path, true)
	} else {
		// Prepare for path param parsing
		isRoot := true
		parts := strings.Count(path, sep)
		path := path[1:] // remove leading slash
		nextSepIndex := 0

		for nextSepIndex >= 0 {
			var part string

			// next occurrence of /. This reduces 6900 ns/op
			index := -1
			for i := 0; i < len(path); i++ {
				if path[i] == sepChar {
					index = i
					break
				}
			}
			nextSepIndex = index

			if nextSepIndex < 0 {
				part = path // remaining path is the part
			} else {
				part = path[:nextSepIndex]   // split part until next /
				path = path[nextSepIndex+1:] // remaining part is the path
			}

			child, p = findPath(n, child, part, isRoot)

			// collect path params
			if len(p.Key) > 0 {

				// lazy initialization
				if params == nil {
					params = make([]Param, 0, parts)
				}

				paramsSize = len(params)
				params = params[:paramsSize+1] // extend array by need
				params[paramsSize] = p
				paramsSize++
			}

			isRoot = false
		}
	}

	if child == nil {
		return nil, nil
	}

	return child.handle, params[:paramsSize]
}

func findPath(root, child *node, part string, isRoot bool) (*node, Param) {
	var p Param

	if isRoot {
		child, p = findChild(root, part)
	} else {
		if child != nil {
			child, p = findChild(child, part)
		}
	}

	return child, p
}

func addPath(n, child *node, part string, i int) *node {
	if i == 0 {
		child = n.insertChild(part)
	} else {
		child = child.insertChild(part)
	}

	return child
}

func (n *node) insertChild(part string) (*node) {
	if existingChild, _ := findChild(n, part); existingChild != nil {
		return existingChild
	}

	child := &node{
		path: part,
	}

	n.children = append(n.children, child)

	return child
}

func findChild(n *node, path string) (*node, Param) {
	p := Param{}

	for _, child := range n.children {
		// actual match
		if child.path == path {
			return child, p
		}

		// path param match
		if child.path[0] == pathParamSepChar {
			p.Key = child.path[1:]
			p.Value = path

			return child, p
		}
	}

	return nil, p
}

func isIndex(path string) bool {
	return sep == path
}
