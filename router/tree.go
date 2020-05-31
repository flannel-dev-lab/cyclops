package KCRouter

import (
	"net/http"
	"net/url"
	"strings"
)

// node contains the struct that holds information related to single node
type node struct {
	// children contains the children for existing node
	children []*node
	// component contains the name of the individual node
	component string
	// isNamedParam tells if a param is a named parameter or not
	isNamedParam bool
	// methods contains a map to http method and a handler for it
	methods map[string]http.HandlerFunc
}

// addNode adds a path to existing tree by splitting the path on "/"
// each split is considered a node and handler is added to the last node
func (n *node) addNode(method, path string, handler http.HandlerFunc) {
	components := strings.Split(path, "/")[1:]
	count := len(components)

	for count != 0 {
		aNode, component := n.traverse(components, nil)

		// Means node already exists, so we are updating one
		if aNode.component == component && count == 1 {
			// TODO add stub to create a panic when registration happens in the type
			// TODO /foo/bar/:hello & /foo/bar/boom
			aNode.methods[method] = handler
			return
		}

		newNode := node{component: component, isNamedParam: false, methods: make(map[string]http.HandlerFunc)}

		// Check if it is a named param
		if len(component) > 0 && component[0] == ':' {
			newNode.isNamedParam = true
		}

		// this is the last component of the url resource, so it gets the handler
		if count == 1 {
			newNode.methods[method] = handler
		}

		// Adds child to the current node
		aNode.children = append(aNode.children, &newNode)
		count--
	}
}

// traverse, traverses through the root node from the path components and returns
// the node and the component name
func (n *node) traverse(components []string, params url.Values) (*node, string) {
	component := components[0]
	if len(n.children) > 0 { // no children, then bail out.
		for _, child := range n.children {
			if component == child.component || child.isNamedParam {
				if child.isNamedParam && params != nil {
					params.Add(child.component[1:], component)
				}
				next := components[1:]
				if len(next) > 0 {
					return child.traverse(next, params)
				} else {
					return child, component
				}
			}
		}
	}

	return n, component
}
