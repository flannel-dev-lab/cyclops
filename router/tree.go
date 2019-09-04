package router

import (
	"strings"
)

type node struct {
	children     []*node
	component    string
	isNamedParam bool
	methods      map[string]Handle
}

func (n *node) addNode(method, path string, handler Handle) {
	if path[0] != '/' {
		panic("Path has to start with a /.")
	}

	components := strings.Split(path, "/")[1:]

	for componentCount := len(components); componentCount >= 1; componentCount-- {
		aNode, component := n.searchTree(components, nil)

		if aNode.component == component && componentCount == 1 { // update an existing node.
			aNode.methods[method] = handler
			return
		}

		newNode := node{component: component, isNamedParam: false, methods: make(map[string]Handle)}

		if len(component) > 0 && component[0] == ':' { // check if it is a named param.
			newNode.isNamedParam = true
		}

		if componentCount == 1 {
			newNode.methods[method] = handler
		}

		aNode.children = append(aNode.children, &newNode)
	}

}

func (n *node) searchTree(components []string, params map[string]string) (*node, string) {
	component := components[0]

	if len(n.children) > 0 { // no children, then bail out.
		for _, child := range n.children {
			if component == child.component || child.isNamedParam {
				if child.isNamedParam && params != nil {
					params[child.component[1:]] = component
				}
				next := components[1:]
				if len(next) > 0 {
					return child.searchTree(next, params)
				} else {
					return child, component
				}
			}
		}
	}

	return n, component
}