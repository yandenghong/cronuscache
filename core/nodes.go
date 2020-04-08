package core

// NodeSelector is the interface that must be implemented to locate
// the node that owns a specific key.
type NodeSelector interface {
	SelectNode(key string) (node NodeGetter, ok bool)
}

// NodeGetter is the interface that must be implemented by a node.
type NodeGetter interface {
	Get(group string, key string) ([]byte, error)
}
