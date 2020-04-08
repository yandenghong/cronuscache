package core

import (
	"log"
	"fmt"
	"sync"
)

type Getter interface {
	Get(key string) ([]byte,error)
}

type GetterFunc func(key string) ([]byte,error)

func (f GetterFunc) Get(key string) ([]byte,error) {
	return f(key)
}

// A Group is a cache namespace and associated data loaded spread over
type Group struct {
	name      string
	getter    Getter
	mainCache cache
	nodes     NodeSelector
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

// NewGroup create a new instance of Group
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = g
	return g
}

// GetGroup returns the named group previously created with NewGroup, or
// nil if there's no such group.
func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	g := groups[name]
	return g
}

// Get value for a key from cache
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	if v, ok := g.mainCache.get(key); ok {
		log.Println("[CronusCache] hit")
		return v, nil
	}

	return g.load(key)
}

func (g *Group) load(key string) (value ByteView, err error) {
	if g.nodes != nil {
		if node, ok := g.nodes.SelectNode(key); ok {
			if value, err = g.getFromNode(node, key); err == nil {
				return value, nil
			}
			log.Println("[CronusCache] Failed to get from node", err)
		}
	}
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err

	}
	value := ByteView{b: copyBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}

// RegisterNodes registers a NodeSelector for choosing remote node
func (g *Group) RegisterNodes(nodes NodeSelector) {
	if g.nodes != nil {
		panic("RegisterNodes called more than once")
	}
	g.nodes = nodes
}

func (g *Group) getFromNode(node NodeGetter, key string) (ByteView, error) {
	bytes, err := node.Get(g.name, key)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: bytes}, nil
}